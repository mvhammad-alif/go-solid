package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"go-solid/internal/app"
)

func main() {
	// Initialize the HTTP server
	e, err := app.InitHTTPServer()
	if err != nil {
		panic(err)
	}

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	log.Println("HTTP Server starting on :1323")
	log.Println("Available endpoints:")
	log.Println("  GET /sync - Sync posts from external API")
	log.Println("  GET /items - Get all posts")

	go func() {
		if err := e.Start(":1323"); err != nil {
			log.Printf("Server error: %v", err)
		}
	}()

	// Wait for termination signal
	<-quit
	log.Println("Shutting down server...")
	log.Println("Server stopped")
}
