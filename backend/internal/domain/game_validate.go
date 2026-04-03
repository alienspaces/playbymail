package domain

import (
	"fmt"

	"gitlab.com/alienspaces/playbymail/core/domain"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
	"gitlab.com/alienspaces/playbymail/internal/record/mecha_record"
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

	switch gameRec.GameType {
	case game_record.GameTypeAdventure:
		issues, err = m.validateAdventureGameReadyForInstance(gameID)
		if err != nil {
			l.Warn("failed validating adventure game >%s< >%v<", gameID, err)
			return nil, err
		}
	case game_record.GameTypeMecha:
		issues, err = m.validateMechaGameReadyForInstance(gameID)
		if err != nil {
			l.Warn("failed validating mecha game >%s< >%v<", gameID, err)
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
//
// State transitions can originate from two sources that must both be considered:
//   - The object's own effects (e.g. a door being pushed open)
//   - Cross-object effects from other objects (e.g. a lever that reveals a hidden door)
//
// Ignoring cross-object effects produces false-positive warnings on objects like the
// Hidden Door, whose states are driven entirely by the Lever in another location.
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

	// Index states by the object they belong to.
	statesByObjectID := make(map[string][]*adventure_game_record.AdventureGameLocationObjectState)
	for _, s := range allStateRecs {
		statesByObjectID[s.AdventureGameLocationObjectID] = append(statesByObjectID[s.AdventureGameLocationObjectID], s)
	}

	// effectsByObjectID holds effects that an object owns — these are the interactions
	// available to players when they act on that object (e.g. push, inspect, burn).
	effectsByObjectID := make(map[string][]*adventure_game_record.AdventureGameLocationObjectEffect)
	for _, e := range allEffectRecs {
		effectsByObjectID[e.AdventureGameLocationObjectID] = append(effectsByObjectID[e.AdventureGameLocationObjectID], e)
	}

	// crossEffectsByTargetID holds effects whose result targets another object via
	// result_adventure_game_location_object_id. These are cross-object effects such as
	// change_object_state, reveal_object, and hide_object. An object may have states
	// that are only reachable via these external transitions (e.g. a door revealed by a
	// lever), so they must be considered separately from the object's own effects.
	crossEffectsByTargetID := make(map[string][]*adventure_game_record.AdventureGameLocationObjectEffect)
	for _, e := range allEffectRecs {
		if e.ResultAdventureGameLocationObjectID.Valid {
			targetID := e.ResultAdventureGameLocationObjectID.String
			crossEffectsByTargetID[targetID] = append(crossEffectsByTargetID[targetID], e)
		}
	}

	for _, obj := range objectRecs {
		states := statesByObjectID[obj.ID]
		if len(states) == 0 {
			continue
		}

		ownEffects := effectsByObjectID[obj.ID]
		crossEffects := crossEffectsByTargetID[obj.ID]

		// Check: object has states but no initial state.
		if !obj.InitialAdventureGameLocationObjectStateID.Valid {
			issues = append(issues, GameValidationIssue{
				Field:    "location_objects",
				Message:  fmt.Sprintf("Object %q has states defined but no initial state is set", obj.Name),
				Severity: ValidationSeverityWarning,
			})
		}

		// Build sets of state IDs that are reachable (can be entered) or required
		// (must be active for an effect to fire) based on the object's own effects.
		reachableStateIDs := make(map[string]bool)
		requiredStateIDs := make(map[string]bool)
		if obj.InitialAdventureGameLocationObjectStateID.Valid {
			reachableStateIDs[obj.InitialAdventureGameLocationObjectStateID.String] = true
		}
		for _, e := range ownEffects {
			if e.ResultAdventureGameLocationObjectStateID.Valid {
				reachableStateIDs[e.ResultAdventureGameLocationObjectStateID.String] = true
			}
			if e.RequiredAdventureGameLocationObjectStateID.Valid {
				requiredStateIDs[e.RequiredAdventureGameLocationObjectStateID.String] = true
			}
		}

		// Cross-object effects from other objects can also set this object's state
		// (e.g. a lever using change_object_state to reveal a hidden door). Include
		// those result states in the reachable set so we don't falsely warn that
		// states like "revealed" are unreachable just because no own-effect produces them.
		for _, e := range crossEffects {
			if e.ResultAdventureGameLocationObjectStateID.Valid {
				reachableStateIDs[e.ResultAdventureGameLocationObjectStateID.String] = true
			}
		}

		// Objects with a remove_object effect are intentionally destructible: when the
		// object is consumed (e.g. a burning rack), the terminal state before removal is
		// by design and should not be flagged as a dead-end. Detect this once per object.
		isDestructible := false
		for _, e := range ownEffects {
			if e.EffectType == adventure_game_record.AdventureGameLocationObjectEffectEffectTypeRemoveObject {
				isDestructible = true
				break
			}
		}

		// If cross-object effects target this object, dead-end suppression applies.
		// External effects (e.g. hide_object, change_object_state) can unconditionally
		// transition this object regardless of its current state, so a state with no
		// own required-state effects is not necessarily a designer mistake.
		hasInboundCrossEffects := len(crossEffects) > 0

		for _, s := range states {
			// Warn about unreachable states (not initial, not a result of any own or
			// cross-object effect). This still fires even for objects with cross-object
			// effects, because an unreachable state is always a configuration mistake.
			if !reachableStateIDs[s.ID] {
				issues = append(issues, GameValidationIssue{
					Field:    "location_objects",
					Message:  fmt.Sprintf("Object %q state %q is unreachable: it is not the initial state and no effect produces it", obj.Name, s.Name),
					Severity: ValidationSeverityWarning,
				})
			}

			// Warn about dead-end states: no own effect requires this state, meaning
			// once the object reaches this state players cannot interact further.
			// Only warn when:
			//   - there are effects with required states (otherwise all states look like dead-ends), and
			//   - no cross-object effects target this object (external forces can still move it), and
			//   - the object is not destructible (terminal state before removal is intentional).
			if len(requiredStateIDs) > 0 && !requiredStateIDs[s.ID] && !hasInboundCrossEffects && !isDestructible {
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

	if rec.GameType != game_record.GameTypeAdventure && rec.GameType != game_record.GameTypeMecha {
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

func (m *Domain) validateMechaGameReadyForInstance(gameID string) ([]GameValidationIssue, error) {
	var issues []GameValidationIssue

	sectorRecs, err := m.GetManyMechaSectorRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: mecha_record.FieldMechaSectorGameID, Val: gameID},
		},
		Limit: 1,
	})
	if err != nil {
		return nil, err
	}

	if len(sectorRecs) == 0 {
		issues = append(issues, GameValidationIssue{
			Field:    "sectors",
			Message:  "Mecha must have at least one sector before creating an instance",
			Severity: ValidationSeverityError,
		})
		return issues, nil
	}

	startingSectorRecs, err := m.GetManyMechaSectorRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: mecha_record.FieldMechaSectorGameID, Val: gameID},
			{Col: mecha_record.FieldMechaSectorIsStartingSector, Val: true},
		},
		Limit: 1,
	})
	if err != nil {
		return nil, err
	}

	if len(startingSectorRecs) == 0 {
		issues = append(issues, GameValidationIssue{
			Field:    "starting_sector",
			Message:  "Mecha must have at least one starting sector before creating an instance",
			Severity: ValidationSeverityError,
		})
	}

	chassisRecs, err := m.GetManyMechaChassisRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: mecha_record.FieldMechaChassisGameID, Val: gameID},
		},
		Limit: 1,
	})
	if err != nil {
		return nil, err
	}

	if len(chassisRecs) == 0 {
		issues = append(issues, GameValidationIssue{
			Field:    "chassis",
			Message:  "Mecha must have at least one chassis defined before creating an instance",
			Severity: ValidationSeverityError,
		})
	}

	// Require exactly one player starter lance with at least one mech.
	starterLanceRecs, err := m.GetManyMechaLanceRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: mecha_record.FieldMechaLanceGameID, Val: gameID},
			{Col: mecha_record.FieldMechaLanceLanceType, Val: mecha_record.LanceTypeStarter},
		},
		Limit: 1,
	})
	if err != nil {
		return nil, err
	}

	if len(starterLanceRecs) == 0 {
		issues = append(issues, GameValidationIssue{
			Field:    "player_starter_lance",
			Message:  "Mecha game must have a player starter lance defined",
			Severity: ValidationSeverityError,
		})
	} else {
		starterMechs, err := m.GetManyMechaLanceMechRecs(&coresql.Options{
			Params: []coresql.Param{
				{Col: mecha_record.FieldMechaLanceMechMechaLanceID, Val: starterLanceRecs[0].ID},
			},
			Limit: 1,
		})
		if err != nil {
			return nil, err
		}
		if len(starterMechs) == 0 {
			issues = append(issues, GameValidationIssue{
				Field:    "player_starter_lance",
				Message:  "Player starter lance must have at least one mech",
				Severity: ValidationSeverityError,
			})
		}
	}

	return issues, nil
}
