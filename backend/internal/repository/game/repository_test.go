package game_test

// NOTE: repository tests are run is the public space so we are
// able to use common setup and teardown tooling for all repositories

import (
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit"
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
		name string
		rec  func(d harness.Data, t *testing.T) *game_record.Game
		err  bool
	}{
		{
			name: "Without ID",
			rec: func(d harness.Data, t *testing.T) *game_record.Game {
				accountRec, err := d.GetAccountRecByRef(harness.AccountOneRef)
				require.NoError(t, err, "GetAccountRecByRef returns without error")
				return &game_record.Game{
					AccountID:         accountRec.ID,
					Name:              fmt.Sprintf("%s %s", gofakeit.Name(), gofakeit.Name()),
					Description:       "Test game description",
					GameType:          game_record.GameTypeAdventure,
					TurnDurationHours: 168, // 1 week
				}
			},
			err: false,
		},
		{
			name: "With ID",
			rec: func(d harness.Data, t *testing.T) *game_record.Game {
				accountRec, err := d.GetAccountRecByRef(harness.AccountOneRef)
				require.NoError(t, err, "GetAccountRecByRef returns without error")
				rec := &game_record.Game{
					AccountID:         accountRec.ID,
					Name:              fmt.Sprintf("%s %s", gofakeit.Name(), gofakeit.Name()),
					Description:       "Test game description",
					GameType:          game_record.GameTypeAdventure,
					TurnDurationHours: 336, // 2 weeks
				}
				id, _ := uuid.NewRandom()
				rec.ID = id.String()
				return rec
			},
			err: false,
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
			r := h.Domain.(*domain.Domain).GameRepository()
			require.NotNil(t, r, "Repository is not nil")

			rec := tc.rec(h.Data, t)

			_, err = r.CreateOne(rec)
			if tc.err == true {
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
		name string
		id   func(d harness.Data, t *testing.T) string
		err  bool
	}{
		{
			name: "With ID",
			id: func(d harness.Data, t *testing.T) string {
				gameRec, err := d.GetGameRecByRef(harness.GameOneRef)
				require.NoError(t, err, "GetGameRecByRef returns without error")
				return gameRec.ID
			},
			err: false,
		},
		{
			name: "Without ID",
			id: func(d harness.Data, t *testing.T) string {
				return ""
			},
			err: true,
		},
	}

	for _, tc := range tests {

		t.Logf("Run test >%s<", tc.name)

		t.Run(tc.name, func(t *testing.T) {

			// harness setup
			_, err := h.Setup()
			require.NoError(t, err, "Setup returns without error")
			defer func() {
				err = h.Teardown()
				require.NoError(t, err, "Teardown returns without error")
			}()

			// repository
			r := h.Domain.(*domain.Domain).GameRepository()
			require.NotNil(t, r, "Repository is not nil")

			rec, err := r.GetOne(tc.id(h.Data, t), nil)
			if tc.err == true {
				require.Error(t, err, "GetOne returns error")
				return
			}
			require.NoError(t, err, "GetOne returns without error")
			require.NotNil(t, rec, "GetOne returns record")
			require.NotEmpty(t, rec.ID, "Record ID is not empty")
		})
	}
}

func TestUpdateOne(t *testing.T) {
	h := deps.NewHarness(t)

	tests := []struct {
		name string
		rec  func(d harness.Data, t *testing.T) *game_record.Game
		err  bool
	}{
		{
			name: "With ID",
			rec: func(d harness.Data, t *testing.T) *game_record.Game {
				gameRec, err := d.GetGameRecByRef(harness.GameOneRef)
				require.NoError(t, err, "GetGameRecByRef returns without error")
				return gameRec
			},
			err: false,
		},
		{
			name: "Without ID",
			rec: func(d harness.Data, t *testing.T) *game_record.Game {
				gameRec, err := d.GetGameRecByRef(harness.GameOneRef)
				require.NoError(t, err, "GetGameRecByRef returns without error")
				rec := gameRec
				rec.ID = ""
				return rec
			},
			err: true,
		},
	}

	for _, tc := range tests {

		t.Logf("Run test >%s<", tc.name)

		t.Run(tc.name, func(t *testing.T) {

			// harness setup
			_, err := h.Setup()
			require.NoError(t, err, "Setup returns without error")
			defer func() {
				err = h.Teardown()
				require.NoError(t, err, "Teardown returns without error")
			}()

			// repository
			r := h.Domain.(*domain.Domain).GameRepository()
			require.NotNil(t, r, "Repository is not nil")

			rec := tc.rec(h.Data, t)

			_, err = r.UpdateOne(rec)
			if tc.err == true {
				require.Error(t, err, "UpdateOne returns error")
				return
			}
			require.NoError(t, err, "UpdateOne returns without error")
			require.NotEmpty(t, rec.UpdatedAt, "UpdateOne returns record with UpdatedAt")
		})
	}
}

func TestDeleteOne(t *testing.T) {
	h := deps.NewHarness(t)

	tests := []struct {
		name string
		id   func(d harness.Data, t *testing.T) string
		err  bool
	}{
		{
			name: "With ID",
			id: func(d harness.Data, t *testing.T) string {
				gameRec, err := d.GetGameRecByRef(harness.GameOneRef)
				require.NoError(t, err, "GetGameRecByRef returns without error")
				return gameRec.ID
			},
			err: false,
		},
		{
			name: "Without ID",
			id: func(d harness.Data, t *testing.T) string {
				return ""
			},
			err: true,
		},
	}

	for _, tc := range tests {

		t.Logf("Run test >%s<", tc.name)

		t.Run(tc.name, func(t *testing.T) {

			// harness setup
			_, err := h.Setup()
			require.NoError(t, err, "Setup returns without error")
			defer func() {
				err = h.Teardown()
				require.NoError(t, err, "Teardown returns without error")
			}()

			// repository
			r := h.Domain.(*domain.Domain).GameRepository()
			require.NotNil(t, r, "Repository is not nil")

			err = r.DeleteOne(tc.id(h.Data, t))
			if tc.err == true {
				require.Error(t, err, "DeleteOne returns error")
				return
			}
			require.NoError(t, err, "DeleteOne returns without error")

			rec, err := r.GetOne(tc.id(h.Data, t), nil)
			require.Error(t, err, "GetOne returns error")
			require.Nil(t, rec, "GetOne does not return record")
		})
	}
}
