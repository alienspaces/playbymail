package jobworker_test

import (
	"bytes"
	"html/template"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/playbymail/internal/utils/testutil"
)

// minimalContentTmpl is a template stub that only defines the `content` block,
// leaving the `footer` block to the base template default so we can test the
// conditional account link logic in the base footer.
const minimalContentTmpl = `{{define "content"}}<p>Test content</p>{{end}}`

// renderBaseWithInlineContent renders the base email template combined with an
// inline content template, using the base template's default footer block.
func renderBaseWithInlineContent(t *testing.T, data any) string {
	t.Helper()

	cfg, _, _, _, _ := testutil.NewDefaultDependencies(t)

	baseTmplPath := filepath.Join(cfg.TemplatesPath, "email", "base.email.html")
	tmpl, err := template.ParseFiles(baseTmplPath)
	require.NoError(t, err, "base template parses without error")

	tmpl, err = tmpl.Parse(minimalContentTmpl)
	require.NoError(t, err, "inline content template parses without error")

	var buf bytes.Buffer
	require.NoError(t, tmpl.ExecuteTemplate(&buf, "base", data))

	return buf.String()
}

func TestBaseEmailTemplate_ConditionalAccountLink(t *testing.T) {
	type tmplData struct {
		SupportEmail string
		AccountURL   string
		Year         int
	}

	t.Run("account URL is omitted when AccountURL is empty", func(t *testing.T) {
		html := renderBaseWithInlineContent(t, tmplData{
			SupportEmail: "support@example.com",
			AccountURL:   "",
			Year:         2026,
		})

		require.NotEmpty(t, html)
		require.NotContains(t, html, "Manage your account",
			"account link should NOT appear when AccountURL is empty")
	})

	t.Run("account URL is rendered when AccountURL is set", func(t *testing.T) {
		html := renderBaseWithInlineContent(t, tmplData{
			SupportEmail: "support@example.com",
			AccountURL:   "http://example.com/account",
			Year:         2026,
		})

		require.NotEmpty(t, html)
		require.True(t, strings.Contains(html, "Manage your account"),
			"account link SHOULD appear when AccountURL is set")
		require.Contains(t, html, "http://example.com/account",
			"account URL SHOULD be present in the rendered output")
	})

	t.Run("support email is always rendered", func(t *testing.T) {
		html := renderBaseWithInlineContent(t, tmplData{
			SupportEmail: "support@example.com",
			AccountURL:   "",
			Year:         2026,
		})

		require.Contains(t, html, "support@example.com")
	})

	t.Run("turn sheet notification template carries AccountURL to base footer", func(t *testing.T) {
		cfg, _, _, _, _ := testutil.NewDefaultDependencies(t)

		baseTmplPath := filepath.Join(cfg.TemplatesPath, "email", "base.email.html")
		specificTmplPath := filepath.Join(cfg.TemplatesPath, "email", "turn_sheet_notification.email.html")

		tmpl, err := template.ParseFiles(baseTmplPath, specificTmplPath)
		require.NoError(t, err)

		data := struct {
			GameName       string
			TurnNumber     int
			TurnSheetURL   string
			ExpirationDate string
			ExpirationTime string
			SupportEmail   string
			AccountURL     string
			Year           int
		}{
			GameName:       "Test Game",
			TurnNumber:     2,
			TurnSheetURL:   "http://example.com/turn-sheets",
			ExpirationDate: "2026-04-01",
			ExpirationTime: "23:59",
			SupportEmail:   "support@example.com",
			AccountURL:     "http://example.com/account",
			Year:           2026,
		}

		var buf bytes.Buffer
		require.NoError(t, tmpl.ExecuteTemplate(&buf, "base", data))

		html := buf.String()
		require.Contains(t, html, "Manage your account",
			"turn sheet notification should include account link")
		require.Contains(t, html, "http://example.com/account")
	})
}
