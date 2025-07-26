package config

import (
	"github.com/caarlos0/env/v10"
)

const (
	AppEnvDevelop    = "develop"
	AppEnvProduction = "production"
)

type Config struct {
	// App environment
	AppEnv string `env:"APP_ENV"`

	// Backend port (default: 8080)
	Port string `env:"PORT" envDefault:"8080"`

	// App home (frontend application) (default: ./frontend/dist)
	AppHome string `env:"APP_HOME" envDefault:"./frontend/dist"`

	// Assets (javascript, css, images, etc.) (default: ./frontend/src/assets)
	AssetsPath string `env:"ASSETS_PATH" envDefault:"./frontend/src/assets"`

	// Schemas (json) (default: ./backend/schemas)
	SchemaPath string `env:"SCHEMA_PATH" envDefault:"./backend/schemas"`

	// Templates (html, email, etc.) (default: ./backend/templates)
	TemplatesPath string `env:"TEMPLATES_PATH" envDefault:"./backend/templates"`

	// CORS
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

	// Forward Email
	ForwardEmailAPIKey string `env:"FORWARDEMAIL_API_KEY"`

	// HMAC key for generating tokens
	TokenHMACKey string `env:"TOKEN_HMAC_KEY"`
}

// Parse parses environment variables into the provided struct using env.Parse.
func Parse[T any](cfg *T) error {
	return env.Parse(cfg)
}
