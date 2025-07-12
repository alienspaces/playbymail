package game_location_link_requirement_test

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

func newHarness(t *testing.T) *harness.Testing {
	dcfg := harness.DataConfig{
		GameConfigs: []harness.GameConfig{
			{
				Reference: harness.GameOneRef,
				Record: &record.Game{
					Name:     "Test Game",
					GameType: record.GameTypeAdventure,
				},
				GameItemConfigs: []harness.GameItemConfig{
					{
						Reference: harness.GameItemOneRef,
						Record: &record.GameItem{
							Name:        "Test Item",
							Description: "Test item for link requirement",
						},
					},
				},
				GameLocationConfigs: []harness.GameLocationConfig{
					{
						Reference: harness.GameLocationOneRef,
						Record: &record.GameLocation{
							Name:        "Test Location",
							Description: "Test location for link requirement",
						},
					},
				},
				GameLocationLinkConfigs: []harness.GameLocationLinkConfig{
					{
						Reference:       harness.GameLocationLinkOneRef,
						FromLocationRef: harness.GameLocationOneRef,
						ToLocationRef:   harness.GameLocationOneRef,
						Record: &record.GameLocationLink{
							Description: "Test link",
							Name:        "Test Link",
						},
						GameLocationLinkRequirementConfigs: []harness.GameLocationLinkRequirementConfig{
							{
								Reference:   harness.GameLocationLinkRequirementOneRef,
								GameItemRef: harness.GameItemOneRef,
								Record: &record.GameLocationLinkRequirement{
									Quantity: 1,
								},
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
	return h
}

func TestCreateOne(t *testing.T) {
	h := newHarness(t)
	tests := []struct {
		name   string
		rec    func(data *harness.Data, t *testing.T) *record.GameLocationLinkRequirement
		hasErr bool
	}{
		{
			name: "valid",
			rec: func(data *harness.Data, t *testing.T) *record.GameLocationLinkRequirement {
				game, err := data.GetGameRecByRef(harness.GameOneRef)
				require.NoError(t, err)
				link, err := data.GetGameLocationLinkRecByRef(harness.GameLocationLinkOneRef)
				require.NoError(t, err)
				item, err := data.GetGameItemRecByRef(harness.GameItemOneRef)
				require.NoError(t, err)
				return &record.GameLocationLinkRequirement{
					GameID:             game.ID,
					GameLocationLinkID: link.ID,
					GameItemID:         item.ID,
					Quantity:           1,
				}
			},
			hasErr: false,
		},
		{
			name: "missing foreign key",
			rec: func(data *harness.Data, t *testing.T) *record.GameLocationLinkRequirement {
				return &record.GameLocationLinkRequirement{
					GameID:             uuid.NewString(),
					GameLocationLinkID: uuid.NewString(),
					GameItemID:         uuid.NewString(),
					Quantity:           1,
				}
			},
			hasErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := h.Setup()
			require.NoError(t, err)
			defer func() { _ = h.Teardown() }()
			r := h.Domain.(*domain.Domain).GameLocationLinkRequirementRepository()
			rec := tt.rec(&h.Data, t)
			created, err := r.CreateOne(rec)
			if tt.hasErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotEmpty(t, created.ID)
			}
		})
	}
}

func TestGetOne(t *testing.T) {
	h := newHarness(t)
	tests := []struct {
		name   string
		id     func(data *harness.Data, t *testing.T) string
		hasErr bool
	}{
		{
			name: "valid",
			id: func(data *harness.Data, t *testing.T) string {
				rec, err := data.GetGameLocationLinkRequirementRecByRef(harness.GameLocationLinkRequirementOneRef)
				require.NoError(t, err)
				return rec.ID
			},
			hasErr: false,
		},
		{
			name: "not found",
			id: func(data *harness.Data, t *testing.T) string {
				return uuid.NewString()
			},
			hasErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := h.Setup()
			require.NoError(t, err)
			defer func() { _ = h.Teardown() }()
			r := h.Domain.(*domain.Domain).GameLocationLinkRequirementRepository()
			id := tt.id(&h.Data, t)
			rec, err := r.GetOne(id, nil)
			if tt.hasErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, rec)
			}
		})
	}
}

func TestUpdateOne(t *testing.T) {
	h := newHarness(t)
	tests := []struct {
		name   string
		rec    func(data *harness.Data, t *testing.T) *record.GameLocationLinkRequirement
		hasErr bool
	}{
		{
			name: "valid",
			rec: func(data *harness.Data, t *testing.T) *record.GameLocationLinkRequirement {
				rec, err := data.GetGameLocationLinkRequirementRecByRef(harness.GameLocationLinkRequirementOneRef)
				require.NoError(t, err)
				rec.Quantity = 2 // simulate update
				return rec
			},
			hasErr: false,
		},
		{
			name: "not found",
			rec: func(data *harness.Data, t *testing.T) *record.GameLocationLinkRequirement {
				rec, err := data.GetGameLocationLinkRequirementRecByRef(harness.GameLocationLinkRequirementOneRef)
				require.NoError(t, err)
				rec.ID = uuid.NewString()
				return rec
			},
			hasErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := h.Setup()
			require.NoError(t, err)
			defer func() { _ = h.Teardown() }()
			r := h.Domain.(*domain.Domain).GameLocationLinkRequirementRepository()
			rec := tt.rec(&h.Data, t)
			updated, err := r.UpdateOne(rec)
			if tt.hasErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, rec.Quantity, updated.Quantity)
			}
		})
	}
}

func TestDeleteOne(t *testing.T) {
	h := newHarness(t)
	tests := []struct {
		name   string
		id     func(data *harness.Data, t *testing.T) string
		hasErr bool
	}{
		{
			name: "valid",
			id: func(data *harness.Data, t *testing.T) string {
				rec, err := data.GetGameLocationLinkRequirementRecByRef(harness.GameLocationLinkRequirementOneRef)
				require.NoError(t, err)
				return rec.ID
			},
			hasErr: false,
		},
		{
			name: "not found",
			id: func(data *harness.Data, t *testing.T) string {
				return uuid.NewString()
			},
			hasErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := h.Setup()
			require.NoError(t, err)
			defer func() { _ = h.Teardown() }()
			r := h.Domain.(*domain.Domain).GameLocationLinkRequirementRepository()
			id := tt.id(&h.Data, t)
			err = r.DeleteOne(id)
			if tt.hasErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
