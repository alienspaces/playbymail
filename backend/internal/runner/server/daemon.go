package runner

import (
	"context"

	"gitlab.com/alienspaces/playbymail/internal/utils/logging"
)

// runDaemonFunc is a long running background process that manages the server game loop.
func (rnr *Runner) runDaemonFunc(ctx context.Context, args map[string]any) error {
	l := logging.LoggerWithFunctionContext(rnr.Log, "runner", "runDaemonFunc")

	if rnr.JobClient == nil {
		l.Warn("(playbymail) runner does not have a job client, not running")
		return nil
	}

	if err := rnr.JobClient.Start(ctx); err != nil {
		return err
	}

	l.Info("(playbymail) job client started")

	<-rnr.JobClient.Stopped()

	l.Info("(playbymail) job client stopped")

	return nil
}
