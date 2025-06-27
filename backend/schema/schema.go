// Package schema contains json schema's and related structures for marshalling and unmarshalling data
package schema

// NOTE: When and if a v2 is ever needed specific directories would be structured as:
//
// ./api/v1/*.json
// ./api/v1/api.go
// ./api/v2/*.json
// ./api/v2/api.go
// etc ..

// Request -
type Request struct{}

type QueryParamsPagination struct {
	PageNumber  int      `json:"page_number"`
	PageSize    int      `json:"page_size"`
	SortColumns []string `json:"sort_column"`
}

func (qp QueryParamsPagination) GetPageNumber() int {
	return qp.PageNumber
}

func (qp QueryParamsPagination) GetPageSize() int {
	return qp.PageSize
}

func (qp QueryParamsPagination) GetSortColumns() []string {
	return qp.SortColumns
}

// Response -
type Response struct {
	Error      *ResponseError      `json:"error,omitempty"`
	Pagination *ResponsePagination `json:"pagination,omitempty"`
}

// ResponseError -
type ResponseError struct {
	Code   string `json:"code"`
	Detail string `json:"detail"`
}

// ResponsePagination -
type ResponsePagination struct {
	Number int `json:"page_number"`
	Size   int `json:"page_size"`
	Count  int `json:"page_count"`
}
