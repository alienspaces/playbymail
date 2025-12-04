package game_test

import (
	"bytes"
	"image"
	"image/png"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	corerecord "gitlab.com/alienspaces/playbymail/core/record"
	"gitlab.com/alienspaces/playbymail/core/server"
	"gitlab.com/alienspaces/playbymail/internal/harness"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
	"gitlab.com/alienspaces/playbymail/internal/runner/server/game"
	"gitlab.com/alienspaces/playbymail/internal/utils/testutil"
	"gitlab.com/alienspaces/playbymail/schema/api/game_schema"
)

// createTestImage creates a simple PNG image for testing
func createTestImage(width, height int) []byte {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	var buf bytes.Buffer
	png.Encode(&buf, img)
	return buf.Bytes()
}

func Test_uploadGameTurnSheetImageHandler(t *testing.T) {
	t.Parallel()

	th := testutil.NewTestHarness(t)
	require.NotNil(t, th, "newTestHarness returns without error")

	_, err := th.Setup()
	require.NoError(t, err, "Test data setup returns without error")

	defer func() {
		err = th.Teardown()
		require.NoError(t, err, "Test data teardown returns without error")
	}()

	gameRec, err := th.Data.GetGameRecByRef(harness.GameOneRef)
	require.NoError(t, err, "GetGameRecByRef returns without error")

	testImageData := createTestImage(2480, 3508) // A4 @ 300 DPI

	type testCase struct {
		testutil.TestCase
		expectResponse func(d harness.Data) game_schema.GameImageResponse
	}

	testCaseResponseDecoder := testutil.TestCaseResponseDecoderGeneric[game_schema.GameImageResponse]

	testCases := []testCase{
		{
			TestCase: testutil.TestCase{
				Name: "upload turn sheet image with default turn_sheet_type returns created image",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[game.UploadGameTurnSheetImage]
				},
				RequestPathParams: func(d harness.Data) map[string]string {
					return map[string]string{
						":game_id": gameRec.ID,
					}
				},
				RequestForms: func(d harness.Data) map[string]any {
					return map[string]any{
						"image": testImageData,
					}
				},
				ResponseDecoder: testCaseResponseDecoder,
				ResponseCode:    http.StatusCreated,
				ShouldTxCommit:  true, // Need to commit to verify image was saved
			},
			expectResponse: func(d harness.Data) game_schema.GameImageResponse {
				return game_schema.GameImageResponse{
					Data: &game_schema.GameImageResponseData{
						GameID:   gameRec.ID,
						Type:     game_record.GameImageTypeTurnSheetBackground,
						MimeType: game_record.GameImageMimeTypePNG,
						FileSize: len(testImageData),
						Width:    2480,
						Height:   3508,
					},
				}
			},
		},
		{
			TestCase: testutil.TestCase{
				Name: "upload turn sheet image with explicit turn_sheet_type returns created image",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[game.UploadGameTurnSheetImage]
				},
				RequestPathParams: func(d harness.Data) map[string]string {
					return map[string]string{
						":game_id": gameRec.ID,
					}
				},
				RequestQueryParams: func(d harness.Data) map[string]any {
					return map[string]any{
						"turn_sheet_type": adventure_game_record.AdventureGameTurnSheetTypeInventoryManagement,
					}
				},
				RequestForms: func(d harness.Data) map[string]any {
					return map[string]any{
						"image": testImageData,
					}
				},
				ResponseDecoder: testCaseResponseDecoder,
				ResponseCode:    http.StatusCreated,
				ShouldTxCommit:  true, // Need to commit to verify image was saved
			},
			expectResponse: func(d harness.Data) game_schema.GameImageResponse {
				return game_schema.GameImageResponse{
					Data: &game_schema.GameImageResponseData{
						GameID:   gameRec.ID,
						Type:     game_record.GameImageTypeTurnSheetBackground,
						MimeType: game_record.GameImageMimeTypePNG,
						FileSize: len(testImageData),
						Width:    2480,
						Height:   3508,
					},
				}
			},
		},
		{
			TestCase: testutil.TestCase{
				Name: "upload turn sheet image with invalid game ID returns not found error",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[game.UploadGameTurnSheetImage]
				},
				RequestPathParams: func(d harness.Data) map[string]string {
					return map[string]string{
						":game_id": "00000000-0000-0000-0000-000000000000",
					}
				},
				RequestForms: func(d harness.Data) map[string]any {
					return map[string]any{
						"image": testImageData,
					}
				},
				ResponseCode: http.StatusNotFound,
			},
		},
		{
			TestCase: testutil.TestCase{
				Name: "upload turn sheet image without image file returns invalid data error",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[game.UploadGameTurnSheetImage]
				},
				RequestPathParams: func(d harness.Data) map[string]string {
					return map[string]string{
						":game_id": gameRec.ID,
					}
				},
				RequestForms: func(d harness.Data) map[string]any {
					return map[string]any{}
				},
				ResponseCode: http.StatusBadRequest,
			},
		},
	}

	for _, testCase := range testCases {
		t.Logf("Running test >%s<", testCase.Name)

		t.Run(testCase.Name, func(t *testing.T) {
			testFunc := func(method string, body interface{}) {
				if testCase.TestResponseCode() != http.StatusCreated {
					return
				}

				require.NotNil(t, body, "Response body is not nil")

				aResp := body.(game_schema.GameImageResponse).Data
				xResp := testCase.expectResponse(th.Data).Data

				require.NotEmpty(t, aResp.ID, "Response ID is not empty")
				require.Equal(t, xResp.GameID, aResp.GameID, "Game ID matches")
				require.Equal(t, xResp.Type, aResp.Type, "Type matches")
				require.Equal(t, xResp.MimeType, aResp.MimeType, "MIME type matches")
				require.Equal(t, xResp.FileSize, aResp.FileSize, "File size matches")
				require.Equal(t, xResp.Width, aResp.Width, "Width matches")
				require.Equal(t, xResp.Height, aResp.Height, "Height matches")

				// Add image to teardown data so it gets cleaned up properly
				imgRec := &game_record.GameImage{
					Record:        corerecord.Record{ID: aResp.ID},
					GameID:        aResp.GameID,
					Type:          aResp.Type,
					TurnSheetType: aResp.TurnSheetType,
				}
				th.AddGameImageRecToTeardown(imgRec)
			}

			testutil.RunTestCase(t, th, &testCase, testFunc)
		})
	}
}

func Test_getGameTurnSheetImageHandler(t *testing.T) {
	t.Parallel()

	th := testutil.NewTestHarness(t)
	require.NotNil(t, th, "newTestHarness returns without error")

	_, err := th.Setup()
	require.NoError(t, err, "Test data setup returns without error")

	defer func() {
		err = th.Teardown()
		require.NoError(t, err, "Test data teardown returns without error")
	}()

	gameRec, err := th.Data.GetGameRecByRef(harness.GameOneRef)
	require.NoError(t, err, "GetGameRecByRef returns without error")

	// Get preconfigured images from harness
	joinGameImageRec, err := th.Data.GetGameImageRecByRef(harness.GameImageJoinGameRef)
	require.NoError(t, err, "GetGameImageRecByRef returns join game image")

	inventoryImageRec, err := th.Data.GetGameImageRecByRef(harness.GameImageInventoryRef)
	require.NoError(t, err, "GetGameImageRecByRef returns inventory image")

	type testCase struct {
		testutil.TestCase
		expectResponse func(d harness.Data) game_schema.GameImageCollectionResponse
	}

	testCaseResponseDecoder := testutil.TestCaseResponseDecoderGeneric[game_schema.GameImageCollectionResponse]

	testCases := []testCase{
		{
			TestCase: testutil.TestCase{
				Name: "get turn sheet images without turn_sheet_type returns all images",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[game.GetManyGameTurnSheetImages]
				},
				RequestPathParams: func(d harness.Data) map[string]string {
					return map[string]string{
						":game_id": gameRec.ID,
					}
				},
				ResponseDecoder: testCaseResponseDecoder,
				ResponseCode:    http.StatusOK,
			},
			expectResponse: func(d harness.Data) game_schema.GameImageCollectionResponse {
				return game_schema.GameImageCollectionResponse{
					Data: []*game_schema.GameImageResponseData{
						{
							GameID:        joinGameImageRec.GameID,
							Type:          joinGameImageRec.Type,
							TurnSheetType: joinGameImageRec.TurnSheetType,
							MimeType:      joinGameImageRec.MimeType,
							FileSize:      joinGameImageRec.FileSize,
							Width:         joinGameImageRec.Width,
							Height:        joinGameImageRec.Height,
						},
						{
							GameID:        inventoryImageRec.GameID,
							Type:          inventoryImageRec.Type,
							TurnSheetType: inventoryImageRec.TurnSheetType,
							MimeType:      inventoryImageRec.MimeType,
							FileSize:      inventoryImageRec.FileSize,
							Width:         inventoryImageRec.Width,
							Height:        inventoryImageRec.Height,
						},
					},
				}
			},
		},
		{
			TestCase: testutil.TestCase{
				Name: "get turn sheet image with explicit turn_sheet_type returns image",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[game.GetManyGameTurnSheetImages]
				},
				RequestPathParams: func(d harness.Data) map[string]string {
					return map[string]string{
						":game_id": gameRec.ID,
					}
				},
				RequestQueryParams: func(d harness.Data) map[string]any {
					return map[string]any{
						"turn_sheet_type": adventure_game_record.AdventureGameTurnSheetTypeJoinGame,
					}
				},
				ResponseDecoder: testCaseResponseDecoder,
				ResponseCode:    http.StatusOK,
			},
			expectResponse: func(d harness.Data) game_schema.GameImageCollectionResponse {
				return game_schema.GameImageCollectionResponse{
					Data: []*game_schema.GameImageResponseData{
						{
							GameID:        joinGameImageRec.GameID,
							Type:          joinGameImageRec.Type,
							TurnSheetType: joinGameImageRec.TurnSheetType,
							MimeType:      joinGameImageRec.MimeType,
							FileSize:      joinGameImageRec.FileSize,
							Width:         joinGameImageRec.Width,
							Height:        joinGameImageRec.Height,
						},
					},
				}
			},
		},
		{
			TestCase: testutil.TestCase{
				Name: "get turn sheet image for inventory turn_sheet_type returns image",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[game.GetManyGameTurnSheetImages]
				},
				RequestPathParams: func(d harness.Data) map[string]string {
					return map[string]string{
						":game_id": gameRec.ID,
					}
				},
				RequestQueryParams: func(d harness.Data) map[string]any {
					return map[string]any{
						"turn_sheet_type": adventure_game_record.AdventureGameTurnSheetTypeInventoryManagement,
					}
				},
				ResponseDecoder: testCaseResponseDecoder,
				ResponseCode:    http.StatusOK,
			},
			expectResponse: func(d harness.Data) game_schema.GameImageCollectionResponse {
				return game_schema.GameImageCollectionResponse{
					Data: []*game_schema.GameImageResponseData{
						{
							GameID:        inventoryImageRec.GameID,
							Type:          inventoryImageRec.Type,
							TurnSheetType: inventoryImageRec.TurnSheetType,
							MimeType:      inventoryImageRec.MimeType,
							FileSize:      inventoryImageRec.FileSize,
							Width:         inventoryImageRec.Width,
							Height:        inventoryImageRec.Height,
						},
					},
				}
			},
		},
		{
			TestCase: testutil.TestCase{
				Name: "get turn sheet image for non-existent turn_sheet_type returns no image",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[game.GetManyGameTurnSheetImages]
				},
				RequestPathParams: func(d harness.Data) map[string]string {
					return map[string]string{
						":game_id": gameRec.ID,
					}
				},
				RequestQueryParams: func(d harness.Data) map[string]any {
					return map[string]any{
						"turn_sheet_type": adventure_game_record.AdventureGameTurnSheetTypeLocationChoice,
					}
				},
				ResponseDecoder: testCaseResponseDecoder,
				ResponseCode:    http.StatusOK,
			},
			expectResponse: func(d harness.Data) game_schema.GameImageCollectionResponse {
				return game_schema.GameImageCollectionResponse{
					Data: []*game_schema.GameImageResponseData{},
				}
			},
		},
	}

	for _, testCase := range testCases {
		t.Logf("Running test >%s<", testCase.Name)

		t.Run(testCase.Name, func(t *testing.T) {
			testFunc := func(method string, body interface{}) {
				require.NotNil(t, body, "Response body is not nil")

				aResp := body.(game_schema.GameImageCollectionResponse)
				xResp := testCase.expectResponse(th.Data)

				require.Equal(t, len(xResp.Data), len(aResp.Data), "Number of images matches")
				if len(xResp.Data) > 0 {
					require.NotEmpty(t, aResp.Data, "Response contains images")
					require.Equal(t, len(xResp.Data), len(aResp.Data), "Number of images matches")
					// Check that all expected images are present (order may vary)
					for _, expectedImg := range xResp.Data {
						found := false
						for _, actualImg := range aResp.Data {
							if expectedImg.TurnSheetType == actualImg.TurnSheetType {
								require.Equal(t, expectedImg.Type, actualImg.Type, "Type matches")
								require.Equal(t, expectedImg.MimeType, actualImg.MimeType, "MIME type matches")
								require.Equal(t, expectedImg.FileSize, actualImg.FileSize, "File size matches")
								require.Equal(t, expectedImg.Width, actualImg.Width, "Width matches")
								require.Equal(t, expectedImg.Height, actualImg.Height, "Height matches")
								require.NotEmpty(t, actualImg.ID, "Image ID is not empty")
								found = true
								break
							}
						}
						require.True(t, found, "Expected image with turn_sheet_type %s found", expectedImg.TurnSheetType)
					}
				} else {
					require.Empty(t, aResp.Data, "Response contains no images")
				}
			}

			testutil.RunTestCase(t, th, &testCase, testFunc)
		})
	}
}

func Test_deleteGameTurnSheetImageHandler(t *testing.T) {
	t.Parallel()

	th := testutil.NewTestHarness(t)
	require.NotNil(t, th, "newTestHarness returns without error")

	_, err := th.Setup()
	require.NoError(t, err, "Test data setup returns without error")

	defer func() {
		err = th.Teardown()
		require.NoError(t, err, "Test data teardown returns without error")
	}()

	gameRec, err := th.Data.GetGameRecByRef(harness.GameOneRef)
	require.NoError(t, err, "GetGameRecByRef returns without error")

	// Get preconfigured join game image
	joinGameImageRec, err := th.Data.GetGameImageRecByRef(harness.GameImageJoinGameRef)
	require.NoError(t, err, "GetGameImageRecByRef returns join game image")

	testCases := []testutil.TestCase{
		{
			Name: "delete turn sheet image by ID returns no content",
			HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
				return rnr.GetHandlerConfig()[game.DeleteOneGameTurnSheetImage]
			},
			RequestPathParams: func(d harness.Data) map[string]string {
				// Use preconfigured join game image
				return map[string]string{
					":game_id":       gameRec.ID,
					":game_image_id": joinGameImageRec.ID,
				}
			},
			ResponseCode: http.StatusNoContent,
		},
		{
			Name: "delete turn sheet image for non-existent image ID returns not found",
			HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
				return rnr.GetHandlerConfig()[game.DeleteOneGameTurnSheetImage]
			},
			RequestPathParams: func(d harness.Data) map[string]string {
				return map[string]string{
					":game_id":       gameRec.ID,
					":game_image_id": "00000000-0000-0000-0000-000000000000",
				}
			},
			ResponseCode: http.StatusNotFound,
		},
	}

	for _, testCase := range testCases {
		t.Logf("Running test >%s<", testCase.Name)

		t.Run(testCase.Name, func(t *testing.T) {
			testutil.RunTestCase(t, th, &testCase, nil)
		})
	}
}
