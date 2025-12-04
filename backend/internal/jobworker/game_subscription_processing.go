package jobworker

import (
	"context"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"

	corejobworker "gitlab.com/alienspaces/playbymail/core/jobworker"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/core/type/storer"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/jobqueue"
	"gitlab.com/alienspaces/playbymail/internal/jobworker/adventure_game"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
	"gitlab.com/alienspaces/playbymail/internal/utils/config"
)

// ProcessSubscriptionWorkerArgs defines the job payload for processing join game turn sheets
// when a game subscription is approved.
type ProcessSubscriptionWorkerArgs struct {
	GameSubscriptionID string
}

func (ProcessSubscriptionWorkerArgs) Kind() string {
	return "join-game-turn-sheet"
}

func (ProcessSubscriptionWorkerArgs) InsertOpts() river.InsertOpts {
	return river.InsertOpts{Queue: jobqueue.QueueDefault}
}

// ProcessSubscriptionProcessor defines the interface for processing join game turn sheets
// for different game types.
type ProcessSubscriptionProcessor interface {
	// ProcessProcessSubscription processes a join game turn sheet and creates the
	// necessary game entities (game instance, character, character instance, etc.)
	ProcessProcessSubscription(ctx context.Context, subscriptionRec *game_record.GameSubscription, turnSheetRec *game_record.GameTurnSheet) error
}

// ProcessSubscriptionWorker processes join game turn sheets when a game subscription
// is approved, creating the necessary game entities for the player to participate.
type ProcessSubscriptionWorker struct {
	river.WorkerDefaults[ProcessSubscriptionWorkerArgs]
	JobWorker
}

func NewProcessSubscriptionWorker(l logger.Logger, cfg config.Config, s storer.Storer) (*ProcessSubscriptionWorker, error) {
	jw, err := NewJobWorker(l, cfg, s)
	if err != nil {
		return nil, err
	}

	return &ProcessSubscriptionWorker{
		JobWorker: *jw,
	}, nil
}

func (w *ProcessSubscriptionWorker) Work(ctx context.Context, j *river.Job[ProcessSubscriptionWorkerArgs]) error {
	l := w.JobWorker.Log.WithFunctionContext("ProcessSubscriptionWorker/Work")

	l.Info("running job ID >%s< Args >%#v<", strconv.FormatInt(j.ID, 10), j.Args)

	c, d, err := w.beginJob(ctx)
	if err != nil {
		return err
	}
	defer func() {
		d.Tx.Rollback(context.Background())
	}()

	_, err = w.DoWork(ctx, d, c, j)
	if err != nil {
		l.Error("ProcessSubscriptionWorker job ID >%s< Args >%#v< failed >%v<", strconv.FormatInt(j.ID, 10), j.Args, err)
		return err
	}

	return corejobworker.CompleteJob(ctx, d.Tx, j)
}

type ProcessSubscriptionDoWorkResult struct {
	GameSubscriptionID string
	ProcessedAt        string
}

func (w *ProcessSubscriptionWorker) DoWork(ctx context.Context, m *domain.Domain, c *river.Client[pgx.Tx], j *river.Job[ProcessSubscriptionWorkerArgs]) (*ProcessSubscriptionDoWorkResult, error) {
	l := w.JobWorker.Log.WithFunctionContext("ProcessSubscriptionWorker/DoWork")

	l.Info("processing join game turn sheet for subscription ID >%s<", j.Args.GameSubscriptionID)

	// Get the subscription record
	subscriptionRec, err := m.GetGameSubscriptionRec(j.Args.GameSubscriptionID, nil)
	if err != nil {
		l.Warn("failed to get game subscription record >%v<", err)
		return nil, err
	}

	// Verify subscription is active
	if subscriptionRec.Status != game_record.GameSubscriptionStatusActive {
		l.Warn("subscription is not active, current status >%s<", subscriptionRec.Status)
		return nil, nil
	}

	// Get the game record to determine game type
	gameRec, err := m.GetGameRec(subscriptionRec.GameID, nil)
	if err != nil {
		l.Warn("failed to get game record >%v<", err)
		return nil, err
	}

	// Find the join game turn sheet for this subscription
	turnSheetRec, err := w.findProcessSubscription(m, subscriptionRec)
	if err != nil {
		l.Warn("failed to find join game turn sheet >%v<", err)
		return nil, err
	}

	if turnSheetRec == nil {
		l.Warn("no join game turn sheet found for subscription ID >%s<", j.Args.GameSubscriptionID)
		return nil, nil
	}

	// Initialize all game type processors
	processors, err := w.initializeProcessors(l, m)
	if err != nil {
		l.Warn("failed to initialize processors >%v<", err)
		return nil, err
	}

	// Get the appropriate processor for this game type
	processor, exists := processors[gameRec.GameType]
	if !exists {
		l.Warn("unsupported game type >%s< for subscription ID >%s<", gameRec.GameType, j.Args.GameSubscriptionID)
		return nil, nil
	}

	// Process join game turn sheet using the game-specific processor
	err = processor.ProcessProcessSubscription(ctx, subscriptionRec, turnSheetRec)
	if err != nil {
		l.Warn("failed to process join game turn sheet for subscription ID >%s< >%v<", j.Args.GameSubscriptionID, err)
		return nil, err
	}

	l.Info("completed processing join game turn sheet for subscription ID >%s<", j.Args.GameSubscriptionID)

	return &ProcessSubscriptionDoWorkResult{
		GameSubscriptionID: j.Args.GameSubscriptionID,
		ProcessedAt:        "now",
	}, nil
}

// findProcessSubscription finds the join game turn sheet for a subscription
func (w *ProcessSubscriptionWorker) findProcessSubscription(m *domain.Domain, subscriptionRec *game_record.GameSubscription) (*game_record.GameTurnSheet, error) {
	// Get all turn sheets for this account
	turnSheetRecs, err := m.GetGameTurnSheetRecsByAccount(subscriptionRec.AccountID)
	if err != nil {
		return nil, err
	}

	// Find the join game turn sheet for this game
	for _, turnSheetRec := range turnSheetRecs {
		if turnSheetRec.GameID == subscriptionRec.GameID &&
			turnSheetRec.SheetType == adventure_game_record.AdventureGameTurnSheetTypeJoinGame {
			return turnSheetRec, nil
		}
	}

	return nil, nil
}

// initializeProcessors creates and registers all available game type processors
func (w *ProcessSubscriptionWorker) initializeProcessors(l logger.Logger, d *domain.Domain) (map[string]ProcessSubscriptionProcessor, error) {
	processors := make(map[string]ProcessSubscriptionProcessor)

	// Register adventure game processor
	adventureProcessor, err := adventure_game.NewAdventureGameProcessSubscriptionProcessor(l, d)
	if err != nil {
		return nil, err
	}
	processors[game_record.GameTypeAdventure] = adventureProcessor

	// TODO: Add new game type processors here
	// Example: processors[game_record.GameTypeStrategy] = strategyProcessor
	// Example: processors[game_record.GameTypePuzzle] = puzzleProcessor

	return processors, nil
}
