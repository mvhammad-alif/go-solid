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
	// Fetch posts from external API
	posts, err := u.repo.FetchPostsFromAPI(ctx)
	if err != nil {
		return err
	}

	// Process and store the fetched posts
	for _, post := range posts {
		if err := u.repo.CreatePost(ctx, &post); err != nil {
			// Continue with other posts even if one fails
			// Log the error but don't fail the entire sync
			continue
		}
	}

	return nil
}

func (u *Usecase) GetItems(ctx context.Context) ([]entity.Post, error) {
	return u.repo.GetPostsWithCache(ctx)
}
