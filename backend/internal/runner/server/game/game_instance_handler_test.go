package game_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/playbymail/core/server"
	"gitlab.com/alienspaces/playbymail/internal/harness"
	"gitlab.com/alienspaces/playbymail/internal/runner/server/game"
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
				Name: "API key with open access \\ search many game instances \\ returns expected instances",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[game.SearchManyGameInstances]
				},
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
				Name: "API key with open access \\ search many game instances with pagination \\ returns expected instances",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[game.SearchManyGameInstances]
				},
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
			testFunc := func(method string, body interface{}) {
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
				Name: "API key with open access \\ get many game instances \\ returns expected instances",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[game.GetManyGameInstances]
				},
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
				Name: "API key with open access \\ get many game instances with pagination \\ returns expected instances",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[game.GetManyGameInstances]
				},
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
			testFunc := func(method string, body interface{}) {
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
				Name: "API key with open access \\ get one game instance with valid instance ID \\ returns expected instance",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[game.GetOneGameInstance]
				},
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
			testFunc := func(method string, body interface{}) {
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
				Name: "API key with open access \\ create game instance with valid properties \\ returns created instance",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[game.CreateOneGameInstance]
				},
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
			testFunc := func(method string, body interface{}) {
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
				Name: "API key with open access \\ update game instance with valid properties \\ returns updated instance",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[game.UpdateOneGameInstance]
				},
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
			testFunc := func(method string, body interface{}) {
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
				Name: "API key with open access \\ delete game instance with valid instance ID \\ returns no content",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[game.DeleteOneGameInstance]
				},
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
			testFunc := func(method string, body interface{}) {
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
				Name: "API key with open access \\ start game instance with valid instance ID \\ returns started instance",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[game.StartGameInstance]
				},
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
			testFunc := func(method string, body interface{}) {
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
				Name: "API key with open access \\ pause game instance with non-started status \\ returns error",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[game.PauseGameInstance]
				},
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
			testFunc := func(method string, body interface{}) {
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
				Name: "API key with open access \\ resume game instance with non-paused status \\ returns error",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[game.ResumeGameInstance]
				},
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
			testFunc := func(method string, body interface{}) {
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
				Name: "API key with open access \\ cancel game instance with valid instance ID \\ returns cancelled instance",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[game.CancelGameInstance]
				},
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
			testFunc := func(method string, body interface{}) {
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
