package queryparam

import (
	"net/url"
	"strconv"
	"strings"

	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/core/jsonschema"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
)

type Operator string

const (
	OpNotEqual         Operator = "ne"
	OpGreaterThanEqual Operator = "gte"
	OpGreaterThan      Operator = "gt"
	OpLessThanEqual    Operator = "lte"
	OpLessThan         Operator = "lt"
	OpLike             Operator = "lk"
	OpILike            Operator = "ilk"
)

type QueryParams struct {
	Params      map[string][]QueryParam
	SortColumns []SortColumn
	PageSize    int
	PageNumber  int
}

// GetParam returns the first query param for the given key
func (qp *QueryParams) GetParam(key string) []QueryParam {
	return qp.Params[key]
}

// GetParamValue returns the value of the first query param for the given key
func (qp *QueryParams) GetParamValue(key string) any {
	return qp.Params[key][0].Val
}

// GetParamValueString returns the value of the first query param for the given key as a string
func (qp *QueryParams) GetParamValueString(key string) string {
	return qp.Params[key][0].Val.(string)
}

// GetParamValueInt returns the value of the first query param for the given key as an int
func (qp *QueryParams) GetParamValueInt(key string) int {
	return qp.Params[key][0].Val.(int)
}

// GetParamValueFloat returns the value of the first query param for the given key as a float64
func (qp *QueryParams) GetParamValueFloat(key string) float64 {
	return qp.Params[key][0].Val.(float64)
}

// GetParamValueBool returns the value of the first query param for the given key as a bool
func (qp *QueryParams) GetParamValueBool(key string) bool {
	return qp.Params[key][0].Val.(bool)
}

// GetParamValues returns all values for the given key
func (qp *QueryParams) GetParamValues(key string) []any {
	params, exists := qp.Params[key]
	if !exists || len(params) == 0 {
		return []any{}
	}
	values := make([]any, len(params))
	for i, param := range params {
		values[i] = param.Val
	}
	return values
}

// GetParamValuesString returns all values for the given key as strings
func (qp *QueryParams) GetParamValuesString(key string) []string {
	params, exists := qp.Params[key]
	if !exists || len(params) == 0 {
		return []string{}
	}
	values := make([]string, len(params))
	for i, param := range params {
		values[i] = param.Val.(string)
	}
	return values
}

// GetParamValuesInt returns all values for the given key as ints
func (qp *QueryParams) GetParamValuesInt(key string) []int {
	params, exists := qp.Params[key]
	if !exists || len(params) == 0 {
		return []int{}
	}
	values := make([]int, len(params))
	for i, param := range params {
		values[i] = param.Val.(int)
	}
	return values
}

// GetParamValuesFloat returns all values for the given key as float64s
func (qp *QueryParams) GetParamValuesFloat(key string) []float64 {
	params, exists := qp.Params[key]
	if !exists || len(params) == 0 {
		return []float64{}
	}
	values := make([]float64, len(params))
	for i, param := range params {
		values[i] = param.Val.(float64)
	}
	return values
}

// GetParamValuesBool returns all values for the given key as bools
func (qp *QueryParams) GetParamValuesBool(key string) []bool {
	params, exists := qp.Params[key]
	if !exists || len(params) == 0 {
		return []bool{}
	}
	values := make([]bool, len(params))
	for i, param := range params {
		values[i] = param.Val.(bool)
	}
	return values
}

type QueryParam struct {
	Val any
	Op  Operator
}

type SortColumn struct {
	Col          string
	IsDescending bool
}

const (
	PageSize      = "page_size"
	PageNumber    = "page_number"
	SortColumnKey = "sort_column"
)

const (
	DefaultPageSize              = "10"
	DefaultPageSizeInt           = 10
	DefaultPageNumber            = "1"
	DefaultPageNumberInt         = 1
	DefaultOrderDescendingColumn = "created_at"
)

func BuildQueryParams(l logger.Logger, q url.Values, queryParamTypes map[string]jsonschema.JSONType) (*QueryParams, error) {
	if len(q) == 0 {
		return &QueryParams{
			Params: map[string][]QueryParam{},
			SortColumns: []SortColumn{
				{
					Col:          DefaultOrderDescendingColumn,
					IsDescending: true,
				},
			},
			PageSize:   DefaultPageSizeInt,
			PageNumber: DefaultPageNumberInt,
		}, nil
	}

	qp := make(map[string][]string, len(q))
	for key, value := range q {
		if len(value) == 0 {
			continue
		}

		qp[strings.TrimSuffix(key, "[]")] = value
	}

	qp, sortColumns, err := extractSortColumns(qp)
	if err != nil {
		l.Warn("failed to resolve sort_column params", err)
		return nil, err
	}

	qp, pageSize, err := extractPageSize(qp)
	if err != nil {
		l.Warn("failed to resolve page_size >%v<", err)
		return nil, err
	}

	qp, pageNumber, err := extractPageNumber(qp)
	if err != nil {
		l.Warn("failed to resolve page_number >%v<", err)
		return nil, err
	}

	return &QueryParams{
		Params:      resolveQueryParamOps(qp, queryParamTypes),
		SortColumns: sortColumns,
		PageSize:    pageSize,
		PageNumber:  pageNumber,
	}, nil
}

// extractPageSize mutates qp
func extractPageSize(qp map[string][]string) (map[string][]string, int, error) {
	qp, pageSize, err := extractIntQueryParam(qp, PageSize, DefaultPageSize)
	if err != nil {
		return qp, 0, err
	}
	if pageSize < 1 {
		return qp, 0, coreerror.NewParamError("Query parameter >%s< is less than 1 >%d<", PageSize, pageSize)
	}

	return qp, pageSize, nil
}

// extractPageNumber mutates qp
func extractPageNumber(qp map[string][]string) (map[string][]string, int, error) {
	qp, pageNumber, err := extractIntQueryParam(qp, PageNumber, DefaultPageNumber)
	if err != nil {
		return qp, 0, err
	}
	if pageNumber < 1 {
		return qp, 0, coreerror.NewParamError("Query parameter >%s< is less than 1 >%d<", PageNumber, pageNumber)
	}

	return qp, pageNumber, nil
}

// extractIntQueryParam extracts the value associated with the key and removes
// the key, mutating the params map. The params map value is expected to be a
// string slice.
func extractIntQueryParam(qp map[string][]string, key string, defaultValue string) (map[string][]string, int, error) {
	qp, valueStr := extractQueryParam(qp, key)
	if valueStr == nil {
		valueStr = []string{defaultValue}
	}

	if len(valueStr) != 1 {
		return qp, 0, coreerror.NewParamError("query parameter >%s< should be a single value but is >%+v<", key, valueStr)
	}

	valueInt, err := strconv.Atoi(valueStr[0])
	if err != nil {
		return qp, 0, coreerror.NewParamError("query parameter >%s< has an invalid value >%+v<", key, valueStr)
	}

	return qp, valueInt, nil
}

// extractSortColumns mutates qp
func extractSortColumns(qp map[string][]string) (map[string][]string, []SortColumn, error) {
	qp, sortColumnValues := extractQueryParam(qp, SortColumnKey)

	var sortColumns []SortColumn
	if sortColumnValues == nil {
		sortColumns = []SortColumn{
			{
				Col:          DefaultOrderDescendingColumn,
				IsDescending: true,
			},
		}

		return qp, sortColumns, nil
	}

	for _, col := range sortColumnValues {
		isDescending := strings.HasPrefix(col, "-")
		if isDescending {
			col = strings.TrimPrefix(col, "-")
		}

		sortColumns = append(sortColumns, SortColumn{
			Col:          col,
			IsDescending: isDescending,
		})
	}

	return qp, sortColumns, nil
}

func extractQueryParam(qp map[string][]string, key string) (map[string][]string, []string) {
	value, ok := qp[key]
	if !ok {
		return qp, nil
	}

	delete(qp, key)
	return qp, value
}

func resolveQueryParamOps(qp url.Values, queryParamTypes map[string]jsonschema.JSONType) map[string][]QueryParam {
	result := map[string][]QueryParam{}

	for key, values := range qp {
		var op Operator
		// assume colons in query params only appear to separate the query param
		// from the operator
		split := strings.Split(key, ":")
		if len(split) > 1 {
			key = split[0]
			op = Operator(split[1])
		}

		rawKey := strings.ReplaceAll(key, "[]", "")
		if (queryParamTypes[rawKey].IsArray || strings.HasSuffix(key, "[]")) && (op == OpILike || op == OpLike) {
			result[rawKey] = append(result[rawKey], QueryParam{
				Op:  op,
				Val: values,
			})

			continue
		}

		for _, v := range values {
			p := QueryParam{
				Op:  op,
				Val: v,
			}
			result[rawKey] = append(result[rawKey], p)
		}
	}

	return result
}
