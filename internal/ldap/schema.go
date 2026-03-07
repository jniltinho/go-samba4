package ldap

// User represents a mapped Active Directory User object
type User struct {
	DN                 string
	SAMAccountName     string
	UserPrincipalName  string
	DisplayName        string
	GivenName          string
	SN                 string
	Mail               string
	TelephoneNumber    string
	Title              string
	Department         string
	UserAccountControl int
	MemberOf           []string
}

// Group represents a mapped Active Directory Group object
type Group struct {
	DN             string
	SAMAccountName string
	Description    string
	GroupType      int
	Member         []string
}

// OU represents an Organizational Unit in AD
type OU struct {
	DN          string
	Name        string
	Description string
}

// Constants for common AD Attributes
const (
	AttrSAMAccountName     = "sAMAccountName"
	AttrUserPrincipalName  = "userPrincipalName"
	AttrDisplayName        = "displayName"
	AttrGivenName          = "givenName"
	AttrSN                 = "sn"
	AttrMail               = "mail"
	AttrTelephoneNumber    = "telephoneNumber"
	AttrTitle              = "title"
	AttrDepartment         = "department"
	AttrUserAccountControl = "userAccountControl"
	AttrMemberOf           = "memberOf"
	AttrDescription        = "description"
	AttrGroupType          = "groupType"
	AttrMember             = "member"
	AttrObjectClass        = "objectClass"
	AttrOU                 = "ou"
)
