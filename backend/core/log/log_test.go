package log

import (
	"testing"

	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/playbymail/core/config"
)

func TestNewLogger(t *testing.T) {
	cfg := config.Config{}
	err := config.Parse(&cfg)
	require.NoError(t, err, "config.Parse returns without error")

	l, err := NewLogger(cfg)
	require.NoError(t, err, "NewLogger returns without error")
	require.NotNil(t, l, "NewLogger is not nil")

	l.Debug("Test level >%s<", "debug")
	l.Info("Test level >%s<", "info")
	l.Warn("Test level >%s<", "warn")
	l.Error("Test level >%s<", "error")

	l.Context("correlation-id", "abcdefg")

	l.Debug("Test level >%s<", "debug")

	l.Context("correlation-id", "hijklmn")

	l.Debug("Test level >%s<", "debug")

	l.Level(ErrorLevel)

	l.Debug("Test level >%s<", "debug")
	l.Info("Test level >%s<", "info")
	l.Warn("Test level >%s<", "warn")
	l.Error("Test level >%s<", "error")
}

func TestNewLoggerWithConfig(t *testing.T) {

	l, err := NewLoggerWithConfig(config.Config{
		LogLevel:    "debug",
		LogIsPretty: true,
	})
	require.NoError(t, err, "NewLogger returns without error")
	require.NotNil(t, l, "NewLogger is not nil")

	l.Debug("Test level >%s<", "debug")
	l.Info("Test level >%s<", "info")
	l.Warn("Test level >%s<", "warn")
	l.Error("Test level >%s<", "error")

	l.Context("correlation-id", "abcdefg")

	l.Debug("Test level >%s<", "debug")

	l.Context("correlation-id", "hijklmn")

	l.Debug("Test level >%s<", "debug")

	l.Level(ErrorLevel)

	l.Debug("Test level >%s<", "debug")
	l.Info("Test level >%s<", "info")
	l.Warn("Test level >%s<", "warn")
	l.Error("Test level >%s<", "error")
}
