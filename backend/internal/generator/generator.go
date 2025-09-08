package generator

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"gitlab.com/alienspaces/playbymail/core/log"
	"gitlab.com/alienspaces/playbymail/internal/record/account_record"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
)

// PDFGenerator handles PDF generation from HTML templates
type PDFGenerator struct {
	logger      *log.Log
	templateDir string
	outputDir   string
}

// TemplateData represents the data structure for templates
type TemplateData struct {
	// Core context records - always present for every turn
	AccountRec      *account_record.Account
	GameRec         *game_record.Game
	GameInstanceRec *game_record.GameInstance

	// Background images
	BackgroundTop    string
	BackgroundMiddle string
	BackgroundBottom string

	// Content sections
	Header  map[string]any
	Content map[string]any
	Footer  map[string]any

	// Turn sheet specific data - can be cast to specific types per turn sheet type
	TurnSheetData any

	// Unique identification code for scanning
	TurnSheetCode string
}

// NewPDFGenerator creates a new PDF generator
func NewPDFGenerator(l *log.Log, templateDir, outputDir string) *PDFGenerator {
	return &PDFGenerator{
		logger:      l,
		templateDir: templateDir,
		outputDir:   outputDir,
	}
}

// GenerateHTML generates HTML from a template
func (g *PDFGenerator) GenerateHTML(ctx context.Context, templatePath string, data TemplateData) (string, error) {
	g.logger.Debug("generating HTML template=%s", templatePath)

	// Load and parse HTML template
	tmpl, err := g.loadTemplate(templatePath)
	if err != nil {
		g.logger.Error("failed to load template template=%s error=%v", templatePath, err)
		return "", fmt.Errorf("failed to load template: %w", err)
	}

	// Execute template with data
	var html bytes.Buffer
	if err := tmpl.Execute(&html, data); err != nil {
		g.logger.Error("failed to execute template template=%s error=%v", templatePath, err)
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	g.logger.Debug("HTML generated successfully template=%s size=%d", templatePath, html.Len())
	return html.String(), nil
}

// GeneratePDF creates a PDF from an HTML template
func (g *PDFGenerator) GeneratePDF(ctx context.Context, templatePath string, data TemplateData) ([]byte, error) {
	g.logger.Info("starting PDF generation template=%s game=%s turn=%d account=%s", templatePath, data.GameRec.Name, data.GameInstanceRec.CurrentTurn, data.AccountRec.Name)

	// Generate HTML using the existing method
	html, err := g.GenerateHTML(ctx, templatePath, data)
	if err != nil {
		g.logger.Error("failed to generate HTML template=%s error=%v", templatePath, err)
		return nil, fmt.Errorf("failed to generate HTML: %w", err)
	}

	// Convert HTML to PDF using chromedp
	g.logger.Debug("converting HTML to PDF html_size=%d", len(html))
	pdfData, err := g.htmlToPDF(ctx, html)
	if err != nil {
		g.logger.Error("failed to convert HTML to PDF error=%v", err)
		return nil, fmt.Errorf("failed to convert HTML to PDF: %w", err)
	}
	g.logger.Info("PDF generated successfully template=%s pdf_size=%d game=%s turn=%d account=%s", templatePath, len(pdfData), data.GameRec.Name, data.GameInstanceRec.CurrentTurn, data.AccountRec.Name)

	return pdfData, nil
}

// GeneratePDFToFile creates a PDF and saves it to a file
func (g *PDFGenerator) GeneratePDFToFile(ctx context.Context, templatePath string, data TemplateData, filename string) error {
	g.logger.Info("starting PDF generation to file template=%s filename=%s game=%s turn=%d account=%s", templatePath, filename, data.GameRec.Name, data.GameInstanceRec.CurrentTurn, data.AccountRec.Name)

	pdfData, err := g.GeneratePDF(ctx, templatePath, data)
	if err != nil {
		g.logger.Error("failed to generate PDF template=%s filename=%s error=%v", templatePath, filename, err)
		return err
	}

	// Ensure output directory exists
	g.logger.Debug("ensuring output directory exists output_dir=%s", g.outputDir)
	if err := os.MkdirAll(g.outputDir, 0755); err != nil {
		g.logger.Error("failed to create output directory output_dir=%s error=%v", g.outputDir, err)
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Write PDF to file
	filepath := filepath.Join(g.outputDir, filename)
	g.logger.Debug("writing PDF to file filepath=%s pdf_size=%d", filepath, len(pdfData))
	if err := os.WriteFile(filepath, pdfData, 0644); err != nil {
		g.logger.Error("failed to write PDF file filepath=%s error=%v", filepath, err)
		return fmt.Errorf("failed to write PDF file: %w", err)
	}
	g.logger.Info("PDF saved to file successfully filepath=%s pdf_size=%d game=%s turn=%d account=%s", filepath, len(pdfData), data.GameRec.Name, data.GameInstanceRec.CurrentTurn, data.AccountRec.Name)

	return nil
}

// loadTemplate loads and parses an HTML template
func (g *PDFGenerator) loadTemplate(templatePath string) (*template.Template, error) {
	fullPath := filepath.Join(g.templateDir, templatePath)
	templateDir := filepath.Dir(fullPath)

	g.logger.Debug("loading template from path template_path=%s full_path=%s template_dir=%s", templatePath, fullPath, templateDir)

	// Parse template with custom functions
	tmpl := template.New(filepath.Base(templatePath)).Funcs(template.FuncMap{
		"add": func(a, b int) int { return a + b },
		"sub": func(a, b int) int { return a - b },
		"mul": func(a, b int) int { return a * b },
		"div": func(a, b int) int { return a / b },
	})

	// Parse all templates in the directory to support includes
	globPattern := filepath.Join(templateDir, "*.template")
	g.logger.Debug("parsing template glob pattern glob_pattern=%s", globPattern)

	tmpl, err := tmpl.ParseGlob(globPattern)
	if err != nil {
		g.logger.Error("failed to parse template glob glob_pattern=%s error=%v", globPattern, err)
		return nil, err
	}

	g.logger.Debug("template parsed successfully template_path=%s", templatePath)
	return tmpl, nil
}

// htmlToPDF converts HTML to PDF using chromedp
func (g *PDFGenerator) htmlToPDF(ctx context.Context, html string) ([]byte, error) {
	g.logger.Debug("starting HTML to PDF conversion html_size=%d", len(html))

	// Create a new context with Chrome options - minimal for reliability
	opts := []chromedp.ExecAllocatorOption{
		chromedp.Flag("headless", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("disable-extensions", true),
		chromedp.Flag("disable-plugins", true),
		chromedp.Flag("disable-background-timer-throttling", true),
		chromedp.Flag("disable-backgrounding-occluded-windows", true),
		chromedp.Flag("disable-renderer-backgrounding", true),
		chromedp.Flag("disable-features", "TranslateUI"),
		chromedp.Flag("disable-ipc-flooding-protection", true),
	}
	g.logger.Debug("Chrome options configured options_count=%d", len(opts))

	// Check if Chrome is available
	chromePath := os.Getenv("GOOGLE_CHROME_SHIM")
	if chromePath == "" {
		g.logger.Debug("GOOGLE_CHROME_SHIM not set, searching for Chrome in common locations")
		// Try to find Chrome in common locations
		commonPaths := []string{
			"/usr/bin/google-chrome",
			"/usr/bin/chromium-browser",
			"/usr/bin/chromium",
			"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
			"/Applications/Chromium.app/Contents/MacOS/Chromium",
		}

		for _, path := range commonPaths {
			if _, err := os.Stat(path); err == nil {
				chromePath = path
				g.logger.Debug("found Chrome at path chrome_path=%s", chromePath)
				break
			}
		}
	} else {
		g.logger.Debug("using Chrome from GOOGLE_CHROME_SHIM chrome_path=%s", chromePath)
	}

	if chromePath == "" {
		g.logger.Error("Chrome not found in any common locations")
		return nil, fmt.Errorf("Chrome not found. Please install Chrome or set GOOGLE_CHROME_SHIM environment variable")
	}

	// Set Chrome executable path
	opts = append(opts, chromedp.ExecPath(chromePath))
	g.logger.Debug("Chrome executable path set chrome_path=%s", chromePath)

	// Create allocator without timeout (browser needs time to start)
	allocCtx, cancel := chromedp.NewExecAllocator(ctx, opts...)
	defer cancel()

	// Create a new context for the run (no timeout on first Run call)
	runCtx, runCancel := chromedp.NewContext(allocCtx)
	defer runCancel()
	g.logger.Debug("Chrome context created")

	// Write HTML to temporary file for Chrome to load
	tmpFile, err := os.CreateTemp("", "pdf_generation_*.html")
	if err != nil {
		g.logger.Error("failed to create temp file error=%v", err)
		return nil, fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	if _, err := tmpFile.WriteString(html); err != nil {
		g.logger.Error("failed to write HTML to temp file error=%v", err)
		return nil, fmt.Errorf("failed to write HTML to temp file: %w", err)
	}
	tmpFile.Close()

	// Generate PDF using Chrome's print to PDF
	g.logger.Debug("starting Chrome PDF generation")
	var pdfData []byte

	err = chromedp.Run(runCtx,
		chromedp.ActionFunc(func(ctx context.Context) error {
			g.logger.Debug("Chrome browser started, navigating to HTML file")
			return chromedp.Navigate("file://" + tmpFile.Name()).Do(ctx)
		}),
		chromedp.ActionFunc(func(ctx context.Context) error {
			g.logger.Debug("waiting for page to load and render")
			// Simple wait for page to load - no specific element waiting
			time.Sleep(3 * time.Second)
			return nil
		}),
		chromedp.ActionFunc(func(ctx context.Context) error {
			g.logger.Debug("calling Chrome PrintToPDF")
			var err error
			pdfData, _, err = page.PrintToPDF().
				WithPrintBackground(true).
				WithPaperWidth(8.5).
				WithPaperHeight(11.0).
				WithMarginTop(0.5).
				WithMarginBottom(0.5).
				WithMarginLeft(0.5).
				WithMarginRight(0.5).
				Do(ctx)
			if err != nil {
				g.logger.Error("Chrome PrintToPDF failed error=%v", err)
			} else {
				g.logger.Debug("Chrome PrintToPDF completed successfully pdf_size=%d", len(pdfData))
			}
			return err
		}),
	)

	if err != nil {
		g.logger.Error("failed to generate PDF with Chrome error=%v", err)
		return nil, fmt.Errorf("failed to generate PDF: %w", err)
	}

	g.logger.Debug("PDF generation completed successfully pdf_size=%d", len(pdfData))
	return pdfData, nil
}
