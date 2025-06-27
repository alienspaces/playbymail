package sendgrid

import (
	b64 "encoding/base64"
	"fmt"
	"log"

	sendgridgo "github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"gitlab.com/alienspaces/playbymail/core/config"
	"gitlab.com/alienspaces/playbymail/core/type/emailer"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
)

const (
	packageName = "sendgrid"
)

type Sendgrid struct {
	log            logger.Logger
	config         config.Config
	sendgridClient *sendgridgo.Client
}

var _ emailer.Emailer = &Sendgrid{}

func New(l logger.Logger, c config.Config) (*Sendgrid, error) {
	e := &Sendgrid{
		config:         c,
		log:            l,
		sendgridClient: &sendgridgo.Client{},
	}

	err := e.Init()
	if err != nil {
		l.Warn("failed sendgrid init >%v<", err)
		return nil, err
	}

	return e, nil
}

func (e *Sendgrid) Init() error {
	l := e.logger("Init")
	l.Info("Initialising")

	apiKey := e.config.SendgridAPIKey
	if apiKey == "" {
		err := fmt.Errorf("missing SENDGRID_API_KEY, failed to initialise sendgrid")
		l.Warn(err.Error())
		return err
	}

	e.sendgridClient = sendgridgo.NewSendClient(apiKey)
	return nil
}

func (e *Sendgrid) Send(msg *emailer.Message) error {
	l := e.logger("Send")
	l.Info("Sending from >%s< to >%s< cc >%v< bcc >%v<", msg.From, msg.To, msg.CC, msg.BCC)

	mailer := mail.NewV3Mail()
	mailer.SetFrom(mail.NewEmail("", msg.From))
	p := e.SetPersonalization(SetPersonalizationArgs{
		Tos:     msg.To,
		CCs:     msg.CC,
		BCCs:    msg.BCC,
		Subject: msg.Subject,
	})
	mailer.AddPersonalizations(p)

	mailer.AddAttachment(e.ConvertAttachments(msg.Attachments)...)

	content := mail.NewContent("text/html", msg.Body)
	mailer.AddContent(content)

	response, err := e.sendgridClient.Send(mailer)
	if err != nil || (response.StatusCode != 200 && response.StatusCode != 202) {
		l.Warn("failed to send email through sendgrid >%v< Response >%+v<", err, response)
		log.Println(err)
		return err
	}

	l.Info("Successfully sent email through sendgrid >%+v<", response)
	return nil
}

func (e *Sendgrid) SendEmail(from, to, subject, emailBody string) error {
	l := e.logger("SendEmail")
	l.Info("Sending from >%s< to >%s< cc >%v< bcc >%v<", from, to)

	mailer := mail.NewV3Mail()
	mailer.SetFrom(mail.NewEmail("", from))
	p := e.SetPersonalization(SetPersonalizationArgs{
		Tos:     []string{to},
		Subject: subject,
	})
	mailer.AddPersonalizations(p)

	content := mail.NewContent("text/plain", emailBody)
	mailer.AddContent(content)

	response, err := e.sendgridClient.Send(mailer)
	if err != nil || (response.StatusCode != 200 && response.StatusCode != 202) {
		l.Warn("failed to send email through sendgrid >%v< Response >%+v<", err, response)
		log.Println(err)
		return err
	}

	l.Info("Successfully sent email through sendgrid >%+v<", response)

	return nil
}

func (e *Sendgrid) SendHTMLEmail(from, to, subject, emailBody string) error {
	l := e.logger("SendHTMLEmail")

	l.Info("Sending from >%s< to >%s< cc >%v< bcc >%v<", from, to)

	mailer := mail.NewV3Mail()
	mailer.SetFrom(mail.NewEmail("", from))
	p := e.SetPersonalization(SetPersonalizationArgs{
		Tos:     []string{to},
		Subject: subject,
	})
	mailer.AddPersonalizations(p)

	content := mail.NewContent("text/html", emailBody)
	mailer.AddContent(content)

	response, err := e.sendgridClient.Send(mailer)
	if err != nil || (response.StatusCode != 200 && response.StatusCode != 202) {
		l.Warn("failed to send email through sendgrid >%v< Response >%+v<", err, response)
		log.Println(err)
		return err
	}

	l.Info("Successfully sent email through sendgrid >%+v<", response)
	return nil
}

func (e *Sendgrid) ConvertAttachments(attachments []emailer.Attachment) []*mail.Attachment {
	var mailAttachments []*mail.Attachment
	for _, attachment := range attachments {
		a := mail.NewAttachment()
		a.SetContent(b64.StdEncoding.EncodeToString(attachment.Content))
		a.SetType(attachment.ContentType)
		a.SetFilename(attachment.Name)
		a.SetDisposition("attachment")
		mailAttachments = append(mailAttachments, a)
	}

	return mailAttachments
}

type SetPersonalizationArgs struct {
	Tos     []string
	CCs     []string
	BCCs    []string
	Subject string
}

func (e *Sendgrid) SetPersonalization(args SetPersonalizationArgs) *mail.Personalization {
	p := mail.NewPersonalization()

	p.AddTos(e.ConvertEmails(args.Tos)...)
	p.AddCCs(e.ConvertEmails(args.CCs)...)
	p.AddBCCs(e.ConvertEmails(args.BCCs)...)
	p.Subject = args.Subject

	return p
}

func (e *Sendgrid) ConvertEmails(emails []string) []*mail.Email {
	var emailers []*mail.Email

	for _, email := range emails {
		emailers = append(emailers, &mail.Email{Address: email})
	}

	return emailers
}

func (e *Sendgrid) logger(functionName string) logger.Logger {
	if e.log == nil {
		return nil
	}
	return e.log.WithPackageContext(packageName).WithFunctionContext(functionName)
}
