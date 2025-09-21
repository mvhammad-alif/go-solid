package tools

import (
	"fmt"
	"log"
	"time"

	"github.com/robfig/cron/v3"
)

type CronService struct {
	cron *cron.Cron
	jobs []CronJob
}

type CronJob struct {
	Name     string
	Schedule string
	Job      func()
}

func NewCronService(jobs []CronJob) *CronService {
	return &CronService{
		cron: cron.New(cron.WithLocation(time.UTC)),
		jobs: jobs,
	}
}

func (c *CronService) Start() error {
	for _, job := range c.jobs {
		_, err := c.cron.AddFunc(job.Schedule, job.Job)
		if err != nil {
			return fmt.Errorf("failed to add cron job '%s': %w", job.Name, err)
		}
		log.Printf("Scheduled cron job: %s (%s)", job.Name, job.Schedule)
	}

	log.Printf("Starting cron service with %d jobs...", len(c.jobs))
	c.cron.Start()
	return nil
}

func (c *CronService) Stop() {
	log.Println("Stopping cron service...")
	c.cron.Stop()
}
