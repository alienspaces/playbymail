package runner

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/playbymail/core/server"
	"gitlab.com/alienspaces/playbymail/internal/harness"
	"gitlab.com/alienspaces/playbymail/schema"
)

func Test_getLocationHandler(t *testing.T) {
	t.Parallel()

	th := newTestHarness(t)
	require.NotNil(t, th, "newTestHarness returns without error")

	_, err := th.Setup()
	require.NoError(t, err, "Test data setup returns without error")
	defer func() {
		err = th.Teardown()
		require.NoError(t, err, "Test data teardown returns without error")
	}()

	type testCase struct {
		TestCase
		collectionRequest     bool
		collectionRecordCount int
	}

	testCaseCollectionResponseDecoder := testCaseResponseDecoderGeneric[schema.GameLocationCollectionResponse]
	testCaseResponseDecoder := testCaseResponseDecoderGeneric[schema.GameLocationResponse]

	// Setup: get a location for reference
	locationRec, err := th.Data.GetGameLocationRecByRef("game-location-one")
	require.NoError(t, err, "GetGameLocationRecByRef returns without error")

	testCases := []testCase{
		{
			TestCase: TestCase{
				Name: "API key with open access \\ get many locations \\ returns expected locations",
				HandlerConfig: func(rnr *Runner) server.HandlerConfig {
					return rnr.HandlerConfig[getManyGameLocations]
				},
				RequestQueryParams: func(d harness.Data) map[string]any {
					return map[string]any{
						"id":          locationRec.ID,
						"page_size":   10,
						"page_number": 1,
					}
				},
				ResponseDecoder: testCaseCollectionResponseDecoder,
				ResponseCode:    http.StatusOK,
			},
			collectionRequest:     true,
			collectionRecordCount: 1,
		},
		{
			TestCase: TestCase{
				Name: "API key with open access \\ get one location with valid location ID \\ returns expected location",
				HandlerConfig: func(rnr *Runner) server.HandlerConfig {
					return rnr.HandlerConfig[getOneGameLocation]
				},
				RequestPathParams: func(d harness.Data) map[string]string {
					params := map[string]string{
						":game_location_id": locationRec.ID,
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

				responses := []*schema.GameLocationResponseData{}
				if testCase.collectionRequest {
					responses = body.(schema.GameLocationCollectionResponse)
				} else {
					responses = append(responses, body.(schema.GameLocationResponse).GameLocationResponseData)
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

			RunTestCase(t, th, &testCase, testFunc)
		})
	}
}

func Test_createUpdateDeleteLocationHandler(t *testing.T) {
	t.Parallel()

	th := newTestHarness(t)
	require.NotNil(t, th, "newTestHarness returns without error")

	_, err := th.Setup()
	require.NoError(t, err, "Test data setup returns without error")
	defer func() {
		err = th.Teardown()
		require.NoError(t, err, "Test data teardown returns without error")
	}()

	type testCase struct {
		TestCase
		expectResponse func(d harness.Data, req schema.GameLocationRequest) schema.GameLocationResponse
	}

	testCaseResponseDecoder := testCaseResponseDecoderGeneric[schema.GameLocationResponse]

	gameRec, err := th.Data.GetGameRecByRef(harness.GameOneRef)
	require.NoError(t, err, "GetGameRecByRef returns without error")
	locationRec, err := th.Data.GetGameLocationRecByRef("game-location-one")
	require.NoError(t, err, "GetGameLocationRecByRef returns without error")

	testCases := []testCase{
		{
			TestCase: TestCase{
				Name: "API key with open access \\ create location with valid properties \\ returns created location",
				HandlerConfig: func(rnr *Runner) server.HandlerConfig {
					return rnr.HandlerConfig[createGameLocation]
				},
				RequestBody: func(d harness.Data) interface{} {
					return schema.GameLocationRequest{
						GameID:      gameRec.ID,
						Name:        "Test Location",
						Description: "Test Description",
					}
				},
				ResponseDecoder: testCaseResponseDecoder,
				ResponseCode:    http.StatusCreated,
			},
			expectResponse: func(d harness.Data, req schema.GameLocationRequest) schema.GameLocationResponse {
				return schema.GameLocationResponse{
					GameLocationResponseData: &schema.GameLocationResponseData{
						GameID:      req.GameID,
						Name:        req.Name,
						Description: req.Description,
					},
				}
			},
		},
		{
			TestCase: TestCase{
				Name: "API key with open access \\ update location with valid properties \\ returns updated location",
				HandlerConfig: func(rnr *Runner) server.HandlerConfig {
					return rnr.HandlerConfig[updateGameLocation]
				},
				RequestPathParams: func(d harness.Data) map[string]string {
					params := map[string]string{
						":game_location_id": locationRec.ID,
					}
					return params
				},
				RequestBody: func(d harness.Data) interface{} {
					return schema.GameLocationRequest{
						GameID:      gameRec.ID,
						Name:        "Updated Location",
						Description: "Updated Description",
					}
				},
				ResponseDecoder: testCaseResponseDecoder,
				ResponseCode:    http.StatusOK,
			},
			expectResponse: func(d harness.Data, req schema.GameLocationRequest) schema.GameLocationResponse {
				return schema.GameLocationResponse{
					GameLocationResponseData: &schema.GameLocationResponseData{
						GameID:      req.GameID,
						Name:        req.Name,
						Description: req.Description,
					},
				}
			},
		},
		{
			TestCase: TestCase{
				Name: "API key with open access \\ delete location with valid location ID \\ returns no content",
				HandlerConfig: func(rnr *Runner) server.HandlerConfig {
					return rnr.HandlerConfig[deleteGameLocation]
				},
				RequestPathParams: func(d harness.Data) map[string]string {
					params := map[string]string{
						":game_location_id": locationRec.ID,
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
				resp, ok := body.(schema.GameLocationResponse)
				require.True(t, ok, "Response body is of type schema.GameLocationResponse")
				lResp := resp.GameLocationResponseData
				require.NotNil(t, lResp, "LocationResponseData is not nil")
				require.NotEmpty(t, lResp.ID, "Location ID is not empty")
				require.NotEmpty(t, lResp.GameID, "Location GameID is not empty")
				require.NotEmpty(t, lResp.Name, "Location Name is not empty")
				if testCase.expectResponse != nil {
					xResp := testCase.expectResponse(
						th.Data,
						testCase.TestRequestBody(th.Data).(schema.GameLocationRequest),
					).GameLocationResponseData
					require.Equal(t, xResp.GameID, lResp.GameID, "Location GameID matches expected")
					require.Equal(t, xResp.Name, lResp.Name, "Location Name equals expected")
					require.Equal(t, xResp.Description, lResp.Description, "Location Description equals expected")
				}
			}

			RunTestCase(t, th, &testCase, testFunc)
		})
	}
}
