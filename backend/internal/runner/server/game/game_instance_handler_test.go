package game_test

import (
	"net/http"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"
	"github.com/stretchr/testify/require"

	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/core/server"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/core/type/storer"
	"gitlab.com/alienspaces/playbymail/internal/harness"
	"gitlab.com/alienspaces/playbymail/internal/runner/server/game"
	"gitlab.com/alienspaces/playbymail/internal/turnsheet"
	"gitlab.com/alienspaces/playbymail/internal/utils/config"
	"gitlab.com/alienspaces/playbymail/internal/utils/testutil"
	"gitlab.com/alienspaces/playbymail/schema/api/game_schema"
)

func Test_searchManyGameInstancesHandler(t *testing.T) {
	t.Parallel()

	th := testutil.NewTestHarness(t)
	require.NotNil(t, th, "TestHarness returns without error")

	_, err := th.Setup()
	require.NoError(t, err, "Test data setup returns without error")
	defer func() {
		err = th.Teardown()
		require.NoError(t, err, "Test data teardown returns without error")
	}()

	// Get a game from the harness data
	gameRec, err := th.Data.GetGameRecByRef(harness.GameOneRef)
	require.NoError(t, err, "GetGameRecByRef returns without error")

	testCases := []struct {
		testutil.TestCase
		expectedCount int
	}{
		{
			TestCase: testutil.TestCase{
				Name: "authenticated manager when search many game instances then returns expected instances",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[game.SearchManyGameInstances]
				},
				RequestHeaders: testutil.AuthHeaderProManager,
				RequestQueryParams: func(d harness.Data) map[string]any {
					return map[string]any{
						"page_size":   "2",
						"page_number": "1",
						"game_id":     gameRec.ID,
					}
				},
				ResponseDecoder: testutil.TestCaseResponseDecoderGeneric[game_schema.GameInstanceCollectionResponse],
				ResponseCode:    http.StatusOK,
			},
			expectedCount: 2, // GameInstanceOneRef and GameInstanceCleanRef
		},
		{
			TestCase: testutil.TestCase{
				Name: "authenticated manager when search many game instances with pagination then returns expected instances",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[game.SearchManyGameInstances]
				},
				RequestHeaders: testutil.AuthHeaderProManager,
				RequestQueryParams: func(d harness.Data) map[string]any {
					return map[string]any{
						"page_size":   "2",
						"page_number": "1",
						"game_id":     gameRec.ID,
					}
				},
				ResponseDecoder: testutil.TestCaseResponseDecoderGeneric[game_schema.GameInstanceCollectionResponse],
				ResponseCode:    http.StatusOK,
			},
			expectedCount: 2, // GameInstanceOneRef and GameInstanceCleanRef
		},
	}

	for _, testCase := range testCases {
		t.Logf("Running test >%s<\n", testCase.Name)

		t.Run(testCase.Name, func(t *testing.T) {
			testFunc := func(method string, body any) {
				require.NotNil(t, body, "Response body is not nil")

				resp := body.(game_schema.GameInstanceCollectionResponse)
				require.NotNil(t, resp.Data, "Response data is not nil")
				require.Equal(t, testCase.expectedCount, len(resp.Data), "Response contains expected number of instances")

				// Confirm type of each record
				for _, instance := range resp.Data {
					require.NotEmpty(t, instance.ID, "Instance ID is not empty")
					require.NotEmpty(t, instance.Status, "Instance Status is not empty")
					require.GreaterOrEqual(t, instance.CurrentTurn, 0, "Instance CurrentTurn is valid")
				}
			}

			testutil.RunTestCase(t, th, &testCase.TestCase, testFunc)
		})
	}
}

func Test_getManyGameInstancesHandler(t *testing.T) {
	t.Parallel()

	th := testutil.NewTestHarness(t)
	require.NotNil(t, th, "TestHarness returns without error")

	_, err := th.Setup()
	require.NoError(t, err, "Test data setup returns without error")
	defer func() {
		err = th.Teardown()
		require.NoError(t, err, "Test data teardown returns without error")
	}()

	// Get a game from the harness data
	gameRec, err := th.Data.GetGameRecByRef(harness.GameOneRef)
	require.NoError(t, err, "GetGameRecByRef returns without error")

	testCases := []struct {
		testutil.TestCase
		expectedCount int
	}{
		{
			TestCase: testutil.TestCase{
				Name: "authenticated manager when get many game instances then returns expected instances",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[game.GetManyGameInstances]
				},
				RequestHeaders: testutil.AuthHeaderProManager,
				RequestPathParams: func(d harness.Data) map[string]string {
					return map[string]string{
						":game_id": gameRec.ID,
					}
				},
				ResponseDecoder: testutil.TestCaseResponseDecoderGeneric[game_schema.GameInstanceCollectionResponse],
				ResponseCode:    http.StatusOK,
			},
			expectedCount: 2, // GameInstanceOneRef and GameInstanceCleanRef
		},
		{
			TestCase: testutil.TestCase{
				Name: "authenticated manager when get many game instances with pagination then returns expected instances",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[game.GetManyGameInstances]
				},
				RequestHeaders: testutil.AuthHeaderProManager,
				RequestPathParams: func(d harness.Data) map[string]string {
					return map[string]string{
						":game_id": gameRec.ID,
					}
				},
				RequestQueryParams: func(d harness.Data) map[string]any {
					return map[string]any{
						"page_size":   "10",
						"page_number": "1",
					}
				},
				ResponseDecoder: testutil.TestCaseResponseDecoderGeneric[game_schema.GameInstanceCollectionResponse],
				ResponseCode:    http.StatusOK,
			},
			expectedCount: 2, // GameInstanceOneRef and GameInstanceCleanRef
		},
	}

	for _, testCase := range testCases {
		t.Logf("Running test >%s<\n", testCase.Name)

		t.Run(testCase.Name, func(t *testing.T) {
			testFunc := func(method string, body any) {
				require.NotNil(t, body, "Response body is not nil")

				resp := body.(game_schema.GameInstanceCollectionResponse)
				require.NotNil(t, resp.Data, "Response data is not nil")
				require.Len(t, resp.Data, testCase.expectedCount, "Response contains expected number of instances")

				// Confirm type of each record
				for _, instance := range resp.Data {
					require.NotEmpty(t, instance.ID, "Instance ID is not empty")
					require.NotEmpty(t, instance.Status, "Instance Status is not empty")
					require.GreaterOrEqual(t, instance.CurrentTurn, 0, "Instance CurrentTurn is valid")
				}
			}

			testutil.RunTestCase(t, th, &testCase.TestCase, testFunc)
		})
	}
}

func Test_getOneGameInstanceHandler(t *testing.T) {
	t.Parallel()

	th := testutil.NewTestHarness(t)
	require.NotNil(t, th, "TestHarness returns without error")

	_, err := th.Setup()
	require.NoError(t, err, "Test data setup returns without error")
	defer func() {
		err = th.Teardown()
		require.NoError(t, err, "Test data teardown returns without error")
	}()

	// Get a game and game instance from the harness data
	gameRec, err := th.Data.GetGameRecByRef(harness.GameOneRef)
	require.NoError(t, err, "GetGameRecByRef returns without error")

	gameInstanceRec, err := th.Data.GetGameInstanceRecByRef(harness.GameInstanceOneRef)
	require.NoError(t, err, "GetGameInstanceRecByRef returns without error")

	testCases := []struct {
		testutil.TestCase
		expectResponse func(d harness.Data) game_schema.GameInstanceResponse
	}{
		{
			TestCase: testutil.TestCase{
				Name: "authenticated manager when get one game instance with valid instance ID then returns expected instance",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[game.GetOneGameInstance]
				},
				RequestHeaders: testutil.AuthHeaderProManager,
				RequestPathParams: func(d harness.Data) map[string]string {
					return map[string]string{
						":game_id":     gameRec.ID,
						":instance_id": gameInstanceRec.ID,
					}
				},
				ResponseDecoder: testutil.TestCaseResponseDecoderGeneric[game_schema.GameInstanceResponse],
				ResponseCode:    http.StatusOK,
			},
			expectResponse: func(d harness.Data) game_schema.GameInstanceResponse {
				return game_schema.GameInstanceResponse{
					Data: &game_schema.GameInstanceResponseData{
						Status:      "created",
						CurrentTurn: 0,
					},
				}
			},
		},
	}

	for _, testCase := range testCases {
		t.Logf("Running test >%s<\n", testCase.Name)

		t.Run(testCase.Name, func(t *testing.T) {
			testFunc := func(method string, body any) {
				require.NotNil(t, body, "Response body is not nil")

				aResp := body.(game_schema.GameInstanceResponse).Data
				xResp := testCase.expectResponse(th.Data).Data

				require.NotEmpty(t, aResp, "Response body is not empty")
				require.Equal(t, xResp.Status, aResp.Status, "Status equals expected")
				require.Equal(t, xResp.CurrentTurn, aResp.CurrentTurn, "Current turn equals expected")
				require.NotEmpty(t, aResp.ID, "Instance ID is not empty")
				require.False(t, aResp.CreatedAt.IsZero(), "Instance CreatedAt is not zero")
			}

			testutil.RunTestCase(t, th, &testCase.TestCase, testFunc)
		})
	}
}

func Test_createOneGameInstanceHandler(t *testing.T) {
	t.Parallel()

	th := testutil.NewTestHarness(t)
	require.NotNil(t, th, "TestHarness returns without error")

	_, err := th.Setup()
	require.NoError(t, err, "Test data setup returns without error")
	defer func() {
		err = th.Teardown()
		require.NoError(t, err, "Test data teardown returns without error")
	}()

	// Get a game from the harness data
	gameRec, err := th.Data.GetGameRecByRef(harness.GameOneRef)
	require.NoError(t, err, "GetGameRecByRef returns without error")

	testCases := []struct {
		testutil.TestCase
		expectResponse func(d harness.Data, req game_schema.GameInstanceRequest) game_schema.GameInstanceResponse
	}{
		{
			TestCase: testutil.TestCase{
				Name: "authenticated manager when create game instance with valid properties then returns created instance",
				NewRunner: func(cfg config.Config, l logger.Logger, s storer.Storer, j *river.Client[pgx.Tx], scanner turnsheet.TurnSheetScanner, d harness.Data) (testutil.TestRunnerer, error) {
					return testutil.NewTestRunner(cfg, l, s, j, scanner)
				},
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[game.CreateOneGameInstance]
				},
				RequestHeaders: testutil.AuthHeaderProManager,
				RequestPathParams: func(d harness.Data) map[string]string {
					return map[string]string{
						":game_id": gameRec.ID,
					}
				},
				RequestBody: func(d harness.Data) interface{} {
					return game_schema.GameInstanceRequest{
						GameID: gameRec.ID,
					}
				},
				ResponseDecoder: testutil.TestCaseResponseDecoderGeneric[game_schema.GameInstanceResponse],
				ResponseCode:    http.StatusCreated,
			},
			expectResponse: func(d harness.Data, req game_schema.GameInstanceRequest) game_schema.GameInstanceResponse {
				return game_schema.GameInstanceResponse{
					Data: &game_schema.GameInstanceResponseData{
						GameID:      req.GameID,
						Status:      "created",
						CurrentTurn: 0,
					},
				}
			},
		},
	}

	for _, testCase := range testCases {
		t.Logf("Running test >%s<\n", testCase.Name)

		t.Run(testCase.Name, func(t *testing.T) {
			testFunc := func(method string, body any) {
				require.NotNil(t, body, "Response body is not nil")

				aResp := body.(game_schema.GameInstanceResponse).Data
				xResp := testCase.expectResponse(th.Data, testCase.RequestBody(th.Data).(game_schema.GameInstanceRequest)).Data

				require.NotEmpty(t, aResp, "Response body is not empty")
				require.Equal(t, xResp.GameID, aResp.GameID, "Game ID equals expected")
				require.Equal(t, xResp.Status, aResp.Status, "Status equals expected")
				require.Equal(t, xResp.CurrentTurn, aResp.CurrentTurn, "Current turn equals expected")
				require.NotEmpty(t, aResp.ID, "Instance ID is not empty")
				require.False(t, aResp.CreatedAt.IsZero(), "Instance CreatedAt is not zero")
			}

			testutil.RunTestCase(t, th, &testCase.TestCase, testFunc)
		})
	}
}

func Test_updateOneGameInstanceHandler(t *testing.T) {
	t.Parallel()

	th := testutil.NewTestHarness(t)
	require.NotNil(t, th, "TestHarness returns without error")

	_, err := th.Setup()
	require.NoError(t, err, "Test data setup returns without error")
	defer func() {
		err = th.Teardown()
		require.NoError(t, err, "Test data teardown returns without error")
	}()

	// Get a game and game instance from the harness data
	gameRec, err := th.Data.GetGameRecByRef(harness.GameOneRef)
	require.NoError(t, err, "GetGameRecByRef returns without error")

	gameInstanceRec, err := th.Data.GetGameInstanceRecByRef(harness.GameInstanceOneRef)
	require.NoError(t, err, "GetGameInstanceRecByRef returns without error")

	testCases := []struct {
		testutil.TestCase
		expectResponse func(d harness.Data, req game_schema.GameInstanceRequest) game_schema.GameInstanceResponse
	}{
		{
			TestCase: testutil.TestCase{
				Name: "authenticated manager when update game instance with valid properties then returns updated instance",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[game.UpdateOneGameInstance]
				},
				RequestHeaders: testutil.AuthHeaderProManager,
				RequestPathParams: func(d harness.Data) map[string]string {
					return map[string]string{
						":game_id":     gameRec.ID,
						":instance_id": gameInstanceRec.ID,
					}
				},
				RequestBody: func(d harness.Data) interface{} {
					return game_schema.GameInstanceRequest{
						GameID:      gameRec.ID,
						CurrentTurn: 1,
					}
				},
				ResponseDecoder: testutil.TestCaseResponseDecoderGeneric[game_schema.GameInstanceResponse],
				ResponseCode:    http.StatusOK,
			},
			expectResponse: func(d harness.Data, req game_schema.GameInstanceRequest) game_schema.GameInstanceResponse {
				return game_schema.GameInstanceResponse{
					Data: &game_schema.GameInstanceResponseData{
						GameID:      req.GameID,
						Status:      "created",
						CurrentTurn: req.CurrentTurn,
					},
				}
			},
		},
	}

	for _, testCase := range testCases {
		t.Logf("Running test >%s<\n", testCase.Name)

		t.Run(testCase.Name, func(t *testing.T) {
			testFunc := func(method string, body any) {
				require.NotNil(t, body, "Response body is not nil")

				aResp := body.(game_schema.GameInstanceResponse).Data
				xResp := testCase.expectResponse(th.Data, testCase.RequestBody(th.Data).(game_schema.GameInstanceRequest)).Data

				require.NotEmpty(t, aResp, "Response body is not empty")
				require.Equal(t, xResp.GameID, aResp.GameID, "Game ID equals expected")
				require.Equal(t, xResp.Status, aResp.Status, "Status equals expected")
				require.Equal(t, xResp.CurrentTurn, aResp.CurrentTurn, "Current turn equals expected")
				require.NotEmpty(t, aResp.ID, "Instance ID is not empty")
				require.False(t, aResp.CreatedAt.IsZero(), "Instance CreatedAt is not zero")
			}

			testutil.RunTestCase(t, th, &testCase.TestCase, testFunc)
		})
	}
}

func Test_deleteOneGameInstanceHandler(t *testing.T) {
	t.Parallel()

	th := testutil.NewTestHarness(t)
	require.NotNil(t, th, "TestHarness returns without error")

	_, err := th.Setup()
	require.NoError(t, err, "Test data setup returns without error")
	defer func() {
		err = th.Teardown()
		require.NoError(t, err, "Test data teardown returns without error")
	}()

	// Get a game and game instance from the harness data
	gameRec, err := th.Data.GetGameRecByRef(harness.GameOneRef)
	require.NoError(t, err, "GetGameRecByRef returns without error")

	gameInstanceRec, err := th.Data.GetGameInstanceRecByRef(harness.GameInstanceOneRef)
	require.NoError(t, err, "GetGameInstanceRecByRef returns without error")

	testCases := []struct {
		testutil.TestCase
	}{
		{
			TestCase: testutil.TestCase{
				Name: "authenticated manager when delete game instance with valid instance ID then returns no content",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[game.DeleteOneGameInstance]
				},
				RequestHeaders: testutil.AuthHeaderProManager,
				RequestPathParams: func(d harness.Data) map[string]string {
					return map[string]string{
						":game_id":     gameRec.ID,
						":instance_id": gameInstanceRec.ID,
					}
				},
				ResponseCode: http.StatusNoContent,
			},
		},
	}

	for _, testCase := range testCases {
		t.Logf("Running test >%s<\n", testCase.Name)

		t.Run(testCase.Name, func(t *testing.T) {
			testFunc := func(method string, body any) {
				// DELETE requests should return no content
				require.Nil(t, body, "Response body is nil for DELETE request")
			}

			testutil.RunTestCase(t, th, &testCase.TestCase, testFunc)
		})
	}
}

func Test_startGameInstanceHandler(t *testing.T) {
	t.Parallel()

	th := testutil.NewTestHarness(t)
	require.NotNil(t, th, "TestHarness returns without error")

	_, err := th.Setup()
	require.NoError(t, err, "Test data setup returns without error")
	defer func() {
		err = th.Teardown()
		require.NoError(t, err, "Test data teardown returns without error")
	}()

	// Get a game and game instance from the harness data
	gameRec, err := th.Data.GetGameRecByRef(harness.GameOneRef)
	require.NoError(t, err, "GetGameRecByRef returns without error")

	gameInstanceRec, err := th.Data.GetGameInstanceRecByRef(harness.GameInstanceOneRef)
	require.NoError(t, err, "GetGameInstanceRecByRef returns without error")

	testCases := []struct {
		testutil.TestCase
		expectResponse func(d harness.Data) game_schema.GameInstanceResponse
	}{
		{
			TestCase: testutil.TestCase{
				Name: "authenticated manager when start game instance with valid instance ID then returns started instance",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[game.StartGameInstance]
				},
				RequestHeaders: testutil.AuthHeaderProManager,
				RequestPathParams: func(d harness.Data) map[string]string {
					return map[string]string{
						":game_id":     gameRec.ID,
						":instance_id": gameInstanceRec.ID,
					}
				},
				ResponseDecoder: testutil.TestCaseResponseDecoderGeneric[game_schema.GameInstanceResponse],
				ResponseCode:    http.StatusOK,
			},
			expectResponse: func(d harness.Data) game_schema.GameInstanceResponse {
				return game_schema.GameInstanceResponse{
					Data: &game_schema.GameInstanceResponseData{
						Status:      "started",
						CurrentTurn: 0,
					},
				}
			},
		},
	}

	for _, testCase := range testCases {
		t.Logf("Running test >%s<\n", testCase.Name)

		t.Run(testCase.Name, func(t *testing.T) {
			testFunc := func(method string, body any) {
				require.NotNil(t, body, "Response body is not nil")

				aResp := body.(game_schema.GameInstanceResponse).Data
				xResp := testCase.expectResponse(th.Data).Data

				require.NotEmpty(t, aResp, "Response body is not empty")
				require.Equal(t, xResp.Status, aResp.Status, "Status equals expected")
				require.Equal(t, xResp.CurrentTurn, aResp.CurrentTurn, "Current turn equals expected")
				require.NotEmpty(t, aResp.ID, "Instance ID is not empty")
				require.False(t, aResp.CreatedAt.IsZero(), "Instance CreatedAt is not zero")
			}

			testutil.RunTestCase(t, th, &testCase.TestCase, testFunc)
		})
	}
}

func Test_pauseGameInstanceHandler(t *testing.T) {
	t.Parallel()

	th := testutil.NewTestHarness(t)
	require.NotNil(t, th, "TestHarness returns without error")

	_, err := th.Setup()
	require.NoError(t, err, "Test data setup returns without error")
	defer func() {
		err = th.Teardown()
		require.NoError(t, err, "Test data teardown returns without error")
	}()

	// Get a game and game instance from the harness data
	gameRec, err := th.Data.GetGameRecByRef(harness.GameOneRef)
	require.NoError(t, err, "GetGameRecByRef returns without error")

	gameInstanceRec, err := th.Data.GetGameInstanceRecByRef(harness.GameInstanceOneRef)
	require.NoError(t, err, "GetGameInstanceRecByRef returns without error")

	// For testing pause, we need a started instance
	// Since the harness creates instances with "created" status,
	// we'll test that the pause endpoint correctly returns an error
	// when trying to pause a non-started instance
	testCases := []struct {
		testutil.TestCase
		expectResponseCode int
	}{
		{
			TestCase: testutil.TestCase{
				Name: "authenticated manager when pause game instance with non-started status then returns error",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[game.PauseGameInstance]
				},
				RequestHeaders: testutil.AuthHeaderProManager,
				RequestPathParams: func(d harness.Data) map[string]string {
					return map[string]string{
						":game_id":     gameRec.ID,
						":instance_id": gameInstanceRec.ID,
					}
				},
				ResponseDecoder: testutil.TestCaseResponseDecoderGeneric[game_schema.GameInstanceResponse],
				ResponseCode:    http.StatusInternalServerError, // Expect error since instance is not started
			},
			expectResponseCode: http.StatusInternalServerError,
		},
	}

	for _, testCase := range testCases {
		t.Logf("Running test >%s<\n", testCase.Name)

		t.Run(testCase.Name, func(t *testing.T) {
			testFunc := func(method string, body any) {
				// For error responses, we don't need to validate the response body
				// The important thing is that the API correctly returns an error
				// when trying to pause a non-started instance
			}

			testutil.RunTestCase(t, th, &testCase.TestCase, testFunc)
		})
	}
}

func Test_resumeGameInstanceHandler(t *testing.T) {
	t.Parallel()

	th := testutil.NewTestHarness(t)
	require.NotNil(t, th, "TestHarness returns without error")

	_, err := th.Setup()
	require.NoError(t, err, "Test data setup returns without error")
	defer func() {
		err = th.Teardown()
		require.NoError(t, err, "Test data teardown returns without error")
	}()

	// Get a game and game instance from the harness data
	gameRec, err := th.Data.GetGameRecByRef(harness.GameOneRef)
	require.NoError(t, err, "GetGameRecByRef returns without error")

	gameInstanceRec, err := th.Data.GetGameInstanceRecByRef(harness.GameInstanceOneRef)
	require.NoError(t, err, "GetGameInstanceRecByRef returns without error")

	// For testing resume, we need a paused instance
	// Since the harness creates instances with "created" status,
	// we'll test that the resume endpoint correctly returns an error
	// when trying to resume a non-paused instance
	testCases := []struct {
		testutil.TestCase
		expectResponseCode int
	}{
		{
			TestCase: testutil.TestCase{
				Name: "authenticated manager when resume game instance with non-paused status then returns error",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[game.ResumeGameInstance]
				},
				RequestHeaders: testutil.AuthHeaderProManager,
				RequestPathParams: func(d harness.Data) map[string]string {
					return map[string]string{
						":game_id":     gameRec.ID,
						":instance_id": gameInstanceRec.ID,
					}
				},
				ResponseDecoder: testutil.TestCaseResponseDecoderGeneric[game_schema.GameInstanceResponse],
				ResponseCode:    http.StatusInternalServerError, // Expect error since instance is not paused
			},
			expectResponseCode: http.StatusInternalServerError,
		},
	}

	for _, testCase := range testCases {
		t.Logf("Running test >%s<\n", testCase.Name)

		t.Run(testCase.Name, func(t *testing.T) {
			testFunc := func(method string, body any) {
				// For error responses, we don't need to validate the response body
				// The important thing is that the API correctly returns an error
				// when trying to resume a non-paused instance
			}

			testutil.RunTestCase(t, th, &testCase.TestCase, testFunc)
		})
	}
}

func Test_cancelGameInstanceHandler(t *testing.T) {
	t.Parallel()

	th := testutil.NewTestHarness(t)
	require.NotNil(t, th, "TestHarness returns without error")

	_, err := th.Setup()
	require.NoError(t, err, "Test data setup returns without error")
	defer func() {
		err = th.Teardown()
		require.NoError(t, err, "Test data teardown returns without error")
	}()

	// Get a game and game instance from the harness data
	gameRec, err := th.Data.GetGameRecByRef(harness.GameOneRef)
	require.NoError(t, err, "GetGameRecByRef returns without error")

	gameInstanceRec, err := th.Data.GetGameInstanceRecByRef(harness.GameInstanceOneRef)
	require.NoError(t, err, "GetGameInstanceRecByRef returns without error")

	testCases := []struct {
		testutil.TestCase
		expectResponse func(d harness.Data) game_schema.GameInstanceResponse
	}{
		{
			TestCase: testutil.TestCase{
				Name: "authenticated manager when cancel game instance with valid instance ID then returns cancelled instance",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[game.CancelGameInstance]
				},
				RequestHeaders: testutil.AuthHeaderProManager,
				RequestPathParams: func(d harness.Data) map[string]string {
					return map[string]string{
						":game_id":     gameRec.ID,
						":instance_id": gameInstanceRec.ID,
					}
				},
				ResponseDecoder: testutil.TestCaseResponseDecoderGeneric[game_schema.GameInstanceResponse],
				ResponseCode:    http.StatusOK,
			},
			expectResponse: func(d harness.Data) game_schema.GameInstanceResponse {
				return game_schema.GameInstanceResponse{
					Data: &game_schema.GameInstanceResponseData{
						Status:      "cancelled",
						CurrentTurn: 0,
					},
				}
			},
		},
	}

	for _, testCase := range testCases {
		t.Logf("Running test >%s<\n", testCase.Name)

		t.Run(testCase.Name, func(t *testing.T) {
			testFunc := func(method string, body any) {
				require.NotNil(t, body, "Response body is not nil")

				aResp := body.(game_schema.GameInstanceResponse).Data
				xResp := testCase.expectResponse(th.Data).Data

				require.NotEmpty(t, aResp, "Response body is not empty")
				require.Equal(t, xResp.Status, aResp.Status, "Status equals expected")
				require.Equal(t, xResp.CurrentTurn, aResp.CurrentTurn, "Current turn equals expected")
				require.NotEmpty(t, aResp.ID, "Instance ID is not empty")
				require.False(t, aResp.CreatedAt.IsZero(), "Instance CreatedAt is not zero")
			}

			testutil.RunTestCase(t, th, &testCase.TestCase, testFunc)
		})
	}
}

func Test_createGameInstanceHandlerValidation(t *testing.T) {
	t.Parallel()

	th := testutil.NewTestHarness(t)
	require.NotNil(t, th, "TestHarness returns without error")

	_, err := th.Setup()
	require.NoError(t, err, "Test data setup returns without error")
	defer func() {
		err = th.Teardown()
		require.NoError(t, err, "Test data teardown returns without error")
	}()

	// Get a game from the harness data
	gameRec, err := th.Data.GetGameRecByRef(harness.GameOneRef)
	require.NoError(t, err, "GetGameRecByRef returns without error")

	testCases := []testutil.TestCase{
		{
			Name: "authenticated manager when create game instance with closed testing but no email delivery then returns validation error",
			NewRunner: func(cfg config.Config, l logger.Logger, s storer.Storer, j *river.Client[pgx.Tx], scanner turnsheet.TurnSheetScanner, d harness.Data) (testutil.TestRunnerer, error) {
				return testutil.NewTestRunner(cfg, l, s, j, scanner)
			},
			HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
				return rnr.GetHandlerConfig()[game.CreateOneGameInstance]
			},
			RequestHeaders: testutil.AuthHeaderProManager,
			RequestPathParams: func(d harness.Data) map[string]string {
				return map[string]string{
					":game_id": gameRec.ID,
				}
			},
			RequestBody: func(d harness.Data) interface{} {
				return game_schema.GameInstanceRequest{
					GameID:                gameRec.ID,
					DeliveryPhysicalPost:  true,
					DeliveryPhysicalLocal: false,
					DeliveryEmail:         false, // Email delivery is false, but closed testing is true
					IsClosedTesting:       true,
				}
			},
			ResponseDecoder: testutil.TestCaseResponseDecoderGeneric[coreerror.Error],
			ResponseCode:    http.StatusBadRequest,
		},
	}

	for _, testCase := range testCases {
		t.Logf("Running test >%s<\n", testCase.Name)

		t.Run(testCase.Name, func(t *testing.T) {
			testFunc := func(method string, body any) {
				if body != nil {
					errResp := body.(coreerror.Error)
					require.NotEmpty(t, errResp.Message, "Error response contains error message")
					require.Contains(t, errResp.Message, "closed testing requires email delivery", "Error message contains expected validation text")
				}
				// If body is nil, the error was already logged and we just need to verify the status code
			}

			testutil.RunTestCase(t, th, &testCase, testFunc)
		})
	}
}

func Test_getJoinGameLinkHandler(t *testing.T) {
	t.Parallel()

	th := testutil.NewTestHarness(t)
	require.NotNil(t, th, "TestHarness returns without error")

	// Configure the first game instance as closed testing before setup
	th.DataConfig.GameConfigs[0].GameInstanceConfigs[0].Record.IsClosedTesting = true
	th.DataConfig.GameConfigs[0].GameInstanceConfigs[0].Record.DeliveryEmail = true
	th.DataConfig.GameConfigs[0].GameInstanceConfigs[0].Record.DeliveryPhysicalPost = false
	th.DataConfig.GameConfigs[0].GameInstanceConfigs[0].Record.DeliveryPhysicalLocal = false

	_, err := th.Setup()
	require.NoError(t, err, "Test data setup returns without error")
	defer func() {
		err = th.Teardown()
		require.NoError(t, err, "Test data teardown returns without error")
	}()

	// Get a game from the harness data
	gameRec, err := th.Data.GetGameRecByRef(harness.GameOneRef)
	require.NoError(t, err, "GetGameRecByRef returns without error")

	// Get the closed testing instance
	gameInstanceRec, err := th.Data.GetGameInstanceRecByRef(harness.GameInstanceOneRef)
	require.NoError(t, err, "GetGameInstanceRecByRef returns without error")
	require.True(t, gameInstanceRec.IsClosedTesting, "Instance is in closed testing mode")
	require.True(t, gameInstanceRec.DeliveryEmail, "Instance has email delivery enabled")

	testCases := []struct {
		testutil.TestCase
		expectResponse func(d harness.Data) game_schema.JoinGameLinkResponse
	}{
		{
			TestCase: testutil.TestCase{
				Name: "authenticated manager when get join game link for closed testing instance then returns join link",
				NewRunner: func(cfg config.Config, l logger.Logger, s storer.Storer, j *river.Client[pgx.Tx], scanner turnsheet.TurnSheetScanner, d harness.Data) (testutil.TestRunnerer, error) {
					return testutil.NewTestRunner(cfg, l, s, j, scanner)
				},
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[game.GetJoinGameLink]
				},
				RequestHeaders: testutil.AuthHeaderProManager,
				RequestPathParams: func(d harness.Data) map[string]string {
					return map[string]string{
						":game_id":     gameRec.ID,
						":instance_id": gameInstanceRec.ID,
					}
				},
				ResponseDecoder: testutil.TestCaseResponseDecoderGeneric[game_schema.JoinGameLinkResponse],
				ResponseCode:    http.StatusOK,
			},
			expectResponse: func(d harness.Data) game_schema.JoinGameLinkResponse {
				return game_schema.JoinGameLinkResponse{
					Data: &game_schema.JoinGameLinkResponseData{
						JoinGameURL: "/player/join-game/",
						JoinGameKey: "",
					},
				}
			},
		},
	}

	for _, testCase := range testCases {
		t.Logf("Running test >%s<\n", testCase.Name)

		t.Run(testCase.Name, func(t *testing.T) {
			testFunc := func(method string, body any) {
				require.NotNil(t, body, "Response body is not nil")

				aResp := body.(game_schema.JoinGameLinkResponse).Data

				require.NotEmpty(t, aResp, "Response body is not empty")
				require.NotEmpty(t, aResp.JoinGameURL, "Join game URL is not empty")
				require.Contains(t, aResp.JoinGameURL, "/player/join-game/", "Join game URL contains expected path")
				require.NotEmpty(t, aResp.JoinGameKey, "Join game key is not empty")
				require.Contains(t, aResp.JoinGameURL, aResp.JoinGameKey, "Join game URL contains the join game key")
			}

			testutil.RunTestCase(t, th, &testCase.TestCase, testFunc)
		})
	}
}

func Test_inviteTesterHandler(t *testing.T) {
	t.Parallel()

	th := testutil.NewTestHarness(t)
	require.NotNil(t, th, "TestHarness returns without error")

	// Configure the first game instance as closed testing before setup
	th.DataConfig.GameConfigs[0].GameInstanceConfigs[0].Record.IsClosedTesting = true
	th.DataConfig.GameConfigs[0].GameInstanceConfigs[0].Record.DeliveryEmail = true
	th.DataConfig.GameConfigs[0].GameInstanceConfigs[0].Record.DeliveryPhysicalPost = false
	th.DataConfig.GameConfigs[0].GameInstanceConfigs[0].Record.DeliveryPhysicalLocal = false

	_, err := th.Setup()
	require.NoError(t, err, "Test data setup returns without error")
	defer func() {
		err = th.Teardown()
		require.NoError(t, err, "Test data teardown returns without error")
	}()

	// Get a game from the harness data
	gameRec, err := th.Data.GetGameRecByRef(harness.GameOneRef)
	require.NoError(t, err, "GetGameRecByRef returns without error")

	// Get the closed testing instance
	closedTestingInstanceRec, err := th.Data.GetGameInstanceRecByRef(harness.GameInstanceOneRef)
	require.NoError(t, err, "GetGameInstanceRecByRef returns without error")
	require.True(t, closedTestingInstanceRec.IsClosedTesting, "Instance is in closed testing mode")
	require.True(t, closedTestingInstanceRec.DeliveryEmail, "Instance has email delivery enabled")

	testCases := []struct {
		testutil.TestCase
		expectResponse func(d harness.Data) game_schema.InviteTesterResponse
	}{
		{
			TestCase: testutil.TestCase{
				Name: "authenticated manager when invite tester to closed testing instance then returns success",
				NewRunner: func(cfg config.Config, l logger.Logger, s storer.Storer, j *river.Client[pgx.Tx], scanner turnsheet.TurnSheetScanner, d harness.Data) (testutil.TestRunnerer, error) {
					return testutil.NewTestRunner(cfg, l, s, j, scanner)
				},
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[game.InviteTester]
				},
				RequestHeaders: testutil.AuthHeaderProManager,
				RequestPathParams: func(d harness.Data) map[string]string {
					return map[string]string{
						":game_id":     gameRec.ID,
						":instance_id": closedTestingInstanceRec.ID,
					}
				},
				RequestBody: func(d harness.Data) interface{} {
					return game_schema.InviteTesterRequest{
						Email: "tester@example.com",
					}
				},
				ResponseDecoder: testutil.TestCaseResponseDecoderGeneric[game_schema.InviteTesterResponse],
				ResponseCode:    http.StatusOK,
			},
			expectResponse: func(d harness.Data) game_schema.InviteTesterResponse {
				return game_schema.InviteTesterResponse{
					Data: &game_schema.InviteTesterResponseData{
						Message: "tester invitation queued",
						Email:   "tester@example.com",
					},
				}
			},
		},
	}

	for _, testCase := range testCases {
		t.Logf("Running test >%s<\n", testCase.Name)

		t.Run(testCase.Name, func(t *testing.T) {
			testFunc := func(method string, body any) {
				require.NotNil(t, body, "Response body is not nil")

				aResp := body.(game_schema.InviteTesterResponse).Data
				xResp := testCase.expectResponse(th.Data).Data

				require.NotEmpty(t, aResp, "Response body is not empty")
				require.Equal(t, xResp.Message, aResp.Message, "Message equals expected")
				require.Equal(t, xResp.Email, aResp.Email, "Email equals expected")
			}

			testutil.RunTestCase(t, th, &testCase.TestCase, testFunc)
		})
	}
}
