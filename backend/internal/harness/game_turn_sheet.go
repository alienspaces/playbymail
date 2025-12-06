package harness

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/riverqueue/river"

	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/jobworker"
	"gitlab.com/alienspaces/playbymail/internal/jobworker/adventure_game"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
	"gitlab.com/alienspaces/playbymail/internal/turn_sheet"
)

// GenerateTurnSheetsForGameInstance runs the adventure game job worker to create turn sheets
// for the specified game instance reference and returns the created records.
func (t *Testing) GenerateTurnSheetsForGameInstance(ctx context.Context, gameInstanceRef string) ([]*game_record.GameTurnSheet, error) {
	l := t.Logger("GenerateTurnSheetsForGameInstance")

	gameInstanceID, ok := t.Data.Refs.GameInstanceRefs[gameInstanceRef]
	if !ok {
		l.Warn("game instance ref >%s< not found", gameInstanceRef)
		return nil, fmt.Errorf("game instance ref not found: %s", gameInstanceRef)
	}

	gameInstanceRec, err := t.Data.GetGameInstanceRecByID(gameInstanceID)
	if err != nil {
		l.Warn("failed to get game instance by id >%s< >%v<", gameInstanceID, err)
		return nil, fmt.Errorf("failed to get game instance: %w", err)
	}

	tx, err := t.Store.BeginTx()
	if err != nil {
		l.Warn("failed to begin transaction for turn sheet generation >%v<", err)
		return nil, fmt.Errorf("failed to begin transaction for turn sheet generation: %w", err)
	}
	defer func() {
		if tx != nil {
			tx.Rollback(ctx)
		}
	}()

	jobDomain, err := domain.NewDomain(t.Log, t.Config)
	if err != nil {
		l.Warn("failed to create domain for turn sheet generation >%v<", err)
		return nil, fmt.Errorf("failed to create domain for turn sheet generation: %w", err)
	}

	err = jobDomain.Init(tx)
	if err != nil {
		l.Warn("failed to init domain for turn sheet generation >%v<", err)
		return nil, fmt.Errorf("failed to init domain for turn sheet generation: %w", err)
	}

	processor, err := adventure_game.NewAdventureGame(t.Log, jobDomain)
	if err != nil {
		l.Warn("failed to create adventure game processor >%v<", err)
		return nil, fmt.Errorf("failed to create adventure game processor: %w", err)
	}

	turnSheetRecs, err := processor.CreateTurnSheets(ctx, gameInstanceRec)
	if err != nil {
		l.Warn("failed to create turn sheets via job worker >%v<", err)
		return nil, fmt.Errorf("failed to create turn sheets: %w", err)
	}

	for _, rec := range turnSheetRecs {
		t.Data.AddGameTurnSheetRec(rec)
		t.teardownData.AddGameTurnSheetRec(rec)
	}

	if len(turnSheetRecs) > 0 {
		if err := t.assignTurnSheetReferencesForTurnNumber(gameInstanceRef, turnSheetRecs[0].TurnNumber, turnSheetRecs); err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		l.Warn("failed to commit turn sheet generation transaction >%v<", err)
		return nil, fmt.Errorf("failed to commit turn sheet generation: %w", err)
	}
	tx = nil

	l.Info("created %d turn sheets for game instance >%s<", len(turnSheetRecs), gameInstanceRef)

	return turnSheetRecs, nil
}

func (t *Testing) generateTurnSheetsForGameInstanceInTx(ctx context.Context, gameInstanceRec *game_record.GameInstance, gameInstanceRef string) ([]*game_record.GameTurnSheet, error) {
	l := t.Logger("generateTurnSheetsForGameInstanceInTx")

	processor, err := adventure_game.NewAdventureGame(t.Log, t.Domain.(*domain.Domain))
	if err != nil {
		l.Warn("failed to create adventure game processor >%v<", err)
		return nil, fmt.Errorf("failed to create adventure game processor: %w", err)
	}

	turnSheetRecs, err := processor.CreateTurnSheets(ctx, gameInstanceRec)
	if err != nil {
		l.Warn("failed to create turn sheets via processor >%v<", err)
		return nil, fmt.Errorf("failed to create turn sheets: %w", err)
	}

	if len(turnSheetRecs) == 0 {
		return nil, fmt.Errorf("no turn sheets were created for game instance %s on turn %d", gameInstanceRec.ID, gameInstanceRec.CurrentTurn+1)
	}

	for _, rec := range turnSheetRecs {
		t.Data.AddGameTurnSheetRec(rec)
		t.teardownData.AddGameTurnSheetRec(rec)
	}

	if err := t.assignTurnSheetReferencesForTurnNumber(gameInstanceRef, turnSheetRecs[0].TurnNumber, turnSheetRecs); err != nil {
		return nil, err
	}

	l.Info("created %d turn sheets for game instance >%s< in current transaction", len(turnSheetRecs), gameInstanceRef)

	return turnSheetRecs, nil
}

func (t *Testing) processGameTurnForInstanceInTx(ctx context.Context, gameInstanceID string) error {
	l := t.Logger("processGameTurnForInstanceInTx")

	dom := t.Domain.(*domain.Domain)

	processor, err := adventure_game.NewAdventureGame(t.Log, dom)
	if err != nil {
		l.Warn("failed to create adventure game processor >%v<", err)
		return fmt.Errorf("failed to create adventure game processor: %w", err)
	}

	instanceRec, err := dom.BeginTurnProcessing(gameInstanceID)
	if err != nil {
		l.Warn("failed to begin turn processing >%v<", err)
		return fmt.Errorf("failed to begin turn processing: %w", err)
	}

	if err := processor.ProcessTurnSheets(ctx, instanceRec); err != nil {
		l.Warn("failed to process turn sheets >%v<", err)
		return fmt.Errorf("failed to process turn sheets: %w", err)
	}

	instanceRec, err = dom.CompleteTurn(gameInstanceID)
	if err != nil {
		l.Warn("failed to complete turn >%v<", err)
		return fmt.Errorf("failed to complete turn: %w", err)
	}

	if _, err := processor.CreateTurnSheets(ctx, instanceRec); err != nil {
		l.Warn("failed to create next turn sheets >%v<", err)
		return fmt.Errorf("failed to create next turn sheets: %w", err)
	}

	l.Info("processed turn for game instance >%s< turn >%d<", gameInstanceID, instanceRec.CurrentTurn)

	return nil
}

func (t *Testing) assignTurnSheetReferencesForTurnNumber(gameInstanceRef string, turnNumber int, turnSheetRecs []*game_record.GameTurnSheet) error {
	if len(turnSheetRecs) == 0 || turnNumber == 0 {
		return nil
	}
	turnCfg, err := t.DataConfig.findGameTurnConfig(gameInstanceRef, turnNumber)
	if err != nil {
		return err
	}
	if turnCfg == nil {
		return nil
	}
	return t.assignTurnSheetReferencesForTurn(*turnCfg, turnSheetRecs)
}

func (t *Testing) assignTurnSheetReferencesForTurn(turnCfg GameTurnConfig, turnSheetRecs []*game_record.GameTurnSheet) error {
	if len(turnCfg.AdventureGameTurnSheetConfigs) == 0 {
		return nil
	}

	dom := t.Domain.(*domain.Domain)
	characterTurnSheetMap := make(map[string]string)

	for _, rec := range turnSheetRecs {
		linkRecs, err := dom.GetManyAdventureGameTurnSheetRecs(&coresql.Options{
			Params: []coresql.Param{
				{
					Col: adventure_game_record.FieldAdventureGameTurnSheetGameTurnSheetID,
					Val: rec.ID,
				},
			},
			Limit: 1,
		})
		if err != nil {
			return fmt.Errorf("failed to lookup adventure game turn sheet link for turn sheet ID %s: %w", rec.ID, err)
		}
		if len(linkRecs) == 0 {
			continue
		}
		linkRec := linkRecs[0]
		characterTurnSheetMap[linkRec.AdventureGameCharacterInstanceID] = linkRec.GameTurnSheetID
		t.Data.AddAdventureGameTurnSheetRec(linkRec)
		t.teardownData.AddAdventureGameTurnSheetRec(linkRec)
	}

	for _, cfg := range turnCfg.AdventureGameTurnSheetConfigs {
		if cfg.GameTurnSheetConfig.Reference == "" || cfg.GameCharacterInstanceRef == "" {
			continue
		}
		characterInstanceID, ok := t.Data.Refs.AdventureGameCharacterInstanceRefs[cfg.GameCharacterInstanceRef]
		if !ok {
			return fmt.Errorf("adventure game character instance ref >%s< not found", cfg.GameCharacterInstanceRef)
		}
		gameTurnSheetID, ok := characterTurnSheetMap[characterInstanceID]
		if !ok {
			return fmt.Errorf("no turn sheet found for character ref >%s< on turn %d", cfg.GameCharacterInstanceRef, turnCfg.TurnNumber)
		}
		t.Data.Refs.GameTurnSheetRefs[cfg.GameTurnSheetConfig.Reference] = gameTurnSheetID
	}

	return nil
}

// processJoinGameSubscriptionInSetup processes a join game subscription during harness Setup
// This uses the existing harness transaction and domain
func (t *Testing) processJoinGameSubscriptionInSetup(ctx context.Context, subscriptionRef string, joinGameScanData any) (*game_record.GameTurnSheet, error) {
	l := t.Logger("processJoinGameSubscriptionInSetup")

	// Use the existing harness domain (already initialized with transaction during Setup)
	m := t.Domain.(*domain.Domain)

	// Get subscription by reference
	subscriptionID, ok := t.Data.Refs.GameSubscriptionRefs[subscriptionRef]
	if !ok {
		l.Warn("subscription ref >%s< not found", subscriptionRef)
		return nil, fmt.Errorf("subscription ref not found: %s", subscriptionRef)
	}

	subscriptionRec, err := t.Data.GetGameSubscriptionRecByID(subscriptionID)
	if err != nil {
		l.Warn("failed to get subscription by id >%s< >%v<", subscriptionID, err)
		return nil, fmt.Errorf("failed to get subscription by id: %w", err)
	}

	// Get account and game records
	accountRec, err := t.Data.GetAccountRecByID(subscriptionRec.AccountID)
	if err != nil {
		l.Warn("failed to get account >%v<", err)
		return nil, fmt.Errorf("failed to get account: %w", err)
	}

	gameRec, err := t.Data.GetGameRecByID(subscriptionRec.GameID)
	if err != nil {
		l.Warn("failed to get game >%v<", err)
		return nil, fmt.Errorf("failed to get game: %w", err)
	}

	// Get account contact for address information
	accountContactRec, err := t.Data.GetAccountContactRecByAccountID(accountRec.ID)
	if err != nil {
		l.Warn("failed to get account contact >%v<", err)
		return nil, fmt.Errorf("failed to get account contact: %w", err)
	}

	// Type assert to adventure game scan data
	adventureScanData, ok := joinGameScanData.(*turn_sheet.AdventureGameJoinGameScanData)
	if !ok {
		l.Warn("invalid scan data type for join game subscription, expected *turn_sheet.AdventureGameJoinGameScanData")
		return nil, fmt.Errorf("invalid scan data type for join game subscription")
	}

	// Construct join game scan data
	scanData := turn_sheet.AdventureGameJoinGameScanData{
		JoinGameScanData: turn_sheet.JoinGameScanData{
			Email:              adventureScanData.Email,
			Name:               adventureScanData.Name,
			PostalAddressLine1: adventureScanData.PostalAddressLine1,
			PostalAddressLine2: adventureScanData.PostalAddressLine2,
			StateProvince:      adventureScanData.StateProvince,
			Country:            adventureScanData.Country,
			PostalCode:         adventureScanData.PostalCode,
		},
		CharacterName: adventureScanData.CharacterName,
	}

	// Use defaults from account/contact if not provided
	if scanData.Email == "" {
		scanData.Email = accountRec.Email
	}
	if scanData.Name == "" && accountContactRec != nil {
		scanData.Name = accountContactRec.Name
	}
	if scanData.PostalAddressLine1 == "" && accountContactRec != nil {
		scanData.PostalAddressLine1 = accountContactRec.PostalAddressLine1
	}
	if scanData.PostalAddressLine2 == "" && accountContactRec != nil {
		scanData.PostalAddressLine2 = accountContactRec.PostalAddressLine2.String
	}
	if scanData.StateProvince == "" && accountContactRec != nil {
		scanData.StateProvince = accountContactRec.StateProvince
	}
	if scanData.Country == "" && accountContactRec != nil {
		scanData.Country = accountContactRec.Country
	}
	if scanData.PostalCode == "" && accountContactRec != nil {
		scanData.PostalCode = accountContactRec.PostalCode
	}

	// Validate scan data
	if err := scanData.Validate(); err != nil {
		l.Warn("invalid join game scan data >%v<", err)
		return nil, fmt.Errorf("invalid join game scan data: %w", err)
	}

	// For join game turn sheets, the code is just the game ID
	// (join game codes are simple identifiers, not full encoded identifiers)
	turnSheetCode := gameRec.ID

	// Get join game data
	joinData, err := turn_sheet.GetTurnSheetJoinGameData(gameRec, turnSheetCode)
	if err != nil {
		l.Warn("failed to get join game data >%v<", err)
		return nil, fmt.Errorf("failed to get join game data: %w", err)
	}

	sheetDataBytes, err := json.Marshal(joinData)
	if err != nil {
		l.Warn("failed to marshal join game sheet data >%v<", err)
		return nil, fmt.Errorf("failed to marshal join game sheet data: %w", err)
	}

	scanDataBytes, err := json.Marshal(scanData)
	if err != nil {
		l.Warn("failed to marshal join game scan data >%v<", err)
		return nil, fmt.Errorf("failed to marshal join game scan data: %w", err)
	}

	// Create join game turn sheet record (using the domain's existing transaction)
	turnSheetRec := &game_record.GameTurnSheet{
		GameID:           gameRec.ID,
		AccountID:        accountRec.ID,
		TurnNumber:       0,
		SheetType:        adventure_game_record.AdventureGameTurnSheetTypeJoinGame,
		SheetOrder:       1,
		SheetData:        json.RawMessage(sheetDataBytes),
		ScannedData:      json.RawMessage(scanDataBytes),
		ProcessingStatus: game_record.TurnSheetProcessingStatusPending,
	}
	turnSheetRec.ScannedAt = sql.NullTime{Time: time.Now(), Valid: true}

	createdTurnSheetRec, err := m.CreateGameTurnSheetRec(turnSheetRec)
	if err != nil {
		l.Warn("failed to create game turn sheet record >%#v< >%v<", turnSheetRec, err)
		return nil, fmt.Errorf("failed to create game turn sheet record: %w", err)
	}

	// Add to data store
	t.Data.AddGameTurnSheetRec(createdTurnSheetRec)
	t.teardownData.AddGameTurnSheetRec(createdTurnSheetRec)

	// Process the subscription using the job worker
	worker, err := jobworker.NewGameSubscriptionProcessingWorker(t.Log, t.Config, t.Store)
	if err != nil {
		l.Warn("failed to create game subscription processing worker >%v<", err)
		return nil, fmt.Errorf("failed to create game subscription processing worker: %w", err)
	}

	// Create job args
	jobArgs := jobworker.GameSubscriptionProcessingWorkerArgs{
		GameSubscriptionID: subscriptionRec.ID,
	}

	// Create a minimal job struct for DoWork
	job := &river.Job[jobworker.GameSubscriptionProcessingWorkerArgs]{
		Args: jobArgs,
	}

	// Call DoWork directly using the harness's JobClient with the same domain/transaction
	_, err = worker.DoWork(ctx, m, t.JobClient, job)
	if err != nil {
		l.Warn("failed to process subscription >%v<", err)
		return nil, fmt.Errorf("failed to process subscription: %w", err)
	}

	l.Info("processed join game subscription >%s<", subscriptionRec.ID)

	return createdTurnSheetRec, nil
}

func (t *Testing) getTurnSheetsForTurn(gameInstanceID string, turnNumber int) ([]*game_record.GameTurnSheet, error) {
	dom := t.Domain.(*domain.Domain)
	turnSheets, err := dom.GetGameTurnSheetRecsByGameInstance(gameInstanceID, turnNumber)
	if err != nil {
		return nil, err
	}
	if len(turnSheets) == 0 {
		return nil, fmt.Errorf("no turn sheets found for game instance >%s< turn >%d<", gameInstanceID, turnNumber)
	}
	for _, rec := range turnSheets {
		t.Data.AddGameTurnSheetRec(rec)
		t.teardownData.AddGameTurnSheetRec(rec)
	}
	return turnSheets, nil
}

func (t *Testing) applyConfiguredScanData(ctx context.Context, turnCfg GameTurnConfig) (bool, error) {
	if len(turnCfg.AdventureGameTurnSheetConfigs) == 0 {
		return false, nil
	}

	shouldProcess := true

	for _, cfg := range turnCfg.AdventureGameTurnSheetConfigs {
		ref := cfg.GameTurnSheetConfig.Reference
		if cfg.GameTurnSheetConfig.ScanDataConfig == nil {
			shouldProcess = false
			continue
		}
		if ref == "" {
			return false, fmt.Errorf("turn %d has scan data but missing turn sheet reference", turnCfg.TurnNumber)
		}
		turnSheetID, ok := t.Data.Refs.GameTurnSheetRefs[ref]
		if !ok {
			return false, fmt.Errorf("turn sheet reference >%s< not found for turn %d", ref, turnCfg.TurnNumber)
		}
		if err := t.applyScanDataToTurnSheet(ctx, turnSheetID, cfg.GameTurnSheetConfig.ScanDataConfig); err != nil {
			return false, err
		}
	}

	return shouldProcess, nil
}

func (t *Testing) applyScanDataToTurnSheet(_ context.Context, turnSheetID string, scanDataConfig any) error {
	l := t.Logger("applyScanDataToTurnSheet")

	m := t.Domain.(*domain.Domain)

	turnSheetRec, err := m.GetGameTurnSheetRec(turnSheetID, nil)
	if err != nil {
		l.Warn("failed to get turn sheet >%s< >%v<", turnSheetID, err)
		return fmt.Errorf("failed to get turn sheet: %w", err)
	}

	var scanDataBytes []byte
	switch turnSheetRec.SheetType {
	case adventure_game_record.AdventureGameTurnSheetTypeLocationChoice:
		locationChoiceScanData, ok := scanDataConfig.(*turn_sheet.LocationChoiceScanData)
		if !ok {
			l.Warn("invalid scan data type for location choice turn sheet, expected *turn_sheet.LocationChoiceScanData")
			return fmt.Errorf("invalid scan data type for location choice turn sheet")
		}
		scanDataBytes, err = json.Marshal(locationChoiceScanData)
		if err != nil {
			l.Warn("failed to marshal location choice scan data >%v<", err)
			return fmt.Errorf("failed to marshal location choice scan data: %w", err)
		}
	default:
		return fmt.Errorf("unsupported sheet type: %s", turnSheetRec.SheetType)
	}

	if err := m.MarkGameTurnSheetAsCompleted(turnSheetRec.ID, scanDataBytes); err != nil {
		l.Warn("failed to mark turn sheet as completed >%v<", err)
		return fmt.Errorf("failed to mark turn sheet as completed: %w", err)
	}

	l.Info("applied scan data to turn sheet >%s<", turnSheetRec.ID)

	return nil
}

// TODO: REMOVE THIS FUNCTION IF NOT USED

// func (t *Testing) processGameTurnForInstance(ctx context.Context, gameInstanceID string, turnNumber int) error {
// 	l := t.Logger("processGameTurnForInstance")

// 	worker, err := jobworker.NewGameTurnProcessingWorker(t.Log, t.Config, t.Store)
// 	if err != nil {
// 		l.Warn("failed to create game turn processing worker >%v<", err)
// 		return fmt.Errorf("failed to create game turn processing worker: %w", err)
// 	}

// 	tx, err := t.Store.BeginTx()
// 	if err != nil {
// 		l.Warn("failed to begin transaction for turn processing >%v<", err)
// 		return fmt.Errorf("failed to begin transaction for turn processing: %w", err)
// 	}
// 	defer tx.Rollback(ctx)

// 	jobDomain, err := domain.NewDomain(t.Log, t.Config)
// 	if err != nil {
// 		l.Warn("failed to create domain for turn processing >%v<", err)
// 		return fmt.Errorf("failed to create domain for turn processing: %w", err)
// 	}

// 	if err := jobDomain.Init(tx); err != nil {
// 		l.Warn("failed to init domain for turn processing >%v<", err)
// 		return fmt.Errorf("failed to init domain for turn processing: %w", err)
// 	}

// 	job := &river.Job[jobworker.GameTurnProcessingWorkerArgs]{
// 		Args: jobworker.GameTurnProcessingWorkerArgs{
// 			GameInstanceID: gameInstanceID,
// 			TurnNumber:     turnNumber,
// 		},
// 	}

// 	if _, err := worker.DoWork(ctx, jobDomain, t.JobClient, job); err != nil {
// 		l.Warn("failed to process turn >%d< for game instance >%s< >%v<", turnNumber, gameInstanceID, err)
// 		return fmt.Errorf("failed to process turn: %w", err)
// 	}

// 	if err := tx.Commit(ctx); err != nil {
// 		l.Warn("failed to commit turn processing transaction >%v<", err)
// 		return fmt.Errorf("failed to commit turn processing transaction: %w", err)
// 	}

// 	return nil
// }
