package adventure_game_item_instance_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/playbymail/core/nullstring"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/harness"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
	"gitlab.com/alienspaces/playbymail/internal/utils/deps"
)

func TestCreateOne(t *testing.T) {
	h := deps.NewHarness(t)

	tests := []struct {
		name   string
		rec    func(d harness.Data, t *testing.T) *adventure_game_record.AdventureGameItemInstance
		hasErr bool
	}{
		{
			name: "Without ID",
			rec: func(d harness.Data, t *testing.T) *adventure_game_record.AdventureGameItemInstance {
				gameRec, err := d.GetGameRecByRef(harness.GameOneRef)
				require.NoError(t, err, "GetGameRecByRef returns without error")
				itemRec, err := d.GetGameItemRecByRef(harness.GameItemOneRef)
				require.NoError(t, err, "GetGameItemRecByRef returns without error")
				gameInstanceRec, err := d.GetGameInstanceRecByRef(harness.GameInstanceOneRef)
				require.NoError(t, err, "GetGameInstanceRecByRef returns without error")
				locationInstanceRec, err := d.GetGameLocationInstanceRecByRef(harness.GameLocationInstanceOneRef)
				require.NoError(t, err, "GetGameLocationRecByRef returns without error")
				return &adventure_game_record.AdventureGameItemInstance{
					GameID:                          gameRec.ID,
					AdventureGameItemID:             itemRec.ID,
					AdventureGameInstanceID:         gameInstanceRec.ID,
					AdventureGameLocationInstanceID: nullstring.FromString(locationInstanceRec.ID),
				}
			},
			hasErr: false,
		},
		{
			name: "With ID",
			rec: func(d harness.Data, t *testing.T) *adventure_game_record.AdventureGameItemInstance {
				gameRec, err := d.GetGameRecByRef(harness.GameOneRef)
				require.NoError(t, err, "GetGameRecByRef returns without error")
				itemRec, err := d.GetGameItemRecByRef(harness.GameItemOneRef)
				require.NoError(t, err, "GetGameItemRecByRef returns without error")
				gameInstanceRec, err := d.GetGameInstanceRecByRef(harness.GameInstanceOneRef)
				require.NoError(t, err, "GetGameInstanceRecByRef returns without error")
				locationInstanceRec, err := d.GetGameLocationInstanceRecByRef(harness.GameLocationInstanceOneRef)
				require.NoError(t, err, "GetGameLocationRecByRef returns without error")
				rec := &adventure_game_record.AdventureGameItemInstance{
					GameID:                          gameRec.ID,
					AdventureGameItemID:             itemRec.ID,
					AdventureGameInstanceID:         gameInstanceRec.ID,
					AdventureGameLocationInstanceID: nullstring.FromString(locationInstanceRec.ID),
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

			r := h.Domain.(*domain.Domain).AdventureGameItemInstanceRepository()
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
				rec, err := h.Data.GetGameItemInstanceRecByRef(harness.GameItemInstanceOneRef)
				require.NoError(t, err, "getGameItemInstanceRecByRef returns without error")
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

			r := h.Domain.(*domain.Domain).AdventureGameItemInstanceRepository()
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
		rec    func(d harness.Data, t *testing.T) *adventure_game_record.AdventureGameItemInstance
		hasErr bool
	}{
		{
			name: "With ID",
			rec: func(d harness.Data, t *testing.T) *adventure_game_record.AdventureGameItemInstance {
				rec, err := h.Data.GetGameItemInstanceRecByRef(harness.GameItemInstanceOneRef)
				require.NoError(t, err, "getGameItemInstanceRecByRef returns without error")
				return rec
			},
			hasErr: false,
		},
		{
			name: "Without ID",
			rec: func(d harness.Data, t *testing.T) *adventure_game_record.AdventureGameItemInstance {
				rec, err := h.Data.GetGameItemInstanceRecByRef(harness.GameItemInstanceOneRef)
				require.NoError(t, err, "getGameItemInstanceRecByRef returns without error")
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

			r := h.Domain.(*domain.Domain).AdventureGameItemInstanceRepository()
			require.NotNil(t, r, "Repository is not nil")

			rec := tc.rec(h.Data, t)

			updated, err := r.UpdateOne(rec)
			if tc.hasErr {
				require.Error(t, err, "UpdateOne returns error")
				return
			}
			require.NoError(t, err, "UpdateOne returns without error")
			require.NotEmpty(t, updated.UpdatedAt, "UpdateOne returns record with UpdatedAt")
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
				rec, err := h.Data.GetGameItemInstanceRecByRef(harness.GameItemInstanceOneRef)
				require.NoError(t, err, "getGameItemInstanceRecByRef returns without error")
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

			r := h.Domain.(*domain.Domain).AdventureGameItemInstanceRepository()
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
