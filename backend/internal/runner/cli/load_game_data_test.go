package runner

import (
	"flag"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v2"

	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/harness"
	"gitlab.com/alienspaces/playbymail/internal/record/account_record"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
	"gitlab.com/alienspaces/playbymail/internal/runner/cli/demo_scenarios"
)

// newTestRunnerWithDomain creates a Runner with the domain and transaction initialised.
func newTestRunnerWithDomain(t *testing.T) *Runner {
	t.Helper()
	cfg, l, s, j, scanner := newDefaultDependencies(t)
	r := newTestRunner(t, cfg, l, s, j, scanner)
	err := r.InitDomain()
	require.NoError(t, err, "InitDomain returns without error")
	return r
}

// removeDemoAccounts cleans up any demo accounts created during a test.
// It removes contacts, account subscriptions, the account_user, and the
// parent account record in the correct FK order.
func removeDemoAccounts(t *testing.T, rnr *Runner) {
	t.Helper()
	dm, ok := rnr.Domain.(*domain.Domain)
	require.True(t, ok)

	for _, def := range demo_scenarios.DemoAccountDefs {
		rec, err := dm.GetAccountUserRecByEmail(def.Email)
		if err != nil || rec == nil {
			continue
		}

		// Remove contacts
		contacts, _ := dm.GetManyAccountUserContactRecs(&coresql.Options{
			Params: []coresql.Param{{Col: account_record.FieldAccountUserContactAccountUserID, Val: rec.ID}},
		})
		for _, c := range contacts {
			_ = dm.RemoveAccountUserContactRec(c.ID)
		}

		// Remove account subscriptions (linked by account_id)
		subs, _ := dm.GetManyAccountSubscriptionRecs(&coresql.Options{
			Params: []coresql.Param{{Col: account_record.FieldAccountSubscriptionAccountID, Val: rec.AccountID}},
		})
		for _, s := range subs {
			_ = dm.RemoveAccountSubscriptionRec(s.ID)
		}

		// Remove account user, then parent account
		_ = dm.RemoveAccountUserRec(rec.ID)
		_ = dm.RemoveAccountRec(rec.AccountID)
	}

	err := rnr.Domain.Commit()
	require.NoError(t, err, "commit demo account removal")

	err = rnr.InitDomainTx()
	require.NoError(t, err, "re-init tx after demo account removal")
}

// newCliContext builds a *cli.Context with the given string and bool flags.
func newCliContext(stringFlags map[string]string, boolFlags map[string]bool) *cli.Context {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	for k := range stringFlags {
		fs.String(k, "", "")
	}
	for k := range boolFlags {
		fs.Bool(k, false, "")
	}

	var args []string
	for k, v := range stringFlags {
		args = append(args, fmt.Sprintf("--%s=%s", k, v))
	}
	for k, v := range boolFlags {
		if v {
			args = append(args, fmt.Sprintf("--%s", k))
		}
	}
	_ = fs.Parse(args)

	app := &cli.App{}
	return cli.NewContext(app, fs, nil)
}

// --- Tests ---

func TestEnsureDemoAccounts(t *testing.T) {
	rnr := newTestRunnerWithDomain(t)
	defer removeDemoAccounts(t, rnr)

	n := len(demo_scenarios.DemoAccountDefs)

	// First call: creates accounts
	demoRecs, err := rnr.ensureDemoAccounts()
	require.NoError(t, err, "ensureDemoAccounts succeeds on first call")
	require.Len(t, demoRecs.Accounts, n, "returns one account per def")
	require.Len(t, demoRecs.AccountUsers, n, "returns one account user per def")
	require.Len(t, demoRecs.AccountUserContacts, n, "returns one account user contact per def")

	for i, def := range demo_scenarios.DemoAccountDefs {
		require.Equal(t, def.Email, demoRecs.AccountUsers[i].Email, "account user email matches def")
		require.NotEmpty(t, demoRecs.AccountUsers[i].ID, "account user ID non-empty for %s", def.Ref)
		require.Equal(t, demoRecs.AccountUsers[i].ID, demoRecs.AccountUserContacts[i].AccountUserID, "contact belongs to account user")
		require.Equal(t, demoRecs.AccountUsers[i].AccountID, demoRecs.Accounts[i].ID, "account user belongs to account")
	}

	// Second call: idempotent -- same IDs returned
	demoRecs2, err := rnr.ensureDemoAccounts()
	require.NoError(t, err, "ensureDemoAccounts succeeds on second call")
	require.Len(t, demoRecs2.AccountUsers, n, "second call returns same count")
	for i := 0; i < n; i++ {
		require.Equal(t, demoRecs.AccountUsers[i].ID, demoRecs2.AccountUsers[i].ID, "second call returns identical account user IDs")
	}
}

func TestRemoveGameByName(t *testing.T) {
	rnr := newTestRunnerWithDomain(t)
	defer removeDemoAccounts(t, rnr)

	// Load the full demo scenario so we have a game to remove
	demoRecs, err := rnr.ensureDemoAccounts()
	require.NoError(t, err)

	config := demo_scenarios.AdventureGameConfig()
	for i := range config.AccountUserGameSubscriptionConfigs {
		if i < len(demoRecs.AccountUsers) && config.AccountUserGameSubscriptionConfigs[i].Record != nil {
			config.AccountUserGameSubscriptionConfigs[i].Record.AccountID = demoRecs.AccountUsers[i].AccountID
			config.AccountUserGameSubscriptionConfigs[i].Record.AccountUserID = demoRecs.AccountUsers[i].ID
		}
	}

	th, err := newTestHarness(t, rnr, config)
	require.NoError(t, err)
	th.ShouldCommitData = true
	_, err = th.Setup()
	require.NoError(t, err)

	// Re-init runner domain tx so it can see the committed harness data
	err = rnr.InitDomainTx()
	require.NoError(t, err)

	gameName := demo_scenarios.DemoAdventureGameName

	dm := rnr.Domain.(*domain.Domain)

	// Verify game exists
	games, err := dm.GetManyGameRecs(&coresql.Options{
		Params: []coresql.Param{{Col: game_record.FieldGameName, Val: gameName}},
	})
	require.NoError(t, err)
	require.NotEmpty(t, games, "game should exist before removal")

	// Remove game
	err = rnr.removeGameByName(gameName)
	require.NoError(t, err, "removeGameByName succeeds")

	err = rnr.Domain.Commit()
	require.NoError(t, err)
	err = rnr.InitDomainTx()
	require.NoError(t, err)

	dm = rnr.Domain.(*domain.Domain)

	// Verify game is gone
	games, err = dm.GetManyGameRecs(&coresql.Options{
		Params: []coresql.Param{{Col: game_record.FieldGameName, Val: gameName}},
	})
	require.NoError(t, err)
	require.Empty(t, games, "game should not exist after removal")

	// No-op on nonexistent game
	err = rnr.removeGameByName("nonexistent-game-name")
	require.NoError(t, err, "removeGameByName is a no-op for nonexistent game")
}

func TestRemoveGameAndDependents(t *testing.T) {
	rnr := newTestRunnerWithDomain(t)
	defer removeDemoAccounts(t, rnr)

	// Load the full scenario
	demoRecs, err := rnr.ensureDemoAccounts()
	require.NoError(t, err)

	config := demo_scenarios.AdventureGameConfig()
	for i := range config.AccountUserGameSubscriptionConfigs {
		if i < len(demoRecs.AccountUsers) && config.AccountUserGameSubscriptionConfigs[i].Record != nil {
			config.AccountUserGameSubscriptionConfigs[i].Record.AccountID = demoRecs.AccountUsers[i].AccountID
			config.AccountUserGameSubscriptionConfigs[i].Record.AccountUserID = demoRecs.AccountUsers[i].ID
		}
	}

	th, err := newTestHarness(t, rnr, config)
	require.NoError(t, err)
	th.ShouldCommitData = true
	_, err = th.Setup()
	require.NoError(t, err)

	// Re-init runner domain tx so it can see the committed harness data
	err = rnr.InitDomainTx()
	require.NoError(t, err)

	// Capture the game ID before removal
	require.NotEmpty(t, th.Data.GameRecs, "harness should have created a game")
	gameID := th.Data.GameRecs[0].ID

	dm := rnr.Domain.(*domain.Domain)

	// Verify dependent records exist before removal
	locs, err := dm.GetManyAdventureGameLocationRecs(&coresql.Options{
		Params: []coresql.Param{{Col: adventure_game_record.FieldAdventureGameLocationGameID, Val: gameID}},
	})
	require.NoError(t, err)
	require.NotEmpty(t, locs, "locations should exist before removal")

	items, err := dm.GetManyAdventureGameItemRecs(&coresql.Options{
		Params: []coresql.Param{{Col: adventure_game_record.FieldAdventureGameItemGameID, Val: gameID}},
	})
	require.NoError(t, err)
	require.NotEmpty(t, items, "items should exist before removal")

	links, err := dm.GetManyAdventureGameLocationLinkRecs(&coresql.Options{
		Params: []coresql.Param{{Col: adventure_game_record.FieldAdventureGameLocationLinkGameID, Val: gameID}},
	})
	require.NoError(t, err)
	require.NotEmpty(t, links, "links should exist before removal")

	// Demo scenario may not create instances/subs (no AccountUserGameSubscriptionConfigs);
	// we still verify they are gone after removal below.
	_, err = dm.GetManyGameInstanceRecs(&coresql.Options{
		Params: []coresql.Param{{Col: game_record.FieldGameInstanceGameID, Val: gameID}},
	})
	require.NoError(t, err)
	_, err = dm.GetManyGameSubscriptionRecs(&coresql.Options{
		Params: []coresql.Param{{Col: game_record.FieldGameSubscriptionGameID, Val: gameID}},
	})
	require.NoError(t, err)

	// Remove the game and all dependents
	err = rnr.removeGameByName(demo_scenarios.DemoAdventureGameName)
	require.NoError(t, err)

	err = rnr.Domain.Commit()
	require.NoError(t, err)

	err = rnr.InitDomainTx()
	require.NoError(t, err)

	dm = rnr.Domain.(*domain.Domain)

	// Verify all dependents are gone
	locs, err = dm.GetManyAdventureGameLocationRecs(&coresql.Options{
		Params: []coresql.Param{{Col: adventure_game_record.FieldAdventureGameLocationGameID, Val: gameID}},
	})
	require.NoError(t, err)
	require.Empty(t, locs, "locations should be removed")

	items, err = dm.GetManyAdventureGameItemRecs(&coresql.Options{
		Params: []coresql.Param{{Col: adventure_game_record.FieldAdventureGameItemGameID, Val: gameID}},
	})
	require.NoError(t, err)
	require.Empty(t, items, "items should be removed")

	links, err = dm.GetManyAdventureGameLocationLinkRecs(&coresql.Options{
		Params: []coresql.Param{{Col: adventure_game_record.FieldAdventureGameLocationLinkGameID, Val: gameID}},
	})
	require.NoError(t, err)
	require.Empty(t, links, "links should be removed")

	instances, err := dm.GetManyGameInstanceRecs(&coresql.Options{
		Params: []coresql.Param{{Col: game_record.FieldGameInstanceGameID, Val: gameID}},
	})
	require.NoError(t, err)
	require.Empty(t, instances, "instances should be removed")

	subs, err := dm.GetManyGameSubscriptionRecs(&coresql.Options{
		Params: []coresql.Param{{Col: game_record.FieldGameSubscriptionGameID, Val: gameID}},
	})
	require.NoError(t, err)
	require.Empty(t, subs, "subscriptions should be removed")

	images, err := dm.GetManyGameImageRecs(&coresql.Options{
		Params: []coresql.Param{{Col: game_record.FieldGameImageGameID, Val: gameID}},
	})
	require.NoError(t, err)
	require.Empty(t, images, "images should be removed")

	creatures, err := dm.GetManyAdventureGameCreatureRecs(&coresql.Options{
		Params: []coresql.Param{{Col: adventure_game_record.FieldAdventureGameCreatureGameID, Val: gameID}},
	})
	require.NoError(t, err)
	require.Empty(t, creatures, "creatures should be removed")

	// Game record itself
	games, err := dm.GetManyGameRecs(&coresql.Options{
		Params: []coresql.Param{{Col: game_record.FieldGameName, Val: demo_scenarios.DemoAdventureGameName}},
	})
	require.NoError(t, err)
	require.Empty(t, games, "game record should be removed")
}

func TestLoadGameData(t *testing.T) {

	t.Run("missing game returns error", func(t *testing.T) {
		rnr := newTestRunnerWithDomain(t)
		ctx := newCliContext(nil, nil)
		err := rnr.loadGameData(ctx)
		require.Error(t, err)
		require.Contains(t, err.Error(), "--game is required")
	})

	t.Run("unknown game returns error", func(t *testing.T) {
		rnr := newTestRunnerWithDomain(t)
		ctx := newCliContext(map[string]string{"game": "Nonexistent Game"}, nil)
		err := rnr.loadGameData(ctx)
		require.Error(t, err)
		require.Contains(t, err.Error(), "unknown demo game")
	})

	t.Run("load draft", func(t *testing.T) {
		rnr := newTestRunnerWithDomain(t)
		defer cleanupLoadedGame(t, rnr)

		ctx := newCliContext(map[string]string{"game": demo_scenarios.DemoAdventureGameName}, nil)
		err := rnr.loadGameData(ctx)
		require.NoError(t, err, "loadGameData succeeds for draft")

		err = rnr.InitDomainTx()
		require.NoError(t, err)

		dm := rnr.Domain.(*domain.Domain)
		games, err := dm.GetManyGameRecs(&coresql.Options{
			Params: []coresql.Param{{Col: game_record.FieldGameName, Val: demo_scenarios.DemoAdventureGameName}},
		})
		require.NoError(t, err)
		require.Len(t, games, 1, "one game should exist")
		require.Equal(t, game_record.GameStatusDraft, games[0].Status, "game should be draft")
	})

	t.Run("case insensitive lookup", func(t *testing.T) {
		rnr := newTestRunnerWithDomain(t)
		defer cleanupLoadedGame(t, rnr)

		ctx := newCliContext(map[string]string{"game": "the door beneath the staircase"}, nil)
		err := rnr.loadGameData(ctx)
		require.NoError(t, err, "case-insensitive game name should work")
	})

	t.Run("duplicate without replace returns error", func(t *testing.T) {
		rnr := newTestRunnerWithDomain(t)
		defer cleanupLoadedGame(t, rnr)

		ctx := newCliContext(map[string]string{"game": demo_scenarios.DemoAdventureGameName}, nil)
		err := rnr.loadGameData(ctx)
		require.NoError(t, err, "first load should succeed")

		err = rnr.loadGameData(ctx)
		require.Error(t, err, "duplicate load should fail")
		require.Contains(t, err.Error(), "already exists")
		require.Contains(t, err.Error(), "--replace")
	})

	t.Run("load with publish", func(t *testing.T) {
		rnr := newTestRunnerWithDomain(t)
		defer cleanupLoadedGame(t, rnr)

		ctx := newCliContext(
			map[string]string{"game": demo_scenarios.DemoAdventureGameName},
			map[string]bool{"publish": true},
		)
		err := rnr.loadGameData(ctx)
		require.NoError(t, err, "loadGameData succeeds with publish")

		err = rnr.InitDomainTx()
		require.NoError(t, err)

		dm := rnr.Domain.(*domain.Domain)
		games, err := dm.GetManyGameRecs(&coresql.Options{
			Params: []coresql.Param{{Col: game_record.FieldGameName, Val: demo_scenarios.DemoAdventureGameName}},
		})
		require.NoError(t, err)
		require.Len(t, games, 1, "one game should exist")
		require.Equal(t, game_record.GameStatusPublished, games[0].Status, "game should be published")
	})

	t.Run("replace", func(t *testing.T) {
		rnr := newTestRunnerWithDomain(t)
		defer cleanupLoadedGame(t, rnr)

		ctx := newCliContext(map[string]string{"game": demo_scenarios.DemoAdventureGameName}, nil)
		err := rnr.loadGameData(ctx)
		require.NoError(t, err)

		err = rnr.InitDomainTx()
		require.NoError(t, err)

		dm := rnr.Domain.(*domain.Domain)
		games, err := dm.GetManyGameRecs(&coresql.Options{
			Params: []coresql.Param{{Col: game_record.FieldGameName, Val: demo_scenarios.DemoAdventureGameName}},
		})
		require.NoError(t, err)
		require.Len(t, games, 1)
		originalID := games[0].ID

		replaceCtx := newCliContext(
			map[string]string{"game": demo_scenarios.DemoAdventureGameName},
			map[string]bool{"replace": true},
		)
		err = rnr.loadGameData(replaceCtx)
		require.NoError(t, err, "loadGameData with replace succeeds")

		err = rnr.InitDomainTx()
		require.NoError(t, err)

		dm = rnr.Domain.(*domain.Domain)
		games, err = dm.GetManyGameRecs(&coresql.Options{
			Params: []coresql.Param{{Col: game_record.FieldGameName, Val: demo_scenarios.DemoAdventureGameName}},
		})
		require.NoError(t, err)
		require.Len(t, games, 1, "exactly one game should exist after replace")
		require.NotEqual(t, originalID, games[0].ID, "game ID should differ after replace")
	})
}

func TestListDemoGames(t *testing.T) {
	games := ListDemoGames()
	require.NotEmpty(t, games, "ListDemoGames should return at least one entry")

	var found *DemoGameSummary
	for i := range games {
		if games[i].Name == demo_scenarios.DemoAdventureGameName {
			found = &games[i]
			break
		}
	}
	require.NotNil(t, found, "should contain %q", demo_scenarios.DemoAdventureGameName)
	require.Equal(t, game_record.GameTypeAdventure, found.GameType)
	require.NotEmpty(t, found.Description)

	for i := 1; i < len(games); i++ {
		require.True(t, games[i-1].Name <= games[i].Name, "games should be sorted by name")
	}
}

func TestLookupDemoGame(t *testing.T) {
	_, ok := LookupDemoGame(demo_scenarios.DemoAdventureGameName)
	require.True(t, ok, "exact name should match")

	_, ok = LookupDemoGame("the door beneath the staircase")
	require.True(t, ok, "lowercase should match")

	_, ok = LookupDemoGame("THE DOOR BENEATH THE STAIRCASE")
	require.True(t, ok, "uppercase should match")

	_, ok = LookupDemoGame("Nonexistent Game")
	require.False(t, ok, "unknown name should not match")
}

// --- internal helpers ---

// newTestHarness creates a harness using the runner's existing store and config.
func newTestHarness(t *testing.T, rnr *Runner, config harness.DataConfig) (*harness.Testing, error) {
	t.Helper()
	return harness.NewTesting(rnr.Config, rnr.Log, rnr.Store, rnr.JobClient, rnr.Scanner, config)
}

// cleanupLoadedGame removes demo game data and accounts so the test is idempotent.
func cleanupLoadedGame(t *testing.T, rnr *Runner) {
	t.Helper()

	// Re-init domain if needed
	if rnr.Domain == nil {
		err := rnr.InitDomain()
		if err != nil {
			return
		}
	} else {
		_ = rnr.InitDomainTx()
	}

	_ = rnr.removeGameByName(demo_scenarios.DemoAdventureGameName)
	_ = rnr.Domain.Commit()
	_ = rnr.InitDomainTx()

	removeDemoAccounts(t, rnr)
}
