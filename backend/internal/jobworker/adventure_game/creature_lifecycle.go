package adventure_game

import (
	"context"
	"database/sql"
	"fmt"

	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
	"gitlab.com/alienspaces/playbymail/internal/turnsheet"
)


// RunCreatureLifecycle processes body decay and creature respawn for all game instances.
// This should run BEFORE CreateTurnSheets so that dead bodies that have decayed
// are removed from the game state and respawned creatures appear for the next turn.
func (p *AdventureGame) RunCreatureLifecycle(ctx context.Context, gameInstanceRec *game_record.GameInstance) error {
	l := p.Logger.WithFunctionContext("AdventureGame/RunCreatureLifecycle")
	l.Info("running creature lifecycle for game instance >%s< turn >%d<", gameInstanceRec.ID, gameInstanceRec.CurrentTurn)

	creatureInstances, err := p.Domain.GetManyAdventureGameCreatureInstanceRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: adventure_game_record.FieldAdventureGameCreatureInstanceGameInstanceID, Val: gameInstanceRec.ID},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to get creature instances: %w", err)
	}

	for _, ci := range creatureInstances {
		if ci.Health > 0 || !ci.DiedAtTurn.Valid {
			continue
		}

		creatureDef, err := p.Domain.GetAdventureGameCreatureRec(ci.AdventureGameCreatureID, nil)
		if err != nil {
			l.Warn("failed to get creature definition >%s< >%v<", ci.AdventureGameCreatureID, err)
			continue
		}

		turnsDeadFor := int64(gameInstanceRec.CurrentTurn) - ci.DiedAtTurn.Int64

		// Respawn takes priority over decay.
		if creatureDef.RespawnTurns > 0 && turnsDeadFor >= int64(creatureDef.RespawnTurns) {
			if err := p.respawnCreature(ctx, gameInstanceRec, ci, creatureDef); err != nil {
				l.Warn("failed to respawn creature >%s< >%v<", ci.ID, err)
			}
			continue
		}

		// Decay: soft-delete the instance if body has expired.
		if turnsDeadFor > int64(creatureDef.BodyDecayTurns) {
			if err := p.Domain.DeleteAdventureGameCreatureInstanceRec(ci.ID); err != nil {
				l.Warn("failed to delete decayed creature instance >%s< >%v<", ci.ID, err)
			} else {
				l.Info("decayed creature instance >%s< (%s) after %d turns", ci.ID, creatureDef.Name, turnsDeadFor)
			}
		}
	}

	return nil
}

// respawnCreature restores a dead creature instance to full health and generates
// a world event for players at that location.
func (p *AdventureGame) respawnCreature(_ context.Context, gameInstanceRec *game_record.GameInstance, ci *adventure_game_record.AdventureGameCreatureInstance, creatureDef *adventure_game_record.AdventureGameCreature) error {
	l := p.Logger.WithFunctionContext("AdventureGame/respawnCreature")

	ci.Health = creatureDef.MaxHealth
	ci.DiedAtTurn = sql.NullInt64{}
	if _, err := p.Domain.UpdateAdventureGameCreatureInstanceRec(ci); err != nil {
		return fmt.Errorf("failed to update respawned creature instance: %w", err)
	}

	l.Info("respawned creature >%s< (%s) at location instance >%s<", ci.ID, creatureDef.Name, ci.AdventureGameLocationInstanceID)

	// Generate world event for all characters at this location.
	characterInstances, err := p.Domain.GetManyAdventureGameCharacterInstanceRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: adventure_game_record.FieldAdventureGameCharacterInstanceGameInstanceID, Val: gameInstanceRec.ID},
			{Col: adventure_game_record.FieldAdventureGameCharacterInstanceAdventureGameLocationInstanceID, Val: ci.AdventureGameLocationInstanceID},
		},
	})
	if err != nil {
		l.Warn("failed to get character instances for respawn world event >%v<", err)
		return nil
	}

	for _, charInst := range characterInstances {
		event := turnsheet.TurnEvent{
			Category: turnsheet.TurnEventCategoryWorld,
			Icon:     turnsheet.TurnEventIconWorld,
			Message:  fmt.Sprintf("A %s emerges from the shadows.", creatureDef.Name),
		}
		if err := turnsheet.AppendTurnEvent(charInst, event); err != nil {
			l.Warn("failed to append respawn world event for character >%s< >%v<", charInst.ID, err)
			continue
		}
		if _, err := p.Domain.UpdateAdventureGameCharacterInstanceRec(charInst); err != nil {
			l.Warn("failed to save respawn world event for character >%s< >%v<", charInst.ID, err)
		}
	}

	return nil
}
