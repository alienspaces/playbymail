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
			Name: "public request returns catalog with subscription entries that have available instances",
			HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
				return rnr.GetHandlerConfig()[catalog.GetCatalogGames]
			},
			ResponseDecoder: testutil.TestCaseResponseDecoderGeneric[catalog_schema.CatalogCollectionResponse],
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

				res := body.(catalog_schema.CatalogCollectionResponse)
				require.NotNil(t, res.Data, "Response data is not nil")

				require.NotEmpty(t, res.Data, "Catalog contains at least one subscription entry")

				for _, entry := range res.Data {
					require.NotEmpty(t, entry.GameSubscriptionID, "GameSubscriptionID is not empty")
					require.NotEmpty(t, entry.GameName, "GameName is not empty")
					require.NotEmpty(t, entry.GameType, "GameType is not empty")
					require.Greater(t, entry.TurnDurationHours, 0, "TurnDurationHours is greater than 0")
				}
			}
			testutil.RunTestCase(t, th, &tc, testFunc)
		})
	}
}
