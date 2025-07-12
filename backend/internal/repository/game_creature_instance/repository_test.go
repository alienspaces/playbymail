package game_creature_instance_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/harness"
	"gitlab.com/alienspaces/playbymail/internal/record"
	"gitlab.com/alienspaces/playbymail/internal/utils/config"
	"gitlab.com/alienspaces/playbymail/internal/utils/deps"
)

func TestGameCreatureInstanceRepository_CRUD(t *testing.T) {
	const (
		gameRef         = harness.GameOneRef
		gameInstanceRef = harness.GameInstanceOneRef
		gameCreatureRef = harness.GameCreatureOneRef
		gameLocationRef = harness.GameLocationOneRef
	)
	dcfg := harness.DataConfig{
		GameConfigs: []harness.GameConfig{
			{
				Reference: gameRef,
				Record:    &record.Game{},
				GameCreatureConfigs: []harness.GameCreatureConfig{
					{
						Reference: gameCreatureRef,
						Record:    &record.GameCreature{},
					},
				},
				GameLocationConfigs: []harness.GameLocationConfig{
					{
						Reference: gameLocationRef,
						Record:    &record.GameLocation{},
					},
				},
				GameInstanceConfigs: []harness.GameInstanceConfig{
					{
						Reference: gameInstanceRef,
						Record:    &record.GameInstance{},
						GameLocationInstanceConfigs: []harness.GameLocationInstanceConfig{
							{
								Reference:       gameLocationRef,
								GameLocationRef: gameLocationRef,
								Record:          &record.GameLocationInstance{},
							},
						},
						GameCreatureInstanceConfigs: []harness.GameCreatureInstanceConfig{
							{
								Reference:       harness.GameCreatureInstanceOneRef,
								GameCreatureRef: gameCreatureRef,
								GameLocationRef: gameLocationRef,
								Record:          &record.GameCreatureInstance{},
							},
						},
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

	_, err = h.Setup()
	require.NoError(t, err)
	defer func() {
		err = h.Teardown()
		require.NoError(t, err)
	}()

	repo := h.Domain.(*domain.Domain).GameCreatureInstanceRepository()
	require.NotNil(t, repo)

	rec, err := h.Data.GetGameCreatureInstanceRecByRef(harness.GameCreatureInstanceOneRef)
	require.NoError(t, err)

	// Update
	rec.IsAlive = false
	_, err = repo.UpdateOne(rec)
	require.NoError(t, err)

	got, err := repo.GetOne(rec.ID, nil)
	require.NoError(t, err)
	require.Equal(t, false, got.IsAlive)

	// Delete
	err = repo.DeleteOne(rec.ID)
	require.NoError(t, err)

	_, err = repo.GetOne(rec.ID, nil)
	require.Error(t, err)
}
