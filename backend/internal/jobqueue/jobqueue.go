package jobqueue

import (
	"github.com/riverqueue/river"
	"gitlab.com/alienspaces/playbymail/core/collection/set"
)

const (
	QueueDefault string = river.QueueDefault
	QueueGame    string = "game"
)

var Queues set.Set[string] = set.New(
	QueueDefault,
	QueueGame,
)
