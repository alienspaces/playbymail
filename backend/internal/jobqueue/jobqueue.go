package jobqueue

import (
	"github.com/riverqueue/river"
	"gitlab.com/alienspaces/playbymail/core/collection/set"
)

const (
	QueueDefault string = river.QueueDefault
)

var Queues set.Set[string] = set.New(
	QueueDefault,
)
