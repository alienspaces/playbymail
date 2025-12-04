package harness

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/riverqueue/river"

	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/jobqueue"
	"gitlab.com/alienspaces/playbymail/internal/jobworker"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
	"gitlab.com/alienspaces/playbymail/internal/turn_sheet"
)

// ProcessTurn executes turn processing workers for a game instance
// This simulates automated play by processing all scanned turn sheets
func (t *Testing) ProcessTurn(ctx context.Context, gameInstanceRef string) error {
	l := t.Logger("ProcessTurn")

	// Get game instance by reference
	gameInstanceID, ok := t.Data.Refs.GameInstanceRefs[gameInstanceRef]
	if !ok {
		l.Warn("game instance ref >%s< not found", gameInstanceRef)
		return fmt.Errorf("game instance ref not found: %s", gameInstanceRef)
	}

	// Create turn processing worker
	worker, err := jobworker.NewGameTurnProcessingWorker(t.Log, t.Config, t.Store)
	if err != nil {
		l.Warn("failed to create turn processing worker >%v<", err)
		return fmt.Errorf("failed to create turn processing worker: %w", err)
	}

	// Get game instance record to get current turn
	tx, err := t.Store.BeginTx()
	if err != nil {
		l.Warn("failed to begin transaction >%v<", err)
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	m, err := domain.NewDomain(t.Log, t.Config)
	if err != nil {
		l.Warn("failed to create domain >%v<", err)
		return fmt.Errorf("failed to create domain: %w", err)
	}

	err = m.Init(tx)
	if err != nil {
		l.Warn("failed to init domain >%v<", err)
		return fmt.Errorf("failed to init domain: %w", err)
	}

	gameInstanceRec, err := m.GetGameInstanceRec(gameInstanceID, nil)
	if err != nil {
		l.Warn("failed to get game instance >%s< >%v<", gameInstanceID, err)
		return fmt.Errorf("failed to get game instance: %w", err)
	}

	// Create turn processing worker
	// We'll call DoWork directly with a manually created domain
	// This bypasses the Work method's transaction management
	jobArgs := jobworker.GameTurnProcessingWorkerArgs{
		GameInstanceID: gameInstanceID,
		TurnNumber:     gameInstanceRec.CurrentTurn,
	}

	// Create a minimal job struct for DoWork
	// DoWork only uses j.Args, so we can create a minimal struct
	// We'll use the job ID from a newly created job
	// Actually, we don't need to insert a job - we can call DoWork directly
	// But DoWork expects a *river.Job which we can't easily create
	// So we'll insert a job to get a proper job struct, then use it
	insertResult, err := t.JobClient.Insert(ctx, jobArgs, &river.InsertOpts{
		Queue: jobqueue.QueueGame,
	})
	if err != nil {
		l.Warn("failed to insert turn processing job >%v<", err)
		return fmt.Errorf("failed to insert turn processing job: %w", err)
	}

	// Create a new domain for DoWork
	// DoWork needs a domain with a transaction
	txForDoWork, err := t.Store.BeginTx()
	if err != nil {
		l.Warn("failed to begin transaction for DoWork >%v<", err)
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer txForDoWork.Rollback(ctx)

	mForDoWork, err := domain.NewDomain(t.Log, t.Config)
	if err != nil {
		l.Warn("failed to create domain for DoWork >%v<", err)
		return fmt.Errorf("failed to create domain: %w", err)
	}

	err = mForDoWork.Init(txForDoWork)
	if err != nil {
		l.Warn("failed to init domain for DoWork >%v<", err)
		return fmt.Errorf("failed to init domain: %w", err)
	}

	// Create a minimal job struct for DoWork
	// We need to convert JobRow to Job
	// JobRow has an Args field that we can use
	// Actually, we can use reflection or type assertion
	// Or we can create a minimal job with just Args
	// Since DoWork only uses j.Args, we can create a minimal struct
	job := &river.Job[jobworker.GameTurnProcessingWorkerArgs]{
		Args: jobArgs,
	}
	// Set the ID from the insert result if available
	if insertResult != nil && insertResult.Job != nil {
		// We can't directly access the ID from JobRow
		// So we'll just use Args for DoWork
	}

	// Call DoWork directly with the manually created domain
	_, err = worker.DoWork(ctx, mForDoWork, t.JobClient, job)
	if err != nil {
		l.Warn("failed to process turn >%v<", err)
		return fmt.Errorf("failed to process turn: %w", err)
	}

	// Commit the transaction after DoWork succeeds
	// Note: DoWork doesn't commit - we need to commit ourselves
	err = txForDoWork.Commit(ctx)
	if err != nil {
		l.Warn("failed to commit transaction >%v<", err)
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Commit the original transaction
	err = tx.Commit(ctx)
	if err != nil {
		l.Warn("failed to commit original transaction >%v<", err)
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	l.Info("processed turn for game instance >%s< turn >%d<", gameInstanceID, gameInstanceRec.CurrentTurn)

	return nil
}

// AutomatedPlayConfig defines which turn sheets should be processed and their choices
type AutomatedPlayConfig struct {
	// TurnSheetRefs maps turn sheet references to their scan data choices
	// For location choice turn sheets, the value should be a location instance ID to choose
	// For join game turn sheets, this can be empty or used for other purposes
	TurnSheetChoices map[string]AutomatedPlayChoice
}

// AutomatedPlayChoice defines what choice to make for a turn sheet
type AutomatedPlayChoice struct {
	// LocationInstanceID is used for location choice turn sheets
	// Specifies which location instance ID to choose
	LocationInstanceID string
}

// SimulateTurnSheetUpload simulates uploading a turn sheet with contrived scan data
// It marks the turn sheet as scanned with the provided scan data
func (t *Testing) SimulateTurnSheetUpload(ctx context.Context, turnSheetRef string, choice AutomatedPlayChoice) error {
	l := t.Logger("SimulateTurnSheetUpload")

	// Get turn sheet by reference
	turnSheetRec, err := t.Data.GetGameTurnSheetRecByRef(turnSheetRef)
	if err != nil {
		l.Warn("failed to get turn sheet by ref >%s< >%v<", turnSheetRef, err)
		return fmt.Errorf("failed to get turn sheet by ref: %w", err)
	}

	// Get the sheet data to understand what type of turn sheet this is
	var sheetData map[string]any
	if err := json.Unmarshal(turnSheetRec.SheetData, &sheetData); err != nil {
		l.Warn("failed to unmarshal sheet data >%v<", err)
		return fmt.Errorf("failed to parse sheet data: %w", err)
	}

	// Construct scan data based on sheet type
	var scanDataBytes []byte
	switch turnSheetRec.SheetType {
	case adventure_game_record.AdventureGameTurnSheetTypeLocationChoice:
		scanDataBytes, err = t.constructLocationChoiceScanData(turnSheetRec, choice)
		if err != nil {
			l.Warn("failed to construct location choice scan data >%v<", err)
			return fmt.Errorf("failed to construct location choice scan data: %w", err)
		}
	case adventure_game_record.AdventureGameTurnSheetTypeJoinGame:
		// Join game turn sheets are typically handled during upload, not automated play
		// For now, we'll skip these or handle them separately if needed
		l.Info("skipping join game turn sheet >%s< for automated play", turnSheetRec.ID)
		return nil
	default:
		l.Warn("unsupported sheet type >%s< for automated play", turnSheetRec.SheetType)
		return fmt.Errorf("unsupported sheet type: %s", turnSheetRec.SheetType)
	}

	// Mark turn sheet as completed with scan data
	tx, err := t.Store.BeginTx()
	if err != nil {
		l.Warn("failed to begin transaction >%v<", err)
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	m, err := domain.NewDomain(t.Log, t.Config)
	if err != nil {
		l.Warn("failed to create domain >%v<", err)
		return fmt.Errorf("failed to create domain: %w", err)
	}

	err = m.Init(tx)
	if err != nil {
		l.Warn("failed to init domain >%v<", err)
		return fmt.Errorf("failed to init domain: %w", err)
	}

	err = m.MarkGameTurnSheetAsCompleted(turnSheetRec.ID, scanDataBytes)
	if err != nil {
		l.Warn("failed to mark turn sheet as completed >%v<", err)
		return fmt.Errorf("failed to mark turn sheet as completed: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		l.Warn("failed to commit transaction >%v<", err)
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	l.Info("simulated turn sheet upload for turn sheet >%s<", turnSheetRec.ID)

	return nil
}

// constructLocationChoiceScanData constructs valid LocationChoiceScanData from
// LocationChoiceData and a choice
func (t *Testing) constructLocationChoiceScanData(turnSheetRec *game_record.GameTurnSheet, choice AutomatedPlayChoice) ([]byte, error) {
	l := t.Logger("constructLocationChoiceScanData")

	// Parse the sheet data to get available location options
	var locationChoiceData turn_sheet.LocationChoiceData
	if err := json.Unmarshal(turnSheetRec.SheetData, &locationChoiceData); err != nil {
		l.Warn("failed to unmarshal location choice data >%v<", err)
		return nil, fmt.Errorf("failed to parse location choice data: %w", err)
	}

	// If a specific location instance ID was provided, use it
	// Otherwise, use the first available option
	chosenLocationID := choice.LocationInstanceID
	if chosenLocationID == "" && len(locationChoiceData.LocationOptions) > 0 {
		// Find the first valid location instance ID from the options
		// We need to map from location IDs in options to location instance IDs
		// For now, we'll use the LocationID from the first option as a fallback
		// TODO: This needs to be more sophisticated - LocationOptions contain LocationID
		// but we need LocationInstanceID. This is a limitation we need to address.
		l.Warn("no location instance ID specified, using first option")
		if len(locationChoiceData.LocationOptions) > 0 {
			// The LocationID in LocationOption is actually the location instance ID
			// based on how it's used in the processor
			chosenLocationID = locationChoiceData.LocationOptions[0].LocationID
		}
	}

	if chosenLocationID == "" {
		return nil, fmt.Errorf("no valid location choice available")
	}

	// Validate the choice is one of the available options
	isValidChoice := false
	for _, option := range locationChoiceData.LocationOptions {
		if option.LocationID == chosenLocationID {
			isValidChoice = true
			break
		}
	}

	if !isValidChoice {
		return nil, fmt.Errorf("invalid location choice: %s is not in available options", chosenLocationID)
	}

	// Construct the scan data
	scanData := turn_sheet.LocationChoiceScanData{
		Choices: []string{chosenLocationID},
	}

	scanDataBytes, err := json.Marshal(scanData)
	if err != nil {
		l.Warn("failed to marshal scan data >%v<", err)
		return nil, fmt.Errorf("failed to marshal scan data: %w", err)
	}

	return scanDataBytes, nil
}
