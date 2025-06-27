package logging

import (
	"gitlab.com/alienspaces/playbymail/core/type/logger"
)

// LoggerWithFunctionContext - Returns a logger with package context and provided function context
// This is a common utility used by both API and CLI runners
func LoggerWithFunctionContext(l logger.Logger, packageName, functionName string) logger.Logger {
	if l == nil {
		return nil
	}
	return l.WithPackageContext(packageName).WithFunctionContext(functionName)
}
