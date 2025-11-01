package adventure_game_creature_test

import (
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/harness"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
	"gitlab.com/alienspaces/playbymail/internal/utils/deps"
)

func TestCreateOne(t *testing.T) {
	h := deps.NewHarness(t)

	tests := []struct {
		name   string
		rec    func(d harness.Data, t *testing.T) *adventure_game_record.AdventureGameCreature
		hasErr bool
	}{
		{
			name: "Without ID",
			rec: func(d harness.Data, t *testing.T) *adventure_game_record.AdventureGameCreature {
				gameRec, err := d.GetGameRecByRef(harness.GameOneRef)
				require.NoError(t, err, "GetGameRecByRef returns without error")
				return &adventure_game_record.AdventureGameCreature{
					Name:   harness.UniqueName(gofakeit.Name()),
					GameID: gameRec.ID,
				}
			},
			hasErr: false,
		},
		{
			name: "With ID",
			rec: func(d harness.Data, t *testing.T) *adventure_game_record.AdventureGameCreature {
				gameRec, err := d.GetGameRecByRef(harness.GameOneRef)
				require.NoError(t, err, "GetGameRecByRef returns without error")
				rec := &adventure_game_record.AdventureGameCreature{
					Name:   harness.UniqueName(gofakeit.Name()),
					GameID: gameRec.ID,
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
			_, err := h.Setup()
			require.NoError(t, err, "Setup returns without error")
			defer func() {
				err = h.Teardown()
				require.NoError(t, err, "Teardown returns without error")
			}()

			r := h.Domain.(*domain.Domain).AdventureGameCreatureRepository()
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
			name: "With ID",
			id: func(d harness.Data, t *testing.T) string {
				rec, err := d.GetAdventureGameCreatureRecByRef(harness.GameCreatureOneRef)
				require.NoError(t, err, "GetGameCreatureRecByRef returns without error")
				return rec.ID
			},
			hasErr: false,
		},
		{
			name: "Without ID",
			id: func(d harness.Data, t *testing.T) string {
				return ""
			},
			hasErr: true,
		},
	}

	for _, tc := range tests {
		t.Logf("Run test >%s<", tc.name)
		t.Run(tc.name, func(t *testing.T) {
			_, err := h.Setup()
			require.NoError(t, err, "Setup returns without error")
			defer func() {
				err = h.Teardown()
				require.NoError(t, err, "Teardown returns without error")
			}()

			r := h.Domain.(*domain.Domain).AdventureGameCreatureRepository()
			require.NotNil(t, r, "Repository is not nil")

			rec, err := r.GetOne(tc.id(h.Data, t), nil)
			if tc.hasErr {
				require.Error(t, err, "GetOne returns error")
				require.Nil(t, rec, "GetOne does not return record")
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
		name   string
		rec    func(d harness.Data, t *testing.T) *adventure_game_record.AdventureGameCreature
		hasErr bool
	}{
		{
			name: "With ID",
			rec: func(d harness.Data, t *testing.T) *adventure_game_record.AdventureGameCreature {
				rec, err := d.GetAdventureGameCreatureRecByRef(harness.GameCreatureOneRef)
				require.NoError(t, err, "GetGameCreatureRecByRef returns without error")
				return rec
			},
			hasErr: false,
		},
		{
			name: "Without ID",
			rec: func(d harness.Data, t *testing.T) *adventure_game_record.AdventureGameCreature {
				rec, err := d.GetAdventureGameCreatureRecByRef(harness.GameCreatureOneRef)
				require.NoError(t, err, "GetGameCreatureRecByRef returns without error")
				rec.ID = ""
				return rec
			},
			hasErr: true,
		},
	}

	for _, tc := range tests {
		t.Logf("Run test >%s<", tc.name)
		t.Run(tc.name, func(t *testing.T) {
			_, err := h.Setup()
			require.NoError(t, err, "Setup returns without error")
			defer func() {
				err = h.Teardown()
				require.NoError(t, err, "Teardown returns without error")
			}()

			r := h.Domain.(*domain.Domain).AdventureGameCreatureRepository()
			require.NotNil(t, r, "Repository is not nil")

			rec := tc.rec(h.Data, t)

			_, err = r.UpdateOne(rec)
			if tc.hasErr {
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
		name   string
		id     func(d harness.Data, t *testing.T) string
		hasErr bool
	}{
		{
			name: "With ID",
			id: func(d harness.Data, t *testing.T) string {
				rec, err := d.GetAdventureGameCreatureRecByRef(harness.GameCreatureOneRef)
				require.NoError(t, err, "GetGameCreatureRecByRef returns without error")
				return rec.ID
			},
			hasErr: false,
		},
		{
			name: "Without ID",
			id: func(d harness.Data, t *testing.T) string {
				return ""
			},
			hasErr: true,
		},
	}

	for _, tc := range tests {
		t.Logf("Run test >%s<", tc.name)
		t.Run(tc.name, func(t *testing.T) {
			_, err := h.Setup()
			require.NoError(t, err, "Setup returns without error")
			defer func() {
				err = h.Teardown()
				require.NoError(t, err, "Teardown returns without error")
			}()

			r := h.Domain.(*domain.Domain).AdventureGameCreatureRepository()
			require.NotNil(t, r, "Repository is not nil")

			err = r.DeleteOne(tc.id(h.Data, t))
			if tc.hasErr {
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
