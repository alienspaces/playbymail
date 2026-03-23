package runner

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v2"

	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/jobworker"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
)

// resendTurnSheetEmail regenerates turn sheet tokens for all players in a game
// instance and queues a notification email for each one. Use this when a player
// reports that their email link has expired or never arrived.
func (rnr *Runner) resendTurnSheetEmail(c *cli.Context) error {
	l := loggerWithFunctionContext(rnr.Log, "resendTurnSheetEmail")

	gameInstanceID := c.String("game-instance-id")
	if gameInstanceID == "" {
		return fmt.Errorf("--game-instance-id is required")
	}

	l.Info("** Resend Turn Sheet Emails for game instance >%s< **", gameInstanceID)

	if err := rnr.InitDomain(); err != nil {
		l.Warn("failed domain init >%v<", err)
		return err
	}

	dm, ok := rnr.Domain.(*domain.Domain)
	if !ok {
		return fmt.Errorf("domain type assertion failed")
	}

	gameInstanceRec, err := dm.GetGameInstanceRec(gameInstanceID, nil)
	if err != nil {
		l.Warn("failed to get game instance >%s< >%v<", gameInstanceID, err)
		return fmt.Errorf("game instance not found: %w", err)
	}

	subInstRecs, err := dm.GetManyGameSubscriptionInstanceRecs(&coresql.Options{
		Params: []coresql.Param{
			{
				Col: game_record.FieldGameSubscriptionInstanceGameInstanceID,
				Val: gameInstanceID,
			},
		},
	})
	if err != nil {
		l.Warn("failed to get game subscription instances >%v<", err)
		return err
	}

	if len(subInstRecs) == 0 {
		fmt.Printf("No players found for game instance %s\n", gameInstanceID)
		return nil
	}

	fmt.Printf("\nResending turn sheet emails for game instance %s (turn %d)...\n\n",
		gameInstanceID, gameInstanceRec.CurrentTurn)

	sent := 0
	for _, si := range subInstRecs {
		args := jobworker.SendTurnSheetNotificationEmailWorkerArgs{
			GameSubscriptionInstanceID: si.ID,
			TurnNumber:                 gameInstanceRec.CurrentTurn,
		}
		_, err := rnr.JobClient.Insert(context.Background(), args, nil)
		if err != nil {
			l.Warn("failed to queue email for subscription instance >%s< >%v<", si.ID, err)
			fmt.Printf("  FAILED  subscription instance %s: %v\n", si.ID, err)
			continue
		}
		sent++
		fmt.Printf("  QUEUED  subscription instance %s\n", si.ID)
	}

	fmt.Printf("\n%d/%d email(s) queued for turn %d\n\n", sent, len(subInstRecs), gameInstanceRec.CurrentTurn)

	return nil
}
