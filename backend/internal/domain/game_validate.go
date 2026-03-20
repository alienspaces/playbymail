package domain

import (
	"fmt"

	"gitlab.com/alienspaces/playbymail/core/domain"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
)

const (
	ValidationSeverityError   = "error"
	ValidationSeverityWarning = "warning"
)

type GameValidationIssue struct {
	Field    string `json:"field"`
	Message  string `json:"message"`
	Severity string `json:"severity"`
}

// ValidateGameReadyForInstance checks whether a game is ready to create an
// instance. It collects all issues rather than stopping at the first, so that
// game designers can see everything that needs fixing at once.
func (m *Domain) ValidateGameReadyForInstance(gameID string) ([]GameValidationIssue, error) {
	l := m.Logger("ValidateGameReadyForInstance")

	gameRec, err := m.GetGameRec(gameID, nil)
	if err != nil {
		return nil, err
	}

	var issues []GameValidationIssue

	if gameRec.GameType == game_record.GameTypeAdventure {
		issues, err = m.validateAdventureGameReadyForInstance(gameID)
		if err != nil {
			l.Warn("failed validating adventure game >%s< >%v<", gameID, err)
			return nil, err
		}
	}

	return issues, nil
}

func (m *Domain) validateAdventureGameReadyForInstance(gameID string) ([]GameValidationIssue, error) {
	var issues []GameValidationIssue

	locationRecs, err := m.GetManyAdventureGameLocationRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: adventure_game_record.FieldAdventureGameLocationGameID, Val: gameID},
		},
		Limit: 1,
	})
	if err != nil {
		return nil, err
	}

	if len(locationRecs) == 0 {
		issues = append(issues, GameValidationIssue{
			Field:    "locations",
			Message:  "Adventure game must have at least one location",
			Severity: ValidationSeverityError,
		})
		return issues, nil
	}

	startingLocationRecs, err := m.GetManyAdventureGameLocationRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: adventure_game_record.FieldAdventureGameLocationGameID, Val: gameID},
			{Col: adventure_game_record.FieldAdventureGameLocationIsStartingLocation, Val: true},
		},
		Limit: 1,
	})
	if err != nil {
		return nil, err
	}

	if len(startingLocationRecs) == 0 {
		issues = append(issues, GameValidationIssue{
			Field:    "starting_location",
			Message:  "Adventure game must have at least one starting location before creating an instance",
			Severity: ValidationSeverityError,
		})
	}

	stateIssues, err := m.validateAdventureGameObjectStateGraph(gameID)
	if err != nil {
		return nil, err
	}
	issues = append(issues, stateIssues...)

	return issues, nil
}

// validateAdventureGameObjectStateGraph checks the state transition graph for each
// location object and emits warnings for designer-configuration problems:
//   - Objects with states defined but no initial state set
//   - States that can never be reached (not the initial state, not a result of any effect)
//   - States that are dead-ends (no effects require this state, so once reached nothing changes)
func (m *Domain) validateAdventureGameObjectStateGraph(gameID string) ([]GameValidationIssue, error) {
	var issues []GameValidationIssue

	// Load all objects for this game.
	objectRecs, err := m.GetManyAdventureGameLocationObjectRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: adventure_game_record.FieldAdventureGameLocationObjectGameID, Val: gameID},
		},
	})
	if err != nil {
		return nil, err
	}
	if len(objectRecs) == 0 {
		return nil, nil
	}

	// Load all states for this game.
	allStateRecs, err := m.GetManyAdventureGameLocationObjectStateRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: adventure_game_record.FieldAdventureGameLocationObjectStateGameID, Val: gameID},
		},
	})
	if err != nil {
		return nil, err
	}

	// Load all effects for this game.
	allEffectRecs, err := m.GetManyAdventureGameLocationObjectEffectRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: adventure_game_record.FieldAdventureGameLocationObjectEffectGameID, Val: gameID},
		},
	})
	if err != nil {
		return nil, err
	}

	// Index states and effects by object ID for efficient lookup.
	statesByObjectID := make(map[string][]*adventure_game_record.AdventureGameLocationObjectState)
	for _, s := range allStateRecs {
		statesByObjectID[s.AdventureGameLocationObjectID] = append(statesByObjectID[s.AdventureGameLocationObjectID], s)
	}

	effectsByObjectID := make(map[string][]*adventure_game_record.AdventureGameLocationObjectEffect)
	for _, e := range allEffectRecs {
		effectsByObjectID[e.AdventureGameLocationObjectID] = append(effectsByObjectID[e.AdventureGameLocationObjectID], e)
	}

	for _, obj := range objectRecs {
		states := statesByObjectID[obj.ID]
		if len(states) == 0 {
			continue
		}

		effects := effectsByObjectID[obj.ID]

		// Check: object has states but no initial state.
		if !obj.InitialAdventureGameLocationObjectStateID.Valid {
			issues = append(issues, GameValidationIssue{
				Field:    "location_objects",
				Message:  fmt.Sprintf("Object %q has states defined but no initial state is set", obj.Name),
				Severity: ValidationSeverityWarning,
			})
		}

		// Build sets of IDs referenced in effects.
		reachableStateIDs := make(map[string]bool) // result of any effect
		requiredStateIDs := make(map[string]bool)   // required by any effect
		if obj.InitialAdventureGameLocationObjectStateID.Valid {
			reachableStateIDs[obj.InitialAdventureGameLocationObjectStateID.String] = true
		}
		for _, e := range effects {
			if e.ResultAdventureGameLocationObjectStateID.Valid {
				reachableStateIDs[e.ResultAdventureGameLocationObjectStateID.String] = true
			}
			if e.RequiredAdventureGameLocationObjectStateID.Valid {
				requiredStateIDs[e.RequiredAdventureGameLocationObjectStateID.String] = true
			}
		}

		for _, s := range states {
			// Warn about unreachable states (not initial, not a result of any effect).
			if !reachableStateIDs[s.ID] {
				issues = append(issues, GameValidationIssue{
					Field:    "location_objects",
					Message:  fmt.Sprintf("Object %q state %q is unreachable: it is not the initial state and no effect produces it", obj.Name, s.Name),
					Severity: ValidationSeverityWarning,
				})
			}

			// Warn about dead-end states (no effects require this state, meaning
			// once an object reaches this state no further effects can be triggered).
			// Only warn if there are effects at all and at least one has a required state,
			// so we don't spam on simple objects with unconditional effects.
			if len(requiredStateIDs) > 0 && !requiredStateIDs[s.ID] {
				issues = append(issues, GameValidationIssue{
					Field:    "location_objects",
					Message:  fmt.Sprintf("Object %q state %q is a dead-end: no effects require it, so players cannot interact further once the object reaches this state", obj.Name, s.Name),
					Severity: ValidationSeverityWarning,
				})
			}
		}
	}

	return issues, nil
}

type validateGameArgs struct {
	nextRec *game_record.Game
	currRec *game_record.Game
}

func (m *Domain) populateGameValidateArgs(currRec, nextRec *game_record.Game) (*validateGameArgs, error) {
	args := &validateGameArgs{
		currRec: currRec,
		nextRec: nextRec,
	}

	// Note: account_id validation and subscription limit checks are now handled
	// at the handler level where we have access to authenticated account information.
	// Games no longer have account_id field, so we can't validate it here.

	return args, nil
}

func (m *Domain) validateGameRecForCreate(rec *game_record.Game) error {
	args, err := m.populateGameValidateArgs(nil, rec)
	if err != nil {
		return err
	}
	return validateGameRecForCreate(args)
}

func (m *Domain) validateGameRecForUpdate(currRec, nextRec *game_record.Game) error {
	args, err := m.populateGameValidateArgs(currRec, nextRec)
	if err != nil {
		return err
	}
	return validateGameRecForUpdate(args)
}

func (m *Domain) validateGameRecForDelete(rec *game_record.Game) error {
	args, err := m.populateGameValidateArgs(nil, rec)
	if err != nil {
		return err
	}
	return validateGameRecForDelete(args)
}

func validateGameRecForCreate(args *validateGameArgs) error {
	rec := args.nextRec

	if rec.Status == "" {
		return InvalidField(game_record.FieldGameStatus, "", "status is required")
	}

	// Note: Subscription limit validation has been moved to handler level
	// where we have access to authenticated account information.
	// Games no longer have account_id field.

	return validateGameRec(args)
}

func validateGameRecForUpdate(args *validateGameArgs) error {
	currRec := args.currRec
	nextRec := args.nextRec

	// Prevent status from being changed back to draft once published
	if currRec.Status == game_record.GameStatusPublished && nextRec.Status != game_record.GameStatusPublished {
		return InvalidField(game_record.FieldGameStatus, currRec.Status, "published game status cannot be changed")
	}

	// Validate status transitions
	// Only allow: draft -> published, draft -> draft
	if currRec.Status == game_record.GameStatusDraft {
		if nextRec.Status != game_record.GameStatusDraft && nextRec.Status != game_record.GameStatusPublished {
			return InvalidField(game_record.FieldGameStatus, nextRec.Status, "invalid status transition from draft")
		}
	}

	return validateGameRec(args)
}

func validateGameRec(args *validateGameArgs) error {
	rec := args.nextRec

	// Note: account_id validation removed - games no longer have account_id field

	if err := domain.ValidateStringField(game_record.FieldGameName, rec.Name); err != nil {
		return err
	}

	if rec.GameType != game_record.GameTypeAdventure {
		return InvalidField(game_record.FieldGameType, rec.GameType, "game type is not valid")
	}

	if rec.TurnDurationHours <= 0 {
		return InvalidField(game_record.FieldGameTurnDurationHours, fmt.Sprintf("%d", rec.TurnDurationHours), "turn duration hours must be greater than 0")
	}

	if err := domain.ValidateStringField(game_record.FieldGameDescription, rec.Description); err != nil {
		return err
	}

	if rec.Status != game_record.GameStatusDraft && rec.Status != game_record.GameStatusPublished {
		return InvalidField(game_record.FieldGameStatus, rec.Status, "status is not valid")
	}

	return nil
}

func validateGameRecForDelete(args *validateGameArgs) error {
	rec := args.nextRec

	if err := domain.ValidateUUIDField(game_record.FieldGameID, rec.ID); err != nil {
		return err
	}

	return nil
}
