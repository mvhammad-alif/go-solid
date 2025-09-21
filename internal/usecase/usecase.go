package usecase

import (
	"context"
	"go-solid/internal/entity"
)

type PostUsecase interface {
	SyncPosts(ctx context.Context) error
	GetItems(ctx context.Context) ([]entity.Post, error)
}
