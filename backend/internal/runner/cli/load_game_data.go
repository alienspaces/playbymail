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
	seedRefs, err := rnr.ensureDemoAccounts()
	if err != nil {
		l.Warn("failed ensuring demo accounts >%v<", err)
		return err
	}
	config.SeedAccountRefs = seedRefs

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

// ensureDemoAccounts checks each demo account, user account and user account contact records exist.
func (rnr *Runner) ensureDemoAccounts() (map[string]string, error) {
	l := loggerWithFunctionContext(rnr.Log, "ensureDemoAccounts")

	dm, ok := rnr.Domain.(*domain.Domain)
	if !ok {
		return nil, fmt.Errorf("domain type assertion failed")
	}

	refs := make(map[string]string, len(demo_scenarios.DemoAccountDefs))

	for _, def := range demo_scenarios.DemoAccountDefs {
		existing, err := dm.GetAccountUserRecByEmail(def.Email)
		if err != nil {
			return nil, fmt.Errorf("failed looking up account by email >%s<: %w", def.Email, err)
		}
		if existing != nil {
			l.Info("demo account already exists ref >%s< email >%s< ID >%s<", def.Ref, def.Email, existing.ID)
			refs[def.Ref] = existing.ID
			continue
		}

		l.Info("creating demo account ref >%s< email >%s<", def.Ref, def.Email)

		_, accountUserRec, _, _, err := dm.UpsertAccount(
			&account_record.Account{},
			&account_record.AccountUser{
				Email: def.Email,
			},
			nil,
		)
		if err != nil {
			return nil, fmt.Errorf("failed creating demo account >%s<: %w", def.Email, err)
		}

		refs[def.Ref] = accountUserRec.ID
		l.Info("created demo account ref >%s< ID >%s<", def.Ref, accountUserRec.ID)
	}

	// Commit the account creation and re-init so the harness starts with a clean tx
	if err := rnr.Domain.Commit(); err != nil {
		return nil, fmt.Errorf("failed committing demo accounts: %w", err)
	}

	if err := rnr.InitDomainTx(); err != nil {
		return nil, fmt.Errorf("failed re-init domain tx after account creation: %w", err)
	}

	return refs, nil
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

// TODO: This function is potentially brittle and should be implemented by a domain or harnessmethod
// that has a test backing it.

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

	// 4. Adventure game item placements
	itemPlacements, err := dm.GetManyAdventureGameItemPlacementRecs(byGame)
	if err != nil {
		return fmt.Errorf("failed getting item placements: %w", err)
	}

	for _, rec := range itemPlacements {
		if err := dm.RemoveAdventureGameItemPlacementRec(rec.ID); err != nil {
			return fmt.Errorf("failed removing item placement >%s<: %w", rec.ID, err)
		}
	}

	// 5. Adventure game creature placements
	creaturePlacements, err := dm.GetManyAdventureGameCreaturePlacementRecs(byGame)
	if err != nil {
		return fmt.Errorf("failed getting creature placements: %w", err)
	}

	for _, rec := range creaturePlacements {
		if err := dm.RemoveAdventureGameCreaturePlacementRec(rec.ID); err != nil {
			return fmt.Errorf("failed removing creature placement >%s<: %w", rec.ID, err)
		}
	}

	// 6. Adventure game location link requirements
	linkReqs, err := dm.GetManyAdventureGameLocationLinkRequirementRecs(byGame)
	if err != nil {
		return fmt.Errorf("failed getting link requirements: %w", err)
	}

	for _, rec := range linkReqs {
		if err := dm.RemoveAdventureGameLocationLinkRequirementRec(rec.ID); err != nil {
			return fmt.Errorf("failed removing link requirement >%s<: %w", rec.ID, err)
		}
	}

	// 7. Adventure game location links
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

	// 8. Adventure game locations
	locs, err := dm.GetManyAdventureGameLocationRecs(byGame)
	if err != nil {
		return fmt.Errorf("failed getting locations: %w", err)
	}

	for _, rec := range locs {
		if err := dm.RemoveAdventureGameLocationRec(rec.ID); err != nil {
			return fmt.Errorf("failed removing location >%s<: %w", rec.ID, err)
		}
	}

	// 9. Adventure game items
	items, err := dm.GetManyAdventureGameItemRecs(byGame)
	if err != nil {
		return fmt.Errorf("failed getting items: %w", err)
	}

	for _, rec := range items {
		if err := dm.RemoveAdventureGameItemRec(rec.ID); err != nil {
			return fmt.Errorf("failed removing item >%s<: %w", rec.ID, err)
		}
	}

	// 10. Adventure game creatures
	creatures, err := dm.GetManyAdventureGameCreatureRecs(byGame)
	if err != nil {
		return fmt.Errorf("failed getting creatures: %w", err)
	}

	for _, rec := range creatures {
		if err := dm.RemoveAdventureGameCreatureRec(rec.ID); err != nil {
			return fmt.Errorf("failed removing creature >%s<: %w", rec.ID, err)
		}
	}

	// 11. Adventure game characters
	chars, err := dm.GetManyAdventureGameCharacterRecs(byGame)
	if err != nil {
		return fmt.Errorf("failed getting characters: %w", err)
	}

	for _, rec := range chars {
		if err := dm.RemoveAdventureGameCharacterRec(rec.ID); err != nil {
			return fmt.Errorf("failed removing character >%s<: %w", rec.ID, err)
		}
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

	// Adventure game item instances
	itemInsts, err := dm.GetManyAdventureGameItemInstanceRecs(byInstance)
	if err != nil {
		return fmt.Errorf("failed getting item instances: %w", err)
	}
	for _, rec := range itemInsts {
		if err := dm.RemoveAdventureGameItemInstanceRec(rec.ID); err != nil {
			return fmt.Errorf("failed removing item instance >%s<: %w", rec.ID, err)
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
