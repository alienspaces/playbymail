package jobworker

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"

	corejobworker "gitlab.com/alienspaces/playbymail/core/jobworker"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/core/type/storer"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/jobworker/adventure_game"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
	"gitlab.com/alienspaces/playbymail/internal/utils/config"
)

// GameTurnProcessingWorkerArgs defines the arguments for processing a game instance turn
type GameTurnProcessingWorkerArgs struct {
	GameInstanceID string `json:"game_instance_id"`
	TurnNumber     int    `json:"turn_number"`
}

func (GameTurnProcessingWorkerArgs) Kind() string { return "game_turn_processing" }

// GameTurnProcessingWorker processes a game instance turn
type GameTurnProcessingWorker struct {
	river.WorkerDefaults[GameTurnProcessingWorkerArgs]
	JobWorker
}

func NewGameTurnProcessingWorker(l logger.Logger, cfg config.Config, s storer.Storer) (*GameTurnProcessingWorker, error) {
	jw, err := NewJobWorker(l, cfg, s)
	if err != nil {
		return nil, err
	}

	return &GameTurnProcessingWorker{
		JobWorker: *jw,
	}, nil
}

func (w *GameTurnProcessingWorker) Work(ctx context.Context, j *river.Job[GameTurnProcessingWorkerArgs]) error {
	l := w.JobWorker.Log.WithFunctionContext("GameTurnProcessingWorker/Work")

	l.Info("running job ID >%s< Args >%#v<", strconv.FormatInt(j.ID, 10), j.Args)

	c, m, err := w.beginJob(ctx)
	if err != nil {
		return err
	}
	defer func() {
		m.Tx.Rollback(context.Background())
	}()

	_, err = w.DoWork(ctx, m, c, j)
	if err != nil {
		l.Error("GameTurnProcessingWorker job ID >%s< Args >%#v< failed >%v<", strconv.FormatInt(j.ID, 10), j.Args, err)
		return err
	}

	return corejobworker.CompleteJob(ctx, m.Tx, j)
}

type GameTurnProcessingDoWorkResult struct {
	GameInstanceID string
	TurnNumber     int
	ProcessedAt    time.Time
}

func (w *GameTurnProcessingWorker) DoWork(ctx context.Context, m *domain.Domain, c *river.Client[pgx.Tx], j *river.Job[GameTurnProcessingWorkerArgs]) (*GameTurnProcessingDoWorkResult, error) {
	l := w.JobWorker.Log.WithFunctionContext("GameTurnProcessingWorker/DoWork")

	l.Info("processing game turn for instance >%s< turn >%d<", j.Args.GameInstanceID, j.Args.TurnNumber)

	// Get the game instance
	gameInstanceRec, err := m.GetGameInstanceRec(j.Args.GameInstanceID, nil)
	if err != nil {
		l.Warn("failed to get game instance ID >%s<; cannot process game turn >%v<", j.Args.GameInstanceID, err)
		return nil, err
	}

	// Verify we're processing the correct turn
	if gameInstanceRec.CurrentTurn != j.Args.TurnNumber {
		l.Warn("turn number mismatch for game instance ID >%s<: expected >%d< but instance is at turn >%d<", j.Args.GameInstanceID, j.Args.TurnNumber, gameInstanceRec.CurrentTurn)
		return nil, fmt.Errorf("turn number mismatch for game instance ID >%s<", j.Args.GameInstanceID)
	}

	// Begin turn processing
	gameInstanceRec, err = m.BeginTurnProcessing(j.Args.GameInstanceID)
	if err != nil {
		l.Warn("failed to begin turn processing for game instance ID >%s<; cannot process game turn >%v<", j.Args.GameInstanceID, err)
		return nil, err
	}

	// Process player turns based on game type
	l.Info("processing turn logic for game instance ID >%s< turn >%d<", j.Args.GameInstanceID, j.Args.TurnNumber)

	// Get the game to determine the game type
	gameRec, err := m.GetGameRec(gameInstanceRec.GameID, nil)
	if err != nil {
		l.Warn("failed to get game ID >%s< for game instance ID >%s<; cannot process game turn >%v<", gameInstanceRec.GameID, j.Args.GameInstanceID, err)
		return nil, err
	}

	// Route to appropriate game type processor and process the turn for all players
	switch gameRec.GameType {
	case game_record.GameTypeAdventure:
		adventureProcessor := adventure_game.NewAdventureGame(l, m)
		err = adventureProcessor.ProcessTurnSheets(ctx, gameInstanceRec)
		if err != nil {
			l.Warn("failed to process adventure game turn for game instance ID >%s< turn >%d<; cannot process game turn >%v<", j.Args.GameInstanceID, j.Args.TurnNumber, err)
			return nil, err
		}
	default:
		l.Warn("unsupported game type >%s< for game instance ID >%s<; cannot process game turn >%v<", gameRec.GameType, j.Args.GameInstanceID, err)
		return nil, fmt.Errorf("unsupported game type: %s for game instance ID >%s<", gameRec.GameType, j.Args.GameInstanceID)
	}

	// TODO: Generate post run processing jobs such as notifications, etc.

	// Complete the turn
	gameInstanceRec, err = m.CompleteTurn(j.Args.GameInstanceID)
	if err != nil {
		l.Warn("failed to complete turn for game instance ID >%s< turn >%d<; cannot process game turn >%v<", j.Args.GameInstanceID, j.Args.TurnNumber, err)
		return nil, err
	}

	l.Info("completed turn processing for game instance >%s< turn >%d<", gameInstanceRec.ID, j.Args.TurnNumber)

	// Generate new turn sheets for the next turn
	l.Info("generating new turn sheets for game instance >%s< turn >%d<", gameInstanceRec.ID, gameInstanceRec.CurrentTurn)

	// Route to appropriate game type processor and process the turn for all players
	switch gameRec.GameType {
	case game_record.GameTypeAdventure:
		adventureProcessor := adventure_game.NewAdventureGame(l, m)
		err = adventureProcessor.GenerateTurnSheets(ctx, gameInstanceRec)
		if err != nil {
			l.Warn("failed to generate new turn sheets for game instance ID >%s< turn >%d<; cannot process game turn >%v<", j.Args.GameInstanceID, j.Args.TurnNumber, err)
			return nil, err
		}
	default:
		l.Warn("unsupported game type >%s< for game instance ID >%s<; cannot generate new turn sheets >%v<", gameRec.GameType, j.Args.GameInstanceID, err)
		return nil, fmt.Errorf("unsupported game type: %s for game instance ID >%s<", gameRec.GameType, j.Args.GameInstanceID)
	}

	return &GameTurnProcessingDoWorkResult{
		GameInstanceID: j.Args.GameInstanceID,
		TurnNumber:     j.Args.TurnNumber,
		ProcessedAt:    time.Now(),
	}, nil
}
