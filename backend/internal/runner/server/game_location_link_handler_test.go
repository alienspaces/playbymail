package runner

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/playbymail/core/server"
	"gitlab.com/alienspaces/playbymail/internal/harness"
	"gitlab.com/alienspaces/playbymail/schema"
)

func Test_gameGameLocationLinkHandler(t *testing.T) {
	t.Parallel()

	th := newTestHarness(t)
	require.NotNil(t, th, "newTestHarness returns without error")

	_, err := th.Setup()
	require.NoError(t, err, "Test data setup returns without error")
	defer func() {
		err = th.Teardown()
		require.NoError(t, err, "Test data teardown returns without error")
	}()

	// Setup: get a location link and locations for reference
	linkRec, err := th.Data.GetGameLocationLinkRecByRef("link-one-two")
	require.NoError(t, err, "GetGameLocationLinkRecByRef returns without error")
	fromLoc, err := th.Data.GetGameLocationRecByID(linkRec.FromGameLocationID)
	require.NoError(t, err, "GetGameLocationRecByID returns without error")
	toLoc, err := th.Data.GetGameLocationRecByID(linkRec.ToGameLocationID)
	require.NoError(t, err, "GetGameLocationRecByID returns without error")

	testCaseCollectionResponseDecoder := testCaseResponseDecoderGeneric[schema.GameLocationLinkCollectionResponse]
	testCaseResponseDecoder := testCaseResponseDecoderGeneric[schema.GameLocationLinkResponse]

	testCases := []struct {
		TestCase
		collectionRequest     bool
		collectionRecordCount int
	}{
		{
			TestCase: TestCase{
				Name: "API key with open access \\ get many location links \\ returns expected links",
				HandlerConfig: func(rnr *Runner) server.HandlerConfig {
					return rnr.HandlerConfig[getManyGameLocationLinks]
				},
				RequestQueryParams: func(d harness.Data) map[string]any {
					return map[string]any{
						"from_game_location_id": fromLoc.ID,
						"page_size":             10,
						"page_number":           1,
					}
				},
				ResponseDecoder: testCaseCollectionResponseDecoder,
				ResponseCode:    http.StatusOK,
			},
			collectionRequest:     true,
			collectionRecordCount: 1,
		},
		{
			TestCase: TestCase{
				Name: "API key with open access \\ get one location link with valid ID \\ returns expected link",
				HandlerConfig: func(rnr *Runner) server.HandlerConfig {
					return rnr.HandlerConfig[getOneGameLocationLink]
				},
				RequestPathParams: func(d harness.Data) map[string]string {
					return map[string]string{
						":location_link_id": linkRec.ID,
					}
				},
				ResponseDecoder: testCaseResponseDecoder,
				ResponseCode:    http.StatusOK,
			},
			collectionRequest: false,
		},
		{
			TestCase: TestCase{
				Name: "API key with open access \\ create location link with valid properties \\ returns created link",
				HandlerConfig: func(rnr *Runner) server.HandlerConfig {
					return rnr.HandlerConfig[createGameLocationLink]
				},
				RequestBody: func(d harness.Data) interface{} {
					return schema.GameLocationLinkRequest{
						FromGameLocationID: toLoc.ID,
						ToGameLocationID:   fromLoc.ID,
						Description:        "Test Link",
						Name:               "Test Link Name",
					}
				},
				ResponseDecoder: testCaseResponseDecoder,
				ResponseCode:    http.StatusCreated,
			},
			collectionRequest: false,
		},
		{
			TestCase: TestCase{
				Name: "API key with open access \\ delete location link with valid ID \\ returns no content",
				HandlerConfig: func(rnr *Runner) server.HandlerConfig {
					return rnr.HandlerConfig[deleteGameLocationLink]
				},
				RequestPathParams: func(d harness.Data) map[string]string {
					return map[string]string{
						":location_link_id": linkRec.ID,
					}
				},
				ResponseCode: http.StatusNoContent,
			},
			collectionRequest: false,
		},
	}

	for _, testCase := range testCases {
		t.Logf("Running test >%s<", testCase.Name)

		t.Run(testCase.Name, func(t *testing.T) {
			testFunc := func(method string, body interface{}) {
				if testCase.TestResponseCode() == http.StatusNoContent {
					// No content expected
					return
				}

				require.NotNil(t, body, "Response body is not nil")

				var responses []*schema.GameLocationLinkResponseData
				if testCase.collectionRequest {
					responses = body.(schema.GameLocationLinkCollectionResponse).Data
				} else {
					responses = append(responses, body.(schema.GameLocationLinkResponse).Data)
				}

				if testCase.collectionRequest {
					require.Equal(t, testCase.collectionRecordCount, len(responses), "Response record count length equals expected")
				}

				for _, d := range responses {
					require.NotEmpty(t, d.ID, "GameLocationLink ID is not empty")
					require.NotEmpty(t, d.FromGameLocationID, "FromGameLocationID is not empty")
					require.NotEmpty(t, d.ToGameLocationID, "ToGameLocationID is not empty")
					require.NotEmpty(t, d.Description, "Description is not empty")
				}
			}

			RunTestCase(t, th, &testCase.TestCase, testFunc)
		})
	}
}
