package harness

import (
	"fmt"

	"github.com/brianvoe/gofakeit"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
)

func (t *Testing) createAdventureGameCharacterRec(charConfig AdventureGameCharacterConfig, gameRec *game_record.Game) (*adventure_game_record.AdventureGameCharacter, error) {
	l := t.Logger("createAdventureGameCharacterRec")

	if gameRec == nil {
		return nil, fmt.Errorf("game record is nil for adventure game character record >%#v<", charConfig)
	}

	if charConfig.AccountRef == "" {
		return nil, fmt.Errorf("adventure game character record >%#v< must have an AccountRef set", charConfig)
	}

	var rec *adventure_game_record.AdventureGameCharacter
	if charConfig.Record != nil {
		recCopy := *charConfig.Record
		rec = &recCopy
	} else {
		rec = &adventure_game_record.AdventureGameCharacter{}
	}

	rec = t.applyAdventureGameCharacterRecDefaultValues(rec)

	rec.GameID = gameRec.ID

	// Get account record
	accountRec, err := t.Data.GetAccountRecByRef(charConfig.AccountRef)
	if err != nil {
		l.Warn("failed resolving account ref >%s<: %v", charConfig.AccountRef, err)
		return nil, err
	}
	rec.AccountID = accountRec.ID

	// Create record
	l.Debug("creating game_character record >%#v<", rec)

	rec, err = t.Domain.(*domain.Domain).CreateAdventureGameCharacterRec(rec)
	if err != nil {
		l.Warn("failed creating game_character record >%v<", err)
		return nil, err
	}

	// Add to data store
	t.Data.AddAdventureGameCharacterRec(rec)

	// Add to teardown data store
	t.teardownData.AddAdventureGameCharacterRec(rec)

	// Add to references store
	if charConfig.Reference != "" {
		t.Data.Refs.AdventureGameCharacterRefs[charConfig.Reference] = rec.ID
	}

	return rec, nil
}

func (t *Testing) applyAdventureGameCharacterRecDefaultValues(rec *adventure_game_record.AdventureGameCharacter) *adventure_game_record.AdventureGameCharacter {
	if rec == nil {
		rec = &adventure_game_record.AdventureGameCharacter{}
	}
	if rec.Name == "" {
		rec.Name = UniqueName(gofakeit.Name())
	}
	return rec
}
