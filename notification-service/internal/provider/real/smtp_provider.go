package real

import (
	"context"
	"fmt"
	"net/smtp"
)

// SMTPProvider sends real emails via an SMTP server.
// It is selected when PROVIDER_MODE=REAL.
type SMTPProvider struct {
	host string
	port string
	user string
	pass string
	from string
}

func New(host, port, user, pass, from string) *SMTPProvider {
	return &SMTPProvider{host: host, port: port, user: user, pass: pass, from: from}
}

func (p *SMTPProvider) Send(_ context.Context, to, subject, body string) error {
	addr := fmt.Sprintf("%s:%s", p.host, p.port)
	auth := smtp.PlainAuth("", p.user, p.pass, p.host)
	msg := []byte(fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s",
		p.from, to, subject, body,
	))
	return smtp.SendMail(addr, auth, p.from, []string{to}, msg)
}
