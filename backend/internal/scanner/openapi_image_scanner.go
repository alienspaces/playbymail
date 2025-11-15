package scanner

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	openAIResponsesEndpoint       = "https://api.openai.com/v1/responses"
	openAIImageTranscriptionModel = "gpt-4o-mini"
	openAIImagePrompt             = "Transcribe all text visible in this turn sheet image."
	defaultImageMimeType          = "image/png"
)

type openAIRequest struct {
	Model string        `json:"model"`
	Input []openAIInput `json:"input"`
}

type openAIInput struct {
	Role    string               `json:"role"`
	Content []openAIInputContent `json:"content"`
}

type openAIInputContent struct {
	Type     string `json:"type"`
	Text     string `json:"text,omitempty"`
	ImageURL string `json:"image_url,omitempty"`
}

type openAIResponse struct {
	Output []openAIMessage `json:"output"`
	Error  *openAIError    `json:"error,omitempty"`
}

type openAIMessage struct {
	Content []openAIMessageContent `json:"content"`
}

type openAIMessageContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type openAIError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
}

// extractTextViaOpenAI sends the scanned image to OpenAI's image understanding
// endpoint and returns the transcribed text. The caller must ensure an API key
// is configured before calling.
func (s *ImageScanner) extractTextViaOpenAI(ctx context.Context, imageData []byte) (string, error) {
	if s.cfg.OpenAIAPIKey == "" {
		return "", fmt.Errorf("OpenAI API key not configured")
	}

	if len(imageData) == 0 {
		return "", fmt.Errorf("empty image data provided")
	}

	mimeType := http.DetectContentType(imageData)
	if mimeType == "" || mimeType == "application/octet-stream" {
		mimeType = defaultImageMimeType
	}

	imageB64 := base64.StdEncoding.EncodeToString(imageData)
	imageURI := fmt.Sprintf("data:%s;base64,%s", mimeType, imageB64)

	reqPayload := openAIRequest{
		Model: s.modelName(),
		Input: []openAIInput{
			{
				Role: "user",
				Content: []openAIInputContent{
					{
						Type: "input_text",
						Text: openAIImagePrompt,
					},
					{
						Type:     "input_image",
						ImageURL: imageURI,
					},
				},
			},
		},
	}

	body, err := json.Marshal(reqPayload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal OpenAI request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, openAIResponsesEndpoint, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("failed to create OpenAI request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.cfg.OpenAIAPIKey))

	client := &http.Client{Timeout: 45 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("OpenAI request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read OpenAI response: %w", err)
	}

	if resp.StatusCode >= 300 {
		apiErr := parseOpenAIAPIError(respBody)
		if apiErr != "" {
			return "", fmt.Errorf("OpenAI API error: %s", apiErr)
		}
		return "", fmt.Errorf("OpenAI API request failed with status %d", resp.StatusCode)
	}

	var apiResp openAIResponse
	err = json.Unmarshal(respBody, &apiResp)
	if err != nil {
		return "", fmt.Errorf("failed to decode OpenAI response: %w", err)
	}

	if apiResp.Error != nil && apiResp.Error.Message != "" {
		return "", fmt.Errorf("OpenAI API error: %s", apiResp.Error.Message)
	}

	text := extractOpenAIText(apiResp)
	if text == "" {
		return "", fmt.Errorf("OpenAI did not return any text")
	}

	return text, nil
}

func parseOpenAIAPIError(body []byte) string {
	var errResp openAIResponse
	if err := json.Unmarshal(body, &errResp); err == nil {
		if errResp.Error != nil && errResp.Error.Message != "" {
			return errResp.Error.Message
		}
	}
	return ""
}

func extractOpenAIText(resp openAIResponse) string {
	var builder strings.Builder

	for _, output := range resp.Output {
		for _, content := range output.Content {
			if content.Type == "output_text" && strings.TrimSpace(content.Text) != "" {
				if builder.Len() > 0 {
					builder.WriteString("\n")
				}
				builder.WriteString(strings.TrimSpace(content.Text))
			}
		}
	}

	return strings.TrimSpace(builder.String())
}

func (s *ImageScanner) extractStructuredViaOpenAI(ctx context.Context, req StructuredScanRequest) ([]byte, error) {
	if s.cfg.OpenAIAPIKey == "" {
		return nil, fmt.Errorf("OpenAI API key not configured")
	}

	if len(req.FilledImage) == 0 {
		return nil, fmt.Errorf("filled image data provided is empty")
	}

	if req.ExpectedJSONSchema == nil {
		return nil, fmt.Errorf("expected JSON schema must be provided")
	}

	skeleton, err := json.Marshal(req.ExpectedJSONSchema)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal expected JSON schema: %w", err)
	}

	contents := []openAIInputContent{}

	if strings.TrimSpace(req.Instructions) != "" {
		contents = append(contents, openAIInputContent{
			Type: "input_text",
			Text: req.Instructions,
		})
	}

	for _, ctxLine := range req.AdditionalContext {
		if strings.TrimSpace(ctxLine) == "" {
			continue
		}
		contents = append(contents, openAIInputContent{
			Type: "input_text",
			Text: ctxLine,
		})
	}

	if len(req.TemplateImage) > 0 {
		contents = append(contents, openAIInputContent{
			Type:     "input_image",
			ImageURL: encodeImageDataURI(req.TemplateImage, req.TemplateImageMIME),
		})
	}

	contents = append(contents, openAIInputContent{
		Type:     "input_image",
		ImageURL: encodeImageDataURI(req.FilledImage, req.FilledImageMIME),
	})

	contents = append(contents,
		openAIInputContent{
			Type: "input_text",
			Text: fmt.Sprintf("Return a strict JSON object matching this structure: %s", string(skeleton)),
		},
		openAIInputContent{
			Type: "input_text",
			Text: "Respond with JSON only. Do not include markdown, commentary, or extra keys.",
		},
	)

	reqPayload := openAIRequest{
		Model: s.modelName(),
		Input: []openAIInput{
			{
				Role:    "user",
				Content: contents,
			},
		},
	}

	body, err := json.Marshal(reqPayload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal OpenAI structured request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, openAIResponsesEndpoint, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create OpenAI structured request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.cfg.OpenAIAPIKey))

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("OpenAI structured request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read OpenAI structured response: %w", err)
	}

	if resp.StatusCode >= 300 {
		apiErr := parseOpenAIAPIError(respBody)
		if apiErr != "" {
			return nil, fmt.Errorf("OpenAI API error: %s", apiErr)
		}
		return nil, fmt.Errorf("OpenAI API request failed with status %d", resp.StatusCode)
	}

	var apiResp openAIResponse
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode OpenAI structured response: %w", err)
	}

	if apiResp.Error != nil && apiResp.Error.Message != "" {
		return nil, fmt.Errorf("OpenAI API error: %s", apiResp.Error.Message)
	}

	text := extractOpenAIText(apiResp)
	if text == "" {
		return nil, fmt.Errorf("OpenAI did not return structured data")
	}

	return []byte(text), nil
}

func encodeImageDataURI(data []byte, mime string) string {
	if mime == "" {
		mime = http.DetectContentType(data)
	}
	if mime == "" {
		mime = defaultImageMimeType
	}
	if !strings.HasPrefix(mime, "image/") {
		mime = defaultImageMimeType
	}
	return fmt.Sprintf("data:%s;base64,%s", mime, base64.StdEncoding.EncodeToString(data))
}

func (s *ImageScanner) modelName() string {
	if trimmed := strings.TrimSpace(s.cfg.OpenAIImageModel); trimmed != "" {
		return trimmed
	}
	return openAIImageTranscriptionModel
}
