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
			name: "valid when all headers present then ok",
			args: args{
				data:            func() []byte { return []byte(`str,integer,float`) },
				requiredHeaders: requiredHeaders,
			},
		},
		{
			name: "valid when all headers present with newlines then ok",
			args: args{
				data: func() []byte {
					return []byte(`



str,integer,float`)
				},
				requiredHeaders: requiredHeaders,
			},
		},
		{
			name: "valid when all required headers, including extraneous headers then ok",
			args: args{
				data:            func() []byte { return []byte(`str,integer,float,additional`) },
				requiredHeaders: requiredHeaders,
			},
		},
		{
			name: "invalid when no headers then err",
			args: args{
				data:            func() []byte { return []byte{} },
				requiredHeaders: requiredHeaders,
			},
			err: &ErrHeader{Err: ErrMissingHeaders},
		},
		{
			name: "invalid when data as headers then err",
			args: args{
				data:            func() []byte { return []byte(`this,1,2.0,is missing all headers`) },
				requiredHeaders: requiredHeaders,
			},
			err: &ErrHeader{Err: ErrMissingHeaders},
		},
		{
			name: "invalid when duplicate headers then err",
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
			name: "invalid when missing required headers then err",
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
