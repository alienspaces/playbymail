package player_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/playbymail/core/server"
	"gitlab.com/alienspaces/playbymail/internal/harness"
	"gitlab.com/alienspaces/playbymail/internal/runner/server/player"
	"gitlab.com/alienspaces/playbymail/internal/utils/testutil"
	"gitlab.com/alienspaces/playbymail/schema/api/player_schema"
)

func Test_getJoinGameInfoHandler(t *testing.T) {
	t.Parallel()

	th := testutil.NewTestHarness(t)
	require.NotNil(t, th, "NewTestHarness returns without error")

	_, err := th.Setup()
	require.NoError(t, err, "Test data setup returns without error")
	defer func() {
		err = th.Teardown()
		require.NoError(t, err, "Test data teardown returns without error")
	}()

	instanceRec, err := th.Data.GetGameInstanceRecByRef(harness.GameInstanceOneRef)
	require.NoError(t, err, "GetGameInstanceRecByRef returns without error")

	type testCase struct {
		testutil.TestCase
	}

	testCases := []testCase{
		{
			TestCase: testutil.TestCase{
				Name: "public request with valid game instance returns join game info",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[player.GetJoinGameInfo]
				},
				RequestPathParams: func(d harness.Data) map[string]string {
					return map[string]string{":game_instance_id": instanceRec.ID}
				},
				ResponseDecoder: testutil.TestCaseResponseDecoderGeneric[player_schema.JoinGameInfoResponse],
				ResponseCode:    http.StatusOK,
			},
		},
		{
			TestCase: testutil.TestCase{
				Name: "public request with non-existent game instance returns not found",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[player.GetJoinGameInfo]
				},
				RequestPathParams: func(d harness.Data) map[string]string {
					return map[string]string{":game_instance_id": "00000000-0000-0000-0000-000000000000"}
				},
				ResponseCode: http.StatusNotFound,
			},
		},
	}

	for _, tc := range testCases {
		t.Logf("Running test >%s<", tc.Name)

		t.Run(tc.Name, func(t *testing.T) {
			testFunc := func(method string, body any) {
				if tc.TestResponseCode() != http.StatusOK {
					return
				}

				require.NotNil(t, body, "Response body is not nil")
				res := body.(player_schema.JoinGameInfoResponse)
				require.NotNil(t, res.Data, "Response data is not nil")
				require.NotEmpty(t, res.Data.GameID, "Game ID is not empty")
				require.NotEmpty(t, res.Data.GameName, "Game name is not empty")
				require.NotNil(t, res.Data.Instance, "Instance is not nil")
				require.Equal(t, instanceRec.ID, res.Data.Instance.ID, "Instance ID matches")
			}
			testutil.RunTestCase(t, th, &tc, testFunc)
		})
	}
}

func Test_verifyJoinGameEmailHandler(t *testing.T) {
	t.Parallel()

	th := testutil.NewTestHarness(t)
	require.NotNil(t, th, "NewTestHarness returns without error")

	_, err := th.Setup()
	require.NoError(t, err, "Test data setup returns without error")
	defer func() {
		err = th.Teardown()
		require.NoError(t, err, "Test data teardown returns without error")
	}()

	instanceRec, err := th.Data.GetGameInstanceRecByRef(harness.GameInstanceOneRef)
	require.NoError(t, err, "GetGameInstanceRecByRef returns without error")

	// Get an account that exists in the harness
	existingAccount, err := th.Data.GetAccountUserRecByRef(harness.StandardAccountRef)
	require.NoError(t, err, "GetAccountUserRecByRef returns without error")

	type testCase struct {
		testutil.TestCase
		expectHasAccount bool
	}

	testCases := []testCase{
		{
			TestCase: testutil.TestCase{
				Name: "known email returns has_account true",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[player.VerifyJoinGameEmail]
				},
				RequestPathParams: func(d harness.Data) map[string]string {
					return map[string]string{":game_instance_id": instanceRec.ID}
				},
				RequestBody: func(d harness.Data) any {
					return player_schema.JoinGameVerifyEmailRequest{Email: existingAccount.Email}
				},
				ResponseDecoder: testutil.TestCaseResponseDecoderGeneric[player_schema.JoinGameVerifyEmailResponse],
				ResponseCode:    http.StatusOK,
			},
			expectHasAccount: true,
		},
		{
			TestCase: testutil.TestCase{
				Name: "unknown email returns has_account false",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[player.VerifyJoinGameEmail]
				},
				RequestPathParams: func(d harness.Data) map[string]string {
					return map[string]string{":game_instance_id": instanceRec.ID}
				},
				RequestBody: func(d harness.Data) any {
					return player_schema.JoinGameVerifyEmailRequest{Email: "newplayer@example.com"}
				},
				ResponseDecoder: testutil.TestCaseResponseDecoderGeneric[player_schema.JoinGameVerifyEmailResponse],
				ResponseCode:    http.StatusOK,
			},
			expectHasAccount: false,
		},
	}

	for _, tc := range testCases {
		t.Logf("Running test >%s<", tc.Name)

		t.Run(tc.Name, func(t *testing.T) {
			testFunc := func(method string, body any) {
				if tc.TestResponseCode() != http.StatusOK {
					return
				}
				require.NotNil(t, body, "Response body is not nil")
				res := body.(player_schema.JoinGameVerifyEmailResponse)
				require.NotNil(t, res.Data, "Response data is not nil")
				require.Equal(t, tc.expectHasAccount, res.Data.HasAccount, "HasAccount equals expected")
			}
			testutil.RunTestCase(t, th, &tc, testFunc)
		})
	}
}

func Test_submitJoinGameHandler(t *testing.T) {
	t.Parallel()

	th := testutil.NewTestHarness(t)
	require.NotNil(t, th, "NewTestHarness returns without error")

	_, err := th.Setup()
	require.NoError(t, err, "Test data setup returns without error")
	defer func() {
		err = th.Teardown()
		require.NoError(t, err, "Test data teardown returns without error")
	}()

	// Use GameInstanceCleanRef which has a single manager subscription (fewer conflicts)
	instanceRec, err := th.Data.GetGameInstanceRecByRef(harness.GameInstanceCleanRef)
	require.NoError(t, err, "GetGameInstanceRecByRef returns without error")

	type testCase struct {
		testutil.TestCase
		ShouldTxCommit bool
	}

	testCases := []testCase{
		{
			TestCase: testutil.TestCase{
				Name: "new player email creates account and subscription and assigns to instance",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[player.SubmitJoinGame]
				},
				RequestPathParams: func(d harness.Data) map[string]string {
					return map[string]string{":game_instance_id": instanceRec.ID}
				},
			RequestBody: func(d harness.Data) any {
				return player_schema.JoinGameSubmitRequest{
					Email:              "newjoiner@example.com",
					Name:               "New Joiner",
					PostalAddressLine1: "123 Test Street",
					StateProvince:      "VIC",
					Country:            "Australia",
					PostalCode:         "3000",
					DeliveryEmail:      true,
				}
			},
				ResponseDecoder: testutil.TestCaseResponseDecoderGeneric[player_schema.JoinGameSubmitResponse],
				ResponseCode:    http.StatusCreated,
			},
		},
		{
			TestCase: testutil.TestCase{
				Name: "non-existent game instance returns not found",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[player.SubmitJoinGame]
				},
				RequestPathParams: func(d harness.Data) map[string]string {
					return map[string]string{":game_instance_id": "00000000-0000-0000-0000-000000000000"}
				},
			RequestBody: func(d harness.Data) any {
				return player_schema.JoinGameSubmitRequest{
					Email:              "someone@example.com",
					Name:               "Someone",
					PostalAddressLine1: "456 Other Road",
					StateProvince:      "NSW",
					Country:            "Australia",
					PostalCode:         "2000",
				}
			},
				ResponseCode: http.StatusNotFound,
			},
		},
	}

	for _, tc := range testCases {
		t.Logf("Running test >%s<", tc.Name)

		t.Run(tc.Name, func(t *testing.T) {
			testFunc := func(method string, body any) {
				if tc.TestResponseCode() != http.StatusCreated {
					return
				}
				require.NotNil(t, body, "Response body is not nil")
				res := body.(player_schema.JoinGameSubmitResponse)
				require.NotNil(t, res.Data, "Response data is not nil")
				require.NotEmpty(t, res.Data.GameSubscriptionID, "GameSubscriptionID is not empty")
				require.Equal(t, instanceRec.ID, res.Data.GameInstanceID, "GameInstanceID matches")
				require.NotEmpty(t, res.Data.GameID, "GameID is not empty")

				// Verify the response body parses correctly
				bodyBytes, err := json.Marshal(res)
				require.NoError(t, err, "json.Marshal returns without error")
				require.NotEmpty(t, bodyBytes, "Body bytes are not empty")
			}
			testutil.RunTestCase(t, th, &tc, testFunc)
		})
	}
}
