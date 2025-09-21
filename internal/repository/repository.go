package repository

import (
	"context"
	"go-solid/internal/entity"
)

type PostRepository interface {
	CreatePost(ctx context.Context, post *entity.Post) error
	GetPosts(ctx context.Context) ([]entity.Post, error)
	GetPostsWithCache(ctx context.Context) ([]entity.Post, error)
	FetchPostsFromAPI(ctx context.Context) ([]entity.Post, error)
}
