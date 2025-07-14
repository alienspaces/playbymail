package config

import (
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"gitlab.com/alienspaces/playbymail/core/config"
)

// Config includes core server Config along with additional service
// specific configuration.
type Config struct {
	config.Config

	// Emailer
	EmailerFaked bool `env:"EMAILER_FAKED" envDefault:"false"`
}

// Parse parses environment variables into the provided struct using env.Parse.
func Parse() (Config, error) {
	var cfg Config
	err := config.Parse(&cfg)

	// Debug
	fmt.Println("cfg", spew.Sdump(cfg))

	return cfg, err
}
