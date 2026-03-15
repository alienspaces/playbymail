package player_test

import (
	"net/http"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"
	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/playbymail/core/server"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/core/type/storer"
	"gitlab.com/alienspaces/playbymail/internal/harness"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
	"gitlab.com/alienspaces/playbymail/internal/runner/server/player"
	"gitlab.com/alienspaces/playbymail/internal/turnsheet"
	"gitlab.com/alienspaces/playbymail/internal/utils/config"
	"gitlab.com/alienspaces/playbymail/internal/utils/testutil"
)

// gameSubscriptionInstanceIDForPlayer returns the game_subscription_instance ID for the pro-player's subscription
// to game-instance-one, which the harness always creates.
func gameSubscriptionInstanceIDForPlayer(t *testing.T, data harness.Data) string {
	t.Helper()
	return gameSubscriptionInstanceIDForSubscriptionRef(t, data, harness.GameSubscriptionPlayerOneRef)
}

// gameSubscriptionInstanceIDForStandardPlayer returns the game_subscription_instance ID for the standard account's
// player subscription to game-instance-one (has the single character/turn sheet in default harness).
func gameSubscriptionInstanceIDForStandardPlayer(t *testing.T, data harness.Data) string {
	t.Helper()
	return gameSubscriptionInstanceIDForSubscriptionRef(t, data, harness.GameSubscriptionPlayerThreeRef)
}

func gameSubscriptionInstanceIDForSubscriptionRef(t *testing.T, data harness.Data, subscriptionRef string) string {
	t.Helper()
	playerSubRec, err := data.GetGameSubscriptionRecByRef(subscriptionRef)
	require.NoError(t, err)
	for _, gsi := range data.GameSubscriptionInstanceRecs {
		if gsi.GameSubscriptionID == playerSubRec.ID {
			return gsi.ID
		}
	}
	t.Fatalf("no game_subscription_instance found for player subscription %s", playerSubRec.ID)
	return ""
}

// turnSheetIDForAccountRef returns a non-completed turn sheet ID for the given account ref's subscription (same account/game).
// If sheetType is non-empty, returns a turn sheet of that type.
func turnSheetIDForAccountRef(t *testing.T, data harness.Data, accountRef string, subscriptionRef string, sheetType string) string {
	t.Helper()
	playerSubRec, err := data.GetGameSubscriptionRecByRef(subscriptionRef)
	require.NoError(t, err)
	accountUserRec, err := data.GetAccountUserRecByRef(accountRef)
	require.NoError(t, err)
	for _, ts := range data.GameTurnSheetRecs {
		if ts.AccountUserID == accountUserRec.ID && ts.AccountID == playerSubRec.AccountID && ts.GameID == playerSubRec.GameID && !ts.IsCompleted {
			if sheetType == "" || ts.SheetType == sheetType {
				return ts.ID
			}
		}
	}
	if sheetType != "" {
		t.Fatalf("no incomplete turn sheet of type %q found for account ref %s", sheetType, accountRef)
	}
	t.Fatalf("no incomplete turn sheet found for account ref %s", accountRef)
	return ""
}

// turnSheetIDForStandardPlayer returns a non-completed turn sheet ID for the standard account's player subscription (default harness has one character).
func turnSheetIDForStandardPlayer(t *testing.T, data harness.Data, sheetType string) string {
	t.Helper()
	return turnSheetIDForAccountRef(t, data, harness.AccountUserStandardRef, harness.GameSubscriptionPlayerThreeRef, sheetType)
}

func Test_getGameSubscriptionInstanceTurnSheetListHandler(t *testing.T) {
	t.Parallel()

	th := testutil.NewTestHarness(t)
	require.NotNil(t, th)

	_, err := th.Setup()
	require.NoError(t, err)
	defer func() {
		err = th.Teardown()
		require.NoError(t, err)
	}()

	testCases := []struct {
		testutil.TestCase
	}{
		{
			TestCase: testutil.TestCase{
				Name: "authenticated standard player gets turn sheet list for their gsi",
				NewRunner: func(cfg config.Config, l logger.Logger, s storer.Storer, j *river.Client[pgx.Tx], scanner turnsheet.TurnSheetScanner, d harness.Data) (testutil.TestRunnerer, error) {
					return testutil.NewTestRunner(cfg, l, s, j, scanner)
				},
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[player.GetGameSubscriptionInstanceTurnSheetList]
				},
				RequestHeaders: testutil.AuthHeaderStandard,
				RequestPathParams: func(d harness.Data) map[string]string {
					return map[string]string{
						":game_subscription_instance_id": gameSubscriptionInstanceIDForStandardPlayer(t, d),
					}
				},
				ResponseCode: http.StatusOK,
			},
		},
		{
			TestCase: testutil.TestCase{
				Name: "unauthenticated request returns unauthorized",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[player.GetGameSubscriptionInstanceTurnSheetList]
				},
				RequestPathParams: func(d harness.Data) map[string]string {
					return map[string]string{
						":game_subscription_instance_id": gameSubscriptionInstanceIDForStandardPlayer(t, d),
					}
				},
				ResponseCode: http.StatusUnauthorized,
			},
		},
	}

	for _, testCase := range testCases {
		t.Logf("Running test >%s<\n", testCase.Name)
		t.Run(testCase.Name, func(t *testing.T) {
			testutil.RunTestCase(t, th, &testCase.TestCase, nil)
		})
	}
}

func Test_downloadGameSubscriptionInstanceTurnSheetPDFHandler(t *testing.T) {
	t.Parallel()

	th := testutil.NewTestHarness(t)
	require.NotNil(t, th)

	_, err := th.Setup()
	require.NoError(t, err)
	defer func() {
		err = th.Teardown()
		require.NoError(t, err)
	}()

	testCases := []struct {
		testutil.TestCase
	}{
		{
			TestCase: testutil.TestCase{
				Name: "unauthenticated request returns unauthorized",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[player.DownloadGameSubscriptionInstanceTurnSheetPDF]
				},
				RequestPathParams: func(d harness.Data) map[string]string {
					return map[string]string{
						":game_subscription_instance_id": gameSubscriptionInstanceIDForPlayer(t, d),
						":game_turn_sheet_id":            "00000000-0000-0000-0000-000000000000",
					}
				},
				ResponseCode: http.StatusUnauthorized,
			},
		},
	}

	for _, testCase := range testCases {
		t.Logf("Running test >%s<\n", testCase.Name)
		t.Run(testCase.Name, func(t *testing.T) {
			testutil.RunTestCase(t, th, &testCase.TestCase, nil)
		})
	}
}

func Test_saveGameSubscriptionInstanceTurnSheetHandler(t *testing.T) {
	t.Parallel()

	th := testutil.NewTestHarness(t)
	require.NotNil(t, th)

	_, err := th.Setup()
	require.NoError(t, err)
	defer func() {
		err = th.Teardown()
		require.NoError(t, err)
	}()

	testCases := []struct {
		testutil.TestCase
	}{
		{
			TestCase: testutil.TestCase{
				Name: "authenticated player saves location_choice turn sheet with valid scanned_data then returns 200",
				NewRunner: func(cfg config.Config, l logger.Logger, s storer.Storer, j *river.Client[pgx.Tx], scanner turnsheet.TurnSheetScanner, d harness.Data) (testutil.TestRunnerer, error) {
					return testutil.NewTestRunner(cfg, l, s, j, scanner)
				},
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[player.SaveGameSubscriptionInstanceTurnSheet]
				},
				RequestHeaders: testutil.AuthHeaderStandard,
				RequestPathParams: func(d harness.Data) map[string]string {
					return map[string]string{
						":game_subscription_instance_id": gameSubscriptionInstanceIDForStandardPlayer(t, d),
						":game_turn_sheet_id":            turnSheetIDForStandardPlayer(t, d, adventure_game_record.AdventureGameTurnSheetTypeLocationChoice),
					}
				},
				RequestBody: func(d harness.Data) any {
					return map[string]any{
						"scanned_data": map[string]any{"location_choice": "loc-any"},
					}
				},
				ResponseCode: http.StatusOK,
			},
		},
		{
			TestCase: testutil.TestCase{
				Name: "authenticated player saves inventory_management turn sheet with valid scanned_data then returns 200",
				NewRunner: func(cfg config.Config, l logger.Logger, s storer.Storer, j *river.Client[pgx.Tx], scanner turnsheet.TurnSheetScanner, d harness.Data) (testutil.TestRunnerer, error) {
					return testutil.NewTestRunner(cfg, l, s, j, scanner)
				},
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[player.SaveGameSubscriptionInstanceTurnSheet]
				},
				RequestHeaders: testutil.AuthHeaderStandard,
				RequestPathParams: func(d harness.Data) map[string]string {
					return map[string]string{
						":game_subscription_instance_id": gameSubscriptionInstanceIDForStandardPlayer(t, d),
						":game_turn_sheet_id":            turnSheetIDForStandardPlayer(t, d, adventure_game_record.AdventureGameTurnSheetTypeInventoryManagement),
					}
				},
				RequestBody: func(d harness.Data) any {
					return map[string]any{
						"scanned_data": map[string]any{
							"pick_up": []string{},
							"drop":    []string{},
							"equip":   []string{"item-1"},
							"unequip": []string{},
						},
					}
				},
				ResponseCode: http.StatusOK,
			},
		},
		{
			TestCase: testutil.TestCase{
				Name: "authenticated player saves location_choice with invalid scanned_data then returns 400",
				NewRunner: func(cfg config.Config, l logger.Logger, s storer.Storer, j *river.Client[pgx.Tx], scanner turnsheet.TurnSheetScanner, d harness.Data) (testutil.TestRunnerer, error) {
					return testutil.NewTestRunner(cfg, l, s, j, scanner)
				},
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[player.SaveGameSubscriptionInstanceTurnSheet]
				},
				RequestHeaders: testutil.AuthHeaderStandard,
				RequestPathParams: func(d harness.Data) map[string]string {
					return map[string]string{
						":game_subscription_instance_id": gameSubscriptionInstanceIDForStandardPlayer(t, d),
						":game_turn_sheet_id":            turnSheetIDForStandardPlayer(t, d, adventure_game_record.AdventureGameTurnSheetTypeLocationChoice),
					}
				},
				RequestBody: func(d harness.Data) any {
					// location_choice must be string; array is invalid per schema
					return map[string]any{
						"scanned_data": map[string]any{"location_choice": []string{"a"}},
					}
				},
				ResponseCode: http.StatusBadRequest,
			},
		},
		{
			TestCase: testutil.TestCase{
				Name: "unauthenticated request returns unauthorized",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[player.SaveGameSubscriptionInstanceTurnSheet]
				},
				RequestPathParams: func(d harness.Data) map[string]string {
					return map[string]string{
						":game_subscription_instance_id": gameSubscriptionInstanceIDForStandardPlayer(t, d),
						":game_turn_sheet_id":            turnSheetIDForStandardPlayer(t, d, ""),
					}
				},
				RequestBody: func(d harness.Data) any {
					return map[string]any{"scanned_data": map[string]any{}}
				},
				ResponseCode: http.StatusUnauthorized,
			},
		},
	}

	for _, testCase := range testCases {
		t.Logf("Running test >%s<\n", testCase.Name)
		t.Run(testCase.Name, func(t *testing.T) {
			testutil.RunTestCase(t, th, &testCase.TestCase, nil)
		})
	}
}

func Test_submitGameSubscriptionInstanceTurnSheetsHandler(t *testing.T) {
	t.Parallel()

	th := testutil.NewTestHarness(t)
	require.NotNil(t, th)

	_, err := th.Setup()
	require.NoError(t, err)
	defer func() {
		err = th.Teardown()
		require.NoError(t, err)
	}()

	testCases := []struct {
		testutil.TestCase
	}{
		{
			TestCase: testutil.TestCase{
				Name: "authenticated standard player submits turn sheets for their gsi",
				NewRunner: func(cfg config.Config, l logger.Logger, s storer.Storer, j *river.Client[pgx.Tx], scanner turnsheet.TurnSheetScanner, d harness.Data) (testutil.TestRunnerer, error) {
					return testutil.NewTestRunner(cfg, l, s, j, scanner)
				},
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[player.SubmitGameSubscriptionInstanceTurnSheets]
				},
				RequestHeaders: testutil.AuthHeaderStandard,
				RequestPathParams: func(d harness.Data) map[string]string {
					return map[string]string{
						":game_subscription_instance_id": gameSubscriptionInstanceIDForStandardPlayer(t, d),
					}
				},
				ResponseCode: http.StatusOK,
			},
		},
		{
			TestCase: testutil.TestCase{
				Name: "unauthenticated request returns unauthorized",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[player.SubmitGameSubscriptionInstanceTurnSheets]
				},
				RequestPathParams: func(d harness.Data) map[string]string {
					return map[string]string{
						":game_subscription_instance_id": gameSubscriptionInstanceIDForStandardPlayer(t, d),
					}
				},
				ResponseCode: http.StatusUnauthorized,
			},
		},
	}

	for _, testCase := range testCases {
		t.Logf("Running test >%s<\n", testCase.Name)
		t.Run(testCase.Name, func(t *testing.T) {
			testutil.RunTestCase(t, th, &testCase.TestCase, nil)
		})
	}
}
