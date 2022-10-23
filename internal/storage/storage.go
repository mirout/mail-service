package storage

import (
	"context"
	"github.com/google/uuid"
	"mail-service/internal/model"
	"time"
)

type User interface {
	CreateUser(ctx context.Context, user model.User) (uuid.UUID, error)
	GetUser(ctx context.Context, id uuid.UUID) (model.User, error)
	GetUserByEmail(ctx context.Context, email string) (model.User, error)
}

type Group interface {
	CreateGroup(ctx context.Context, user model.Group) (uuid.UUID, error)
	GetGroupById(ctx context.Context, id uuid.UUID) (model.Group, error)
	GetGroupByName(ctx context.Context, email string) (model.Group, error)
	AddUserToGroup(ctx context.Context, userID, groupID uuid.UUID) error
	RemoveUserFromGroup(ctx context.Context, userID, groupID uuid.UUID) error
	GetUsersByGroup(ctx context.Context, groupID uuid.UUID) ([]model.User, error)
}

type Mail interface {
	CreateMail(ctx context.Context, mail model.Mail) (uuid.UUID, error)
	MarkAsSent(ctx context.Context, id uuid.UUID, time time.Time) error
	MarkAsWatched(ctx context.Context, id uuid.UUID) error
	GetMailById(ctx context.Context, id uuid.UUID) (model.Mail, error)
	GetMailsBySentTo(ctx context.Context, userID uuid.UUID) ([]model.Mail, error)
	GetMailWithUser(ctx context.Context, id uuid.UUID) (model.MailWithUser, error)
}
