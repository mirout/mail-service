package email

import (
	"bytes"
	"context"
	"fmt"
	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
	"html/template"
	"mail-service/internal/model"
	"mail-service/internal/storage"
	"time"
)

type MailSender interface {
	SendMailToUser(ctx context.Context, mail model.Mail) error
	CreateDelayedMail(ctx context.Context, mail model.Mail, delay time.Time) error
}

type SmtpConfig struct {
	Addr     string
	Username string
	Password string
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

type MailWorker struct {
	smtpClient *smtp.Client
	author     string

	mails storage.Mail
	users storage.User

	queue DelayedQueue
}

func NewWorker(smtpConfig SmtpConfig, mails storage.Mail, users storage.User, q DelayedQueue) (*MailWorker, error) {
	cl, err := smtp.Dial(smtpConfig.Addr)
	if err != nil {
		return nil, fmt.Errorf("can't dial: %w", err)
	}

	err = cl.StartTLS(nil)
	if err != nil {
		return nil, fmt.Errorf("can't start tls: %w", err)
	}

	auth := sasl.NewPlainClient("", smtpConfig.Username, smtpConfig.Password)
	err = cl.Auth(auth)
	if err != nil {
		return nil, fmt.Errorf("can't auth: %w", err)
	}

	return &MailWorker{
		smtpClient: cl,
		author:     smtpConfig.Username,
		mails:      mails,
		users:      users,
		queue:      q,
	}, nil
}

func (m *MailWorker) Close() error {
	err := m.smtpClient.Quit()
	if err != nil {
		return m.smtpClient.Close()
	}
	return nil
}

func (m *MailWorker) SendMailToUser(ctx context.Context, mail model.Mail) error {
	id, err := m.mails.CreateMail(ctx, mail)
	if err != nil {
		return fmt.Errorf("can't create mail: %w", err)
	}

	user, err := m.users.GetUser(ctx, mail.ToUserId)
	if err != nil {
		return fmt.Errorf("can't get user: %w", err)
	}

	err = m.Send(user, mail)
	if err != nil {
		return fmt.Errorf("can't send mail: %w", err)
	}

	err = m.mails.MarkAsSent(ctx, id, time.Now())
	if err != nil {
		return fmt.Errorf("can't mark mail as sent: %w", err)
	}

	return nil
}

func buildHtml(user model.User, mail model.Mail) (bytes.Buffer, error) {
	tmpl, err := template.ParseFiles("templates/template.html")
	if err != nil {
		return bytes.Buffer{}, fmt.Errorf("can't parse template: %w", err)
	}

	var b bytes.Buffer
	err = tmpl.Execute(&b, struct {
		FirstName string
		LastName  string
		Body      string
		Id        string
	}{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Body:      mail.Body,
		Id:        mail.ID.String(),
	})
	if err != nil {
		return bytes.Buffer{}, fmt.Errorf("can't execute template: %w", err)
	}
	return b, nil

}

func (m *MailWorker) Send(user model.User, mail model.Mail) error {
	err := m.smtpClient.Mail(m.author, nil)
	if err != nil {
		return fmt.Errorf("can't set author: %w", err)
	}

	err = m.smtpClient.Rcpt(user.Email)
	if err != nil {
		return fmt.Errorf("can't set recipient: %w", err)
	}

	wc, err := m.smtpClient.Data()
	if err != nil {
		return fmt.Errorf("can't create data writer: %w", err)
	}

	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body, err := buildHtml(user, mail)
	if err != nil {
		return fmt.Errorf("can't build html: %w", err)
	}
	msg := fmt.Sprintf("Subject: %s\n%s\n%s\n", mail.Subject, mime, body.String())

	_, err = fmt.Fprintf(wc, msg)
	if err != nil {
		return fmt.Errorf("can't write message: %w", err)
	}

	err = wc.Close()
	if err != nil {
		return fmt.Errorf("can't close writer: %w", err)
	}
	return nil
}

func (m *MailWorker) CreateDelayedMail(ctx context.Context, mail model.Mail, delay time.Time) error {
	id, err := m.mails.CreateMail(ctx, mail)
	if err != nil {
		return fmt.Errorf("can't create mail: %w", err)
	}
	err = m.queue.Enqueue(ctx, Mail{ID: id}, delay.Unix())
	if err != nil {
		return fmt.Errorf("can't enqueue mail: %w", err)
	}
	return nil
}

func (m *MailWorker) Run() {
	ch := m.queue.GetReadyChannel()

	for {
		mails := <-ch
		for _, mail := range mails {
			mail, err := m.mails.GetMailById(context.Background(), mail.ID)
			if err != nil {
				fmt.Printf("can't get mail: %v", err)
				continue
			}
			user, err := m.users.GetUser(context.Background(), mail.ToUserId)
			if err != nil {
				fmt.Printf("can't get user: %v", err)
				continue
			}
			err = m.Send(user, mail)
			if err != nil {
				fmt.Printf("can't send mail: %v", err)
				continue
			}
			err = m.mails.MarkAsSent(context.Background(), mail.ID, time.Now())
			if err != nil {
				fmt.Printf("can't mark mail as sent: %v", err)
				continue
			}
		}
	}
}
