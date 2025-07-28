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

func Test_getGameLocationInstanceHandler(t *testing.T) {
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

	testCaseCollectionResponseDecoder := testutil.TestCaseResponseDecoderGeneric[api.AdventureGameLocationInstanceCollectionResponse]
	testCaseResponseDecoder := testutil.TestCaseResponseDecoderGeneric[api.AdventureGameLocationInstanceResponse]

	// Setup: get a location instance for reference
	locationInstanceRec, err := th.Data.GetGameLocationInstanceRecByRef(harness.GameLocationInstanceOneRef)
	require.NoError(t, err, "GetGameLocationInstanceRecByLocationRef returns without error")
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

				var responses []*api.AdventureGameLocationInstanceResponseData
				if testCase.collectionRequest {
					responses = body.(api.AdventureGameLocationInstanceCollectionResponse).Data
				} else {
					responses = append(responses, body.(api.AdventureGameLocationInstanceResponse).Data)
				}

				if testCase.collectionRequest {
					for i, d := range responses {
						t.Logf("Record %d: ID=%s, GameInstanceID=%s, GameLocationID=%s", i, d.ID, d.GameInstanceID, d.GameLocationID)
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
					require.NotEmpty(t, d.GameLocationID, "LocationInstance GameLocationID is not empty")
				}
			}

			testutil.RunTestCase(t, th, &testCase, testFunc)
		})
	}
}

func Test_createUpdateDeleteGameLocationInstanceHandler(t *testing.T) {
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
		expectResponse func(d harness.Data, req api.AdventureGameLocationInstanceRequest) api.AdventureGameLocationInstanceResponse
	}

	testCaseResponseDecoder := testutil.TestCaseResponseDecoderGeneric[api.AdventureGameLocationInstanceResponse]

	gameInstanceRec, err := th.Data.GetGameInstanceRecByRef(harness.GameInstanceOneRef)
	require.NoError(t, err, "GetGameInstanceRecByRef returns without error")
	require.NotNil(t, gameInstanceRec, "Game instance exists for reference")

	locationRec, err := th.Data.GetGameLocationRecByRef(harness.GameLocationOneRef)
	require.NoError(t, err, "GetGameLocationRecByRef returns without error")

	locationInstanceRec, err := th.Data.GetGameLocationInstanceRecByRef(harness.GameLocationInstanceOneRef)
	require.NoError(t, err, "GetGameLocationInstanceRecByRef returns without error")
	require.NotNil(t, locationInstanceRec, "Location instance exists for reference")
	gameLocationInstance := locationInstanceRec

	testCases := []testCase{
		{
			TestCase: testutil.TestCase{
				Name: "API key with open access \\ create location instance with valid properties \\ returns created instance",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[adventure_game.CreateOneAdventureGameLocationInstance]
				},
				RequestPathParams: func(d harness.Data) map[string]string {
					return map[string]string{
						":game_instance_id": gameInstanceRec.ID,
					}
				},
				RequestBody: func(d harness.Data) any {
					return api.AdventureGameLocationInstanceRequest{
						GameID:         gameInstanceRec.GameID,
						GameInstanceID: gameInstanceRec.ID,
						GameLocationID: locationRec.ID,
					}
				},
				ResponseDecoder: testCaseResponseDecoder,
				ResponseCode:    http.StatusCreated,
			},
			expectResponse: func(d harness.Data, req api.AdventureGameLocationInstanceRequest) api.AdventureGameLocationInstanceResponse {
				return api.AdventureGameLocationInstanceResponse{
					Data: &api.AdventureGameLocationInstanceResponseData{
						GameInstanceID: req.GameInstanceID,
						GameLocationID: req.GameLocationID,
						GameID:         gameInstanceRec.GameID,
					},
				}
			},
		},
		{
			TestCase: testutil.TestCase{
				Name: "API key with open access \\ update location instance with valid properties \\ returns updated instance",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[adventure_game.UpdateOneAdventureGameLocationInstance]
				},
				RequestPathParams: func(d harness.Data) map[string]string {
					return map[string]string{
						":game_instance_id":     gameInstanceRec.ID,
						":location_instance_id": gameLocationInstance.ID,
					}
				},
				RequestBody: func(d harness.Data) any {
					return api.AdventureGameLocationInstanceRequest{
						GameID:         gameInstanceRec.GameID,
						GameInstanceID: gameInstanceRec.ID,
						GameLocationID: locationRec.ID,
					}
				},
				ResponseDecoder: testCaseResponseDecoder,
				ResponseCode:    http.StatusOK,
			},
			expectResponse: func(d harness.Data, req api.AdventureGameLocationInstanceRequest) api.AdventureGameLocationInstanceResponse {
				return api.AdventureGameLocationInstanceResponse{
					Data: &api.AdventureGameLocationInstanceResponseData{
						GameInstanceID: req.GameInstanceID,
						GameLocationID: req.GameLocationID,
						GameID:         gameInstanceRec.GameID,
					},
				}
			},
		},
		{
			TestCase: testutil.TestCase{
				Name: "API key with open access \\ delete location instance with valid ID \\ returns no content",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[adventure_game.DeleteOneAdventureGameLocationInstance]
				},
				RequestPathParams: func(d harness.Data) map[string]string {
					return map[string]string{
						":game_instance_id":     gameInstanceRec.ID,
						":location_instance_id": gameLocationInstance.ID,
					}
				},
				ResponseCode: http.StatusNoContent,
			},
			expectResponse: func(d harness.Data, req api.AdventureGameLocationInstanceRequest) api.AdventureGameLocationInstanceResponse {
				return api.AdventureGameLocationInstanceResponse{}
			},
		},
	}

	for _, testCase := range testCases {
		t.Logf("Running test >%s<", testCase.Name)

		t.Run(testCase.Name, func(t *testing.T) {
			testFunc := func(method string, body interface{}) {
				if testCase.TestResponseCode() == http.StatusNoContent {
					// No content expected
					return
				}

				require.NotNil(t, body, "Response body is not nil")
				resp, ok := body.(api.AdventureGameLocationInstanceResponse)
				require.True(t, ok, "Response body is of type api.AdventureGameLocationInstanceResponse")
				gliResp := resp.Data
				require.NotNil(t, gliResp, "LocationInstanceResponseData is not nil")
				require.NotEmpty(t, gliResp.ID, "LocationInstance ID is not empty")
				require.NotEmpty(t, gliResp.GameInstanceID, "LocationInstance GameInstanceID is not empty")
				require.NotEmpty(t, gliResp.GameLocationID, "LocationInstance GameLocationID is not empty")
				if testCase.expectResponse != nil {
					xResp := testCase.expectResponse(
						th.Data,
						testCase.TestRequestBody(th.Data).(api.AdventureGameLocationInstanceRequest),
					).Data
					require.Equal(t, xResp.GameInstanceID, gliResp.GameInstanceID, "LocationInstance GameInstanceID matches expected")
					require.Equal(t, xResp.GameLocationID, gliResp.GameLocationID, "LocationInstance GameLocationID equals expected")
				}
			}

			testutil.RunTestCase(t, th, &testCase, testFunc)
		})
	}
}
