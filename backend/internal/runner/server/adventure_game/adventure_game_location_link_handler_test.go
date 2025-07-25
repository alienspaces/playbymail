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
	"gitlab.com/alienspaces/playbymail/schema"
)

func Test_adventureGameLocationLinkHandler(t *testing.T) {
	t.Parallel()

	th := deps.NewHandlerTestHarness(t)
	require.NotNil(t, th, "NewTestHarness returns without error")

	// Add a new location to the game that we can use in the test
	th.DataConfig.GameConfigs[0].GameLocationConfigs = append(
		th.DataConfig.GameConfigs[0].GameLocationConfigs,
		harness.GameLocationConfig{
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
	locationOneRec, err := th.Data.GetGameLocationRecByRef(harness.GameLocationOneRef)
	require.NoError(t, err, "GetGameLocationRecByRef returns without error")
	locationThreeRec, err := th.Data.GetGameLocationRecByRef(harness.GameLocationThreeRef)
	require.NoError(t, err, "GetGameLocationRecByRef returns without error")
	linkRec, err := th.Data.GetGameLocationLinkRecByRef(harness.GameLocationLinkOneRef)
	require.NoError(t, err, "GetGameLocationLinkRecByRef returns without error")

	testCaseCollectionResponseDecoder := testutil.TestCaseResponseDecoderGeneric[schema.AdventureGameLocationLinkCollectionResponse]
	testCaseResponseDecoder := testutil.TestCaseResponseDecoderGeneric[schema.AdventureGameLocationLinkResponse]

	testCases := []struct {
		testutil.TestCase
	}{
		{
			TestCase: testutil.TestCase{
				Name: "API key with open access \\ get many location links \\ returns expected links",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[adventure_game.SearchManyAdventureGameLocationLinks]
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
				Name: "API key with open access \\ get one location link with valid ID \\ returns expected link",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[adventure_game.GetOneAdventureGameLocationLink]
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
				Name: "API key with open access \\ create location link with valid properties \\ returns created link",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[adventure_game.CreateOneAdventureGameLocationLink]
				},
				RequestPathParams: func(d harness.Data) map[string]string {
					return map[string]string{
						":game_id": gameRec.ID,
					}
				},
				RequestBody: func(d harness.Data) any {
					return schema.AdventureGameLocationLinkRequest{
						FromGameLocationID: locationOneRec.ID,
						ToGameLocationID:   locationThreeRec.ID,
						Name:               "Test Link",
					}
				},
				ResponseDecoder: testCaseResponseDecoder,
				ResponseCode:    http.StatusCreated,
			},
		},
		{
			TestCase: testutil.TestCase{
				Name: "API key with open access \\ update location link with valid properties \\ returns updated link",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[adventure_game.UpdateOneAdventureGameLocationLink]
				},
				RequestPathParams: func(d harness.Data) map[string]string {
					return map[string]string{
						":game_id":          gameRec.ID,
						":location_link_id": linkRec.ID,
					}
				},
				RequestBody: func(d harness.Data) any {
					return schema.AdventureGameLocationLinkRequest{
						Name:               "Updated Test Link",
						Description:        "Updated Test Description",
						FromGameLocationID: linkRec.FromAdventureGameLocationID,
						ToGameLocationID:   linkRec.ToAdventureGameLocationID,
					}
				},
				ResponseDecoder: testCaseResponseDecoder,
				ResponseCode:    http.StatusCreated,
			},
		},
		{
			TestCase: testutil.TestCase{
				Name: "API key with open access \\ delete location link with valid ID \\ returns no content",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[adventure_game.DeleteOneAdventureGameLocationLink]
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
