package repositor

import (
	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/collection/set"
)

// RLSConstraint defines a row-level security constraint that is automatically
// applied when a repository record has a matching column name. Constraints use
// SQL templates with named args that are resolved from RLS identifiers.
type RLSConstraint struct {
	Column                 string   // Column name that triggers this constraint (e.g., "game_instance_id", "game_id")
	SQLTemplate            string   // SQL fragment using named args (e.g., "IN (SELECT game_instance_id FROM game_subscription WHERE account_id = :account_id AND status = 'active')")
	RequiredRLSIdentifiers []string // Required RLS identifiers (e.g., ["account_id"]) that must be present in RLS identifiers for the constraint to be applied
	SkipSelfMapping        bool     // When true, prevents domain from mapping this constraint to the table's own "id" column (use when SQLTemplate references Column as a literal column name)
}

// Repositor -
type Repositor interface {
	TableName() string
	Attributes() []string
	ArrayFields() set.Set[string]
	Tx() pgx.Tx
	SetRLSIdentifiers(identifiers map[string][]string)
	SetRLSConstraints(constraints []RLSConstraint)
}
