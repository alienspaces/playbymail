package turnsheet_test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/playbymail/internal/turnsheet"
	"gitlab.com/alienspaces/playbymail/internal/utils/testutil"
)

// TestRenderAllSheets generates HTML and PDF output files for every turn sheet type using
// the shared dev fixtures (turnsheet.DevFixtures). This is the single authoritative rendering
// test — one source of truth for sample data used by both automated tests and the dev server.
//
// Run via:
//
//	./tools/render-turnsheets
func TestRenderAllSheets(t *testing.T) {
	cfg, l, _, _, _ := testutil.NewDefaultDependencies(t)

	cfg.TemplatesPath = "../../templates"
	cfg.SaveTestFiles = true

	ctx := context.Background()

	type formatCase struct {
		format turnsheet.DocumentFormat
		ext    string
	}
	formats := []formatCase{
		{turnsheet.DocumentFormatHTML, "html"},
		{turnsheet.DocumentFormatPDF, "pdf"},
	}

	for _, f := range turnsheet.DevFixtures() {
		f := f
		t.Run(f.OutputBaseName, func(t *testing.T) {
			var code string
			if f.IsJoinSheet {
				code = generateTestJoinTurnSheetCode(t)
			} else {
				code = generateTestTurnSheetCode(t)
			}

			bg := loadTestBackgroundImage(t, "testdata/"+f.BackgroundFile)
			data, err := json.Marshal(f.MakeData(bg, code))
			require.NoError(t, err, "should marshal fixture data")

			processor, err := f.NewProcessor(l, cfg)
			require.NoError(t, err, "should create processor")

			for _, fc := range formats {
				output, err := processor.GenerateTurnSheet(ctx, l, fc.format, data)
				require.NoError(t, err, "should render %s as %s", f.OutputBaseName, fc.ext)
				require.NotEmpty(t, output, "output should not be empty")

				if cfg.SaveTestFiles {
					path := fmt.Sprintf("testdata/%s.%s", f.OutputBaseName, fc.ext)
					err = os.WriteFile(path, output, 0644)
					require.NoError(t, err, "should write %s", path)
					t.Logf("wrote %s", path)
				}
			}
		})
	}
}
