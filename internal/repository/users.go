package repository

import (
	"cofee-shop-mongo/models"
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type UserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{
		collection: db.Collection("users"),
	}
}

func (r *UserRepository) CreateUser(ctx context.Context, user models.User) (string, error) {
	const op = "repository.CreateUser"
	_, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return user.Username, nil
}

func (r *UserRepository) GetAllUsers(ctx context.Context) ([]models.User, error) {
	const op = "repository.GetAllUsers"
	var users []models.User

	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *UserRepository) GetUserById(ctx context.Context, userId string) (models.User, error) {
	const op = "repository.GetUserById"
	var user models.User

	err := r.collection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.User{}, fmt.Errorf("%s: %w", op, ErrNotFound)
		}
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}
	return user, nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	const op = "repository.GetUserByEmail"
	var user models.User

	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.User{}, fmt.Errorf("%s: %w", op, ErrNotFound)
		}
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}
	return user, nil
}

func (r *UserRepository) UpdateUserById(ctx context.Context, userId string, user models.User) error {
	const op = "repository.UpdateUserById"
	filter := bson.M{"user_id": userId}
	update := bson.M{"$set": bson.M{
		"username": user.Username,
		"email":    user.Email,
		"password": user.Password,
		"role":     user.Role,
	}}

	res, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if res.MatchedCount == 0 {
		return fmt.Errorf("%s: %w", op, ErrNotFound)
	}
	return nil
}

func (r *UserRepository) DeleteUserById(ctx context.Context, userId string) error {
	const op = "repository.DeleteUserById"
	res, err := r.collection.DeleteOne(ctx, bson.M{"user_id": userId})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if res.DeletedCount == 0 {
		return fmt.Errorf("%s: %w", op, ErrNotFound)
	}
	return nil
}
