package runner

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/playbymail/core/server"
	"gitlab.com/alienspaces/playbymail/internal/harness"
	"gitlab.com/alienspaces/playbymail/schema"
)

func Test_getGameLocationInstanceHandler(t *testing.T) {
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

	testCaseCollectionResponseDecoder := testCaseResponseDecoderGeneric[schema.GameLocationInstanceCollectionResponse]
	testCaseResponseDecoder := testCaseResponseDecoderGeneric[schema.GameLocationInstanceResponse]

	// Setup: get a location instance for reference
	instanceRec, err := th.Data.GetGameLocationInstanceRecByRef(harness.GameLocationInstanceOneRef)
	require.NoError(t, err, "GetGameLocationInstanceRecByLocationRef returns without error")
	require.NotNil(t, instanceRec, "Location instance exists for reference")

	rec := instanceRec

	testCases := []testCase{
		{
			TestCase: TestCase{
				Name: "API key with open access \\ get many location instances \\ returns expected instances",
				HandlerConfig: func(rnr *Runner) server.HandlerConfig {
					return rnr.HandlerConfig[getManyGameLocationInstances]
				},
				RequestQueryParams: func(d harness.Data) map[string]any {
					return map[string]any{
						"game_location_id": rec.GameLocationID,
						"page_size":        10,
						"page_number":      1,
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
				Name: "API key with open access \\ get one location instance with valid ID \\ returns expected instance",
				HandlerConfig: func(rnr *Runner) server.HandlerConfig {
					return rnr.HandlerConfig[getOneGameLocationInstance]
				},
				RequestPathParams: func(d harness.Data) map[string]string {
					return map[string]string{
						":game_location_instance_id": rec.ID,
					}
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

				var responses []*schema.GameLocationInstanceResponseData
				if testCase.collectionRequest {
					responses = body.(schema.GameLocationInstanceCollectionResponse).Data
				} else {
					responses = append(responses, body.(schema.GameLocationInstanceResponse).Data)
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
					require.NotEmpty(t, d.ID, "LocationInstance ID is not empty")
					require.NotEmpty(t, d.GameInstanceID, "LocationInstance GameInstanceID is not empty")
					require.NotEmpty(t, d.GameLocationID, "LocationInstance GameLocationID is not empty")
				}
			}

			RunTestCase(t, th, &testCase, testFunc)
		})
	}
}

func Test_createUpdateDeleteGameLocationInstanceHandler(t *testing.T) {
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
		expectResponse func(d harness.Data, req schema.GameLocationInstanceRequest) schema.GameLocationInstanceResponse
	}

	testCaseResponseDecoder := testCaseResponseDecoderGeneric[schema.GameLocationInstanceResponse]

	gameInstanceRec, err := th.Data.GetGameInstanceRecByRef(harness.GameInstanceOneRef)
	require.NoError(t, err, "GetGameInstanceRecByRef returns without error")
	require.NotNil(t, gameInstanceRec, "Game instance exists for reference")

	locationRec, err := th.Data.GetGameLocationRecByRef(harness.GameLocationOneRef)
	require.NoError(t, err, "GetGameLocationRecByRef returns without error")

	gameLocationInstanceRec, err := th.Data.GetGameLocationInstanceRecByRef(harness.GameLocationInstanceOneRef)
	require.NoError(t, err, "GetGameLocationInstanceRecByRef returns without error")
	require.NotNil(t, gameLocationInstanceRec, "Location instance exists for reference")
	gameLocationInstance := gameLocationInstanceRec

	testCases := []testCase{
		{
			TestCase: TestCase{
				Name: "API key with open access \\ create location instance with valid properties \\ returns created instance",
				HandlerConfig: func(rnr *Runner) server.HandlerConfig {
					return rnr.HandlerConfig[createGameLocationInstance]
				},
				RequestBody: func(d harness.Data) interface{} {
					return schema.GameLocationInstanceRequest{
						GameInstanceID: gameInstanceRec.ID,
						GameLocationID: locationRec.ID,
						GameID:         gameInstanceRec.GameID,
					}
				},
				ResponseDecoder: testCaseResponseDecoder,
				ResponseCode:    http.StatusCreated,
			},
			expectResponse: func(d harness.Data, req schema.GameLocationInstanceRequest) schema.GameLocationInstanceResponse {
				return schema.GameLocationInstanceResponse{
					Data: &schema.GameLocationInstanceResponseData{
						GameInstanceID: req.GameInstanceID,
						GameLocationID: req.GameLocationID,
						GameID:         req.GameID,
					},
				}
			},
		},
		{
			TestCase: TestCase{
				Name: "API key with open access \\ update location instance with valid properties \\ returns updated instance",
				HandlerConfig: func(rnr *Runner) server.HandlerConfig {
					return rnr.HandlerConfig[updateGameLocationInstance]
				},
				RequestPathParams: func(d harness.Data) map[string]string {
					return map[string]string{
						":game_location_instance_id": gameLocationInstance.ID,
					}
				},
				RequestBody: func(d harness.Data) interface{} {
					return schema.GameLocationInstanceRequest{
						GameInstanceID: gameInstanceRec.ID,
						GameLocationID: locationRec.ID,
						GameID:         gameInstanceRec.GameID,
					}
				},
				ResponseDecoder: testCaseResponseDecoder,
				ResponseCode:    http.StatusOK,
			},
			expectResponse: func(d harness.Data, req schema.GameLocationInstanceRequest) schema.GameLocationInstanceResponse {
				return schema.GameLocationInstanceResponse{
					Data: &schema.GameLocationInstanceResponseData{
						GameInstanceID: req.GameInstanceID,
						GameLocationID: req.GameLocationID,
						GameID:         req.GameID,
					},
				}
			},
		},
		{
			TestCase: TestCase{
				Name: "API key with open access \\ delete location instance with valid ID \\ returns no content",
				HandlerConfig: func(rnr *Runner) server.HandlerConfig {
					return rnr.HandlerConfig[deleteGameLocationInstance]
				},
				RequestPathParams: func(d harness.Data) map[string]string {
					return map[string]string{
						":game_location_instance_id": gameLocationInstance.ID,
					}
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
				resp, ok := body.(schema.GameLocationInstanceResponse)
				require.True(t, ok, "Response body is of type schema.GameLocationInstanceResponse")
				gliResp := resp.Data
				require.NotNil(t, gliResp, "LocationInstanceResponseData is not nil")
				require.NotEmpty(t, gliResp.ID, "LocationInstance ID is not empty")
				require.NotEmpty(t, gliResp.GameInstanceID, "LocationInstance GameInstanceID is not empty")
				require.NotEmpty(t, gliResp.GameLocationID, "LocationInstance GameLocationID is not empty")
				if testCase.expectResponse != nil {
					xResp := testCase.expectResponse(
						th.Data,
						testCase.TestRequestBody(th.Data).(schema.GameLocationInstanceRequest),
					).Data
					require.Equal(t, xResp.GameInstanceID, gliResp.GameInstanceID, "LocationInstance GameInstanceID matches expected")
					require.Equal(t, xResp.GameLocationID, gliResp.GameLocationID, "LocationInstance GameLocationID equals expected")
				}
			}

			RunTestCase(t, th, &testCase, testFunc)
		})
	}
}
