package runner

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.com/alienspaces/playbymail/core/server"
	"gitlab.com/alienspaces/playbymail/internal/harness"
	"gitlab.com/alienspaces/playbymail/schema"
)

func Test_gameSubscriptionHandler(t *testing.T) {
	th := newTestHarness(t)
	require.NotNil(t, th, "newTestHarness returns without error")
	_, err := th.Setup()
	require.NoError(t, err, "Test data setup returns without error")
	defer func() { _ = th.Teardown() }()

	collectionDecoder := testCaseResponseDecoderGeneric[schema.GameSubscriptionCollectionResponse]
	singleDecoder := testCaseResponseDecoderGeneric[schema.GameSubscriptionResponse]

	testCases := []TestCase{
		{
			Name: "GET many game subscriptions",
			HandlerConfig: func(rnr *Runner) server.HandlerConfig {
				return rnr.HandlerConfig[getManyGameSubscriptions]
			},
			ResponseDecoder: collectionDecoder,
			ResponseCode:    http.StatusOK,
		},
		{
			Name: "POST create game subscription",
			HandlerConfig: func(rnr *Runner) server.HandlerConfig {
				return rnr.HandlerConfig[createGameSubscription]
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
			HandlerConfig: func(rnr *Runner) server.HandlerConfig {
				return rnr.HandlerConfig[getOneGameSubscription]
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
			HandlerConfig: func(rnr *Runner) server.HandlerConfig {
				return rnr.HandlerConfig[updateGameSubscription]
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
			HandlerConfig: func(rnr *Runner) server.HandlerConfig {
				return rnr.HandlerConfig[deleteGameSubscription]
			},
			RequestPathParams: func(d harness.Data) map[string]string {
				return map[string]string{":game_subscription_id": d.GameSubscriptionRecs[0].ID}
			},
			ResponseCode: http.StatusNoContent,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			RunTestCase(t, th, &testCase, func(method string, body interface{}) {
				if testCase.TestResponseCode() == http.StatusNoContent {
					return
				}
				require.NotNil(t, body, "Response body is not nil")
			})
		})
	}
}
