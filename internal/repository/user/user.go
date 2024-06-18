package user

import (
	"context"
	"errors"
	"go-solid/internal/entity"
	"go-solid/internal/repository"
)

type Repository struct{}

func NewRepository() repository.UserRepository {
	return &Repository{}
}

var mapUser = map[int64]*entity.User{
	1: {
		ID:    1,
		Name:  "John Doe",
		Email: "john.doe@gmail.com",
	},
	2: {
		ID:    2,
		Name:  "John Doe 2",
		Email: "john.doe+2@gmail.com",
	},
}

func (r Repository) GetByID(ctx context.Context, userID int64) (*entity.User, error) {
	if user, ok := mapUser[userID]; ok {
		return user, nil
	}

	return nil, errors.New("user not found")
}
