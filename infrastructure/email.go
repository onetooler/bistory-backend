package infrastructure

import (
	"bytes"
	"fmt"
	"html/template"
	"os"

	"github.com/onetooler/bistory-backend/config"
	"github.com/onetooler/bistory-backend/logger"
	"github.com/onetooler/bistory-backend/util"
	"gopkg.in/gomail.v2"
)

type EmailSender interface {
	SendEmail(to, subject, template string, body any) error
}

type emailSender struct {
	account   string
	host      string
	port      int
	username  string
	password  string
	dialer    *gomail.Dialer
	templates map[string]*template.Template
}

type disabledEmailsender struct{}

func NewEmailSender(logger logger.Logger, conf *config.Config, templates map[string]*template.Template) EmailSender {
	if !conf.Email.Enabled {
		return &disabledEmailsender{}
	}

	logger.GetZapLogger().Infof("Try email smtp connection")
	d := gomail.NewDialer(conf.Email.Host, conf.Email.Port, conf.Email.Username, conf.Email.Password)
	closer, err := d.Dial()
	if err != nil {
		logger.GetZapLogger().Errorf("Failure email smtp connection. error: %s", err.Error())
		os.Exit(config.ErrExitStatus)
	}
	util.Check(closer.Close)
	logger.GetZapLogger().Infof("Success email smtp connection, %s:%s", conf.Email.Host, conf.Email.Port)

	return &emailSender{
		account:   conf.Email.Account,
		host:      conf.Email.Host,
		port:      conf.Email.Port,
		username:  conf.Email.Username,
		password:  conf.Email.Password,
		dialer:    d,
		templates: templates,
	}
}

func (e emailSender) SendEmail(to, subject, template string, body any) error {
	t, ok := e.templates[template]
	if !ok {
		return fmt.Errorf("template not found: %s", template)
	}
	var buf bytes.Buffer
	err := t.Execute(&buf, body)
	if err != nil {
		return err
	}

	msg := gomail.NewMessage()
	msg.SetHeader("From", e.account)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", buf.String())

	return e.dialer.DialAndSend(msg)
}

func (e disabledEmailsender) SendEmail(to, subject, template string, body any) error {
	return fmt.Errorf("email sender disabled by config")
}
