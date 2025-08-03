package jsonschema

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/xeipuuv/gojsonschema"
)

type cacheKey string

var schemaCache = map[cacheKey]*gojsonschema.Schema{}
var mu = sync.Mutex{}

// SchemaWithReferences -
type SchemaWithReferences struct {
	Main               Schema
	References         []Schema
	cacheKey           cacheKey
	referenceFullPaths []string
}

func (s *SchemaWithReferences) IsEmpty() bool {
	return s.Main.Name == "" && s.Main.Location == ""
}

func (s *SchemaWithReferences) GetReferencesFullPaths() []string {
	if len(s.referenceFullPaths) != 0 {
		return s.referenceFullPaths
	}
	var paths []string
	for _, r := range s.References {
		paths = append(paths, r.GetFullPath())
	}
	s.referenceFullPaths = paths
	return s.referenceFullPaths
}

// Schema -
type Schema struct {
	Location string
	Name     string
	fullPath string
}

func (s *Schema) GetFullPath() string {
	return s.fullPath
}

// Utility functions
func ValidateJSON(schema SchemaWithReferences, document []byte) error {
	result, err := Validate(schema, document)
	if err != nil {
		return err
	}
	return MapError(result)
}

func Validate(schema SchemaWithReferences, data interface{}) (*gojsonschema.Result, error) {

	s, err := Compile(schema)
	if err != nil {
		return nil, err
	}

	var dataLoader gojsonschema.JSONLoader
	switch d := data.(type) {
	case nil:
		return nil, fmt.Errorf("data is nil")
	case []byte:
		dataLoader = gojsonschema.NewStringLoader(string(d))
	case string:
		dataLoader = gojsonschema.NewStringLoader(d)
	default:
		dataLoader = gojsonschema.NewGoLoader(d)
	}

	return s.Validate(dataLoader)
}

func MapError(result *gojsonschema.Result) error {
	if result.Valid() {
		return nil
	}

	var errStr string

	for _, e := range result.Errors() {
		if errStr == "" {
			errStr = e.String()
			continue
		}
		errStr = fmt.Sprintf("%s, %s", errStr, e.String())
	}

	return errors.New(errStr)
}

// Compile caches JSON schema compilation
func Compile(sr SchemaWithReferences) (*gojsonschema.Schema, error) {
	key := generateCacheKey(sr)
	cached, ok := schemaCache[key]
	if !ok {
		mu.Lock()
		defer mu.Unlock()

		if cached, ok = schemaCache[key]; ok {
			return cached, nil
		}
	} else {
		return cached, nil
	}

	s, err := compile(sr)
	if err != nil {
		return nil, err
	}

	schemaCache[key] = s

	return s, nil
}

func generateCacheKey(s SchemaWithReferences) cacheKey {
	if s.cacheKey != "" {
		return s.cacheKey
	}

	var refs []string
	for _, r := range s.References {
		refs = append(refs, r.GetFullPath())
	}

	key := s.Main.GetFullPath() + strings.Join(refs, "-")
	s.cacheKey = cacheKey(key)
	return s.cacheKey
}

// Internal non-caching JSON schema compilation
func compile(sr SchemaWithReferences) (*gojsonschema.Schema, error) {

	sl := gojsonschema.NewSchemaLoader()
	sl.Validate = true
	sl.AutoDetect = false
	sl.Draft = gojsonschema.Draft7

	for _, ref := range sr.References {
		// Read the schema file content directly
		content, err := os.ReadFile(ref.GetFullPath())
		if err != nil {
			return nil, fmt.Errorf("failed reading reference schema file >%s< err >%w<", ref.GetFullPath(), err)
		}
		loader := gojsonschema.NewStringLoader(string(content))
		err = sl.AddSchemas(loader)
		if err != nil {
			return nil, fmt.Errorf("failed adding reference schema >%s< err >%w<", ref.GetFullPath(), err)
		}
	}

	// Read the main schema file content directly
	content, err := os.ReadFile(sr.Main.GetFullPath())
	if err != nil {
		return nil, fmt.Errorf("failed reading main schema file >%s< err >%w<", sr.Main.GetFullPath(), err)
	}
	loader := gojsonschema.NewStringLoader(string(content))
	s, err := sl.Compile(loader)
	if err != nil {
		return nil, fmt.Errorf("failed adding main schema >%s< err >%w<, are you sure you've loaded all required reference schemas?", sr.Main.GetFullPath(), err)
	}

	return s, nil
}

func ResolveSchemaLocation(schemaPath string, cfg SchemaWithReferences) SchemaWithReferences {

	if cfg.Main.Location != "" {
		cfg.Main.fullPath = fmt.Sprintf("%s/%s/%s", schemaPath, cfg.Main.Location, cfg.Main.Name)
	} else {
		cfg.Main.fullPath = fmt.Sprintf("%s/%s", schemaPath, cfg.Main.Name)
	}

	for i := range cfg.References {
		if cfg.References[i].Location != "" {
			cfg.References[i].fullPath = fmt.Sprintf("%s/%s/%s", schemaPath, cfg.References[i].Location, cfg.References[i].Name)
		} else {
			cfg.References[i].fullPath = fmt.Sprintf("%s/%s", schemaPath, cfg.References[i].Name)
		}
	}

	return cfg
}
