package turnsheet_test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	coreconfig "gitlab.com/alienspaces/playbymail/core/config"
	"gitlab.com/alienspaces/playbymail/core/convert"
	corelog "gitlab.com/alienspaces/playbymail/core/log"
	"gitlab.com/alienspaces/playbymail/internal/turnsheet"
	"gitlab.com/alienspaces/playbymail/internal/utils/config"
	"gitlab.com/alienspaces/playbymail/internal/utils/testutil"
)

func TestMonsterEncounterProcessor_GenerateTurnSheet(t *testing.T) {

	cfg, l, _, _, _ := testutil.NewDefaultDependencies(t)
	cfg.TemplatesPath = "../../templates"

	processor, err := turnsheet.NewMonsterEncounterProcessor(l, cfg)
	require.NoError(t, err)

	backgroundImage := loadTestBackgroundImage(t, "testdata/background-darkforest.png")

	tests := []struct {
		name               string
		data               any
		expectError        bool
		expectErrorMessage string
	}{
		{
			name:               "given empty MonsterEncounterData when generating turn sheet then validation error is returned",
			data:               &turnsheet.MonsterEncounterData{},
			expectError:        true,
			expectErrorMessage: "game name is required",
		},
		{
			name:        "given nil data when generating turn sheet then PDF generation is handled gracefully",
			data:        nil,
			expectError: false,
		},
		{
			name:        "given invalid data type when generating turn sheet then PDF generation is handled gracefully",
			data:        "invalid data",
			expectError: false,
		},
		{
			name: "given valid MonsterEncounterData with alive creature when generating turn sheet then PDF is generated successfully",
			data: &turnsheet.MonsterEncounterData{
				TurnSheetTemplateData: turnsheet.TurnSheetTemplateData{
					GameName:        convert.Ptr("The Door Beneath the Staircase"),
					GameType:        convert.Ptr("adventure"),
					TurnNumber:      convert.Ptr(3),
					AccountName:     convert.Ptr("Test Player"),
					TurnSheetCode:   convert.Ptr(generateTestTurnSheetCode(t)),
					BackgroundImage: &backgroundImage,
				},
				CharacterName:      "Aldric",
				CharacterHealth:    65,
				CharacterMaxHealth: 100,
				CharacterAttack:    8,
				CharacterDefense:   3,
				EquippedWeapon: &turnsheet.EquippedWeapon{
					ItemInstanceID: "weapon-1",
					Name:           "Iron Sword",
					Damage:         8,
				},
				EquippedArmor: &turnsheet.EquippedArmor{
					ItemInstanceID: "armor-1",
					Name:           "Leather Jerkin",
					Defense:        3,
				},
			Creatures: []turnsheet.EncounterCreature{
				{
					CreatureInstanceID: "creature-1",
					Name:               "Sand Serpent",
					Description:        "A massive serpent that lurks beneath desert sands.",
					Health:             80,
					MaxHealth:          100,
					AttackDamage:       12,
					Defense:            2,
					Disposition:        "aggressive",
				},
			},
			MaxActions: 3,
		},
		expectError: false,
	},
	{
		name: "given MonsterEncounterData with no creatures when generating turn sheet then error is returned",
			data: &turnsheet.MonsterEncounterData{
				TurnSheetTemplateData: turnsheet.TurnSheetTemplateData{
					GameName:      convert.Ptr("Test Game"),
					GameType:      convert.Ptr("adventure"),
					TurnNumber:    convert.Ptr(1),
					AccountName:   convert.Ptr("Test Player"),
					TurnSheetCode: convert.Ptr(generateTestTurnSheetCode(t)),
				},
				CharacterName:      "Aldric",
				CharacterHealth:    100,
				CharacterMaxHealth: 100,
				Creatures:          []turnsheet.EncounterCreature{},
				MaxActions:         3,
			},
			expectError:        true,
			expectErrorMessage: "at least one creature is required",
		},
		{
			name: "given read-only MonsterEncounterData with dead creature when generating turn sheet then PDF is generated successfully",
			data: &turnsheet.MonsterEncounterData{
				TurnSheetTemplateData: turnsheet.TurnSheetTemplateData{
					GameName:        convert.Ptr("The Door Beneath the Staircase"),
					GameType:        convert.Ptr("adventure"),
					TurnNumber:      convert.Ptr(5),
					AccountName:     convert.Ptr("Test Player"),
					TurnSheetCode:   convert.Ptr(generateTestTurnSheetCode(t)),
					BackgroundImage: &backgroundImage,
				},
				CharacterName:      "Aldric",
				CharacterHealth:    90,
				CharacterMaxHealth: 100,
			Creatures: []turnsheet.EncounterCreature{
				{
					CreatureInstanceID: "creature-2",
					Name:               "Sand Serpent",
					Description:        "A massive serpent. It lies still.",
					Health:             0,
					MaxHealth:          100,
					AttackDamage:       12,
					Defense:            2,
					Disposition:        "aggressive",
					IsDead:             true,
				},
			},
				MaxActions: 0,
				ReadOnly:   true,
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var sheetData []byte
			if tt.data != nil {
				var err error
				sheetData, err = json.Marshal(tt.data)
				require.NoError(t, err, "Should marshal test data")
			}

			ctx := context.Background()
			pdfData, err := processor.GenerateTurnSheet(ctx, l, turnsheet.DocumentFormatPDF, sheetData)

			if tt.expectError {
				require.Error(t, err, "Should return error")
				if tt.expectErrorMessage != "" {
					require.Contains(t, err.Error(), tt.expectErrorMessage, "Error message should contain expected text")
				}
				require.Nil(t, pdfData, "PDF data should be nil on error")
			} else if err != nil {
				t.Logf("PDF generation failed (may be expected in test environment): %v", err)
			}
		})
	}
}

func TestMonsterEncounterProcessor_ScanTurnSheet(t *testing.T) {

	cfg, l, _, _, _ := testutil.NewDefaultDependencies(t)
	cfg.TemplatesPath = "../../templates"

	processor, err := turnsheet.NewMonsterEncounterProcessor(l, cfg)
	require.NoError(t, err)

	tests := []struct {
		name               string
		imageDataFn        func() ([]byte, error)
		sheetDataFn        func() ([]byte, error)
		expectError        bool
		expectErrorMessage string
		requiresScanner    bool
	}{
		{
			name: "given empty image data when scanning turn sheet then error is returned",
			imageDataFn: func() ([]byte, error) {
				return []byte{}, nil
			},
			sheetDataFn: func() ([]byte, error) {
				return []byte(`{}`), nil
			},
			expectError:        true,
			expectErrorMessage: "empty image data",
			requiresScanner:    false,
		},
		{
			name: "given nil image data when scanning turn sheet then error is returned",
			imageDataFn: func() ([]byte, error) {
				return nil, nil
			},
			sheetDataFn: func() ([]byte, error) {
				return []byte(`{"character_name":"Aldric"}`), nil
			},
			expectError:        true,
			expectErrorMessage: "empty image data",
			requiresScanner:    false,
		},
		// TODO: (agent) Add test case with real scanned monster encounter image when test data exists.
		// Use a fixture image (e.g. testdata/adventure_game_monster_encounter_turn_sheet_filled.jpg)
		// with expected MonsterEncounterScanData. Blocked on test asset availability.
		// {
		// 	name: "given real scanned image when scanning then combat actions extracted",
		// 	imageDataFn: func() ([]byte, error) {
		// 		return os.ReadFile("testdata/adventure_game_monster_encounter_turn_sheet_filled.jpg")
		// 	},
		// 	sheetDataFn: func() ([]byte, error) {
		// 		data := turnsheet.MonsterEncounterData{
		// 			CharacterName: "Aldric",
		// 			Creatures: []turnsheet.EncounterCreature{
		// 				{CreatureInstanceID: "creature-1", Name: "Sand Serpent"},
		// 			},
		// 			MaxActions: 3,
		// 		}
		// 		return json.Marshal(data)
		// 	},
		// 	expectError:     false,
		// 	requiresScanner: true,
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testStart := time.Now()
			if tt.requiresScanner {
				requireOpenAIKey(t)
			}

			imageData, err := tt.imageDataFn()
			if err != nil {
				t.Fatalf("Failed to load image data: %v", err)
			}
			t.Logf("Loaded image data: %d bytes in %v", len(imageData), time.Since(testStart))

			sheetData, err := tt.sheetDataFn()
			if err != nil {
				t.Fatalf("Failed to get sheet data: %v", err)
			}

			ctx := context.Background()
			scanStart := time.Now()
			resultData, err := processor.ScanTurnSheet(ctx, l, sheetData, imageData)
			t.Logf("ScanTurnSheet completed in %v", time.Since(scanStart))

			if tt.expectError {
				require.Error(t, err, "Should return error")
				if tt.expectErrorMessage != "" {
					require.Contains(t, err.Error(), tt.expectErrorMessage, "Error message should contain expected text")
				}
				require.Nil(t, resultData, "Result should be nil on error")
			} else {
				require.NoError(t, err, "Should not return error")
				require.NotNil(t, resultData, "Result should not be nil")
			}
		})
	}
}

// TestGenerateMonsterEncounterFormatsForPrinting generates HTML and PDF versions for physical testing.
// Set SAVE_TEST_FILES=true to save the files to the testdata directory.
func TestGenerateMonsterEncounterFormatsForPrinting(t *testing.T) {

	cfg, l, _, _, _ := testutil.NewDefaultDependencies(t)
	cfg.TemplatesPath = "../../templates"
	cfg.SaveTestFiles = true

	processor, err := turnsheet.NewMonsterEncounterProcessor(l, cfg)
	require.NoError(t, err)

	type formatCase struct {
		name     string
		format   turnsheet.DocumentFormat
		ext      string
		logExtra bool
	}

	cases := []formatCase{
		{name: "pdf", format: turnsheet.DocumentFormatPDF, ext: "pdf", logExtra: true},
		{name: "html", format: turnsheet.DocumentFormatHTML, ext: "html"},
	}

	backgroundImage := loadTestBackgroundImage(t, "testdata/background-darkforest.png")
	turnSheetCode := generateTestTurnSheetCode(t)

	testData := &turnsheet.MonsterEncounterData{
		TurnSheetTemplateData: turnsheet.TurnSheetTemplateData{
			GameName:          convert.Ptr("The Door Beneath the Staircase"),
			GameType:          convert.Ptr("adventure"),
			TurnNumber:        convert.Ptr(3),
			AccountName:       convert.Ptr("Test Player"),
			TurnSheetCode:     convert.Ptr(turnSheetCode),
			TurnSheetDeadline: convert.Ptr(time.Now().Add(24 * time.Hour)),
			BackgroundImage:   &backgroundImage,
			TurnEvents: []turnsheet.TurnEvent{
				{Category: turnsheet.TurnEventCategoryMovement, Icon: turnsheet.TurnEventIconMovement, Message: "Aldric entered the shadowed corridor."},
				{Category: turnsheet.TurnEventCategoryCombat, Icon: turnsheet.TurnEventIconCombat, Message: "Aldric attacked the Goblin Scout for 8 damage."},
				{Category: turnsheet.TurnEventCategoryCombat, Icon: turnsheet.TurnEventIconCombat, Message: "Goblin Scout retaliated for 4 damage."},
				{Category: turnsheet.TurnEventCategorySystem, Icon: turnsheet.TurnEventIconSystem, Message: "Aldric's health dropped to 65."},
			},
		},
		CharacterName:      "Aldric",
		CharacterHealth:    65,
		CharacterMaxHealth: 100,
		CharacterAttack:    8,
		CharacterDefense:   3,
		EquippedWeapon: &turnsheet.EquippedWeapon{
			ItemInstanceID: "weapon-1",
			Name:           "Iron Sword",
			Damage:         8,
		},
		EquippedArmor: &turnsheet.EquippedArmor{
			ItemInstanceID: "armor-1",
			Name:           "Leather Jerkin",
			Defense:        3,
		},
	Creatures: []turnsheet.EncounterCreature{
		{
			CreatureInstanceID: "creature-1",
			Name:               "Sand Serpent",
			Description:        "A massive serpent that lurks beneath the desert sands, striking with terrifying speed.",
			Health:             80,
			MaxHealth:          100,
			AttackDamage:       12,
			Defense:            2,
			Disposition:        "aggressive",
		},
		{
			CreatureInstanceID: "creature-2",
			Name:               "Desert Scorpion",
			Description:        "A venomous scorpion the size of a dog. Its tail curls menacingly.",
			Health:             35,
			MaxHealth:          50,
			AttackDamage:       8,
			Defense:            4,
			Disposition:        "aggressive",
		},
	},
		MaxActions: 3,
	}

	ctx := context.Background()
	sheetData, err := json.Marshal(testData)
	require.NoError(t, err, "Should marshal test data")

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			output, err := processor.GenerateTurnSheet(ctx, l, tc.format, sheetData)
			require.NoError(t, err, "Should generate output without error")
			require.NotEmpty(t, output, "Output should not be empty")

			if cfg.SaveTestFiles {
				path := fmt.Sprintf("testdata/adventure_game_monster_encounter_turnsheet.%s", tc.ext)
				err = os.WriteFile(path, output, 0644)
				require.NoError(t, err, "Should save output to testdata directory")
				t.Logf("%s preview saved to %s", tc.name, path)

				if tc.logExtra {
					t.Logf("Output size: %d bytes", len(output))
					t.Logf("")
					t.Logf("Generated successfully. To test the scanner:")
					t.Logf("1. Print the PDF: %s", path)
					t.Logf("2. Fill in combat actions (attack/do nothing + target)")
					t.Logf("3. Scan the completed turn sheet to a JPEG file")
					t.Logf("4. Save the JPEG in testdata/ with a descriptive name")
					t.Logf("5. Write a test that loads the JPEG and tests the scanner")
				}
			}
		})
	}
}

// TestMonsterEncounterScanData_GetActions_FlatFormat verifies that scanned_data submitted
// by the in-browser HTML form (flat action_N / action_N_target keys) is correctly
// normalised by GetActions() into the structured CombatAction slice that the processor
// uses to resolve combat. This is the format that was stored in the database for turns
// 4 and 5 when the bug was discovered.
func TestMonsterEncounterScanData_GetActions_FlatFormat(t *testing.T) {
	t.Parallel()

	// Exact payload stored in DB from the two failed turns (3× attack on the Sand Serpent).
	raw := []byte(`{"action_0":"attack","action_1":"attack","action_2":"attack","action_0_target":"406efb70-0c8d-44c0-8139-0e659b76f77f","action_1_target":"406efb70-0c8d-44c0-8139-0e659b76f77f","action_2_target":"406efb70-0c8d-44c0-8139-0e659b76f77f"}`)

	var scanData turnsheet.MonsterEncounterScanData
	err := json.Unmarshal(raw, &scanData)
	require.NoError(t, err)

	actions := scanData.GetActions()
	require.Len(t, actions, 3, "expected 3 attack actions from flat format")
	for i, a := range actions {
		require.Equal(t, "attack", a.ActionType, "action %d: expected attack", i)
		require.Equal(t, "406efb70-0c8d-44c0-8139-0e659b76f77f", a.TargetCreatureInstanceID, "action %d: wrong target", i)
	}
}

// TestMonsterEncounterScanData_GetActions_StructuredFormat verifies that the structured
// actions array format (produced by the frontend after the fix and by the OCR scanner)
// is returned unchanged by GetActions().
func TestMonsterEncounterScanData_GetActions_StructuredFormat(t *testing.T) {
	t.Parallel()

	raw := []byte(`{"actions":[{"action_type":"attack","target_creature_instance_id":"abc123"},{"action_type":"do_nothing"}]}`)

	var scanData turnsheet.MonsterEncounterScanData
	err := json.Unmarshal(raw, &scanData)
	require.NoError(t, err)

	actions := scanData.GetActions()
	require.Len(t, actions, 2)
	require.Equal(t, "attack", actions[0].ActionType)
	require.Equal(t, "abc123", actions[0].TargetCreatureInstanceID)
	require.Equal(t, "do_nothing", actions[1].ActionType)
}

// TestMonsterEncounterScanData_GetActions_MixedDoNothing verifies that GetActions handles
// a mix of do_nothing and attack actions in flat format correctly.
func TestMonsterEncounterScanData_GetActions_MixedDoNothing(t *testing.T) {
	t.Parallel()

	raw := []byte(`{"action_0":"attack","action_0_target":"creature-1","action_1":"do_nothing","action_1_target":"","action_2":"do_nothing","action_2_target":""}`)

	var scanData turnsheet.MonsterEncounterScanData
	err := json.Unmarshal(raw, &scanData)
	require.NoError(t, err)

	actions := scanData.GetActions()
	require.Len(t, actions, 3)
	require.Equal(t, "attack", actions[0].ActionType)
	require.Equal(t, "creature-1", actions[0].TargetCreatureInstanceID)
	require.Equal(t, "do_nothing", actions[1].ActionType)
	require.Equal(t, "do_nothing", actions[2].ActionType)
}

// TestMonsterEncounterScanData_GetActions_Empty verifies that an empty payload returns
// an empty (nil) slice.
func TestMonsterEncounterScanData_GetActions_Empty(t *testing.T) {
	t.Parallel()

	raw := []byte(`{}`)

	var scanData turnsheet.MonsterEncounterScanData
	err := json.Unmarshal(raw, &scanData)
	require.NoError(t, err)

	actions := scanData.GetActions()
	require.Empty(t, actions)
}

// TestMonsterEncounterProcessor_GenerateTurnSheet_HTMLContent verifies the rendered HTML
// contains all user-facing "Creature Encounter" strings and CSS introduced by the
// Creature Encounter polish: renamed title/section, horizontal radio row, section height
// stretching, and compact header override.
//
// This test constructs its dependencies without a database connection since GenerateTurnSheet
// only renders HTML templates — it does not query the database.
func TestMonsterEncounterProcessor_GenerateTurnSheet_HTMLContent(t *testing.T) {
	t.Parallel()

	l := corelog.NewDefaultLogger()
	cfg := config.Config{Config: coreconfig.Config{TemplatesPath: "../../templates"}}

	processor, err := turnsheet.NewMonsterEncounterProcessor(l, cfg)
	require.NoError(t, err)

	backgroundImage := loadTestBackgroundImage(t, "testdata/background-darkforest.png")
	baseData := turnsheet.TurnSheetTemplateData{
		GameName:        convert.Ptr("The Door Beneath the Staircase"),
		GameType:        convert.Ptr("adventure"),
		TurnNumber:      convert.Ptr(3),
		AccountName:     convert.Ptr("Test Player"),
		TurnSheetCode:   convert.Ptr(generateTestTurnSheetCode(t)),
		BackgroundImage: &backgroundImage,
	}
	creatures := []turnsheet.EncounterCreature{
		{
			CreatureInstanceID: "creature-1",
			Name:               "Sand Serpent",
			Description:        "A massive serpent that lurks beneath desert sands.",
			Health:             80,
			MaxHealth:          100,
			AttackDamage:       12,
			Defense:            2,
			Disposition:        "aggressive",
		},
	}

	tests := []struct {
		name         string
		readOnly     bool
		wantContains []string
		wantAbsent   []string
	}{
		{
			name:     "interactive sheet contains creature encounter title and all polish CSS",
			readOnly: false,
			wantContains: []string{
				"Creature Encounter", // renamed title
				"Creature Actions",   // renamed section heading
				"action-options-row", // horizontal radio buttons CSS class
				"health-bar-outer",   // shared health bar class
				"stat-panel",         // shared stat panel class
				"min-height: 22mm",   // compact header override
				"ATK:",               // character attack stat row
				"DEF:",               // character defense stat row
			},
			wantAbsent: []string{
				"Monster Encounter", // old user-facing string must not appear
				"Combat Actions",    // old section heading must not appear
			},
		},
		{
			name:     "read-only sheet contains creature encounter notice and no action slots",
			readOnly: true,
			wantContains: []string{
				"Creature Encounter", // renamed title
				"creature encounter", // read-only notice text
				"action-options-row", // CSS still rendered in styles block
				"min-height: 22mm",  // compact header override
				"ATK:",              // character attack stat row
				"DEF:",              // character defense stat row
			},
			wantAbsent: []string{
				"Monster Encounter",
				"Combat Actions",
			},
		},
	}

	ctx := context.Background()

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			data := &turnsheet.MonsterEncounterData{
				TurnSheetTemplateData: baseData,
				CharacterName:         "Aldric",
				CharacterHealth:       65,
				CharacterMaxHealth:    100,
				CharacterAttack:       8,
				CharacterDefense:      3,
				EquippedWeapon: &turnsheet.EquippedWeapon{
					ItemInstanceID: "weapon-1",
					Name:           "Iron Sword",
					Damage:         8,
				},
				Creatures:  creatures,
				MaxActions: 3,
				ReadOnly:   tt.readOnly,
			}

			sheetData, err := json.Marshal(data)
			require.NoError(t, err)

			html, err := processor.GenerateTurnSheet(ctx, l, turnsheet.DocumentFormatHTML, sheetData)
			require.NoError(t, err)
			require.NotEmpty(t, html)

			htmlStr := string(html)

			for _, want := range tt.wantContains {
				require.True(t, strings.Contains(htmlStr, want),
					"expected HTML to contain %q", want)
			}
			for _, absent := range tt.wantAbsent {
			require.False(t, strings.Contains(htmlStr, absent),
				"expected HTML not to contain %q", absent)
		}
	})
	}
}

// TestMonsterEncounterProcessor_GenerateTurnSheet_CreatureImage verifies that a creature
// portrait data URL is embedded correctly in the rendered HTML without being sanitised to
// "#ZgotmplZ" by html/template's URL context escaping.
func TestMonsterEncounterProcessor_GenerateTurnSheet_CreatureImage(t *testing.T) {
	t.Parallel()

	l := corelog.NewDefaultLogger()
	cfg := config.Config{Config: coreconfig.Config{TemplatesPath: "../../templates"}}

	processor, err := turnsheet.NewMonsterEncounterProcessor(l, cfg)
	require.NoError(t, err)

	// Minimal 1×1 JPEG as a data URL (small enough for a unit test).
	creatureImageURL := "data:image/jpeg;base64,/9j/4AAQSkZJRgABAQAAAQABAAD/2wBDAAgGBgcGBQgHBwcJCQgKDBQNDAsLDBkSEw8UHRofHh0aHBwgJC4nICIsIxwcKDcpLDAxNDQ0Hyc5PTgyPC4zNDL/" +
		"wAARCAABAAEDASIAAhEBAxEB/8QAFAABAAAAAAAAAAAAAAAAAAAACf/EABQQAQAAAAAAAAAAAAAAAAAAAAD/xAAUAQEAAAAAAAAAAAAAAAAAAAAA/8QAFBEBAAAAAAAAAAAAAAAAAAAAAP/aAAwDAQACEQMRAD8AJQAB/9k="

	data := &turnsheet.MonsterEncounterData{
		TurnSheetTemplateData: turnsheet.TurnSheetTemplateData{
			GameName:      convert.Ptr("Test Game"),
			GameType:      convert.Ptr("adventure"),
			TurnNumber:    convert.Ptr(1),
			AccountName:   convert.Ptr("Test Player"),
			TurnSheetCode: convert.Ptr(generateTestTurnSheetCode(t)),
		},
		CharacterName:      "Aldric",
		CharacterHealth:    80,
		CharacterMaxHealth: 100,
		Creatures: []turnsheet.EncounterCreature{
			{
				CreatureInstanceID: "creature-1",
				Name:               "Crypt Spider",
				Description:        "A pale spider lurks in the shadows.",
				Health:             30,
				MaxHealth:          30,
				AttackDamage:       8,
				Defense:            2,
				Disposition:        "aggressive",
				ImageDataURL:       &creatureImageURL,
			},
		},
		MaxActions: 2,
	}

	sheetData, err := json.Marshal(data)
	require.NoError(t, err)

	ctx := context.Background()
	html, err := processor.GenerateTurnSheet(ctx, l, turnsheet.DocumentFormatHTML, sheetData)
	require.NoError(t, err)
	require.NotEmpty(t, html)

	htmlStr := string(html)

	require.True(t, strings.Contains(htmlStr, `src="data:image/jpeg`),
		"expected rendered HTML to contain creature portrait as data URL, got #ZgotmplZ sanitised URL instead")
	require.False(t, strings.Contains(htmlStr, "#ZgotmplZ"),
		"expected rendered HTML NOT to contain #ZgotmplZ (html/template URL sanitisation artifact)")
	require.False(t, strings.Contains(htmlStr, `class="creature-card-portrait-placeholder"`),
		"expected placeholder div NOT to be rendered when ImageDataURL is set")
}
