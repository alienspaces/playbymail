package logger

// Logger -
type Logger interface {
	NewInstance() (Logger, error)
	Context(key, value string)
	WithApplicationContext(value string) Logger
	WithDurationContext(value string) Logger
	WithPackageContext(value string) Logger
	WithFunctionContext(value string) Logger
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}
