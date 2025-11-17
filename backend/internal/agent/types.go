package agent

// TextExtractionRequest contains data for text extraction from images
type TextExtractionRequest struct {
	ImageData []byte
	ImageMIME string
	Prompt    string // Optional custom prompt (defaults to standard OCR prompt)
}

// StructuredExtractionRequest contains data for structured data extraction
type StructuredExtractionRequest struct {
	Instructions      string
	AdditionalContext []string
	TemplateImage     *ImageData // Optional reference image
	FilledImage       ImageData  // Required: image to extract from
	ExpectedSchema    map[string]any
}

// ImageData represents an image with its MIME type
type ImageData struct {
	Data []byte
	MIME string
}

// ContentGenerationRequest contains data for text generation
type ContentGenerationRequest struct {
	SystemPrompt string
	UserPrompt   string
	MaxTokens    int
	Temperature  float64
}

// TextAnalysisRequest contains data for text analysis
type TextAnalysisRequest struct {
	Text    string
	Task    string // "sentiment", "summarize", "extract_entities", etc.
	Options map[string]any
}

// TextAnalysisResult contains analysis results
type TextAnalysisResult struct {
	Task       string
	Result     map[string]any
	Confidence float64
}
