package harness

import (
	"fmt"

	"github.com/brianvoe/gofakeit"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/record"
)

func (t *Testing) createGameLocationLinkRec(linkConfig GameLocationLinkConfig, gameRec *record.Game) (*record.GameLocationLink, error) {
	l := t.Logger("createGameLocationLinkRec")

	if gameRec == nil {
		return nil, fmt.Errorf("game record is nil for game_location_link record >%#v<", linkConfig)
	}

	if linkConfig.FromLocationRef == "" && linkConfig.ToLocationRef == "" {
		return nil, fmt.Errorf("game_location_link record >%#v< must have either FromLocationRef or ToLocationRef set", linkConfig)
	}

	var rec *record.GameLocationLink
	if linkConfig.Record != nil {
		recCopy := *linkConfig.Record
		rec = &recCopy
	} else {
		rec = &record.GameLocationLink{}
	}

	rec = t.applyGameLocationLinkRecDefaultValues(rec)

	// Set game_id from parent game
	rec.GameID = gameRec.ID

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

	// Create record
	l.Info("creating location link record >%#v<", rec)

	rec, err := t.Domain.(*domain.Domain).CreateGameLocationLinkRec(rec)
	if err != nil {
		l.Warn("failed creating location link record >%v<", err)
		return nil, err
	}

	// Add to data
	t.Data.AddGameLocationLinkRec(rec)

	// Add to teardown data
	t.teardownData.AddGameLocationLinkRec(rec)

	// Add to references store
	if linkConfig.Reference != "" {
		t.Data.Refs.GameLocationLinkRefs[linkConfig.Reference] = rec.ID
	}

	return rec, nil
}

func (t *Testing) applyGameLocationLinkRecDefaultValues(rec *record.GameLocationLink) *record.GameLocationLink {
	if rec == nil {
		rec = &record.GameLocationLink{}
	}
	if rec.Description == "" {
		rec.Description = gofakeit.Sentence(5)
	}
	if rec.Name == "" {
		rec.Name = "Link " + gofakeit.Word()
	}
	return rec
}
