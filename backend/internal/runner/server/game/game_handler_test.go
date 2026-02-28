package game_test

import (
	"net/http"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"
	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/playbymail/core/server"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/core/type/storer"
	"gitlab.com/alienspaces/playbymail/internal/harness"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
	"gitlab.com/alienspaces/playbymail/internal/runner/server/game"
	"gitlab.com/alienspaces/playbymail/internal/turnsheet"
	"gitlab.com/alienspaces/playbymail/internal/utils/config"
	"gitlab.com/alienspaces/playbymail/internal/utils/testutil"
	"gitlab.com/alienspaces/playbymail/schema/api/game_schema"
)

func Test_getGameHandler(t *testing.T) {
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
		collectionRequest     bool
		collectionRecordCount int
	}

	testCaseCollectionResponseDecoder := testutil.TestCaseResponseDecoderGeneric[game_schema.GameCollectionResponse]
	testCaseResponseDecoder := testutil.TestCaseResponseDecoderGeneric[game_schema.GameResponse]

	// Setup: get a game for reference
	gameRec, err := th.Data.GetGameRecByRef(harness.GameOneRef)
	require.NoError(t, err, "GetGameRecByRef returns without error")

	testCases := []testCase{
		{
			TestCase: testutil.TestCase{
				Name: "authenticated user when get many games then returns expected games",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[game.GetManyGames]
				},
				RequestHeaders: testutil.AuthHeaderStandard,
				RequestQueryParams: func(d harness.Data) map[string]any {
					return map[string]any{
						"id":          gameRec.ID,
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
			TestCase: testutil.TestCase{
				Name: "authenticated user when get one game with valid game ID then returns expected game",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[game.GetOneGame]
				},
				RequestHeaders: testutil.AuthHeaderStandard,
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
			testFunc := func(method string, body any) {
				if testCase.TestResponseCode() != http.StatusOK {
					return
				}

				require.NotNil(t, body, "Response body is not nil")

				var responses []*game_schema.GameResponseData
				if testCase.collectionRequest {
					responses = body.(game_schema.GameCollectionResponse).Data
				} else {
					responses = append(responses, body.(game_schema.GameResponse).Data)
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

			testutil.RunTestCase(t, th, &testCase, testFunc)
		})
	}
}

func Test_createUpdateDeleteGameHandler(t *testing.T) {
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
		expectResponse func(d harness.Data, req game_schema.GameRequest) game_schema.GameResponse
	}

	testCaseResponseDecoder := testutil.TestCaseResponseDecoderGeneric[game_schema.GameResponse]

	testCases := []testCase{
		{
			TestCase: testutil.TestCase{
				Name: "authenticated designer when create game with valid properties then returns created game",
				NewRunner: func(cfg config.Config, l logger.Logger, s storer.Storer, j *river.Client[pgx.Tx], scanner turnsheet.TurnSheetScanner, d harness.Data) (testutil.TestRunnerer, error) {
					return testutil.NewTestRunner(cfg, l, s, j, scanner)
				},
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[game.CreateOneGame]
				},
				RequestHeaders: testutil.AuthHeaderProDesigner,
				RequestBody: func(d harness.Data) any {
					return game_schema.GameRequest{
						Name:              "Test Game",
						GameType:          game_record.GameTypeAdventure,
						TurnDurationHours: 168, // 1 week
						Description:       "A test adventure game description",
					}
				},
				ResponseDecoder: testCaseResponseDecoder,
				ResponseCode:    http.StatusCreated,
			},
			expectResponse: func(d harness.Data, req game_schema.GameRequest) game_schema.GameResponse {
				return game_schema.GameResponse{
					Data: &game_schema.GameResponseData{
						Name: req.Name,
					},
				}
			},
		},
		{
			TestCase: testutil.TestCase{
				Name: "authenticated designer when update game with valid properties then returns updated game",
				NewRunner: func(cfg config.Config, l logger.Logger, s storer.Storer, j *river.Client[pgx.Tx], scanner turnsheet.TurnSheetScanner, d harness.Data) (testutil.TestRunnerer, error) {
					return testutil.NewTestRunner(cfg, l, s, j, scanner)
				},
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[game.UpdateOneGame]
				},
				RequestHeaders: testutil.AuthHeaderProDesigner,
				RequestPathParams: func(d harness.Data) map[string]string {
					gameRec, err := d.GetGameRecByRef(harness.GameOneRef)
					require.NoError(t, err, "GetGameRecByRef returns without error")
					params := map[string]string{
						":game_id": gameRec.ID,
					}
					return params
				},
				RequestBody: func(d harness.Data) any {
					gameRec, err := d.GetGameRecByRef(harness.GameOneRef)
					require.NoError(t, err, "GetGameRecByRef returns without error")
					return game_schema.GameRequest{
						Name:              gameRec.Name,
						GameType:          gameRec.GameType,
						TurnDurationHours: gameRec.TurnDurationHours,
						Description:       gameRec.Description,
					}
				},
				ResponseDecoder: testCaseResponseDecoder,
				ResponseCode:    http.StatusOK,
			},
			expectResponse: func(d harness.Data, req game_schema.GameRequest) game_schema.GameResponse {
				return game_schema.GameResponse{
					Data: &game_schema.GameResponseData{
						Name: req.Name,
					},
				}
			},
		},
		{
			TestCase: testutil.TestCase{
				Name: "authenticated designer when delete game with valid game ID then returns no content",
				NewRunner: func(cfg config.Config, l logger.Logger, s storer.Storer, j *river.Client[pgx.Tx], scanner turnsheet.TurnSheetScanner, d harness.Data) (testutil.TestRunnerer, error) {
					return testutil.NewTestRunner(cfg, l, s, j, scanner)
				},
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[game.DeleteOneGame]
				},
				RequestHeaders: testutil.AuthHeaderProDesigner,
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
			testFunc := func(method string, body any) {
				if testCase.TestResponseCode() != http.StatusCreated {
					return
				}

				require.NotNil(t, body, "Response body is not nil")

				aResp := body.(game_schema.GameResponse).Data
				xResp := testCase.expectResponse(
					th.Data,
					testCase.TestRequestBody(th.Data).(game_schema.GameRequest),
				).Data

				require.NotEmpty(t, aResp, "Response body is not empty")
				require.NotEmpty(t, aResp.ID, "Game ID is not empty")
				require.Equal(t, xResp.Name, aResp.Name, "Game Name equals expected")
				require.False(t, aResp.CreatedAt.IsZero(), "Game CreatedAt is not zero")
			}

			testutil.RunTestCase(t, th, &testCase, testFunc)
		})
	}
}
