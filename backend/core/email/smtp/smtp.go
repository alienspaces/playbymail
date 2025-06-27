package smtp

import (
	"fmt"
	"io"
	"log"
	"net/smtp"

	"gitlab.com/alienspaces/playbymail/core/config"
	"gitlab.com/alienspaces/playbymail/core/type/emailer"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
)

const (
	packageName = "smtp"
	// MimeHTML - Content type
	MimeHTML = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
)

// SMTP -
type SMTP struct {
	log        logger.Logger
	config     config.Config
	smtpClient *smtp.Client
	host       string
}

var _ emailer.Emailer = &SMTP{}

// NewSMTP -
func NewSMTP(l logger.Logger, c config.Config) (*SMTP, error) {
	e := &SMTP{
		config:     c,
		log:        l,
		smtpClient: &smtp.Client{},
	}

	err := e.Init()
	if err != nil {
		l.Warn("failed email init >%v<", err)
		return nil, err
	}

	return e, nil
}

func (e *SMTP) Init() error {
	l := e.logger("Connect")
	l.Info("Initialising")

	host := e.config.SMTPHost
	if host == "" {
		err := fmt.Errorf("missing host, cannot init initialise emailer")
		l.Warn(err.Error())
		return err
	}

	e.host = host

	l.Info("SMPT host >%v<", e.host)

	err := e.Connect()
	if err != nil {
		l.Warn("failed connecting to host >%s< >%v<", e.host, err)
		return err
	}

	return nil
}

func (e *SMTP) Connect() error {
	l := e.logger("Connect")
	l.Debug("Connecting")

	c, err := smtp.Dial(e.host)
	if err != nil {
		l.Warn("failed dialing host >%v<", err)
		return err
	}

	e.smtpClient = c
	return nil
}

func (e *SMTP) Send(msg *emailer.Message) (err error) {
	l := e.logger("Send")
	l.Info("Sending host >%s< from >%s< to >%s< cc >%v< bcc >%v<", e.host, msg.From, msg.To, msg.CC, msg.BCC)

	err = e.Connect()
	if err != nil {
		l.Warn("failed connecting to host >%s< >%v<", e.host, err)
		return err
	}
	defer e.Quit()

	c := e.smtpClient

	if err := c.Mail(msg.From); err != nil {
		l.Warn("failed setting from >%s< >%v<", msg.From, err)
		return err
	}

	for _, t := range msg.To {
		if err := c.Rcpt(t); err != nil {
			l.Warn("failed setting recipient >%s< >%v<", t, err)
			return err
		}
	}

	var wc io.WriteCloser
	wc, err = c.Data()
	if err != nil {
		return err
	}

	defer func() {
		closeErr := wc.Close()
		if closeErr != nil {
			l.Warn("failed closing client data writer >%v<", closeErr)
		}
		if err != nil {
			return
		}
		err = closeErr
	}()

	var msgBytes []byte
	msgBytes, err = msg.Bytes()
	if err != nil {
		return err
	}

	l.Info("Message length >%d<", len(msgBytes))

	_, err = wc.Write(msgBytes)
	if err != nil {
		return err
	}

	return err
}

func (e *SMTP) Quit() error {
	l := e.logger("Quit")

	l.Debug("Quitting")

	// if not connected return nil
	if e.smtpClient.Text == nil {
		return nil
	}
	err := e.smtpClient.Quit()
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func (e *SMTP) logger(functionName string) logger.Logger {
	if e.log == nil {
		return nil
	}
	return e.log.WithPackageContext(packageName).WithFunctionContext(functionName)
}
