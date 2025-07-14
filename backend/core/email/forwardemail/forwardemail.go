package forwardemail

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"gitlab.com/alienspaces/playbymail/core/config"
	"gitlab.com/alienspaces/playbymail/core/type/emailer"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
)

const (
	packageName = "forwardemail"
	apiBaseURL  = "https://api.forwardemail.net/v1/emails"
)

type ForwardEmail struct {
	log    logger.Logger
	config config.Config
	apiKey string
}

var _ emailer.Emailer = &ForwardEmail{}

func New(l logger.Logger, c config.Config) (*ForwardEmail, error) {
	f := &ForwardEmail{
		log:    l,
		config: c,
		apiKey: c.ForwardEmailAPIKey, // Add this to your config
	}
	if err := f.Init(); err != nil {
		return nil, err
	}
	return f, nil
}

func (f *ForwardEmail) Init() error {
	if f.apiKey == "" {
		return fmt.Errorf("missing Forward Email API key")
	}
	return nil
}

type forwardEmailRequest struct {
	From    string   `json:"from"`
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	Text    string   `json:"text,omitempty"`
	HTML    string   `json:"html,omitempty"`
}

func (f *ForwardEmail) Send(msg *emailer.Message) error {
	l := f.logger("Send")
	l.Info("sending from >%s< to >%v<", msg.From, msg.To)

	// For now, only support plain text body
	reqBody := forwardEmailRequest{
		From:    msg.From,
		To:      msg.To,
		Subject: msg.Subject,
		Text:    msg.Body,
	}
	body, err := json.Marshal(reqBody)
	if err != nil {
		l.Warn("failed to marshal request body >%v<", err)
		return err
	}

	req, err := http.NewRequest("POST", apiBaseURL, bytes.NewReader(body))
	if err != nil {
		l.Warn("failed to create request >%v<", err)
		return err
	}
	req.SetBasicAuth(f.apiKey, "")
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		l.Warn("failed to send request >%v<", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// Read and log the response body for more details
		respBody := new(bytes.Buffer)
		_, _ = respBody.ReadFrom(resp.Body)
		l.Warn("forwardemail API returned status %d, body: %s", resp.StatusCode, respBody.String())
		return fmt.Errorf("forwardemail API error: %d, body: %s", resp.StatusCode, respBody.String())
	}

	l.Info("successfully sent email via forwardemail")
	return nil
}

func (f *ForwardEmail) logger(functionName string) logger.Logger {
	if f.log == nil {
		return nil
	}
	return f.log.WithPackageContext(packageName).WithFunctionContext(functionName)
}
