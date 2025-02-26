package service

import (
	"cofee-shop-mongo/internal/auth"
	"cofee-shop-mongo/internal/config"
	"cofee-shop-mongo/internal/repository"
	"cofee-shop-mongo/internal/utils"
	"errors"
	"fmt"

	"cofee-shop-mongo/models"
	"context"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type AuthRepository interface {
	GetUserByEmail(ctx context.Context, email string) (models.User, error)
	CreateUser(ctx context.Context, user models.User) (string, error)
}

type AuthService struct {
	Repo AuthRepository
	JWT  config.JWTConfig
}

func NewAuthService(Repo AuthRepository, JWTConfig config.JWTConfig) *AuthService {
	return &AuthService{Repo, JWTConfig}
}

func (s *AuthService) LoginUser(ctx context.Context, payload models.UserLoginPayload) (string, error) {
	const op = "service.LoginUser"
	// look up the user from database to take the  password
	// get user by email from the database in order to get password
	user, err := s.Repo.GetUserByEmail(ctx, payload.Email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return "", errors.New("invalid email or password")
		}
		return "", fmt.Errorf("%s: failed to fetch user: %w", op, err)
	}

	// compare passwords from request and database
	if !auth.VerifyPassword(user.Password, payload.Password) {
		return "", errors.New("invalid email or password")
	}

	//generate jwt token and sign it (auth.CreateJWT already signs it for you)
	tokenString, err := auth.CreateJWT(user.UserID, user.Role, s.JWT.JWTExpirationInSeconds)
	if err != nil {
		return "", fmt.Errorf("%s: failed to generate JWT token: %w", op, err)
	}

	return tokenString, nil
}

func (s *AuthService) RegisterUser(ctx context.Context, payload models.RegisterUserPayload) (string, error) {
	// check if user already exists
	const op = "service.RegisterUser"
	_, err := s.Repo.GetUserByEmail(ctx, payload.Email)
	if err == nil {
		return "", ErrAlreadyExists
	} else if !errors.Is(err, repository.ErrNotFound) {
		fmt.Println("sdaf")
		return "", fmt.Errorf("%s: failed to check existing user: %w", op, err)
	}

	// hash password
	hashedPassword, err := auth.HashPassword(payload.Password)
	if err != nil {
		return "", fmt.Errorf("%s: failed to hash password: %w", op, err)
	}

	// create user instance
	user := models.User{
		UserID:   utils.GenerateRandomString(5),
		Username: payload.Username,
		Email:    payload.Email,
		Password: hashedPassword,
		Role:     "client",
	}

	// call create register method to register new user
	_, err = s.Repo.CreateUser(ctx, user)
	if err != nil {
		return "", fmt.Errorf("%s: failed to create user: %w", op, err)
	}

	return user.UserID, nil
}
