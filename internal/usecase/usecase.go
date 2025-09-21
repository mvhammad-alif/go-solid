package usecase

import (
	"context"
	"go-solid/internal/entity"
)

type UserUsecase interface {
	GetUserDetail(ctx context.Context, userID int64) (*entity.User, error)
}

type PostUsecase interface {
	SyncPosts(ctx context.Context) error
	GetItems(ctx context.Context) ([]entity.Post, error)
}
