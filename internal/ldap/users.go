package ldap

import (
	"encoding/binary"
	"fmt"
	"unicode/utf16"

	goldap "github.com/go-ldap/ldap/v3"
)

// GetAllUsers attributes from LDAP
func (c *Client) GetAllUsers(filter string) ([]User, error) {
	if filter == "" {
		filter = "(objectClass=user)"
	}

	searchRequest := goldap.NewSearchRequest(
		c.config.BaseDN,
		goldap.ScopeWholeSubtree, goldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(&(objectCategory=person)%s)", filter),
		[]string{
			AttrSAMAccountName, AttrUserPrincipalName, AttrDisplayName,
			AttrGivenName, AttrSN, AttrMail, AttrTelephoneNumber, AttrTitle,
			AttrDepartment, AttrUserAccountControl, AttrMemberOf,
		},
		nil,
	)

	sr, err := c.Search(searchRequest)
	if err != nil {
		return nil, fmt.Errorf("error searching for users: %w", err)
	}

	var users []User
	for _, entry := range sr.Entries {
		userAccountControl := 0
		fmt.Sscanf(entry.GetAttributeValue(AttrUserAccountControl), "%d", &userAccountControl)

		users = append(users, User{
			DN:                 entry.DN,
			SAMAccountName:     entry.GetAttributeValue(AttrSAMAccountName),
			UserPrincipalName:  entry.GetAttributeValue(AttrUserPrincipalName),
			DisplayName:        entry.GetAttributeValue(AttrDisplayName),
			GivenName:          entry.GetAttributeValue(AttrGivenName),
			SN:                 entry.GetAttributeValue(AttrSN),
			Mail:               entry.GetAttributeValue(AttrMail),
			TelephoneNumber:    entry.GetAttributeValue(AttrTelephoneNumber),
			Title:              entry.GetAttributeValue(AttrTitle),
			Department:         entry.GetAttributeValue(AttrDepartment),
			UserAccountControl: userAccountControl,
			MemberOf:           entry.GetAttributeValues(AttrMemberOf),
		})
	}

	return users, nil
}

// GetUserBySAM fetches a single user by sAMAccountName.
func (c *Client) GetUserBySAM(sam string) (*User, error) {
	safeSAM := goldap.EscapeFilter(sam)
	searchRequest := goldap.NewSearchRequest(
		c.config.BaseDN,
		goldap.ScopeWholeSubtree, goldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(&(objectCategory=person)(objectClass=user)(sAMAccountName=%s))", safeSAM),
		[]string{
			AttrSAMAccountName, AttrUserPrincipalName, AttrDisplayName,
			AttrGivenName, AttrSN, AttrMail, AttrTelephoneNumber, AttrTitle,
			AttrDepartment, AttrUserAccountControl, AttrMemberOf,
		},
		nil,
	)

	sr, err := c.Search(searchRequest)
	if err != nil {
		return nil, fmt.Errorf("error searching for user %q: %w", sam, err)
	}

	if len(sr.Entries) == 0 {
		return nil, fmt.Errorf("user %q not found", sam)
	}

	entry := sr.Entries[0]
	uac := 0
	fmt.Sscanf(entry.GetAttributeValue(AttrUserAccountControl), "%d", &uac)

	user := &User{
		DN:                 entry.DN,
		SAMAccountName:     entry.GetAttributeValue(AttrSAMAccountName),
		UserPrincipalName:  entry.GetAttributeValue(AttrUserPrincipalName),
		DisplayName:        entry.GetAttributeValue(AttrDisplayName),
		GivenName:          entry.GetAttributeValue(AttrGivenName),
		SN:                 entry.GetAttributeValue(AttrSN),
		Mail:               entry.GetAttributeValue(AttrMail),
		TelephoneNumber:    entry.GetAttributeValue(AttrTelephoneNumber),
		Title:              entry.GetAttributeValue(AttrTitle),
		Department:         entry.GetAttributeValue(AttrDepartment),
		UserAccountControl: uac,
		MemberOf:           entry.GetAttributeValues(AttrMemberOf),
	}

	return user, nil
}

// CreateUser creates a new AD user in the given OU DN.
// Steps: 1) Add disabled account, 2) Set password, 3) Enable account.
func (c *Client) CreateUser(u User, password, ouDN string) error {
	dn := fmt.Sprintf("CN=%s,%s", goldap.EscapeFilter(u.DisplayName), ouDN)

	addRequest := goldap.NewAddRequest(dn, nil)
	addRequest.Attribute("objectClass", []string{"top", "person", "organizationalPerson", "user"})
	addRequest.Attribute(AttrSAMAccountName, []string{u.SAMAccountName})
	addRequest.Attribute(AttrUserAccountControl, []string{"514"}) // disabled initially

	if u.DisplayName != "" {
		addRequest.Attribute(AttrDisplayName, []string{u.DisplayName})
	}
	if u.GivenName != "" {
		addRequest.Attribute(AttrGivenName, []string{u.GivenName})
	}
	if u.SN != "" {
		addRequest.Attribute(AttrSN, []string{u.SN})
	}
	if u.UserPrincipalName != "" {
		addRequest.Attribute(AttrUserPrincipalName, []string{u.UserPrincipalName})
	}
	if u.Mail != "" {
		addRequest.Attribute(AttrMail, []string{u.Mail})
	}
	if u.TelephoneNumber != "" {
		addRequest.Attribute(AttrTelephoneNumber, []string{u.TelephoneNumber})
	}
	if u.Title != "" {
		addRequest.Attribute(AttrTitle, []string{u.Title})
	}
	if u.Department != "" {
		addRequest.Attribute(AttrDepartment, []string{u.Department})
	}

	if err := c.conn.Add(addRequest); err != nil {
		return fmt.Errorf("failed to create user %q: %w", u.SAMAccountName, err)
	}

	// Set password (requires LDAPS)
	if password != "" {
		if err := c.SetPassword(dn, password); err != nil {
			// Attempt to clean up the created object on password failure
			_ = c.conn.Del(goldap.NewDelRequest(dn, nil))
			return fmt.Errorf("failed to set password for %q: %w", u.SAMAccountName, err)
		}

		// Set final UAC (use provided value, default to 512=enabled)
		finalUAC := 512
		if u.UserAccountControl > 0 {
			finalUAC = u.UserAccountControl
		}
		modEnable := goldap.NewModifyRequest(dn, nil)
		modEnable.Replace(AttrUserAccountControl, []string{fmt.Sprintf("%d", finalUAC)})
		if err := c.conn.Modify(modEnable); err != nil {
			return fmt.Errorf("failed to set account control %q: %w", u.SAMAccountName, err)
		}
	}

	return nil
}

// UpdateUser modifies an existing user's attributes by DN.
// When u.UserAccountControl > 0, it also updates the UAC (enable/disable).
func (c *Client) UpdateUser(dn string, u User) error {
	modRequest := goldap.NewModifyRequest(dn, nil)

	replaceIfNotEmpty := func(attr, val string) {
		if val != "" {
			modRequest.Replace(attr, []string{val})
		}
	}

	replaceIfNotEmpty(AttrDisplayName, u.DisplayName)
	replaceIfNotEmpty(AttrGivenName, u.GivenName)
	replaceIfNotEmpty(AttrSN, u.SN)
	replaceIfNotEmpty(AttrMail, u.Mail)
	replaceIfNotEmpty(AttrTelephoneNumber, u.TelephoneNumber)
	replaceIfNotEmpty(AttrTitle, u.Title)
	replaceIfNotEmpty(AttrDepartment, u.Department)
	replaceIfNotEmpty(AttrUserPrincipalName, u.UserPrincipalName)

	// UAC > 0 means the caller explicitly wants to set account control flags
	if u.UserAccountControl > 0 {
		modRequest.Replace(AttrUserAccountControl, []string{fmt.Sprintf("%d", u.UserAccountControl)})
	}

	if err := c.conn.Modify(modRequest); err != nil {
		return fmt.Errorf("failed to update user %q: %w", dn, err)
	}

	return nil
}

// DeleteUser removes an AD object by its DN.
func (c *Client) DeleteUser(dn string) error {
	delRequest := goldap.NewDelRequest(dn, nil)
	if err := c.conn.Del(delRequest); err != nil {
		return fmt.Errorf("failed to delete user %q: %w", dn, err)
	}
	return nil
}

// SetPassword sets a user's password via unicodePwd (requires LDAPS or StartTLS).
// The password is encoded as UTF-16LE surrounded by double quotes, as required by AD.
func (c *Client) SetPassword(dn, newPassword string) error {
	encoded, err := encodePassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to encode password: %w", err)
	}

	modRequest := goldap.NewModifyRequest(dn, nil)
	modRequest.Replace("unicodePwd", []string{string(encoded)})

	if err := c.conn.Modify(modRequest); err != nil {
		return fmt.Errorf("failed to set password: %w", err)
	}

	return nil
}

// encodePassword returns the AD-compatible unicodePwd value:
// UTF-16LE encoding of the password surrounded by double quotes.
func encodePassword(password string) ([]byte, error) {
	quoted := `"` + password + `"`
	runes := utf16.Encode([]rune(quoted))
	buf := make([]byte, len(runes)*2)
	for i, r := range runes {
		binary.LittleEndian.PutUint16(buf[i*2:], r)
	}
	return buf, nil
}
