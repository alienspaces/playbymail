package config

import (
	"github.com/caarlos0/env/v10"
)

const (
	AppEnvDevelop    = "develop"
	AppEnvProduction = "production"
)

type Config struct {
	AppEnv             string   `env:"APP_ENV"`
	Port               string   `env:"PORT" envDefault:"8080"`
	AppHome            string   `env:"APP_HOME" envDefault:"./frontend/dist"`
	AssetsPath         string   `env:"ASSETS_PATH" envDefault:"./frontend/src/assets"`
	CORSAllowedOrigins []string `env:"CORS_ALLOWED_ORIGINS" envDefault:"*"`
	CORSAllowedHeaders []string `env:"CORS_ALLOWED_HEADERS" envDefault:"Content-Type,Authorization"`

	// Database
	DatabaseURL                string `env:"DATABASE_URL,required"`
	DatabaseMaxOpenConnections int    `env:"DATABASE_MAX_OPEN_CONNECTIONS" envDefault:"180"`
	DatabaseMaxIdleConnections int    `env:"DATABASE_MAX_IDLE_CONNECTIONS" envDefault:"45"`
	DatabaseMaxIdleTimeMins    int    `env:"DATABASE_MAX_IDLE_TIME_MINS" envDefault:"15"`

	// Log
	LogLevel    string `env:"LOG_LEVEL" envDefault:"debug"`
	LogIsPretty bool   `env:"LOG_IS_PRETTY" envDefault:"true"`

	// SMTP
	SMTPHost string `env:"SMTP_HOST"`

	// Sendgrid
	SendgridAPIKey string `env:"SENGRID_API_KEY"`
}

// Parse parses environment variables into the provided struct using env.Parse.
func Parse[T any](cfg *T) error {
	return env.Parse(cfg)
}
