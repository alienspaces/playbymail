package repository

import (
	"fmt"

	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/collection/set"
	"gitlab.com/alienspaces/playbymail/core/collection/slice"
	"gitlab.com/alienspaces/playbymail/core/tag"
)

type TestRecorder[Rec any] interface {
	*Rec
	ToNamedArgs() pgx.NamedArgs
}

// TestToNamedArgs is a test helper for core and service records to test
// a records ToNamedArgs method implementation.
func TestToNamedArgs[Rec any, RecPtr TestRecorder[Rec]](rec RecPtr) error {
	args := rec.ToNamedArgs()
	namedKeys := tag.GetFieldTagValues(*rec, "db")
	if len(namedKeys) == 0 {
		return fmt.Errorf("failed getting expected args >%#v< rec >%#v<", namedKeys, rec)
	}

	expected := set.New(slice.FromMapKeys(args)...)
	actual := set.New(namedKeys...)
	diff := set.SymmetricDifference(expected, actual)
	if len(diff) > 0 {
		return fmt.Errorf("named keys >%v< not found in args >%v<", diff, args)
	}

	// sanity check on GetFieldTagValues not returning every key; and
	// ensure there are no duplicate named args
	if len(args) != len(namedKeys) {
		return fmt.Errorf("length of named args >%d< does not equal expected length >%d<",
			len(args),
			len(namedKeys),
		)
	}

	return nil
}
