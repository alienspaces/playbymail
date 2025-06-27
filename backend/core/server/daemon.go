package server

import (
	"context"
)

// runDaemon is the default daemon process which starts and stops a job client.
func (rnr *Runner) runDaemon(pctx context.Context, args map[string]interface{}) error {
	l := Logger(rnr.Log, "runDaemon")

	if rnr.JobClient == nil {
		l.Warn("(core) runner does not have a job client, not running")
		return nil
	}

	if err := rnr.JobClient.Start(pctx); err != nil {
		return err
	}

	l.Info("(core) job client started")

	<-rnr.JobClient.Stopped()

	l.Info("(core) job client stopped")

	return nil
}
