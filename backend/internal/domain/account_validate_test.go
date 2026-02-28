package domain_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/harness"
	"gitlab.com/alienspaces/playbymail/internal/record/account_record"
	"gitlab.com/alienspaces/playbymail/internal/utils/config"
	"gitlab.com/alienspaces/playbymail/internal/utils/deps"
)

func TestCreateAccountUserRec_Validation(t *testing.T) {
	dataConfig := harness.DataConfig{}

	cfg, err := config.Parse()
	require.NoError(t, err)

	l, s, j, scanner, err := deps.NewDefaultDependencies(cfg)
	require.NoError(t, err)

	th, err := harness.NewTesting(cfg, l, s, j, scanner, dataConfig)
	require.NoError(t, err)

	th.ShouldCommitData = false

	_, err = th.Setup()
	require.NoError(t, err)
	defer func() {
		err = th.Teardown()
		require.NoError(t, err)
	}()

	testCases := []struct {
		name        string
		rec         *account_record.AccountUser
		expectError bool
		errorCode   coreerror.Code
	}{
		{
			name: "succeeds with valid email and active status",
			rec: &account_record.AccountUser{
				Email:  harness.UniqueEmail("valid@example.com"),
				Status: account_record.AccountUserStatusActive,
			},
			expectError: false,
		},
		{
			name: "succeeds with valid email and defaults to active status",
			rec: &account_record.AccountUser{
				Email: harness.UniqueEmail("default-status@example.com"),
			},
			expectError: false,
		},
		{
			name: "succeeds with pending_approval status",
			rec: &account_record.AccountUser{
				Email:  harness.UniqueEmail("pending@example.com"),
				Status: account_record.AccountUserStatusPendingApproval,
			},
			expectError: false,
		},
		{
			name:        "fails when record is nil",
			rec:         nil,
			expectError: true,
			errorCode:   coreerror.ErrorCodeInvalidData,
		},
		{
			name: "fails when email is empty",
			rec: &account_record.AccountUser{
				Email:  "",
				Status: account_record.AccountUserStatusActive,
			},
			expectError: true,
			errorCode:   coreerror.ErrorCodeInvalidData,
		},
		{
			name: "fails with invalid status",
			rec: &account_record.AccountUser{
				Email:  harness.UniqueEmail("badstatus@example.com"),
				Status: "nonexistent_status",
			},
			expectError: true,
			errorCode:   coreerror.ErrorCodeInvalidData,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := th.Domain.(*domain.Domain)

			rec, err := m.CreateAccountUserRec(tc.rec)

			if tc.expectError {
				require.Error(t, err)
				if tc.errorCode != "" {
					var coreErr coreerror.Error
					require.ErrorAs(t, err, &coreErr)
					require.Equal(t, tc.errorCode, coreErr.ErrorCode)
				}
				require.Nil(t, rec)
			} else {
				require.NoError(t, err)
				require.NotNil(t, rec)
				require.NotEmpty(t, rec.ID)
			}
		})
	}
}

func TestUpdateAccountUserRec_Validation(t *testing.T) {
	dataConfig := harness.DataConfig{
		AccountConfigs: []harness.AccountConfig{
			{
				Reference: "test-account",
				Record: &account_record.AccountUser{
					Email:  harness.UniqueEmail("update-test@example.com"),
					Status: account_record.AccountUserStatusActive,
				},
			},
		},
	}

	cfg, err := config.Parse()
	require.NoError(t, err)

	l, s, j, scanner, err := deps.NewDefaultDependencies(cfg)
	require.NoError(t, err)

	th, err := harness.NewTesting(cfg, l, s, j, scanner, dataConfig)
	require.NoError(t, err)

	th.ShouldCommitData = false

	_, err = th.Setup()
	require.NoError(t, err)
	defer func() {
		err = th.Teardown()
		require.NoError(t, err)
	}()

	m := th.Domain.(*domain.Domain)

	accountUser, err := th.Data.GetAccountUserRecByRef("test-account")
	require.NoError(t, err)

	t.Run("fails when email is changed", func(t *testing.T) {
		updateRec := &account_record.AccountUser{
			Email:  "different@example.com",
			Status: account_record.AccountUserStatusActive,
		}
		updateRec.ID = accountUser.ID

		_, err := m.UpdateAccountUserRec(updateRec)
		require.Error(t, err)
		var coreErr coreerror.Error
		require.ErrorAs(t, err, &coreErr)
		require.Equal(t, coreerror.ErrorCodeInvalidData, coreErr.ErrorCode)
	})

	t.Run("fails with invalid status", func(t *testing.T) {
		updateRec := &account_record.AccountUser{
			Email:  accountUser.Email,
			Status: "invalid_status",
		}
		updateRec.ID = accountUser.ID

		_, err := m.UpdateAccountUserRec(updateRec)
		require.Error(t, err)
	})
}
