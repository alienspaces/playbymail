package game_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"
	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/playbymail/core/nullstring"
	"gitlab.com/alienspaces/playbymail/core/server"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/core/type/storer"
	"gitlab.com/alienspaces/playbymail/internal/harness"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
	"gitlab.com/alienspaces/playbymail/internal/runner/server/game"
	"gitlab.com/alienspaces/playbymail/internal/turn_sheet"
	"gitlab.com/alienspaces/playbymail/internal/utils/testutil"
	"gitlab.com/alienspaces/playbymail/internal/utils/turnsheet"
)

func Test_uploadTurnSheetHandler(t *testing.T) {
	t.Parallel()

	// Create a custom harness config with turn sheets for existing players
	testDataConfig := harness.DefaultDataConfig()

	// Add a turn sheet for existing player in game instance one
	// Construct sheet data using LocationChoiceData struct
	locationChoiceData := turn_sheet.LocationChoiceData{
		LocationName:        "Test Location",
		LocationDescription: "A test location for turn sheet testing",
		LocationOptions: []turn_sheet.LocationOption{
			{
				LocationID:              "location_1",
				LocationLinkName:        "Location One",
				LocationLinkDescription: "First location option",
			},
			{
				LocationID:              "location_2",
				LocationLinkName:        "Location Two",
				LocationLinkDescription: "Second location option",
			},
		},
	}
	sheetDataBytes, err2 := json.Marshal(locationChoiceData)
	require.NoError(t, err2, "marshal location choice data returns without error")
	testDataConfig.GameConfigs[0].GameInstanceConfigs[0].AdventureGameTurnSheetConfigs = []harness.AdventureGameTurnSheetConfig{
		{
			GameTurnSheetConfig: harness.GameTurnSheetConfig{
				Reference:        harness.GameTurnSheetOneRef,
				AccountRef:       harness.AccountOneRef,
				TurnNumber:       1,
				SheetType:        adventure_game_record.AdventureSheetTypeLocationChoice,
				SheetOrder:       1,
				SheetData:        string(sheetDataBytes),
				ProcessingStatus: game_record.TurnSheetProcessingStatusPending,
			},
			GameCharacterInstanceRef: harness.GameCharacterInstanceOneRef,
		},
	}

	// Add a second game instance with "started" status for join game test with started game
	now := time.Now()

	// Create a new slice with existing instances plus the new one
	existingInstances := testDataConfig.GameConfigs[0].GameInstanceConfigs
	testDataConfig.GameConfigs[0].GameInstanceConfigs = make([]harness.GameInstanceConfig, 0, len(existingInstances)+1)
	testDataConfig.GameConfigs[0].GameInstanceConfigs = append(testDataConfig.GameConfigs[0].GameInstanceConfigs, existingInstances...)
	testDataConfig.GameConfigs[0].GameInstanceConfigs = append(testDataConfig.GameConfigs[0].GameInstanceConfigs, harness.GameInstanceConfig{
		Reference: harness.GameInstanceTwoRef,
		Record: &game_record.GameInstance{
			Status:    game_record.GameInstanceStatusStarted,
			StartedAt: sql.NullTime{Time: now, Valid: true},
		},
	})

	th := testutil.NewTestHarnessWithConfig(t, testDataConfig)
	require.NotNil(t, th, "newTestHarness returns without error")

	_, err := th.Setup()
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
				NewRunner: func(l logger.Logger, s storer.Storer, j *river.Client[pgx.Tx], d harness.Data) (testutil.TestRunnerer, error) {
					rnr, err := testutil.NewTestRunner(l, s, j)
					if err != nil {
						return nil, err
					}

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

					// Override runner methods to return mock data using harness data
					rnr.GetTurnSheetCodeFromImageFunc = func(ctx context.Context, imageData []byte) (string, error) {
						return turnSheetCode, nil
					}
					rnr.GetTurnSheetScanDataFunc = func(ctx context.Context, sheetType string, sheetData []byte, imageData []byte) ([]byte, error) {
						mockData := map[string]any{
							"choices": []string{"location_1"},
						}
						return json.Marshal(mockData)
					}
					return rnr, nil
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
				NewRunner: func(l logger.Logger, s storer.Storer, j *river.Client[pgx.Tx], d harness.Data) (testutil.TestRunnerer, error) {
					rnr, err := testutil.NewTestRunner(l, s, j)
					if err != nil {
						return nil, err
					}

					// Get game from harness (game instance status is "created" by default)
					gameID, ok := d.Refs.GameRefs[harness.GameOneRef]
					if !ok {
						return nil, fmt.Errorf("game ref %s not found", harness.GameOneRef)
					}

					// Generate a valid join game turn sheet code (no turn sheet record needed)
					turnSheetCode, err := turnsheet.GenerateJoinTurnSheetCode(gameID)
					if err != nil {
						return nil, fmt.Errorf("failed to generate join turn sheet code: %w", err)
					}

					// Override runner methods to return mock data
					rnr.GetTurnSheetCodeFromImageFunc = func(ctx context.Context, imageData []byte) (string, error) {
						return turnSheetCode, nil
					}
					rnr.GetTurnSheetScanDataFunc = func(ctx context.Context, sheetType string, sheetData []byte, imageData []byte) ([]byte, error) {
						// Mock join game scan data for new player
						mockData := map[string]any{
							"email": "newplayer@example.com",
							"name":  "New Player",
						}
						return json.Marshal(mockData)
					}
					return rnr, nil
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
				return body.TurnSheetID != "" && body.SheetType == adventure_game_record.AdventureSheetTypeJoinGame
			},
		},
		{
			TestCase: testutil.TestCase{
				Name: "existing game that has started and new player \\ upload join game turn sheet \\ returns processed join game turn sheet",
				NewRunner: func(l logger.Logger, s storer.Storer, j *river.Client[pgx.Tx], d harness.Data) (testutil.TestRunnerer, error) {
					rnr, err := testutil.NewTestRunner(l, s, j)
					if err != nil {
						return nil, err
					}

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

					// Generate a valid join game turn sheet code
					turnSheetCode, err := turnsheet.GenerateJoinTurnSheetCode(gameID)
					if err != nil {
						return nil, fmt.Errorf("failed to generate join turn sheet code: %w", err)
					}

					// Override runner methods to return mock data
					rnr.GetTurnSheetCodeFromImageFunc = func(ctx context.Context, imageData []byte) (string, error) {
						return turnSheetCode, nil
					}
					rnr.GetTurnSheetScanDataFunc = func(ctx context.Context, sheetType string, sheetData []byte, imageData []byte) ([]byte, error) {
						// Mock join game scan data for new player
						mockData := map[string]any{
							"email": "newplayer2@example.com",
							"name":  "New Player 2",
						}
						return json.Marshal(mockData)
					}
					return rnr, nil
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
				return body.TurnSheetID != "" && body.SheetType == adventure_game_record.AdventureSheetTypeJoinGame
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
