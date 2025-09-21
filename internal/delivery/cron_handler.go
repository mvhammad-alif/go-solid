package delivery

import (
	"context"
	"log"
	"time"

	"go-solid/internal/usecase"
)

type CronHandler struct {
	postUsecase usecase.PostUsecase
}

func NewCronHandler(postUsecase usecase.PostUsecase) *CronHandler {
	return &CronHandler{
		postUsecase: postUsecase,
	}
}

func (c *CronHandler) SyncPostsJob() {
	log.Println("Running scheduled sync posts job...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	start := time.Now()

	err := c.postUsecase.SyncPosts(ctx)
	duration := time.Since(start)

	if err != nil {
		log.Printf("ERROR: Sync posts job failed after %v: %v", duration, err)
	} else {
		log.Printf("SUCCESS: Sync posts job completed in %v", duration)
	}
}
