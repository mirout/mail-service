package storage

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"mail-service/internal/model"
	"time"
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

func (s *SqlStorage) Close() error {
	return s.db.Close()
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

func (s *SqlStorage) GetUser(ctx context.Context, id uuid.UUID) (model.User, error) {
	var user model.User

	if err := s.db.GetContext(ctx, &user, `
		SELECT * FROM users WHERE id = $1
	`, id); err != nil {
		return model.User{}, fmt.Errorf("can't get user: %w", err)
	}

	return user, nil
}

func (s *SqlStorage) GetUserByEmail(ctx context.Context, email string) (model.User, error) {
	var user model.User

	if err := s.db.GetContext(ctx, &user, `
		SELECT * FROM users WHERE email = $1
	`, email); err != nil {
		return model.User{}, fmt.Errorf("can't get user: %w", err)
	}

	return user, nil
}

func (s *SqlStorage) CreateGroup(ctx context.Context, group model.Group) (uuid.UUID, error) {
	result, err := s.db.NamedQueryContext(ctx, `
		INSERT INTO groups (name)
		VALUES (:name)
		RETURNING id
	`, group)
	if err != nil {
		return uuid.Nil, fmt.Errorf("can't create group: %w", err)
	}

	var id uuid.UUID

	if result.Next() {
		if err = result.Scan(&id); err != nil {
			return uuid.Nil, fmt.Errorf("can't get id: %w", err)
		}
	}

	return id, nil
}

func (s *SqlStorage) GetGroupById(ctx context.Context, id uuid.UUID) (model.Group, error) {
	var group model.Group

	if err := s.db.GetContext(ctx, &group, `
		SELECT * FROM groups WHERE id = $1
	`, id); err != nil {
		return model.Group{}, fmt.Errorf("can't get group: %w", err)
	}

	return group, nil
}

func (s *SqlStorage) GetGroupByName(ctx context.Context, name string) (model.Group, error) {
	var group model.Group

	if err := s.db.GetContext(ctx, &group, `
		SELECT * FROM groups WHERE name = $1
	`, name); err != nil {
		return model.Group{}, fmt.Errorf("can't get group: %w", err)
	}

	return group, nil
}

func (s *SqlStorage) AddUserToGroup(ctx context.Context, userID, groupID uuid.UUID) error {
	if _, err := s.db.ExecContext(ctx, `
		INSERT INTO users_groups (user_id, group_id)
		VALUES ($1, $2)
	`, userID, groupID); err != nil {
		return fmt.Errorf("can't add user to group: %w", err)
	}

	return nil
}

func (s *SqlStorage) RemoveUserFromGroup(ctx context.Context, userID, groupID uuid.UUID) error {
	if _, err := s.db.ExecContext(ctx, `
		DELETE FROM users_groups WHERE user_id = $1 AND group_id = $2
	`, userID, groupID); err != nil {
		return fmt.Errorf("can't remove user from group: %w", err)
	}

	return nil
}

func (s *SqlStorage) GetUsersByGroup(ctx context.Context, groupID uuid.UUID) ([]model.User, error) {
	var users []model.User

	if err := s.db.SelectContext(ctx, &users, `
		SELECT u.* FROM users u
		INNER JOIN users_groups ug ON u.id = ug.user_id
		WHERE ug.group_id = $1
	`, groupID); err != nil {
		return nil, fmt.Errorf("can't get users by group: %w", err)
	}

	return users, nil
}

func (s *SqlStorage) CreateMail(ctx context.Context, mail model.Mail) (uuid.UUID, error) {
	result, err := s.db.NamedQueryContext(ctx, `
		INSERT INTO mails (subject, body, to_user_id)
		VALUES (:subject, :body, :to_user_id)
		RETURNING id
	`, mail)
	if err != nil {
		return uuid.Nil, fmt.Errorf("can't create mail: %w", err)
	}

	var id uuid.UUID

	if result.Next() {
		if err = result.Scan(&id); err != nil {
			return uuid.Nil, fmt.Errorf("can't get id: %w", err)
		}
	}

	return id, nil
}

func (s *SqlStorage) MarkAsSent(ctx context.Context, mailID uuid.UUID, time time.Time) error {
	if _, err := s.db.ExecContext(ctx, `
		UPDATE mails SET sent_at = $1 WHERE id = $2
	`, time, mailID); err != nil {
		return fmt.Errorf("can't mark as sent: %w", err)
	}

	return nil
}

func (s *SqlStorage) GetMailById(ctx context.Context, id uuid.UUID) (model.Mail, error) {
	var mail model.Mail

	if err := s.db.GetContext(ctx, &mail, `
		SELECT * FROM mails WHERE id = $1
	`, id); err != nil {
		return model.Mail{}, fmt.Errorf("can't get mail: %w", err)
	}

	return mail, nil
}

func (s *SqlStorage) GetMailsByReceiver(ctx context.Context, receiverID uuid.UUID) ([]model.Mail, error) {
	var mails []model.Mail

	if err := s.db.SelectContext(ctx, &mails, `
		SELECT * FROM mails WHERE to_user_id = $1
	`, receiverID); err != nil {
		return nil, fmt.Errorf("can't get mails by receiver: %w", err)
	}

	return mails, nil
}

func (s *SqlStorage) GetMailWithUser(ctx context.Context, id uuid.UUID) (model.MailWithUser, error) {
	var mail model.MailWithUser

	if err := s.db.GetContext(ctx, &mail, `
		SELECT m.*, u.email, u.first_name, u.last_name FROM mails m
		INNER JOIN users u ON m.to_user_id = u.id
		WHERE m.id = $1
	`, id); err != nil {
		return model.MailWithUser{}, fmt.Errorf("can't get mail with user: %w", err)
	}

	return mail, nil
}
