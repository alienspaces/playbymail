package fake

import (
	"gitlab.com/alienspaces/playbymail/core/config"
	"gitlab.com/alienspaces/playbymail/core/type/emailer"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
)

// Fake -
type Fake struct {
	log    logger.Logger
	config config.Config
	host   string
}

var _ emailer.Emailer = &Fake{}

// New -
func New(l logger.Logger, c config.Config) (*Fake, error) {
	e := &Fake{
		config: c,
		log:    l,
	}
	return e, nil
}

func (e *Fake) Send(msg *emailer.Message) error {
	l := e.logger("Send")
	l.Info("Sending host >%s< from >%s< to >%s< cc >%v< bcc >%v<", e.host, msg.From, msg.To, msg.CC, msg.BCC)
	return nil
}

func (e *Fake) logger(functionName string) logger.Logger {
	if e.log == nil {
		return nil
	}
	return e.log.WithPackageContext("(fake)").WithFunctionContext(functionName)
}
