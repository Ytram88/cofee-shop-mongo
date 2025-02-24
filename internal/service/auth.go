package service

import (
	"cofee-shop-mongo/internal/auth"
	"cofee-shop-mongo/internal/config"
	"cofee-shop-mongo/internal/utils"

	"cofee-shop-mongo/models"
	"context"
	"errors"
	"fmt"
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
	// look up the user from database to take the  password
	// get user by email from the database in order to get password
	user, err := s.Repo.GetUserByEmail(ctx, payload.Email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return "", errors.New("invalid password or email")
		}
		return "", err
	}
	// compare passwords from request and database
	if !auth.VerifyPassword(payload.Password, payload.Password) {
		return "", errors.New("invalid password or email")
	}
	//generate jwt token and sign it (auth.CreateJWT already signs it for you)
	tokenString, err := auth.CreateJWT([]byte(s.JWT.JWTSecret), user.UserID, s.JWT.JWTExpirationInSeconds)
	if err != nil {
		return "", err
	}

	// return it back
	return tokenString, nil

}
func (s *AuthService) RegisterUser(ctx context.Context, payload models.RegisterUserPayload) error {
	// check if user already exists
	_, err := s.Repo.GetUserByEmail(ctx, payload.Email)
	if err != nil && err != mongo.ErrNoDocuments {
		return err
	}
	if err == nil {
		return fmt.Errorf("user with this email already exists")
	}
	// hash password
	hashedPassword, err := auth.HashPassword(payload.Password)
	if err != nil {
		return err
	}
	//create a user instance
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
		return errors.New("failed to create user")
	}
	return nil
}
