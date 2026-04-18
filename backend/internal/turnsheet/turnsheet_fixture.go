package turnsheet

import (
	"context"

	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/utils/config"
)

// TurnSheetProcessor is implemented by all turn sheet processor types.
// It is used by DevFixture.NewProcessor to provide a unified rendering interface
// for the automated rendering tests.
type TurnSheetProcessor interface {
	GenerateTurnSheet(ctx context.Context, l logger.Logger, format DocumentFormat, data []byte) ([]byte, error)
}

// DevFixture describes a single turn sheet type's sample rendering fixture.
//
// It is the single source of truth for sample/demo turn sheet data and is used by:
//   - The rendering tests (TestRenderAllSheets) to generate HTML and PDF output files.
//   - The dev server (cmd/turnsheet-dev) to render HTML for the live gallery.
//
// Each turn sheet type registers its fixture via a dedicated *_fixture.go file alongside
// its implementation (*_processor.go) and tests (*_test.go).
type DevFixture struct {
	// TemplatePath is relative to the templates root, e.g. "turnsheet/adventure_game_location_choice.template".
	TemplatePath string
	// OutputBaseName is the base filename without extension, e.g. "adventure_game_location_choice_turnsheet".
	OutputBaseName string
	// BackgroundFile is the filename within the testdata directory, e.g. "background-cliffpath.png".
	BackgroundFile string
	// IsJoinSheet indicates the fixture uses a join-game turn sheet code rather than a play code.
	// Used by the rendering test to call the correct code generator.
	IsJoinSheet bool
	// MakeData constructs the template data for this fixture using the provided background image
	// as a base64 data URI and the turn sheet code string.
	MakeData func(backgroundDataURI string, turnSheetCode string) any
	// NewProcessor constructs the processor for this sheet type.
	// Used by the rendering tests; the dev server calls GenerateHTML directly.
	NewProcessor func(l logger.Logger, cfg config.Config) (TurnSheetProcessor, error)
}

// strPtr / intPtr are small helpers used by the individual *_fixture.go files.
func strPtr(v string) *string { return &v }
func intPtr(v int) *int       { return &v }

// DevFixtures returns all turn sheet sample fixtures in gallery display order.
// Each fixture is defined in its own *_fixture.go file alongside the sheet implementation.
func DevFixtures() []DevFixture {
	return []DevFixture{
		AdventureGameLocationChoiceFixture(),
		AdventureGameInventoryManagementFixture(),
		AdventureGameMonsterEncounterFixture(),
		AdventureGameJoinGameFixture(),
		MechaGameSquadManagementFixture(),
		MechaGameOrdersFixture(),
		MechaGameJoinGameFixture(),
	}
}
