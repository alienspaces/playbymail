package domain

import (
	"errors"

	"github.com/jackc/pgx/v5"
	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/record"
)

// GetManyGameLocationLinkRequirementRecs -
func (m *Domain) GetManyGameLocationLinkRequirementRecs(opts *coresql.Options) ([]*record.GameLocationLinkRequirement, error) {
	l := m.Logger("GetManyGameLocationLinkRequirementRecs")
	l.Debug("getting many game_location_link_requirement records opts >%#v<", opts)
	r := m.GameLocationLinkRequirementRepository()
	recs, err := r.GetMany(opts)
	if err != nil {
		return nil, databaseError(err)
	}
	return recs, nil
}

// GetGameLocationLinkRequirementRec -
func (m *Domain) GetGameLocationLinkRequirementRec(recID string, lock *coresql.Lock) (*record.GameLocationLinkRequirement, error) {
	l := m.Logger("GetGameLocationLinkRequirementRec")
	l.Debug("getting game_location_link_requirement record ID >%s<", recID)
	if err := domain.ValidateUUIDField("id", recID); err != nil {
		return nil, err
	}
	r := m.GameLocationLinkRequirementRepository()
	rec, err := r.GetOne(recID, lock)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, coreerror.NewNotFoundError(record.TableGameLocationLinkRequirement, recID)
	} else if err != nil {
		return nil, databaseError(err)
	}
	return rec, nil
}

// CreateGameLocationLinkRequirementRec -
func (m *Domain) CreateGameLocationLinkRequirementRec(rec *record.GameLocationLinkRequirement) (*record.GameLocationLinkRequirement, error) {
	l := m.Logger("CreateGameLocationLinkRequirementRec")
	l.Debug("creating game_location_link_requirement record >%#v<", rec)
	r := m.GameLocationLinkRequirementRepository()
	// Add validation here if needed
	var err error
	rec, err = r.CreateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}
	return rec, nil
}

// UpdateGameLocationLinkRequirementRec -
func (m *Domain) UpdateGameLocationLinkRequirementRec(next *record.GameLocationLinkRequirement) (*record.GameLocationLinkRequirement, error) {
	l := m.Logger("UpdateGameLocationLinkRequirementRec")
	_, err := m.GetGameLocationLinkRequirementRec(next.ID, coresql.ForUpdateNoWait)
	if err != nil {
		return next, err
	}
	l.Debug("updating game_location_link_requirement record >%#v<", next)
	// Add validation here if needed
	r := m.GameLocationLinkRequirementRepository()
	next, err = r.UpdateOne(next)
	if err != nil {
		return next, databaseError(err)
	}
	return next, nil
}

// DeleteGameLocationLinkRequirementRec -
func (m *Domain) DeleteGameLocationLinkRequirementRec(recID string) error {
	l := m.Logger("DeleteGameLocationLinkRequirementRec")
	l.Debug("deleting game_location_link_requirement record ID >%s<", recID)
	_, err := m.GetGameLocationLinkRequirementRec(recID, coresql.ForUpdateNoWait)
	if err != nil {
		return err
	}
	r := m.GameLocationLinkRequirementRepository()
	// Add validation here if needed
	if err := r.DeleteOne(recID); err != nil {
		return databaseError(err)
	}
	return nil
}

// RemoveGameLocationLinkRequirementRec -
func (m *Domain) RemoveGameLocationLinkRequirementRec(recID string) error {
	l := m.Logger("RemoveGameLocationLinkRequirementRec")
	l.Debug("removing game_location_link_requirement record ID >%s<", recID)
	_, err := m.GetGameLocationLinkRequirementRec(recID, coresql.ForUpdateNoWait)
	if err != nil {
		return err
	}
	r := m.GameLocationLinkRequirementRepository()
	// Add validation here if needed
	if err := r.RemoveOne(recID); err != nil {
		return databaseError(err)
	}
	return nil
}
