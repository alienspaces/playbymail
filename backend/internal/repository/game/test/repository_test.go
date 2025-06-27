package test

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
	"gitlab.com/alienspaces/playbymail/internal/record"
	"gitlab.com/alienspaces/playbymail/internal/utils/config"
	"gitlab.com/alienspaces/playbymail/internal/utils/deps"
)

func TestCreateOne(t *testing.T) {

	// harness
	dcfg := harness.DataConfig{}

	cfg, err := config.Parse()
	require.NoError(t, err, "Parse returns without error")

	l, s, j, err := deps.Default(cfg)
	require.NoError(t, err, "Default dependencies returns without error")

	h, err := harness.NewTesting(l, s, j, dcfg)
	require.NoError(t, err, "NewTesting returns without error")

	tests := []struct {
		name string
		rec  func() *record.Game
		err  bool
	}{
		{
			name: "Without ID",
			rec: func() *record.Game {
				return &record.Game{
					Name: fmt.Sprintf("%s %s", gofakeit.Name(), gofakeit.Name()),
				}
			},
			err: false,
		},
		{
			name: "With ID",
			rec: func() *record.Game {
				rec := &record.Game{
					Name: fmt.Sprintf("%s %s", gofakeit.Name(), gofakeit.Name()),
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
			_, err = h.Setup()
			require.NoError(t, err, "Setup returns without error")
			defer func() {
				err = h.Teardown()
				require.NoError(t, err, "Teardown returns without error")
			}()

			// repository
			r := h.Domain.(*domain.Domain).GameRepository()
			require.NotNil(t, r, "Repository is not nil")

			rec := tc.rec()

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

	// harness
	dcfg := harness.DataConfig{
		GameConfig: []harness.GameConfig{
			{
				Record: &record.Game{},
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
		id   func() string
		err  bool
	}{
		{
			name: "With ID",
			id: func() string {
				return h.Data.GameRecs[0].ID
			},
			err: false,
		},
		{
			name: "Without ID",
			id: func() string {
				return ""
			},
			err: true,
		},
	}

	for _, tc := range tests {

		t.Logf("Run test >%s<", tc.name)

		t.Run(tc.name, func(t *testing.T) {

			// harness setup
			_, err = h.Setup()
			require.NoError(t, err, "Setup returns without error")
			defer func() {
				err = h.Teardown()
				require.NoError(t, err, "Teardown returns without error")
			}()

			// repository
			r := h.Domain.(*domain.Domain).GameRepository()
			require.NotNil(t, r, "Repository is not nil")

			rec, err := r.GetOne(tc.id(), nil)
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

	// harness
	dcfg := harness.DataConfig{
		GameConfig: []harness.GameConfig{
			{
				Record: &record.Game{},
			},
		},
	}

	cfg, err := config.Parse()
	require.NoError(t, err, "Parse returns without error")

	l, s, j, err := deps.Default(cfg)
	require.NoError(t, err, "Default dependencies returns without error")

	h, err := harness.NewTesting(l, s, j, dcfg)
	require.NoError(t, err, "NewTesting returns without error")

	require.NoError(t, err, "NewTesting returns without error")

	tests := []struct {
		name string
		rec  func() *record.Game
		err  bool
	}{
		{
			name: "With ID",
			rec: func() *record.Game {
				return h.Data.GameRecs[0]
			},
			err: false,
		},
		{
			name: "Without ID",
			rec: func() *record.Game {
				rec := h.Data.GameRecs[0]
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
			_, err = h.Setup()
			require.NoError(t, err, "Setup returns without error")
			defer func() {
				err = h.Teardown()
				require.NoError(t, err, "Teardown returns without error")
			}()

			// repository
			r := h.Domain.(*domain.Domain).GameRepository()
			require.NotNil(t, r, "Repository is not nil")

			rec := tc.rec()

			_, err := r.UpdateOne(rec)
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

	// harness
	dcfg := harness.DataConfig{
		GameConfig: []harness.GameConfig{
			{
				Record: &record.Game{},
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
		id   func() string
		err  bool
	}{
		{
			name: "With ID",
			id: func() string {
				return h.Data.GameRecs[0].ID
			},
			err: false,
		},
		{
			name: "Without ID",
			id: func() string {
				return ""
			},
			err: true,
		},
	}

	for _, tc := range tests {

		t.Logf("Run test >%s<", tc.name)

		t.Run(tc.name, func(t *testing.T) {

			// harness setup
			_, err = h.Setup()
			require.NoError(t, err, "Setup returns without error")
			defer func() {
				err = h.Teardown()
				require.NoError(t, err, "Teardown returns without error")
			}()

			// repository
			r := h.Domain.(*domain.Domain).GameRepository()
			require.NotNil(t, r, "Repository is not nil")

			err = r.DeleteOne(tc.id())
			if tc.err == true {
				require.Error(t, err, "DeleteOne returns error")
				return
			}
			require.NoError(t, err, "DeleteOne returns without error")

			rec, err := r.GetOne(tc.id(), nil)
			require.Error(t, err, "GetOne returns error")
			require.Nil(t, rec, "GetOne does not return record")
		})
	}
}
