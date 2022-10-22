package storage

import (
	"context"
	"github.com/google/uuid"
	"mail-service/internal/model"
)

type UserService interface {
	CreateUser(ctx context.Context, user model.User) (uuid.UUID, error)
	GetUser(ctx context.Context, id uuid.UUID) (model.User, error)
	GetUserByEmail(ctx context.Context, email string) (model.User, error)
}

type GroupService interface {
	CreateGroup(ctx context.Context, user model.Group) (uuid.UUID, error)
	GetGroupById(ctx context.Context, id uuid.UUID) (model.Group, error)
	GetGroupByName(ctx context.Context, email string) (model.Group, error)
	AddUserToGroup(ctx context.Context, userID, groupID uuid.UUID) error
	RemoveUserFromGroup(ctx context.Context, userID, groupID uuid.UUID) error
	GetUsersByGroup(ctx context.Context, groupID uuid.UUID) ([]model.User, error)
}
