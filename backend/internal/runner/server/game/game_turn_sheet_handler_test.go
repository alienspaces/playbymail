package game_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"
	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/playbymail/core/nullstring"
	"gitlab.com/alienspaces/playbymail/core/nulltime"
	"gitlab.com/alienspaces/playbymail/core/server"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/core/type/storer"
	"gitlab.com/alienspaces/playbymail/internal/harness"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
	"gitlab.com/alienspaces/playbymail/internal/runner/server/game"
	"gitlab.com/alienspaces/playbymail/internal/turn_sheet"
	"gitlab.com/alienspaces/playbymail/internal/utils/config"
	"gitlab.com/alienspaces/playbymail/internal/utils/testutil"
	"gitlab.com/alienspaces/playbymail/internal/utils/turnsheet"
)

func Test_uploadTurnSheetHandler(t *testing.T) {
	t.Parallel()

	// Create a custom harness config
	testDataConfig := harness.DefaultDataConfig()

	// Configure expected turn sheets so references can be resolved after job workers run
	err := testDataConfig.ReplaceGameTurnConfigs(harness.GameInstanceOneRef, []harness.GameTurnConfig{
		{
			TurnNumber: 1,
			AdventureGameTurnSheetConfigs: []harness.AdventureGameTurnSheetConfig{
				{
					GameTurnSheetConfig: harness.GameTurnSheetConfig{
						Reference:        harness.GameTurnSheetOneRef,
						AccountRef:       harness.AccountOneRef,
						SheetType:        adventure_game_record.AdventureGameTurnSheetTypeLocationChoice,
						ProcessingStatus: game_record.TurnSheetProcessingStatusPending,
					},
					GameCharacterInstanceRef: harness.GameCharacterInstanceOneRef,
				},
			},
		},
	})
	require.NoError(t, err, "ReplaceAdventureGameTurnSheetConfigs returns without error")

	// Add a second game instance with "started" status for join game test with started game
	now := time.Now()

	err = testDataConfig.AppendGameInstanceConfigs(harness.GameOneRef, []harness.GameInstanceConfig{
		{
			Reference: harness.GameInstanceTwoRef,
			Record: &game_record.GameInstance{
				Status:    game_record.GameInstanceStatusStarted,
				StartedAt: nulltime.FromTime(now),
			},
		},
	})
	require.NoError(t, err, "AppendGameInstanceConfigs returns without error")

	th := testutil.NewTestHarnessWithConfig(t, testDataConfig)
	require.NotNil(t, th, "newTestHarness returns without error")

	_, err = th.Setup()
	require.NoError(t, err, "Test data setup returns without error")
	defer func() {
		err = th.Teardown()
		require.NoError(t, err, "Test data teardown returns without error")
	}()

	type testCase struct {
		testutil.TestCase
		expectResponse func(d harness.Data, body game.TurnSheetUploadResponse) bool
	}

	testCaseResponseDecoder := testutil.TestCaseResponseDecoderGeneric[game.TurnSheetUploadResponse]

	testCases := []testCase{
		{
			TestCase: testutil.TestCase{
				Name: "authenticated user \\ upload empty image \\ returns error",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[game.UploadTurnSheet]
				},
				RequestBody: func(d harness.Data) any {
					return []byte{}
				},
				RequestHeaders: func(d harness.Data) map[string]string {
					return map[string]string{
						"Content-Type": "image/jpeg",
					}
				},
				ResponseDecoder: testCaseResponseDecoder,
				ResponseCode:    http.StatusBadRequest,
			},
			expectResponse: func(d harness.Data, body game.TurnSheetUploadResponse) bool {
				return true // Error case, response structure may vary
			},
		},
		{
			TestCase: testutil.TestCase{
				Name: "existing game with existing player \\ upload location choice turn sheet \\ returns processed turn sheet",
				NewRunner: func(cfg config.Config, l logger.Logger, s storer.Storer, j *river.Client[pgx.Tx], scanner turn_sheet.TurnSheetScanner, d harness.Data) (testutil.TestRunnerer, error) {
					// Get turn sheet from harness
					turnSheetRec, err := d.GetGameTurnSheetRecByRef(harness.GameTurnSheetOneRef)
					if err != nil {
						return nil, fmt.Errorf("failed to get turn sheet from harness: %w", err)
					}

					gameID := turnSheetRec.GameID
					gameInstanceID := nullstring.ToString(turnSheetRec.GameInstanceID)
					accountID := turnSheetRec.AccountID

					// Generate a valid turn sheet code using harness turn sheet
					turnSheetCode, err := turnsheet.GenerateTurnSheetCode(gameID, gameInstanceID, accountID, turnSheetRec.ID)
					if err != nil {
						return nil, fmt.Errorf("failed to generate turn sheet code: %w", err)
					}

					// Create mock scanner with test data
					mockScanner := &testutil.MockTurnSheetScanner{
						GetTurnSheetCodeFromImageFunc: func(ctx context.Context, l logger.Logger, imageData []byte) (string, error) {
							return turnSheetCode, nil
						},
						GetTurnSheetScanDataFunc: func(ctx context.Context, l logger.Logger, sheetType string, sheetData []byte, imageData []byte) ([]byte, error) {
							mockData := map[string]any{
								"choices": []string{"location_1"},
							}
							return json.Marshal(mockData)
						},
					}

					return testutil.NewTestRunnerWithTurnSheetScanner(cfg, l, s, j, mockScanner)
				},
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[game.UploadTurnSheet]
				},
				RequestBody: func(d harness.Data) any {
					// For mocked tests, any image data works
					return []byte("fake-image-data")
				},
				RequestHeaders: func(d harness.Data) map[string]string {
					return map[string]string{
						"Content-Type": "image/jpeg",
					}
				},
				ResponseDecoder: testCaseResponseDecoder,
				ResponseCode:    http.StatusOK,
			},
			expectResponse: func(d harness.Data, body game.TurnSheetUploadResponse) bool {
				// Response should have turn sheet ID and scanned data
				return body.TurnSheetID != "" && body.SheetType != ""
			},
		},
		{
			TestCase: testutil.TestCase{
				Name: "existing game that has not started and new player \\ upload join game turn sheet \\ returns processed join game turn sheet",
				NewRunner: func(cfg config.Config, l logger.Logger, s storer.Storer, j *river.Client[pgx.Tx], scanner turn_sheet.TurnSheetScanner, d harness.Data) (testutil.TestRunnerer, error) {
					// Get game from harness (game instance status is "created" by default)
					gameID, ok := d.Refs.GameRefs[harness.GameOneRef]
					if !ok {
						return nil, fmt.Errorf("game ref %s not found", harness.GameOneRef)
					}

					// Get manager subscription for this game
					managerSubscriptionID, ok := d.Refs.GameSubscriptionRefs[harness.GameSubscriptionManagerOneRef]
					if !ok {
						return nil, fmt.Errorf("manager subscription ref %s not found", harness.GameSubscriptionManagerOneRef)
					}

					// Generate a valid join game turn sheet code (no turn sheet record needed)
					turnSheetCode, err := turnsheet.GenerateJoinTurnSheetCode(gameID, managerSubscriptionID)
					if err != nil {
						return nil, fmt.Errorf("failed to generate join turn sheet code: %w", err)
					}

					// Create mock scanner with test data
					mockScanner := &testutil.MockTurnSheetScanner{
						GetTurnSheetCodeFromImageFunc: func(ctx context.Context, l logger.Logger, imageData []byte) (string, error) {
							return turnSheetCode, nil
						},
						GetTurnSheetScanDataFunc: func(ctx context.Context, l logger.Logger, sheetType string, sheetData []byte, imageData []byte) ([]byte, error) {
							// Mock join game scan data for new player
							mockData := map[string]any{
								"email":                "newplayer@example.com",
								"name":                 "New Player",
								"postal_address_line1": "123 Huntsmans Road",
								"state_province":       "VIC",
								"country":              "Australia",
								"postal_code":          "12345",
							}
							return json.Marshal(mockData)
						},
					}

					return testutil.NewTestRunnerWithTurnSheetScanner(cfg, l, s, j, mockScanner)
				},
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[game.UploadTurnSheet]
				},
				RequestBody: func(d harness.Data) any {
					return []byte("fake-image-data")
				},
				RequestHeaders: func(d harness.Data) map[string]string {
					return map[string]string{
						"Content-Type": "image/jpeg",
					}
				},
				ResponseDecoder: testCaseResponseDecoder,
				ResponseCode:    http.StatusAccepted,
			},
			expectResponse: func(d harness.Data, body game.TurnSheetUploadResponse) bool {
				// Response should have turn sheet ID for join game
				return body.TurnSheetID != "" && body.SheetType == adventure_game_record.AdventureGameTurnSheetTypeJoinGame
			},
		},
		{
			TestCase: testutil.TestCase{
				Name: "existing game that has started and new player \\ upload join game turn sheet \\ returns processed join game turn sheet",
				NewRunner: func(cfg config.Config, l logger.Logger, s storer.Storer, j *river.Client[pgx.Tx], scanner turn_sheet.TurnSheetScanner, d harness.Data) (testutil.TestRunnerer, error) {
					// Get game from harness (use game instance two which has "started" status)
					gameID, ok := d.Refs.GameRefs[harness.GameOneRef]
					if !ok {
						return nil, fmt.Errorf("game ref %s not found", harness.GameOneRef)
					}

					// Verify game instance two is started
					gameInstanceRec, err := d.GetGameInstanceRecByRef(harness.GameInstanceTwoRef)
					if err != nil {
						return nil, fmt.Errorf("failed to get game instance two: %w", err)
					}
					if gameInstanceRec.Status != game_record.GameInstanceStatusStarted {
						return nil, fmt.Errorf("game instance two should have 'started' status, got '%s'", gameInstanceRec.Status)
					}

					// Get manager subscription for this game
					managerSubscriptionID, ok := d.Refs.GameSubscriptionRefs[harness.GameSubscriptionManagerOneRef]
					if !ok {
						return nil, fmt.Errorf("manager subscription ref %s not found", harness.GameSubscriptionManagerOneRef)
					}

					// Generate a valid join game turn sheet code
					turnSheetCode, err := turnsheet.GenerateJoinTurnSheetCode(gameID, managerSubscriptionID)
					if err != nil {
						return nil, fmt.Errorf("failed to generate join turn sheet code: %w", err)
					}

					// Create mock scanner with test data
					mockScanner := &testutil.MockTurnSheetScanner{
						GetTurnSheetCodeFromImageFunc: func(ctx context.Context, l logger.Logger, imageData []byte) (string, error) {
							return turnSheetCode, nil
						},
						GetTurnSheetScanDataFunc: func(ctx context.Context, l logger.Logger, sheetType string, sheetData []byte, imageData []byte) ([]byte, error) {
							// Mock join game scan data for new player
							mockData := map[string]any{
								"email":                "newplayer2@example.com",
								"name":                 "New Player 2",
								"postal_address_line1": "123 Huntsmans Road",
								"state_province":       "VIC",
								"country":              "Australia",
								"postal_code":          "12345",
							}
							return json.Marshal(mockData)
						},
					}

					return testutil.NewTestRunnerWithTurnSheetScanner(cfg, l, s, j, mockScanner)
				},
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[game.UploadTurnSheet]
				},
				RequestBody: func(d harness.Data) any {
					return []byte("fake-image-data")
				},
				RequestHeaders: func(d harness.Data) map[string]string {
					return map[string]string{
						"Content-Type": "image/jpeg",
					}
				},
				ResponseDecoder: testCaseResponseDecoder,
				ResponseCode:    http.StatusAccepted,
			},
			expectResponse: func(d harness.Data, body game.TurnSheetUploadResponse) bool {
				// Response should have turn sheet ID for join game
				return body.TurnSheetID != "" && body.SheetType == adventure_game_record.AdventureGameTurnSheetTypeJoinGame
			},
		},
	}

	for _, testCase := range testCases {
		t.Logf("Running test >%s<", testCase.Name)

		t.Run(testCase.Name, func(t *testing.T) {
			testFunc := func(method string, body interface{}) {
				if testCase.TestResponseCode() == http.StatusOK {
					require.NotNil(t, body, "Response body is not nil")
					response := body.(game.TurnSheetUploadResponse)
					require.True(t, testCase.expectResponse(th.Data, response), "Response matches expected structure")
					require.NotEmpty(t, response.TurnSheetID, "Turn sheet ID is not empty")
				}
			}

			testutil.RunTestCase(t, th, &testCase.TestCase, testFunc)
		})
	}
}

func Test_downloadJoinGameTurnSheetsHandler(t *testing.T) {
	t.Parallel()

	th := testutil.NewTestHarness(t)
	require.NotNil(t, th, "newTestHarness returns without error")

	_, err := th.Setup()
	require.NoError(t, err, "Test data setup returns without error")
	defer func() {
		err = th.Teardown()
		require.NoError(t, err, "Test data teardown returns without error")
	}()

	type testCase struct {
		testutil.TestCase
		expectPDFResponse func(d harness.Data, body []byte) bool
	}

	testCasePDFResponseDecoder := func(body io.Reader) (any, error) {
		pdfData, err := io.ReadAll(body)
		if err != nil {
			return nil, err
		}
		return pdfData, nil
	}

	testCases := []testCase{
		{
			TestCase: testutil.TestCase{
				Name: "authenticated manager \\ download join game turn sheet for adventure game \\ returns PDF",
				NewRunner: func(cfg config.Config, l logger.Logger, s storer.Storer, j *river.Client[pgx.Tx], scanner turn_sheet.TurnSheetScanner, d harness.Data) (testutil.TestRunnerer, error) {
					// Get manager account (AccountTwoRef has manager subscription)
					managerAccountRec, err := d.GetAccountRecByRef(harness.AccountTwoRef)
					if err != nil {
						return nil, fmt.Errorf("failed to get manager account: %w", err)
					}

					rnr, err := testutil.NewTestRunnerWithAccountID(cfg, l, s, j, scanner, managerAccountRec.ID, managerAccountRec.Email)
					if err != nil {
						return nil, err
					}

					return rnr, nil
				},
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[game.DownloadJoinGameTurnSheets]
				},
				RequestPathParams: func(d harness.Data) map[string]string {
					gameID, ok := d.Refs.GameRefs[harness.GameOneRef]
					require.True(t, ok, "game ref exists")
					return map[string]string{
						":game_id": gameID,
					}
				},
				ResponseDecoder: testCasePDFResponseDecoder,
				ResponseCode:    http.StatusOK,
			},
			expectPDFResponse: func(d harness.Data, body []byte) bool {
				// PDF should not be empty and should start with PDF magic bytes
				return len(body) > 0 && len(body) > 4 && string(body[0:4]) == "%PDF"
			},
		},
		{
			TestCase: testutil.TestCase{
				Name: "authenticated player \\ download join game turn sheet with game_subscription_id query param \\ returns PDF",
				NewRunner: func(cfg config.Config, l logger.Logger, s storer.Storer, j *river.Client[pgx.Tx], scanner turn_sheet.TurnSheetScanner, d harness.Data) (testutil.TestRunnerer, error) {
					// Use player account (AccountThreeRef)
					playerAccountRec, err := d.GetAccountRecByRef(harness.AccountThreeRef)
					if err != nil {
						return nil, fmt.Errorf("failed to get player account: %w", err)
					}

					rnr, err := testutil.NewTestRunnerWithAccountID(cfg, l, s, j, scanner, playerAccountRec.ID, playerAccountRec.Email)
					if err != nil {
						return nil, err
					}

					return rnr, nil
				},
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[game.DownloadJoinGameTurnSheets]
				},
				RequestPathParams: func(d harness.Data) map[string]string {
					gameID, ok := d.Refs.GameRefs[harness.GameOneRef]
					require.True(t, ok, "game ref exists")
					return map[string]string{
						":game_id": gameID,
					}
				},
				RequestQueryParams: func(d harness.Data) map[string]any {
					// Provide manager subscription ID as query parameter
					managerSubscriptionID, ok := d.Refs.GameSubscriptionRefs[harness.GameSubscriptionManagerOneRef]
					require.True(t, ok, "manager subscription ref exists")
					return map[string]any{
						"game_subscription_id": managerSubscriptionID,
					}
				},
				ResponseDecoder: testCasePDFResponseDecoder,
				ResponseCode:    http.StatusOK,
			},
			expectPDFResponse: func(d harness.Data, body []byte) bool {
				// PDF should not be empty and should start with PDF magic bytes
				return len(body) > 0 && len(body) > 4 && string(body[0:4]) == "%PDF"
			},
		},
		{
			TestCase: testutil.TestCase{
				Name: "authenticated non-manager \\ download join game turn sheet without game_subscription_id \\ returns error",
				NewRunner: func(cfg config.Config, l logger.Logger, s storer.Storer, j *river.Client[pgx.Tx], scanner turn_sheet.TurnSheetScanner, d harness.Data) (testutil.TestRunnerer, error) {
					// Use player account (AccountThreeRef) which is not a manager
					playerAccountRec, err := d.GetAccountRecByRef(harness.AccountThreeRef)
					if err != nil {
						return nil, fmt.Errorf("failed to get player account: %w", err)
					}

					rnr, err := testutil.NewTestRunnerWithAccountID(cfg, l, s, j, scanner, playerAccountRec.ID, playerAccountRec.Email)
					if err != nil {
						return nil, err
					}

					return rnr, nil
				},
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[game.DownloadJoinGameTurnSheets]
				},
				RequestPathParams: func(d harness.Data) map[string]string {
					gameID, ok := d.Refs.GameRefs[harness.GameOneRef]
					require.True(t, ok, "game ref exists")
					return map[string]string{
						":game_id": gameID,
					}
				},
				ResponseDecoder: testutil.TestCaseResponseDecoderGeneric[map[string]any],
				ResponseCode:    http.StatusBadRequest,
			},
			expectPDFResponse: func(d harness.Data, body []byte) bool {
				return true // Error case, not used for error responses
			},
		},
		{
			TestCase: testutil.TestCase{
				Name: "authenticated user \\ download join game turn sheet with invalid game_subscription_id \\ returns error",
				NewRunner: func(cfg config.Config, l logger.Logger, s storer.Storer, j *river.Client[pgx.Tx], scanner turn_sheet.TurnSheetScanner, d harness.Data) (testutil.TestRunnerer, error) {
					// Use any account
					accountRec, err := d.GetAccountRecByRef(harness.AccountOneRef)
					if err != nil {
						return nil, fmt.Errorf("failed to get account: %w", err)
					}

					rnr, err := testutil.NewTestRunnerWithAccountID(cfg, l, s, j, scanner, accountRec.ID, accountRec.Email)
					if err != nil {
						return nil, err
					}

					return rnr, nil
				},
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[game.DownloadJoinGameTurnSheets]
				},
				RequestPathParams: func(d harness.Data) map[string]string {
					gameID, ok := d.Refs.GameRefs[harness.GameOneRef]
					require.True(t, ok, "game ref exists")
					return map[string]string{
						":game_id": gameID,
					}
				},
				RequestQueryParams: func(d harness.Data) map[string]any {
					// Provide invalid subscription ID
					return map[string]any{
						"game_subscription_id": "00000000-0000-0000-0000-000000000000",
					}
				},
				ResponseDecoder: testutil.TestCaseResponseDecoderGeneric[map[string]any],
				ResponseCode:    http.StatusNotFound,
			},
			expectPDFResponse: func(d harness.Data, body []byte) bool {
				return true // Error case, not used for error responses
			},
		},
		{
			TestCase: testutil.TestCase{
				Name: "authenticated user \\ download join game turn sheet with game_subscription_id for wrong game \\ returns error",
				NewRunner: func(cfg config.Config, l logger.Logger, s storer.Storer, j *river.Client[pgx.Tx], scanner turn_sheet.TurnSheetScanner, d harness.Data) (testutil.TestRunnerer, error) {
					// Use any account
					accountRec, err := d.GetAccountRecByRef(harness.AccountOneRef)
					if err != nil {
						return nil, fmt.Errorf("failed to get account: %w", err)
					}

					rnr, err := testutil.NewTestRunnerWithAccountID(cfg, l, s, j, scanner, accountRec.ID, accountRec.Email)
					if err != nil {
						return nil, err
					}

					return rnr, nil
				},
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[game.DownloadJoinGameTurnSheets]
				},
				RequestPathParams: func(d harness.Data) map[string]string {
					gameID, ok := d.Refs.GameRefs[harness.GameOneRef]
					require.True(t, ok, "game ref exists")
					return map[string]string{
						":game_id": gameID,
					}
				},
				RequestQueryParams: func(d harness.Data) map[string]any {
					// Provide player subscription ID (not manager, and belongs to same game but wrong type)
					playerSubscriptionID, ok := d.Refs.GameSubscriptionRefs[harness.GameSubscriptionPlayerOneRef]
					require.True(t, ok, "player subscription ref exists")
					return map[string]any{
						"game_subscription_id": playerSubscriptionID,
					}
				},
				ResponseDecoder: testutil.TestCaseResponseDecoderGeneric[map[string]any],
				ResponseCode:    http.StatusBadRequest,
			},
			expectPDFResponse: func(d harness.Data, body []byte) bool {
				return true // Error case, not used for error responses
			},
		},
		{
			TestCase: testutil.TestCase{
				Name: "authenticated user \\ download join game turn sheet with invalid game ID \\ returns error",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[game.DownloadJoinGameTurnSheets]
				},
				RequestPathParams: func(d harness.Data) map[string]string {
					return map[string]string{
						":game_id": "00000000-0000-0000-0000-000000000000",
					}
				},
				ResponseDecoder: testutil.TestCaseResponseDecoderGeneric[map[string]any],
				ResponseCode:    http.StatusNotFound,
			},
			expectPDFResponse: func(d harness.Data, body []byte) bool {
				return true // Error case, not used for error responses
			},
		},
	}

	for _, testCase := range testCases {
		t.Logf("Running test >%s<", testCase.Name)

		t.Run(testCase.Name, func(t *testing.T) {
			testFunc := func(method string, body interface{}) {
				if testCase.TestResponseCode() == http.StatusOK {
					require.NotNil(t, body, "Response body is not nil")
					pdfData, ok := body.([]byte)
					require.True(t, ok, "Response body is []byte")
					require.True(t, testCase.expectPDFResponse(th.Data, pdfData), "Response matches expected structure")
					require.Greater(t, len(pdfData), 0, "PDF data is not empty")
					// Verify PDF magic bytes
					require.GreaterOrEqual(t, len(pdfData), 4, "PDF data is at least 4 bytes")
					require.Equal(t, "%PDF", string(pdfData[0:4]), "Response is a valid PDF")
				}
			}

			testutil.RunTestCase(t, th, &testCase.TestCase, testFunc)
		})
	}
}
