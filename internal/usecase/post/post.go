package post

import (
	"context"
	"go-solid/internal/entity"
	"go-solid/internal/repository"
	"go-solid/internal/usecase"
)

type Usecase struct {
	repo repository.PostRepository
}

func NewUsecase(repo repository.PostRepository) usecase.PostUsecase {
	return &Usecase{repo: repo}
}

func (u *Usecase) SyncPosts(ctx context.Context) error {
	return u.repo.SyncPosts(ctx)
}

func (u *Usecase) GetItems(ctx context.Context) ([]entity.Post, error) {
	return u.repo.GetPosts(ctx)
}
