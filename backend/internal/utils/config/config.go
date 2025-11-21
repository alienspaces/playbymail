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

	// Save test PDFs (set SAVE_TEST_FILES=true to save, KEEP_TEST_FILES=true to keep after tests)
	SaveTestFiles bool `env:"SAVE_TEST_FILES" envDefault:"false"`

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
