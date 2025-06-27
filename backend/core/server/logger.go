package server

import (
	"gitlab.com/alienspaces/playbymail/core/type/logger"
)

func Logger(l logger.Logger, functionName string) logger.Logger {
	return l.WithFunctionContext(functionName)
}
