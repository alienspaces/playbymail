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

func Test_getAdventureGameLocationInstancesHandler(t *testing.T) {
	t.Parallel()

	th := deps.NewHandlerTestHarness(t)
	require.NotNil(t, th, "newTestHarness returns without error")

	_, err := th.Setup()
	require.NoError(t, err, "Test data setup returns without error")
	defer func() {
		err = th.Teardown()
		require.NoError(t, err, "Test data teardown returns without error")
	}()

	type testCase struct {
		testutil.TestCase
		collectionRequest     bool
		collectionRecordCount int
	}

	testCaseCollectionResponseDecoder := testutil.TestCaseResponseDecoderGeneric[adventure_game_schema.AdventureGameLocationInstanceCollectionResponse]
	testCaseResponseDecoder := testutil.TestCaseResponseDecoderGeneric[adventure_game_schema.AdventureGameLocationInstanceResponse]

	// Setup: get a location instance for reference
	locationInstanceRec, err := th.Data.GetAdventureGameLocationInstanceRecByGameInstanceAndLocationRef(harness.GameInstanceOneRef, harness.GameLocationOneRef)
	require.NoError(t, err, "GetAdventureGameLocationInstanceRecByGameInstanceAndLocationRef returns without error")
	require.NotNil(t, locationInstanceRec, "Location instance exists for reference")

	// Get the game instance to get the game_id
	gameInstanceRec, err := th.Data.GetGameInstanceRecByID(locationInstanceRec.GameInstanceID)
	require.NoError(t, err, "GetGameInstanceRecByID returns without error")
	require.NotNil(t, gameInstanceRec, "Game instance exists for reference")

	t.Logf("Test filter game_location_id: %s", locationInstanceRec.AdventureGameLocationID)

	testCases := []testCase{
		{
			TestCase: testutil.TestCase{
				Name: "API key with open access \\ get many location instances \\ returns expected instances",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[adventure_game.GetManyAdventureGameLocationInstances]
				},
				RequestPathParams: func(d harness.Data) map[string]string {
					return map[string]string{
						":game_instance_id": gameInstanceRec.ID,
					}
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
			collectionRequest:     true,
			collectionRecordCount: 2,
		},
		{
			TestCase: testutil.TestCase{
				Name: "API key with open access \\ get one location instance with valid ID \\ returns expected instance",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[adventure_game.GetOneAdventureGameLocationInstance]
				},
				RequestPathParams: func(d harness.Data) map[string]string {
					return map[string]string{
						":game_instance_id":     gameInstanceRec.ID,
						":location_instance_id": locationInstanceRec.ID,
					}
				},
				ResponseDecoder: testCaseResponseDecoder,
				ResponseCode:    http.StatusOK,
			},
			collectionRequest:     false,
			collectionRecordCount: 0,
		},
	}

	for _, testCase := range testCases {
		t.Logf("Running test >%s<", testCase.Name)

		t.Run(testCase.Name, func(t *testing.T) {
			testFunc := func(method string, body interface{}) {
				if testCase.TestResponseCode() != http.StatusOK {
					return
				}

				require.NotNil(t, body, "Response body is not nil")

				var responses []*adventure_game_schema.AdventureGameLocationInstanceResponseData
				if testCase.collectionRequest {
					responses = body.(adventure_game_schema.AdventureGameLocationInstanceCollectionResponse).Data
				} else {
					responses = append(responses, body.(adventure_game_schema.AdventureGameLocationInstanceResponse).Data)
				}

				if testCase.collectionRequest {
					for i, d := range responses {
						t.Logf("Record %d: ID=%s, GameInstanceID=%s, AdventureGameLocationID=%s", i, d.ID, d.GameInstanceID, d.AdventureGameLocationID)
					}
					require.Equal(t, testCase.collectionRecordCount, len(responses), "Response record count length equals expected")
				}

				if testCase.collectionRequest && testCase.collectionRecordCount == 0 {
					require.Empty(t, responses, "Response body should be empty")
				} else {
					require.NotEmpty(t, responses, "Response body is not empty")
				}

				for _, d := range responses {
					require.NotEmpty(t, d.ID, "LocationInstance ID is not empty")
					require.NotEmpty(t, d.GameInstanceID, "LocationInstance GameInstanceID is not empty")
					require.NotEmpty(t, d.AdventureGameLocationID, "LocationInstance AdventureGameLocationID is not empty")
				}
			}

			testutil.RunTestCase(t, th, &testCase, testFunc)
		})
	}
}
