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

func TestGetManyGameParameterConfigurationsHandler(t *testing.T) {
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
	}

	testCaseResponseDecoder := testutil.TestCaseResponseDecoderGeneric[game_schema.GameParameterConfigurationCollectionResponse]

	testCases := []testCase{
		{
			TestCase: testutil.TestCase{
				Name: "API key with open access \\ get all configurations \\ returns expected configurations",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[game.GetManyGameParameterConfigurations]
				},
				RequestQueryParams: func(d harness.Data) map[string]any {
					return map[string]any{}
				},
				ResponseDecoder: testCaseResponseDecoder,
				ResponseCode:    http.StatusOK,
			},
		},
		{
			TestCase: testutil.TestCase{
				Name: "API key with open access \\ filter by adventure game type \\ returns adventure configurations",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[game.GetManyGameParameterConfigurations]
				},
				RequestQueryParams: func(d harness.Data) map[string]any {
					return map[string]any{
						"game_type": "adventure",
					}
				},
				ResponseDecoder: testCaseResponseDecoder,
				ResponseCode:    http.StatusOK,
			},
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

				response := body.(game_schema.GameParameterConfigurationCollectionResponse)
				require.GreaterOrEqual(t, len(response.Data), 1, "Should have at least one configuration")

				// Verify structure
				if len(response.Data) > 0 {
					config := response.Data[0]
					require.NotEmpty(t, config.GameType, "GameType should not be empty")
					require.NotEmpty(t, config.ConfigKey, "ConfigKey should not be empty")
					require.NotEmpty(t, config.ValueType, "ValueType should not be empty")
				}
			}

			testutil.RunTestCase(t, th, &testCase.TestCase, testFunc)
		})
	}
}
