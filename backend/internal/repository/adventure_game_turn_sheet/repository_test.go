package adventure_game_turn_sheet_test

import (
	"database/sql"
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/harness"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
	"gitlab.com/alienspaces/playbymail/internal/utils/deps"
)

func TestCreateOne(t *testing.T) {
	h := deps.NewHarness(t)

	tests := []struct {
		name   string
		rec    func(d harness.Data, t *testing.T) *adventure_game_record.AdventureGameTurnSheet
		hasErr bool
	}{
		{
			name: "Without ID",
			rec: func(d harness.Data, t *testing.T) *adventure_game_record.AdventureGameTurnSheet {
				// Create test records first
				gameRec, err := d.GetGameRecByRef(harness.GameOneRef)
				require.NoError(t, err, "GetGameRecByRef returns without error")
				characterInstanceRec, err := d.GetAdventureGameCharacterInstanceRecByRef(harness.GameCharacterInstanceOneRef)
				require.NoError(t, err, "GetGameCharacterInstanceRecByRef returns without error")

				// Create a game turn sheet first
				gameInstanceRec, err := d.GetGameInstanceRecByRef(harness.GameInstanceOneRef)
				require.NoError(t, err, "GetGameInstanceRecByRef returns without error")
				accountRec, err := d.GetAccountRecByRef(harness.AccountOneRef)
				require.NoError(t, err, "GetAccountRecByRef returns without error")

				gameTurnSheetRec := &game_record.GameTurnSheet{
					GameID:           gameRec.ID,
					AccountID:        accountRec.ID,
					TurnNumber:       1,
					SheetType:        "location_choice",
					SheetOrder:       1,
					SheetData:        json.RawMessage(`{"locations": ["north", "south", "east", "west"]}`),
					IsCompleted:      false,
					ProcessingStatus: "pending",
				}
				gameTurnSheetRec.GameInstanceID = sql.NullString{String: gameInstanceRec.ID, Valid: true}

				gameTurnSheetRepo := h.Domain.(*domain.Domain).GameTurnSheetRepository()
				_, err = gameTurnSheetRepo.CreateOne(gameTurnSheetRec)
				require.NoError(t, err, "CreateOne game turn sheet returns without error")

				return &adventure_game_record.AdventureGameTurnSheet{
					GameID:                           gameRec.ID,
					AdventureGameCharacterInstanceID: characterInstanceRec.ID,
					GameTurnSheetID:                  gameTurnSheetRec.ID,
				}
			},
			hasErr: false,
		},
		{
			name: "With ID",
			rec: func(d harness.Data, t *testing.T) *adventure_game_record.AdventureGameTurnSheet {
				gameRec, err := d.GetGameRecByRef(harness.GameOneRef)
				require.NoError(t, err, "GetGameRecByRef returns without error")
				characterInstanceRec, err := d.GetAdventureGameCharacterInstanceRecByRef(harness.GameCharacterInstanceOneRef)
				require.NoError(t, err, "GetGameCharacterInstanceRecByRef returns without error")

				// Create a game turn sheet first
				gameInstanceRec, err := d.GetGameInstanceRecByRef(harness.GameInstanceOneRef)
				require.NoError(t, err, "GetGameInstanceRecByRef returns without error")
				accountRec, err := d.GetAccountRecByRef(harness.AccountOneRef)
				require.NoError(t, err, "GetAccountRecByRef returns without error")

				gameTurnSheetRec := &game_record.GameTurnSheet{
					GameID:           gameRec.ID,
					AccountID:        accountRec.ID,
					TurnNumber:       2,
					SheetType:        "combat",
					SheetOrder:       1,
					SheetData:        json.RawMessage(`{"enemies": ["goblin", "orc"], "weapons": ["sword", "bow"]}`),
					IsCompleted:      false,
					ProcessingStatus: "pending",
				}
				gameTurnSheetRec.GameInstanceID = sql.NullString{String: gameInstanceRec.ID, Valid: true}

				gameTurnSheetRepo := h.Domain.(*domain.Domain).GameTurnSheetRepository()
				_, err = gameTurnSheetRepo.CreateOne(gameTurnSheetRec)
				require.NoError(t, err, "CreateOne game turn sheet returns without error")

				rec := &adventure_game_record.AdventureGameTurnSheet{
					GameID:                           gameRec.ID,
					AdventureGameCharacterInstanceID: characterInstanceRec.ID,
					GameTurnSheetID:                  gameTurnSheetRec.ID,
				}
				id, _ := uuid.NewRandom()
				rec.ID = id.String()
				return rec
			},
			hasErr: false,
		},
	}

	for _, tc := range tests {

		t.Logf("Run test >%s<", tc.name)

		t.Run(tc.name, func(t *testing.T) {

			// Test harness
			_, err := h.Setup()
			require.NoError(t, err, "Setup returns without error")
			defer func() {
				err = h.Teardown()
				require.NoError(t, err, "Teardown returns without error")
			}()

			// repository
			r := h.Domain.(*domain.Domain).AdventureGameTurnSheetRepository()
			require.NotNil(t, r, "Repository is not nil")

			rec := tc.rec(h.Data, t)

			_, err = r.CreateOne(rec)
			if tc.hasErr {
				require.Error(t, err, "CreateOne returns error")
				return
			}
			require.NoError(t, err, "CreateOne returns without error")
			require.NotEmpty(t, rec.CreatedAt, "CreateOne returns record with CreatedAt")
		})
	}
}

func TestGetOne(t *testing.T) {
	h := deps.NewHarness(t)

	tests := []struct {
		name   string
		id     func(d harness.Data, t *testing.T) string
		hasErr bool
	}{
		{
			name: "With valid ID",
			id: func(d harness.Data, t *testing.T) string {
				// Create test records first
				gameRec, err := d.GetGameRecByRef(harness.GameOneRef)
				require.NoError(t, err, "GetGameRecByRef returns without error")
				characterInstanceRec, err := d.GetAdventureGameCharacterInstanceRecByRef(harness.GameCharacterInstanceOneRef)
				require.NoError(t, err, "GetGameCharacterInstanceRecByRef returns without error")

				gameInstanceRec, err := d.GetGameInstanceRecByRef(harness.GameInstanceOneRef)
				require.NoError(t, err, "GetGameInstanceRecByRef returns without error")
				accountRec, err := d.GetAccountRecByRef(harness.AccountOneRef)
				require.NoError(t, err, "GetAccountRecByRef returns without error")

				gameTurnSheetRec := &game_record.GameTurnSheet{
					GameID:           gameRec.ID,
					AccountID:        accountRec.ID,
					TurnNumber:       1,
					SheetType:        "inventory",
					SheetOrder:       1,
					SheetData:        json.RawMessage(`{"items": ["sword", "shield", "potion"]}`),
					IsCompleted:      false,
					ProcessingStatus: "pending",
				}
				gameTurnSheetRec.GameInstanceID = sql.NullString{String: gameInstanceRec.ID, Valid: true}

				gameTurnSheetRepo := h.Domain.(*domain.Domain).GameTurnSheetRepository()
				_, err = gameTurnSheetRepo.CreateOne(gameTurnSheetRec)
				require.NoError(t, err, "CreateOne game turn sheet returns without error")

				rec := &adventure_game_record.AdventureGameTurnSheet{
					GameID:                           gameRec.ID,
					AdventureGameCharacterInstanceID: characterInstanceRec.ID,
					GameTurnSheetID:                  gameTurnSheetRec.ID,
				}

				r := h.Domain.(*domain.Domain).AdventureGameTurnSheetRepository()
				_, err = r.CreateOne(rec)
				require.NoError(t, err, "CreateOne returns without error")
				return rec.ID
			},
			hasErr: false,
		},
		{
			name: "With invalid ID",
			id: func(d harness.Data, t *testing.T) string {
				return "invalid-uuid"
			},
			hasErr: true,
		},
	}

	for _, tc := range tests {

		t.Logf("Run test >%s<", tc.name)

		t.Run(tc.name, func(t *testing.T) {

			// Test harness
			_, err := h.Setup()
			require.NoError(t, err, "Setup returns without error")
			defer func() {
				err = h.Teardown()
				require.NoError(t, err, "Teardown returns without error")
			}()

			// repository
			r := h.Domain.(*domain.Domain).AdventureGameTurnSheetRepository()
			require.NotNil(t, r, "Repository is not nil")

			id := tc.id(h.Data, t)

			rec, err := r.GetOne(id, nil)
			if tc.hasErr {
				require.Error(t, err, "GetOne returns error")
				return
			}
			require.NoError(t, err, "GetOne returns without error")
			require.NotNil(t, rec, "GetOne returns record")
			require.Equal(t, id, rec.ID, "GetOne returns correct record")
		})
	}
}

func TestDeleteOne(t *testing.T) {
	h := deps.NewHarness(t)

	tests := []struct {
		name   string
		id     func(d harness.Data, t *testing.T) string
		hasErr bool
	}{
		{
			name: "With valid ID",
			id: func(d harness.Data, t *testing.T) string {
				// Create test records first
				gameRec, err := d.GetGameRecByRef(harness.GameOneRef)
				require.NoError(t, err, "GetGameRecByRef returns without error")
				characterInstanceRec, err := d.GetAdventureGameCharacterInstanceRecByRef(harness.GameCharacterInstanceOneRef)
				require.NoError(t, err, "GetGameCharacterInstanceRecByRef returns without error")

				gameInstanceRec, err := d.GetGameInstanceRecByRef(harness.GameInstanceOneRef)
				require.NoError(t, err, "GetGameInstanceRecByRef returns without error")
				accountRec, err := d.GetAccountRecByRef(harness.AccountOneRef)
				require.NoError(t, err, "GetAccountRecByRef returns without error")

				gameTurnSheetRec := &game_record.GameTurnSheet{
					GameID:           gameRec.ID,
					AccountID:        accountRec.ID,
					TurnNumber:       3,
					SheetType:        "combat",
					SheetOrder:       1,
					SheetData:        json.RawMessage(`{"enemies": ["dragon"], "weapons": ["sword", "bow"]}`),
					IsCompleted:      false,
					ProcessingStatus: "pending",
				}
				gameTurnSheetRec.GameInstanceID = sql.NullString{String: gameInstanceRec.ID, Valid: true}

				gameTurnSheetRepo := h.Domain.(*domain.Domain).GameTurnSheetRepository()
				_, err = gameTurnSheetRepo.CreateOne(gameTurnSheetRec)
				require.NoError(t, err, "CreateOne game turn sheet returns without error")

				rec := &adventure_game_record.AdventureGameTurnSheet{
					GameID:                           gameRec.ID,
					AdventureGameCharacterInstanceID: characterInstanceRec.ID,
					GameTurnSheetID:                  gameTurnSheetRec.ID,
				}

				r := h.Domain.(*domain.Domain).AdventureGameTurnSheetRepository()
				_, err = r.CreateOne(rec)
				require.NoError(t, err, "CreateOne returns without error")
				return rec.ID
			},
			hasErr: false,
		},
	}

	for _, tc := range tests {

		t.Logf("Run test >%s<", tc.name)

		t.Run(tc.name, func(t *testing.T) {

			// Test harness
			_, err := h.Setup()
			require.NoError(t, err, "Setup returns without error")
			defer func() {
				err = h.Teardown()
				require.NoError(t, err, "Teardown returns without error")
			}()

			// repository
			r := h.Domain.(*domain.Domain).AdventureGameTurnSheetRepository()
			require.NotNil(t, r, "Repository is not nil")

			id := tc.id(h.Data, t)

			err = r.DeleteOne(id)
			if tc.hasErr {
				require.Error(t, err, "DeleteOne returns error")
				return
			}
			require.NoError(t, err, "DeleteOne returns without error")
		})
	}
}
