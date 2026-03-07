package ldap

import (
	"fmt"

	goldap "github.com/go-ldap/ldap/v3"
)

// GetAllGroups queries groups from AD
func (c *Client) GetAllGroups(filter string) ([]Group, error) {
	if filter == "" {
		filter = "(objectClass=group)"
	}

	searchRequest := goldap.NewSearchRequest(
		c.config.BaseDN,
		goldap.ScopeWholeSubtree, goldap.NeverDerefAliases, 0, 0, false,
		filter,
		[]string{
			AttrSAMAccountName, AttrDescription, AttrGroupType, AttrMember,
		},
		nil,
	)

	sr, err := c.Search(searchRequest)
	if err != nil {
		return nil, fmt.Errorf("error searching for groups: %w", err)
	}

	var groups []Group
	for _, entry := range sr.Entries {
		groupType := 0
		fmt.Sscanf(entry.GetAttributeValue(AttrGroupType), "%d", &groupType)

		groups = append(groups, Group{
			DN:             entry.DN,
			SAMAccountName: entry.GetAttributeValue(AttrSAMAccountName),
			Description:    entry.GetAttributeValue(AttrDescription),
			GroupType:      groupType,
			Member:         entry.GetAttributeValues(AttrMember),
		})
	}

	return groups, nil
}
