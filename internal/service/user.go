package service

import (
	"cofee-shop-mongo/models"
	"context"
	"fmt"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user models.User) (string, error)
	GetAllUsers(ctx context.Context) ([]models.User, error)
	GetUserById(ctx context.Context, userId string) (models.User, error)
	UpdateUserById(ctx context.Context, userId string, user models.User) error
	DeleteUserById(ctx context.Context, userId string) error
}

type UserService struct {
	Repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{Repo: repo}
}

func (s *UserService) CreateUser(ctx context.Context, user models.User) (string, error) {
	const op = "service.CreateUser"

	id, err := s.Repo.CreateUser(ctx, user)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *UserService) GetAllUsers(ctx context.Context) ([]models.User, error) {
	const op = "service.GetAllUsers"

	users, err := s.Repo.GetAllUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return users, nil
}

func (s *UserService) GetUserById(ctx context.Context, userId string) (models.User, error) {
	const op = "service.GetUserById"

	user, err := s.Repo.GetUserById(ctx, userId)
	if err != nil {
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (s *UserService) UpdateUserById(ctx context.Context, userId string, user models.User) error {
	const op = "service.UpdateUserById"
	
	err := s.Repo.UpdateUserById(ctx, userId, user)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *UserService) DeleteUserById(ctx context.Context, userId string) error {
	const op = "service.DeleteUserById"

	err := s.Repo.DeleteUserById(ctx, userId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
