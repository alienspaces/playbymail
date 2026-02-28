package emailer

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"mime/multipart"
	"strings"
)

type Message struct {
	From        string
	To          []string
	CC          []string
	BCC         []string
	Subject     string
	Body        string
	Attachments []Attachment
}

type Attachment struct {
	Name        string
	Content     []byte
	ContentType string
}

func NewMessage(from string, to []string, cc []string, bcc []string, subject string, body string, attachments []Attachment) *Message {
	m := &Message{
		From:        from,
		To:          to,
		CC:          cc,
		BCC:         bcc,
		Subject:     subject,
		Body:        body,
		Attachments: []Attachment{},
	}
	if len(attachments) != 0 {
		m.Attachments = attachments
	}
	return m
}

func (m *Message) Attach(attachment Attachment) error {
	m.Attachments = append(m.Attachments, attachment)
	return nil
}

func (m *Message) Bytes() ([]byte, error) {

	buf := bytes.NewBuffer(nil)
	withAttachments := len(m.Attachments) > 0

	fmt.Fprintf(buf, "From:%s\r\n", m.From)

	fmt.Fprintf(buf, "To:%s\r\n", strings.Join(m.To, ","))

	fmt.Fprintf(buf, "Subject:%s\r\n", m.Subject)

	if len(m.CC) > 0 {
		fmt.Fprintf(buf, "Cc:%s\r\n", strings.Join(m.CC, ","))
	}

	if len(m.BCC) > 0 {
		fmt.Fprintf(buf, "Bcc:%s\r\n", strings.Join(m.BCC, ","))
	}

	buf.WriteString("MIME-Version: 1.0\r\n")
	writer := multipart.NewWriter(buf)
	boundary := writer.Boundary()
	if withAttachments {
		fmt.Fprintf(buf, "Content-Type: multipart/mixed; boundary=%s\r\n\r\n", boundary)
		fmt.Fprintf(buf, "--%s\r\n", boundary)
	}

	// TODO CX-103: Support different content types for message body
	if len(m.Body) > 0 {
		buf.WriteString("Content-Type: text/plain; charset=utf-8\r\n")
		buf.WriteString("\r\n")
		buf.WriteString(m.Body)
	}

	if withAttachments {
		for _, attachment := range m.Attachments {
			fmt.Fprintf(buf, "\r\n--%s\r\n", boundary)
			fmt.Fprintf(buf, "Content-Type:%s\r\n", attachment.ContentType)
			buf.WriteString("Content-Transfer-Encoding: base64\r\n")
			fmt.Fprintf(buf, "Content-Disposition: attachment; filename=%s\r\n", attachment.Name)
			buf.WriteString("\r\n")

			b := make([]byte, base64.StdEncoding.EncodedLen(len(attachment.Content)))
			base64.StdEncoding.Encode(b, attachment.Content)
			buf.Write(b)

			fmt.Fprintf(buf, "\r\n--%s", boundary)
		}

		buf.WriteString("--")
	}

	return buf.Bytes(), nil
}
