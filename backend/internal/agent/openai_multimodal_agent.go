package agent

import (
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/utils/config"
)

type multimodalAgent struct {
	VisionAgent
	TextAgent
}

// NewMultiModalAgent creates a new OpenAI MultiModalAgent implementation
func NewOpenAIMultimodalAgent(l logger.Logger, cfg config.Config) MultiModalAgent {
	l = l.WithFunctionContext("NewOpenAIMultimodalAgent")

	l.Info("instantiating openai multi modal agent")

	return &multimodalAgent{
		VisionAgent: NewOpenAIVisionAgent(l, cfg),
		TextAgent:   NewOpenAITextAgent(l, cfg),
	}
}
