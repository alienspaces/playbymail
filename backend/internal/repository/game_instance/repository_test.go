package game_instance_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/harness"
	"gitlab.com/alienspaces/playbymail/internal/record"
	"gitlab.com/alienspaces/playbymail/internal/utils/deps"
)

func TestCreateOne(t *testing.T) {
	h := deps.NewHarness(t)

	tests := []struct {
		name   string
		rec    func(d harness.Data, t *testing.T) *record.GameInstance
		hasErr bool
	}{
		{
			name: "Valid",
			rec: func(d harness.Data, t *testing.T) *record.GameInstance {
				gameRec, err := d.GetGameRecByRef(harness.GameOneRef)
				require.NoError(t, err)
				return &record.GameInstance{GameID: gameRec.ID}
			},
			hasErr: false,
		},
		{
			name: "Missing GameID",
			rec: func(d harness.Data, t *testing.T) *record.GameInstance {
				return &record.GameInstance{GameID: ""}
			},
			hasErr: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := h.Setup()
			require.NoError(t, err)
			defer func() {
				err = h.Teardown()
				require.NoError(t, err)
			}()
			repo := h.Domain.(*domain.Domain).GameInstanceRepository()
			rec := tc.rec(h.Data, t)
			_, err = repo.CreateOne(rec)
			if tc.hasErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotEmpty(t, rec.ID)
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
			name: "Valid",
			id: func(d harness.Data, t *testing.T) string {
				rec, err := d.GetGameInstanceRecByRef(harness.GameInstanceOneRef)
				require.NoError(t, err)
				return rec.ID
			},
			hasErr: false,
		},
		{
			name:   "Missing ID",
			id:     func(d harness.Data, t *testing.T) string { return "" },
			hasErr: true,
		},
		{
			name:   "Invalid ID",
			id:     func(d harness.Data, t *testing.T) string { return uuid.New().String() },
			hasErr: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := h.Setup()
			require.NoError(t, err)
			defer func() {
				err = h.Teardown()
				require.NoError(t, err)
			}()
			repo := h.Domain.(*domain.Domain).GameInstanceRepository()
			rec, err := repo.GetOne(tc.id(h.Data, t), nil)
			if tc.hasErr {
				require.Error(t, err)
				require.Nil(t, rec)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, rec)
			require.NotEmpty(t, rec.ID)
		})
	}
}

func TestUpdateOne(t *testing.T) {
	h := deps.NewHarness(t)

	tests := []struct {
		name   string
		rec    func(d harness.Data, t *testing.T) *record.GameInstance
		hasErr bool
	}{
		{
			name: "Valid",
			rec: func(d harness.Data, t *testing.T) *record.GameInstance {
				rec, err := d.GetGameInstanceRecByRef(harness.GameInstanceOneRef)
				require.NoError(t, err)
				return rec
			},
			hasErr: false,
		},
		{
			name: "Missing ID",
			rec: func(d harness.Data, t *testing.T) *record.GameInstance {
				rec, err := d.GetGameInstanceRecByRef(harness.GameInstanceOneRef)
				require.NoError(t, err)
				rec.ID = ""
				return rec
			},
			hasErr: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := h.Setup()
			require.NoError(t, err)
			defer func() {
				err = h.Teardown()
				require.NoError(t, err)
			}()
			repo := h.Domain.(*domain.Domain).GameInstanceRepository()
			rec := tc.rec(h.Data, t)
			_, err = repo.UpdateOne(rec)
			if tc.hasErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotEmpty(t, rec.UpdatedAt)
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
			name: "Valid",
			id: func(d harness.Data, t *testing.T) string {
				rec, err := d.GetGameInstanceRecByRef(harness.GameInstanceOneRef)
				require.NoError(t, err)
				return rec.ID
			},
			hasErr: false,
		},
		{
			name:   "Missing ID",
			id:     func(d harness.Data, t *testing.T) string { return "" },
			hasErr: true,
		},
		{
			name:   "Invalid ID",
			id:     func(d harness.Data, t *testing.T) string { return uuid.New().String() },
			hasErr: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := h.Setup()
			require.NoError(t, err)
			defer func() {
				err = h.Teardown()
				require.NoError(t, err)
			}()
			repo := h.Domain.(*domain.Domain).GameInstanceRepository()
			id := tc.id(h.Data, t)
			err = repo.DeleteOne(id)
			if tc.hasErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			rec, err := repo.GetOne(id, nil)
			require.Error(t, err)
			require.Nil(t, rec)
		})
	}
}
