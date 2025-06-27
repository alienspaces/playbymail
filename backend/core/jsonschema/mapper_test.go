package jsonschema

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// Define test structs
type TestStruct struct {
	Name      string   `json:"name"`
	Age       int      `json:"age"`
	Admin     bool     `json:"admin,omitempty"`
	StrSlice  []string `json:"str_slice"`
	NumSlice  []int    `json:"num_slice"`
	BoolSlice []bool   `json:"bool_slice"`
}

type NestedStruct struct {
	Score    float64 `json:"score"`
	Multitag string  `json:"multitag,omitempty" db:"-"`
	Untagged string  `db:"untagged"`
	TestStruct
}

func Test_CreateJSONTypeMap(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected map[string]JSONType
	}{
		{
			name:  "struct \\ ok",
			input: TestStruct{},
			expected: map[string]JSONType{
				"name":       {ElemType: "string"},
				"age":        {ElemType: "number"},
				"admin":      {ElemType: "boolean"},
				"str_slice":  {ElemType: "string", IsArray: true},
				"num_slice":  {ElemType: "number", IsArray: true},
				"bool_slice": {ElemType: "boolean", IsArray: true},
			},
		},
		{
			name:  "pointer to struct \\ ok",
			input: &TestStruct{},
			expected: map[string]JSONType{
				"name":       {ElemType: "string"},
				"age":        {ElemType: "number"},
				"admin":      {ElemType: "boolean"},
				"str_slice":  {ElemType: "string", IsArray: true},
				"num_slice":  {ElemType: "number", IsArray: true},
				"bool_slice": {ElemType: "boolean", IsArray: true},
			},
		},
		{
			name:  "nested struct \\ ok",
			input: NestedStruct{},
			expected: map[string]JSONType{
				"name":       {ElemType: "string"},
				"age":        {ElemType: "number"},
				"admin":      {ElemType: "boolean"},
				"score":      {ElemType: "number"},
				"str_slice":  {ElemType: "string", IsArray: true},
				"num_slice":  {ElemType: "number", IsArray: true},
				"bool_slice": {ElemType: "boolean", IsArray: true},
				"multitag":   {ElemType: "string"},
			},
		},
		{
			name:  "pointer to nested struct \\ ok",
			input: &NestedStruct{},
			expected: map[string]JSONType{
				"name":       {ElemType: "string"},
				"age":        {ElemType: "number"},
				"admin":      {ElemType: "boolean"},
				"score":      {ElemType: "number"},
				"str_slice":  {ElemType: "string", IsArray: true},
				"num_slice":  {ElemType: "number", IsArray: true},
				"bool_slice": {ElemType: "boolean", IsArray: true},
				"multitag":   {ElemType: "string"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CreateJSONTypeMap(tt.input)
			require.Equal(t, tt.expected, result)
		})
	}
}
