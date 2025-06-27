package smtp

import (
	"fmt"

	"gitlab.com/alienspaces/playbymail/core/config"
	"gitlab.com/alienspaces/playbymail/core/type/emailer"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
)

// FakeSMTP -
type FakeSMTP struct {
	SMTP
}

var _ emailer.Emailer = &FakeSMTP{}

// NewFakeSMTP -
func NewFakeSMTP(l logger.Logger, c config.Config) (*FakeSMTP, error) {
	e := &FakeSMTP{
		SMTP: SMTP{
			config: c,
			log:    l,
		},
	}
	return e, nil
}

func (e *FakeSMTP) Send(msg *emailer.Message) error {
	l := e.logger("Send")
	l.Info("Sending host >%s< from >%s< to >%s< cc >%v< bcc >%v<", e.host, msg.From, msg.To, msg.CC, msg.BCC)
	return nil
}

func (e *FakeSMTP) logger(functionName string) logger.Logger {
	if e.log == nil {
		return nil
	}
	return e.log.WithPackageContext(fmt.Sprintf("(fake) %s", packageName)).WithFunctionContext(functionName)
}
