package repository

import (
	"context"
	"go-solid/internal/entity"
)

type UserRepository interface {
	GetByID(ctx context.Context, userID int64) (*entity.User, error)
}

type PostRepository interface {
	CreatePost(ctx context.Context, post *entity.Post) error
	GetPosts(ctx context.Context) ([]entity.Post, error)
	TagPost(ctx context.Context, userIDs []int64) error
	SyncPosts(ctx context.Context) error
}
