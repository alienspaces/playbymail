package generator

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

// GeneratePNG renders the specified template data as a PNG screenshot. This is
// primarily used to provide blank reference images for OCR pipelines.
func (g *PDFGenerator) GeneratePNG(ctx context.Context, templatePath string, data any) ([]byte, error) {
	l := g.logger.WithFunctionContext("PDFGenerator/GeneratePNG")

	html, err := g.GenerateHTML(ctx, templatePath, data)
	if err != nil {
		return nil, err
	}

	pngData, err := g.htmlToPNG(ctx, html)
	if err != nil {
		return nil, err
	}

	l.Info("PNG generated successfully template=%s png_size=%d", templatePath, len(pngData))
	return pngData, nil
}

var pngSignature = []byte{0x89, 'P', 'N', 'G', 0x0D, 0x0A, 0x1A, 0x0A}

func (g *PDFGenerator) htmlToPNG(ctx context.Context, html string) ([]byte, error) {
	l := g.logger.WithFunctionContext("PDFGenerator/htmlToPNG")

	l.Info("starting HTML to PNG conversion html_size=%d", len(html))

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

	chromePath := os.Getenv("GOOGLE_CHROME_SHIM")
	if chromePath == "" {
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
				break
			}
		}
	}

	if chromePath == "" {
		return nil, fmt.Errorf("chrome not found. Please install Chrome or set GOOGLE_CHROME_SHIM environment variable")
	}

	opts = append(opts, chromedp.ExecPath(chromePath))

	allocCtx, cancelAlloc := chromedp.NewExecAllocator(ctx, opts...)
	defer cancelAlloc()

	runCtx, cancelRun := chromedp.NewContext(allocCtx)
	defer cancelRun()

	tmpFile, err := os.CreateTemp("", "turn_sheet_preview_*.html")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(html); err != nil {
		return nil, fmt.Errorf("failed to write HTML to temp file: %w", err)
	}
	tmpFile.Close()

	var pngData []byte

	err = chromedp.Run(runCtx,
		chromedp.Navigate("file://"+tmpFile.Name()),
		chromedp.ActionFunc(func(ctx context.Context) error {
			return emulation.SetDeviceMetricsOverride(1240, 1754, 1.0, false).
				WithScreenOrientation(&emulation.ScreenOrientation{
					Type:  emulation.OrientationTypePortraitPrimary,
					Angle: 0,
				}).Do(ctx)
		}),
		chromedp.ActionFunc(func(ctx context.Context) error {
			time.Sleep(2 * time.Second)
			return nil
		}),
		chromedp.ActionFunc(func(ctx context.Context) error {
			buf, err := page.CaptureScreenshot().
				WithFormat(page.CaptureScreenshotFormatPng).
				WithQuality(100).
				WithFromSurface(true).
				Do(ctx)
			if err != nil {
				return err
			}
			pngData = buf
			return nil
		}),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to capture PNG screenshot: %w", err)
	}

	if len(pngData) < len(pngSignature) || !bytes.Equal(pngData[:len(pngSignature)], pngSignature) {
		return nil, fmt.Errorf("chrome screenshot did not return PNG bytes (signature=%v)", pngData[:min(len(pngData), len(pngSignature))])
	}

	return pngData, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
