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
	"gitlab.com/alienspaces/playbymail/internal/runner/server/game"
	"gitlab.com/alienspaces/playbymail/internal/turnsheet"
	"gitlab.com/alienspaces/playbymail/internal/utils/config"
	"gitlab.com/alienspaces/playbymail/internal/utils/testutil"
	"gitlab.com/alienspaces/playbymail/schema/api/player_schema"
)

func Test_getJoinInfoHandler(t *testing.T) {
	t.Parallel()

	th := testutil.NewTestHarness(t)
	require.NotNil(t, th)

	_, err := th.Setup()
	require.NoError(t, err)
	defer func() {
		err = th.Teardown()
		require.NoError(t, err)
	}()

	testCases := []testutil.TestCase{
		{
			Name: "public request returns join info for active manager subscription",
			NewRunner: func(cfg config.Config, l logger.Logger, s storer.Storer, j *river.Client[pgx.Tx], scanner turnsheet.TurnSheetScanner, d harness.Data) (testutil.TestRunnerer, error) {
				return testutil.NewTestRunner(cfg, l, s, j, scanner)
			},
			HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
				return rnr.GetHandlerConfig()[game.GetJoinInfo]
			},
			RequestPathParams: func(d harness.Data) map[string]string {
				return map[string]string{
					":game_subscription_id": managerSubscriptionID(t, d),
				}
			},
			ResponseDecoder: testutil.TestCaseResponseDecoderGeneric[player_schema.JoinGameInfoResponse],
			ResponseCode:    http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			testFunc := func(method string, body any) {
				if tc.TestResponseCode() != http.StatusOK {
					return
				}
				require.NotNil(t, body)
				res := body.(player_schema.JoinGameInfoResponse)
				require.NotNil(t, res.Data)
				require.NotEmpty(t, res.Data.GameSubscriptionID)
				require.NotEmpty(t, res.Data.GameName)
				require.NotEmpty(t, res.Data.GameType)
			}
			testutil.RunTestCase(t, th, &tc, testFunc)
		})
	}
}

func Test_submitJoinHandler(t *testing.T) {
	t.Parallel()

	th := testutil.NewTestHarness(t)
	require.NotNil(t, th)

	_, err := th.Setup()
	require.NoError(t, err)
	defer func() {
		err = th.Teardown()
		require.NoError(t, err)
	}()

	type testCase struct {
		testutil.TestCase
		expectInstanceID bool
		expectStatus     string
	}

	testCases := []testCase{
		{
			TestCase: testutil.TestCase{
				Name: "authenticated request creates active subscription with instance",
				NewRunner: func(cfg config.Config, l logger.Logger, s storer.Storer, j *river.Client[pgx.Tx], scanner turnsheet.TurnSheetScanner, d harness.Data) (testutil.TestRunnerer, error) {
					return testutil.NewTestRunner(cfg, l, s, j, scanner)
				},
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[game.SubmitJoin]
				},
				RequestHeaders: testutil.AuthHeaderProManager,
				RequestPathParams: func(d harness.Data) map[string]string {
					return map[string]string{
						":game_subscription_id": managerSubscriptionID(t, d),
					}
				},
				RequestBody: func(d harness.Data) interface{} {
					return player_schema.JoinGameSubmitRequest{
						Email:              "newplayer@example.com",
						Name:               "New Player",
						DeliveryMethod:     "post",
						PostalAddressLine1: "123 Main St",
						StateProvince:      "CA",
						Country:            "US",
						PostalCode:         "90210",
					}
				},
				ResponseDecoder: testutil.TestCaseResponseDecoderGeneric[player_schema.JoinGameSubmitResponse],
				ResponseCode:    http.StatusCreated,
			},
			expectInstanceID: true,
			expectStatus:     "active",
		},
		{
			TestCase: testutil.TestCase{
				Name: "unauthenticated request creates pending_approval subscription with instance",
				NewRunner: func(cfg config.Config, l logger.Logger, s storer.Storer, j *river.Client[pgx.Tx], scanner turnsheet.TurnSheetScanner, d harness.Data) (testutil.TestRunnerer, error) {
					return testutil.NewTestRunner(cfg, l, s, j, scanner)
				},
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[game.SubmitJoin]
				},
				RequestPathParams: func(d harness.Data) map[string]string {
					return map[string]string{
						":game_subscription_id": managerSubscriptionID(t, d),
					}
				},
				RequestBody: func(d harness.Data) interface{} {
					return player_schema.JoinGameSubmitRequest{
						Email:          "newplayer2@example.com",
						Name:           "New Player Two",
						DeliveryMethod: "email",
					}
				},
				ResponseDecoder: testutil.TestCaseResponseDecoderGeneric[player_schema.JoinGameSubmitResponse],
				ResponseCode:    http.StatusCreated,
			},
			expectInstanceID: true,
			expectStatus:     "pending_approval",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			tc := tc
			testFunc := func(method string, body any) {
				if tc.TestResponseCode() != http.StatusCreated {
					return
				}
				require.NotNil(t, body)
				res := body.(player_schema.JoinGameSubmitResponse)
				require.NotNil(t, res.Data)
				require.NotEmpty(t, res.Data.GameSubscriptionID)
				require.NotEmpty(t, res.Data.GameID)
				require.Equal(t, tc.expectStatus, res.Data.Status)
				if tc.expectInstanceID {
					require.NotEmpty(t, res.Data.GameInstanceID)
				} else {
					require.Empty(t, res.Data.GameInstanceID)
				}
			}
			testutil.RunTestCase(t, th, &tc.TestCase, testFunc)
		})
	}
}

func Test_submitJoinHandler_duplicateJoin(t *testing.T) {
	t.Parallel()

	th := testutil.NewTestHarness(t)
	require.NotNil(t, th)

	_, err := th.Setup()
	require.NoError(t, err)
	defer func() {
		err = th.Teardown()
		require.NoError(t, err)
	}()

	// AccountUserStandardRef already has GameSubscriptionPlayerOneRef (a player subscription
	// for GameOneRef) in the harness data. Attempting to join again via the manager link
	// as that user should be rejected by the duplicate-join guard without needing to commit
	// any additional state.
	tc := testutil.TestCase{
		Name: "authenticated user with existing player subscription is rejected",
		NewRunner: func(cfg config.Config, l logger.Logger, s storer.Storer, j *river.Client[pgx.Tx], scanner turnsheet.TurnSheetScanner, d harness.Data) (testutil.TestRunnerer, error) {
			return testutil.NewTestRunner(cfg, l, s, j, scanner)
		},
		HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
			return rnr.GetHandlerConfig()[game.SubmitJoin]
		},
		RequestHeaders: testutil.AuthHeaderStandard,
		RequestPathParams: func(d harness.Data) map[string]string {
			return map[string]string{":game_subscription_id": managerSubscriptionID(t, d)}
		},
		RequestBody: func(d harness.Data) interface{} {
			return player_schema.JoinGameSubmitRequest{
				Email:          "standard@example.com",
				Name:           "Standard User",
				CharacterName:  "Standard Character",
				DeliveryMethod: "email",
			}
		},
		ResponseDecoder: testutil.TestCaseResponseDecoderGeneric[player_schema.JoinGameSubmitResponse],
		ResponseCode:    http.StatusBadRequest,
	}
	testutil.RunTestCase(t, th, &tc, func(method string, body any) {})
}
