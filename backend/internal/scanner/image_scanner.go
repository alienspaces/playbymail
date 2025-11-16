package scanner

import (
	"context"
	"fmt"
	"time"

	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/utils/config"
)

// ImageScanner is responsible for generic OCR and text extraction from scanned
// images.
//
// The image scanner provides:
//   - Generic text extraction (OCR output for turn sheet processors to parse)
//
// Turn sheet processors are responsible for:
//   - Parsing extracted text into structured data
//   - Extracting turn sheet codes from OCR text
//   - Validating player choices
//   - Updating game state
type StructuredScanRequest struct {
	Instructions       string
	AdditionalContext  []string
	TemplateImage      []byte
	TemplateImageMIME  string
	FilledImage        []byte
	FilledImageMIME    string
	ExpectedJSONSchema map[string]any
}

type StructuredExtractorFunc func(context.Context, StructuredScanRequest) ([]byte, error)
type TextExtractorFunc func(context.Context, []byte) (string, error)

type ImageScanner struct {
	logger              logger.Logger
	cfg                 config.Config
	structuredExtractor StructuredExtractorFunc
	textExtractor       TextExtractorFunc
}

// NewImageScanner creates a new image scanner that leverages OpenAI-hosted OCR
// flows for both generic text extraction and structured data extraction.
func NewImageScanner(l logger.Logger, cfg config.Config) *ImageScanner {
	s := &ImageScanner{
		logger: l,
		cfg:    cfg,
	}
	s.structuredExtractor = s.extractStructuredViaOpenAI
	s.textExtractor = s.extractTextViaOpenAI
	return s
}

// SetStructuredExtractor allows tests to override the structured extraction
// pathway with a deterministic implementation.
func (s *ImageScanner) SetStructuredExtractor(fn StructuredExtractorFunc) {
	if fn != nil {
		s.structuredExtractor = fn
	}
}

// SetTextExtractor allows tests to override the raw text extraction workflow.
func (s *ImageScanner) SetTextExtractor(fn TextExtractorFunc) {
	if fn != nil {
		s.textExtractor = fn
	}
}

// ExtractTextFromImage performs OCR on image data and returns extracted text
func (s *ImageScanner) ExtractTextFromImage(ctx context.Context, imageData []byte) (string, error) {
	l := s.logger.WithFunctionContext("ImageScanner/ExtractTextFromImage")

	start := time.Now()
	l.Info("extracting text from image data image_size=%d", len(imageData))

	if len(imageData) == 0 {
		return "", fmt.Errorf("empty image data provided")
	}

	if s.textExtractor == nil {
		return "", fmt.Errorf("text extractor not configured")
	}

	text, err := s.textExtractor(ctx, imageData)
	duration := time.Since(start)
	if err != nil {
		l.Warn("text extraction failed after %v error=%v", duration, err)
		return "", err
	}

	l.Info("extracted text length >%d< characters duration=%v", len(text), duration)
	return text, nil
}

// ExtractStructuredData sends the structured extraction request to the hosted
// OCR provider and returns a JSON payload matching the expected schema.
func (s *ImageScanner) ExtractStructuredData(ctx context.Context, req StructuredScanRequest) ([]byte, error) {
	l := s.logger.WithFunctionContext("ImageScanner/ExtractStructuredData")

	start := time.Now()
	l.Info("extracting structured data filled_image_size=%d template_image_size=%d", len(req.FilledImage), len(req.TemplateImage))

	if len(req.FilledImage) == 0 {
		return nil, fmt.Errorf("filled image data is required")
	}

	if req.ExpectedJSONSchema == nil {
		return nil, fmt.Errorf("expected JSON schema is required")
	}

	if s.structuredExtractor == nil {
		return nil, fmt.Errorf("structured extractor not configured")
	}

	resp, err := s.structuredExtractor(ctx, req)
	duration := time.Since(start)
	if err != nil {
		l.Warn("structured extraction failed after %v error=%v", duration, err)
		return nil, err
	}

	l.Info("extracted structured data response_size=%d duration=%v", len(resp), duration)
	return resp, nil
}
