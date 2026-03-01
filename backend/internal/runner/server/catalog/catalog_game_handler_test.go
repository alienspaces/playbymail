package catalog_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/playbymail/core/server"
	"gitlab.com/alienspaces/playbymail/internal/runner/server/catalog"
	"gitlab.com/alienspaces/playbymail/internal/utils/testutil"
	"gitlab.com/alienspaces/playbymail/schema/api/catalog_schema"
)

func Test_getCatalogGamesHandler(t *testing.T) {
	t.Parallel()

	th := testutil.NewTestHarness(t)
	require.NotNil(t, th, "NewTestHarness returns without error")

	_, err := th.Setup()
	require.NoError(t, err, "Test data setup returns without error")
	defer func() {
		err = th.Teardown()
		require.NoError(t, err, "Test data teardown returns without error")
	}()

	testCases := []testutil.TestCase{
		{
			Name: "public request returns catalog with games that have available instances",
			HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
				return rnr.GetHandlerConfig()[catalog.GetCatalogGames]
			},
			// No RequestHeaders — public endpoint requires no authentication
			ResponseDecoder: testutil.TestCaseResponseDecoderGeneric[catalog_schema.CatalogGameCollectionResponse],
			ResponseCode:    http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Logf("Running test >%s<", tc.Name)

		t.Run(tc.Name, func(t *testing.T) {
			testFunc := func(method string, body any) {
				if tc.TestResponseCode() != http.StatusOK {
					return
				}

				require.NotNil(t, body, "Response body is not nil")

				res := body.(catalog_schema.CatalogGameCollectionResponse)
				require.NotNil(t, res.Data, "Response data is not nil")

				// The default harness creates GameOneRef with an active manager subscription
				// (ProManagerAccountRef / GameSubscriptionManagerOneRef) and game instances
				// in "created" status -- so at least one game should appear in the catalog.
				require.NotEmpty(t, res.Data, "Catalog contains at least one game")

				// Verify the structure of each catalog entry -- do not assert specific IDs
				// because the full test suite runs tests in parallel and other harnesses may
				// also commit games that are visible to this handler's transaction.
				for _, game := range res.Data {
					require.NotEmpty(t, game.ID, "Game ID is not empty")
					require.NotEmpty(t, game.Name, "Game name is not empty")
					require.NotEmpty(t, game.GameType, "Game type is not empty")
					require.Greater(t, game.TurnDurationHours, 0, "Turn duration hours is greater than 0")
					require.NotNil(t, game.AvailableInstances, "Available instances is not nil")
					require.NotEmpty(t, game.AvailableInstances, "Game has at least one available instance")

					for _, inst := range game.AvailableInstances {
						require.NotEmpty(t, inst.ID, "Instance ID is not empty")
					}
				}
			}
			testutil.RunTestCase(t, th, &tc, testFunc)
		})
	}
}
