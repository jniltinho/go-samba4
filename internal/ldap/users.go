package ldap

import (
	"fmt"

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
