package game_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/playbymail/core/collection/set"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/core/server"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/harness"
	game "gitlab.com/alienspaces/playbymail/internal/runner/server/game"
	"gitlab.com/alienspaces/playbymail/internal/utils/testutil"
	"gitlab.com/alienspaces/playbymail/schema/api/game_schema"
)

func Test_getManyGameInstanceParametersHandler(t *testing.T) {
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
		expectResponse func(d harness.Data, req game_schema.GameInstanceParameterRequest) game_schema.GameInstanceParameterCollectionResponse
	}{
		{
			TestCase: testutil.TestCase{
				Name: "authenticated manager when get many game instance parameters then returns expected parameters",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[game.GetManyGameInstanceParameters]
				},
				RequestHeaders: func(d harness.Data) map[string]string {
					return testutil.AuthHeaderProManager(d)
				},
				RequestPathParams: func(d harness.Data) map[string]string {
					return map[string]string{
						":game_id":     gameRec.ID,
						":instance_id": gameInstanceRec.ID,
					}
				},
				ResponseDecoder: testutil.TestCaseResponseDecoderGeneric[game_schema.GameInstanceParameterCollectionResponse],
				ResponseCode:    http.StatusOK,
			},
			expectResponse: func(d harness.Data, req game_schema.GameInstanceParameterRequest) game_schema.GameInstanceParameterCollectionResponse {
				return game_schema.GameInstanceParameterCollectionResponse{
					Data: []*game_schema.GameInstanceParameter{
						{
							ParameterKey:   domain.AdventureGameParameterCharacterLives,
							ParameterValue: "3",
						},
					},
				}
			},
		},
		{
			TestCase: testutil.TestCase{
				Name: "authenticated manager when get many game instance parameters with pagination then returns expected parameters",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[game.GetManyGameInstanceParameters]
				},
				RequestHeaders: func(d harness.Data) map[string]string {
					return testutil.AuthHeaderProManager(d)
				},
				RequestPathParams: func(d harness.Data) map[string]string {
					return map[string]string{
						":game_id":     gameRec.ID,
						":instance_id": gameInstanceRec.ID,
					}
				},
				RequestQueryParams: func(d harness.Data) map[string]any {
					return map[string]any{
						"page_size":   "10",
						"page_number": "1",
					}
				},
				ResponseDecoder: testutil.TestCaseResponseDecoderGeneric[game_schema.GameInstanceParameterCollectionResponse],
				ResponseCode:    http.StatusOK,
			},
			expectResponse: func(d harness.Data, req game_schema.GameInstanceParameterRequest) game_schema.GameInstanceParameterCollectionResponse {
				return game_schema.GameInstanceParameterCollectionResponse{
					Data: []*game_schema.GameInstanceParameter{
						{
							ParameterKey:   domain.AdventureGameParameterCharacterLives,
							ParameterValue: "3",
						},
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

				aResp := body.(game_schema.GameInstanceParameterCollectionResponse).Data
				xResp := testCase.expectResponse(
					th.Data,
					game_schema.GameInstanceParameterRequest{}, // Not used for GET requests
				).Data

				require.NotEmpty(t, aResp, "Response body is not empty")
				require.Len(t, aResp, 1, "Response contains exactly one parameter")
				require.Equal(t, xResp[0].ParameterKey, aResp[0].ParameterKey, "Parameter key equals expected")
				require.Equal(t, xResp[0].ParameterValue, aResp[0].ParameterValue, "Parameter value equals expected")
				require.NotEmpty(t, aResp[0].ID, "Parameter ID is not empty")
				require.False(t, aResp[0].CreatedAt.IsZero(), "Parameter CreatedAt is not zero")
			}

			testutil.RunTestCase(t, th, &testCase.TestCase, testFunc)
		})
	}
}

func Test_getOneGameInstanceParameterHandler(t *testing.T) {
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

	// Get the existing parameter from the harness data
	existingParam, err := th.Data.GetGameInstanceParameterRecByRef(harness.GameInstanceParameterOneRef)
	require.NoError(t, err, "GetGameInstanceParameterRecByRef returns without error")

	testCases := []struct {
		testutil.TestCase
		expectResponse func(d harness.Data, req game_schema.GameInstanceParameterRequest) game_schema.GameInstanceParameterResponse
	}{
		{
			TestCase: testutil.TestCase{
				Name: "authenticated manager when get one game instance parameter with valid parameter ID then returns expected parameter",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[game.GetOneGameInstanceParameter]
				},
				RequestHeaders: func(d harness.Data) map[string]string {
					return testutil.AuthHeaderProManager(d)
				},
				RequestPathParams: func(d harness.Data) map[string]string {
					return map[string]string{
						":game_id":      gameRec.ID,
						":instance_id":  gameInstanceRec.ID,
						":parameter_id": existingParam.ID,
					}
				},
				ResponseDecoder: testutil.TestCaseResponseDecoderGeneric[game_schema.GameInstanceParameterResponse],
				ResponseCode:    http.StatusOK,
			},
			expectResponse: func(d harness.Data, req game_schema.GameInstanceParameterRequest) game_schema.GameInstanceParameterResponse {
				return game_schema.GameInstanceParameterResponse{
					Data: &game_schema.GameInstanceParameter{
						ParameterKey:   domain.AdventureGameParameterCharacterLives,
						ParameterValue: "3",
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

				aResp := body.(game_schema.GameInstanceParameterResponse).Data
				xResp := testCase.expectResponse(
					th.Data,
					game_schema.GameInstanceParameterRequest{}, // Not used for GET requests
				).Data

				require.NotEmpty(t, aResp, "Response body is not empty")
				require.Equal(t, xResp.ParameterKey, aResp.ParameterKey, "Parameter key equals expected")
				require.Equal(t, xResp.ParameterValue, aResp.ParameterValue, "Parameter value equals expected")
				require.NotEmpty(t, aResp.ID, "Parameter ID is not empty")
				require.False(t, aResp.CreatedAt.IsZero(), "Parameter CreatedAt is not zero")
			}

			testutil.RunTestCase(t, th, &testCase.TestCase, testFunc)
		})
	}
}

func Test_createUpdateDeleteGameInstanceParameterHandler(t *testing.T) {
	t.Parallel()

	th := testutil.NewTestHarness(t)
	require.NotNil(t, th, "TestHarness returns without error")

	_, err := th.Setup()
	require.NoError(t, err, "Test data setup returns without error")
	defer func() {
		err = th.Teardown()
		require.NoError(t, err, "Test data teardown returns without error")
	}()

	// Get a game and clean game instance from the harness data
	gameRec, err := th.Data.GetGameRecByRef(harness.GameOneRef)
	require.NoError(t, err, "GetGameRecByRef returns without error")

	gameInstanceRec, err := th.Data.GetGameInstanceRecByRef(harness.GameInstanceCleanRef)
	require.NoError(t, err, "GetGameInstanceRecByRef returns without error")

	testCases := []struct {
		testutil.TestCase
		expectResponse func(d harness.Data, req game_schema.GameInstanceParameterRequest) game_schema.GameInstanceParameterResponse
	}{
		{
			TestCase: testutil.TestCase{
				Name: "authenticated manager when create game instance parameter with character_lives value 5 then returns created parameter",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[game.CreateOneGameInstanceParameter]
				},
				RequestHeaders: func(d harness.Data) map[string]string {
					return testutil.AuthHeaderProManager(d)
				},
				RequestPathParams: func(d harness.Data) map[string]string {
					return map[string]string{
						":game_id":     gameRec.ID,
						":instance_id": gameInstanceRec.ID,
					}
				},
				RequestBody: func(d harness.Data) any {
					return game_schema.GameInstanceParameterRequest{
						ParameterKey:   domain.AdventureGameParameterCharacterLives, // Use actual supported parameter
						ParameterValue: "5",
					}
				},
				ResponseDecoder: testutil.TestCaseResponseDecoderGeneric[game_schema.GameInstanceParameterResponse],
				ResponseCode:    http.StatusCreated,
			},
			expectResponse: func(d harness.Data, req game_schema.GameInstanceParameterRequest) game_schema.GameInstanceParameterResponse {
				return game_schema.GameInstanceParameterResponse{
					Data: &game_schema.GameInstanceParameter{
						ParameterKey:   req.ParameterKey,
						ParameterValue: fmt.Sprintf("%v", req.ParameterValue), // Convert any type to string
					},
				}
			},
		},
		{
			TestCase: testutil.TestCase{
				Name: "authenticated manager when create game instance parameter with character_lives value 3 then returns created parameter",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[game.CreateOneGameInstanceParameter]
				},
				RequestHeaders: func(d harness.Data) map[string]string {
					return testutil.AuthHeaderProManager(d)
				},
				RequestPathParams: func(d harness.Data) map[string]string {
					return map[string]string{
						":game_id":     gameRec.ID,
						":instance_id": gameInstanceRec.ID,
					}
				},
				RequestBody: func(d harness.Data) any {
					return game_schema.GameInstanceParameterRequest{
						ParameterKey:   domain.AdventureGameParameterCharacterLives, // Use actual supported parameter
						ParameterValue: 3,                                           // Integer value
					}
				},
				ResponseDecoder: testutil.TestCaseResponseDecoderGeneric[game_schema.GameInstanceParameterResponse],
				ResponseCode:    http.StatusCreated,
			},
			expectResponse: func(d harness.Data, req game_schema.GameInstanceParameterRequest) game_schema.GameInstanceParameterResponse {
				return game_schema.GameInstanceParameterResponse{
					Data: &game_schema.GameInstanceParameter{
						ParameterKey:   req.ParameterKey,
						ParameterValue: fmt.Sprintf("%v", req.ParameterValue), // Convert any type to string
					},
				}
			},
		},
		{
			TestCase: testutil.TestCase{
				Name: "authenticated manager when create game instance parameter with character_lives value 5 then returns created parameter",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[game.CreateOneGameInstanceParameter]
				},
				RequestHeaders: func(d harness.Data) map[string]string {
					return testutil.AuthHeaderProManager(d)
				},
				RequestPathParams: func(d harness.Data) map[string]string {
					return map[string]string{
						":game_id":     gameRec.ID,
						":instance_id": gameInstanceRec.ID,
					}
				},
				RequestBody: func(d harness.Data) any {
					return game_schema.GameInstanceParameterRequest{
						ParameterKey:   domain.AdventureGameParameterCharacterLives, // Use actual supported parameter
						ParameterValue: 5,                                           // Integer value
					}
				},
				ResponseDecoder: testutil.TestCaseResponseDecoderGeneric[game_schema.GameInstanceParameterResponse],
				ResponseCode:    http.StatusCreated,
			},
			expectResponse: func(d harness.Data, req game_schema.GameInstanceParameterRequest) game_schema.GameInstanceParameterResponse {
				return game_schema.GameInstanceParameterResponse{
					Data: &game_schema.GameInstanceParameter{
						ParameterKey:   req.ParameterKey,
						ParameterValue: fmt.Sprintf("%v", req.ParameterValue), // Convert any type to string
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

				aResp := body.(game_schema.GameInstanceParameterResponse).Data
				xResp := testCase.expectResponse(
					th.Data,
					testCase.TestRequestBody(th.Data).(game_schema.GameInstanceParameterRequest),
				).Data

				require.NotEmpty(t, aResp, "Response body is not empty")
				require.Equal(t, xResp.ParameterKey, aResp.ParameterKey, "Parameter key equals expected")
				require.Equal(t, xResp.ParameterValue, aResp.ParameterValue, "Parameter value equals expected")
				require.NotEmpty(t, aResp.ID, "Parameter ID is not empty")
				require.False(t, aResp.CreatedAt.IsZero(), "Parameter CreatedAt is not zero")
			}

			testutil.RunTestCase(t, th, &testCase.TestCase, testFunc)
		})
	}
}

func Test_updateGameInstanceParameterHandler(t *testing.T) {
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

	// Get the existing parameter from the harness data
	existingParam, err := th.Data.GetGameInstanceParameterRecByRef(harness.GameInstanceParameterOneRef)
	require.NoError(t, err, "GetGameInstanceParameterRecByRef returns without error")

	testCases := []struct {
		testutil.TestCase
		expectResponse func(d harness.Data, req game_schema.GameInstanceParameterRequest) game_schema.GameInstanceParameterResponse
	}{
		{
			TestCase: testutil.TestCase{
				Name: "authenticated manager when update game instance parameter with valid properties then returns updated parameter",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[game.UpdateOneGameInstanceParameter]
				},
				RequestHeaders: func(d harness.Data) map[string]string {
					return testutil.AuthHeaderProManager(d)
				},
				RequestPathParams: func(d harness.Data) map[string]string {
					return map[string]string{
						":game_id":      gameRec.ID,
						":instance_id":  gameInstanceRec.ID,
						":parameter_id": existingParam.ID,
					}
				},
				RequestBody: func(d harness.Data) any {
					return game_schema.GameInstanceParameterRequest{
						ParameterKey:   domain.AdventureGameParameterCharacterLives,
						ParameterValue: "5", // Update from 3 to 5
					}
				},
				ResponseDecoder: testutil.TestCaseResponseDecoderGeneric[game_schema.GameInstanceParameterResponse],
				ResponseCode:    http.StatusOK,
			},
			expectResponse: func(d harness.Data, req game_schema.GameInstanceParameterRequest) game_schema.GameInstanceParameterResponse {
				return game_schema.GameInstanceParameterResponse{
					Data: &game_schema.GameInstanceParameter{
						ParameterKey:   req.ParameterKey,
						ParameterValue: fmt.Sprintf("%v", req.ParameterValue), // Convert any type to string
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

				aResp := body.(game_schema.GameInstanceParameterResponse).Data
				xResp := testCase.expectResponse(
					th.Data,
					testCase.TestRequestBody(th.Data).(game_schema.GameInstanceParameterRequest),
				).Data

				require.NotEmpty(t, aResp, "Response body is not empty")
				require.Equal(t, xResp.ParameterKey, aResp.ParameterKey, "Parameter key equals expected")
				require.Equal(t, xResp.ParameterValue, aResp.ParameterValue, "Parameter value equals expected")
				require.NotEmpty(t, aResp.ID, "Parameter ID is not empty")
				require.False(t, aResp.CreatedAt.IsZero(), "Parameter CreatedAt is not zero")
			}

			testutil.RunTestCase(t, th, &testCase.TestCase, testFunc)
		})
	}
}

func Test_deleteGameInstanceParameterHandler(t *testing.T) {
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

	// Get the existing parameter from the harness data
	existingParam, err := th.Data.GetGameInstanceParameterRecByRef(harness.GameInstanceParameterOneRef)
	require.NoError(t, err, "GetGameInstanceParameterRecByRef returns without error")

	testCases := []struct {
		testutil.TestCase
	}{
		{
			TestCase: testutil.TestCase{
				Name: "authenticated manager when delete game instance parameter with valid parameter ID then returns no content",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[game.DeleteOneGameInstanceParameter]
				},
				RequestHeaders: func(d harness.Data) map[string]string {
					return testutil.AuthHeaderProManager(d)
				},
				RequestPathParams: func(d harness.Data) map[string]string {
					return map[string]string{
						":game_id":      gameRec.ID,
						":instance_id":  gameInstanceRec.ID,
						":parameter_id": existingParam.ID,
					}
				},
				ResponseCode: http.StatusNoContent,
			},
		},
	}

	for _, testCase := range testCases {
		t.Logf("Running test >%s<\n", testCase.Name)

		t.Run(testCase.Name, func(t *testing.T) {
			testutil.RunTestCase(t, th, &testCase.TestCase, nil)
		})
	}
}

func Test_gameInstanceParameterHandlerValidation(t *testing.T) {
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
		ResponseError func(d harness.Data) testutil.ExpectedErrorResponse
	}{
		{
			TestCase: testutil.TestCase{
				Name: "authenticated manager when create game instance parameter with missing parameter_key then returns validation error",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[game.CreateOneGameInstanceParameter]
				},
				RequestHeaders: func(d harness.Data) map[string]string {
					return testutil.AuthHeaderProManager(d)
				},
				RequestPathParams: func(d harness.Data) map[string]string {
					return map[string]string{
						":game_id":     gameRec.ID,
						":instance_id": gameInstanceRec.ID,
					}
				},
				RequestBody: func(d harness.Data) any {
					return game_schema.GameInstanceParameterRequest{
						ParameterKey:   "", // Missing parameter key
						ParameterValue: "test_value",
					}
				},
				ResponseDecoder: testutil.TestCaseResponseDecoderGeneric[coreerror.Error],
				ResponseCode:    http.StatusBadRequest,
			},
			ResponseError: func(d harness.Data) testutil.ExpectedErrorResponse {
				return testutil.ExpectedErrorResponse{
					Errs: set.FromSlice([]coreerror.Code{"invalid_parameter_key"}),
				}
			},
		},
		{
			TestCase: testutil.TestCase{
				Name: "authenticated manager when create game instance parameter with empty parameter_key then returns validation error",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[game.CreateOneGameInstanceParameter]
				},
				RequestHeaders: func(d harness.Data) map[string]string {
					return testutil.AuthHeaderProManager(d)
				},
				RequestPathParams: func(d harness.Data) map[string]string {
					return map[string]string{
						":game_id":     gameRec.ID,
						":instance_id": gameInstanceRec.ID,
					}
				},
				RequestBody: func(d harness.Data) any {
					return game_schema.GameInstanceParameterRequest{
						ParameterKey:   "", // Empty parameter key
						ParameterValue: "test_value",
					}
				},
				ResponseDecoder: testutil.TestCaseResponseDecoderGeneric[coreerror.Error],
				ResponseCode:    http.StatusBadRequest,
			},
			ResponseError: func(d harness.Data) testutil.ExpectedErrorResponse {
				return testutil.ExpectedErrorResponse{
					Errs: set.FromSlice([]coreerror.Code{"invalid_parameter_key"}),
				}
			},
		},
		{
			TestCase: testutil.TestCase{
				Name: "authenticated manager when create game instance parameter with missing parameter_value then returns validation error",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[game.CreateOneGameInstanceParameter]
				},
				RequestHeaders: func(d harness.Data) map[string]string {
					return testutil.AuthHeaderProManager(d)
				},
				RequestPathParams: func(d harness.Data) map[string]string {
					return map[string]string{
						":game_id":     gameRec.ID,
						":instance_id": gameInstanceRec.ID,
					}
				},
				RequestBody: func(d harness.Data) any {
					return game_schema.GameInstanceParameterRequest{
						ParameterKey:   domain.AdventureGameParameterCharacterLives, // Valid parameter key
						ParameterValue: "",                                          // Missing parameter value
					}
				},
				ResponseDecoder: testutil.TestCaseResponseDecoderGeneric[coreerror.Error],
				ResponseCode:    http.StatusBadRequest,
			},
			ResponseError: func(d harness.Data) testutil.ExpectedErrorResponse {
				return testutil.ExpectedErrorResponse{
					Errs: set.FromSlice([]coreerror.Code{"invalid_parameter_value"}),
				}
			},
		},
	}

	for _, testCase := range testCases {
		t.Logf("Running test >%s<\n", testCase.Name)

		t.Run(testCase.Name, func(t *testing.T) {
			testutil.RunTestCase(t, th, &testCase.TestCase, nil)
		})
	}
}
