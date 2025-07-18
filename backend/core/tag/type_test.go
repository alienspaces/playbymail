package tag

import (
	"database/sql"
	"testing"
	"time"

	"github.com/lib/pq"
	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/playbymail/core/collection/set"
)

// Record is reproduced here instead of using repository.Record to prevent an import cycle
type Record struct {
	ID        string       `db:"id"`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at"`
	DeletedAt sql.NullTime `db:"deleted_at"`
}

type Country struct {
	IsSupportedCountry         bool           `db:"is_supported_country"`
	IsSupportedContractCountry bool           `db:"is_supported_contract_country"`
	SupportedLocales           pq.StringArray `db:"supported_locales"`
	Currencies                 pq.StringArray `db:"currencies"`
	SupportedCurrencies        pq.StringArray `db:"supported_currencies"`
	SystemCurrency             sql.NullString `db:"system_currency"`
	ByteSlice                  []byte         `db:"byte_slice"`
	Record
	Nested struct {
		NestedSupportedLocales    pq.StringArray `db:"nested_supported_locales"`
		NestedCurrencies          pq.StringArray `db:"nested_currencies"`
		NestedSupportedCurrencies pq.StringArray `db:"nested_supported_currencies"`
		privateNested             struct {
			privateNestedSupportedLocales    pq.StringArray `db:"private_nested_supported_locales"`
			privateNestedCurrencies          pq.StringArray `db:"private_nested_currencies"`
			privateNestedSupportedCurrencies pq.StringArray `db:"private_nested_supported_currencies"`
		}
	}
}

func TestGetArrayFieldNamesFromStruct(t *testing.T) {
	tests := []struct {
		name   string
		entity any
		want   set.Set[string]
	}{
		{
			name:   "country record",
			entity: Country{},
			want: set.New(
				"supported_locales",
				"currencies",
				"supported_currencies",
				"nested_supported_locales",
				"nested_currencies",
				"nested_supported_currencies",
				"private_nested_supported_locales",
				"private_nested_currencies",
				"private_nested_supported_currencies",
			),
		},
		{
			name:   "empty struct",
			entity: struct{}{},
			want:   set.Set[string]{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := GetArrayFieldTagValues(tt.entity, "db")
			for a := range actual {
				require.Contains(t, tt.want, a, "struct should have expected array field names")
			}
		})
	}
}
