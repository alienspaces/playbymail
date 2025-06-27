package config

import (
	"gitlab.com/alienspaces/playbymail/core/config"
)

// Config includes core server Config along with additional service
// specific configuration.
type Config struct {
	config.Config
}

// Parse parses environment variables into the provided struct using env.Parse.
func Parse() (Config, error) {
	var cfg Config
	err := config.Parse(&cfg)
	return cfg, err
}
