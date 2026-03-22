package domain_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/harness"
	"gitlab.com/alienspaces/playbymail/internal/record/account_record"
	"gitlab.com/alienspaces/playbymail/internal/utils/config"
	"gitlab.com/alienspaces/playbymail/internal/utils/deps"
)

func TestVerifyAccountUserVerificationToken_StatusPromotion(t *testing.T) {
	cfg, err := config.Parse()
	require.NoError(t, err)

	l, s, j, scanner, err := deps.NewDefaultDependencies(cfg)
	require.NoError(t, err)

	th, err := harness.NewTesting(cfg, l, s, j, scanner, harness.DataConfig{})
	require.NoError(t, err)

	th.ShouldCommitData = false

	_, err = th.Setup()
	require.NoError(t, err)
	defer func() {
		err = th.Teardown()
		require.NoError(t, err)
	}()

	m := th.Domain.(*domain.Domain)

	t.Run("pending_approval user is promoted to active after verification", func(t *testing.T) {
		_, accountUserRec, _, _, err := m.UpsertAccount(
			&account_record.Account{},
			&account_record.AccountUser{
				Email: harness.UniqueEmail("pending-verify@example.com"),
				// Status intentionally empty — CreateAccountUserRec defaults to pending_approval
			},
			nil,
		)
		require.NoError(t, err)
		require.Equal(t, account_record.AccountUserStatusPendingApproval, accountUserRec.Status, "new user starts pending_approval")

		token, err := m.GenerateAccountUserVerificationToken(accountUserRec)
		require.NoError(t, err)
		require.NotEmpty(t, token)

		_, err = m.VerifyAccountUserVerificationToken(token, false)
		require.NoError(t, err)

		updated, err := m.GetAccountUserRec(accountUserRec.ID, nil)
		require.NoError(t, err)
		require.Equal(t, account_record.AccountUserStatusActive, updated.Status, "user should be active after verification")
	})

	t.Run("already active user remains active after verification", func(t *testing.T) {
		_, accountUserRec, _, _, err := m.UpsertAccount(
			&account_record.Account{},
			&account_record.AccountUser{
				Email:  harness.UniqueEmail("active-verify@example.com"),
				Status: account_record.AccountUserStatusActive,
			},
			nil,
		)
		require.NoError(t, err)
		require.Equal(t, account_record.AccountUserStatusActive, accountUserRec.Status, "user starts active")

		token, err := m.GenerateAccountUserVerificationToken(accountUserRec)
		require.NoError(t, err)

		_, err = m.VerifyAccountUserVerificationToken(token, false)
		require.NoError(t, err)

		updated, err := m.GetAccountUserRec(accountUserRec.ID, nil)
		require.NoError(t, err)
		require.Equal(t, account_record.AccountUserStatusActive, updated.Status, "active user stays active after verification")
	})
}
