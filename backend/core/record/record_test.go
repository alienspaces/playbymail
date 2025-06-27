package record

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewRecordID(t *testing.T) {
	ids := map[string]struct{}{}

	for range make([]int, 10000) {
		id := NewRecordID()

		_, ok := ids[id]
		ids[id] = struct{}{}

		require.False(t, ok, "NewRecordID should be unique")

		uuidLength := 36
		require.NotEmpty(t, id, "NewRecordID should not be empty")
		require.Equal(t, uuidLength, len(id), "NewRecordID has expected length")
	}
}

func Test_RFC3339Microseconds(t *testing.T) {
	tests := []struct {
		name string
		time func(*testing.T) time.Time
		want string
	}{
		{
			name: "positive offset",
			time: func(t *testing.T) time.Time {
				parse, err := time.Parse(RFC3339Microseconds, "2006-01-02T15:04:05.000000+08:00")
				require.NoError(t, err, "time.Parse should not error")
				return parse
			},
			want: "2006-01-02T15:04:05.000000+08:00",
		},
		{
			name: "utc",
			time: func(t *testing.T) time.Time {
				parse, err := time.Parse(RFC3339Microseconds, "2006-01-02T15:04:05.123456Z")
				require.NoError(t, err, "time.Parse should not error")
				return parse
			},
			want: "2006-01-02T15:04:05.123456Z",
		},
		{
			name: "negative offset",
			time: func(t *testing.T) time.Time {
				parse, err := time.Parse(RFC3339Microseconds, "2024-09-02T09:34:45.725187-04:00")
				require.NoError(t, err, "time.Parse should not error")
				return parse
			},
			want: "2024-09-02T09:34:45.725187-04:00",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, tt.time(t).Format(RFC3339Microseconds))
		})
	}
}
