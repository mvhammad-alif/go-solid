package tools

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/robfig/cron/v3"
	"go-solid/internal/usecase"
)

type CronService struct {
	cron    *cron.Cron
	usecase usecase.PostUsecase
}

func NewCronService(postUsecase usecase.PostUsecase) *CronService {
	return &CronService{
		cron:    cron.New(cron.WithLocation(time.UTC)),
		usecase: postUsecase,
	}
}

func (c *CronService) Start() error {
	// Add sync posts job every 15 minutes
	_, err := c.cron.AddFunc("*/15 * * * *", func() {
		c.syncPostsJob()
	})
	if err != nil {
		return fmt.Errorf("failed to add sync posts job: %w", err)
	}

	log.Println("Starting cron service...")
	log.Println("Scheduled jobs:")
	log.Println("  - Sync posts: every 15 minutes")

	c.cron.Start()
	return nil
}

func (c *CronService) Stop() {
	log.Println("Stopping cron service...")
	c.cron.Stop()
}

func (c *CronService) syncPostsJob() {
	log.Println("Running scheduled sync posts job...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	start := time.Now()

	err := c.usecase.SyncPosts(ctx)
	duration := time.Since(start)

	if err != nil {
		log.Printf("ERROR: Sync posts job failed after %v: %v", duration, err)
	} else {
		log.Printf("SUCCESS: Sync posts job completed in %v", duration)
	}
}
