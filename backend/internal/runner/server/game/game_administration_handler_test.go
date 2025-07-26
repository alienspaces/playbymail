package game_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.com/alienspaces/playbymail/core/server"
	"gitlab.com/alienspaces/playbymail/internal/harness"
	"gitlab.com/alienspaces/playbymail/internal/runner/server/game"
	"gitlab.com/alienspaces/playbymail/internal/utils/testutil"
	"gitlab.com/alienspaces/playbymail/schema"
)

func Test_gameAdministrationHandler(t *testing.T) {
	th := testutil.NewTestHarness(t)
	require.NotNil(t, th, "newTestHarness returns without error")
	_, err := th.Setup()
	require.NoError(t, err, "Test data setup returns without error")
	defer func() { _ = th.Teardown() }()

	collectionDecoder := testutil.TestCaseResponseDecoderGeneric[schema.GameAdministrationCollectionResponse]
	singleDecoder := testutil.TestCaseResponseDecoderGeneric[schema.GameAdministrationResponse]

	testCases := []testutil.TestCase{
		{
			Name: "GET many game administrations",
			HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
				return rnr.GetHandlerConfig()[game.GetManyGameAdministrations]
			},
			ResponseDecoder: collectionDecoder,
			ResponseCode:    http.StatusOK,
		},
		{
			Name: "POST create game administration",
			HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
				return rnr.GetHandlerConfig()[game.CreateGameAdministration]
			},
			RequestBody: func(d harness.Data) any {
				gameRec, _ := d.GetGameRecByRef(harness.GameOneRef)
				accountRec, _ := d.GetAccountRecByRef(harness.AccountTwoRef)
				return schema.GameAdministrationRequest{
					GameID:             gameRec.ID,
					AccountID:          accountRec.ID,
					GrantedByAccountID: accountRec.ID,
				}
			},
			ResponseDecoder: singleDecoder,
			ResponseCode:    http.StatusCreated,
		},
		{
			Name: "GET one game administration",
			HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
				return rnr.GetHandlerConfig()[game.GetOneGameAdministration]
			},
			RequestPathParams: func(d harness.Data) map[string]string {
				return map[string]string{":game_administration_id": d.GameAdministrationRecs[0].ID}
			},
			ResponseDecoder: singleDecoder,
			ResponseCode:    http.StatusOK,
		},
		{
			Name: "PUT update game administration",
			HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
				return rnr.GetHandlerConfig()[game.UpdateGameAdministration]
			},
			RequestPathParams: func(d harness.Data) map[string]string {
				return map[string]string{":game_administration_id": d.GameAdministrationRecs[0].ID}
			},
			RequestBody: func(d harness.Data) any {
				return schema.GameAdministrationRequest{
					GameID:             d.GameAdministrationRecs[0].GameID,
					AccountID:          d.GameAdministrationRecs[0].AccountID,
					GrantedByAccountID: d.GameAdministrationRecs[0].GrantedByAccountID,
				}
			},
			ResponseDecoder: singleDecoder,
			ResponseCode:    http.StatusOK,
		},
		{
			Name: "DELETE game administration",
			HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
				return rnr.GetHandlerConfig()[game.DeleteGameAdministration]
			},
			RequestPathParams: func(d harness.Data) map[string]string {
				return map[string]string{":game_administration_id": d.GameAdministrationRecs[0].ID}
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
