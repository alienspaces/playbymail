package harness

import (
	"fmt"

	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/record"
)

func (t *Testing) createGameAdministrationRec(administrationConfig GameAdministrationConfig, gameRec *record.Game) (*record.GameAdministration, error) {
	l := t.Logger("createGameAdministrationRec")

	if gameRec == nil {
		return nil, fmt.Errorf("game record is nil for game_administration record >%#v<", administrationConfig)
	}

	if administrationConfig.AccountRef == "" {
		return nil, fmt.Errorf("game_administration record >%#v< must have an AccountRef set", administrationConfig)
	}

	if administrationConfig.GrantedByAccountRef == "" {
		return nil, fmt.Errorf("game_administration record >%#v< must have a GrantedByAccountRef set", administrationConfig)
	}

	var rec *record.GameAdministration
	if administrationConfig.Record != nil {
		recCopy := *administrationConfig.Record
		rec = &recCopy
	} else {
		rec = &record.GameAdministration{}
	}

	rec = t.applyGameAdministrationRecDefaultValues(rec)

	rec.GameID = gameRec.ID

	// Get account record
	accountRec, err := t.Data.GetAccountRecByRef(administrationConfig.AccountRef)
	if err != nil {
		l.Warn("failed resolving account ref >%s<: %v", administrationConfig.AccountRef, err)
		return nil, err
	}
	rec.AccountID = accountRec.ID

	// Get granted by account record
	grantedByAccountRec, err := t.Data.GetAccountRecByRef(administrationConfig.GrantedByAccountRef)
	if err != nil {
		l.Warn("failed resolving granted by account ref >%s<: %v", administrationConfig.GrantedByAccountRef, err)
		return nil, err
	}
	rec.GrantedByAccountID = grantedByAccountRec.ID

	// Create record
	l.Info("creating game_administration record >%#v<", rec)

	rec, err = t.Domain.(*domain.Domain).CreateGameAdministrationRec(rec)
	if err != nil {
		l.Warn("failed creating game_administration record >%v<", err)
		return nil, err
	}

	// Add to data store
	t.Data.AddGameAdministrationRec(rec)

	// Add to teardown data store
	t.teardownData.AddGameAdministrationRec(rec)

	// Add to references store
	if administrationConfig.Reference != "" {
		t.Data.Refs.GameAdministrationRefs[administrationConfig.Reference] = rec.ID
	}

	return rec, nil
}

func (t *Testing) applyGameAdministrationRecDefaultValues(rec *record.GameAdministration) *record.GameAdministration {
	if rec == nil {
		rec = &record.GameAdministration{}
	}
	return rec
}
