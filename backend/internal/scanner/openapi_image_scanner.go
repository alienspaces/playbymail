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

// OpenAPI
// https://platform.openai.com/docs/api-reference/responses/create

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

// extractTextViaOpenAI sends the scanned image to OpenAI's responses endpoint
// (https://api.openai.com/v1/responses) and returns the transcribed text.
// The caller must ensure an API key is configured before calling.
func (s *ImageScanner) extractTextViaOpenAI(ctx context.Context, imageData []byte) (string, error) {
	l := s.logger.WithFunctionContext("ImageScanner/extractTextViaOpenAI")

	start := time.Now()

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

	// Optimize image before encoding
	optimized, optimizedMIME, err := optimizeImageForOCR(l, imageData, mimeType)
	if err != nil {
		l.Warn("image optimization failed, using original >%v<", err)
		optimized = imageData
		optimizedMIME = mimeType
	}

	imageB64 := base64.StdEncoding.EncodeToString(optimized)
	imageURI := fmt.Sprintf("data:%s;base64,%s", optimizedMIME, imageB64)

	l.Info("prepared image data URI original_size=%d optimized_size=%d mime_type=%s base64_size=%d",
		len(imageData), len(optimized), optimizedMIME, len(imageB64))

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
	apiStart := time.Now()
	model := s.modelName()

	l.Info("sending OpenAI text extraction request endpoint=%s model=%s request_size=%d",
		openAIResponsesEndpoint, model, len(body))

	resp, err := client.Do(httpReq)
	apiDuration := time.Since(apiStart)
	if err != nil {
		l.Warn("OpenAI HTTP request failed after %v error=%v", apiDuration, err)
		return "", fmt.Errorf("OpenAI request failed: %w", err)
	}
	defer resp.Body.Close()

	l.Info("received OpenAI response status=%d duration=%v content_length=%s",
		resp.StatusCode, apiDuration, resp.Header.Get("Content-Length"))

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read OpenAI response: %w", err)
	}

	l.Debug("read response body size=%d", len(respBody))

	if resp.StatusCode >= 300 {
		apiErr := parseOpenAIAPIError(respBody)
		errorPreview := string(respBody)
		if len(errorPreview) > 200 {
			errorPreview = errorPreview[:200] + "..."
		}
		l.Warn("OpenAI API error response status=%d body_preview=%s", resp.StatusCode, errorPreview)
		if apiErr != "" {
			return "", fmt.Errorf("OpenAI API error: %s", apiErr)
		}
		return "", fmt.Errorf("OpenAI API request failed with status %d", resp.StatusCode)
	}

	var apiResp openAIResponse
	err = json.Unmarshal(respBody, &apiResp)
	if err != nil {
		l.Warn("failed to decode OpenAI response body_size=%d error=%v", len(respBody), err)
		return "", fmt.Errorf("failed to decode OpenAI response: %w", err)
	}

	if apiResp.Error != nil && apiResp.Error.Message != "" {
		l.Warn("OpenAI API returned error type=%s message=%s", apiResp.Error.Type, apiResp.Error.Message)
		return "", fmt.Errorf("OpenAI API error: %s", apiResp.Error.Message)
	}

	text := extractOpenAIText(apiResp)
	if text == "" {
		l.Warn("OpenAI response contained no text output_count=%d", len(apiResp.Output))
		return "", fmt.Errorf("OpenAI did not return any text")
	}

	totalDuration := time.Since(start)

	l.Info("OpenAI text extraction completed text_length=%d total_duration=%v api_duration=%v prep_duration=%v",
		len(text), totalDuration, apiDuration, apiStart.Sub(start))

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
	l := s.logger.WithFunctionContext("ImageScanner/extractStructuredViaOpenAI")
	start := time.Now()

	if s.cfg.OpenAIAPIKey == "" {
		return nil, fmt.Errorf("OpenAI API key not configured")
	}

	if len(req.FilledImage) == 0 {
		return nil, fmt.Errorf("filled image data provided is empty")
	}

	if req.ExpectedJSONSchema == nil {
		return nil, fmt.Errorf("expected JSON schema must be provided")
	}

	l.Debug("preparing structured extraction request has_template=%v", len(req.TemplateImage) > 0)

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
		optimizedTemplate, optimizedTemplateMIME, err := optimizeImageForOCR(l, req.TemplateImage, req.TemplateImageMIME)
		if err != nil {
			l.Warn("template image optimization failed, using original >%v<", err)
			optimizedTemplate = req.TemplateImage
			optimizedTemplateMIME = req.TemplateImageMIME
		}
		dataURI := encodeImageDataURI(optimizedTemplate, optimizedTemplateMIME)
		if dataURI == "" {
			l.Warn("failed to encode template image data URI, skipping template image")
		} else {
			l.Debug("template image encoded template_size=%d template_mime=%s data_uri_length=%d",
				len(optimizedTemplate), optimizedTemplateMIME, len(dataURI))
			contents = append(contents, openAIInputContent{
				Type:     "input_image",
				ImageURL: dataURI,
			})
		}
	}

	optimizedFilled, optimizedFilledMIME, err := optimizeImageForOCR(l, req.FilledImage, req.FilledImageMIME)
	if err != nil {
		l.Warn("filled image optimization failed, using original >%v<", err)
		optimizedFilled = req.FilledImage
		optimizedFilledMIME = req.FilledImageMIME
	}
	dataURI := encodeImageDataURI(optimizedFilled, optimizedFilledMIME)
	if dataURI == "" {
		return nil, fmt.Errorf("failed to encode filled image data URI")
	}
	l.Debug("filled image encoded filled_size=%d filled_mime=%s data_uri_length=%d",
		len(optimizedFilled), optimizedFilledMIME, len(dataURI))

	contents = append(contents, openAIInputContent{
		Type:     "input_image",
		ImageURL: dataURI,
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
	apiStart := time.Now()
	model := s.modelName()
	imageCount := 0
	textCount := 0
	for _, c := range contents {
		switch c.Type {
		case "input_image":
			imageCount++
		case "input_text":
			textCount++
		}
	}
	l.Info("sending OpenAI structured extraction request endpoint=%s model=%s request_size=%d content_items=%d images=%d text=%d",
		openAIResponsesEndpoint, model, len(body), len(contents), imageCount, textCount)
	resp, err := client.Do(httpReq)
	apiDuration := time.Since(apiStart)
	if err != nil {
		l.Warn("OpenAI HTTP request failed after %v error=%v", apiDuration, err)
		return nil, fmt.Errorf("OpenAI structured request failed: %w", err)
	}
	defer resp.Body.Close()

	l.Info("received OpenAI response status=%d duration=%v content_length=%s",
		resp.StatusCode, apiDuration, resp.Header.Get("Content-Length"))

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read OpenAI structured response: %w", err)
	}
	l.Debug("read response body size=%d", len(respBody))

	if resp.StatusCode >= 300 {
		apiErr := parseOpenAIAPIError(respBody)
		errorPreview := string(respBody)
		if len(errorPreview) > 200 {
			errorPreview = errorPreview[:200] + "..."
		}
		l.Warn("OpenAI API error response status=%d body_preview=%s", resp.StatusCode, errorPreview)
		if apiErr != "" {
			return nil, fmt.Errorf("OpenAI API error: %s", apiErr)
		}
		return nil, fmt.Errorf("OpenAI API request failed with status %d", resp.StatusCode)
	}

	var apiResp openAIResponse
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		l.Warn("failed to decode OpenAI structured response body_size=%d error=%v", len(respBody), err)
		return nil, fmt.Errorf("failed to decode OpenAI structured response: %w", err)
	}

	if apiResp.Error != nil && apiResp.Error.Message != "" {
		l.Warn("OpenAI API returned error type=%s message=%s", apiResp.Error.Type, apiResp.Error.Message)
		return nil, fmt.Errorf("OpenAI API error: %s", apiResp.Error.Message)
	}

	text := extractOpenAIText(apiResp)
	if text == "" {
		l.Warn("OpenAI response contained no text output_count=%d", len(apiResp.Output))
		return nil, fmt.Errorf("OpenAI did not return structured data")
	}

	totalDuration := time.Since(start)
	prepDuration := apiStart.Sub(start)
	l.Info("OpenAI structured extraction completed response_size=%d total_duration=%v api_duration=%v prep_duration=%v",
		len(text), totalDuration, apiDuration, prepDuration)
	return []byte(text), nil
}

func encodeImageDataURI(data []byte, mime string) string {
	// Validate image data before encoding
	if len(data) == 0 {
		return ""
	}

	// If MIME type not provided, detect it
	if mime == "" {
		mime = http.DetectContentType(data)
	}

	// Validate MIME type matches actual image format
	if mime == "" || mime == "application/octet-stream" {
		// Try to detect from image signature
		if len(data) >= 3 && data[0] == 0xFF && data[1] == 0xD8 && data[2] == 0xFF {
			mime = "image/jpeg"
		} else if len(data) >= 8 && bytes.Equal(data[0:8], []byte{0x89, 'P', 'N', 'G', 0x0D, 0x0A, 0x1A, 0x0A}) {
			mime = "image/png"
		} else {
			mime = defaultImageMimeType
		}
	}

	// Ensure MIME type is valid
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
