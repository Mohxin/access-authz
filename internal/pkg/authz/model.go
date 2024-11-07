package authz

type UserAccess struct {
	Context          Context
	Roles            []string
	PermissionGroups map[string][]string
}

type Context struct {
	ID   string
	Type string
	Tag  string
}

type User struct {
	ID          string
	Email       string
	CDSID       string
	CountryCode string
	Partners    []Partner
}

type Partner struct {
	ID               string
	RoleCode         string
	Name             string
	Type             string
	DistributorID    string
	ParmaPartnerCode string
	Market           string
	IsPrimary        bool
	Active           bool
	Roles            []string
}
