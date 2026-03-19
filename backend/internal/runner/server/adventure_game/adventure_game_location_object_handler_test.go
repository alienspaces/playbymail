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

func Test_adventureGameLocationObjectHandler(t *testing.T) {
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

	locationRec, err := th.Data.GetAdventureGameLocationRecByRef(harness.GameLocationOneRef)
	require.NoError(t, err, "GetAdventureGameLocationRecByRef returns without error")

	locationObjectRec, err := th.Data.GetAdventureGameLocationObjectRecByRef(harness.GameLocationObjectOneRef)
	require.NoError(t, err, "GetAdventureGameLocationObjectRecByRef returns without error")

	testCaseCollectionResponseDecoder := testutil.TestCaseResponseDecoderGeneric[adventure_game_schema.AdventureGameLocationObjectCollectionResponse]
	testCaseResponseDecoder := testutil.TestCaseResponseDecoderGeneric[adventure_game_schema.AdventureGameLocationObjectResponse]

	testCases := []struct {
		testutil.TestCase
	}{
		{
			TestCase: testutil.TestCase{
				Name: "authenticated user when search many location objects then returns expected objects",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[adventure_game.SearchManyAdventureGameLocationObjects]
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
				Name: "authenticated user when get many location objects for game then returns expected objects",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[adventure_game.GetManyAdventureGameLocationObjects]
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
				Name: "authenticated user when get one location object then returns expected object",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[adventure_game.GetOneAdventureGameLocationObject]
				},
				RequestHeaders: testutil.AuthHeaderStandard,
				RequestPathParams: func(d harness.Data) map[string]string {
					return map[string]string{
						":game_id":            gameRec.ID,
						":location_object_id": locationObjectRec.ID,
					}
				},
				ResponseDecoder: testCaseResponseDecoder,
				ResponseCode:    http.StatusOK,
			},
		},
		{
			TestCase: testutil.TestCase{
				Name: "authenticated designer when create location object with valid properties then returns created object",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[adventure_game.CreateOneAdventureGameLocationObject]
				},
				RequestHeaders: testutil.AuthHeaderProDesigner,
				RequestBody: func(d harness.Data) any {
					return adventure_game_schema.AdventureGameLocationObjectRequest{
						AdventureGameLocationID: locationRec.ID,
						Name:                    "Test Object",
						Description:             "A test object for handler tests",
						InitialState:            "intact",
						IsHidden:                false,
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
				Name: "authenticated designer when update location object then returns updated object",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[adventure_game.UpdateOneAdventureGameLocationObject]
				},
				RequestHeaders: testutil.AuthHeaderProDesigner,
				RequestBody: func(d harness.Data) any {
					return adventure_game_schema.AdventureGameLocationObjectRequest{
						AdventureGameLocationID: locationRec.ID,
						Name:                    "Updated Object",
						Description:             "An updated object description",
						InitialState:            "closed",
						IsHidden:                true,
					}
				},
				RequestPathParams: func(d harness.Data) map[string]string {
					return map[string]string{
						":game_id":            gameRec.ID,
						":location_object_id": locationObjectRec.ID,
					}
				},
				ResponseDecoder: testCaseResponseDecoder,
				ResponseCode:    http.StatusOK,
			},
		},
		{
			TestCase: testutil.TestCase{
				Name: "authenticated designer when delete location object then returns no content",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[adventure_game.DeleteOneAdventureGameLocationObject]
				},
				RequestHeaders: testutil.AuthHeaderProDesigner,
				RequestPathParams: func(d harness.Data) map[string]string {
					return map[string]string{
						":game_id":            gameRec.ID,
						":location_object_id": locationObjectRec.ID,
					}
				},
				ResponseCode: http.StatusNoContent,
			},
		},
		{
			TestCase: testutil.TestCase{
				Name: "unauthenticated user when create location object then returns unauthorized",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[adventure_game.CreateOneAdventureGameLocationObject]
				},
				RequestBody: func(d harness.Data) any {
					return adventure_game_schema.AdventureGameLocationObjectRequest{
						AdventureGameLocationID: locationRec.ID,
						Name:                    "Test Object",
						Description:             "A test object",
					}
				},
				RequestPathParams: func(d harness.Data) map[string]string {
					return map[string]string{":game_id": gameRec.ID}
				},
				ResponseCode: http.StatusUnauthorized,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.TestName(), func(t *testing.T) {
			testutil.RunTestCase(t, th, &tc.TestCase, nil)
		})
	}
}
