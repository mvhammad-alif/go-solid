package post

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"go-solid/internal/config"
	"go-solid/internal/database"
	"go-solid/internal/entity"
	"go-solid/internal/repository"
	"net/http"
	"time"

	"github.com/cenkalti/backoff/v4"
	_ "github.com/go-sql-driver/mysql"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(cfg *config.Config) (repository.PostRepository, error) {
	db, err := sql.Open("mysql", cfg.GetDatabaseDSN())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Run migrations
	migration := database.NewMigration(db)
	if err := migration.CreateTables(); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return &Repository{db: db}, nil
}

func (r *Repository) GetPosts(ctx context.Context) ([]entity.Post, error) {
	query := "SELECT id, user_id, title, body FROM posts ORDER BY created_at DESC"

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query posts: %w", err)
	}
	defer rows.Close()

	var posts []entity.Post
	for rows.Next() {
		var post entity.Post
		err := rows.Scan(&post.ID, &post.UserID, &post.Title, &post.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to scan post: %w", err)
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return posts, nil
}

func (r *Repository) CreatePost(ctx context.Context, post *entity.Post) error {
	query := "INSERT INTO posts (id, user_id, title, body) VALUES (?, ?, ?, ?) ON DUPLICATE KEY UPDATE title = ?, body = ?"

	_, err := r.db.ExecContext(ctx, query,
		post.ID, post.UserID, post.Title, post.Body,
		post.Title, post.Body,
	)

	if err != nil {
		return fmt.Errorf("failed to create post: %w", err)
	}

	return nil
}

func (r *Repository) TagPost(ctx context.Context, userIDs []int64) error {
	// For now, this is a placeholder implementation
	// In a real application, you might want to create a separate table for post tags
	// or use a many-to-many relationship table
	time.Sleep(1 * time.Second)
	return nil
}

func (r *Repository) FetchPostsFromAPI(ctx context.Context) ([]entity.Post, error) {
	var posts []entity.Post

	// Create exponential backoff strategy
	backoffStrategy := backoff.NewExponentialBackOff()
	backoffStrategy.InitialInterval = 1 * time.Second
	backoffStrategy.MaxInterval = 30 * time.Second
	backoffStrategy.MaxElapsedTime = 5 * time.Minute
	backoffStrategy.Multiplier = 2

	// Operation to retry - fetching posts from external API
	operation := func() error {
		resp, err := http.Get("https://jsonplaceholder.typicode.com/posts")
		if err != nil {
			return backoff.Permanent(fmt.Errorf("failed to make HTTP request: %w", err))
		}
		defer resp.Body.Close()

		// Don't retry on client errors (4xx) - these are permanent
		if resp.StatusCode >= 400 && resp.StatusCode < 500 {
			return backoff.Permanent(fmt.Errorf("client error: %d %s", resp.StatusCode, resp.Status))
		}

		// Don't retry on server errors (5xx) - these might be temporary, so we allow retry
		if resp.StatusCode >= 500 {
			return fmt.Errorf("server error: %d %s", resp.StatusCode, resp.Status)
		}

		// Only process successful responses (2xx)
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			return backoff.Permanent(fmt.Errorf("unexpected status code: %d %s", resp.StatusCode, resp.Status))
		}

		return json.NewDecoder(resp.Body).Decode(&posts)
	}

	// Execute with retry logic
	if err := backoff.Retry(operation, backoff.WithContext(backoffStrategy, ctx)); err != nil {
		return nil, fmt.Errorf("failed to fetch posts after retries: %w", err)
	}

	return posts, nil
}
