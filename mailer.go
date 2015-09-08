package mailer

import (
	"bytes"
	"fmt"
	"path"
	"strings"

	"net/smtp"
	"pault.ag/go/config"
	"text/template"
)

type MailerRC struct {
	Sender   string `flag:"mailer-sender"    description:"SMTP Sender"`
	Password string `flag:"mailer-password"  description:"SMTP Password"`
	Host     string `flag:"mailer-server"    description:"SMTP Server"`
	Port     int    `flag:"mailer-port"      description:"SMTP Port"`
}

type Mailer struct {
	Config MailerRC
	Root   string
}

type MailerData struct {
	From string
	To   string
	Data interface{}
}

func (m *Mailer) Mail(to []string, mailTemplate string, data interface{}) error {
	auth := smtp.PlainAuth(
		"",
		m.Config.Sender,
		m.Config.Password,
		m.Config.Host,
	)

	byteBuffer := bytes.Buffer{}

	t, err := template.ParseFiles(path.Join(m.Root, mailTemplate))
	if err != nil {
		return err
	}

	if err := t.Execute(&byteBuffer, MailerData{
		From: m.Config.Sender,
		To:   strings.Join(to, ", "),
		Data: data,
	}); err != nil {
		return err
	}

	err = smtp.SendMail(
		fmt.Sprintf("%s:%d", m.Config.Host, m.Config.Port),
		auth,
		m.Config.Sender,
		to,
		byteBuffer.Bytes(),
	)
	return err
}

func NewMailer(root string) (*Mailer, error) {
	mailerRC := MailerRC{}
	if err := config.Load("mailer", &mailerRC); err != nil {
		return nil, err
	}
	return &Mailer{
		Config: mailerRC,
		Root:   root,
	}, nil
}
