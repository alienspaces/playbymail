package log

import (
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog"

	"gitlab.com/alienspaces/playbymail/core/config"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
)

const (
	ContextKeyCorrelationID = "correlation-id"
)

// Log -
type Log struct {
	log    zerolog.Logger
	fields map[string]interface{}
	config config.Config
}

var _ logger.Logger = &Log{}

// Level -
type Level uint32

const (
	// DebugLevel -
	DebugLevel = 5
	// InfoLevel -
	InfoLevel = 4
	// WarnLevel -
	WarnLevel = 3
	// ErrorLevel -
	ErrorLevel = 2
)

var levelMap = map[Level]zerolog.Level{
	// DebugLevel -
	DebugLevel: zerolog.DebugLevel,
	// InfoLevel -
	InfoLevel: zerolog.InfoLevel,
	// WarnLevel -
	WarnLevel: zerolog.WarnLevel,
	// ErrorLevel -
	ErrorLevel: zerolog.ErrorLevel,
}

func NewDefaultLogger() *Log {
	l := Log{
		fields: make(map[string]interface{}),
		config: config.Config{
			LogLevel:    "info",
			LogIsPretty: false,
		},
	}

	l.Init()
	return &l
}

// NewLogger returns a logger
func NewLogger(cfg config.Config) (*Log, error) {

	l := Log{
		fields: make(map[string]interface{}),
		config: cfg,
	}

	l.Init()
	return &l, nil
}

// NewLoggerWithConfig returns a logger with the provided configuration
func NewLoggerWithConfig(cfg config.Config) (*Log, error) {

	l := Log{
		fields: make(map[string]interface{}),
		config: cfg,
	}

	l.Init()
	return &l, nil
}

// Init initializes logger
func (l *Log) Init() {

	l.log = zerolog.New(os.Stdout).With().Timestamp().Logger()

	// Pretty
	if l.config.LogIsPretty {
		output := zerolog.ConsoleWriter{
			Out: os.Stdout,
			// The following adds colour to the value of additional log fields,
			// a nice shade of purple actually..
			FormatFieldValue: func(i interface{}) string {
				if i != nil {
					return fmt.Sprintf("\x1b[%dm%v\x1b[0m", 35, i)
				}
				return ""
			},
		}
		l.log = l.log.Output(output)
	}

	// Level
	level := strings.ToLower(l.config.LogLevel)

	switch level {
	case "debug":
		l.log = l.log.Level(zerolog.DebugLevel)
	case "info":
		l.log = l.log.Level(zerolog.InfoLevel)
	case "warn":
		l.log = l.log.Level(zerolog.WarnLevel)
	case "error":
		l.log = l.log.Level(zerolog.ErrorLevel)
	default:
		l.log = l.log.Level(zerolog.DebugLevel)
	}
}

// NewInstance - Create a new log instance based off configuration of this
// instance
func (l *Log) NewInstance() (logger.Logger, error) {
	return &Log{
		fields: make(map[string]interface{}),
		config: l.config,
		log:    l.log.With().Logger(),
	}, nil
}

// Printf -
func (l *Log) Printf(format string, args ...interface{}) {
	l.log.Printf(format, args...)
}

// Level -
func (l *Log) Level(level Level) {
	if lvl, ok := levelMap[level]; ok {
		l.log = l.log.Level(lvl)
	}
}

// Context - set logging
func (l *Log) Context(key, value string) {
	if value == "" {
		delete(l.fields, key)
		return
	}
	l.fields[key] = value
}

// WithApplicationContext - Shallow copied logger instance with new application
// context and existing context
func (l *Log) WithApplicationContext(value string) logger.Logger {
	ctxLog := *l
	fields := map[string]any{}

	for k, v := range ctxLog.fields {
		fields[k] = v
	}

	fields["application"] = value
	ctxLog.fields = fields
	return &ctxLog
}

// WithPackageContext - Shallow copied logger instance with new package context
// and existing context
func (l *Log) WithPackageContext(value string) logger.Logger {
	ctxLog := *l
	fields := map[string]any{}

	for k, v := range ctxLog.fields {
		fields[k] = v
	}

	fields["package"] = value
	ctxLog.fields = fields
	return &ctxLog
}

// WithFunctionContext - Shallow copied logger instance with new function context
// and existing context
func (l *Log) WithFunctionContext(value string) logger.Logger {
	ctxLog := *l
	fields := map[string]any{}

	for k, v := range ctxLog.fields {
		fields[k] = v
	}

	fields["function"] = value
	ctxLog.fields = fields
	return &ctxLog
}

// WithDurationContext - Shallow copied logger instance with new timing context
// and existing context
func (l *Log) WithDurationContext(value string) logger.Logger {
	ctxLog := *l
	fields := map[string]any{}

	for k, v := range ctxLog.fields {
		fields[k] = v
	}

	fields["duration"] = value
	ctxLog.fields = fields
	return &ctxLog
}

// Debug -
func (l *Log) Debug(msg string, args ...interface{}) {
	ctxLog := l.log.With().Fields(l.fields).Logger()
	ctxLog.Debug().Msgf(msg, args...)
}

// Info -
func (l *Log) Info(msg string, args ...interface{}) {
	ctxLog := l.log.With().Fields(l.fields).Logger()
	ctxLog.Info().Msgf(msg, args...)
}

// Warn -
func (l *Log) Warn(msg string, args ...interface{}) {
	ctxLog := l.log.With().Fields(l.fields).Logger()
	ctxLog.Warn().Msgf(msg, args...)
}

// Error -
func (l *Log) Error(msg string, args ...interface{}) {
	ctxLog := l.log.With().Fields(l.fields).Logger()
	ctxLog.Error().Msgf(msg, args...)
}
