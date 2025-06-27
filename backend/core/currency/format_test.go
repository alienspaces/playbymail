package currency

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLowestDenominationMonetaryUnitToFormatted(t *testing.T) {
	t.Parallel()

	type args struct {
		value    string
		currency string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "AUD",
			args: args{
				currency: AUD,
			},
			want: "$1,234,567.89",
		},
		{
			name: "CAD",
			args: args{
				currency: CAD,
			},
			want: "$1,234,567.89",
		},
		{
			name: "CHF",
			args: args{
				currency: CHF,
			},
			want: "CHF 1'234'567.89",
		},
		{
			name: "CNY",
			args: args{
				currency: CNY,
			},
			want: "¥1,234,567.89",
		},
		{
			name: "DKK",
			args: args{
				currency: DKK,
			},
			want: "1.234.567,89 kr",
		},
		{
			name: "EUR",
			args: args{
				currency: EUR,
			},
			want: "€1.234.567,89",
		},
		{
			name: "GBP",
			args: args{
				currency: GBP,
			},
			want: "£1,234,567.89",
		},
		{
			name: "JPY",
			args: args{
				currency: JPY,
			},
			want: "¥123,456,789",
		},
		{
			name: "KRW",
			args: args{
				currency: KRW,
			},
			want: "₩123,456,789",
		},
		{
			name: "MXN",
			args: args{
				currency: MXN,
			},
			want: "$1,234,567.89",
		},
		{
			name: "NOK",
			args: args{
				currency: NOK,
			},
			want: "1 234 567,89 kr",
		},
		{
			name: "NZD",
			args: args{
				currency: NZD,
			},
			want: "$1,234,567.89",
		},
		{
			name: "PLN",
			args: args{
				currency: PLN,
			},
			want: "1 234 567,89 zł",
		},
		{
			name: "SEK",
			args: args{
				currency: SEK,
			},
			want: "1 234 567,89 kr",
		},
		{
			name: "SGD",
			args: args{
				currency: SGD,
			},
			want: "$1,234,567.89",
		},
		{
			name: "TWD",
			args: args{
				currency: TWD,
			},
			want: "NT$123,456,789",
		},
		{
			name: "USD",
			args: args{
				currency: USD,
			},
			want: "$1,234,567.89",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			v := "123456789"
			if tt.args.value != "" {
				v = tt.args.value
			}

			got, err := LowestDenominationMonetaryUnitToFormatted(v, tt.args.currency)
			require.NoError(t, err)

			require.Equal(t, tt.want, got)
		})
	}
}
