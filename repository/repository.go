package repository

import (
	"context"

	"github.com/f3rcho/rest-posts/models"
)

type Repository interface {
	InsertUser(ctx context.Context, user *models.User) error
	GetUserByID(ctx context.Context, ID string) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	InsertPost(ctx context.Context, post *models.Post) error
	GetPostById(ctx context.Context, ID string) (*models.Post, error)
	DeletePostById(ctx context.Context, ID string, userID string) error
	UpdatePost(ctx context.Context, post *models.Post, userId string) error
	ListPosts(ctx context.Context, pagination *models.PaginationDTO) ([]*models.Post, error)
	Close() error
}

var implementation Repository

func SetRespository(repository Repository) {
	implementation = repository
}

// user
func InsertUser(ctx context.Context, user *models.User) error {
	return implementation.InsertUser(ctx, user)
}
func GetUserByID(ctx context.Context, ID string) (*models.User, error) {
	return implementation.GetUserByID(ctx, ID)
}
func GetUserByEmail(ctx context.Context, Email string) (*models.User, error) {
	return implementation.GetUserByEmail(ctx, Email)
}

// post
func InsertPost(ctx context.Context, post *models.Post) error {
	return implementation.InsertPost(ctx, post)
}
func ListPosts(ctx context.Context, pagination *models.PaginationDTO) ([]*models.Post, error) {
	return implementation.ListPosts(ctx, pagination)
}
func GetPostById(ctx context.Context, ID string) (*models.Post, error) {
	return implementation.GetPostById(ctx, ID)
}
func DeletePostById(ctx context.Context, ID, userID string) error {
	return implementation.DeletePostById(ctx, ID, userID)
}
func UpdatePost(ctx context.Context, post *models.Post, userId string) error {
	return implementation.UpdatePost(ctx, post, userId)
}

func Close() error {
	return implementation.Close()
}
