package v1

import (
	"github.com/volvo-cars/connect-access-control/internal/pkg/authz"
	"github.com/volvo-cars/connect-access-control/internal/pkg/store"
)

func toClient(client store.Client) Client {
	return Client{
		ID:                 client.ID,
		Name:               client.Name,
		Description:        client.Description,
		WhitelistedDomains: client.WhitelistedDomains,
		DependantScopes:    client.DependantScopes,
	}
}

func toClients(clients []store.Client) []Client {
	arr := make([]Client, len(clients))
	for i, client := range clients {
		arr[i] = toClient(client)
	}

	return arr
}

func toRole(role store.Role) Role {
	return Role{
		ID:          role.ID,
		Name:        role.Name,
		Description: role.Description,
	}
}

func toRoles(roles []store.Role) []Role {
	arr := make([]Role, len(roles))
	for i, role := range roles {
		arr[i] = toRole(role)
	}

	return arr
}

func toScope(scope store.Scope) Scope {
	permissionGroups := toPermissionGroups(scope.PermissionGroups)
	return Scope{
		Key:              scope.Key,
		Label:            scope.Label,
		Description:      scope.Description,
		PermissionGroups: permissionGroups,
	}
}

func toScopes(scopes []store.Scope) []Scope {
	arr := make([]Scope, len(scopes))
	for i, scope := range scopes {
		arr[i] = toScope(scope)
	}

	return arr
}

func toPermissionGroup(group store.PermissionGroup) PermissionGroup {
	return PermissionGroup{
		Key:         group.Key,
		Label:       group.Label,
		Description: group.Description,
	}
}

func toPermissionGroups(groups []store.PermissionGroup) []PermissionGroup {
	arr := make([]PermissionGroup, len(groups))
	for i, group := range groups {
		arr[i] = toPermissionGroup(group)
	}

	return arr
}

func toRoleMapping(mapping store.RoleMapping) RoleMapping {
	arr := make([]Mapping, len(mapping.Mapping))
	for i, m := range mapping.Mapping {
		arr[i] = Mapping{
			Filter: Filter{
				Market:      m.Filter.Market,
				UserType:    m.Filter.UserType,
				PartnerType: m.Filter.PartnerType,
			},
			PermissionGroups: m.PermissionGroups,
		}
	}

	return RoleMapping{
		RoleID:  mapping.RoleID,
		Mapping: arr,
	}
}

func toRoleMappings(mappings []store.RoleMapping) []RoleMapping {
	arr := make([]RoleMapping, len(mappings))
	for i, mapping := range mappings {
		arr[i] = toRoleMapping(mapping)
	}

	return arr
}

func toUser(user authz.User) User {
	partners := toPartners(user.Partners)
	return User{
		ID:          user.ID,
		Email:       user.Email,
		CDSID:       user.CDSID,
		CountryCode: user.CountryCode,
		Partners:    partners,
	}
}

func toPartner(partner authz.Partner) Partner {
	return Partner{
		ID:               partner.ID,
		RoleCode:         partner.RoleCode,
		Name:             partner.Name,
		Type:             partner.Type,
		DistributorID:    partner.DistributorID,
		ParmaPartnerCode: partner.ParmaPartnerCode,
		Market:           partner.Market,
		Active:           partner.Active,
		Roles:            partner.Roles,
	}
}

func toPartners(partners []authz.Partner) []Partner {
	arr := make([]Partner, len(partners))
	for i, partner := range partners {
		arr[i] = toPartner(partner)
	}

	return arr
}

func toUserAccess(access authz.UserAccess) UserAccess {
	return UserAccess{
		Context: Context{
			ID:   access.Context.ID,
			Type: access.Context.Type,
			Tag:  access.Context.Tag,
		},
		Roles:            access.Roles,
		PermissionGroups: access.PermissionGroups,
	}
}

func toUserAccesses(accesses []authz.UserAccess) []UserAccess {
	arr := make([]UserAccess, len(accesses))
	for i, access := range accesses {
		arr[i] = toUserAccess(access)
	}

	return arr
}
