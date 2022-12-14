package model

import (
	"database/sql"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Email     string    `json:"email" db:"email"`
	FirstName string    `json:"first_name" db:"first_name"`
	LastName  string    `json:"last_name" db:"last_name"`
	CreatedAt string    `json:"created_at" db:"created_at"`
}

type Group struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	CreatedAt string    `json:"created_at" db:"created_at"`
}

type Mail struct {
	ID        uuid.UUID      `json:"id" db:"id"`
	ToUserId  uuid.UUID      `json:"to_user_id" db:"to_user_id"`
	Subject   string         `json:"subject" db:"subject"`
	Body      string         `json:"body" db:"body"`
	CreatedAt string         `json:"created_at" db:"created_at"`
	SentAt    sql.NullString `json:"sent_at" db:"sent_at"`
	Watched   bool           `json:"watched" db:"watched"`
}

type MailJson struct {
	Subject string `json:"subject"`
	Body    string `json:"body"`
	SendAt  string `json:"send_at"`
}

func (m *MailJson) Validate() error {
	return validation.ValidateStruct(m,
		validation.Field(&m.Subject, validation.Required),
		validation.Field(&m.Body, validation.Required),
		validation.Field(&m.SendAt, validation.Date(time.RFC3339)),
	)
}

type MailWithUser struct {
	ID        uuid.UUID      `json:"id" db:"id"`
	FirstName string         `json:"first_name" db:"first_name"`
	LastName  string         `json:"last_name" db:"last_name"`
	Email     string         `json:"email" db:"email"`
	Subject   string         `json:"subject" db:"subject"`
	Body      string         `json:"body" db:"body"`
	CreatedAt string         `json:"created_at" db:"created_at"`
	SentAt    sql.NullString `json:"sent_at" db:"sent_at"`
}
