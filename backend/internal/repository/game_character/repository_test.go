package game_character_test

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

func newHarness(t *testing.T) *harness.Testing {
	dcfg := harness.DataConfig{
		AccountConfigs: []harness.AccountConfig{{
			Reference: harness.AccountOneRef,
			Record:    &record.Account{Name: "Test Account", Email: gofakeit.Email()},
		}},
		GameConfigs: []harness.GameConfig{{
			Reference: harness.GameOneRef,
			Record:    &record.Game{Name: "Test Game"},
			GameCharacterConfigs: []harness.GameCharacterConfig{{
				Reference:  harness.GameCharacterOneRef,
				AccountRef: harness.AccountOneRef,
				Record:     &record.GameCharacter{Name: "Test Character"},
			}},
		}},
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
	h.ShouldCommitData = false
	tests := []struct {
		name   string
		rec    func(d harness.Data, t *testing.T) *record.GameCharacter
		hasErr bool
	}{
		{
			name: "valid",
			rec: func(d harness.Data, t *testing.T) *record.GameCharacter {
				gameRec, err := d.GetGameRecByRef(harness.GameOneRef)
				require.NoError(t, err)
				accountRec, err := d.GetAccountRecByRef(harness.AccountOneRef)
				require.NoError(t, err)
				return &record.GameCharacter{
					GameID:    gameRec.ID,
					AccountID: accountRec.ID,
					Name:      "New Character",
				}
			},
			hasErr: false,
		},
		{
			name: "missing name",
			rec: func(d harness.Data, t *testing.T) *record.GameCharacter {
				gameRec, err := d.GetGameRecByRef(harness.GameOneRef)
				require.NoError(t, err)
				accountRec, err := d.GetAccountRecByRef(harness.AccountOneRef)
				require.NoError(t, err)
				return &record.GameCharacter{
					GameID:    gameRec.ID,
					AccountID: accountRec.ID,
					Name:      "",
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
			r := h.Domain.(*domain.Domain).GameCharacterRepository()
			rec := tt.rec(h.Data, t)
			_, err = r.CreateOne(rec)
			if tt.hasErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotEmpty(t, rec.ID)
			}
		})
	}
}

func TestGetOne(t *testing.T) {
	h := newHarness(t)
	tests := []struct {
		name   string
		id     func(d harness.Data, t *testing.T) string
		hasErr bool
	}{
		{
			name: "existing",
			id: func(d harness.Data, t *testing.T) string {
				rec, err := d.GetGameCharacterRecByRef(harness.GameCharacterOneRef)
				require.NoError(t, err)
				return rec.ID
			},
			hasErr: false,
		},
		{
			name: "not found",
			id: func(d harness.Data, t *testing.T) string {
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
			r := h.Domain.(*domain.Domain).GameCharacterRepository()
			rec, err := r.GetOne(tt.id(h.Data, t), nil)
			if tt.hasErr {
				require.Error(t, err)
				require.Nil(t, rec)
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
		rec    func(d harness.Data, t *testing.T) *record.GameCharacter
		hasErr bool
	}{
		{
			name: "valid update",
			rec: func(d harness.Data, t *testing.T) *record.GameCharacter {
				rec, err := d.GetGameCharacterRecByRef(harness.GameCharacterOneRef)
				require.NoError(t, err)
				rec.Name = "Updated Name"
				return rec
			},
			hasErr: false,
		},
		{
			name: "missing name",
			rec: func(d harness.Data, t *testing.T) *record.GameCharacter {
				rec, err := d.GetGameCharacterRecByRef(harness.GameCharacterOneRef)
				require.NoError(t, err)
				rec.Name = ""
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
			r := h.Domain.(*domain.Domain).GameCharacterRepository()
			rec := tt.rec(h.Data, t)
			updated, err := r.UpdateOne(rec)
			if tt.hasErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, rec.Name, updated.Name)
			}
		})
	}
}

func TestDeleteOne(t *testing.T) {
	h := newHarness(t)
	tests := []struct {
		name   string
		id     func(d harness.Data, t *testing.T) string
		hasErr bool
	}{
		{
			name: "existing",
			id: func(d harness.Data, t *testing.T) string {
				rec, err := d.GetGameCharacterRecByRef(harness.GameCharacterOneRef)
				require.NoError(t, err)
				return rec.ID
			},
			hasErr: false,
		},
		{
			name: "not found",
			id: func(d harness.Data, t *testing.T) string {
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
			r := h.Domain.(*domain.Domain).GameCharacterRepository()
			err = r.DeleteOne(tt.id(h.Data, t))
			if tt.hasErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				rec, err := r.GetOne(tt.id(h.Data, t), nil)
				require.Error(t, err)
				require.Nil(t, rec)
			}
		})
	}
}
