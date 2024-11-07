package store

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
	"sync"

	"github.com/volvo-cars/connect-access-control/internal/pkg/utils"
)

var (
	ErrDirEmpty            = errors.New("directory is empty")
	ErrScopeNotFound       = errors.New("scope not found")
	ErrClientNotFound      = errors.New("client not found")
	ErrRoleNotFound        = errors.New("role not found")
	ErrRoleMappingNotFound = errors.New("role mapping not found")
	ErrTypeUnsupported     = errors.New("type unsupported")
)

type AccessControlStore struct {
	rootDir      string
	Clients      *KV[string, Client]
	Scopes       *KV[string, Scope]
	Roles        *KV[string, Role]
	RoleMappings *KV[string, []Mapping]
}

func NewAccessControlStore(rootDir string) *AccessControlStore {
	return &AccessControlStore{
		rootDir:      rootDir,
		Clients:      NewKV[string, Client](),
		Scopes:       NewKV[string, Scope](),
		Roles:        NewKV[string, Role](),
		RoleMappings: NewKV[string, []Mapping](),
	}
}

// GetClient retrieves a client from the in-memory database by its key.
func (store *AccessControlStore) GetClient(clientID string) (Client, error) {
	key := clientKey(clientID)
	client, exists := store.Clients.Get(key)
	if !exists {
		return Client{}, ErrClientNotFound
	}

	return client, nil
}

// GetClients retrieves all clients from the in-memory database.
func (store *AccessControlStore) GetClients() ([]Client, error) {
	return store.Clients.Values(), nil
}

// GetScope retrieves a scope from the in-memory database by its key.
func (store *AccessControlStore) GetScope(scopeID string) (Scope, error) {
	key := ScopeKey(scopeID)
	scope, exists := store.Scopes.Get(key)
	if !exists {
		return Scope{}, ErrScopeNotFound
	}

	return scope, nil
}

// GetScopes retrieves all scopes from the in-memory database.
func (store *AccessControlStore) GetScopes() ([]Scope, error) {
	return store.Scopes.Values(), nil
}

// GetRole retrieves a role from the in-memory database by its key.
func (store *AccessControlStore) GetRole(roleID string) (Role, error) {
	key := roleKey(roleID)
	role, exists := store.Roles.Get(key)
	if !exists {
		return Role{}, ErrRoleNotFound
	}

	return role, nil
}

// GetRoles retrieves all roles from the in-memory database.
func (store *AccessControlStore) GetRoles() ([]Role, error) {
	return store.Roles.Values(), nil
}

// GetRoleMapping retrieves a role mapping from the in-memory database by its scope and role keys.
func (store *AccessControlStore) GetRoleMapping(scopeID, roleID string) (RoleMapping, error) {
	key := roleMappingKey(scopeID, roleID)
	mapping, exists := store.RoleMappings.Get(key)
	if !exists {
		return RoleMapping{}, ErrRoleMappingNotFound
	}

	return RoleMapping{
		RoleID:  roleID,
		Mapping: mapping,
	}, nil
}

// GetRoleMappings retrieves all role mappings from the in-memory database by its scope key.
func (store *AccessControlStore) GetRoleMappings(scopeID string) ([]RoleMapping, error) {
	arr := make([]RoleMapping, 0)

	store.RoleMappings.Filter(func(key string, value []Mapping) bool {
		if !startWith(key, fmt.Sprintf("scope:%s/role:", scopeID)) {
			return true
		}

		arr = append(arr, RoleMapping{
			RoleID:  extractRoleID(key),
			Mapping: value,
		})
		return true
	})

	return arr, nil
}

func (store *AccessControlStore) Process() error {
	if err := store.processClients(); err != nil && !errors.Is(err, ErrDirEmpty) {
		return fmt.Errorf("failed to load clients error: %w", err)
	}

	if err := store.processRoles(); err != nil {
		return fmt.Errorf("failed to load roles error: %w", err)
	}

	if err := store.processScopes(); err != nil {
		return fmt.Errorf("failed to load scopes error: %w", err)
	}

	return nil
}

func (store *AccessControlStore) processClients() error {
	clientsDir := path.Join(store.rootDir, "clients")
	dirs, err := utils.ReadDirNames(clientsDir)
	if err != nil {
		return fmt.Errorf("failed to read directory [%s]: %w", clientsDir, err)
	}

	if len(dirs) == 0 {
		return ErrDirEmpty
	}

	for _, dir := range dirs {
		clientFile := path.Join(dir, "client.yaml")
		definition, err := utils.YAMLUnmarshal[ClientDefinition](clientFile)
		if err != nil {
			return fmt.Errorf("failed to unmarshal scope file [%s]: %w", clientFile, err)
		}

		client := definition.Client
		store.Clients.Set(clientKey(client.ID), client)
	}

	return nil
}

func (store *AccessControlStore) processRoles() error {
	roleFile := path.Join(store.rootDir, "config", "roles.yaml")
	roleDefinition, err := utils.YAMLUnmarshal[RoleDefinition](roleFile)
	if err != nil {
		return fmt.Errorf("failed to unmarshal role file [%s]: %w", roleFile, err)
	}

	roles := roleDefinition.Roles
	for _, role := range roles {
		store.Roles.Set(roleKey(role.ID), role)
	}

	return nil
}

func (store *AccessControlStore) processScopes() error {
	dirs, err := store.scanScopesDir()
	if err != nil {
		return err
	}

	const maxConcurrent int = 5
	var (
		mu   sync.Mutex
		wg   sync.WaitGroup
		errs error
		sem  = make(chan struct{}, maxConcurrent)
	)

	for _, dir := range dirs {
		wg.Add(1)
		sem <- struct{}{} // Acquire semaphore

		go func(dir string) {
			defer wg.Done()
			defer func() { <-sem }() // Release semaphore
			// Load the directory
			if err := store.populateScopes(dir); err != nil && !errors.Is(err, utils.ErrNotFound) {
				mu.Lock()
				errs = errors.Join(errs, err)
				mu.Unlock()
			}
		}(dir)
	}

	wg.Wait()

	return errs
}

func (store *AccessControlStore) populateScopes(dirPath string) error {
	scope, err := store.populateScope(dirPath)
	if err != nil {
		return fmt.Errorf("failed to load scope error: %w", err)
	}

	roleMappings, err := store.populateRoleMappings(dirPath)
	if err != nil && !errors.Is(err, ErrDirEmpty) {
		return fmt.Errorf("failed to load role mapping error: %w", err)
	}

	for _, roleMapping := range roleMappings {
		key := roleMappingKey(scope.Key, roleMapping.RoleID)
		if _, exists := store.RoleMappings.Get(key); exists {
			return fmt.Errorf("duplicate role mapping found for scope [%s] and role [%s]", scope.Key, roleMapping.RoleID)
		}

		store.RoleMappings.Set(key, roleMapping.Mapping)
	}

	return nil
}

func (store *AccessControlStore) populateScope(dirPath string) (Scope, error) {
	scopeFile := path.Join(dirPath, "scope.yaml")
	scopeDefinition, err := utils.YAMLUnmarshal[ScopeDefinition](scopeFile)
	if err != nil {
		return Scope{}, fmt.Errorf("failed to unmarshal scope file [%s]: %w", scopeFile, err)
	}

	scope := scopeDefinition.Scope

	permGroupFile := path.Join(dirPath, "permission-groups.yaml")
	permissionGroupsDefinition, err := utils.YAMLUnmarshal[PermissionGroupDefinition](permGroupFile)
	if errors.Is(err, utils.ErrNotFound) {
		return scope, nil
	}
	if err != nil {
		return scope, fmt.Errorf("failed to unmarshal permission groups file [%s]: %w", permGroupFile, err)
	}

	scope.PermissionGroups = permissionGroupsDefinition.PermissionGroups
	store.Scopes.Set(ScopeKey(scope.Key), scope)

	return scope, nil
}

func (store *AccessControlStore) populateRoleMappings(dirPath string) ([]RoleMapping, error) {
	roleMappingDir := path.Join(dirPath, "role-mapping")
	files, err := os.ReadDir(roleMappingDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory [%s]: %w", roleMappingDir, err)
	}

	if len(files) == 0 {
		return nil, ErrDirEmpty
	}

	roleMappings := make([]RoleMapping, 0, len(files))
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		roleMappingFilePath := path.Join(roleMappingDir, file.Name())
		definition, err := utils.YAMLUnmarshal[RoleMappingDefinition](roleMappingFilePath)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal role mapping file [%s]: %w", roleMappingFilePath, err)
		}

		roleMappings = append(roleMappings, definition.RoleAssignment)
	}

	return roleMappings, nil
}

func (store *AccessControlStore) scanScopesDir() ([]string, error) {
	scopesDirPath := path.Join(store.rootDir, "scopes")
	dirs, err := utils.ReadDirNames(scopesDirPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory [%s]: %w", scopesDirPath, err)
	}

	return dirs, nil
}

func startWith(key, prefix string) bool {
	return strings.HasPrefix(key, prefix)
}

func extractRoleID(key string) string {
	parts := strings.Split(key, "role:")
	return parts[len(parts)-1]
}

func clientKey(role string) string {
	return "client:" + role
}

func roleKey(role string) string {
	return "role:" + role
}

func ScopeKey(scope string) string {
	return "scope:" + scope
}

func roleMappingKey(scope, role string) string {
	return fmt.Sprintf("scope:%s/role:%s", scope, role)
}
