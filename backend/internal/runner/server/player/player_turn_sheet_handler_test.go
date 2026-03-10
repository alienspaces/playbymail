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
	"gitlab.com/alienspaces/playbymail/internal/runner/server/player"
	"gitlab.com/alienspaces/playbymail/internal/turnsheet"
	"gitlab.com/alienspaces/playbymail/internal/utils/config"
	"gitlab.com/alienspaces/playbymail/internal/utils/testutil"
)

// gsiIDForPlayer returns the game_subscription_instance ID for the pro-player's subscription
// to game-instance-one, which the harness always creates.
func gsiIDForPlayer(t *testing.T, data harness.Data) string {
	t.Helper()

	playerSubRec, err := data.GetGameSubscriptionRecByRef(harness.GameSubscriptionPlayerOneRef)
	require.NoError(t, err)

	for _, gsi := range data.GameSubscriptionInstanceRecs {
		if gsi.GameSubscriptionID == playerSubRec.ID {
			return gsi.ID
		}
	}

	t.Fatalf("no game_subscription_instance found for player subscription %s", playerSubRec.ID)
	return ""
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
				Name: "authenticated pro player gets turn sheet list for their gsi",
				NewRunner: func(cfg config.Config, l logger.Logger, s storer.Storer, j *river.Client[pgx.Tx], scanner turnsheet.TurnSheetScanner, d harness.Data) (testutil.TestRunnerer, error) {
					return testutil.NewTestRunner(cfg, l, s, j, scanner)
				},
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[player.GetGameSubscriptionInstanceTurnSheetList]
				},
				RequestHeaders: testutil.AuthHeaderProPlayer,
				RequestPathParams: func(d harness.Data) map[string]string {
					return map[string]string{
						":game_subscription_instance_id": gsiIDForPlayer(t, d),
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
						":game_subscription_instance_id": gsiIDForPlayer(t, d),
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
						":game_subscription_instance_id": gsiIDForPlayer(t, d),
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
				Name: "authenticated pro player submits turn sheets for their gsi",
				NewRunner: func(cfg config.Config, l logger.Logger, s storer.Storer, j *river.Client[pgx.Tx], scanner turnsheet.TurnSheetScanner, d harness.Data) (testutil.TestRunnerer, error) {
					return testutil.NewTestRunner(cfg, l, s, j, scanner)
				},
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[player.SubmitGameSubscriptionInstanceTurnSheets]
				},
				RequestHeaders: testutil.AuthHeaderProPlayer,
				RequestPathParams: func(d harness.Data) map[string]string {
					return map[string]string{
						":game_subscription_instance_id": gsiIDForPlayer(t, d),
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
						":game_subscription_instance_id": gsiIDForPlayer(t, d),
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
