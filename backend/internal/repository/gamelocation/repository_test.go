package gamelocation_test

import (
	"fmt"
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
		rec  func(d harness.Data, t *testing.T) *record.GameLocation
		err  bool
	}{
		{
			name: "Without ID",
			rec: func(d harness.Data, t *testing.T) *record.GameLocation {
				gameRec, err := d.GetGameRecByRef(harness.GameOneRef)
				require.NoError(t, err, "GetGameRecByRef returns without error")
				return &record.GameLocation{
					GameID:      gameRec.ID,
					Name:        fmt.Sprintf("%s %s", gofakeit.Name(), gofakeit.Name()),
					Description: gofakeit.Sentence(5),
				}
			},
			err: false,
		},
		{
			name: "With ID",
			rec: func(d harness.Data, t *testing.T) *record.GameLocation {
				gameRec, err := d.GetGameRecByRef(harness.GameOneRef)
				require.NoError(t, err, "GetGameRecByRef returns without error")
				rec := &record.GameLocation{
					GameID:      gameRec.ID,
					Name:        fmt.Sprintf("%s %s", gofakeit.Name(), gofakeit.Name()),
					Description: gofakeit.Sentence(5),
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

			r := h.Domain.(*domain.Domain).GameLocationRepository()
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
	dcfg := harness.DataConfig{
		GameConfigs: []harness.GameConfig{
			{
				Reference: harness.GameOneRef,
				Record:    &record.Game{},
				GameLocationConfigs: []harness.GameLocationConfig{
					{
						Reference: harness.GameLocationOneRef,
						Record: &record.GameLocation{
							Name:        "Test Location",
							Description: "Test Description",
						},
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
				gameLocationRec, err := d.GetGameLocationRecByRef(harness.GameLocationOneRef)
				require.NoError(t, err, "GetGameLocationRecByRef returns without error")
				return gameLocationRec.ID
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

			r := h.Domain.(*domain.Domain).GameLocationRepository()
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
	dcfg := harness.DataConfig{
		GameConfigs: []harness.GameConfig{
			{
				Reference: harness.GameOneRef,
				Record:    &record.Game{},
				GameLocationConfigs: []harness.GameLocationConfig{
					{
						Reference: harness.GameLocationOneRef,
						Record: &record.GameLocation{
							Name:        "Test Location",
							Description: "Test Description",
						},
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
		rec  func(d harness.Data, t *testing.T) *record.GameLocation
		err  bool
	}{
		{
			name: "With ID",
			rec: func(d harness.Data, t *testing.T) *record.GameLocation {
				gameLocationRec, err := d.GetGameLocationRecByRef(harness.GameLocationOneRef)
				require.NoError(t, err, "GetGameLocationRecByRef returns without error")
				return gameLocationRec
			},
			err: false,
		},
		{
			name: "Without ID",
			rec: func(d harness.Data, t *testing.T) *record.GameLocation {
				return &record.GameLocation{}
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

			r := h.Domain.(*domain.Domain).GameLocationRepository()
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
	dcfg := harness.DataConfig{
		GameConfigs: []harness.GameConfig{
			{
				Reference: harness.GameOneRef,
				Record:    &record.Game{},
				GameLocationConfigs: []harness.GameLocationConfig{
					{
						Reference: harness.GameLocationOneRef,
						Record: &record.GameLocation{
							Name:        "Test Location",
							Description: "Test Description",
						},
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
				gameLocationRec, err := d.GetGameLocationRecByRef(harness.GameLocationOneRef)
				require.NoError(t, err, "GetGameLocationRecByRef returns without error")
				return gameLocationRec.ID
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

			r := h.Domain.(*domain.Domain).GameLocationRepository()
			require.NotNil(t, r, "Repository is not nil")

			err = r.DeleteOne(tc.id(h.Data, t))
			if tc.err == true {
				require.Error(t, err, "DeleteOne returns error")
				return
			}
			require.NoError(t, err, "DeleteOne returns without error")
		})
	}
}
