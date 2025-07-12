package game_instance_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/harness"
	"gitlab.com/alienspaces/playbymail/internal/record"
	"gitlab.com/alienspaces/playbymail/internal/utils/config"
	"gitlab.com/alienspaces/playbymail/internal/utils/deps"
)

const (
	gameRef         = harness.GameOneRef
	gameInstanceRef = harness.GameInstanceOneRef
)

func newHarness(t *testing.T) *harness.Testing {
	dcfg := harness.DataConfig{
		GameConfigs: []harness.GameConfig{
			{
				Reference: gameRef,
				Record:    &record.Game{},
				GameInstanceConfigs: []harness.GameInstanceConfig{
					{
						Reference: gameInstanceRef,
						Record:    &record.GameInstance{},
					},
				},
			},
		},
	}
	cfg, err := config.Parse()
	require.NoError(t, err)
	l, s, j, err := deps.Default(cfg)
	require.NoError(t, err)
	h, err := harness.NewTesting(l, s, j, dcfg)
	require.NoError(t, err)
	return h
}

func TestCreateOne(t *testing.T) {
	tests := []struct {
		name string
		rec  func(d harness.Data, t *testing.T) *record.GameInstance
		err  bool
	}{
		{
			name: "Valid",
			rec: func(d harness.Data, t *testing.T) *record.GameInstance {
				gameRec, err := d.GetGameRecByRef(gameRef)
				require.NoError(t, err)
				return &record.GameInstance{GameID: gameRec.ID}
			},
			err: false,
		},
		{
			name: "Missing GameID",
			rec: func(d harness.Data, t *testing.T) *record.GameInstance {
				return &record.GameInstance{GameID: ""}
			},
			err: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			h := newHarness(t)
			_, err := h.Setup()
			require.NoError(t, err)
			defer func() {
				err = h.Teardown()
				require.NoError(t, err)
			}()
			repo := h.Domain.(*domain.Domain).GameInstanceRepository()
			rec := tc.rec(h.Data, t)
			_, err = repo.CreateOne(rec)
			if tc.err {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotEmpty(t, rec.ID)
		})
	}
}

func TestGetOne(t *testing.T) {
	tests := []struct {
		name string
		id   func(d harness.Data, t *testing.T) string
		err  bool
	}{
		{
			name: "Valid",
			id: func(d harness.Data, t *testing.T) string {
				rec, err := d.GetGameInstanceRecByRef(gameInstanceRef)
				require.NoError(t, err)
				return rec.ID
			},
			err: false,
		},
		{
			name: "Missing ID",
			id:   func(d harness.Data, t *testing.T) string { return "" },
			err:  true,
		},
		{
			name: "Invalid ID",
			id:   func(d harness.Data, t *testing.T) string { return uuid.New().String() },
			err:  true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			h := newHarness(t)
			_, err := h.Setup()
			require.NoError(t, err)
			defer func() {
				err = h.Teardown()
				require.NoError(t, err)
			}()
			repo := h.Domain.(*domain.Domain).GameInstanceRepository()
			rec, err := repo.GetOne(tc.id(h.Data, t), nil)
			if tc.err {
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
	tests := []struct {
		name string
		rec  func(d harness.Data, t *testing.T) *record.GameInstance
		err  bool
	}{
		{
			name: "Valid",
			rec: func(d harness.Data, t *testing.T) *record.GameInstance {
				rec, err := d.GetGameInstanceRecByRef(gameInstanceRef)
				require.NoError(t, err)
				return rec
			},
			err: false,
		},
		{
			name: "Missing ID",
			rec: func(d harness.Data, t *testing.T) *record.GameInstance {
				rec, err := d.GetGameInstanceRecByRef(gameInstanceRef)
				require.NoError(t, err)
				rec.ID = ""
				return rec
			},
			err: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			h := newHarness(t)
			_, err := h.Setup()
			require.NoError(t, err)
			defer func() {
				err = h.Teardown()
				require.NoError(t, err)
			}()
			repo := h.Domain.(*domain.Domain).GameInstanceRepository()
			rec := tc.rec(h.Data, t)
			_, err = repo.UpdateOne(rec)
			if tc.err {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotEmpty(t, rec.UpdatedAt)
		})
	}
}

func TestDeleteOne(t *testing.T) {
	tests := []struct {
		name string
		id   func(d harness.Data, t *testing.T) string
		err  bool
	}{
		{
			name: "Valid",
			id: func(d harness.Data, t *testing.T) string {
				rec, err := d.GetGameInstanceRecByRef(gameInstanceRef)
				require.NoError(t, err)
				return rec.ID
			},
			err: false,
		},
		{
			name: "Missing ID",
			id:   func(d harness.Data, t *testing.T) string { return "" },
			err:  true,
		},
		{
			name: "Invalid ID",
			id:   func(d harness.Data, t *testing.T) string { return uuid.New().String() },
			err:  true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			h := newHarness(t)
			_, err := h.Setup()
			require.NoError(t, err)
			defer func() {
				err = h.Teardown()
				require.NoError(t, err)
			}()
			repo := h.Domain.(*domain.Domain).GameInstanceRepository()
			id := tc.id(h.Data, t)
			err = repo.DeleteOne(id)
			if tc.err {
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
