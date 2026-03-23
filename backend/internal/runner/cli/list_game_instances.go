package runner

import (
	"fmt"
	"sort"

	"github.com/urfave/cli/v2"

	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
)

// listGameInstances prints all game instances with their status, current turn,
// player count, and next turn due date.
func (rnr *Runner) listGameInstances(c *cli.Context) error {
	l := loggerWithFunctionContext(rnr.Log, "listGameInstances")

	l.Info("** List Game Instances **")

	if err := rnr.InitDomain(); err != nil {
		l.Warn("failed domain init >%v<", err)
		return err
	}

	dm, ok := rnr.Domain.(*domain.Domain)
	if !ok {
		return fmt.Errorf("domain type assertion failed")
	}

	instanceRecs, err := dm.GetManyGameInstanceRecs(nil)
	if err != nil {
		l.Warn("failed getting game instance records >%v<", err)
		return err
	}

	// Index games by ID for name lookup
	gameRecs, err := dm.GetManyGameRecs(nil)
	if err != nil {
		l.Warn("failed getting game records >%v<", err)
		return err
	}
	gamesByID := make(map[string]*game_record.Game, len(gameRecs))
	for _, g := range gameRecs {
		gamesByID[g.ID] = g
	}

	// Count subscribed players per game instance
	subInstRecs, err := dm.GetManyGameSubscriptionInstanceRecs(&coresql.Options{})
	if err != nil {
		l.Warn("failed getting game subscription instance records >%v<", err)
		return err
	}
	playerCountByInstanceID := make(map[string]int)
	for _, si := range subInstRecs {
		playerCountByInstanceID[si.GameInstanceID]++
	}

	sort.Slice(instanceRecs, func(i, j int) bool {
		return instanceRecs[i].CreatedAt.After(instanceRecs[j].CreatedAt)
	})

	fmt.Printf("\n%-38s  %-30s  %-10s  %-5s  %-7s  %-22s\n",
		"Game Instance ID", "Game Name", "Status", "Turn", "Players", "Next Turn Due")
	fmt.Printf("%-38s  %-30s  %-10s  %-5s  %-7s  %-22s\n",
		"--------------------------------------",
		"------------------------------",
		"----------",
		"-----",
		"-------",
		"----------------------",
	)

	for _, inst := range instanceRecs {
		gameName := ""
		if g, ok := gamesByID[inst.GameID]; ok {
			gameName = g.Name
			if len(gameName) > 30 {
				gameName = gameName[:27] + "..."
			}
		}

		nextDue := "-"
		if inst.NextTurnDueAt.Valid {
			nextDue = inst.NextTurnDueAt.Time.Local().Format("2006-01-02 15:04 MST")
		}

		fmt.Printf("%-38s  %-30s  %-10s  %-5d  %-7d  %-22s\n",
			inst.ID,
			gameName,
			inst.Status,
			inst.CurrentTurn,
			playerCountByInstanceID[inst.ID],
			nextDue,
		)
	}

	fmt.Printf("\n%d game instance(s) found\n\n", len(instanceRecs))

	return nil
}
