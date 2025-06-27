package fake

import (
	"database/sql"
	"time"

	"gitlab.com/alienspaces/playbymail/core/nulltime"
)

func NowNullTime() sql.NullTime {
	return nulltime.FromTime(time.Now())
}
