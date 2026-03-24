package agent

import (
	"bytes"
	"context"
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

type openAITextAgent struct {
	logger logger.Logger
	cfg    config.Config
	client *http.Client
}

// NewOpenAITextAgent creates a new OpenAI TextAgent implementation that uses the same
// API key, endpoint, and model as the vision agent for consistent configuration.
func NewOpenAITextAgent(l logger.Logger, cfg config.Config) TextAgent {
	l = l.WithFunctionContext("NewOpenAITextAgent")

	l.Info("instantiating openai text agent")

	return &openAITextAgent{
		logger: l,
		cfg:    cfg,
		client: &http.Client{Timeout: 45 * time.Second},
	}
}

// GenerateContent calls the OpenAI Responses API with a text-only prompt, reusing
// the same API key and model configured for turn sheet scanning.
func (a *openAITextAgent) GenerateContent(ctx context.Context, req ContentGenerationRequest) (string, error) {
	l := a.logger.WithFunctionContext("OpenAITextAgent/GenerateContent")
	start := time.Now()

	if a.cfg.OpenAIAPIKey == "" {
		return "", fmt.Errorf("OpenAI API key not configured")
	}

	if req.UserPrompt == "" {
		return "", fmt.Errorf("user prompt is required")
	}

	inputs := []openAIInput{}

	if req.SystemPrompt != "" {
		inputs = append(inputs, openAIInput{
			Role: "system",
			Content: []openAIInputContent{
				{
					Type: "input_text",
					Text: req.SystemPrompt,
				},
			},
		})
	}

	inputs = append(inputs, openAIInput{
		Role: "user",
		Content: []openAIInputContent{
			{
				Type: "input_text",
				Text: req.UserPrompt,
			},
		},
	})

	reqPayload := openAITextRequest{
		Model: a.modelName(),
		Input: inputs,
	}

	if req.Temperature > 0 {
		t := req.Temperature
		reqPayload.Temperature = &t
	}

	if req.MaxTokens > 0 {
		reqPayload.MaxOutputTokens = req.MaxTokens
	}

	body, err := json.Marshal(reqPayload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal OpenAI text request: %w", err)
	}

	apiStart := time.Now()
	model := a.modelName()

	l.Info("sending OpenAI text generation request endpoint=%s model=%s request_size=%d",
		openAIResponsesEndpoint, model, len(body))

	var resp *http.Response
	var respBody []byte
	var apiResp openAIResponse

	backoffConfig := backoff.NewExponentialBackOff()
	backoffConfig.MaxElapsedTime = 2 * time.Minute

	err = backoff.RetryNotify(func() error {
		httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, openAIResponsesEndpoint, bytes.NewReader(body))
		if err != nil {
			return backoff.Permanent(fmt.Errorf("failed to create OpenAI text request: %w", err))
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
			return backoff.Permanent(fmt.Errorf("failed to read OpenAI text response: %w", err))
		}

		l.Debug("read response body size=%d status=%d", len(respBody), resp.StatusCode)

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

		err = json.Unmarshal(respBody, &apiResp)
		if err != nil {
			return backoff.Permanent(fmt.Errorf("failed to decode OpenAI text response: %w", err))
		}

		if apiResp.Error != nil && apiResp.Error.Message != "" {
			if isRetryableError(nil, resp.StatusCode, apiResp.Error.Message) {
				l.Debug("retryable OpenAI API error message=%s", apiResp.Error.Message)
				return fmt.Errorf("OpenAI API error: %s", apiResp.Error.Message)
			}
			return backoff.Permanent(fmt.Errorf("OpenAI API error: %s", apiResp.Error.Message))
		}

		return nil
	}, backoffConfig, func(err error, duration time.Duration) {
		l.Debug("retrying OpenAI text request after %v error=%v", duration, err)
	})

	apiDuration := time.Since(apiStart)
	if err != nil {
		l.Warn("OpenAI text generation request failed after %v error=%v", apiDuration, err)
		return "", err
	}

	l.Info("received OpenAI text response status=%d duration=%v content_length=%d",
		resp.StatusCode, apiDuration, len(respBody))

	text := extractOpenAIText(apiResp)
	if text == "" {
		l.Warn("OpenAI text response contained no text output_count=%d", len(apiResp.Output))
		return "", fmt.Errorf("OpenAI did not return any text")
	}

	totalDuration := time.Since(start)
	l.Info("OpenAI text generation completed text_length=%d total_duration=%v", len(text), totalDuration)

	return text, nil
}

func (a *openAITextAgent) AnalyzeText(ctx context.Context, req TextAnalysisRequest) (TextAnalysisResult, error) {
	return TextAnalysisResult{}, fmt.Errorf("text analysis not yet implemented")
}

func (a *openAITextAgent) modelName() string {
	if trimmed := strings.TrimSpace(a.cfg.OpenAIImageModel); trimmed != "" {
		return trimmed
	}
	return openAIImageTranscriptionModel
}

// openAITextRequest is like openAIRequest but supports temperature and max_output_tokens.
type openAITextRequest struct {
	Model           string        `json:"model"`
	Input           []openAIInput `json:"input"`
	Temperature     *float64      `json:"temperature,omitempty"`
	MaxOutputTokens int           `json:"max_output_tokens,omitempty"`
}
