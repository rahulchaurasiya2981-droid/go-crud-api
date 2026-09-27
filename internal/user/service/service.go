package service

import (
	"context"

	"github.com/rahulchaurasiya2981-droid/go-crud-api/internal/user/dto"
	"github.com/rahulchaurasiya2981-droid/go-crud-api/internal/user/entity"
	"github.com/rahulchaurasiya2981-droid/go-crud-api/internal/user/repository"
)

type Service interface {
	GetUsers(ctx context.Context) ([]dto.UserResponse, error)
	CreateUser(ctx context.Context, request dto.CreateUserRequest) (dto.UserResponse, error)
	DeleteUser(ctx context.Context, id int64) (dto.UserResponse, error)
	UpdateUser(ctx context.Context, id int64, request dto.UpdateUserRequest) (dto.UserResponse, error)
}

type service struct {
	repository repository.Repository
}

func New(repository repository.Repository) Service {
	return &service{repository: repository}
}

func (s *service) GetUsers(ctx context.Context) ([]dto.UserResponse, error) {
	users, err := s.repository.GetUsers(ctx)
	if err != nil {
		return nil, err
	}

	response := make([]dto.UserResponse, len(users))
	for i, user := range users {
		response[i] = toUserResponse(user)
	}

	return response, nil
}

func (s *service) CreateUser(ctx context.Context, request dto.CreateUserRequest) (dto.UserResponse, error) {
	created, err := s.repository.CreateUser(ctx, entity.User{
		Name:  request.Name,
		Email: request.Email,
		Age:   *request.Age,
	})
	if err != nil {
		return dto.UserResponse{}, err
	}

	return toUserResponse(created), nil
}

func (s *service) DeleteUser(ctx context.Context, id int64) (dto.UserResponse, error) {
	deleted, err := s.repository.DeleteUser(ctx, id)
	if err != nil {
		return dto.UserResponse{}, err
	}

	return toUserResponse(deleted), nil
}

func (s *service) UpdateUser(ctx context.Context, id int64, request dto.UpdateUserRequest) (dto.UserResponse, error) {
	updated, err := s.repository.UpdateUser(ctx, id, entity.UserUpdate{
		Name:  request.Name,
		Email: request.Email,
		Age:   request.Age,
	})
	if err != nil {
		return dto.UserResponse{}, err
	}

	return toUserResponse(updated), nil
}

func toUserResponse(user entity.User) dto.UserResponse {
	return dto.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Age:       user.Age,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
