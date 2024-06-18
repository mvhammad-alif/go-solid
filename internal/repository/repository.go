package repository

import (
	"context"
	"go-solid/internal/entity"
)

type UserRepository interface {
	GetByID(ctx context.Context, userID int64) (*entity.User, error)
}
