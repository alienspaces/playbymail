package harness

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/riverqueue/river"

	"gitlab.com/alienspaces/playbymail/core/nullstring"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/jobworker"
	"gitlab.com/alienspaces/playbymail/internal/jobworker/adventure_game"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
	"gitlab.com/alienspaces/playbymail/internal/turnsheet"
	"gitlab.com/alienspaces/playbymail/internal/utils/turnsheetutil"
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

	l.Info("created %d turn sheets for game instance >%s< in current transaction", len(turnSheetRecs), gameInstanceRef)

	return turnSheetRecs, nil
}

// assignTurnSheetRefs maps turn sheet references to their IDs by looking up which character instance
// owns each turn sheet, then matching against TurnSheetRefConfigs.
func (t *Testing) assignTurnSheetRefs(refConfigs []TurnSheetRefConfig, turnSheetRecs []*game_record.GameTurnSheet) error {
	if len(refConfigs) == 0 {
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
	}

	for _, cfg := range refConfigs {
		if cfg.Reference == "" || cfg.GameCharacterInstanceRef == "" {
			continue
		}
		characterInstanceID, ok := t.Data.Refs.AdventureGameCharacterInstanceRefs[cfg.GameCharacterInstanceRef]
		if !ok {
			return fmt.Errorf("adventure game character instance ref >%s< not found", cfg.GameCharacterInstanceRef)
		}
		gameTurnSheetID, ok := characterTurnSheetMap[characterInstanceID]
		if !ok {
			return fmt.Errorf("no turn sheet found for character ref >%s<", cfg.GameCharacterInstanceRef)
		}
		t.Data.Refs.GameTurnSheetRefs[cfg.Reference] = gameTurnSheetID
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
	accountRec, err := t.Data.GetAccountUserRecByID(subscriptionRec.AccountID)
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
	accountUserContactRec, err := t.Data.GetAccountUserContactRecByAccountUserID(accountRec.ID)
	if err != nil {
		l.Warn("failed to get account contact >%v<", err)
		return nil, fmt.Errorf("failed to get account contact: %w", err)
	}

	// Type assert to adventure game scan data
	adventureScanData, ok := joinGameScanData.(*turnsheet.AdventureGameJoinGameScanData)
	if !ok {
		l.Warn("invalid scan data type for join game subscription, expected *turn_sheet.AdventureGameJoinGameScanData")
		return nil, fmt.Errorf("invalid scan data type for join game subscription")
	}

	// Find manager subscription for this game
	var managerSubscriptionID string
	for _, subRec := range t.Data.GameSubscriptionRecs {
		if subRec.GameID == gameRec.ID && subRec.SubscriptionType == game_record.GameSubscriptionTypeManager {
			managerSubscriptionID = subRec.ID
			break
		}
	}

	if managerSubscriptionID == "" {
		l.Warn("no manager subscription found for game >%s<", gameRec.ID)
		return nil, fmt.Errorf("no manager subscription found for game")
	}

	// Construct join game scan data
	scanData := turnsheet.AdventureGameJoinGameScanData{
		JoinGameScanData: turnsheet.JoinGameScanData{
			GameSubscriptionID: managerSubscriptionID,
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
	if scanData.Name == "" && accountUserContactRec != nil {
		scanData.Name = nullstring.ToString(accountUserContactRec.Name)
	}
	if scanData.PostalAddressLine1 == "" && accountUserContactRec != nil {
		scanData.PostalAddressLine1 = nullstring.ToString(accountUserContactRec.PostalAddressLine1)
	}
	if scanData.PostalAddressLine2 == "" && accountUserContactRec != nil {
		scanData.PostalAddressLine2 = nullstring.ToString(accountUserContactRec.PostalAddressLine2)
	}
	if scanData.StateProvince == "" && accountUserContactRec != nil {
		scanData.StateProvince = nullstring.ToString(accountUserContactRec.StateProvince)
	}
	if scanData.Country == "" && accountUserContactRec != nil {
		scanData.Country = nullstring.ToString(accountUserContactRec.Country)
	}
	if scanData.PostalCode == "" && accountUserContactRec != nil {
		scanData.PostalCode = nullstring.ToString(accountUserContactRec.PostalCode)
	}

	// Validate scan data
	if err := scanData.Validate(); err != nil {
		l.Warn("invalid join game scan data >%v<", err)
		return nil, fmt.Errorf("invalid join game scan data: %w", err)
	}

	// Generate join game turn sheet code using the proper encoding format
	// (base64-encoded JSON with checksum, same format as playing codes but
	// with only game ID and manager subscription ID populated)
	turnSheetCode, err := turnsheetutil.GenerateJoinGameTurnSheetCode(managerSubscriptionID)
	if err != nil {
		l.Warn("failed to generate join turn sheet code >%v<", err)
		return nil, fmt.Errorf("failed to generate join turn sheet code: %w", err)
	}

	// Get join game data
	joinData, err := turnsheet.GetTurnSheetJoinGameData(gameRec, turnSheetCode)
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
		AccountUserID:    subscriptionRec.AccountUserID,
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
