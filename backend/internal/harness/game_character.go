package harness

import (
	"fmt"

	"github.com/brianvoe/gofakeit"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/record"
)

func (t *Testing) createGameCharacterRec(charConfig GameCharacterConfig, gameRec *record.Game) (*record.GameCharacter, error) {
	l := t.Logger("createGameCharacterRec")

	if gameRec == nil {
		return nil, fmt.Errorf("game record is nil for game_character record >%#v<", charConfig)
	}

	if charConfig.AccountRef == "" {
		return nil, fmt.Errorf("game_character record >%#v< must have an AccountRef set", charConfig)
	}

	var rec *record.GameCharacter
	if charConfig.Record != nil {
		recCopy := *charConfig.Record
		rec = &recCopy
	} else {
		rec = &record.GameCharacter{}
	}

	rec = t.applyGameCharacterRecDefaultValues(rec)

	rec.GameID = gameRec.ID

	// Get account record
	accountRec, err := t.Data.GetAccountRecByRef(charConfig.AccountRef)
	if err != nil {
		l.Warn("failed resolving account ref >%s<: %v", charConfig.AccountRef, err)
		return nil, err
	}
	rec.AccountID = accountRec.ID

	// Create record
	l.Info("creating game_character record >%#v<", rec)

	rec, err = t.Domain.(*domain.Domain).CreateGameCharacterRec(rec)
	if err != nil {
		l.Warn("failed creating game_character record >%v<", err)
		return nil, err
	}

	// Add to data store
	t.Data.AddGameCharacterRec(rec)

	// Add to teardown data store
	t.teardownData.AddGameCharacterRec(rec)

	// Add to references store
	if charConfig.Reference != "" {
		t.Data.Refs.GameCharacterRefs[charConfig.Reference] = rec.ID
	}

	return rec, nil
}

func (t *Testing) applyGameCharacterRecDefaultValues(rec *record.GameCharacter) *record.GameCharacter {
	if rec == nil {
		rec = &record.GameCharacter{}
	}
	if rec.Name == "" {
		rec.Name = UniqueName(gofakeit.Name())
	}
	return rec
}
