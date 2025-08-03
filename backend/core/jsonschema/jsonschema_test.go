package jsonschema

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xeipuuv/gojsonschema"
)

func TestCompile(t *testing.T) {
	cwd, err := os.Getwd()
	require.NoError(t, err, "Getwd returns without error")

	testdataPath := fmt.Sprintf("%s/testdata", cwd)

	schema := SchemaWithReferences{
		Main: Schema{
			Location: "",
			Name:     "test.main.schema.json",
		},
		References: []Schema{
			{
				Location: "",
				Name:     "test.data.schema.json",
			},
		},
	}

	// Resolve the schema location to set the fullPath fields
	schema = ResolveSchemaLocation(testdataPath, schema)

	s, err := Compile(schema)
	require.NoError(t, err, "Compile returns without error")
	require.NotNil(t, s, "Compile returns a compiled schema")

	schema.References = append(schema.References, Schema{
		Location: "",
		Name:     "test.missing.schema.json",
	})

	// Resolve the schema location again for the updated schema
	schema = ResolveSchemaLocation(testdataPath, schema)

	s, err = Compile(schema)
	require.Error(t, err, "Compile returns with error")
	require.Nil(t, s, "Compile does not return a compiled schema with error")
}

func TestSchemaPathFileURL(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Error getting current directory: %v", err)
	}

	// Absolute path
	absolutePath := filepath.Join(cwd, "../../schema/api/account_schema", "account.request-auth.request.schema.json")
	if _, err := os.Stat(absolutePath); err != nil {
		t.Fatalf("File does not exist at absolute path: %v", err)
	}
	absoluteLoader := gojsonschema.NewReferenceLoader(fmt.Sprintf("file://%s", absolutePath))
	_, err = gojsonschema.NewSchema(absoluteLoader)
	if err != nil {
		t.Errorf("Error loading schema with absolute path: %v", err)
	}

	// Relative path (relative to backend/core/jsonschema)
	relativePath := "../../schema/api/account_schema/account.request-auth.request.schema.json"
	if _, err := os.Stat(relativePath); err != nil {
		t.Fatalf("File does not exist at relative path: %v", err)
	}
	relativeLoader := gojsonschema.NewReferenceLoader(fmt.Sprintf("file://%s", relativePath))
	_, err = gojsonschema.NewSchema(relativeLoader)
	if err != nil {
		t.Errorf("Error loading schema with relative path: %v", err)
	}
}

var bench_CompileResult *gojsonschema.Schema

func BenchmarkCompile(b *testing.B) {
	cwd, err := os.Getwd()
	require.NoError(b, err, "Getpwd returns without error")

	testdataPath := fmt.Sprintf("%s/testdata", cwd)

	var r *gojsonschema.Schema
	schema := SchemaWithReferences{
		Main: Schema{
			Location: "",
			Name:     "test.main.schema.json",
		},
		References: []Schema{
			{
				Location: "",
				Name:     "test.data.schema.json",
			},
		},
	}

	// Resolve the schema location to set the fullPath fields
	schema = ResolveSchemaLocation(testdataPath, schema)

	for n := 0; n < b.N; n++ {
		r, _ = Compile(schema)
	}

	bench_CompileResult = r
}

var bench_compileResult *gojsonschema.Schema

func Benchmark_compile(b *testing.B) {
	cwd, err := os.Getwd()
	require.NoError(b, err, "Getpwd returns without error")

	testdataPath := fmt.Sprintf("%s/testdata", cwd)

	var r *gojsonschema.Schema
	schema := SchemaWithReferences{
		Main: Schema{
			Location: "",
			Name:     "test.main.schema.json",
		},
		References: []Schema{
			{
				Location: "",
				Name:     "test.data.schema.json",
			},
		},
	}

	// Resolve the schema location to set the fullPath fields
	schema = ResolveSchemaLocation(testdataPath, schema)

	for n := 0; n < b.N; n++ {
		r, _ = compile(schema)
	}

	bench_compileResult = r
}

func TestValidate(t *testing.T) {
	//
}
