package adventure_game_character_instance_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/harness"
	adventure_game_record "gitlab.com/alienspaces/playbymail/internal/record/adventure_game"
	"gitlab.com/alienspaces/playbymail/internal/utils/deps"
)

func TestCreateOne(t *testing.T) {
	h := deps.NewHarness(t)

	tests := []struct {
		name   string
		rec    func(d harness.Data, t *testing.T) *adventure_game_record.AdventureGameCharacterInstance
		hasErr bool
	}{
		{
			name: "Without ID",
			rec: func(d harness.Data, t *testing.T) *adventure_game_record.AdventureGameCharacterInstance {
				gameInstanceRec, err := d.GetGameInstanceRecByRef(harness.GameInstanceOneRef)
				require.NoError(t, err)
				charRec, err := d.GetGameCharacterRecByRef(harness.GameCharacterTwoRef)
				require.NoError(t, err)
				locationInstanceRec, err := d.GetGameLocationInstanceRecByRef(harness.GameLocationInstanceOneRef)
				require.NoError(t, err)
				return &adventure_game_record.AdventureGameCharacterInstance{
					GameID:                          gameInstanceRec.GameID,
					AdventureGameInstanceID:         gameInstanceRec.ID,
					AdventureGameCharacterID:        charRec.ID,
					AdventureGameLocationInstanceID: locationInstanceRec.ID,
					Health:                          100,
				}
			},
			hasErr: false,
		},
		{
			name: "With ID",
			rec: func(d harness.Data, t *testing.T) *adventure_game_record.AdventureGameCharacterInstance {
				gameInstanceRec, err := d.GetGameInstanceRecByRef(harness.GameInstanceOneRef)
				require.NoError(t, err)
				charRec, err := d.GetGameCharacterRecByRef(harness.GameCharacterTwoRef)
				require.NoError(t, err)
				locationInstanceRec, err := d.GetGameLocationInstanceRecByRef(harness.GameLocationInstanceOneRef)
				require.NoError(t, err)
				rec := &adventure_game_record.AdventureGameCharacterInstance{
					GameID:                          gameInstanceRec.GameID,
					AdventureGameInstanceID:         gameInstanceRec.ID,
					AdventureGameCharacterID:        charRec.ID,
					AdventureGameLocationInstanceID: locationInstanceRec.ID,
					Health:                          100,
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
			require.NoError(t, err)
			defer func() {
				err = h.Teardown()
				require.NoError(t, err)
			}()

			r := h.Domain.(*domain.Domain).AdventureGameCharacterInstanceRepository()
			require.NotNil(t, r)

			rec := tc.rec(h.Data, t)

			_, err = r.CreateOne(rec)
			if tc.hasErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotEmpty(t, rec.CreatedAt)
			// Health should be 100
			require.Equal(t, 100, rec.Health)
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
				rec, err := h.Data.GetGameCharacterInstanceRecByRef(harness.GameCharacterInstanceOneRef) // Add a reference if available
				require.NoError(t, err)
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
			require.NoError(t, err)
			defer func() {
				err = h.Teardown()
				require.NoError(t, err)
			}()

			r := h.Domain.(*domain.Domain).AdventureGameCharacterInstanceRepository()
			require.NotNil(t, r)

			rec, err := r.GetOne(tc.id(h.Data, t), nil)
			if tc.hasErr {
				require.Error(t, err)
				require.Nil(t, rec)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, rec)
			require.NotEmpty(t, rec.ID)
			// Health should be 100
			require.Equal(t, 100, rec.Health)
		})
	}
}

func TestUpdateOne(t *testing.T) {
	h := deps.NewHarness(t)

	tests := []struct {
		name   string
		rec    func(d harness.Data, t *testing.T) *adventure_game_record.AdventureGameCharacterInstance
		hasErr bool
	}{
		{
			name: "With ID",
			rec: func(d harness.Data, t *testing.T) *adventure_game_record.AdventureGameCharacterInstance {
				rec, err := h.Data.GetGameCharacterInstanceRecByRef(harness.GameCharacterInstanceOneRef) // Add a reference if available
				require.NoError(t, err)
				return rec
			},
			hasErr: false,
		},
		{
			name: "Without ID",
			rec: func(d harness.Data, t *testing.T) *adventure_game_record.AdventureGameCharacterInstance {
				rec, err := h.Data.GetGameCharacterInstanceRecByRef(harness.GameCharacterInstanceOneRef) // Add a reference if available
				require.NoError(t, err)
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
			require.NoError(t, err)
			defer func() {
				err = h.Teardown()
				require.NoError(t, err)
			}()

			r := h.Domain.(*domain.Domain).AdventureGameCharacterInstanceRepository()
			require.NotNil(t, r)

			rec := tc.rec(h.Data, t)
			rec.Health = 75

			updated, err := r.UpdateOne(rec)
			if tc.hasErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotEmpty(t, updated.UpdatedAt)
			require.Equal(t, 75, updated.Health)
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
				rec, err := h.Data.GetGameCharacterInstanceRecByRef(harness.GameCharacterInstanceOneRef) // Add a reference if available
				require.NoError(t, err)
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
			require.NoError(t, err)
			defer func() {
				err = h.Teardown()
				require.NoError(t, err)
			}()

			r := h.Domain.(*domain.Domain).AdventureGameCharacterInstanceRepository()
			require.NotNil(t, r)

			err = r.DeleteOne(tc.id(h.Data, t))
			if tc.hasErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			rec, err := r.GetOne(tc.id(h.Data, t), nil)
			require.Error(t, err)
			require.Nil(t, rec)
		})
	}
}
