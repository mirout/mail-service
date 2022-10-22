package email

import (
	"fmt"
	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
)

type MailSender interface {
	SendMailToUser(user string, message string) error
}

type mailSender struct {
	smtpClient *smtp.Client
	author     string
}

func NewSender(addr, username, password string) (*mailSender, error) {
	cl, err := smtp.Dial(addr)
	if err != nil {
		return nil, fmt.Errorf("can't dial: %w", err)
	}

	err = cl.StartTLS(nil)
	if err != nil {
		return nil, fmt.Errorf("can't start tls: %w", err)
	}

	auth := sasl.NewPlainClient("", username, password)
	err = cl.Auth(auth)
	if err != nil {
		return nil, fmt.Errorf("can't auth: %w", err)
	}

	return &mailSender{
		smtpClient: cl,
		author:     username,
	}, nil
}

func (m *mailSender) SendMailToUser(user string, message string) error {
	err := m.smtpClient.Mail(m.author, nil)
	if err != nil {
		return fmt.Errorf("can't set author: %w", err)
	}

	err = m.smtpClient.Rcpt(user)
	if err != nil {
		return fmt.Errorf("can't set recipient: %w", err)
	}

	wc, err := m.smtpClient.Data()
	if err != nil {
		return fmt.Errorf("can't create data writer: %w", err)
	}

	_, err = fmt.Fprintf(wc, message)
	if err != nil {
		return fmt.Errorf("can't write message: %w", err)
	}

	err = wc.Close()
	if err != nil {
		return fmt.Errorf("can't close writer: %w", err)
	}

	return nil
}
