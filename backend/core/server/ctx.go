package server

type contextKey int

const (
	ctxKeyAuth          contextKey = 1
	ctxKeyCorrelationID contextKey = 2
)
