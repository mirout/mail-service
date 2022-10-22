package storage

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5"
	"github.com/jmoiron/sqlx"
	"mail-service/internal/model"
)

type SqlStorage struct {
	db *sqlx.DB
}

func NewSqlStorage(ctx context.Context, driverName, dataSourceName string) (*SqlStorage, error) {
	db, err := sqlx.Open(driverName, dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("can't open db: %w", err)
	}

	if err = db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("can't ping db: %w", err)
	}

	return &SqlStorage{db: db}, nil
}

func (s *SqlStorage) CreateUser(ctx context.Context, user model.User) (uuid.UUID, error) {
	result, err := s.db.NamedQueryContext(ctx, `
		INSERT INTO users (email, first_name, last_name)
		VALUES (:email, :first_name, :last_name)
		RETURNING id
	`, user)
	if err != nil {
		return uuid.Nil, fmt.Errorf("can't create user: %w", err)
	}

	var id uuid.UUID

	if result.Next() {
		if err = result.Scan(&id); err != nil {
			return uuid.Nil, fmt.Errorf("can't get id: %w", err)
		}
	}

	return id, nil
}
