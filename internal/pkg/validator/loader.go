package validator

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/volvo-cars/connect-access-control/internal/pkg/utils"
	js "github.com/xeipuuv/gojsonschema"
)

type ValidationResult struct {
	FilePath string
	errors   []Error
}

type Error struct {
	Message string
	Field   string
	Value   interface{}
}

// Valid indicates if no errors were found
func (v *ValidationResult) Valid() bool {
	return len(v.errors) == 0
}

// Errors returns the errors that were found
func (v *ValidationResult) Errors() []Error {
	return v.errors
}

type SchemaLoader struct {
	collection *KV[string, string]
}

func NewSchemaLoader() *SchemaLoader {
	return &SchemaLoader{
		collection: NewKV[string, string](),
	}
}

func (l *SchemaLoader) Set(key, schema string) {
	l.collection.Set(key, schema)
}

func (l *SchemaLoader) Get(key string) (string, bool) {
	return l.collection.Get(key)
}

func (l *SchemaLoader) Contains(key string) bool {
	return l.collection.Contains(key)
}

func (l *SchemaLoader) Delete(key string) {
	l.collection.Delete(key)
}

func (l *SchemaLoader) List() map[string]string {
	return l.collection.List()
}

func (l *SchemaLoader) LoadPath(path string) error {
	return l.Load(l.SchemaKey(path), path)
}

func (l *SchemaLoader) SchemaKey(path string) string {
	return strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
}

func (l *SchemaLoader) Load(key, path string) error {
	if l.Contains(key) {
		return nil
	}

	// fmt.Printf(" :-> Loading schema [%s] from [%s]\n", key, path)
	schema, err := utils.YAMLToJSON(path)
	if err != nil {
		return fmt.Errorf("failed to convert schema [%s] to json: %w", path, err)
	}

	l.Set(key, string(schema))
	return nil
}

func (l *SchemaLoader) Validate(schemaPath, filePath string) (*ValidationResult, error) {
	schemaKey := schemaPath

	// fullback with lazy loading schema
	if !l.Contains(schemaKey) {
		if err := l.Load(schemaKey, schemaPath); err != nil {
			return nil, err
		}
	}

	schema, ok := l.Get(schemaKey)
	if !ok {
		return nil, fmt.Errorf("schema [%s] not found", schemaKey)
	}

	return validate(schema, filePath)
}

func validate(schema, filePath string) (*ValidationResult, error) {
	document, err := utils.YAMLToJSON(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to convert document [%s] to json: %w", filePath, err)
	}

	schemaLoader := js.NewStringLoader(schema)
	documentLoader := js.NewBytesLoader(document)
	result, err := js.Validate(schemaLoader, documentLoader)
	if err != nil {
		return nil, fmt.Errorf("failed to validate document [%s] with schema: %w", filePath, err)
	}

	res := &ValidationResult{
		FilePath: filePath,
	}

	if result.Valid() {
		return res, nil
	}

	for _, desc := range result.Errors() {
		err := Error{
			Message: desc.String(),
			Field:   desc.Field(),
			Value:   desc.Value(),
		}

		res.errors = append(res.errors, err)
	}

	return res, nil
}
