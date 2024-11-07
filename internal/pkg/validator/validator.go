package validator

import (
	"errors"
	"fmt"
	"path"
	"sync"

	"github.com/volvo-cars/connect-access-control/internal/pkg/utils"
)

const (
	scopeFile                  = "scope.yaml"
	scopeSchemaFile            = "scope.yaml"
	clientFile                 = "client.yaml"
	clientSchemaFile           = "client.yaml"
	roleMappingSchemaFile      = "role-mapping.yaml"
	permissionGroupsFile       = "permission-groups.yaml"
	permissionGroupsSchemaFile = "permission-groups.yaml"
	scopesDir                  = "scopes"
	clientsDir                 = "clients"
	schemaDir                  = "config/schema"
	roleMappingDir             = "role-mapping"
)

type SchemaValidator struct {
	RootDir   string
	SchemaDir string
	loader    *SchemaLoader
}

func NewSchemaValidator(rootDir, schemaDir string) *SchemaValidator {
	return &SchemaValidator{
		RootDir:   rootDir,
		SchemaDir: schemaDir,
		loader:    NewSchemaLoader(),
	}
}

func (v *SchemaValidator) LoadSchema() error {
	schemaFiles := []string{
		clientSchemaFile,
		scopeSchemaFile,
		roleMappingSchemaFile,
		permissionGroupsSchemaFile,
	}

	for _, fileName := range schemaFiles {
		schemaPath := path.Join(v.RootDir, schemaDir, fileName)
		if err := v.loader.Load(schemaPath, schemaPath); err != nil {
			return err
		}
	}

	return nil
}

func (v *SchemaValidator) Validate() ([]*ValidationResult, error) {
	// pre-load schema files
	if err := v.LoadSchema(); err != nil {
		return nil, err
	}

	validators := []func() ([]*ValidationResult, error){
		v.validateClients,
		v.validateScopes,
	}

	results := make([]*ValidationResult, 0)
	for _, validate := range validators {
		vResults, err := validate()
		if err != nil {
			return nil, err
		}

		if len(vResults) > 0 {
			results = append(results, vResults...)
		}
	}

	return results, nil
}

func (v *SchemaValidator) validateClients() ([]*ValidationResult, error) {
	subDirPath := path.Join(v.RootDir, clientsDir)
	dirNames, err := utils.ReadDirNames(subDirPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory [%s]: %w", subDirPath, err)
	}

	results := make([]*ValidationResult, 0)
	for _, dir := range dirNames {
		result, err := v.validateClient(dir)
		if err != nil {
			return nil, err
		}

		if result != nil {
			results = append(results, result)
		}
	}

	return results, nil
}

func (v *SchemaValidator) validateScopes() ([]*ValidationResult, error) {
	dirs, err := v.scanScopesSubDirNames()
	if err != nil {
		return nil, err
	}

	const maxConcurrent int = 5
	var (
		mu      sync.Mutex
		wg      sync.WaitGroup
		results = make([]*ValidationResult, 0)
		sem     = make(chan struct{}, maxConcurrent)
	)

	for _, dir := range dirs {
		wg.Add(1)
		sem <- struct{}{} // Acquire semaphore

		go func(dir string) {
			defer wg.Done()
			defer func() { <-sem }() // Release semaphore

			vResults, err := v.validateScopeDir(dir)
			if err != nil {
				mu.Lock()
				results = append(results, &ValidationResult{errors: []Error{{Message: err.Error()}}})
				mu.Unlock()
				return
			}

			if len(vResults) > 0 {
				mu.Lock()
				results = append(results, vResults...)
				mu.Unlock()
			}
		}(dir)
	}

	wg.Wait()

	return results, nil
}

func (v *SchemaValidator) validateScopeDir(dirPath string) ([]*ValidationResult, error) {
	validators := []func(string) (*ValidationResult, error){
		v.validateScope,
		v.validatePermissionGroups,
	}

	results := make([]*ValidationResult, 0)
	for _, validate := range validators {
		result, err := validate(dirPath)
		if err != nil && !errors.Is(err, utils.ErrNotFound) {
			return nil, err
		}

		if result != nil {
			results = append(results, result)
		}
	}

	result, err := v.validateRoleMappingsDir(dirPath)
	if err != nil && !errors.Is(err, utils.ErrNotFound) {
		return nil, err
	}

	if len(result) > 0 {
		results = append(results, result...)
	}

	return results, nil
}

func (v *SchemaValidator) validateRoleMappingsDir(dirPath string) ([]*ValidationResult, error) {
	subDirPath := path.Join(dirPath, roleMappingDir)
	fileNames, err := utils.ReadFileNames(subDirPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory [%s]: %w", subDirPath, err)
	}

	results := make([]*ValidationResult, 0)
	for _, file := range fileNames {
		result, err := v.validateRoleMapping(file)
		if err != nil {
			return nil, err
		}

		if result != nil {
			results = append(results, result)
		}
	}

	return results, nil
}

func (v *SchemaValidator) validateScope(dirPath string) (*ValidationResult, error) {
	schemaPath := path.Join(v.RootDir, v.SchemaDir, scopeSchemaFile)
	documentPath := path.Join(dirPath, scopeFile)
	return v.loader.Validate(schemaPath, documentPath)
}

func (v *SchemaValidator) validatePermissionGroups(dirPath string) (*ValidationResult, error) {
	schemaPath := path.Join(v.RootDir, v.SchemaDir, permissionGroupsSchemaFile)
	documentPath := path.Join(dirPath, permissionGroupsFile)
	return v.loader.Validate(schemaPath, documentPath)
}

func (v *SchemaValidator) validateClient(dirPath string) (*ValidationResult, error) {
	schemaPath := path.Join(v.RootDir, v.SchemaDir, clientSchemaFile)
	documentPath := path.Join(dirPath, clientFile)
	return v.loader.Validate(schemaPath, documentPath)
}

func (v *SchemaValidator) validateRoleMapping(documentPath string) (*ValidationResult, error) {
	schemaPath := path.Join(v.RootDir, v.SchemaDir, roleMappingSchemaFile)
	return v.loader.Validate(schemaPath, documentPath)
}

func (v *SchemaValidator) scanScopesSubDirNames() ([]string, error) {
	scopesDirPath := path.Join(v.RootDir, scopesDir)
	dirNames, err := utils.ReadDirNames(scopesDirPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory [%s]: %w", scopesDirPath, err)
	}

	return dirNames, nil
}
