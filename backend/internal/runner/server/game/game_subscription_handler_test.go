package game_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.com/alienspaces/playbymail/core/server"
	"gitlab.com/alienspaces/playbymail/internal/harness"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
	"gitlab.com/alienspaces/playbymail/internal/runner/server/game"
	"gitlab.com/alienspaces/playbymail/internal/utils/testutil"
	"gitlab.com/alienspaces/playbymail/schema/api/game_schema"
)

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
			RequestHeaders: func(d harness.Data) map[string]string {
				return testutil.AuthHeaderProManager(d)
			},
			ResponseDecoder: collectionDecoder,
			ResponseCode:    http.StatusOK,
		},
		{
			Name: "POST create game subscription",
			HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
				return rnr.GetHandlerConfig()[game.CreateOneGameSubscription]
			},
			RequestHeaders: func(d harness.Data) map[string]string {
				return testutil.AuthHeaderProManager(d)
			},
			RequestBody: func(d harness.Data) any {
				gameRec, _ := d.GetGameRecByRef(harness.GameOneRef)
				accountRec, _ := d.GetAccountUserRecByRef(harness.ProManagerAccountRef)
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
			RequestHeaders: func(d harness.Data) map[string]string {
				return testutil.AuthHeaderProManager(d)
			},
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
			RequestHeaders: func(d harness.Data) map[string]string {
				return testutil.AuthHeaderProManager(d)
			},
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
			RequestHeaders: func(d harness.Data) map[string]string {
				return testutil.AuthHeaderProManager(d)
			},
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
