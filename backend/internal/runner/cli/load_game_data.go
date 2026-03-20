package runner

import (
	"fmt"

	"github.com/urfave/cli/v2"

	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/harness"
	"gitlab.com/alienspaces/playbymail/internal/record/account_record"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
	"gitlab.com/alienspaces/playbymail/internal/runner/cli/demo_scenarios"
)

// loadGameData loads a demo game into the target database.
// Game name is required (use --list to see options). Loaded games are draft unless --publish is set.
func (rnr *Runner) loadGameData(c *cli.Context) error {
	l := loggerWithFunctionContext(rnr.Log, "loadGameData")

	if c.Bool("list") {
		return rnr.listDemoGames(c)
	}

	gameName := c.String("game")
	if gameName == "" {
		return fmt.Errorf("--game is required (use --list to see available demo games)")
	}

	entry, ok := LookupDemoGame(gameName)
	if !ok {
		l.Warn("unknown demo game >%s<", gameName)
		return fmt.Errorf("unknown demo game %q (use --list to see available demo games)", gameName)
	}

	l.Info("** Load Demo Game: %s **", gameName)
	config := entry.Config()

	err := rnr.InitDomain()
	if err != nil {
		l.Warn("failed domain init >%v<", err)
		return err
	}

	// Ensure demo accounts exist -- create if missing, reuse if present.
	demoRecs, err := rnr.ensureDemoAccounts()
	if err != nil {
		l.Warn("failed ensuring demo accounts >%v<", err)
		return err
	}

	// Populate each top-level subscription Record with account IDs (same order as DemoAccountDefs).
	for i := range config.AccountUserGameSubscriptionConfigs {
		if i >= len(demoRecs.AccountUsers) {
			break
		}
		rec := config.AccountUserGameSubscriptionConfigs[i].Record
		if rec != nil {
			rec.AccountID = demoRecs.AccountUsers[i].AccountID
			rec.AccountUserID = demoRecs.AccountUsers[i].ID
		}
	}

	for i := range config.GameConfigs {
		gc := &config.GameConfigs[i]
		if gc.Record == nil || gc.Record.Name == "" {
			continue
		}
		existing, err := rnr.gameExistsByName(gc.Record.Name)
		if err != nil {
			return err
		}
		if existing && !c.Bool("replace") {
			return fmt.Errorf("game %q already exists; use --replace to overwrite", gc.Record.Name)
		}
		if existing && c.Bool("replace") {
			if err := rnr.removeGameByName(gc.Record.Name); err != nil {
				l.Warn("failed removing existing game >%s<: %v", gc.Record.Name, err)
				return err
			}
		}
	}

	if c.Bool("replace") {
		if err := rnr.Domain.Commit(); err != nil {
			l.Warn("failed committing game removal >%v<", err)
			return err
		}
		if err := rnr.InitDomainTx(); err != nil {
			l.Warn("failed re-init domain tx after removal >%v<", err)
			return err
		}
	}

	testHarness, err := harness.NewTesting(rnr.Config, rnr.Log, rnr.Store, rnr.JobClient, rnr.Scanner, config)
	if err != nil {
		l.Warn("failed new testing harness >%v<", err)
		return err
	}

	testHarness.ShouldCommitData = true

	_, err = testHarness.Setup()
	if err != nil {
		l.Warn("failed harness setup >%v<", err)
		return err
	}

	// Publish games via the domain layer if requested
	if c.Bool("publish") {
		if err := rnr.InitDomainTx(); err != nil {
			l.Warn("failed re-init domain tx for publish >%v<", err)
			return err
		}
		dm, ok := rnr.Domain.(*domain.Domain)
		if !ok {
			return fmt.Errorf("cannot publish: domain type assertion failed")
		}
		for _, rec := range testHarness.Data.GameRecs {
			rec.Status = game_record.GameStatusPublished
			_, err = dm.UpdateGameRec(rec)
			if err != nil {
				l.Warn("failed publishing game %s >%v<", rec.ID, err)
				return err
			}
			l.Info("published game %s (%s)", rec.Name, rec.ID)
		}
		if err := rnr.Domain.Commit(); err != nil {
			l.Warn("failed committing game publish >%v<", err)
			return err
		}
	}

	l.Info("game data loaded successfully")

	return nil
}

// DemoAccountRecords holds all account, account user, and account user contact
// records for the demo accounts (one per DemoAccountDef, same order).
type DemoAccountRecords struct {
	Accounts            []*account_record.Account
	AccountUsers        []*account_record.AccountUser
	AccountUserContacts []*account_record.AccountUserContact
}

// ensureDemoAccounts ensures each demo account, account user, and account_user_contact exist;
// it returns all such records (create or fetch), one set per DemoAccountDef in order.
func (rnr *Runner) ensureDemoAccounts() (*DemoAccountRecords, error) {
	l := loggerWithFunctionContext(rnr.Log, "ensureDemoAccounts")

	dm, ok := rnr.Domain.(*domain.Domain)
	if !ok {
		return nil, fmt.Errorf("domain type assertion failed")
	}

	n := len(demo_scenarios.DemoAccountDefs)
	accounts := make([]*account_record.Account, 0, n)
	accountUsers := make([]*account_record.AccountUser, 0, n)
	accountUserContacts := make([]*account_record.AccountUserContact, 0, n)

	for _, def := range demo_scenarios.DemoAccountDefs {
		existing, err := dm.GetAccountUserRecByEmail(def.Email)
		if err != nil {
			return nil, fmt.Errorf("failed looking up account by email >%s<: %w", def.Email, err)
		}
		if existing != nil {
			l.Info("demo account already exists ref >%s< email >%s< ID >%s<", def.Ref, def.Email, existing.ID)
			acctRec, err := dm.GetAccountRec(existing.AccountID, nil)
			if err != nil {
				return nil, fmt.Errorf("failed getting account for email >%s<: %w", def.Email, err)
			}
			contactRec, err := dm.GetAccountUserContactRecByAccountUserID(existing.ID, nil)
			if err != nil {
				return nil, fmt.Errorf("failed getting account user contact for email >%s<: %w", def.Email, err)
			}
			accounts = append(accounts, acctRec)
			accountUsers = append(accountUsers, existing)
			accountUserContacts = append(accountUserContacts, contactRec)
			continue
		}

		l.Info("creating demo account ref >%s< email >%s<", def.Ref, def.Email)

		acctRec, accountUserRec, contactRec, _, err := dm.UpsertAccount(
			&account_record.Account{},
			&account_record.AccountUser{
				Email: def.Email,
			},
			nil,
		)
		if err != nil {
			return nil, fmt.Errorf("failed creating demo account >%s<: %w", def.Email, err)
		}

		l.Info("created demo account ref >%s< ID >%s<", def.Ref, accountUserRec.ID)
		accounts = append(accounts, acctRec)
		accountUsers = append(accountUsers, accountUserRec)
		accountUserContacts = append(accountUserContacts, contactRec)
	}

	// Commit the account creation and re-init so the harness starts with a clean tx
	if err := rnr.Domain.Commit(); err != nil {
		return nil, fmt.Errorf("failed committing demo accounts: %w", err)
	}

	if err := rnr.InitDomainTx(); err != nil {
		return nil, fmt.Errorf("failed re-init domain tx after account creation: %w", err)
	}

	return &DemoAccountRecords{
		Accounts:            accounts,
		AccountUsers:        accountUsers,
		AccountUserContacts: accountUserContacts,
	}, nil
}

// gameExistsByName returns true if at least one non-deleted game with the
// given name exists in the database.
func (rnr *Runner) gameExistsByName(name string) (bool, error) {

	dm, ok := rnr.Domain.(*domain.Domain)
	if !ok {
		return false, fmt.Errorf("domain type assertion failed")
	}

	games, err := dm.GetManyGameRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: game_record.FieldGameName, Val: name},
		},
	})
	if err != nil {
		return false, fmt.Errorf("failed checking for existing game >%s<: %w", name, err)
	}

	return len(games) > 0, nil
}

// removeGameByName finds a game by exact name and removes it together with
// all dependent records (subscriptions, instances, adventure game data, images).
// It is a no-op when no game with the given name exists.
func (rnr *Runner) removeGameByName(name string) error {
	l := loggerWithFunctionContext(rnr.Log, "removeGameByName")

	dm, ok := rnr.Domain.(*domain.Domain)
	if !ok {
		return fmt.Errorf("domain type assertion failed")
	}

	games, err := dm.GetManyGameRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: game_record.FieldGameName, Val: name},
		},
	})
	if err != nil {
		return fmt.Errorf("failed looking up game by name >%s<: %w", name, err)
	}

	if len(games) == 0 {
		l.Info("no existing game named >%s<, nothing to remove", name)
		return nil
	}

	for _, gameRec := range games {
		l.Info("removing game >%s< (ID %s)", gameRec.Name, gameRec.ID)
		if err := rnr.removeGameAndDependents(dm, gameRec.ID); err != nil {
			return err
		}
	}

	return nil
}

// TODO: (agent) Move this cascade into domain (e.g. domain.RemoveGameAndDependents(ctx, gameID)) and call it from here; add integration or harness tests that create a game with dependents, call the method, and assert full teardown. Remove this private helper once done.

// removeGameAndDependents cascades removal of all records that belong to a
// single game. The order mirrors the harness RemoveData teardown sequence.
func (rnr *Runner) removeGameAndDependents(dm *domain.Domain, gameID string) error {
	l := loggerWithFunctionContext(rnr.Log, "removeGameAndDependents")

	byGame := &coresql.Options{
		Params: []coresql.Param{{Col: "game_id", Val: gameID}},
	}

	// 1. Adventure game turn sheets (linked by game_id, not game_instance_id)
	agts, err := dm.GetManyAdventureGameTurnSheetRecs(byGame)
	if err != nil {
		return fmt.Errorf("failed getting adventure game turn sheets: %w", err)
	}
	for _, rec := range agts {
		if err := dm.RemoveAdventureGameTurnSheetRec(rec.ID); err != nil {
			return fmt.Errorf("failed removing adventure game turn sheet >%s<: %w", rec.ID, err)
		}
	}

	// 2. Game instances and their dependents
	instances, err := dm.GetManyGameInstanceRecs(&coresql.Options{
		Params: []coresql.Param{{Col: game_record.FieldGameInstanceGameID, Val: gameID}},
	})
	if err != nil {
		return fmt.Errorf("failed getting game instances: %w", err)
	}
	for _, inst := range instances {
		if err := rnr.removeGameInstanceDependents(dm, inst.ID); err != nil {
			return err
		}
		l.Info("removing game instance >%s<", inst.ID)
		if err := dm.RemoveGameInstanceRec(inst.ID); err != nil {
			return fmt.Errorf("failed removing game instance >%s<: %w", inst.ID, err)
		}
	}

	// 3. Game subscriptions (subscription instances were removed with game instances above)
	subs, err := dm.GetManyGameSubscriptionRecs(&coresql.Options{
		Params: []coresql.Param{{Col: game_record.FieldGameSubscriptionGameID, Val: gameID}},
	})
	if err != nil {
		return fmt.Errorf("failed getting game subscriptions: %w", err)
	}
	for _, sub := range subs {
		l.Info("removing game subscription >%s<", sub.ID)
		if err := dm.RemoveGameSubscriptionRec(sub.ID); err != nil {
			return fmt.Errorf("failed removing game subscription >%s<: %w", sub.ID, err)
		}
	}

	// 4–12. Adventure game definition data (objects first — effects reference links, then placements, links, locations, etc.)
	if err := rnr.removeAdventureGameDefinitionData(dm, gameID); err != nil {
		return err
	}

	// 12. Game images
	images, err := dm.GetManyGameImageRecs(&coresql.Options{
		Params: []coresql.Param{{Col: game_record.FieldGameImageGameID, Val: gameID}},
	})
	if err != nil {
		return fmt.Errorf("failed getting game images: %w", err)
	}
	for _, rec := range images {
		if err := dm.RemoveGameImageRec(rec.ID); err != nil {
			return fmt.Errorf("failed removing game image >%s<: %w", rec.ID, err)
		}
	}

	// 13. Game record itself
	l.Info("removing game record >%s<", gameID)
	if err := dm.RemoveGameRec(gameID); err != nil {
		return fmt.Errorf("failed removing game >%s<: %w", gameID, err)
	}

	return nil
}

// removeAdventureGameDefinitionData removes adventure-game-specific definition
// records (placements, links, objects, locations, items, creatures, characters)
// in FK-safe order.
func (rnr *Runner) removeAdventureGameDefinitionData(dm *domain.Domain, gameID string) error {
	byGame := &coresql.Options{
		Params: []coresql.Param{{Col: "game_id", Val: gameID}},
	}

	// Location object effects reference location links via result_adventure_game_location_link_id,
	// so objects must be removed before links.
	if err := rnr.removeAdventureGameLocationObjects(dm, byGame); err != nil {
		return err
	}

	// Item effects reference location links via result_adventure_game_location_link_id AND
	// reference items, so they must be removed before both links and items.
	itemEffects, err := dm.GetManyAdventureGameItemEffectRecs(byGame)
	if err != nil {
		return fmt.Errorf("failed getting item effects: %w", err)
	}
	for _, rec := range itemEffects {
		if err := dm.RemoveAdventureGameItemEffectRec(rec.ID); err != nil {
			return fmt.Errorf("failed removing item effect >%s<: %w", rec.ID, err)
		}
	}

	// Placements
	itemPlacements, err := dm.GetManyAdventureGameItemPlacementRecs(byGame)
	if err != nil {
		return fmt.Errorf("failed getting item placements: %w", err)
	}
	for _, rec := range itemPlacements {
		if err := dm.RemoveAdventureGameItemPlacementRec(rec.ID); err != nil {
			return fmt.Errorf("failed removing item placement >%s<: %w", rec.ID, err)
		}
	}

	creaturePlacements, err := dm.GetManyAdventureGameCreaturePlacementRecs(byGame)
	if err != nil {
		return fmt.Errorf("failed getting creature placements: %w", err)
	}
	for _, rec := range creaturePlacements {
		if err := dm.RemoveAdventureGameCreaturePlacementRec(rec.ID); err != nil {
			return fmt.Errorf("failed removing creature placement >%s<: %w", rec.ID, err)
		}
	}

	// Location links and requirements
	linkReqs, err := dm.GetManyAdventureGameLocationLinkRequirementRecs(byGame)
	if err != nil {
		return fmt.Errorf("failed getting link requirements: %w", err)
	}
	for _, rec := range linkReqs {
		if err := dm.RemoveAdventureGameLocationLinkRequirementRec(rec.ID); err != nil {
			return fmt.Errorf("failed removing link requirement >%s<: %w", rec.ID, err)
		}
	}

	links, err := dm.GetManyAdventureGameLocationLinkRecs(&coresql.Options{
		Params: []coresql.Param{{Col: adventure_game_record.FieldAdventureGameLocationLinkGameID, Val: gameID}},
	})
	if err != nil {
		return fmt.Errorf("failed getting location links: %w", err)
	}
	for _, rec := range links {
		if err := dm.RemoveAdventureGameLocationLinkRec(rec.ID); err != nil {
			return fmt.Errorf("failed removing location link >%s<: %w", rec.ID, err)
		}
	}

	// Locations, items, creatures, characters
	return rnr.removeAdventureGameEntities(dm, byGame)
}

func (rnr *Runner) removeAdventureGameLocationObjects(dm *domain.Domain, byGame *coresql.Options) error {
	objEffects, err := dm.GetManyAdventureGameLocationObjectEffectRecs(byGame)
	if err != nil {
		return fmt.Errorf("failed getting location object effects: %w", err)
	}
	for _, rec := range objEffects {
		if err := dm.RemoveAdventureGameLocationObjectEffectRec(rec.ID); err != nil {
			return fmt.Errorf("failed removing location object effect >%s<: %w", rec.ID, err)
		}
	}

	// Objects and states have a circular FK dependency:
	//   state.adventure_game_location_object_id → object.id
	//   object.initial_adventure_game_location_object_state_id → state.id
	// Break the cycle by clearing the initial_state_id FK on every object before deleting states.
	objs, err := dm.GetManyAdventureGameLocationObjectRecs(byGame)
	if err != nil {
		return fmt.Errorf("failed getting location objects: %w", err)
	}
	for _, rec := range objs {
		if rec.InitialAdventureGameLocationObjectStateID.Valid {
			rec.InitialAdventureGameLocationObjectStateID = adventure_game_record.AdventureGameLocationObject{}.InitialAdventureGameLocationObjectStateID
			if _, err := dm.UpdateAdventureGameLocationObjectRec(rec); err != nil {
				return fmt.Errorf("failed clearing initial state on location object >%s<: %w", rec.ID, err)
			}
		}
	}

	states, err := dm.GetManyAdventureGameLocationObjectStateRecs(byGame)
	if err != nil {
		return fmt.Errorf("failed getting location object states: %w", err)
	}
	for _, rec := range states {
		if err := dm.RemoveAdventureGameLocationObjectStateRec(rec.ID); err != nil {
			return fmt.Errorf("failed removing location object state >%s<: %w", rec.ID, err)
		}
	}

	for _, rec := range objs {
		if err := dm.RemoveAdventureGameLocationObjectRec(rec.ID); err != nil {
			return fmt.Errorf("failed removing location object >%s<: %w", rec.ID, err)
		}
	}
	return nil
}

func (rnr *Runner) removeAdventureGameEntities(dm *domain.Domain, byGame *coresql.Options) error {
	locs, err := dm.GetManyAdventureGameLocationRecs(byGame)
	if err != nil {
		return fmt.Errorf("failed getting locations: %w", err)
	}
	for _, rec := range locs {
		if err := dm.RemoveAdventureGameLocationRec(rec.ID); err != nil {
			return fmt.Errorf("failed removing location >%s<: %w", rec.ID, err)
		}
	}

	items, err := dm.GetManyAdventureGameItemRecs(byGame)
	if err != nil {
		return fmt.Errorf("failed getting items: %w", err)
	}
	for _, rec := range items {
		if err := dm.RemoveAdventureGameItemRec(rec.ID); err != nil {
			return fmt.Errorf("failed removing item >%s<: %w", rec.ID, err)
		}
	}

	creatures, err := dm.GetManyAdventureGameCreatureRecs(byGame)
	if err != nil {
		return fmt.Errorf("failed getting creatures: %w", err)
	}
	for _, rec := range creatures {
		if err := dm.RemoveAdventureGameCreatureRec(rec.ID); err != nil {
			return fmt.Errorf("failed removing creature >%s<: %w", rec.ID, err)
		}
	}

	chars, err := dm.GetManyAdventureGameCharacterRecs(byGame)
	if err != nil {
		return fmt.Errorf("failed getting characters: %w", err)
	}
	for _, rec := range chars {
		if err := dm.RemoveAdventureGameCharacterRec(rec.ID); err != nil {
			return fmt.Errorf("failed removing character >%s<: %w", rec.ID, err)
		}
	}
	return nil
}

// removeGameInstanceDependents removes all records that depend on a single
// game instance: turn sheets, adventure game instance records, subscription
// instances, and instance parameters. Adventure game turn sheets are removed
// at the game level in removeGameAndDependents.
func (rnr *Runner) removeGameInstanceDependents(dm *domain.Domain, instanceID string) error {
	byInstance := &coresql.Options{
		Params: []coresql.Param{{Col: "game_instance_id", Val: instanceID}},
	}

	// Game turn sheets
	turnSheets, err := dm.GetManyGameTurnSheetRecs(byInstance)
	if err != nil {
		return fmt.Errorf("failed getting game turn sheets: %w", err)
	}
	for _, rec := range turnSheets {
		if err := dm.RemoveGameTurnSheetRec(rec.ID); err != nil {
			return fmt.Errorf("failed removing game turn sheet >%s<: %w", rec.ID, err)
		}
	}

	// Adventure game item instances must be removed before character instances
	// (item_instance.adventure_game_character_instance_id FK references character_instance)
	itemInsts, err := dm.GetManyAdventureGameItemInstanceRecs(byInstance)
	if err != nil {
		return fmt.Errorf("failed getting item instances: %w", err)
	}
	for _, rec := range itemInsts {
		if err := dm.RemoveAdventureGameItemInstanceRec(rec.ID); err != nil {
			return fmt.Errorf("failed removing item instance >%s<: %w", rec.ID, err)
		}
	}

	// Adventure game character instances
	charInsts, err := dm.GetManyAdventureGameCharacterInstanceRecs(byInstance)
	if err != nil {
		return fmt.Errorf("failed getting character instances: %w", err)
	}
	for _, rec := range charInsts {
		if err := dm.RemoveAdventureGameCharacterInstanceRec(rec.ID); err != nil {
			return fmt.Errorf("failed removing character instance >%s<: %w", rec.ID, err)
		}
	}

	// Adventure game creature instances
	creatInsts, err := dm.GetManyAdventureGameCreatureInstanceRecs(byInstance)
	if err != nil {
		return fmt.Errorf("failed getting creature instances: %w", err)
	}
	for _, rec := range creatInsts {
		if err := dm.RemoveAdventureGameCreatureInstanceRec(rec.ID); err != nil {
			return fmt.Errorf("failed removing creature instance >%s<: %w", rec.ID, err)
		}
	}

	// Adventure game location object instances (must precede location instances)
	objInsts, err := dm.GetManyAdventureGameLocationObjectInstanceRecs(byInstance)
	if err != nil {
		return fmt.Errorf("failed getting location object instances: %w", err)
	}
	for _, rec := range objInsts {
		if err := dm.RemoveAdventureGameLocationObjectInstanceRec(rec.ID); err != nil {
			return fmt.Errorf("failed removing location object instance >%s<: %w", rec.ID, err)
		}
	}

	// Adventure game location instances
	locInsts, err := dm.GetManyAdventureGameLocationInstanceRecs(byInstance)
	if err != nil {
		return fmt.Errorf("failed getting location instances: %w", err)
	}
	for _, rec := range locInsts {
		if err := dm.RemoveAdventureGameLocationInstanceRec(rec.ID); err != nil {
			return fmt.Errorf("failed removing location instance >%s<: %w", rec.ID, err)
		}
	}

	// Game subscription instances
	subInsts, err := dm.GetManyGameSubscriptionInstanceRecs(&coresql.Options{
		Params: []coresql.Param{{Col: game_record.FieldGameSubscriptionInstanceGameInstanceID, Val: instanceID}},
	})
	if err != nil {
		return fmt.Errorf("failed getting subscription instances: %w", err)
	}
	for _, rec := range subInsts {
		if err := dm.RemoveGameSubscriptionInstanceRec(rec.ID); err != nil {
			return fmt.Errorf("failed removing subscription instance >%s<: %w", rec.ID, err)
		}
	}

	// Game instance parameters
	params, err := dm.GetManyGameInstanceParameterRecs(&coresql.Options{
		Params: []coresql.Param{{Col: game_record.FieldGameInstanceParameterGameInstanceID, Val: instanceID}},
	})
	if err != nil {
		return fmt.Errorf("failed getting instance parameters: %w", err)
	}
	for _, rec := range params {
		if err := dm.RemoveGameInstanceParameterRec(rec.ID); err != nil {
			return fmt.Errorf("failed removing instance parameter >%s<: %w", rec.ID, err)
		}
	}

	return nil
}

// listDemoGames prints registered demo game summaries to stdout.
func (rnr *Runner) listDemoGames(_ *cli.Context) error {
	games := ListDemoGames()
	fmt.Println("Available demo games:")
	fmt.Println()
	for _, g := range games {
		fmt.Printf("  %s (%s)\n", g.Name, g.GameType)
		fmt.Printf("    %s\n\n", g.Description)
	}
	fmt.Println("Usage:")
	fmt.Printf("  db-load-demo-game --game \"<name>\"\n")
	fmt.Printf("  db-load-demo-game --game \"<name>\" --publish\n")
	fmt.Printf("  db-load-demo-game --game \"<name>\" --replace --publish\n")
	fmt.Println()
	return nil
}
