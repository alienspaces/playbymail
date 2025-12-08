package config

import (
	"fmt"

	"gitlab.com/alienspaces/playbymail/core/config"
)

// Config includes core server Config along with additional service specific configuration.
type Config struct {
	config.Config

	// Emailer
	EmailerFaked bool `env:"EMAILER_FAKED" envDefault:"false"`

	// Email addresses
	SupportEmailAddress string `env:"SUPPORT_EMAIL_ADDRESS" envDefault:"support@playbymail.games"`
	NoReplyEmailAddress string `env:"NO_REPLY_EMAIL_ADDRESS" envDefault:"noreply@playbymail.games"`

	// Save test PDFs (set SAVE_TEST_FILES=true to save, KEEP_TEST_FILES=true to keep after tests)
	SaveTestFiles bool `env:"SAVE_TEST_FILES" envDefault:"false"`

	// Test authentication bypass configuration
	// When both are set, requests with the specified header and matching value
	// can use email as verification code. This enables E2E and API testing.
	// - TestBypassHeaderName: HTTP header name (e.g., "X-Test-Bypass")
	// - TestBypassHeaderValue: Required header value to enable bypass
	TestBypassHeaderName  string `env:"TEST_BYPASS_HEADER_NAME" envDefault:""`
	TestBypassHeaderValue string `env:"TEST_BYPASS_HEADER_VALUE" envDefault:""`

	// Agent configuration
	AgentProvider string `env:"AGENT_PROVIDER" envDefault:"openai"` // "openai", "anthropic", "local"

	// OpenAI settings
	OpenAIAPIKey     string `env:"OPENAI_API_KEY" envDefault:""`
	OpenAIImageModel string `env:"OPENAI_IMAGE_MODEL" envDefault:"gpt-4o-mini"`

	// Anthropic settings (future)
	AnthropicAPIKey string `env:"ANTHROPIC_API_KEY" envDefault:""`
	AnthropicModel  string `env:"ANTHROPIC_MODEL" envDefault:"claude-3-opus"`
}

// Parse parses environment variables into the provided struct using env.Parse.
func Parse() (Config, error) {
	var cfg Config
	err := config.Parse(&cfg)
	if err != nil {
		return cfg, err
	}

	// Validate required configuration that doesn't have defaults applied by env.Parse
	if cfg.AppHost == "" {
		return cfg, fmt.Errorf("APP_HOST is required")
	}

	return cfg, nil
}
