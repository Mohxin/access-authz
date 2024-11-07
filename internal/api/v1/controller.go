package v1

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/volvo-cars/connect-access-control/internal/pkg/authz"
	"github.com/volvo-cars/connect-access-control/internal/pkg/store"
	"github.com/volvo-cars/go-render"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type authzClient interface {
	GetUserByCDSID(ctx context.Context, cdsid string) (authz.User, error)
	GetUserAccess(ctx context.Context, cdsid string, scopes []string) ([]authz.UserAccess, error)
}

type authzStore interface {
	GetClient(key string) (store.Client, error)
	GetClients() ([]store.Client, error)
	GetRole(id string) (store.Role, error)
	GetRoles() ([]store.Role, error)
	GetScope(key string) (store.Scope, error)
	GetScopes() ([]store.Scope, error)
	GetRoleMapping(scopeID, roleID string) (store.RoleMapping, error)
	GetRoleMappings(scopeID string) ([]store.RoleMapping, error)
}

type tracer interface {
	Start(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span)
}

type Controller struct {
	tracer      tracer
	authzStore  authzStore
	authzClient authzClient
}

func NewController(svc authzStore, authzClient authzClient) *Controller {
	return &Controller{
		tracer:      otel.Tracer("controller/iam"),
		authzStore:  svc,
		authzClient: authzClient,
	}
}

func (c *Controller) RegisterRoutes(router chi.Router) {
	router.Route("/iam", func(r chi.Router) {
		r.Route("/clients", func(r chi.Router) {
			r.Get("/", c.getClients)
			r.Get("/{clientID}", c.getClient)
		})

		r.Route("/roles", func(r chi.Router) {
			r.Get("/", c.getRoles)
			r.Get("/{roleID}", c.getRole)
		})

		r.Route("/scopes", func(r chi.Router) {
			r.Get("/", c.getScopes)
			r.Get("/{scopeKey}", c.getScope)
			r.Get("/{scopeKey}/mappings", c.getRoleMappings)
			r.Get("/{scopeKey}/mappings/{roleID}", c.getRoleMapping)
		})

		r.Route("/users", func(r chi.Router) {
			r.Get("/{cdsid}", c.getUser)
			r.Get("/{cdsid}/access", c.getUserAccess)
		})
	})
}

// GetClient godoc
//
//	@Summary		get client
//	@Description	get client by ID
//	@Tags			clients
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Client ID"
//	@Success		200	{object}	ClientResponse
//	@Failure		400	{object}	ErrorResponse
//	@Failure		404	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/iam/clients/{id} [get]
func (c *Controller) getClient(w http.ResponseWriter, r *http.Request) {
	_, span := c.tracer.Start(r.Context(), "controller.getClient")
	defer span.End()

	clientID := chi.URLParam(r, "clientID")
	if clientID == "" {
		c.failure(w, r, http.StatusBadRequest, errors.New("field client id is invalid"))
		return
	}

	client, err := c.authzStore.GetClient(clientID)
	if err != nil {
		if errors.Is(err, store.ErrClientNotFound) {
			c.failure(w, r, http.StatusNotFound, err)
			return
		}

		c.failure(w, r, http.StatusInternalServerError, err)
		return
	}

	response := toClient(client)
	render.Success(w, http.StatusOK, response)
}

// GetClients godoc
//
//	@Summary		get clients
//	@Description	get all clients
//	@Tags			clients
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	ClientsResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/iam/clients [get]
func (c *Controller) getClients(w http.ResponseWriter, r *http.Request) {
	_, span := c.tracer.Start(r.Context(), "controller.getClients")
	defer span.End()

	clients, err := c.authzStore.GetClients()
	if err != nil {
		c.failure(w, r, http.StatusInternalServerError, err)
		return
	}

	response := toClients(clients)
	render.Success(w, http.StatusOK, response)
}

// GetRole godoc
//
//	@Summary		get role
//	@Description	get role by ID
//	@Tags			roles
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Role ID"
//	@Success		200	{object}	RoleResponse
//	@Failure		400	{object}	ErrorResponse
//	@Failure		404	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/iam/roles/{id} [get]
func (c *Controller) getRole(w http.ResponseWriter, r *http.Request) {
	_, span := c.tracer.Start(r.Context(), "controller.getRole")
	defer span.End()

	roleID := chi.URLParam(r, "roleID")
	if roleID == "" {
		c.failure(w, r, http.StatusBadRequest, errors.New("field role key is invalid"))
		return
	}

	role, err := c.authzStore.GetRole(roleID)
	if err != nil {
		if errors.Is(err, store.ErrRoleNotFound) {
			c.failure(w, r, http.StatusNotFound, err)
			return
		}

		c.failure(w, r, http.StatusInternalServerError, err)
		return
	}

	response := toRole(role)
	render.Success(w, http.StatusOK, response)
}

// GetRoles godoc
//
//	@Summary		get roles
//	@Description	get all roles
//	@Tags			roles
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	RolesResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/iam/roles [get]
func (c *Controller) getRoles(w http.ResponseWriter, r *http.Request) {
	_, span := c.tracer.Start(r.Context(), "controller.getRoles")
	defer span.End()

	roles, err := c.authzStore.GetRoles()
	if err != nil {
		c.failure(w, r, http.StatusInternalServerError, err)
		return
	}

	response := toRoles(roles)
	render.Success(w, http.StatusOK, response)
}

// GetScope godoc
//
//	@Summary		get scope
//	@Description	get scope by key
//	@Tags			scopes
//	@Accept			json
//	@Produce		json
//	@Param			scopeKey	path		string	true	"Scope key"
//	@Success		200			{object}	ScopeResponse
//	@Failure		400			{object}	ErrorResponse
//	@Failure		404			{object}	ErrorResponse
//	@Failure		500			{object}	ErrorResponse
//	@Router			/iam/scopes/{scopeKey} [get]
func (c *Controller) getScope(w http.ResponseWriter, r *http.Request) {
	_, span := c.tracer.Start(r.Context(), "controller.getScope")
	defer span.End()

	key := chi.URLParam(r, "scopeKey")
	if key == "" {
		c.failure(w, r, http.StatusBadRequest, errors.New("field scope key is invalid"))
		return
	}

	scopes, err := c.authzStore.GetScope(key)
	if err != nil {
		if errors.Is(err, store.ErrScopeNotFound) {
			c.failure(w, r, http.StatusNotFound, err)
			return
		}

		c.failure(w, r, http.StatusInternalServerError, err)
		return
	}

	response := toScope(scopes)
	render.Success(w, http.StatusOK, response)
}

// GetScopes godoc
//
//	@Summary		get scopes
//	@Description	get all scopes
//	@Tags			scopes
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	ScopesResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/iam/scopes [get]
func (c *Controller) getScopes(w http.ResponseWriter, r *http.Request) {
	_, span := c.tracer.Start(r.Context(), "controller.getScopes")
	defer span.End()

	scopes, err := c.authzStore.GetScopes()
	if err != nil {
		c.failure(w, r, http.StatusInternalServerError, err)
		return
	}

	response := toScopes(scopes)
	render.Success(w, http.StatusOK, response)
}

// GetRoleMapping godoc
//
//	@Summary		get role mapping
//	@Description	get role mapping by scope and role ID
//	@Tags			scopes
//	@Accept			json
//	@Produce		json
//	@Param			scopeKey	path		string	true	"Scope key"
//	@Param			roleID		path		string	true	"Role ID"
//	@Success		200			{object}	RoleMappingResponse
//	@Failure		400			{object}	ErrorResponse
//	@Failure		404			{object}	ErrorResponse
//	@Failure		500			{object}	ErrorResponse
//	@Router			/iam/scopes/{scopeKey}/mappings/{roleID} [get]
func (c *Controller) getRoleMapping(w http.ResponseWriter, r *http.Request) {
	_, span := c.tracer.Start(r.Context(), "controller.getRoleMapping")
	defer span.End()

	scopeID := chi.URLParam(r, "scopeKey")
	if scopeID == "" {
		c.failure(w, r, http.StatusBadRequest, errors.New("field scope key is invalid"))
		return
	}

	roleID := chi.URLParam(r, "roleID")
	if roleID == "" {
		c.failure(w, r, http.StatusBadRequest, errors.New("field role id is invalid"))
		return
	}

	mapping, err := c.authzStore.GetRoleMapping(scopeID, roleID)
	if err != nil {
		if errors.Is(err, store.ErrRoleMappingNotFound) {
			c.failure(w, r, http.StatusNotFound, err)
			return
		}

		c.failure(w, r, http.StatusInternalServerError, err)
		return
	}

	response := toRoleMapping(mapping)
	render.Success(w, http.StatusOK, response)
}

// GetRoleMappings godoc
//
//	@Summary		get role mappings
//	@Description	get all role mappings for a scope
//	@Tags			scopes
//	@Accept			json
//	@Produce		json
//	@Param			scopeKey	path		string	true	"Scope key"
//	@Success		200			{object}	RoleMappingsResponse
//	@Failure		400			{object}	ErrorResponse
//	@Failure		500			{object}	ErrorResponse
//	@Router			/iam/scopes/{scopeKey}/mappings [get]
func (c *Controller) getRoleMappings(w http.ResponseWriter, r *http.Request) {
	_, span := c.tracer.Start(r.Context(), "controller.getRoleMappings")
	defer span.End()

	scopeID := chi.URLParam(r, "scopeKey")
	if scopeID == "" {
		c.failure(w, r, http.StatusBadRequest, errors.New("field scope key is invalid"))
		return
	}

	mappings, err := c.authzStore.GetRoleMappings(scopeID)
	if err != nil {
		c.failure(w, r, http.StatusInternalServerError, err)
		return
	}

	response := toRoleMappings(mappings)
	render.Success(w, http.StatusOK, response)
}

// GetUser godoc
//
//	@Summary		get user
//	@Description	get user by CDSID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			cdsid	path		string	true	"User CDSID"
//	@Success		200		{object}	UserResponse
//	@Failure		400		{object}	ErrorResponse
//	@Failure		404		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Router			/iam/users/{cdsid} [get]
func (c *Controller) getUser(w http.ResponseWriter, r *http.Request) {
	ctx, span := c.tracer.Start(r.Context(), "controller.getUser")
	defer span.End()

	cdsid := chi.URLParam(r, "cdsid")
	if cdsid == "" {
		c.failure(w, r, http.StatusBadRequest, errors.New("field cdsid is invalid"))
		return
	}

	user, err := c.authzClient.GetUserByCDSID(ctx, cdsid)
	if err != nil {
		if errors.Is(err, authz.ErrUserNotFound) {
			c.failure(w, r, http.StatusNotFound, err)
			return
		}

		c.failure(w, r, http.StatusInternalServerError, err)
		return
	}

	response := toUser(user)
	render.Success(w, http.StatusOK, response)
}

// GetUserAccess godoc
//
//	@Summary		get user access
//	@Description	get all access for a user
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			cdsid	path		string		true	"User CDSID"
//	@Param			scope	query		[]string	true	"Scope key"
//	@Success		200		{object}	UserAccessResponse
//	@Failure		400		{object}	ErrorResponse
//	@Failure		404		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Router			/iam/users/{cdsid}/access [get]
func (c *Controller) getUserAccess(w http.ResponseWriter, r *http.Request) {
	ctx, span := c.tracer.Start(r.Context(), "controller.getUserAccess")
	defer span.End()

	cdsid := chi.URLParam(r, "cdsid")
	if cdsid == "" {
		c.failure(w, r, http.StatusBadRequest, errors.New("field user id is invalid"))
		return
	}

	scopes := r.URL.Query()["scope"]
	if len(scopes) == 0 {
		c.failure(w, r, http.StatusBadRequest, errors.New("field scope is invalid"))
		return
	}

	userAccess, err := c.authzClient.GetUserAccess(ctx, cdsid, scopes)
	if err != nil {
		if errors.Is(err, authz.ErrUserNotFound) {
			c.failure(w, r, http.StatusNotFound, err)
			return
		}

		c.failure(w, r, http.StatusInternalServerError, err)
		return
	}

	response := toUserAccesses(userAccess)
	render.Success(w, http.StatusOK, response)
}

func (c *Controller) failure(w http.ResponseWriter, r *http.Request, status int, err error) {
	span := trace.SpanFromContext(r.Context())
	span.SetStatus(codes.Error, err.Error())
	span.RecordError(err)

	render.Failure(w, status, err)
}
