package server

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRequestContentType(t *testing.T) {
	tests := []struct {
		name              string
		r                 *http.Request
		separateEncodings bool
		want              string
		want1             map[string]string
	}{
		{
			name: "with encoding then wanting separate encodings",
			r: &http.Request{
				Header: map[string][]string{
					"Content-Type": {"application/json; charset=utf-8"},
				},
			},
			separateEncodings: true,
			want:              "application/json",
			want1:             map[string]string{"charset": "utf-8"},
		},
		{
			name: "with encoding then not watning separate encodings",
			r: &http.Request{
				Header: map[string][]string{
					"Content-Type": {"application/json; charset=utf-8"},
				},
			},
			separateEncodings: false,
			want:              "application/json; charset=utf-8",
			want1:             nil,
		},
		{
			name: "without encoding then wanting separate encodings",
			r: &http.Request{
				Header: map[string][]string{
					"Content-Type": {"application/json"},
				},
			},
			separateEncodings: true,
			want:              "application/json",
			want1:             nil,
		},
		{
			name: "without encoding then not watning separate encodings",
			r: &http.Request{
				Header: map[string][]string{
					"Content-Type": {"application/json"},
				},
			},
			separateEncodings: false,
			want:              "application/json",
			want1:             nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := RequestContentType(tt.r, tt.separateEncodings)
			require.Equal(t, tt.want, got, "ContentType returns expected content type")
			require.Equal(t, tt.want1, got1, "ContentType returns expected content type")
		})
	}
}
