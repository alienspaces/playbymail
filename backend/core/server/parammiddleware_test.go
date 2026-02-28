package server

import (
	"fmt"
	"net/url"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/core/jsonschema"
)

var queryParamTypes = map[string]jsonschema.JSONType{
	"id":       {ElemType: "string"},
	"str":      {ElemType: "string"},
	"num":      {ElemType: "number"},
	"bool":     {ElemType: "boolean"},
	"bool_arr": {ElemType: "boolean", IsArray: true},
	"num_arr":  {ElemType: "number", IsArray: true},
	"str_arr":  {ElemType: "string", IsArray: true},
}

func Test_extractPathParams(t *testing.T) {

	type testcase struct {
		name         string
		path         string
		expectParams []string
	}

	tests := []testcase{
		{
			name:         "two parameters with trailing parameter",
			path:         "/animals/:animal_id/humans/:human_id",
			expectParams: []string{"animal_id", "human_id"},
		},
		{
			name:         "two parameters without trailing parameter",
			path:         "/animals/:animal_id/humans/:human_id/overlords",
			expectParams: []string{"animal_id", "human_id"},
		},
		{
			name:         "one parameter without trailing parameter",
			path:         "/animals/:animal_id/humans",
			expectParams: []string{"animal_id"},
		},
		{
			name:         "no parameters",
			path:         "/animals",
			expectParams: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := extractPathParams(tt.path)
			require.Equal(t, tt.expectParams, params, "Extracted params equals expected")
		})
	}
}

func Test_validateQueryParameters(t *testing.T) {

	l, _, _ := newDefaultDependencies(t)

	type args struct {
		q           url.Values
		paramSchema jsonschema.SchemaWithReferences
	}

	type testcase struct {
		name    string
		args    args
		errCode coreerror.Code
	}
	tests := []testcase{
		{
			name: "nil query params then ok",
			args: args{
				q: nil,
			},
		},
		{
			name: "id query param then ok",
			args: args{
				q: url.Values{
					"id": []string{"a87feca8-d6f0-4794-98c7-037b30219520"},
				},
			},
		},
		{
			name: "string query param then ok",
			args: args{
				q: url.Values{
					"str": []string{"asdf"},
				},
			},
		},
		{
			name: "number query param then ok",
			args: args{
				q: url.Values{
					"num": []string{"123"},
				},
			},
		},
		{
			name: "bool query param then ok",
			args: args{
				q: url.Values{
					"bool": []string{"true"},
				},
			},
		},
		{
			name: "string array query param then ok",
			args: args{
				q: url.Values{
					"str_arr": []string{"1", "2"},
				},
			},
		},
		{
			name: "number array query param then ok",
			args: args{
				q: url.Values{
					"num_arr": []string{"1", "2"},
				},
			},
		},
		{
			name: "bool array query param then ok",
			args: args{
				q: url.Values{
					"bool_arr": []string{"true", "false"},
				},
			},
		},
		{
			name: "string, number, boolean query params then ok",
			args: args{
				q: url.Values{
					"id":   []string{"a87feca8-d6f0-4794-98c7-037b30219520"},
					"str":  []string{"asdf"},
					"num":  []string{"123"},
					"bool": []string{"false"},
				},
			},
		},
		{
			name: "empty string query param value then err",
			args: args{
				q: url.Values{
					"str": []string{""},
				},
			},
			errCode: coreerror.GetRegistryError(coreerror.ErrorCodeSchemaValidation).ErrorCode,
		},
		{
			name: "number query param value below min then err",
			args: args{
				q: url.Values{
					"num": []string{"0"},
				},
			},
			errCode: coreerror.GetRegistryError(coreerror.ErrorCodeSchemaValidation).ErrorCode,
		},
		{
			name: "string query param with array value when string query param is not array type then ok",
			args: args{
				q: url.Values{
					"str[]": []string{"1", "2"},
				},
			},
		},
		{
			name: "number query param with array value when number query param is not array type then ok",
			args: args{
				q: url.Values{
					"num[]": []string{"1", "2"},
				},
			},
		},
		{
			name: "boolean query param with array value when boolean query param is not array type then ok",
			args: args{
				q: url.Values{
					"bool[]": []string{"true", "false"},
				},
			},
		},
		{
			name: "additional property then ok",
			args: args{
				q: url.Values{
					"asdf": []string{"0"},
				},
			},
		},
		{
			name: "multiple errors - empty string query param value, and num query param with array values and some below min then err",
			args: args{
				q: url.Values{
					"str":     []string{""},
					"num_arr": []string{"0", "1"},
				},
			},
			errCode: coreerror.GetRegistryError(coreerror.ErrorCodeSchemaValidation).ErrorCode,
		},
	}

	cwd, err := os.Getwd()
	require.NoError(t, err, "Getwd returns without error")

	testdataPath := fmt.Sprintf("%s/testdata", cwd)

	for i := range tests {
		tests[i].args.paramSchema = jsonschema.SchemaWithReferences{
			Main: jsonschema.Schema{
				Location: "",
				Name:     "test.main.schema.json",
			},
			References: []jsonschema.Schema{
				{
					Location: "",
					Name:     "test.data.schema.json",
				},
			},
		}

		// Resolve the schema location to set the fullPath fields
		tests[i].args.paramSchema = jsonschema.ResolveSchemaLocation(testdataPath, tests[i].args.paramSchema)
	}

	noSchemaTests := []testcase{
		{
			name: "query param with no schema",
			args: args{
				q: url.Values{
					"asdf": []string{"0"},
				},
			},
			errCode: coreerror.GetRegistryError(coreerror.ErrorCodeInvalidParam).ErrorCode,
		},
	}

	tests = append(tests, noSchemaTests...)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err = validateParams(l, tt.args.q, &ValidateParamsConfig{
				Schema:          tt.args.paramSchema,
				queryParamTypes: queryParamTypes,
			})
			if tt.errCode != "" {
				require.Error(t, err, "validateQueryParameters should return err")
				coreerrorErr, conversionErr := coreerror.ToError(err)
				require.Nil(t, conversionErr, "should not have an err that is not wrapped")

				require.Equal(t, tt.errCode, coreerrorErr.ErrorCode)

				e := coreerror.ProcessParamError(err)
				coreerrorErr, conversionErr = coreerror.ToError(e)
				require.Nil(t, conversionErr, "should not have an err that is not wrapped")

				require.Equal(t, coreerror.GetRegistryError(coreerror.ErrorCodeInvalidParam).ErrorCode, coreerrorErr.ErrorCode)
			} else {
				require.NoError(t, err, "validateParams should not return err")
			}
		})
	}
}

func Test_paramsToJSON(t *testing.T) {

	type args struct {
		q url.Values
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "nil query params then empty",
			want: "",
		},
		{
			name: "no query params then empty",
			args: args{
				q: url.Values{},
			},
			want: "",
		},
		{
			name: "string query param with no value then ok",
			args: args{
				q: url.Values{
					"str": []string{},
				},
			},
			want: `{"str":""}`,
		},
		{
			name: "string query param with one string value then ok",
			args: args{
				q: url.Values{
					"str": []string{"a"},
				},
			},
			want: `{"str":"a"}`,
		},
		{
			name: "string query param with operator and one string value then ok",
			args: args{
				q: url.Values{
					"str:ilk": []string{"a"},
				},
			},
			want: `{"str":"a"}`,
		},
		{
			name: "boolean query param with no value then ok",
			args: args{
				q: url.Values{
					"bool": []string{},
				},
			},
			want: `{"bool":""}`,
		},
		{
			name: "boolean query param with one bool value then ok",
			args: args{
				q: url.Values{
					"bool": []string{"false"},
				},
			},
			want: `{"bool":false}`,
		},
		{
			name: "boolean query param with operator and one bool value then ok",
			args: args{
				q: url.Values{
					"bool:ne": []string{"true"},
				},
			},
			want: `{"bool":true}`,
		},
		{
			name: "number query param with one number value then ok",
			args: args{
				q: url.Values{
					"num": []string{"123"},
				},
			},
			want: `{"num":123}`,
		},
		{
			name: "number query param with operator and one number value then ok",
			args: args{
				q: url.Values{
					"num:lte": []string{"123"},
				},
			},
			want: `{"num":123}`,
		},
		{
			name: "string query param with array empty value then ok",
			args: args{
				q: url.Values{
					"str[]": []string{},
				},
			},
			want: `{"str[]":[]}`,
		},
		{
			name: "string query param with array with one string value then ok",
			args: args{
				q: url.Values{
					"str[]": []string{"a"},
				},
			},
			want: `{"str[]":["a"]}`,
		},
		{
			name: "string query param with array with multiple string values then ok",
			args: args{
				q: url.Values{
					"str[]": []string{"a", "az"},
				},
			},
			want: `{"str[]":["a","az"]}`,
		},
		{
			name: "string query param with array with bool, string, number values then ok",
			args: args{
				q: url.Values{
					"str[]": []string{"a", "123", "true"},
				},
			},
			want: `{"str[]":["a","123","true"]}`,
		},
		{
			name: "boolean query param with array with empty value then ok",
			args: args{
				q: url.Values{
					"bool[]": []string{},
				},
			},
			want: `{"bool[]":[]}`,
		},
		{
			name: "boolean query param with array with one boolean value then ok",
			args: args{
				q: url.Values{
					"bool[]": []string{"true"},
				},
			},
			want: `{"bool[]":[true]}`,
		},
		{
			name: "boolean query param with array with multiple boolean values then ok",
			args: args{
				q: url.Values{
					"bool[]": []string{"true", "false"},
				},
			},
			want: `{"bool[]":[true,false]}`,
		},
		{
			name: "number query param with array with empty value then ok",
			args: args{
				q: url.Values{
					"num[]": []string{},
				},
			},
			want: `{"num[]":[]}`,
		},
		{
			name: "number query param with array with one number value then ok",
			args: args{
				q: url.Values{
					"num[]": []string{"1"},
				},
			},
			want: `{"num[]":[1]}`,
		},
		{
			name: "number query param with array with multiple number values then ok",
			args: args{
				q: url.Values{
					"num[]": []string{"1", "2"},
				},
			},
			want: `{"num[]":[1,2]}`,
		},
		{
			name: "string array query param with empty value then ok",
			args: args{
				q: url.Values{
					"str_arr": []string{},
				},
			},
			want: `{"str_arr":[]}`,
		},
		{
			name: "string array query param with one string value then ok",
			args: args{
				q: url.Values{
					"str_arr": []string{"a"},
				},
			},
			want: `{"str_arr":["a"]}`,
		},
		{
			name: "string array query param with multiple string values then ok",
			args: args{
				q: url.Values{
					"str_arr": []string{"a", "az"},
				},
			},
			want: `{"str_arr":["a","az"]}`,
		},
		{
			name: "string array query param with bool, string, number values then ok",
			args: args{
				q: url.Values{
					"str_arr": []string{"a", "123", "true"},
				},
			},
			want: `{"str_arr":["a","123","true"]}`,
		},
		{
			name: "boolean array query param with empty value then ok",
			args: args{
				q: url.Values{
					"bool_arr": []string{},
				},
			},
			want: `{"bool_arr":[]}`,
		},
		{
			name: "boolean array query param with one boolean value then ok",
			args: args{
				q: url.Values{
					"bool_arr": []string{"true"},
				},
			},
			want: `{"bool_arr":[true]}`,
		},
		{
			name: "boolean array query param with multiple boolean values then ok",
			args: args{
				q: url.Values{
					"bool_arr": []string{"true", "false"},
				},
			},
			want: `{"bool_arr":[true,false]}`,
		},
		{
			name: "number array query param with empty value then ok",
			args: args{
				q: url.Values{
					"num_arr": []string{},
				},
			},
			want: `{"num_arr":[]}`,
		},
		{
			name: "number array query param with one number value then ok",
			args: args{
				q: url.Values{
					"num_arr": []string{"1"},
				},
			},
			want: `{"num_arr":[1]}`,
		},
		{
			name: "number array query param with multiple number values then ok",
			args: args{
				q: url.Values{
					"num_arr": []string{"1", "2"},
				},
			},
			want: `{"num_arr":[1,2]}`,
		},
		{
			name: "unknown query param then string query param values",
			args: args{
				q: url.Values{
					"unknown": []string{"abc"},
				},
			},
			want: `{"unknown":"abc"}`,
		},
		{
			name: "unknown array query param with string, number, boolean values then first param value populated",
			args: args{
				q: url.Values{
					"unknown[]": []string{"abc", "123", "false"},
				},
			},
			want: `{"unknown[]":["abc","123","false"]}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := paramsToJSON(tt.args.q, queryParamTypes)
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}
