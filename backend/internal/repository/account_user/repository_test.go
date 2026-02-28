package account_user_test

import (
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/harness"
	"gitlab.com/alienspaces/playbymail/internal/record/account_record"
	"gitlab.com/alienspaces/playbymail/internal/utils/deps"
)

func TestCreateOne(t *testing.T) {
	h := deps.NewHarness(t)

	tests := []struct {
		name   string
		rec    func(d harness.Data, t *testing.T) *account_record.AccountUser
		hasErr bool
	}{
		{
			name: "Without ID",
			rec: func(d harness.Data, t *testing.T) *account_record.AccountUser {
				existing, err := d.GetAccountUserRecByRef(harness.StandardAccountRef)
				require.NoError(t, err)

				return &account_record.AccountUser{
					AccountID: existing.AccountID,
					Email:     fmt.Sprintf("%s@example.com", gofakeit.Name()),
					Status:    account_record.AccountUserStatusActive,
				}
			},
			hasErr: false,
		},
		{
			name: "With ID",
			rec: func(d harness.Data, t *testing.T) *account_record.AccountUser {
				existing, err := d.GetAccountUserRecByRef(harness.StandardAccountRef)
				require.NoError(t, err)

				rec := &account_record.AccountUser{
					AccountID: existing.AccountID,
					Email:     fmt.Sprintf("%s@example.com", gofakeit.Name()),
					Status:    account_record.AccountUserStatusActive,
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

			// Test harness
			_, err := h.Setup()
			require.NoError(t, err, "Setup returns without error")
			defer func() {
				err = h.Teardown()
				require.NoError(t, err, "Teardown returns without error")
			}()

			// repository
			r := h.Domain.(*domain.Domain).AccountUserRepository()
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
				accRec, err := d.GetAccountUserRecByRef(harness.StandardAccountRef)
				require.NoError(t, err, "GetAccountUserRecByRef returns without error")
				return accRec.ID
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

			// harness setup
			_, err := h.Setup()
			require.NoError(t, err, "Setup returns without error")
			defer func() {
				err = h.Teardown()
				require.NoError(t, err, "Teardown returns without error")
			}()

			// repository
			r := h.Domain.(*domain.Domain).AccountUserRepository()
			require.NotNil(t, r, "Repository is not nil")

			rec, err := r.GetOne(tc.id(h.Data, t), nil)
			if tc.hasErr {
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
	h := deps.NewHarness(t)

	tests := []struct {
		name   string
		rec    func(d harness.Data, t *testing.T) *account_record.AccountUser
		hasErr bool
	}{
		{
			name: "With ID",
			rec: func(d harness.Data, t *testing.T) *account_record.AccountUser {
				accRec, err := d.GetAccountUserRecByRef(harness.StandardAccountRef)
				require.NoError(t, err, "GetAccountUserRecByRef returns without error")
				return accRec
			},
			hasErr: false,
		},
		{
			name: "Without ID",
			rec: func(d harness.Data, t *testing.T) *account_record.AccountUser {
				accRec, err := d.GetAccountUserRecByRef(harness.StandardAccountRef)
				require.NoError(t, err, "GetAccountUserRecByRef returns without error")
				rec := accRec
				rec.ID = ""
				return rec
			},
			hasErr: true,
		},
	}

	for _, tc := range tests {

		t.Logf("Run test >%s<", tc.name)

		t.Run(tc.name, func(t *testing.T) {

			// harness setup
			_, err := h.Setup()
			require.NoError(t, err, "Setup returns without error")
			defer func() {
				err = h.Teardown()
				require.NoError(t, err, "Teardown returns without error")
			}()

			// repository
			r := h.Domain.(*domain.Domain).AccountUserRepository()
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
				accRec, err := d.GetAccountUserRecByRef(harness.StandardAccountRef)
				require.NoError(t, err, "GetAccountUserRecByRef returns without error")
				return accRec.ID
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

			// harness setup
			_, err := h.Setup()
			require.NoError(t, err, "Setup returns without error")
			defer func() {
				err = h.Teardown()
				require.NoError(t, err, "Teardown returns without error")
			}()

			// repository
			r := h.Domain.(*domain.Domain).AccountUserRepository()
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
