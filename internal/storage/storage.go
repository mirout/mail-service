package storage

import (
	"context"
	"github.com/google/uuid"
	"mail-service/internal/model"
)

type StorageService interface {
	CreateUser(ctx context.Context, user model.User) (uuid.UUID, error)
}
