package domain_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/core/record"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/harness"
	"gitlab.com/alienspaces/playbymail/internal/record/account_record"
	"gitlab.com/alienspaces/playbymail/internal/utils/config"
	"gitlab.com/alienspaces/playbymail/internal/utils/deps"
)

func TestCreateAccountUserRec_Validation(t *testing.T) {
	dataConfig := harness.DataConfig{
		AccountConfigs: []harness.AccountConfig{
			{
				Reference: "existing-account",
				Record: &account_record.AccountUser{
					Email:  harness.UniqueEmail("existing@example.com"),
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

	existingAccount, err := th.Data.GetAccountUserRecByRef("existing-account")
	require.NoError(t, err)

	m := th.Domain.(*domain.Domain)

	t.Run("succeeds with valid email and active status", func(t *testing.T) {
		rec := &account_record.AccountUser{
			AccountID: existingAccount.AccountID,
			Email:     harness.UniqueEmail("valid@example.com"),
			Status:    account_record.AccountUserStatusActive,
		}
		result, err := m.CreateAccountUserRec(rec)
		require.NoError(t, err)
		require.NotNil(t, result)
		require.NotEmpty(t, result.ID)
	})

	t.Run("succeeds with pending_approval status", func(t *testing.T) {
		rec := &account_record.AccountUser{
			AccountID: existingAccount.AccountID,
			Email:     harness.UniqueEmail("pending@example.com"),
			Status:    account_record.AccountUserStatusPendingApproval,
		}
		result, err := m.CreateAccountUserRec(rec)
		require.NoError(t, err)
		require.NotNil(t, result)
	})

	t.Run("fails when record is nil", func(t *testing.T) {
		result, err := m.CreateAccountUserRec(nil)
		require.Error(t, err)
		var coreErr coreerror.Error
		require.ErrorAs(t, err, &coreErr)
		require.Equal(t, coreerror.ErrorCodeInvalidData, coreErr.ErrorCode)
		require.Nil(t, result)
	})

	t.Run("fails when email is empty", func(t *testing.T) {
		rec := &account_record.AccountUser{
			AccountID: existingAccount.AccountID,
			Email:     "",
			Status:    account_record.AccountUserStatusActive,
		}
		result, err := m.CreateAccountUserRec(rec)
		require.Error(t, err)
		var coreErr coreerror.Error
		require.ErrorAs(t, err, &coreErr)
		require.Equal(t, coreerror.ErrorCodeInvalidData, coreErr.ErrorCode)
		require.NotNil(t, result)
		require.Empty(t, result.ID)
	})

	t.Run("fails with invalid account_id", func(t *testing.T) {
		rec := &account_record.AccountUser{
			AccountID: "not-a-uuid",
			Email:     harness.UniqueEmail("bad-account@example.com"),
			Status:    account_record.AccountUserStatusActive,
		}
		result, err := m.CreateAccountUserRec(rec)
		require.Error(t, err)
		require.NotNil(t, result)
		require.Empty(t, result.ID)
	})

	t.Run("fails with invalid status", func(t *testing.T) {
		rec := &account_record.AccountUser{
			AccountID: record.NewRecordID(),
			Email:     harness.UniqueEmail("badstatus@example.com"),
			Status:    "nonexistent_status",
		}
		result, err := m.CreateAccountUserRec(rec)
		require.Error(t, err)
		var coreErr coreerror.Error
		require.ErrorAs(t, err, &coreErr)
		require.Equal(t, coreerror.ErrorCodeInvalidData, coreErr.ErrorCode)
		require.NotNil(t, result)
		require.Empty(t, result.ID)
	})
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
