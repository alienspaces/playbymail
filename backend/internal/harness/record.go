package harness

import (
	"fmt"

	"github.com/brianvoe/gofakeit"

	corerecord "gitlab.com/alienspaces/playbymail/core/record"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/record"
)

// UniqueName appends a UUID4 to the end of the name to make it unique
// for parallel test execution.
func UniqueName(name string) string {
	if name == "" {
		name = gofakeit.Color()
	}
	return fmt.Sprintf("%s (%s)", name, corerecord.NewRecordID())
}

// NormalName removes the unique UUID4 from the end of the name to make it normal for
// test harness functions that return a record based on its non unique name.
func NormalName(name string) string {
	return name[:len(name)-39]
}

// UniqueEmail appends a UUID4 to the end of the email to make it unique
// for parallel test execution.
func UniqueEmail(email string) string {
	return fmt.Sprintf("%s-%s", email, corerecord.NewRecordID())
}

// Account
func (t *Testing) createAccountRec(accountConfig AccountConfig) (*record.Account, error) {
	l := t.Logger("createAccountRec")

	// Create a new record instance to avoid reusing the same record across tests
	var rec *record.Account
	if accountConfig.Record != nil {
		// Copy the record to avoid modifying the original
		recCopy := *accountConfig.Record
		rec = &recCopy
	} else {
		rec = &record.Account{}
	}

	rec = t.applyAccountRecDefaultValues(rec)

	l.Info("creating account record >%#v<", rec)

	rec, err := t.Domain.(*domain.Domain).CreateAccountRec(rec)
	if err != nil {
		l.Warn("failed creating account record >%v<", err)
		return nil, err
	}

	// Add the account record to the data store
	t.Data.AddAccountRec(rec)

	// Add the account record to the teardown data store
	t.teardownData.AddAccountRec(rec)

	// Add the account record to the data store refs
	if accountConfig.Reference != "" {
		t.Data.Refs.AccountRefs[accountConfig.Reference] = rec.ID
	}

	return rec, nil
}

func (t *Testing) applyAccountRecDefaultValues(rec *record.Account) *record.Account {
	if rec == nil {
		rec = &record.Account{}
	}

	if rec.Email == "" {
		rec.Email = UniqueEmail(gofakeit.Email())
	}

	if rec.Name == "" {
		rec.Name = UniqueName(gofakeit.Name())
	}

	return rec
}

// Game
func (t *Testing) createGameRec(gameConfig GameConfig) (*record.Game, error) {
	l := t.Logger("createGameRec")

	// Create a new record instance to avoid reusing the same record across tests
	var rec *record.Game
	if gameConfig.Record != nil {
		// Copy the record to avoid modifying the original
		recCopy := *gameConfig.Record
		rec = &recCopy
	} else {
		rec = &record.Game{}
	}

	rec = t.applyGameRecDefaultValues(rec)

	l.Info("creating game record >%#v<", rec)

	rec, err := t.Domain.(*domain.Domain).CreateGameRec(rec)
	if err != nil {
		l.Warn("failed creating game record >%v<", err)
		return nil, err
	}

	// Add the game record to the data store
	t.Data.AddGameRec(rec)

	// Add the game record to the teardown data store
	t.teardownData.AddGameRec(rec)

	// Add the game record to the data store refs
	if gameConfig.Reference != "" {
		t.Data.Refs.GameRefs[gameConfig.Reference] = rec.ID
	}

	return rec, nil
}

func (t *Testing) applyGameRecDefaultValues(rec *record.Game) *record.Game {
	if rec == nil {
		rec = &record.Game{}
	}

	if rec.Name == "" {
		rec.Name = UniqueName(gofakeit.Name())
	}

	if rec.GameType == "" {
		rec.GameType = record.GameTypeAdventure
	}

	return rec
}

// GameCharacter
func (t *Testing) createGameCharacterRec(charConfig GameCharacterConfig, gameRec *record.Game) (*record.GameCharacter, error) {
	l := t.Logger("createGameCharacterRec")

	if gameRec == nil {
		return nil, fmt.Errorf("game record is nil for game_character record >%#v<", charConfig)
	}

	accountRec, err := t.Data.GetAccountRecByRef(charConfig.AccountRef)
	if err != nil {
		l.Warn("failed resolving account ref >%s<: %v", charConfig.AccountRef, err)
		return nil, err
	}

	// Create a new record instance to avoid reusing the same record across tests
	var rec *record.GameCharacter
	if charConfig.Record != nil {
		recCopy := *charConfig.Record
		rec = &recCopy
	} else {
		rec = &record.GameCharacter{}
	}

	rec = t.applyGameCharacterRecDefaultValues(rec)

	// Set the game and account IDs
	rec.GameID = gameRec.ID
	rec.AccountID = accountRec.ID

	l.Info("creating game_character record >%#v<", rec)

	rec, err = t.Domain.(*domain.Domain).CreateGameCharacterRec(rec)
	if err != nil {
		l.Warn("failed creating game_character record >%v<", err)
		return nil, err
	}

	// Add the game_character record to the data store
	t.Data.AddGameCharacterRec(rec)

	// Add the game_character record to the teardown data store
	t.teardownData.AddGameCharacterRec(rec)

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

// Location
func (t *Testing) createGameLocationRec(gameLocationConfig GameLocationConfig, gameRec *record.Game) (*record.GameLocation, error) {
	l := t.Logger("createGameLocationRec")

	if gameRec == nil {
		return nil, fmt.Errorf("game record is nil for game location record >%#v<", gameLocationConfig)
	}

	// Create a new record instance to avoid reusing the same record across tests
	var rec *record.GameLocation
	if gameLocationConfig.Record != nil {
		// Copy the record to avoid modifying the original
		recCopy := *gameLocationConfig.Record
		rec = &recCopy
	} else {
		rec = &record.GameLocation{}
	}

	rec = t.applyGameLocationRecDefaultValues(rec)

	// Set the game ID if it is provided
	rec.GameID = gameRec.ID

	l.Info("creating game location record >%#v<", rec)

	rec, err := t.Domain.(*domain.Domain).CreateGameLocationRec(rec)
	if err != nil {
		l.Warn("failed creating game location record >%v<", err)
		return nil, err
	}

	// Add the location record to the data store
	t.Data.AddGameLocationRec(rec)

	// Add the location record to the teardown data store
	t.teardownData.AddGameLocationRec(rec)

	// Add the location record to the data store refs
	if gameLocationConfig.Reference != "" {
		t.Data.Refs.GameLocationRefs[gameLocationConfig.Reference] = rec.ID
	}

	return rec, nil
}

func (t *Testing) applyGameLocationRecDefaultValues(rec *record.GameLocation) *record.GameLocation {
	if rec == nil {
		rec = &record.GameLocation{}
	}
	if rec.Name == "" {
		rec.Name = UniqueName(gofakeit.Name())
	}
	if rec.Description == "" {
		rec.Description = gofakeit.Sentence(10)
	}
	return rec
}

// LocationLink
func (t *Testing) createLocationLinkRec(linkConfig LocationLinkConfig) (*record.LocationLink, error) {
	l := t.Logger("createLocationLinkRec")

	var rec *record.LocationLink
	if linkConfig.Record != nil {
		recCopy := *linkConfig.Record
		rec = &recCopy
	} else {
		rec = &record.LocationLink{}
	}

	rec = t.applyLocationLinkRecDefaultValues(rec)

	// Resolve location references to IDs
	if linkConfig.FromLocationRef != "" {
		fromLoc, err := t.Data.GetGameLocationRecByRef(linkConfig.FromLocationRef)
		if err != nil || fromLoc == nil || fromLoc.ID == "" {
			l.Error("could not resolve FromLocationRef >%s< to a valid location ID", linkConfig.FromLocationRef)
			return nil, fmt.Errorf("could not resolve FromLocationRef >%s< to a valid location ID", linkConfig.FromLocationRef)
		}
		rec.FromGameLocationID = fromLoc.ID
	}
	if linkConfig.ToLocationRef != "" {
		toLoc, err := t.Data.GetGameLocationRecByRef(linkConfig.ToLocationRef)
		if err != nil || toLoc == nil || toLoc.ID == "" {
			l.Error("could not resolve ToLocationRef >%s< to a valid location ID", linkConfig.ToLocationRef)
			return nil, fmt.Errorf("could not resolve ToLocationRef >%s< to a valid location ID", linkConfig.ToLocationRef)
		}
		rec.ToGameLocationID = toLoc.ID
	}

	if rec.FromGameLocationID == "" || rec.ToGameLocationID == "" {
		l.Error("location link must have both FromGameLocationID and ToGameLocationID set, got from: >%s< to: >%s<", rec.FromGameLocationID, rec.ToGameLocationID)
		return nil, fmt.Errorf("location link must have both FromGameLocationID and ToGameLocationID set, got from: >%s< to: >%s<", rec.FromGameLocationID, rec.ToGameLocationID)
	}

	l.Info("creating location link record >%#v<", rec)

	rec, err := t.Domain.(*domain.Domain).CreateLocationLinkRec(rec)
	if err != nil {
		l.Warn("failed creating location link record >%v<", err)
		return nil, err
	}

	t.Data.AddLocationLinkRec(rec)
	t.teardownData.AddLocationLinkRec(rec)

	if linkConfig.Reference != "" {
		t.Data.Refs.LocationLinkRefs[linkConfig.Reference] = rec.ID
	}

	return rec, nil
}

func (t *Testing) applyLocationLinkRecDefaultValues(rec *record.LocationLink) *record.LocationLink {
	if rec == nil {
		rec = &record.LocationLink{}
	}
	if rec.Description == "" {
		rec.Description = gofakeit.Sentence(5)
	}
	if rec.Name == "" {
		rec.Name = "Link " + gofakeit.Word()
	}
	return rec
}
