package runner

import (
	"fmt"
	"sort"
	"time"

	"github.com/urfave/cli/v2"

	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/record/account_record"
)

// listUsers prints all accounts and their users with status and session info.
func (rnr *Runner) listUsers(c *cli.Context) error {
	l := loggerWithFunctionContext(rnr.Log, "listUsers")

	l.Info("** List Users **")

	if err := rnr.InitDomain(); err != nil {
		l.Warn("failed domain init >%v<", err)
		return err
	}

	dm, ok := rnr.Domain.(*domain.Domain)
	if !ok {
		return fmt.Errorf("domain type assertion failed")
	}

	accountUserRecs, err := dm.GetManyAccountUserRecs(nil)
	if err != nil {
		l.Warn("failed getting account user records >%v<", err)
		return err
	}

	accountRecs, err := dm.GetManyAccountRecs(nil)
	if err != nil {
		l.Warn("failed getting account records >%v<", err)
		return err
	}

	accountsByID := make(map[string]*account_record.Account, len(accountRecs))
	for _, a := range accountRecs {
		accountsByID[a.ID] = a
	}

	sort.Slice(accountUserRecs, func(i, j int) bool {
		return accountUserRecs[i].CreatedAt.After(accountUserRecs[j].CreatedAt)
	})

	now := time.Now()

	fmt.Printf("\n%-40s  %-10s  %-18s  %-14s  %-20s  %-20s\n",
		"Email", "Acct Status", "User Status", "Session", "Created", "Last Updated")
	fmt.Printf("%-40s  %-10s  %-18s  %-14s  %-20s  %-20s\n",
		"----------------------------------------",
		"----------",
		"------------------",
		"--------------",
		"--------------------",
		"--------------------",
	)

	for _, u := range accountUserRecs {
		acctStatus := "unknown"
		if a, ok := accountsByID[u.AccountID]; ok {
			acctStatus = a.Status
		}

		sessionInfo := "no session"
		if u.SessionToken.Valid && u.SessionToken.String != "" {
			if u.SessionTokenExpiresAt.Valid && u.SessionTokenExpiresAt.Time.After(now) {
				sessionInfo = "active"
			} else {
				sessionInfo = "expired"
			}
		}

		created := u.CreatedAt.Local().Format("2006-01-02 15:04")

		lastUpdated := "-"
		if u.UpdatedAt.Valid {
			lastUpdated = u.UpdatedAt.Time.Local().Format("2006-01-02 15:04")
		}

		fmt.Printf("%-40s  %-10s  %-18s  %-14s  %-20s  %-20s\n",
			u.Email,
			acctStatus,
			u.Status,
			sessionInfo,
			created,
			lastUpdated,
		)
	}

	fmt.Printf("\n%d user(s) found\n\n", len(accountUserRecs))

	return nil
}
