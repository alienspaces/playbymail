package adventure_game_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/playbymail/core/server"
	"gitlab.com/alienspaces/playbymail/internal/harness"
	"gitlab.com/alienspaces/playbymail/internal/runner/server/adventure_game"
	"gitlab.com/alienspaces/playbymail/internal/utils/deps"
	"gitlab.com/alienspaces/playbymail/internal/utils/testutil"
	"gitlab.com/alienspaces/playbymail/schema/api/adventure_game_schema"
)

func Test_adventureGameLocationObjectEffectHandler(t *testing.T) {
	t.Parallel()

	th := deps.NewHandlerTestHarness(t)
	require.NotNil(t, th, "NewHandlerTestHarness returns without error")

	_, err := th.Setup()
	require.NoError(t, err, "Test data setup returns without error")
	defer func() {
		err = th.Teardown()
		require.NoError(t, err, "Test data teardown returns without error")
	}()

	gameRec, err := th.Data.GetGameRecByRef(harness.GameOneRef)
	require.NoError(t, err, "GetGameRecByRef returns without error")

	locationObjectRec, err := th.Data.GetAdventureGameLocationObjectRecByRef(harness.GameLocationObjectOneRef)
	require.NoError(t, err, "GetAdventureGameLocationObjectRecByRef returns without error")

	locationObjectEffectRec, err := th.Data.GetAdventureGameLocationObjectEffectRecByRef(harness.GameLocationObjectEffectOneRef)
	require.NoError(t, err, "GetAdventureGameLocationObjectEffectRecByRef returns without error")

	gameDraftRec, err := th.Data.GetGameRecByRef(harness.GameDraftRef)
	require.NoError(t, err, "GetGameRecByRef(GameDraftRef) returns without error")

	testCaseCollectionResponseDecoder := testutil.TestCaseResponseDecoderGeneric[adventure_game_schema.AdventureGameLocationObjectEffectCollectionResponse]
	testCaseResponseDecoder := testutil.TestCaseResponseDecoderGeneric[adventure_game_schema.AdventureGameLocationObjectEffectResponse]

	testCases := []struct {
		testutil.TestCase
	}{
		{
			TestCase: testutil.TestCase{
				Name: "authenticated user when search many location object effects then returns expected effects",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[adventure_game.SearchManyAdventureGameLocationObjectEffects]
				},
				RequestHeaders: testutil.AuthHeaderStandard,
				RequestQueryParams: func(d harness.Data) map[string]any {
					return map[string]any{
						"page_size":   10,
						"page_number": 1,
					}
				},
				ResponseDecoder: testCaseCollectionResponseDecoder,
				ResponseCode:    http.StatusOK,
			},
		},
		{
			TestCase: testutil.TestCase{
				Name: "authenticated user when get many location object effects for game then returns expected effects",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[adventure_game.GetManyAdventureGameLocationObjectEffects]
				},
				RequestHeaders: testutil.AuthHeaderStandard,
				RequestPathParams: func(d harness.Data) map[string]string {
					return map[string]string{":game_id": gameRec.ID}
				},
				RequestQueryParams: func(d harness.Data) map[string]any {
					return map[string]any{
						"page_size":   10,
						"page_number": 1,
					}
				},
				ResponseDecoder: testCaseCollectionResponseDecoder,
				ResponseCode:    http.StatusOK,
			},
		},
		{
			TestCase: testutil.TestCase{
				Name: "authenticated user when get one location object effect then returns expected effect",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[adventure_game.GetOneAdventureGameLocationObjectEffect]
				},
				RequestHeaders: testutil.AuthHeaderStandard,
				RequestPathParams: func(d harness.Data) map[string]string {
					return map[string]string{
						":game_id":                    gameRec.ID,
						":location_object_effect_id":  locationObjectEffectRec.ID,
					}
				},
				ResponseDecoder: testCaseResponseDecoder,
				ResponseCode:    http.StatusOK,
			},
		},
		{
			TestCase: testutil.TestCase{
				Name: "authenticated designer when create location object effect with valid properties then returns created effect",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[adventure_game.CreateOneAdventureGameLocationObjectEffect]
				},
				RequestHeaders: testutil.AuthHeaderProDesigner,
				RequestBody: func(d harness.Data) any {
					return adventure_game_schema.AdventureGameLocationObjectEffectRequest{
						AdventureGameLocationObjectID: locationObjectRec.ID,
						ActionType:                    "inspect",
						ResultDescription:             "You examine the object closely.",
						EffectType:                    "info",
						IsRepeatable:                  true,
					}
				},
				RequestPathParams: func(d harness.Data) map[string]string {
					return map[string]string{":game_id": gameRec.ID}
				},
				ResponseDecoder: testCaseResponseDecoder,
				ResponseCode:    http.StatusCreated,
			},
		},
		{
			TestCase: testutil.TestCase{
				Name: "authenticated designer when update location object effect then returns updated effect",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[adventure_game.UpdateOneAdventureGameLocationObjectEffect]
				},
				RequestHeaders: testutil.AuthHeaderProDesigner,
				RequestBody: func(d harness.Data) any {
					return adventure_game_schema.AdventureGameLocationObjectEffectRequest{
						AdventureGameLocationObjectID: locationObjectRec.ID,
						ActionType:                    "touch",
						ResultDescription:             "You touch the object and feel a strange energy.",
						EffectType:                    "info",
						IsRepeatable:                  false,
					}
				},
				RequestPathParams: func(d harness.Data) map[string]string {
					return map[string]string{
						":game_id":                   gameRec.ID,
						":location_object_effect_id": locationObjectEffectRec.ID,
					}
				},
				ResponseDecoder: testCaseResponseDecoder,
				ResponseCode:    http.StatusOK,
			},
		},
		{
			TestCase: testutil.TestCase{
				Name: "authenticated designer when delete location object effect then returns no content",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[adventure_game.DeleteOneAdventureGameLocationObjectEffect]
				},
				RequestHeaders: testutil.AuthHeaderProDesigner,
				RequestPathParams: func(d harness.Data) map[string]string {
					return map[string]string{
						":game_id":                   gameRec.ID,
						":location_object_effect_id": locationObjectEffectRec.ID,
					}
				},
				ResponseCode: http.StatusNoContent,
			},
		},
		{
			TestCase: testutil.TestCase{
				Name: "unauthenticated user when create location object effect then returns unauthorized",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[adventure_game.CreateOneAdventureGameLocationObjectEffect]
				},
				RequestBody: func(d harness.Data) any {
					return adventure_game_schema.AdventureGameLocationObjectEffectRequest{
						AdventureGameLocationObjectID: locationObjectRec.ID,
						ActionType:                    "inspect",
						ResultDescription:             "You examine the object.",
						EffectType:                    "info",
					}
				},
				RequestPathParams: func(d harness.Data) map[string]string {
					return map[string]string{":game_id": gameRec.ID}
				},
				ResponseCode: http.StatusUnauthorized,
			},
		},
		{
			TestCase: testutil.TestCase{
				Name: "authenticated designer without game ownership when create location object effect then returns unauthorized",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[adventure_game.CreateOneAdventureGameLocationObjectEffect]
				},
				RequestHeaders: testutil.AuthHeaderProDesigner,
				RequestBody: func(d harness.Data) any {
					return adventure_game_schema.AdventureGameLocationObjectEffectRequest{
						AdventureGameLocationObjectID: locationObjectRec.ID,
						ActionType:                    "inspect",
						ResultDescription:             "Should not be created",
						EffectType:                    "info",
					}
				},
			RequestPathParams: func(d harness.Data) map[string]string {
				return map[string]string{":game_id": gameDraftRec.ID}
			},
			ResponseCode: http.StatusForbidden,
		},
	},
	}

	for _, tc := range testCases {
		t.Run(tc.TestName(), func(t *testing.T) {
			testutil.RunTestCase(t, th, &tc.TestCase, nil)
		})
	}
}
