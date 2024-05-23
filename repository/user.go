package repository

import (
	"context"

	"github.com/f3rcho/rest-posts/models"
)

type UserRepository interface {
	InsertUser(ctx context.Context, user *models.User) error
	GetUserByID(ctx context.Context, ID int64) (*models.User, error)
}

var implementation UserRepository

func SetRespository(repository UserRepository) {
	implementation = repository
}

func InsertUser(ctx context.Context, user *models.User) error {
	return implementation.InsertUser(ctx, user)
}

func GetUserByID(ctx context.Context, ID int64) (*models.User, error) {
	return implementation.GetUserByID(ctx, ID)
}
