package csv

import (
	"testing"

	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/playbymail/core/collection/set"
)

func TestValidateHeader(t *testing.T) {
	type args struct {
		data            func() []byte
		requiredHeaders set.Set[string]
	}
	requiredHeaders := set.New("str", "integer", "float")
	tests := []struct {
		name string
		args args
		err  *ErrHeader
	}{
		{
			name: "valid \\ all headers present \\ ok",
			args: args{
				data:            func() []byte { return []byte(`str,integer,float`) },
				requiredHeaders: requiredHeaders,
			},
		},
		{
			name: "valid \\ all headers present with newlines \\ ok",
			args: args{
				data: func() []byte {
					return []byte(`



str,integer,float`)
				},
				requiredHeaders: requiredHeaders,
			},
		},
		{
			name: "valid \\ all required headers, including extraneous headers \\ ok",
			args: args{
				data:            func() []byte { return []byte(`str,integer,float,additional`) },
				requiredHeaders: requiredHeaders,
			},
		},
		{
			name: "invalid \\ no headers \\ err",
			args: args{
				data:            func() []byte { return []byte{} },
				requiredHeaders: requiredHeaders,
			},
			err: &ErrHeader{Err: ErrMissingHeaders},
		},
		{
			name: "invalid \\ data as headers \\ err",
			args: args{
				data:            func() []byte { return []byte(`this,1,2.0,is missing all headers`) },
				requiredHeaders: requiredHeaders,
			},
			err: &ErrHeader{Err: ErrMissingHeaders},
		},
		{
			name: "invalid \\ duplicate headers \\ err",
			args: args{
				data:            func() []byte { return []byte(`str,integer,float,str`) },
				requiredHeaders: requiredHeaders,
			},
			err: &ErrHeader{
				Fields: []string{"str"},
				Err:    ErrDuplicateHeaders,
			},
		},
		{
			name: "invalid \\ missing required headers \\ err",
			args: args{
				data:            func() []byte { return []byte(`str,integer`) },
				requiredHeaders: requiredHeaders,
			},
			err: &ErrHeader{
				Fields: []string{"float"},
				Err:    ErrMissingRequiredHeaders,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateHeader(tt.args.data(), tt.args.requiredHeaders)

			if tt.err != nil {
				require.Equal(t, *tt.err, err)

			} else {
				require.NoError(t, err)
			}
		})
	}
}
