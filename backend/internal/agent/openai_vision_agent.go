package agent

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

	"github.com/cenkalti/backoff/v4"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/utils/config"
)

const (
	openAIResponsesEndpoint       = "https://api.openai.com/v1/responses"
	openAIImageTranscriptionModel = "gpt-4o-mini"
	openAIImagePrompt             = "Transcribe ALL text visible in this turn sheet image exactly as it appears. Include every character, especially any long alphanumeric codes or strings at the bottom of the page. Do not summarize or describe the code - output the actual characters."
	defaultImageMimeType          = "image/png"
)

type openAIVisionAgent struct {
	logger logger.Logger
	cfg    config.Config
	client *http.Client
}

// NewOpenAIVisionAgent creates a new OpenAI VisionAgent implementation
func NewOpenAIVisionAgent(l logger.Logger, cfg config.Config) VisionAgent {
	l = l.WithFunctionContext("NewOpenAIVisionAgent")

	l.Info("instantiating openai vision agent")

	return &openAIVisionAgent{
		logger: l,
		cfg:    cfg,
		client: &http.Client{Timeout: 45 * time.Second},
	}
}

func (a *openAIVisionAgent) ExtractText(ctx context.Context, req TextExtractionRequest) (string, error) {
	l := a.logger.WithFunctionContext("OpenAIVisionAgent/ExtractText")
	start := time.Now()

	if a.cfg.OpenAIAPIKey == "" {
		return "", fmt.Errorf("OpenAI API key not configured")
	}

	if len(req.ImageData) == 0 {
		return "", fmt.Errorf("empty image data provided")
	}

	prompt := req.Prompt
	if prompt == "" {
		prompt = openAIImagePrompt
	}

	// Encode image as data URI (image should already be optimized by scanner)
	imageB64 := base64.StdEncoding.EncodeToString(req.ImageData)
	imageURI := fmt.Sprintf("data:%s;base64,%s", req.ImageMIME, imageB64)

	l.Info("prepared image data URI image_size=%d mime_type=%s base64_size=%d",
		len(req.ImageData), req.ImageMIME, len(imageB64))

	reqPayload := openAIRequest{
		Model: a.modelName(),
		Input: []openAIInput{
			{
				Role: "user",
				Content: []openAIInputContent{
					{
						Type: "input_text",
						Text: prompt,
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

	apiStart := time.Now()
	model := a.modelName()

	l.Info("sending OpenAI text extraction request endpoint=%s model=%s request_size=%d",
		openAIResponsesEndpoint, model, len(body))

	var resp *http.Response
	var respBody []byte
	var apiResp openAIResponse

	backoffConfig := backoff.NewExponentialBackOff()
	backoffConfig.MaxElapsedTime = 2 * time.Minute

	err = backoff.RetryNotify(func() error {
		// Create a new request for each retry (body reader can only be read once)
		httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, openAIResponsesEndpoint, bytes.NewReader(body))
		if err != nil {
			return backoff.Permanent(fmt.Errorf("failed to create OpenAI request: %w", err))
		}

		httpReq.Header.Set("Content-Type", "application/json")
		httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.cfg.OpenAIAPIKey))

		resp, err = a.client.Do(httpReq)
		if err != nil {
			l.Debug("OpenAI HTTP request failed error=%v", err)
			return err
		}
		defer resp.Body.Close()

		respBody, err = io.ReadAll(resp.Body)
		if err != nil {
			return backoff.Permanent(fmt.Errorf("failed to read OpenAI response: %w", err))
		}

		l.Debug("read response body size=%d status=%d", len(respBody), resp.StatusCode)

		// Check for HTTP errors
		if resp.StatusCode >= 300 {
			apiErr := parseOpenAIAPIError(respBody)
			errMsg := apiErr
			if errMsg == "" {
				errMsg = fmt.Sprintf("HTTP status %d", resp.StatusCode)
			}

			if isRetryableError(nil, resp.StatusCode, errMsg) {
				l.Debug("retryable error status=%d message=%s", resp.StatusCode, errMsg)
				return fmt.Errorf("OpenAI API error: %s", errMsg)
			}
			return backoff.Permanent(fmt.Errorf("OpenAI API error: %s", errMsg))
		}

		// Parse response
		err = json.Unmarshal(respBody, &apiResp)
		if err != nil {
			return backoff.Permanent(fmt.Errorf("failed to decode OpenAI response: %w", err))
		}

		// Check for OpenAI API errors in response body
		if apiResp.Error != nil && apiResp.Error.Message != "" {
			if isRetryableError(nil, resp.StatusCode, apiResp.Error.Message) {
				l.Debug("retryable OpenAI API error message=%s", apiResp.Error.Message)
				return fmt.Errorf("OpenAI API error: %s", apiResp.Error.Message)
			}
			return backoff.Permanent(fmt.Errorf("OpenAI API error: %s", apiResp.Error.Message))
		}

		// Check for content policy refusals in the response text (before returning from retry loop)
		// This allows us to retry if we get a refusal
		text := extractOpenAIText(apiResp)
		if text != "" && IsContentPolicyRefusal(text) {
			l.Debug("content policy refusal detected in response, will retry")
			return fmt.Errorf("OpenAI content policy refusal detected")
		}

		return nil
	}, backoffConfig, func(err error, duration time.Duration) {
		l.Debug("retrying OpenAI request after %v error=%v", duration, err)
	})

	apiDuration := time.Since(apiStart)
	if err != nil {
		l.Warn("OpenAI request failed after %v error=%v", apiDuration, err)
		return "", err
	}

	l.Info("received OpenAI response status=%d duration=%v content_length=%d",
		resp.StatusCode, apiDuration, len(respBody))

	text := extractOpenAIText(apiResp)
	if text == "" {
		l.Warn("OpenAI response contained no text output_count=%d", len(apiResp.Output))
		return "", fmt.Errorf("OpenAI did not return any text")
	}

	// Note: Content policy refusal check is now done inside the retry loop
	// so we don't need to check again here (it would have already triggered a retry)

	totalDuration := time.Since(start)

	l.Info("OpenAI text extraction completed text_length=%d total_duration=%v api_duration=%v prep_duration=%v",
		len(text), totalDuration, apiDuration, apiStart.Sub(start))

	return text, nil
}

func (a *openAIVisionAgent) ExtractStructuredData(ctx context.Context, req StructuredExtractionRequest) ([]byte, error) {
	l := a.logger.WithFunctionContext("OpenAIVisionAgent/ExtractStructuredData")
	start := time.Now()

	if a.cfg.OpenAIAPIKey == "" {
		return nil, fmt.Errorf("OpenAI API key not configured")
	}

	if len(req.FilledImage.Data) == 0 {
		return nil, fmt.Errorf("filled image data is required")
	}

	if req.ExpectedSchema == nil {
		return nil, fmt.Errorf("expected JSON schema is required")
	}

	l.Debug("preparing structured extraction request has_template=%v", req.TemplateImage != nil)

	skeleton, err := json.Marshal(req.ExpectedSchema)
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

	if req.TemplateImage != nil && len(req.TemplateImage.Data) > 0 {
		dataURI := encodeImageDataURI(req.TemplateImage.Data, req.TemplateImage.MIME)
		if dataURI == "" {
			l.Warn("failed to encode template image data URI, skipping template image")
		} else {
			l.Debug("template image encoded template_size=%d template_mime=%s data_uri_length=%d",
				len(req.TemplateImage.Data), req.TemplateImage.MIME, len(dataURI))
			contents = append(contents, openAIInputContent{
				Type:     "input_image",
				ImageURL: dataURI,
			})
		}
	}

	dataURI := encodeImageDataURI(req.FilledImage.Data, req.FilledImage.MIME)
	if dataURI == "" {
		return nil, fmt.Errorf("failed to encode filled image data URI")
	}
	l.Debug("filled image encoded filled_size=%d filled_mime=%s data_uri_length=%d",
		len(req.FilledImage.Data), req.FilledImage.MIME, len(dataURI))

	contents = append(contents,
		openAIInputContent{
			Type:     "input_image",
			ImageURL: dataURI,
		},
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
		Model: a.modelName(),
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

	apiStart := time.Now()
	model := a.modelName()
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

	var resp *http.Response
	var respBody []byte
	var apiResp openAIResponse

	// Use longer timeout for structured extraction
	client := &http.Client{Timeout: 60 * time.Second}
	backoffConfig := backoff.NewExponentialBackOff()
	backoffConfig.MaxElapsedTime = 3 * time.Minute

	err = backoff.RetryNotify(func() error {
		// Create a new request for each retry (body reader can only be read once)
		httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, openAIResponsesEndpoint, bytes.NewReader(body))
		if err != nil {
			return backoff.Permanent(fmt.Errorf("failed to create OpenAI structured request: %w", err))
		}

		httpReq.Header.Set("Content-Type", "application/json")
		httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.cfg.OpenAIAPIKey))

		resp, err = client.Do(httpReq)
		if err != nil {
			l.Debug("OpenAI HTTP request failed error=%v", err)
			return err
		}
		defer resp.Body.Close()

		respBody, err = io.ReadAll(resp.Body)
		if err != nil {
			return backoff.Permanent(fmt.Errorf("failed to read OpenAI structured response: %w", err))
		}

		l.Debug("read response body size=%d status=%d", len(respBody), resp.StatusCode)

		// Check for HTTP errors
		if resp.StatusCode >= 300 {
			apiErr := parseOpenAIAPIError(respBody)
			errMsg := apiErr
			if errMsg == "" {
				errMsg = fmt.Sprintf("HTTP status %d", resp.StatusCode)
			}

			if isRetryableError(nil, resp.StatusCode, errMsg) {
				l.Debug("retryable error status=%d message=%s", resp.StatusCode, errMsg)
				return fmt.Errorf("OpenAI API error: %s", errMsg)
			}
			return backoff.Permanent(fmt.Errorf("OpenAI API error: %s", errMsg))
		}

		// Parse response
		err = json.Unmarshal(respBody, &apiResp)
		if err != nil {
			return backoff.Permanent(fmt.Errorf("failed to decode OpenAI structured response: %w", err))
		}

		// Check for OpenAI API errors in response body
		if apiResp.Error != nil && apiResp.Error.Message != "" {
			if isRetryableError(nil, resp.StatusCode, apiResp.Error.Message) {
				l.Debug("retryable OpenAI API error message=%s", apiResp.Error.Message)
				return fmt.Errorf("OpenAI API error: %s", apiResp.Error.Message)
			}
			return backoff.Permanent(fmt.Errorf("OpenAI API error: %s", apiResp.Error.Message))
		}

		// Check for content policy refusals in structured extraction response
		text := extractOpenAIText(apiResp)
		if text != "" && IsContentPolicyRefusal(text) {
			l.Debug("content policy refusal detected in structured extraction response, will retry")
			return fmt.Errorf("OpenAI content policy refusal detected")
		}

		return nil
	}, backoffConfig, func(err error, duration time.Duration) {
		l.Debug("retrying OpenAI structured extraction request after %v error=%v", duration, err)
	})

	apiDuration := time.Since(apiStart)
	if err != nil {
		l.Warn("OpenAI structured extraction request failed after %v error=%v", apiDuration, err)
		return nil, err
	}

	l.Info("received OpenAI response status=%d duration=%v content_length=%d",
		resp.StatusCode, apiDuration, len(respBody))

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

// Helper types and functions for OpenAI API
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

func parseOpenAIAPIError(body []byte) string {
	var errResp openAIResponse
	if err := json.Unmarshal(body, &errResp); err == nil {
		if errResp.Error != nil && errResp.Error.Message != "" {
			return errResp.Error.Message
		}
	}
	return ""
}

// isRetryableError determines if an OpenAI error should be retried
func isRetryableError(err error, statusCode int, errorMessage string) bool {
	// Retry on network errors
	if err != nil {
		return true
	}

	// Retry on 5xx server errors
	if statusCode >= 500 && statusCode < 600 {
		return true
	}

	// Retry on rate limit errors
	if statusCode == http.StatusTooManyRequests {
		return true
	}

	// Retry on service unavailable
	if statusCode == http.StatusServiceUnavailable {
		return true
	}

	// Retry on OpenAI errors that suggest retrying
	if errorMessage != "" {
		lowerMsg := strings.ToLower(errorMessage)
		if strings.Contains(lowerMsg, "retry") ||
			strings.Contains(lowerMsg, "error occurred while processing") ||
			strings.Contains(lowerMsg, "rate limit") ||
			strings.Contains(lowerMsg, "server error") ||
			strings.Contains(lowerMsg, "temporary") ||
			strings.Contains(lowerMsg, "content policy") ||
			strings.Contains(lowerMsg, "refusal") {
			return true
		}
	}

	return false
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

func (a *openAIVisionAgent) modelName() string {
	if trimmed := strings.TrimSpace(a.cfg.OpenAIImageModel); trimmed != "" {
		return trimmed
	}
	return openAIImageTranscriptionModel
}
