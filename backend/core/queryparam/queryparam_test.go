package queryparam

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/playbymail/core/jsonschema"
)

var queryParamTypes = map[string]jsonschema.JSONType{
	"customer_countries": {
		ElemType: "string",
		IsArray:  true,
	},
	"contract_country": {
		ElemType: "string",
	},
	"currencies": {
		ElemType: "string",
		IsArray:  true,
	},
	"seller_countries": {
		ElemType: "string",
		IsArray:  true,
	},
	"created_at": {
		ElemType: "string",
	},
	"sort_column": {
		ElemType: "string",
	},
	"page_size": {
		ElemType: "number",
	},
	"page_number": {
		ElemType: "number",
	},
}

func TestBuildQueryParams(t *testing.T) {
	type args struct {
		queryParams url.Values
	}
	tests := []struct {
		name string
		args args
		want *QueryParams
	}{
		{
			name: "multiple values, pagination, sort_column",
			args: args{
				queryParams: url.Values{
					"customer_countries":    []string{"US", "CA", "AU"},
					"contract_country":      []string{"US"},
					"currencies[]:ilk":      []string{"usd", "cad", "aud"},
					"seller_countries[]:lk": []string{"U", "A"},
					"created_at:gte":        []string{"2020-01-01"},
					"created_at:lte":        []string{"2022-01-01"},
					"sort_column":           []string{"-created_at", "updated_at"},
					"page_size":             []string{"10"},
					"page_number":           []string{"1"},
				},
			},
			want: &QueryParams{
				Params: map[string][]QueryParam{
					"customer_countries": {
						{
							Val: "US",
						},
						{
							Val: "CA",
						},
						{
							Val: "AU",
						},
					},
					"contract_country": {
						{
							Val: "US",
						},
					},
					"currencies": {
						{
							Val: []string{"usd", "cad", "aud"},
							Op:  OpILike,
						},
					},
					"seller_countries": {
						{
							Val: []string{"U", "A"},
							Op:  OpLike,
						},
					},
					"created_at": {
						{
							Val: "2020-01-01",
							Op:  OpGreaterThanEqual,
						},
						{
							Val: "2022-01-01",
							Op:  OpLessThanEqual,
						},
					},
				},
				SortColumns: []SortColumn{
					{
						Col:          "created_at",
						IsDescending: true,
					},
					{
						Col:          "updated_at",
						IsDescending: false,
					},
				},
				PageSize:   10,
				PageNumber: 1,
			},
		},
		{
			name: "single value, pagination, sort_column",
			args: args{
				queryParams: url.Values{
					"customer_countries": []string{"US"},
					"contract_country":   []string{"US"},
					"sort_column":        []string{"created_at", "-updated_at"},
					"page_size":          []string{"10"},
					"page_number":        []string{"2"},
				},
			},
			want: &QueryParams{
				Params: map[string][]QueryParam{
					"customer_countries": {
						{
							Val: "US",
						},
					},
					"contract_country": {
						{
							Val: "US",
						},
					},
				},
				SortColumns: []SortColumn{
					{
						Col:          "created_at",
						IsDescending: false,
					},
					{
						Col:          "updated_at",
						IsDescending: true,
					},
				},
				PageSize:   10,
				PageNumber: 2,
			},
		},
		{
			name: "no query params",
			args: args{
				queryParams: url.Values{},
			},
			want: &QueryParams{
				Params: map[string][]QueryParam{},
				SortColumns: []SortColumn{
					{
						Col:          DefaultOrderDescendingColumn,
						IsDescending: true,
					},
				},
				PageSize:   DefaultPageSizeInt,
				PageNumber: DefaultPageNumberInt,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualQueryParams, err := BuildQueryParams(nil, tt.args.queryParams, queryParamTypes)
			require.NoError(t, err, "buildQueryParams should not return error")

			// map[string][]QueryParam is unstable requiring the following comparison code
			for key, expected := range tt.want.Params {
				expectedMap := make(map[string]QueryParam, len(expected))
				for _, e := range expected {
					expectedMap[fmt.Sprintf("%#v+%s", e.Val, e.Op)] = e
				}

				params := actualQueryParams.Params[key]
				require.NotEmpty(t, params, "actual params should not be empty")

				actualMap := make(map[string]QueryParam, len(params))
				for _, a := range params {
					actualMap[fmt.Sprintf("%#v+%s", a.Val, a.Op)] = a
				}

				require.Equalf(t, expectedMap, actualMap, "params")
			}

			require.Equalf(t, tt.want.PageSize, actualQueryParams.PageSize, "page size should equal")
			require.Equalf(t, tt.want.SortColumns, actualQueryParams.SortColumns, "sort columns should equal")
			require.Equalf(t, tt.want.PageNumber, actualQueryParams.PageNumber, "page number should equal")
		})
	}
}
