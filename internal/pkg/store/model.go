package store

type Market string

const (
	MarketSe Market = "SE"
	MarketNo Market = "NO"
	MarketDe Market = "DE"
	MarketGb Market = "GB"
	MarketUs Market = "US"
)

func (m Market) String() string {
	return string(m)
}

type UserType string

const (
	UserTypeInternal UserType = "INTERNAL"
	UserTypeExternal UserType = "EXTERNAL"
)

func (u UserType) String() string {
	return string(u)
}

type PartnerType string

const (
	PartnerTypeNsc   PartnerType = "NSC"
	PartnerTypeParma PartnerType = "PARMA"
)

func (u PartnerType) String() string {
	return string(u)
}

type ClientDefinition struct {
	Client Client `json:"client"`
}

type Client struct {
	ID                 string   `json:"id"`
	Name               string   `json:"name"`
	Description        string   `json:"description"`
	WhitelistedDomains []string `json:"whitelisted_domains"`
	DependantScopes    []string `json:"dependant_scopes"`
}

type RoleDefinition struct {
	Roles []Role `json:"roles"`
}

type Role struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ScopeDefinition struct {
	Scope Scope `json:"scope"`
}

type Scope struct {
	Key              string            `json:"key"`
	Label            string            `json:"label"`
	Description      string            `json:"description"`
	PermissionGroups []PermissionGroup `json:"-"`
}

type PermissionGroupDefinition struct {
	PermissionGroups []PermissionGroup `json:"permission_groups"`
}

type PermissionGroup struct {
	Key         string `json:"key"`
	Label       string `json:"label"`
	Description string `json:"description"`
}

type RoleMappingDefinition struct {
	RoleAssignment RoleMapping `json:"role"`
}

type RoleMapping struct {
	RoleID  string    `json:"id"`
	Mapping []Mapping `json:"mapping"`
}

type Mapping struct {
	Filter           Filter   `json:"filter"`
	PermissionGroups []string `json:"permission_groups"`
}

type Filter struct {
	Market      []string `json:"market"`
	UserType    []string `json:"user_type"`
	PartnerType []string `json:"partner_type"`
}

func ToMarketString(market []Market) []string {
	arr := make([]string, len(market))
	for i, m := range market {
		arr[i] = m.String()
	}

	return arr
}

func ToUserTypeString(userType []UserType) []string {
	arr := make([]string, len(userType))
	for i, u := range userType {
		arr[i] = u.String()
	}

	return arr
}
