package error

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestToError(t *testing.T) {
	type args struct {
		e error
	}
	tests := []struct {
		name    string
		args    args
		want    Error
		wantErr bool
	}{
		{
			name: "invalid JSON",
			args: args{
				e: GetRegistryError(ErrorCodeInvalidData),
			},
			want: GetRegistryError(ErrorCodeInvalidData),
		},
		{
			name: "invalid query param",
			args: args{
				e: GetRegistryError(ErrorCodeInvalidParam),
			},
			want: GetRegistryError(ErrorCodeInvalidParam),
		},
		{
			name: "unauthenticated",
			args: args{
				e: GetRegistryError(ErrorCodeUnauthenticated),
			},
			want: GetRegistryError(ErrorCodeUnauthenticated),
		},
		{

			name: "unauthorized",
			args: args{
				e: GetRegistryError(ErrorCodeUnauthorized),
			},
			want: GetRegistryError(ErrorCodeUnauthorized),
		},
		{
			name: "not found",
			args: args{
				e: GetRegistryError(ErrorCodeNotFound),
			},
			want: GetRegistryError(ErrorCodeNotFound),
		},
		{
			name: "internal",
			args: args{
				e: GetRegistryError(ErrorCodeInternal),
			},
			want: GetRegistryError(ErrorCodeInternal),
		},
		{
			name: "error",
			args: args{
				e: fmt.Errorf("error"),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToError(tt.args.e)
			if tt.wantErr {
				require.NotNil(t, err, "should not be able to convert error type to coreerror.Error")
				require.Zero(t, got, "coreerror.Error should have zero value for error type that cannot be converted to coreerror.Error")
			} else {
				require.Nil(t, err, "should be able to convert error type to coreerror.Error")
				require.Equal(t, tt.want, got, "coreerror.Error converted should be the same as expected")
			}
		})
	}
}
