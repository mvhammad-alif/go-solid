package user

import (
	"context"
	"go-solid/internal/entity"
	"go-solid/internal/repository"
	"go-solid/internal/usecase"
)

type Usecase struct {
	repo repository.UserRepository
}

func NewUsecase(repo repository.UserRepository) usecase.UserUsecase {
	return &Usecase{repo: repo}
}

func (u Usecase) GetUserDetail(ctx context.Context, userID int64) (*entity.User, error) {
	return u.repo.GetByID(ctx, userID)
}
