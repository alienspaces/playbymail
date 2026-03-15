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

func managerSubscriptionID(t *testing.T, d harness.Data) string {
	t.Helper()
	subID, ok := d.Refs.GameSubscriptionRefs[harness.GameSubscriptionManagerOneRef]
	require.True(t, ok, "manager subscription ref exists")
	return subID
}

func Test_inviteHandler(t *testing.T) {
	t.Parallel()

	th := testutil.NewTestHarness(t)
	require.NotNil(t, th, "NewTestHarness returns without error")

	_, err := th.Setup()
	require.NoError(t, err, "Test data setup returns without error")
	defer func() {
		err = th.Teardown()
		require.NoError(t, err, "Test data teardown returns without error")
	}()

	testCases := []struct {
		testutil.TestCase
		expectResponse func(d harness.Data) game_schema.InviteResponse
	}{
		{
			TestCase: testutil.TestCase{
				Name: "authenticated manager invites player via subscription queues invitation",
				NewRunner: func(cfg config.Config, l logger.Logger, s storer.Storer, j *river.Client[pgx.Tx], scanner turnsheet.TurnSheetScanner, d harness.Data) (testutil.TestRunnerer, error) {
					return testutil.NewTestRunner(cfg, l, s, j, scanner)
				},
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[game.Invite]
				},
				RequestHeaders: testutil.AuthHeaderProManager,
				RequestPathParams: func(d harness.Data) map[string]string {
					return map[string]string{
						":game_subscription_id": managerSubscriptionID(t, d),
					}
				},
				RequestBody: func(d harness.Data) interface{} {
					return game_schema.InviteRequest{
						Email: "player@example.com",
					}
				},
				ResponseDecoder: testutil.TestCaseResponseDecoderGeneric[game_schema.InviteResponse],
				ResponseCode:    http.StatusOK,
			},
			expectResponse: func(d harness.Data) game_schema.InviteResponse {
				return game_schema.InviteResponse{
					Data: &game_schema.InviteResponseData{
						Message: "invitation queued",
						Email:   "player@example.com",
					},
				}
			},
		},
		{
			TestCase: testutil.TestCase{
				Name: "unauthenticated request returns unauthorized",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[game.Invite]
				},
				RequestPathParams: func(d harness.Data) map[string]string {
					return map[string]string{
						":game_subscription_id": managerSubscriptionID(t, d),
					}
				},
				RequestBody: func(d harness.Data) interface{} {
					return game_schema.InviteRequest{
						Email: "player@example.com",
					}
				},
				ResponseCode: http.StatusUnauthorized,
			},
		},
	}

	for _, testCase := range testCases {
		t.Logf("Running test >%s<\n", testCase.Name)

		t.Run(testCase.Name, func(t *testing.T) {
			testFunc := func(method string, body any) {
				if testCase.TestResponseCode() != http.StatusOK {
					return
				}
				require.NotNil(t, body, "Response body is not nil")

				aResp := body.(game_schema.InviteResponse).Data
				xResp := testCase.expectResponse(th.Data).Data

				require.NotEmpty(t, aResp, "Response body is not empty")
				require.Equal(t, xResp.Message, aResp.Message, "Message equals expected")
				require.Equal(t, xResp.Email, aResp.Email, "Email equals expected")
			}

			testutil.RunTestCase(t, th, &testCase.TestCase, testFunc)
		})
	}
}

func Test_gameSubscriptionHandler(t *testing.T) {
	t.Parallel()

	th := testutil.NewTestHarness(t)
	require.NotNil(t, th, "newTestHarness returns without error")

	_, err := th.Setup()
	require.NoError(t, err, "Test data setup returns without error")

	defer func() {
		err = th.Teardown()
		require.NoError(t, err, "Test data teardown returns without error")
	}()

	collectionDecoder := testutil.TestCaseResponseDecoderGeneric[game_schema.GameSubscriptionCollectionResponse]
	singleDecoder := testutil.TestCaseResponseDecoderGeneric[game_schema.GameSubscriptionResponse]

	testCases := []testutil.TestCase{
		{
			Name: "GET many game subscriptions",
			HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
				return rnr.GetHandlerConfig()[game.GetManyGameSubscriptions]
			},
			RequestHeaders:  testutil.AuthHeaderProManager,
			ResponseDecoder: collectionDecoder,
			ResponseCode:    http.StatusOK,
		},
		{
			Name: "POST create game subscription",
			HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
				return rnr.GetHandlerConfig()[game.CreateOneGameSubscription]
			},
			RequestHeaders: testutil.AuthHeaderProManager,
			RequestBody: func(d harness.Data) any {
				gameRec, _ := d.GetGameRecByRef(harness.GameOneRef)
				accountRec, _ := d.GetAccountUserRecByRef(harness.AccountUserProManagerRef)
				accountUserContactRec, _ := d.GetAccountUserContactRecByAccountUserID(accountRec.ID)
				return game_schema.GameSubscriptionRequest{
					GameID:               gameRec.ID,
					AccountID:            accountRec.AccountID,
					AccountUserContactID: &accountUserContactRec.ID,
					SubscriptionType:     game_record.GameSubscriptionTypePlayer,
				}
			},
			ResponseDecoder: singleDecoder,
			ResponseCode:    http.StatusCreated,
		},
		{
			Name: "GET one game subscription",
			HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
				return rnr.GetHandlerConfig()[game.GetOneGameSubscription]
			},
			RequestHeaders: testutil.AuthHeaderProManager,
			RequestPathParams: func(d harness.Data) map[string]string {
				subID, ok := d.Refs.GameSubscriptionRefs[harness.GameSubscriptionManagerOneRef]
				require.True(t, ok, "manager subscription ref exists")
				return map[string]string{":game_subscription_id": subID}
			},
			ResponseDecoder: singleDecoder,
			ResponseCode:    http.StatusOK,
		},
		{
			Name: "PUT update game subscription",
			HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
				return rnr.GetHandlerConfig()[game.UpdateOneGameSubscription]
			},
			RequestHeaders: testutil.AuthHeaderProManager,
			RequestPathParams: func(d harness.Data) map[string]string {
				subID, ok := d.Refs.GameSubscriptionRefs[harness.GameSubscriptionManagerOneRef]
				require.True(t, ok, "manager subscription ref exists")
				return map[string]string{":game_subscription_id": subID}
			},
			RequestBody: func(d harness.Data) any {
				sub, _ := d.GetGameSubscriptionRecByRef(harness.GameSubscriptionManagerOneRef)
				return game_schema.GameSubscriptionRequest{
					GameID:           sub.GameID,
					AccountID:        sub.AccountID,
					SubscriptionType: game_record.GameSubscriptionTypeManager,
				}
			},
			ResponseDecoder: singleDecoder,
			ResponseCode:    http.StatusOK,
		},
		{
			Name: "DELETE game subscription",
			HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
				return rnr.GetHandlerConfig()[game.DeleteOneGameSubscription]
			},
			RequestHeaders: testutil.AuthHeaderProManager,
			RequestPathParams: func(d harness.Data) map[string]string {
				subID, ok := d.Refs.GameSubscriptionRefs[harness.GameSubscriptionManagerOneRef]
				require.True(t, ok, "manager subscription ref exists")
				return map[string]string{":game_subscription_id": subID}
			},
			ResponseCode: http.StatusNoContent,
		},
	}

	for _, testCase := range testCases {

		t.Logf("Running test >%s<\n", testCase.Name)

		t.Run(testCase.Name, func(t *testing.T) {
			testutil.RunTestCase(t, th, &testCase, func(method string, body any) {
				if testCase.TestResponseCode() == http.StatusNoContent {
					return
				}
				require.NotNil(t, body, "Response body is not nil")
			})
		})
	}
}
