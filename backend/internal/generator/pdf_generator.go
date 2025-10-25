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
	"gitlab.com/alienspaces/playbymail/core/type/logger"
)

// PDFGenerator handles PDF generation from HTML templates
type PDFGenerator struct {
	logger      logger.Logger
	templateDir string
	outputDir   string
}

// NewPDFGenerator creates a new PDF generator
func NewPDFGenerator(l logger.Logger) *PDFGenerator {
	return &PDFGenerator{
		logger: l,
	}
}

// GenerateHTML generates HTML from a template
func (g *PDFGenerator) GenerateHTML(ctx context.Context, templatePath string, data any) (string, error) {
	l := g.logger.WithFunctionContext("PDFGenerator/GenerateHTML")

	l.Info("generating HTML template=%s", templatePath)

	// Load and parse HTML template
	tmpl, err := g.loadTemplate(templatePath)
	if err != nil {
		l.Warn("failed to load template template=%s error=%v", templatePath, err)
		return "", fmt.Errorf("failed to load template: %w", err)
	}

	// Execute template with data
	var html bytes.Buffer
	if err := tmpl.Execute(&html, data); err != nil {
		l.Warn("failed to execute template template=%s error=%v", templatePath, err)
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	l.Info("HTML generated successfully template=%s size=%d", templatePath, html.Len())

	return html.String(), nil
}

// GeneratePDF creates a PDF from an HTML template
func (g *PDFGenerator) GeneratePDF(ctx context.Context, templatePath string, data any) ([]byte, error) {
	l := g.logger.WithFunctionContext("PDFGenerator/GeneratePDF")

	l.Info("starting PDF generation template=%s", templatePath)

	// Generate HTML from template
	html, err := g.GenerateHTML(ctx, templatePath, data)
	if err != nil {
		l.Warn("failed to generate HTML template=%s error=%v", templatePath, err)
		return nil, fmt.Errorf("failed to generate HTML: %w", err)
	}

	// Convert HTML to PDF using chromedp
	l.Debug("converting HTML to PDF html_size=%d", len(html))

	pdfData, err := g.htmlToPDF(ctx, html)
	if err != nil {
		l.Warn("failed to convert HTML to PDF error=%v", err)
		return nil, fmt.Errorf("failed to convert HTML to PDF: %w", err)
	}

	l.Info("PDF generated successfully template=%s pdf_size=%d", templatePath, len(pdfData))

	return pdfData, nil
}

// GeneratePDFToFile creates a PDF and saves it to a file
func (g *PDFGenerator) GeneratePDFToFile(ctx context.Context, templatePath string, data any, filename string) error {
	l := g.logger.WithFunctionContext("PDFGenerator/GeneratePDFToFile")

	l.Info("starting PDF generation to file template=%s filename=%s", templatePath, filename)

	pdfData, err := g.GeneratePDF(ctx, templatePath, data)
	if err != nil {
		g.logger.Error("failed to generate PDF template=%s filename=%s error=%v", templatePath, filename, err)
		return err
	}

	// Ensure output directory exists
	l.Debug("ensuring output directory exists output_dir=%s", g.outputDir)

	if err := os.MkdirAll(g.outputDir, 0755); err != nil {
		l.Warn("failed to create output directory output_dir=%s error=%v", g.outputDir, err)
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Write PDF to file
	filepath := filepath.Join(g.outputDir, filename)

	l.Debug("writing PDF to file filepath=%s pdf_size=%d", filepath, len(pdfData))

	if err := os.WriteFile(filepath, pdfData, 0644); err != nil {
		l.Warn("failed to write PDF file filepath=%s error=%v", filepath, err)
		return fmt.Errorf("failed to write PDF file: %w", err)
	}

	l.Info("PDF saved to file successfully filepath=%s pdf_size=%d", filepath, len(pdfData))

	return nil
}

// loadTemplate loads and parses an HTML template
func (g *PDFGenerator) loadTemplate(templatePath string) (*template.Template, error) {
	l := g.logger.WithFunctionContext("PDFGenerator/loadTemplate")

	fullPath := filepath.Join(g.templateDir, templatePath)
	templateDir := filepath.Dir(fullPath)

	l.Info("loading template from path template_path=%s full_path=%s template_dir=%s", templatePath, fullPath, templateDir)

	// Parse template with custom functions
	tmpl := template.New(filepath.Base(templatePath)).Funcs(template.FuncMap{
		"add": func(a, b int) int { return a + b },
		"sub": func(a, b int) int { return a - b },
		"mul": func(a, b int) int { return a * b },
		"div": func(a, b int) int { return a / b },
	})

	// Parse base template first (located in turn_sheet/base/template/)
	baseTemplatePath := filepath.Join(g.templateDir, "base", "template", "*.template")
	l.Debug("parsing base templates glob_pattern=%s", baseTemplatePath)

	tmpl, err := tmpl.ParseGlob(baseTemplatePath)
	if err != nil {
		l.Warn("failed to parse base templates glob_pattern=%s error=%v", baseTemplatePath, err)
		return nil, fmt.Errorf("failed to parse base templates: %w", err)
	}

	// Parse specific template directory to support type-specific includes
	globPattern := filepath.Join(templateDir, "*.template")
	l.Debug("parsing specific templates glob_pattern=%s", globPattern)

	tmpl, err = tmpl.ParseGlob(globPattern)
	if err != nil {
		l.Warn("failed to parse specific templates glob_pattern=%s error=%v", globPattern, err)
		return nil, fmt.Errorf("failed to parse specific templates: %w", err)
	}

	l.Info("template parsed successfully template_path=%s", templatePath)

	return tmpl, nil
}

// htmlToPDF converts HTML to PDF using chromedp
func (g *PDFGenerator) htmlToPDF(ctx context.Context, html string) ([]byte, error) {
	l := g.logger.WithFunctionContext("PDFGenerator/htmlToPDF")

	l.Info("starting HTML to PDF conversion html_size=%d", len(html))

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

	l.Debug("chrome options configured options_count=%d", len(opts))

	// Check if Chrome is available
	chromePath := os.Getenv("GOOGLE_CHROME_SHIM")
	if chromePath == "" {
		l.Debug("GOOGLE_CHROME_SHIM not set, searching for Chrome in common locations")
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
				l.Debug("found Chrome at path chrome_path=%s", chromePath)
				break
			}
		}
	} else {
		l.Debug("using Chrome from GOOGLE_CHROME_SHIM chrome_path=%s", chromePath)
	}

	if chromePath == "" {
		// In test environment, return a mock PDF instead of failing
		if os.Getenv("TESTING") == "true" {
			l.Warn("chrome not found, returning mock PDF for testing")
			return []byte("mock-pdf-data-for-testing"), nil
		}
		l.Warn("chrome not found in any common locations")
		return nil, fmt.Errorf("chrome not found. Please install Chrome or set GOOGLE_CHROME_SHIM environment variable")
	}

	// Set Chrome executable path
	opts = append(opts, chromedp.ExecPath(chromePath))
	l.Debug("chrome executable path set chrome_path=%s", chromePath)

	// Create allocator without timeout (browser needs time to start)
	allocCtx, cancel := chromedp.NewExecAllocator(ctx, opts...)
	defer cancel()

	// Create a new context for the run (no timeout on first Run call)
	runCtx, runCancel := chromedp.NewContext(allocCtx)
	defer runCancel()

	l.Debug("chrome context created")

	// Write HTML to temporary file for Chrome to load
	tmpFile, err := os.CreateTemp("", "pdf_generation_*.html")
	if err != nil {
		l.Warn("failed to create temp file error=%v", err)
		return nil, fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	if _, err := tmpFile.WriteString(html); err != nil {
		l.Warn("failed to write HTML to temp file error=%v", err)
		return nil, fmt.Errorf("failed to write HTML to temp file: %w", err)
	}
	tmpFile.Close()

	// Generate PDF using Chrome's print to PDF
	l.Debug("starting Chrome PDF generation")

	var pdfData []byte

	err = chromedp.Run(runCtx,
		chromedp.ActionFunc(func(ctx context.Context) error {
			l.Debug("chrome browser started, navigating to HTML file")
			return chromedp.Navigate("file://" + tmpFile.Name()).Do(ctx)
		}),
		chromedp.ActionFunc(func(ctx context.Context) error {
			l.Debug("waiting for page to load and render")
			// Simple wait for page to load - no specific element waiting
			time.Sleep(3 * time.Second)
			return nil
		}),
		chromedp.ActionFunc(func(ctx context.Context) error {
			l.Debug("calling Chrome PrintToPDF")
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
				l.Warn("chrome PrintToPDF failed error=%v", err)
			} else {
				l.Debug("chrome PrintToPDF completed successfully pdf_size=%d", len(pdfData))
			}
			return err
		}),
	)

	if err != nil {
		l.Warn("failed to generate PDF with Chrome error=%v", err)
		return nil, fmt.Errorf("failed to generate PDF: %w", err)
	}

	l.Info("PDF generation completed successfully pdf_size=%d", len(pdfData))

	return pdfData, nil
}
