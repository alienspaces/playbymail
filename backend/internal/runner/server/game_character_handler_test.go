package runner

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/playbymail/core/server"
	"gitlab.com/alienspaces/playbymail/internal/harness"
	"gitlab.com/alienspaces/playbymail/schema"
)

func Test_gameCharacterHandler(t *testing.T) {
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
	accountRec, err := th.Data.GetAccountRecByRef(harness.AccountThreeRef)
	require.NoError(t, err, "GetAccountRecByRef(AccountThreeRef) returns without error")
	charRec, err := th.Data.GetGameCharacterRecByRef(harness.GameCharacterOneRef)
	require.NoError(t, err, "GetGameCharacterRecByRef returns without error")

	testCaseCollectionResponseDecoder := testCaseResponseDecoderGeneric[schema.GameCharacterCollectionResponse]
	testCaseResponseDecoder := testCaseResponseDecoderGeneric[schema.GameCharacterResponse]

	testCases := []struct {
		TestCase
	}{
		{
			TestCase: TestCase{
				Name: "API key with open access \\ get many game characters \\ returns expected characters",
				HandlerConfig: func(rnr *Runner) server.HandlerConfig {
					return rnr.HandlerConfig[getManyGameCharacters]
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
				Name: "API key with open access \\ create game character with valid properties \\ returns created character",
				HandlerConfig: func(rnr *Runner) server.HandlerConfig {
					return rnr.HandlerConfig[createGameCharacter]
				},
				RequestBody: func(d harness.Data) any {
					return schema.GameCharacterRequest{
						GameID:    gameRec.ID,
						AccountID: accountRec.ID, // Use AccountTwoRef for uniqueness
						Name:      "Test Character",
					}
				},
				ResponseDecoder: testCaseResponseDecoder,
				ResponseCode:    http.StatusCreated,
			},
		},
		{
			TestCase: TestCase{
				Name: "API key with open access \\ get one game character with valid ID \\ returns expected character",
				HandlerConfig: func(rnr *Runner) server.HandlerConfig {
					return rnr.HandlerConfig[getOneGameCharacter]
				},
				RequestPathParams: func(d harness.Data) map[string]string {
					return map[string]string{":game_character_id": charRec.ID}
				},
				ResponseDecoder: testCaseResponseDecoder,
				ResponseCode:    http.StatusOK,
			},
		},
		{
			TestCase: TestCase{
				Name: "API key with open access \\ delete game character with valid ID \\ returns no content",
				HandlerConfig: func(rnr *Runner) server.HandlerConfig {
					return rnr.HandlerConfig[deleteGameCharacter]
				},
				RequestPathParams: func(d harness.Data) map[string]string {
					return map[string]string{":game_character_id": charRec.ID}
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
