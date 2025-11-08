package config

import (
	"gitlab.com/alienspaces/playbymail/core/config"
)

// Config includes core server Config along with additional service specific configuration.
type Config struct {
	config.Config

	// Emailer
	EmailerFaked bool `env:"EMAILER_FAKED" envDefault:"false"`

	// Save test PDFs
	SaveTestFiles bool `env:"SAVE_TEST_FILES" envDefault:"true"`
}

// Parse parses environment variables into the provided struct using env.Parse.
func Parse() (Config, error) {
	var cfg Config
	err := config.Parse(&cfg)
	return cfg, err
}
