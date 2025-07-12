package runner

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/playbymail/core/server"
	"gitlab.com/alienspaces/playbymail/internal/harness"
	"gitlab.com/alienspaces/playbymail/schema"
)

func Test_gameCreatureHandler(t *testing.T) {
	t.Parallel()

	th := newTestHarness(t)
	require.NotNil(t, th, "newTestHarness returns without error")

	_, err := th.Setup()
	require.NoError(t, err, "Test data setup returns without error")
	defer func() {
		err = th.Teardown()
		require.NoError(t, err, "Test data teardown returns without error")
	}()

	gameRec, err := th.Data.GetGameRecByRef(harness.GameOneRef)
	require.NoError(t, err, "GetGameRecByRef returns without error")
	creatureRec, err := th.Data.GetGameCreatureRecByRef(harness.GameCreatureOneRef)
	require.NoError(t, err, "GetGameCreatureRecByRef returns without error")

	testCaseCollectionResponseDecoder := testCaseResponseDecoderGeneric[schema.GameCreatureCollectionResponse]
	testCaseResponseDecoder := testCaseResponseDecoderGeneric[schema.GameCreatureResponse]

	testCases := []struct {
		TestCase
	}{
		{
			TestCase: TestCase{
				Name: "API key with open access \\ get many game creatures \\ returns expected creatures",
				HandlerConfig: func(rnr *Runner) server.HandlerConfig {
					return rnr.HandlerConfig[getManyGameCreatures]
				},
				RequestQueryParams: func(d harness.Data) map[string]any {
					return map[string]any{
						"game_id":     gameRec.ID,
						"page_size":   10,
						"page_number": 1,
					}
				},
				ResponseDecoder: testCaseCollectionResponseDecoder,
				ResponseCode:    http.StatusOK,
			},
		},
		{
			TestCase: TestCase{
				Name: "API key with open access \\ create game creature with valid properties \\ returns created creature",
				HandlerConfig: func(rnr *Runner) server.HandlerConfig {
					return rnr.HandlerConfig[createGameCreature]
				},
				RequestBody: func(d harness.Data) any {
					return schema.GameCreatureRequest{
						GameID: gameRec.ID,
						Name:   "Test Creature",
					}
				},
				ResponseDecoder: testCaseResponseDecoder,
				ResponseCode:    http.StatusCreated,
			},
		},
		{
			TestCase: TestCase{
				Name: "API key with open access \\ get one game creature with valid ID \\ returns expected creature",
				HandlerConfig: func(rnr *Runner) server.HandlerConfig {
					return rnr.HandlerConfig[getOneGameCreature]
				},
				RequestPathParams: func(d harness.Data) map[string]string {
					return map[string]string{":game_creature_id": creatureRec.ID}
				},
				ResponseDecoder: testCaseResponseDecoder,
				ResponseCode:    http.StatusOK,
			},
		},
		{
			TestCase: TestCase{
				Name: "API key with open access \\ delete game creature with valid ID \\ returns no content",
				HandlerConfig: func(rnr *Runner) server.HandlerConfig {
					return rnr.HandlerConfig[deleteGameCreature]
				},
				RequestPathParams: func(d harness.Data) map[string]string {
					return map[string]string{":game_creature_id": creatureRec.ID}
				},
				ResponseCode: http.StatusNoContent,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.TestName(), func(t *testing.T) {
			RunTestCase(t, th, &tc.TestCase, nil)
		})
	}
}
