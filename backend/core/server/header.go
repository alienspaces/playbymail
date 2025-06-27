package server

import (
	"fmt"
	"net/http"
	"strings"
)

const (
	HeaderXPagination = "X-Pagination"
)

const (
	HeaderContentTypeCSV       = "text/csv"
	HeaderContentTypeJSON      = "application/json"
	HeaderContentTypePDF       = "application/pdf"
	HeaderContentTypeXML       = "application/xml"
	HeaderContentTypeMultipart = "multipart/form-data"
)

// Returns the first found Content-Type header with a map of additional encodings
func RequestContentType(r *http.Request, separateEncodings bool) (string, map[string]string) {
	h := r.Header[http.CanonicalHeaderKey("Content-Type")]
	if len(h) == 0 {
		return "", nil
	}
	if !separateEncodings {
		return h[0], nil
	}

	hParts := strings.Split(h[0], ";")
	if len(hParts) == 1 {
		return hParts[0], nil
	}
	e := map[string]string{}
	for _, ce := range hParts[1:] {
		ceParts := strings.Split(ce, "=")
		// Skip invalid encodings..
		if len(ceParts) != 2 {
			continue
		}
		e[strings.Trim(ceParts[0], " ")] = strings.Trim(ceParts[1], " ")
	}
	return hParts[0], e
}

func XPaginationHeader(collectionLen int, pageSize int) func(http.ResponseWriter) error {
	hasMore := collectionLen > pageSize

	return func(w http.ResponseWriter) error {
		w.Header().Set(HeaderXPagination, XPaginationHeaderValue(hasMore))
		return nil
	}
}

func XPaginationHeaderValue(hasMore bool) string {
	return fmt.Sprintf(`{"has_more":%t}`, hasMore)
}
