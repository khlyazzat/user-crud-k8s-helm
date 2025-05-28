package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/khlyazzat/user-crud-k8s-helm/internal/db/models"
	"github.com/khlyazzat/user-crud-k8s-helm/internal/db/repository"
	"github.com/khlyazzat/user-crud-k8s-helm/internal/values"

	db "github.com/khlyazzat/user-crud-k8s-helm/internal/db/postgres"

	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
)

type User interface {
	AddUser(ctx context.Context, user *models.User) (string, error)
	GetUserByID(ctx context.Context, userId string) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) (*models.User, error)
	DeleteUser(ctx context.Context, user *models.User) error
}

type userRepository struct {
	repository.CRUD
	db.DB
}

func (r *userRepository) AddUser(ctx context.Context, user *models.User) (string, error) {
	var id string
	err := r.DB.NewInsert().
		Model(user).
		Returning("id").
		Scan(ctx, &id)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return "", fmt.Errorf("%w: %s", values.ErrEmailExists, pgErr.Detail)
		}
		return "", fmt.Errorf("failed to insert user: %w", err)
	}
	return id, nil
}

func (r *userRepository) GetUser(ctx context.Context, userId uuid.UUID) (*models.User, error) {
	m := &models.User{
		ID: userId,
	}
	err := r.GetByID(ctx, m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (r *userRepository) GetUserByID(ctx context.Context, userId string) (*models.User, error) {
	user := new(models.User)

	err := r.DB.NewSelect().
		Model(user).
		Where("id = ?", userId).
		Limit(1).
		Scan(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("%w: id=%s", values.ErrUserNotFound, userId)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}
	return user, nil
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	user := new(models.User)

	err := r.DB.NewSelect().
		Model(user).
		Where("email = ?", email).
		Limit(1).
		Scan(ctx)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	return user, nil
}

func (r *userRepository) UpdateUser(ctx context.Context, user *models.User) (*models.User, error) {
	err := r.DB.NewUpdate().
		Model(user).
		WherePK().
		Returning("*").
		Scan(ctx)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("%w: id=%s", values.ErrUserNotFound, user.ID)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}
	return user, nil
}

func (r *userRepository) DeleteUser(ctx context.Context, user *models.User) error {
	res, err := r.DB.NewDelete().
		Model(user).
		WherePK().
		Exec(ctx)

	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("%w: id=%s", values.ErrUserNotFound, user.ID)
	}
	return nil
}

func New(db db.DB) User {
	return &userRepository{
		CRUD: repository.NewCRUD(db),
		DB:   db,
	}
}
