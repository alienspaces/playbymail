package generator

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
)

// PDFGenerator handles PDF generation from HTML templates
type PDFGenerator struct {
	logger       logger.Logger
	templatePath string
	outputDir    string
}

// NewPDFGenerator creates a new PDF generator
func NewPDFGenerator(l logger.Logger) (*PDFGenerator, error) {
	l = l.WithFunctionContext("NewPDFGenerator")

	l.Info("creating PDF generator")

	pdfGenerator := &PDFGenerator{
		logger: l,
	}

	return pdfGenerator, nil
}

// SetTemplatePath sets the template path
func (g *PDFGenerator) SetTemplatePath(path string) {
	g.templatePath = path
}

// SetOutputDir sets the output directory
func (g *PDFGenerator) SetOutputDir(dir string) {
	g.outputDir = dir
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

	// Remove leading blank lines that can occur from template whitespace
	htmlStr := html.String()
	htmlStr = strings.TrimLeft(htmlStr, " \t\n\r")

	l.Info("HTML generated successfully template=%s size=%d", templatePath, len(htmlStr))

	// Debug: Save generated HTML to file for inspection
	debugPath := "/tmp/turn_sheet_debug.html"
	if err := os.WriteFile(debugPath, []byte(htmlStr), 0644); err != nil {
		l.Warn("failed to write debug HTML file: %v", err)
	} else {
		l.Info("saved debug HTML to %s", debugPath)
	}

	return htmlStr, nil
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

	// If templatePath is empty, assume we're running from backend/ directory
	templateBase := g.templatePath
	if templateBase == "" {
		templateBase = "."
	}

	fullPath := filepath.Join(templateBase, templatePath)
	templateDir := filepath.Dir(fullPath)

	l.Info("loading template from path template_path=%s full_path=%s template_dir=%s", templatePath, fullPath, templateDir)

	// Parse template with custom functions
	tmpl := template.New(filepath.Base(templatePath)).Funcs(template.FuncMap{
		"add": func(a, b int) int { return a + b },
		"sub": func(a, b int) int { return a - b },
		"mul": func(a, b int) int { return a * b },
		"div": func(a, b int) int { return a / b },
		// safeURL marks a string as safe for use in URL contexts (e.g., data: URLs)
		"safeURL": func(s string) template.URL { return template.URL(s) },
	})

	// Parse base template first (located in turn_sheet/)
	// Make path absolute to avoid issues with working directory
	var baseTemplatePath string
	if templateBase != "" && templateBase != "." {
		baseTemplatePath = filepath.Join(templateBase, "turn_sheet", "base.template")
	} else {
		// We're running from backend/, so use relative path from there
		baseTemplatePath = filepath.Join("templates", "turn_sheet", "base.template")
	}

	l.Info("parsing base template path=%s", baseTemplatePath)

	var err error
	tmpl, err = tmpl.ParseFiles(baseTemplatePath)
	if err != nil {
		l.Warn("failed to parse base template path=%s error=%v", baseTemplatePath, err)
		return nil, fmt.Errorf("failed to parse base template: %w", err)
	}

	// Parse specific template directory to support type-specific includes
	specificTemplatePath := fullPath

	l.Info("parsing specific template path=%s", specificTemplatePath)

	tmpl, err = tmpl.ParseFiles(specificTemplatePath)
	if err != nil {
		l.Warn("failed to parse specific template path=%s error=%v", specificTemplatePath, err)
		return nil, fmt.Errorf("failed to parse specific template: %w", err)
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

	l.Info("chrome options configured options_count=%d", len(opts))

	// Find Chrome executable path
	chromePath := findChromePath(l)
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

	l.Info("chrome executable path set chrome_path=%s", chromePath)

	// Create allocator without timeout (browser needs time to start)
	allocCtx, cancel := chromedp.NewExecAllocator(ctx, opts...)
	defer cancel()

	// Create a new context for the run (no timeout on first Run call)
	runCtx, runCancel := chromedp.NewContext(allocCtx)
	defer runCancel()

	l.Debug("chrome context created")

	// Use data URL to avoid file system access issues in CI environments
	// Base64 encode the HTML and use data URL
	htmlB64 := base64.StdEncoding.EncodeToString([]byte(html))
	dataURL := fmt.Sprintf("data:text/html;base64,%s", htmlB64)
	l.Debug("using data URL for HTML content html_size=%d data_url_size=%d", len(html), len(dataURL))

	// Generate PDF using Chrome's print to PDF
	l.Debug("starting Chrome PDF generation")

	var pdfData []byte

	err := chromedp.Run(runCtx,
		chromedp.ActionFunc(func(ctx context.Context) error {
			l.Debug("chrome browser started, navigating to data URL")
			return chromedp.Navigate(dataURL).Do(ctx)
		}),
		chromedp.ActionFunc(func(ctx context.Context) error {
			l.Debug("waiting for page to load and render")
			// Wait for the page to be ready
			if err := chromedp.WaitReady("body").Do(ctx); err != nil {
				l.Warn("failed to wait for body to be ready >%v<", err)
				// Continue anyway - page might still work
			}

			// Wait for background image to load if it exists
			// First check if the image element exists
			var imageExists bool
			err := chromedp.Evaluate(`document.querySelector('.background-image') !== null`, &imageExists).Do(ctx)
			if err != nil {
				l.Warn("failed to check for background image element >%v<, continuing anyway", err)
			} else if imageExists {
				l.Debug("background image element found, waiting for it to load")
				// Wait for the image to be visible (this ensures it's in the DOM)
				if err := chromedp.WaitVisible(".background-image", chromedp.ByQuery).Do(ctx); err != nil {
					l.Warn("failed to wait for background image to be visible >%v<, continuing anyway", err)
				}

				// Wait for the image to actually load by checking its complete property
				// Poll until the image is loaded or timeout after 10 seconds
				timeout := time.After(10 * time.Second)
				ticker := time.NewTicker(100 * time.Millisecond)

				imageLoaded := false
				for !imageLoaded {
					select {
					case <-timeout:
						l.Warn("timeout waiting for background image to load, continuing anyway")
						ticker.Stop()
						imageLoaded = true // Break out of loop
					case <-ticker.C:
						var imageComplete bool
						err := chromedp.Evaluate(`
							(() => {
								const img = document.querySelector('.background-image');
								return img && img.complete && img.naturalHeight !== 0;
							})()
						`, &imageComplete).Do(ctx)
						if err == nil && imageComplete {
							l.Debug("background image loaded successfully")
							ticker.Stop()
							imageLoaded = true
						}
					}
				}
			} else {
				l.Debug("no background image element found, skipping image load wait")
			}

			// Additional wait to ensure all rendering is complete
			time.Sleep(500 * time.Millisecond)
			return nil
		}),
		chromedp.ActionFunc(func(ctx context.Context) error {
			l.Debug("calling Chrome PrintToPDF")
			var err error
			// A4 paper size: 210mm x 297mm = 8.27in x 11.69in
			pdfData, _, err = page.PrintToPDF().
				WithPrintBackground(true).
				WithPaperWidth(8.27).
				WithPaperHeight(11.69).
				WithMarginTop(0).
				WithMarginBottom(0).
				WithMarginLeft(0).
				WithMarginRight(0).
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
