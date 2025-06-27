package runnable

import (
	"gitlab.com/alienspaces/playbymail/core/type/storer"
)

// Runnable -
type Runnable interface {
	Init(s storer.Storer) error
	Run(args map[string]any) error
}

// StatelessRunnable -
type StatelessRunnable interface {
	Init() error
	Run(args map[string]any) error
}
