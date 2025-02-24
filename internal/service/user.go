package service

import (
	"cofee-shop-mongo/models"
	"context"
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
	return s.Repo.CreateUser(ctx, user)
}

func (s *UserService) GetAllUsers(ctx context.Context) ([]models.User, error) {
	return s.Repo.GetAllUsers(ctx)
}

func (s *UserService) GetUserById(ctx context.Context, userId string) (models.User, error) {
	return s.Repo.GetUserById(ctx, userId)
}

func (s *UserService) UpdateUserById(ctx context.Context, userId string, user models.User) error {
	return s.Repo.UpdateUserById(ctx, userId, user)
}

func (s *UserService) DeleteUserById(ctx context.Context, userId string) error {
	return s.Repo.DeleteUserById(ctx, userId)
}
