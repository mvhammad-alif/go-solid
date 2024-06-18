package usecase

import (
	"context"
	"go-solid/internal/entity"
)

type UserUsecase interface {
	GetUserDetail(ctx context.Context, userID int64) (*entity.User, error)
}
