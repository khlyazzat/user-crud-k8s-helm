package user

import (
	"context"
	"database/sql"
	"errors"

	"github.com/khlyazzat/user-crud-k8s-helm/internal/db/models"
	"github.com/khlyazzat/user-crud-k8s-helm/internal/db/postgres"
	userRepo "github.com/khlyazzat/user-crud-k8s-helm/internal/db/repository/user"
	"github.com/khlyazzat/user-crud-k8s-helm/internal/values"

	"github.com/khlyazzat/user-crud-k8s-helm/internal/dto"
)

type User interface {
	AddUser(ctx context.Context, request *dto.AddUserRequest) (*dto.AddUserResponse, error)
	GetUser(ctx context.Context, request *dto.GetUserRequest) (*dto.GetUserResponse, error)
	DeleteUser(ctx context.Context, request *dto.DeleteUserRequest) error
	UpdateUser(ctx context.Context, userId string, request *dto.UpdateUserRequest) (*dto.UpdateUserResponse, error)
}

type userService struct {
	userRepo userRepo.User
}

func (u *userService) AddUser(ctx context.Context, request *dto.AddUserRequest) (*dto.AddUserResponse, error) {
	user, err := u.userRepo.GetUserByEmail(ctx, request.Email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	if user != nil {
		return nil, values.ErrEmailExists
	}

	userId := ""
	newUser := &models.User{
		Name:  request.Name,
		Email: request.Email,
		Age:   request.Age,
	}

	userId, err = u.userRepo.AddUser(ctx, newUser)
	if err != nil {
		return nil, err
	}
	return &dto.AddUserResponse{
		ID: userId,
	}, nil
}

func (u *userService) GetUser(ctx context.Context, request *dto.GetUserRequest) (*dto.GetUserResponse, error) {
	user, err := u.userRepo.GetUserByID(ctx, request.UserID)
	if err != nil {
		return nil, err
	}
	return &dto.GetUserResponse{
		ID:    user.ID.String(),
		Name:  user.Name,
		Email: user.Email,
		Age:   user.Age,
	}, nil
}

func (u *userService) DeleteUser(ctx context.Context, request *dto.DeleteUserRequest) error {
	user, err := u.userRepo.GetUserByID(ctx, request.UserID)
	if err != nil {
		return err
	}
	err = u.userRepo.DeleteUser(ctx, user)
	if err != nil {
		return err
	}
	return nil

}

func (u *userService) UpdateUser(ctx context.Context, userId string, request *dto.UpdateUserRequest) (*dto.UpdateUserResponse, error) {
	user, err := u.userRepo.GetUserByID(ctx, userId)
	if err != nil {
		return nil, err
	}
	if request.Name != nil {
		user.Name = *request.Name
	}
	if request.Email != nil {
		user.Email = *request.Email
	}
	if request.Age != nil {
		user.Age = *request.Age
	}
	user, err = u.userRepo.UpdateUser(ctx, user)
	if err != nil {
		return nil, err
	}
	return &dto.UpdateUserResponse{
		ID:    user.ID.String(),
		Name:  user.Name,
		Email: user.Email,
		Age:   user.Age,
	}, nil
}

func NewUserService(db postgres.DB) User {
	return &userService{
		userRepo: userRepo.New(db),
	}
}
