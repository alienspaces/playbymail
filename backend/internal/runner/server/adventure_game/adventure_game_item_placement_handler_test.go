package adventure_game_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/playbymail/core/server"
	"gitlab.com/alienspaces/playbymail/internal/harness"
	"gitlab.com/alienspaces/playbymail/internal/runner/server/adventure_game"
	"gitlab.com/alienspaces/playbymail/internal/utils/deps"
	"gitlab.com/alienspaces/playbymail/internal/utils/testutil"
	"gitlab.com/alienspaces/playbymail/schema/api"
)

func Test_adventureGameItemPlacementHandler(t *testing.T) {
	t.Parallel()

	th := deps.NewHandlerTestHarness(t)
	require.NotNil(t, th, "NewTestHarness returns without error")

	_, err := th.Setup()
	require.NoError(t, err, "Test data setup returns without error")
	defer func() {
		err = th.Teardown()
		require.NoError(t, err, "Test data teardown returns without error")
	}()

	gameRec, err := th.Data.GetGameRecByRef(harness.GameOneRef)
	require.NoError(t, err, "GetGameRecByRef returns without error")

	itemRec, err := th.Data.GetGameItemRecByRef(harness.GameItemOneRef)
	require.NoError(t, err, "GetGameItemRecByRef returns without error")

	locationRec, err := th.Data.GetGameLocationRecByRef(harness.GameLocationOneRef)
	require.NoError(t, err, "GetGameLocationRecByRef returns without error")

	testCaseCollectionResponseDecoder := testutil.TestCaseResponseDecoderGeneric[api.AdventureGameItemPlacementCollectionResponse]
	testCaseResponseDecoder := testutil.TestCaseResponseDecoderGeneric[api.AdventureGameItemPlacementResponse]

	testCases := []struct {
		testutil.TestCase
	}{
		{
			TestCase: testutil.TestCase{
				Name: "API key with open access \\ get many item placements \\ returns expected placements",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[adventure_game.SearchManyAdventureGameItemPlacements]
				},
				RequestQueryParams: func(d harness.Data) map[string]any {
					return map[string]any{
						"page_size":   10,
						"page_number": 1,
					}
				},
				ResponseDecoder: testCaseCollectionResponseDecoder,
				ResponseCode:    http.StatusOK,
			},
		},
		{
			TestCase: testutil.TestCase{
				Name: "API key with open access \\ create item placement with valid properties \\ returns created placement",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[adventure_game.CreateOneAdventureGameItemPlacement]
				},
				RequestBody: func(d harness.Data) any {
					return api.AdventureGameItemPlacementRequest{
						AdventureGameItemID:     itemRec.ID,
						AdventureGameLocationID: locationRec.ID,
						InitialCount:            5,
					}
				},
				RequestPathParams: func(d harness.Data) map[string]string {
					return map[string]string{":game_id": gameRec.ID}
				},
				ResponseDecoder: testCaseResponseDecoder,
				ResponseCode:    http.StatusCreated,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.TestName(), func(t *testing.T) {
			testutil.RunTestCase(t, th, &tc.TestCase, nil)
		})
	}
}
