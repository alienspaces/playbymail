package agent

import (
	"context"
	"fmt"

	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/utils/config"
)

type openAITextAgent struct {
	logger logger.Logger
	cfg    config.Config
}

// NewOpenAITextAgent creates a new OpenAI TextAgent implementation
// Note: This is a stub implementation for future text generation features
func NewOpenAITextAgent(l logger.Logger, cfg config.Config) TextAgent {
	return &openAITextAgent{
		logger: l,
		cfg:    cfg,
	}
}

func (a *openAITextAgent) GenerateContent(ctx context.Context, req ContentGenerationRequest) (string, error) {
	return "", fmt.Errorf("text generation not yet implemented")
}

func (a *openAITextAgent) AnalyzeText(ctx context.Context, req TextAnalysisRequest) (TextAnalysisResult, error) {
	return TextAnalysisResult{}, fmt.Errorf("text analysis not yet implemented")
}
