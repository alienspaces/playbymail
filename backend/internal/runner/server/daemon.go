package runner

import "context"

// runDaemonFunc is a long running background process that manages the server game loop.
func (rnr *Runner) runDaemonFunc(ctx context.Context, args map[string]any) error {
	l := loggerWithFunctionContext(rnr.Log, "RunDaemon")

	l.Info("(runner) running playbymail daemon")

	return nil
}
