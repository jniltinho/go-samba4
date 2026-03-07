package ldap

import (
	"fmt"

	goldap "github.com/go-ldap/ldap/v3"
)

// GetAllOUs queries Organizational Units
func (c *Client) GetAllOUs() ([]OU, error) {
	searchRequest := goldap.NewSearchRequest(
		c.config.BaseDN,
		goldap.ScopeWholeSubtree, goldap.NeverDerefAliases, 0, 0, false,
		"(objectClass=organizationalUnit)",
		[]string{
			AttrOU, AttrDescription,
		},
		nil,
	)

	sr, err := c.Search(searchRequest)
	if err != nil {
		return nil, fmt.Errorf("error searching for OUs: %w", err)
	}

	var ous []OU
	for _, entry := range sr.Entries {
		ous = append(ous, OU{
			DN:          entry.DN,
			Name:        entry.GetAttributeValue(AttrOU),
			Description: entry.GetAttributeValue(AttrDescription),
		})
	}

	return ous, nil
}
