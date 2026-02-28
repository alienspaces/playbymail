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

func Test_adventureGameLocationLinkHandler(t *testing.T) {
	t.Parallel()

	th := deps.NewHandlerTestHarness(t)
	require.NotNil(t, th, "NewTestHarness returns without error")

	// Add a new location to the game that we can use in the test
	th.DataConfig.GameConfigs[0].AdventureGameLocationConfigs = append(
		th.DataConfig.GameConfigs[0].AdventureGameLocationConfigs,
		harness.AdventureGameLocationConfig{
			Reference: harness.GameLocationThreeRef,
		},
	)

	_, err := th.Setup()
	require.NoError(t, err, "Test data setup returns without error")
	defer func() {
		err = th.Teardown()
		require.NoError(t, err, "Test data teardown returns without error")
	}()

	gameRec, err := th.Data.GetGameRecByRef(harness.GameOneRef)
	require.NoError(t, err, "GetGameRecByRef returns without error")
	locationOneRec, err := th.Data.GetAdventureGameLocationRecByRef(harness.GameLocationOneRef)
	require.NoError(t, err, "GetGameLocationRecByRef returns without error")
	locationThreeRec, err := th.Data.GetAdventureGameLocationRecByRef(harness.GameLocationThreeRef)
	require.NoError(t, err, "GetGameLocationRecByRef returns without error")
	linkRec, err := th.Data.GetAdventureGameLocationLinkRecByRef(harness.GameLocationLinkOneRef)
	require.NoError(t, err, "GetGameLocationLinkRecByRef returns without error")

	testCaseCollectionResponseDecoder := testutil.TestCaseResponseDecoderGeneric[adventure_game_schema.AdventureGameLocationLinkCollectionResponse]
	testCaseResponseDecoder := testutil.TestCaseResponseDecoderGeneric[adventure_game_schema.AdventureGameLocationLinkResponse]

	testCases := []struct {
		testutil.TestCase
	}{
		{
			TestCase: testutil.TestCase{
				Name: "authenticated user when get many location links then returns expected links",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[adventure_game.SearchManyAdventureGameLocationLinks]
				},
				RequestHeaders: func(d harness.Data) map[string]string {
					return testutil.AuthHeaderStandard(d)
				},
				RequestPathParams: func(d harness.Data) map[string]string {
					return map[string]string{
						":game_id": gameRec.ID,
					}
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
				Name: "authenticated user when get one location link with valid ID then returns expected link",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[adventure_game.GetOneAdventureGameLocationLink]
				},
				RequestHeaders: func(d harness.Data) map[string]string {
					return testutil.AuthHeaderStandard(d)
				},
				RequestPathParams: func(d harness.Data) map[string]string {
					return map[string]string{
						":game_id":          gameRec.ID,
						":location_link_id": linkRec.ID,
					}
				},
				ResponseDecoder: testCaseResponseDecoder,
				ResponseCode:    http.StatusOK,
			},
		},
		{
			TestCase: testutil.TestCase{
				Name: "authenticated designer when create location link with valid properties then returns created link",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[adventure_game.CreateOneAdventureGameLocationLink]
				},
				RequestHeaders: func(d harness.Data) map[string]string {
					return testutil.AuthHeaderProDesigner(d)
				},
				RequestPathParams: func(d harness.Data) map[string]string {
					return map[string]string{
						":game_id": gameRec.ID,
					}
				},
				RequestBody: func(d harness.Data) any {
					return adventure_game_schema.AdventureGameLocationLinkRequest{
						Name:                        "Test Link",
						Description:                 "Test Link Description",
						FromAdventureGameLocationID: locationOneRec.ID,
						ToAdventureGameLocationID:   locationThreeRec.ID,
					}
				},
				ResponseDecoder: testCaseResponseDecoder,
				ResponseCode:    http.StatusCreated,
			},
		},
		{
			TestCase: testutil.TestCase{
				Name: "authenticated designer when update location link with valid properties then returns updated link",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[adventure_game.UpdateOneAdventureGameLocationLink]
				},
				RequestHeaders: func(d harness.Data) map[string]string {
					return testutil.AuthHeaderProDesigner(d)
				},
				RequestPathParams: func(d harness.Data) map[string]string {
					return map[string]string{
						":game_id":          gameRec.ID,
						":location_link_id": linkRec.ID,
					}
				},
				RequestBody: func(d harness.Data) any {
					return adventure_game_schema.AdventureGameLocationLinkRequest{
						Name:                        "Updated Test Link",
						Description:                 "Updated Test Description",
						FromAdventureGameLocationID: linkRec.FromAdventureGameLocationID,
						ToAdventureGameLocationID:   linkRec.ToAdventureGameLocationID,
					}
				},
				ResponseDecoder: testCaseResponseDecoder,
				ResponseCode:    http.StatusCreated,
			},
		},
		{
			TestCase: testutil.TestCase{
				Name: "authenticated designer when delete location link with valid ID then returns no content",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[adventure_game.DeleteOneAdventureGameLocationLink]
				},
				RequestHeaders: func(d harness.Data) map[string]string {
					return testutil.AuthHeaderProDesigner(d)
				},
				RequestPathParams: func(d harness.Data) map[string]string {
					return map[string]string{
						":game_id":          gameRec.ID,
						":location_link_id": linkRec.ID,
					}
				},
				ResponseCode: http.StatusNoContent,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.TestName(), func(t *testing.T) {
			testutil.RunTestCase(t, th, &tc.TestCase, nil)
		})
	}
}
