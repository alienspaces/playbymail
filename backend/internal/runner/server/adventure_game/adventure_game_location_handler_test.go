package adventure_game

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/playbymail/core/server"
	"gitlab.com/alienspaces/playbymail/internal/harness"
	"gitlab.com/alienspaces/playbymail/internal/runner/server/testutil"
	"gitlab.com/alienspaces/playbymail/schema"
)

func Test_getGameLocationHandler(t *testing.T) {
	t.Parallel()

	th := testutil.NewTestHarness(t)
	require.NotNil(t, th, "newTestHarness returns without error")

	_, err := th.Setup()
	require.NoError(t, err, "Test data setup returns without error")
	defer func() {
		err = th.Teardown()
		require.NoError(t, err, "Test data teardown returns without error")
	}()

	gameRec, err := th.Data.GetGameRecByRef(harness.GameOneRef)
	require.NoError(t, err, "GetGameRecByRef returns without error")

	locationRec, err := th.Data.GetGameLocationRecByRef(harness.GameLocationOneRef)
	require.NoError(t, err, "GetGameLocationRecByRef returns without error")

	type testCase struct {
		testutil.TestCase
		collectionRequest     bool
		collectionRecordCount int
	}

	testCaseCollectionResponseDecoder := testutil.TestCaseResponseDecoderGeneric[schema.AdventureGameLocationCollectionResponse]
	testCaseResponseDecoder := testutil.TestCaseResponseDecoderGeneric[schema.AdventureGameLocationResponse]

	testCases := []testCase{
		{
			TestCase: testutil.TestCase{
				Name: "API key with open access \\ get many locations \\ returns expected locations",
				HandlerConfig: func(rnr *server.Runner) server.HandlerConfig {
					return rnr.HandlerConfig[getManyAdventureGameLocations]
				},
				RequestQueryParams: func(d harness.Data) map[string]any {
					return map[string]any{
						"page_size":   10,
						"page_number": 1,
					}
				},
				RequestPathParams: func(d harness.Data) map[string]string {
					return map[string]string{
						":game_id": gameRec.ID,
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
				Name: "API key with open access \\ get one location with valid location ID \\ returns expected location",
				HandlerConfig: func(rnr *server.Runner) server.HandlerConfig {
					return rnr.HandlerConfig[getOneAdventureGameLocation]
				},
				RequestPathParams: func(d harness.Data) map[string]string {
					params := map[string]string{
						":game_id":     gameRec.ID,
						":location_id": locationRec.ID,
					}
					return params
				},
				ResponseDecoder: testCaseResponseDecoder,
				ResponseCode:    http.StatusOK,
			},
			collectionRequest: false,
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

				var responses []*schema.AdventureGameLocationResponseData
				if testCase.collectionRequest {
					responses = body.(schema.AdventureGameLocationCollectionResponse).Data
				} else {
					responses = append(responses, body.(schema.AdventureGameLocationResponse).Data)
				}

				if testCase.collectionRequest {
					require.Equal(t, testCase.collectionRecordCount, len(responses), "Response record count length equals expected")
				}

				if testCase.collectionRequest && testCase.collectionRecordCount == 0 {
					require.Empty(t, responses, "Response body should be empty")
				} else {
					require.NotEmpty(t, responses, "Response body is not empty")
				}

				for _, d := range responses {
					require.NotEmpty(t, d.ID, "Location ID is not empty")
					require.NotEmpty(t, d.GameID, "Location GameID is not empty")
					require.NotEmpty(t, d.Name, "Location Name is not empty")
				}
			}

			testutil.RunTestCase(t, th, &testCase, testFunc)
		})
	}
}

func Test_createUpdateDeleteGameLocationHandler(t *testing.T) {
	t.Parallel()

	th := testutil.NewTestHarness(t)
	require.NotNil(t, th, "newTestHarness returns without error")

	_, err := th.Setup()
	require.NoError(t, err, "Test data setup returns without error")
	defer func() {
		err = th.Teardown()
		require.NoError(t, err, "Test data teardown returns without error")
	}()

	type testCase struct {
		testutil.TestCase
		expectResponse func(d harness.Data, req schema.AdventureGameLocationRequest) schema.AdventureGameLocationResponse
	}

	testCaseResponseDecoder := testutil.TestCaseResponseDecoderGeneric[schema.AdventureGameLocationResponse]

	gameRec, err := th.Data.GetGameRecByRef(harness.GameOneRef)
	require.NoError(t, err, "GetGameRecByRef returns without error")
	locationRec, err := th.Data.GetGameLocationRecByRef(harness.GameLocationOneRef)
	require.NoError(t, err, "GetGameLocationRecByRef returns without error")

	testCases := []testCase{
		{
			TestCase: testutil.TestCase{
				Name: "API key with open access \\ create location with valid properties \\ returns created location",
				HandlerConfig: func(rnr *server.Runner) server.HandlerConfig {
					return rnr.HandlerConfig[createOneAdventureGameLocation]
				},
				RequestBody: func(d harness.Data) interface{} {
					return schema.AdventureGameLocationRequest{
						Name:        "Test Location",
						Description: "Test Description",
					}
				},
				RequestPathParams: func(d harness.Data) map[string]string {
					return map[string]string{
						":game_id": gameRec.ID,
					}
				},
				ResponseDecoder: testCaseResponseDecoder,
				ResponseCode:    http.StatusCreated,
			},
			expectResponse: func(d harness.Data, req schema.AdventureGameLocationRequest) schema.AdventureGameLocationResponse {
				return schema.AdventureGameLocationResponse{
					Data: &schema.AdventureGameLocationResponseData{
						GameID:      gameRec.ID,
						Name:        req.Name,
						Description: req.Description,
					},
				}
			},
		},
		{
			TestCase: testutil.TestCase{
				Name: "API key with open access \\ update location with valid properties \\ returns updated location",
				HandlerConfig: func(rnr *server.Runner) server.HandlerConfig {
					return rnr.HandlerConfig[updateOneAdventureGameLocation]
				},
				RequestPathParams: func(d harness.Data) map[string]string {
					params := map[string]string{
						":game_id":     gameRec.ID,
						":location_id": locationRec.ID,
					}
					return params
				},
				RequestBody: func(d harness.Data) interface{} {
					return schema.AdventureGameLocationRequest{
						Name:        "Updated Location",
						Description: "Updated Description",
					}
				},
				ResponseDecoder: testCaseResponseDecoder,
				ResponseCode:    http.StatusOK,
			},
			expectResponse: func(d harness.Data, req schema.AdventureGameLocationRequest) schema.AdventureGameLocationResponse {
				return schema.AdventureGameLocationResponse{
					Data: &schema.AdventureGameLocationResponseData{
						GameID:      gameRec.ID,
						Name:        req.Name,
						Description: req.Description,
					},
				}
			},
		},
		{
			TestCase: testutil.TestCase{
				Name: "API key with open access \\ delete location with valid location ID \\ returns no content",
				HandlerConfig: func(rnr *server.Runner) server.HandlerConfig {
					return rnr.HandlerConfig[deleteOneAdventureGameLocation]
				},
				RequestPathParams: func(d harness.Data) map[string]string {
					params := map[string]string{
						":game_id":     gameRec.ID,
						":location_id": locationRec.ID,
					}
					return params
				},
				ResponseDecoder: testCaseResponseDecoder,
				ResponseCode:    http.StatusNoContent,
			},
			expectResponse: nil,
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
				resp, ok := body.(schema.AdventureGameLocationResponse)
				require.True(t, ok, "Response body is of type schema.AdventureGameLocationResponse")
				lResp := resp.Data
				require.NotNil(t, lResp, "LocationResponseData is not nil")
				require.NotEmpty(t, lResp.ID, "Location ID is not empty")
				require.NotEmpty(t, lResp.GameID, "Location GameID is not empty")
				require.NotEmpty(t, lResp.Name, "Location Name is not empty")
				if testCase.expectResponse != nil {
					xResp := testCase.expectResponse(
						th.Data,
						testCase.TestRequestBody(th.Data).(schema.AdventureGameLocationRequest),
					).Data
					require.Equal(t, xResp.GameID, lResp.GameID, "Location GameID matches expected")
					require.Equal(t, xResp.Name, lResp.Name, "Location Name equals expected")
					require.Equal(t, xResp.Description, lResp.Description, "Location Description equals expected")
				}
			}

			testutil.RunTestCase(t, th, &testCase, testFunc)
		})
	}
}
