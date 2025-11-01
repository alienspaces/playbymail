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
	"gitlab.com/alienspaces/playbymail/schema/api/adventure_game_schema"
)

func Test_adventureGameCreaturePlacementHandler(t *testing.T) {
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

	creatureRec, err := th.Data.GetAdventureGameCreatureRecByRef(harness.GameCreatureOneRef)
	require.NoError(t, err, "GetGameCreatureRecByRef returns without error")

	locationRec, err := th.Data.GetAdventureGameLocationRecByRef(harness.GameLocationOneRef)
	require.NoError(t, err, "GetGameLocationRecByRef returns without error")

	testCaseCollectionResponseDecoder := testutil.TestCaseResponseDecoderGeneric[adventure_game_schema.AdventureGameCreaturePlacementCollectionResponse]
	testCaseResponseDecoder := testutil.TestCaseResponseDecoderGeneric[adventure_game_schema.AdventureGameCreaturePlacementResponse]

	testCases := []struct {
		testutil.TestCase
	}{
		{
			TestCase: testutil.TestCase{
				Name: "API key with open access \\ get many creature placements \\ returns expected placements",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[adventure_game.SearchManyAdventureGameCreaturePlacements]
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
				Name: "API key with open access \\ create creature placement with valid properties \\ returns created placement",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[adventure_game.CreateOneAdventureGameCreaturePlacement]
				},
				RequestBody: func(d harness.Data) any {
					return adventure_game_schema.AdventureGameCreaturePlacementRequest{
						AdventureGameCreatureID: creatureRec.ID,
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
