package v1

type Response[T any] struct {
	Data T `json:"data"`
} // @name Response

type ErrorResponse struct {
	Error Error `json:"error"`
} // @name ErrorResponse

type Error struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
} // @name Error

type (
	ClientResponse       = Response[Client]        // @name ClientResponse
	ClientsResponse      = Response[[]Client]      // @name ClientsResponse
	RoleResponse         = Response[Role]          // @name RoleResponse
	RolesResponse        = Response[[]Role]        // @name RolesResponse
	ScopeResponse        = Response[Scope]         // @name ScopeResponse
	ScopesResponse       = Response[[]Scope]       // @name ScopesResponse
	RoleMappingResponse  = Response[RoleMapping]   // @name RoleMappingResponse
	RoleMappingsResponse = Response[[]RoleMapping] // @name RoleMappingsResponse
	UserResponse         = Response[User]          // @name UserResponse
	UserAccessResponse   = Response[UserAccess]    // @name UserAccessResponse
)
