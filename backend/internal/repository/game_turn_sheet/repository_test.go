package game_turn_sheet_test

import (
	"database/sql"
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/harness"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
	"gitlab.com/alienspaces/playbymail/internal/utils/deps"
)

func TestCreateOne(t *testing.T) {
	h := deps.NewHarness(t)

	tests := []struct {
		name   string
		rec    func(d harness.Data, t *testing.T) *game_record.GameTurnSheet
		hasErr bool
	}{
		{
			name: "Without ID",
			rec: func(d harness.Data, t *testing.T) *game_record.GameTurnSheet {
				gameRec, err := d.GetGameRecByRef(harness.GameOneRef)
				require.NoError(t, err, "GetGameRecByRef returns without error")
				gameInstanceRec, err := d.GetGameInstanceRecByRef(harness.GameInstanceOneRef)
				require.NoError(t, err, "GetGameInstanceRecByRef returns without error")
				accountRec, err := d.GetAccountUserRecByRef(harness.StandardAccountRef)
				require.NoError(t, err, "GetAccountUserRecByRef returns without error")

				sheetData := map[string]interface{}{
					"type": "location_choice",
					"options": []string{
						"north",
						"south",
						"east",
						"west",
					},
				}
				sheetDataBytes, _ := json.Marshal(sheetData)

				rec := &game_record.GameTurnSheet{
					GameID:           gameRec.ID,
					AccountID:        accountRec.AccountID,
					AccountUserID:    accountRec.ID,
					TurnNumber:       1,
					SheetType:        "location_choice",
					SheetOrder:       1,
					SheetData:        sheetDataBytes,
					IsCompleted:      false,
					ProcessingStatus: "pending",
				}
				rec.GameInstanceID = sql.NullString{String: gameInstanceRec.ID, Valid: true}
				return rec
			},
			hasErr: false,
		},
		{
			name: "With ID",
			rec: func(d harness.Data, t *testing.T) *game_record.GameTurnSheet {
				gameRec, err := d.GetGameRecByRef(harness.GameOneRef)
				require.NoError(t, err, "GetGameRecByRef returns without error")
				gameInstanceRec, err := d.GetGameInstanceRecByRef(harness.GameInstanceOneRef)
				require.NoError(t, err, "GetGameInstanceRecByRef returns without error")
				accountRec, err := d.GetAccountUserRecByRef(harness.StandardAccountRef)
				require.NoError(t, err, "GetAccountUserRecByRef returns without error")

				sheetData := map[string]interface{}{
					"type": "combat",
					"enemies": []map[string]interface{}{
						{"name": "goblin", "hp": 10},
						{"name": "orc", "hp": 20},
					},
				}
				sheetDataBytes, _ := json.Marshal(sheetData)

				rec := &game_record.GameTurnSheet{
					GameID:           gameRec.ID,
					AccountID:        accountRec.AccountID,
					AccountUserID:    accountRec.ID,
					TurnNumber:       2,
					SheetType:        "combat",
					SheetOrder:       1,
					SheetData:        sheetDataBytes,
					IsCompleted:      false,
					ProcessingStatus: "pending",
				}
				rec.GameInstanceID = sql.NullString{String: gameInstanceRec.ID, Valid: true}
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
			r := h.Domain.(*domain.Domain).GameTurnSheetRepository()
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
				// Create a test record first
				gameRec, err := d.GetGameRecByRef(harness.GameOneRef)
				require.NoError(t, err, "GetGameRecByRef returns without error")
				gameInstanceRec, err := d.GetGameInstanceRecByRef(harness.GameInstanceOneRef)
				require.NoError(t, err, "GetGameInstanceRecByRef returns without error")
				accountRec, err := d.GetAccountUserRecByRef(harness.StandardAccountRef)
				require.NoError(t, err, "GetAccountUserRecByRef returns without error")

				sheetData := map[string]interface{}{
					"type":  "inventory",
					"items": []string{"sword", "shield", "potion"},
				}
				sheetDataBytes, _ := json.Marshal(sheetData)

				rec := &game_record.GameTurnSheet{
					GameID:           gameRec.ID,
					AccountID:        accountRec.AccountID,
					AccountUserID:    accountRec.ID,
					TurnNumber:       1,
					SheetType:        "inventory",
					SheetOrder:       1,
					SheetData:        sheetDataBytes,
					IsCompleted:      false,
					ProcessingStatus: "pending",
				}
				rec.GameInstanceID = sql.NullString{String: gameInstanceRec.ID, Valid: true}

				r := h.Domain.(*domain.Domain).GameTurnSheetRepository()
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
			r := h.Domain.(*domain.Domain).GameTurnSheetRepository()
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

func TestUpdateOne(t *testing.T) {
	h := deps.NewHarness(t)

	tests := []struct {
		name   string
		rec    func(d harness.Data, t *testing.T) *game_record.GameTurnSheet
		hasErr bool
	}{
		{
			name: "Update sheet data",
			rec: func(d harness.Data, t *testing.T) *game_record.GameTurnSheet {
				// Create a test record first
				gameRec, err := d.GetGameRecByRef(harness.GameOneRef)
				require.NoError(t, err, "GetGameRecByRef returns without error")
				gameInstanceRec, err := d.GetGameInstanceRecByRef(harness.GameInstanceOneRef)
				require.NoError(t, err, "GetGameInstanceRecByRef returns without error")
				accountRec, err := d.GetAccountUserRecByRef(harness.StandardAccountRef)
				require.NoError(t, err, "GetAccountUserRecByRef returns without error")

				sheetData := map[string]interface{}{
					"type":    "location_choice",
					"options": []string{"north", "south"},
				}
				sheetDataBytes, _ := json.Marshal(sheetData)

				rec := &game_record.GameTurnSheet{
					GameID:           gameRec.ID,
					AccountID:        accountRec.AccountID,
					AccountUserID:    accountRec.ID,
					TurnNumber:       1,
					SheetType:        "location_choice",
					SheetOrder:       1,
					SheetData:        sheetDataBytes,
					IsCompleted:      false,
					ProcessingStatus: "pending",
				}
				rec.GameInstanceID = sql.NullString{String: gameInstanceRec.ID, Valid: true}

				r := h.Domain.(*domain.Domain).GameTurnSheetRepository()
				_, err = r.CreateOne(rec)
				require.NoError(t, err, "CreateOne returns without error")

				// Update the record
				updatedSheetData := map[string]interface{}{
					"type":     "location_choice",
					"options":  []string{"north", "south", "east", "west"},
					"selected": "north",
				}
				updatedSheetDataBytes, _ := json.Marshal(updatedSheetData)
				rec.SheetData = updatedSheetDataBytes
				rec.IsCompleted = true

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
			r := h.Domain.(*domain.Domain).GameTurnSheetRepository()
			require.NotNil(t, r, "Repository is not nil")

			rec := tc.rec(h.Data, t)

			_, err = r.UpdateOne(rec)
			if tc.hasErr {
				require.Error(t, err, "UpdateOne returns error")
				return
			}
			require.NoError(t, err, "UpdateOne returns without error")
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
				// Create a test record first
				gameRec, err := d.GetGameRecByRef(harness.GameOneRef)
				require.NoError(t, err, "GetGameRecByRef returns without error")
				gameInstanceRec, err := d.GetGameInstanceRecByRef(harness.GameInstanceOneRef)
				require.NoError(t, err, "GetGameInstanceRecByRef returns without error")
				accountRec, err := d.GetAccountUserRecByRef(harness.StandardAccountRef)
				require.NoError(t, err, "GetAccountUserRecByRef returns without error")

				sheetData := map[string]interface{}{
					"type": "combat",
					"enemies": []map[string]interface{}{
						{"name": "dragon", "hp": 100},
					},
				}
				sheetDataBytes, _ := json.Marshal(sheetData)

				rec := &game_record.GameTurnSheet{
					GameID:           gameRec.ID,
					AccountID:        accountRec.AccountID,
					AccountUserID:    accountRec.ID,
					TurnNumber:       3,
					SheetType:        "combat",
					SheetOrder:       1,
					SheetData:        sheetDataBytes,
					IsCompleted:      false,
					ProcessingStatus: "pending",
				}
				rec.GameInstanceID = sql.NullString{String: gameInstanceRec.ID, Valid: true}

				r := h.Domain.(*domain.Domain).GameTurnSheetRepository()
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
			r := h.Domain.(*domain.Domain).GameTurnSheetRepository()
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
