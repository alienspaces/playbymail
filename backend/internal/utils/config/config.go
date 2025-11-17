package config

import (
	"gitlab.com/alienspaces/playbymail/core/config"
)

// Config includes core server Config along with additional service specific configuration.
type Config struct {
	config.Config

	// Emailer
	EmailerFaked bool `env:"EMAILER_FAKED" envDefault:"false"`

	// Save test PDFs (set SAVE_TEST_FILES=true to save, KEEP_TEST_FILES=true to keep after tests)
	SaveTestFiles bool `env:"SAVE_TEST_FILES" envDefault:"false"`

	// OpenAI settings
	OpenAIAPIKey     string `env:"OPENAI_API_KEY" envDefault:""`
	OpenAIImageModel string `env:"OPENAI_IMAGE_MODEL" envDefault:"gpt-4o-mini"`
}

// Parse parses environment variables into the provided struct using env.Parse.
func Parse() (Config, error) {
	var cfg Config
	err := config.Parse(&cfg)
	return cfg, err
}
