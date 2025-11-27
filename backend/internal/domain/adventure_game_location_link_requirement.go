package domain

import (
	"errors"

	"github.com/jackc/pgx/v5"
	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
)

// GetManyAdventureGameLocationLinkRequirementRecs -
func (m *Domain) GetManyAdventureGameLocationLinkRequirementRecs(opts *coresql.Options) ([]*adventure_game_record.AdventureGameLocationLinkRequirement, error) {
	l := m.Logger("GetManyAdventureGameLocationLinkRequirementRecs")

	l.Debug("getting many adventure_game_location_link_requirement records opts >%#v<", opts)

	r := m.AdventureGameLocationLinkRequirementRepository()

	recs, err := r.GetMany(opts)
	if err != nil {
		return nil, databaseError(err)
	}

	return recs, nil
}

// GetAdventureGameLocationLinkRequirementRec -
func (m *Domain) GetAdventureGameLocationLinkRequirementRec(recID string, lock *coresql.Lock) (*adventure_game_record.AdventureGameLocationLinkRequirement, error) {
	l := m.Logger("GetAdventureGameLocationLinkRequirementRec")

	l.Debug("getting adventure_game_location_link_requirement record ID >%s<", recID)

	if err := domain.ValidateUUIDField("id", recID); err != nil {
		return nil, err
	}

	r := m.AdventureGameLocationLinkRequirementRepository()

	rec, err := r.GetOne(recID, lock)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, coreerror.NewNotFoundError(adventure_game_record.TableAdventureGameLocationLinkRequirement, recID)
	} else if err != nil {
		return nil, databaseError(err)
	}

	return rec, nil
}

// CreateAdventureGameLocationLinkRequirementRec -
func (m *Domain) CreateAdventureGameLocationLinkRequirementRec(rec *adventure_game_record.AdventureGameLocationLinkRequirement) (*adventure_game_record.AdventureGameLocationLinkRequirement, error) {
	l := m.Logger("CreateAdventureGameLocationLinkRequirementRec")

	l.Debug("creating adventure_game_location_link_requirement record >%#v<", rec)

	r := m.AdventureGameLocationLinkRequirementRepository()

	if err := m.validateAdventureGameLocationLinkRequirementRecForCreate(rec); err != nil {
		l.Warn("failed to validate adventure_game_location_link_requirement record >%v<", err)
		return rec, err
	}

	var err error
	rec, err = r.CreateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}

	return rec, nil
}

// UpdateAdventureGameLocationLinkRequirementRec -
func (m *Domain) UpdateAdventureGameLocationLinkRequirementRec(rec *adventure_game_record.AdventureGameLocationLinkRequirement) (*adventure_game_record.AdventureGameLocationLinkRequirement, error) {
	l := m.Logger("UpdateAdventureGameLocationLinkRequirementRec")

	_, err := m.GetAdventureGameLocationLinkRequirementRec(rec.ID, coresql.ForUpdateNoWait)
	if err != nil {
		return rec, err
	}

	l.Debug("updating adventure_game_location_link_requirement record >%#v<", rec)

	if err := m.validateAdventureGameLocationLinkRequirementRecForUpdate(rec); err != nil {
		l.Warn("failed to validate adventure_game_location_link_requirement record >%v<", err)
		return rec, err
	}

	r := m.AdventureGameLocationLinkRequirementRepository()

	updatedRec, err := r.UpdateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}

	return updatedRec, nil
}

// DeleteAdventureGameLocationLinkRequirementRec -
func (m *Domain) DeleteAdventureGameLocationLinkRequirementRec(recID string) error {
	l := m.Logger("DeleteAdventureGameLocationLinkRequirementRec")

	l.Debug("deleting adventure_game_location_link_requirement record ID >%s<", recID)

	_, err := m.GetAdventureGameLocationLinkRequirementRec(recID, coresql.ForUpdateNoWait)
	if err != nil {
		return err
	}

	r := m.AdventureGameLocationLinkRequirementRepository()

	// Add validation here if needed
	if err := r.DeleteOne(recID); err != nil {
		return databaseError(err)
	}

	return nil
}

// RemoveAdventureGameLocationLinkRequirementRec -
func (m *Domain) RemoveAdventureGameLocationLinkRequirementRec(recID string) error {
	l := m.Logger("RemoveAdventureGameLocationLinkRequirementRec")

	l.Debug("removing adventure_game_location_link_requirement record ID >%s<", recID)

	_, err := m.GetAdventureGameLocationLinkRequirementRec(recID, coresql.ForUpdateNoWait)
	if err != nil {
		return err
	}

	r := m.AdventureGameLocationLinkRequirementRepository()

	// Add validation here if needed
	if err := r.RemoveOne(recID); err != nil {
		return databaseError(err)
	}

	return nil
}
