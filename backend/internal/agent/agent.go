package agent

import (
	"context"
	"strings"

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
	l = l.WithFunctionContext("NewMultiModalAgent")

	l.Info("instantiating multi modal agent")

	return NewOpenAIMultimodalAgent(l, cfg)
}

// IsContentPolicyRefusal detects if the extracted text is a content policy refusal message.
// This is a generic function that works across different AI agent providers, as they may
// return similar refusal patterns when content policies are triggered.
func IsContentPolicyRefusal(text string) bool {
	if text == "" {
		return false
	}

	lowerText := strings.ToLower(strings.TrimSpace(text))

	// Common refusal patterns from AI agents
	refusalPatterns := []string{
		"i'm sorry",
		"i cannot",
		"i can't",
		"i am not able",
		"i'm not able",
		"i cannot assist",
		"i can't assist",
		"unable to assist",
		"cannot help",
		"can't help",
		"not able to help",
		"content policy",
		"against my usage policies",
		"against my guidelines",
		"against my programming",
		"inappropriate",
		"harmful",
		"unsafe",
	}

	// Check if text matches common refusal patterns
	for _, pattern := range refusalPatterns {
		if strings.Contains(lowerText, pattern) {
			// Additional check: refusal messages are typically short and don't contain actual extracted text
			// Real OCR text would be longer and contain alphanumeric characters
			if len(text) < 200 && !ContainsAlphanumericCode(text) {
				return true
			}
		}
	}

	return false
}

// ContainsAlphanumericCode checks if text contains alphanumeric codes (typical of turn sheet codes).
// This helps distinguish between refusal messages and actual OCR output.
func ContainsAlphanumericCode(text string) bool {
	// Look for patterns like: letters and numbers mixed, long alphanumeric strings
	hasLetters := false
	hasNumbers := false
	consecutiveAlnum := 0
	maxConsecutiveAlnum := 0

	for _, r := range text {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
			hasLetters = true
			consecutiveAlnum++
		} else if r >= '0' && r <= '9' {
			hasNumbers = true
			consecutiveAlnum++
		} else {
			if consecutiveAlnum > maxConsecutiveAlnum {
				maxConsecutiveAlnum = consecutiveAlnum
			}
			consecutiveAlnum = 0
		}
	}
	if consecutiveAlnum > maxConsecutiveAlnum {
		maxConsecutiveAlnum = consecutiveAlnum
	}

	// Turn sheet codes are typically 10+ alphanumeric characters
	return hasLetters && hasNumbers && maxConsecutiveAlnum >= 10
}
