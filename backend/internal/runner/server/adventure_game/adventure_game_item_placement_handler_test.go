package adventure_game

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/playbymail/core/server"
	"gitlab.com/alienspaces/playbymail/internal/harness"
	"gitlab.com/alienspaces/playbymail/internal/utils/deps"
	"gitlab.com/alienspaces/playbymail/internal/utils/testutil"
	"gitlab.com/alienspaces/playbymail/schema"
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

	testCaseCollectionResponseDecoder := testutil.TestCaseResponseDecoderGeneric[schema.AdventureGameItemPlacementCollectionResponse]
	testCaseResponseDecoder := testutil.TestCaseResponseDecoderGeneric[schema.AdventureGameItemPlacementResponse]

	testCases := []struct {
		testutil.TestCase
	}{
		{
			TestCase: testutil.TestCase{
				Name: "API key with open access \\ get many item placements \\ returns expected placements",
				HandlerConfig: func(rnr *server.Runner) server.HandlerConfig {
					return rnr.HandlerConfig[searchManyAdventureGameItemPlacements]
				},
				RequestQueryParams: func(d harness.Data) map[string]any {
					return map[string]any{
						"page_size":   10,
						"page_number": 1,
					}
				},
				RequestPathParams: func(d harness.Data) map[string]string {
					return map[string]string{":game_id": gameRec.ID}
				},
				ResponseDecoder: testCaseCollectionResponseDecoder,
				ResponseCode:    http.StatusOK,
			},
		},
		{
			TestCase: testutil.TestCase{
				Name: "API key with open access \\ create item placement with valid properties \\ returns created placement",
				HandlerConfig: func(rnr *server.Runner) server.HandlerConfig {
					return rnr.HandlerConfig[createOneAdventureGameItemPlacement]
				},
				RequestBody: func(d harness.Data) any {
					return schema.AdventureGameItemPlacementRequest{
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
