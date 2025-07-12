package game_creature_test

import (
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/harness"
	"gitlab.com/alienspaces/playbymail/internal/record"
	"gitlab.com/alienspaces/playbymail/internal/utils/config"
	"gitlab.com/alienspaces/playbymail/internal/utils/deps"
)

func TestCreateOne(t *testing.T) {
	dcfg := harness.DataConfig{
		GameConfigs: []harness.GameConfig{
			{
				Reference: harness.GameOneRef,
				Record:    &record.Game{},
				GameCreatureConfigs: []harness.GameCreatureConfig{
					{
						Reference: harness.GameCreatureOneRef,
						Record:    &record.GameCreature{},
					},
				},
			},
		},
	}

	cfg, err := config.Parse()
	require.NoError(t, err, "Parse returns without error")

	l, s, j, err := deps.Default(cfg)
	require.NoError(t, err, "Default dependencies returns without error")

	h, err := harness.NewTesting(l, s, j, dcfg)
	require.NoError(t, err, "NewTesting returns without error")

	h.ShouldCommitData = false

	tests := []struct {
		name string
		rec  func(d harness.Data, t *testing.T) *record.GameCreature
		err  bool
	}{
		{
			name: "Without ID",
			rec: func(d harness.Data, t *testing.T) *record.GameCreature {
				gameRec, err := d.GetGameRecByRef(harness.GameOneRef)
				require.NoError(t, err, "GetGameRecByRef returns without error")
				return &record.GameCreature{
					Name:   harness.UniqueName(gofakeit.Name()),
					GameID: gameRec.ID,
				}
			},
			err: false,
		},
		{
			name: "With ID",
			rec: func(d harness.Data, t *testing.T) *record.GameCreature {
				gameRec, err := d.GetGameRecByRef(harness.GameOneRef)
				require.NoError(t, err, "GetGameRecByRef returns without error")
				rec := &record.GameCreature{
					Name:   harness.UniqueName(gofakeit.Name()),
					GameID: gameRec.ID,
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
			_, err = h.Setup()
			require.NoError(t, err, "Setup returns without error")
			defer func() {
				err = h.Teardown()
				require.NoError(t, err, "Teardown returns without error")
			}()

			r := h.Domain.(*domain.Domain).GameCreatureRepository()
			require.NotNil(t, r, "Repository is not nil")

			rec := tc.rec(h.Data, t)

			_, err = r.CreateOne(rec)
			if tc.err {
				require.Error(t, err, "CreateOne returns error")
				return
			}
			require.NoError(t, err, "CreateOne returns without error")
			require.NotEmpty(t, rec.CreatedAt, "CreateOne returns record with CreatedAt")
		})
	}
}

func TestGetOne(t *testing.T) {
	dcfg := harness.DataConfig{
		GameConfigs: []harness.GameConfig{
			{
				Reference: harness.GameOneRef,
				Record:    &record.Game{},
				GameCreatureConfigs: []harness.GameCreatureConfig{
					{
						Reference: harness.GameCreatureOneRef,
						Record:    &record.GameCreature{},
					},
				},
			},
		},
	}

	cfg, err := config.Parse()
	require.NoError(t, err, "Parse returns without error")

	l, s, j, err := deps.Default(cfg)
	require.NoError(t, err, "Default dependencies returns without error")

	h, err := harness.NewTesting(l, s, j, dcfg)
	require.NoError(t, err, "NewTesting returns without error")

	tests := []struct {
		name string
		id   func(d harness.Data, t *testing.T) string
		err  bool
	}{
		{
			name: "With ID",
			id: func(d harness.Data, t *testing.T) string {
				rec, err := d.GetGameCreatureRecByRef(harness.GameCreatureOneRef)
				require.NoError(t, err, "GetGameCreatureRecByRef returns without error")
				return rec.ID
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
			_, err = h.Setup()
			require.NoError(t, err, "Setup returns without error")
			defer func() {
				err = h.Teardown()
				require.NoError(t, err, "Teardown returns without error")
			}()

			r := h.Domain.(*domain.Domain).GameCreatureRepository()
			require.NotNil(t, r, "Repository is not nil")

			rec, err := r.GetOne(tc.id(h.Data, t), nil)
			if tc.err {
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
	dcfg := harness.DataConfig{
		GameConfigs: []harness.GameConfig{
			{
				Reference: harness.GameOneRef,
				Record:    &record.Game{},
				GameCreatureConfigs: []harness.GameCreatureConfig{
					{
						Reference: harness.GameCreatureOneRef,
						Record:    &record.GameCreature{},
					},
				},
			},
		},
	}

	cfg, err := config.Parse()
	require.NoError(t, err, "Parse returns without error")

	l, s, j, err := deps.Default(cfg)
	require.NoError(t, err, "Default dependencies returns without error")

	h, err := harness.NewTesting(l, s, j, dcfg)
	require.NoError(t, err, "NewTesting returns without error")

	tests := []struct {
		name string
		rec  func(d harness.Data, t *testing.T) *record.GameCreature
		err  bool
	}{
		{
			name: "With ID",
			rec: func(d harness.Data, t *testing.T) *record.GameCreature {
				rec, err := d.GetGameCreatureRecByRef(harness.GameCreatureOneRef)
				require.NoError(t, err, "GetGameCreatureRecByRef returns without error")
				return rec
			},
			err: false,
		},
		{
			name: "Without ID",
			rec: func(d harness.Data, t *testing.T) *record.GameCreature {
				rec, err := d.GetGameCreatureRecByRef(harness.GameCreatureOneRef)
				require.NoError(t, err, "GetGameCreatureRecByRef returns without error")
				rec.ID = ""
				return rec
			},
			err: true,
		},
	}

	for _, tc := range tests {
		t.Logf("Run test >%s<", tc.name)
		t.Run(tc.name, func(t *testing.T) {
			_, err = h.Setup()
			require.NoError(t, err, "Setup returns without error")
			defer func() {
				err = h.Teardown()
				require.NoError(t, err, "Teardown returns without error")
			}()

			r := h.Domain.(*domain.Domain).GameCreatureRepository()
			require.NotNil(t, r, "Repository is not nil")

			rec := tc.rec(h.Data, t)

			_, err := r.UpdateOne(rec)
			if tc.err {
				require.Error(t, err, "UpdateOne returns error")
				return
			}
			require.NoError(t, err, "UpdateOne returns without error")
			require.NotEmpty(t, rec.UpdatedAt, "UpdateOne returns record with UpdatedAt")
		})
	}
}

func TestDeleteOne(t *testing.T) {
	dcfg := harness.DataConfig{
		GameConfigs: []harness.GameConfig{
			{
				Reference: harness.GameOneRef,
				Record:    &record.Game{},
				GameCreatureConfigs: []harness.GameCreatureConfig{
					{
						Reference: harness.GameCreatureOneRef,
						Record:    &record.GameCreature{},
					},
				},
			},
		},
	}

	cfg, err := config.Parse()
	require.NoError(t, err, "Parse returns without error")

	l, s, j, err := deps.Default(cfg)
	require.NoError(t, err, "Default dependencies returns without error")

	h, err := harness.NewTesting(l, s, j, dcfg)
	require.NoError(t, err, "NewTesting returns without error")

	tests := []struct {
		name string
		id   func(d harness.Data, t *testing.T) string
		err  bool
	}{
		{
			name: "With ID",
			id: func(d harness.Data, t *testing.T) string {
				rec, err := d.GetGameCreatureRecByRef(harness.GameCreatureOneRef)
				require.NoError(t, err, "GetGameCreatureRecByRef returns without error")
				return rec.ID
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
			_, err = h.Setup()
			require.NoError(t, err, "Setup returns without error")
			defer func() {
				err = h.Teardown()
				require.NoError(t, err, "Teardown returns without error")
			}()

			r := h.Domain.(*domain.Domain).GameCreatureRepository()
			require.NotNil(t, r, "Repository is not nil")

			err = r.DeleteOne(tc.id(h.Data, t))
			if tc.err {
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
