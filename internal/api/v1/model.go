package v1

type Client struct {
	ID                 string   `json:"id,omitempty"`
	Name               string   `json:"name,omitempty"`
	Description        string   `json:"description,omitempty"`
	WhitelistedDomains []string `json:"whitelisted_domains,omitempty"`
	DependantScopes    []string `json:"dependant_scopes,omitempty"`
} // @name Client

type Role struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
} // @name Role

type Scope struct {
	Key              string            `json:"key,omitempty"`
	Label            string            `json:"label,omitempty"`
	Description      string            `json:"description,omitempty"`
	PermissionGroups []PermissionGroup `json:"permission_groups,omitempty"`
} // @name Scope

type PermissionGroup struct {
	Key         string `json:"key,omitempty"`
	Label       string `json:"label,omitempty"`
	Description string `json:"description,omitempty"`
} // @name PermissionGroup

type RoleMapping struct {
	RoleID  string    `json:"id,omitempty"`
	Mapping []Mapping `json:"mapping,omitempty"`
} // @name RoleMapping

type Mapping struct {
	Filter           Filter   `json:"filter,omitempty"`
	PermissionGroups []string `json:"permission_groups,omitempty"`
} // @name Mapping

type Filter struct {
	Market      []string `json:"market,omitempty"`
	UserType    []string `json:"user_type,omitempty"`
	PartnerType []string `json:"partner_type,omitempty"`
} // @name Filter

type User struct {
	ID          string    `json:"id,omitempty"`
	Email       string    `json:"email,omitempty"`
	CDSID       string    `json:"cdsid,omitempty"`
	CountryCode string    `json:"country_code,omitempty"`
	Partners    []Partner `json:"partners,omitempty"`
} // @name User

type Partner struct {
	ID               string   `json:"id,omitempty"`
	RoleCode         string   `json:"role_code,omitempty"`
	Name             string   `json:"name,omitempty"`
	Type             string   `json:"type,omitempty"`
	DistributorID    string   `json:"distributor_id,omitempty"`
	ParmaPartnerCode string   `json:"parma_partner_code,omitempty"`
	Market           string   `json:"market,omitempty"`
	Active           bool     `json:"active,omitempty"`
	Roles            []string `json:"roles,omitempty"`
} // @name Partner

type UserAccess struct {
	Context          Context             `json:"context,omitempty"`
	Roles            []string            `json:"roles,omitempty"`
	PermissionGroups map[string][]string `json:"permission_groups,omitempty"`
} // @name UserAccess

type Context struct {
	ID   string `json:"id,omitempty"`
	Type string `json:"type,omitempty"`
	Tag  string `json:"tag,omitempty"`
} // @name Context
