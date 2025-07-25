package runner_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.com/alienspaces/playbymail/core/server"
	"gitlab.com/alienspaces/playbymail/internal/harness"
	runner "gitlab.com/alienspaces/playbymail/internal/runner/server"
	"gitlab.com/alienspaces/playbymail/internal/utils/testutil"
	"gitlab.com/alienspaces/playbymail/schema"
)

func Test_gameSubscriptionHandler(t *testing.T) {
	th := testutil.NewTestHarness(t)
	require.NotNil(t, th, "newTestHarness returns without error")
	_, err := th.Setup()
	require.NoError(t, err, "Test data setup returns without error")
	defer func() { _ = th.Teardown() }()

	collectionDecoder := testutil.TestCaseResponseDecoderGeneric[schema.GameSubscriptionCollectionResponse]
	singleDecoder := testutil.TestCaseResponseDecoderGeneric[schema.GameSubscriptionResponse]

	testCases := []testutil.TestCase{
		{
			Name: "GET many game subscriptions",
			HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
				return rnr.GetHandlerConfig()[runner.GetManyGameSubscriptions]
			},
			ResponseDecoder: collectionDecoder,
			ResponseCode:    http.StatusOK,
		},
		{
			Name: "POST create game subscription",
			HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
				return rnr.GetHandlerConfig()[runner.CreateGameSubscription]
			},
			RequestBody: func(d harness.Data) any {
				gameRec, _ := d.GetGameRecByRef(harness.GameOneRef)
				accountRec, _ := d.GetAccountRecByRef(harness.AccountTwoRef)
				return schema.GameSubscriptionRequest{
					GameID:           gameRec.ID,
					AccountID:        accountRec.ID,
					SubscriptionType: "Player",
				}
			},
			ResponseDecoder: singleDecoder,
			ResponseCode:    http.StatusCreated,
		},
		{
			Name: "GET one game subscription",
			HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
				return rnr.GetHandlerConfig()[runner.GetOneGameSubscription]
			},
			RequestPathParams: func(d harness.Data) map[string]string {
				// Use the first created subscription
				return map[string]string{":game_subscription_id": d.GameSubscriptionRecs[0].ID}
			},
			ResponseDecoder: singleDecoder,
			ResponseCode:    http.StatusOK,
		},
		{
			Name: "PUT update game subscription",
			HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
				return rnr.GetHandlerConfig()[runner.UpdateGameSubscription]
			},
			RequestPathParams: func(d harness.Data) map[string]string {
				return map[string]string{":game_subscription_id": d.GameSubscriptionRecs[0].ID}
			},
			RequestBody: func(d harness.Data) any {
				return schema.GameSubscriptionRequest{
					GameID:           d.GameSubscriptionRecs[0].GameID,
					AccountID:        d.GameSubscriptionRecs[0].AccountID,
					SubscriptionType: "Manager",
				}
			},
			ResponseDecoder: singleDecoder,
			ResponseCode:    http.StatusOK,
		},
		{
			Name: "DELETE game subscription",
			HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
				return rnr.GetHandlerConfig()[runner.DeleteGameSubscription]
			},
			RequestPathParams: func(d harness.Data) map[string]string {
				return map[string]string{":game_subscription_id": d.GameSubscriptionRecs[0].ID}
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
