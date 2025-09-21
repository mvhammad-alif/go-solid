package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"go-solid/internal/app"
)

func main() {
	// Initialize only the cron service
	cronService, err := app.InitCronService()
	if err != nil {
		panic(err)
	}

	// Start the cron service
	if err := cronService.Start(); err != nil {
		log.Printf("Failed to start cron service: %v", err)
		panic(err)
	}

	log.Println("Cron service starting...")
	log.Println("Scheduled jobs:")
	log.Println("  - Sync posts: every 15 minutes")

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	log.Println("Cron service is running. Press Ctrl+C to stop.")

	// Wait for termination signal
	<-quit
	log.Println("Shutting down cron service...")

	// Stop the cron service
	cronService.Stop()

	log.Println("Cron service stopped")
}
