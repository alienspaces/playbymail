package runner

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/playbymail/core/server"
	"gitlab.com/alienspaces/playbymail/internal/harness"
	"gitlab.com/alienspaces/playbymail/schema"
)

func Test_getGameHandler(t *testing.T) {
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

	testCaseCollectionResponseDecoder := testCaseResponseDecoderGeneric[schema.GameCollectionResponse]
	testCaseResponseDecoder := testCaseResponseDecoderGeneric[schema.GameResponse]

	testCases := []testCase{
		{
			TestCase: TestCase{
				Name: "API key with open access \\ get many games \\ returns expected games",
				HandlerConfig: func(rnr *Runner) server.HandlerConfig {
					return rnr.HandlerConfig[getManyGames]
				},
				RequestQueryParams: func(d harness.Data) map[string]any {
					gameRec, err := d.GetGameRecByRef(harness.GameOneRef)
					require.NoError(t, err, "GetGameRecByRef returns without error")
					return map[string]any{
						"page_size":   10,
						"page_number": 1,
						"id":          gameRec.ID,
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
				Name: "API key with open access \\ get one game with valid game ID \\ returns expected game",
				HandlerConfig: func(rnr *Runner) server.HandlerConfig {
					return rnr.HandlerConfig[getOneGame]
				},
				RequestPathParams: func(d harness.Data) map[string]string {
					gameRec, err := d.GetGameRecByRef(harness.GameOneRef)
					require.NoError(t, err, "GetGameRecByRef returns without error")
					params := map[string]string{
						":game_id": gameRec.ID,
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

				responses := []*schema.GameResponseData{}
				if testCase.collectionRequest {
					responses = body.(schema.GameCollectionResponse)
				} else {
					responses = append(responses, body.(schema.GameResponse).GameResponseData)
				}

				if testCase.collectionRequest {
					require.Equal(t, testCase.collectionRecordCount, len(responses), "Response record count length equals expected")
				}

				if testCase.collectionRequest && testCase.collectionRecordCount == 0 {
					require.Empty(t, responses, "Response body should be empty")
				} else {
					require.NotEmpty(t, responses, "Response body is not empty")
				}

				// The response is also validated against the response JSON schema
				for _, d := range responses {
					require.NotEmpty(t, d.ID, "Game ID is not empty")
					require.False(t, d.CreatedAt.IsZero(), "Game CreatedAt is not zero")
				}
			}

			RunTestCase(t, th, &testCase, testFunc)
		})
	}
}

func Test_createUpdateDeleteGameHandler(t *testing.T) {
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
		expectResponse func(d harness.Data, req schema.GameRequest) schema.GameResponse
	}

	testCaseResponseDecoder := testCaseResponseDecoderGeneric[schema.GameResponse]

	testCases := []testCase{
		{
			TestCase: TestCase{
				Name: "API key with open access \\ create game with valid properties \\ returns created game",
				HandlerConfig: func(rnr *Runner) server.HandlerConfig {
					return rnr.HandlerConfig[createGame]
				},
				RequestBody: func(d harness.Data) interface{} {
					return schema.GameRequest{
						Name: "Test Game",
					}
				},
				ResponseDecoder: testCaseResponseDecoder,
				ResponseCode:    http.StatusCreated,
			},
			expectResponse: func(d harness.Data, req schema.GameRequest) schema.GameResponse {
				return schema.GameResponse{
					GameResponseData: &schema.GameResponseData{
						Name: req.Name,
					},
				}
			},
		},
		{
			TestCase: TestCase{
				Name: "API key with open access \\ update game with valid properties \\ returns updated game",
				HandlerConfig: func(rnr *Runner) server.HandlerConfig {
					return rnr.HandlerConfig[updateGame]
				},
				RequestPathParams: func(d harness.Data) map[string]string {
					gameRec, err := d.GetGameRecByRef(harness.GameOneRef)
					require.NoError(t, err, "GetGameRecByRef returns without error")
					params := map[string]string{
						":game_id": gameRec.ID,
					}
					return params
				},
				RequestBody: func(d harness.Data) interface{} {
					return schema.GameRequest{
						Name: "Test Game",
					}
				},
				ResponseDecoder: testCaseResponseDecoder,
				ResponseCode:    http.StatusOK,
			},
			expectResponse: func(d harness.Data, req schema.GameRequest) schema.GameResponse {
				return schema.GameResponse{
					GameResponseData: &schema.GameResponseData{
						Name: req.Name,
					},
				}
			},
		},
		{
			TestCase: TestCase{
				Name: "API key with open access \\ delete game with valid game ID \\ returns no content",
				HandlerConfig: func(rnr *Runner) server.HandlerConfig {
					return rnr.HandlerConfig[deleteGame]
				},
				RequestPathParams: func(d harness.Data) map[string]string {
					gameRec, err := d.GetGameRecByRef(harness.GameOneRef)
					require.NoError(t, err, "GetGameRecByRef returns without error")
					params := map[string]string{
						":game_id": gameRec.ID,
					}
					return params
				},
				ResponseDecoder: testCaseResponseDecoder,
				ResponseCode:    http.StatusNoContent,
			},
		},
	}

	for _, testCase := range testCases {
		t.Logf("Running test >%s<\n", testCase.Name)

		t.Run(testCase.Name, func(t *testing.T) {
			testFunc := func(method string, body interface{}) {
				if testCase.TestResponseCode() != http.StatusCreated {
					return
				}

				require.NotNil(t, body, "Response body is not nil")

				aResp := body.(schema.GameResponse).GameResponseData
				xResp := testCase.expectResponse(
					th.Data,
					testCase.TestRequestBody(th.Data).(schema.GameRequest),
				).GameResponseData

				require.NotEmpty(t, aResp, "Response body is not empty")
				require.NotEmpty(t, aResp.ID, "Game ID is not empty")
				require.Equal(t, xResp.Name, aResp.Name, "Game Name equals expected")
				require.False(t, aResp.CreatedAt.IsZero(), "Game CreatedAt is not zero")
			}

			RunTestCase(t, th, &testCase, testFunc)
		})
	}
}
