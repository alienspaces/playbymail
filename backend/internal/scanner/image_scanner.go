package scanner

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/agent"
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

type ImageScanner struct {
	logger logger.Logger
	cfg    config.Config
	agent  agent.MultiModalAgent
}

// NewImageScanner creates a new image scanner that uses an AI agent for OCR.
// The scanner handles image processing and delegates AI operations to the agent.
func NewImageScanner(l logger.Logger, cfg config.Config) (*ImageScanner, error) {
	// Direct instantiation of OpenAI agent (simple approach, no factory pattern)
	agent := agent.NewMultiModalAgent(l, cfg)

	return &ImageScanner{
		logger: l,
		cfg:    cfg,
		agent:  agent,
	}, nil
}

// SetAgent allows tests to override the agent with a mock implementation
func (s *ImageScanner) SetAgent(a agent.MultiModalAgent) {
	if a != nil {
		s.agent = a
	}
}

// ExtractTextFromImage performs OCR on image data and returns extracted text.
// Handles image optimization before delegating to agent.
func (s *ImageScanner) ExtractTextFromImage(ctx context.Context, imageData []byte) (string, error) {
	l := s.logger.WithFunctionContext("ImageScanner/ExtractTextFromImage")
	start := time.Now()

	l.Info("extracting text from image data image_size=%d", len(imageData))

	if len(imageData) == 0 {
		return "", fmt.Errorf("empty image data provided")
	}

	if s.agent == nil {
		return "", fmt.Errorf("agent not configured")
	}

	// Image optimization (scanner responsibility)
	mimeType := http.DetectContentType(imageData)
	if mimeType == "" || mimeType == "application/octet-stream" {
		mimeType = "image/png"
	}

	optimized, optimizedMIME, err := optimizeImageForOCR(l, imageData, mimeType)
	if err != nil {
		l.Warn("image optimization failed, using original >%v<", err)
		optimized = imageData
		optimizedMIME = mimeType
	}

	// Delegate to agent (AI provider responsibility)
	req := agent.TextExtractionRequest{
		ImageData: optimized,
		ImageMIME: optimizedMIME,
	}

	text, err := s.agent.ExtractText(ctx, req)
	duration := time.Since(start)
	if err != nil {
		l.Warn("text extraction failed after %v error=%v", duration, err)
		return "", err
	}

	l.Info("extracted text length >%d< characters duration=%v", len(text), duration)
	return text, nil
}

// ExtractStructuredData sends the structured extraction request to the agent.
// Handles image optimization before delegating to agent.
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

	if s.agent == nil {
		return nil, fmt.Errorf("agent not configured")
	}

	// Optimize images (scanner responsibility)
	filledOptimized, filledMIME, err := optimizeImageForOCR(l, req.FilledImage, req.FilledImageMIME)
	if err != nil {
		l.Warn("filled image optimization failed, using original >%v<", err)
		filledOptimized = req.FilledImage
		filledMIME = req.FilledImageMIME
	}

	var templateImage *agent.ImageData
	if len(req.TemplateImage) > 0 {
		templateOptimized, templateMIME, err := optimizeImageForOCR(l, req.TemplateImage, req.TemplateImageMIME)
		if err != nil {
			l.Warn("template image optimization failed, using original >%v<", err)
			templateOptimized = req.TemplateImage
			templateMIME = req.TemplateImageMIME
		}
		templateImage = &agent.ImageData{
			Data: templateOptimized,
			MIME: templateMIME,
		}
	}

	// Delegate to agent (AI provider responsibility)
	agentReq := agent.StructuredExtractionRequest{
		Instructions:      req.Instructions,
		AdditionalContext: req.AdditionalContext,
		TemplateImage:     templateImage,
		FilledImage: agent.ImageData{
			Data: filledOptimized,
			MIME: filledMIME,
		},
		ExpectedSchema: req.ExpectedJSONSchema,
	}

	resp, err := s.agent.ExtractStructuredData(ctx, agentReq)
	duration := time.Since(start)
	if err != nil {
		l.Warn("structured extraction failed after %v error=%v", duration, err)
		return nil, err
	}

	l.Info("extracted structured data response_size=%d duration=%v", len(resp), duration)
	return resp, nil
}
