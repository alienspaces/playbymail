package agent

import (
	"context"

	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/utils/config"
)

// VisionAgent handles image-based AI operations for OCR and structured extraction
type VisionAgent interface {
	// ExtractText extracts plain text from an image
	ExtractText(ctx context.Context, req TextExtractionRequest) (string, error)

	// ExtractStructuredData extracts structured JSON from images with schema
	ExtractStructuredData(ctx context.Context, req StructuredExtractionRequest) ([]byte, error)
}

// TextAgent handles text-based AI operations (for future content generation)
type TextAgent interface {
	// GenerateContent generates text content based on prompts
	GenerateContent(ctx context.Context, req ContentGenerationRequest) (string, error)

	// AnalyzeText analyzes and processes text
	AnalyzeText(ctx context.Context, req TextAnalysisRequest) (TextAnalysisResult, error)
}

// MultiModalAgent combines vision and text capabilities
type MultiModalAgent interface {
	VisionAgent
	TextAgent
}

// NewMultiModalAgent creates a new MultiModalAgent implementation
func NewMultiModalAgent(l logger.Logger, cfg config.Config) MultiModalAgent {
	return NewOpenAIMultimodalAgent(l, cfg)
}
