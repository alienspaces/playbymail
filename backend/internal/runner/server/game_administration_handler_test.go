package runner

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.com/alienspaces/playbymail/core/server"
	"gitlab.com/alienspaces/playbymail/internal/harness"
	"gitlab.com/alienspaces/playbymail/schema"
)

func Test_gameAdministrationHandler(t *testing.T) {
	th := newTestHarness(t)
	require.NotNil(t, th, "newTestHarness returns without error")
	_, err := th.Setup()
	require.NoError(t, err, "Test data setup returns without error")
	defer func() { _ = th.Teardown() }()

	collectionDecoder := testCaseResponseDecoderGeneric[schema.GameAdministrationCollectionResponse]
	singleDecoder := testCaseResponseDecoderGeneric[schema.GameAdministrationResponse]

	testCases := []TestCase{
		{
			Name: "GET many game administrations",
			HandlerConfig: func(rnr *Runner) server.HandlerConfig {
				return rnr.HandlerConfig[getManyGameAdministrations]
			},
			ResponseDecoder: collectionDecoder,
			ResponseCode:    http.StatusOK,
		},
		{
			Name: "POST create game administration",
			HandlerConfig: func(rnr *Runner) server.HandlerConfig {
				return rnr.HandlerConfig[createGameAdministration]
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
			HandlerConfig: func(rnr *Runner) server.HandlerConfig {
				return rnr.HandlerConfig[getOneGameAdministration]
			},
			RequestPathParams: func(d harness.Data) map[string]string {
				return map[string]string{":game_administration_id": d.GameAdministrationRecs[0].ID}
			},
			ResponseDecoder: singleDecoder,
			ResponseCode:    http.StatusOK,
		},
		{
			Name: "PUT update game administration",
			HandlerConfig: func(rnr *Runner) server.HandlerConfig {
				return rnr.HandlerConfig[updateGameAdministration]
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
			HandlerConfig: func(rnr *Runner) server.HandlerConfig {
				return rnr.HandlerConfig[deleteGameAdministration]
			},
			RequestPathParams: func(d harness.Data) map[string]string {
				return map[string]string{":game_administration_id": d.GameAdministrationRecs[0].ID}
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
