package post

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"go-solid/internal/entity"
	"go-solid/internal/repository"
	"go-solid/internal/tools"
	"net/http"
	"time"

	"github.com/cenkalti/backoff/v4"
	_ "github.com/go-sql-driver/mysql"
)

type Repository struct {
	db    *sql.DB
	redis *tools.RedisClient
}

func NewRepository(db *sql.DB, redisClient *tools.RedisClient) repository.PostRepository {
	return &Repository{
		db:    db,
		redis: redisClient,
	}
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

func (r *Repository) GetPostsWithCache(ctx context.Context) ([]entity.Post, error) {
	// Redis cache key
	cacheKey := "posts:all"

	// Try to get from Redis first
	cachedPosts, err := r.redis.Fetch(ctx, cacheKey)
	if err == nil {
		// Cache hit - decode JSON and return
		var posts []entity.Post
		if err := json.Unmarshal([]byte(cachedPosts), &posts); err == nil {
			return posts, nil
		}
		// If JSON decode fails, continue to database
	}

	// Cache miss or decode error - get from database
	posts, err := r.GetPosts(ctx)
	if err != nil {
		return nil, err
	}

	// Store in Redis with 1-minute TTL
	if len(posts) > 0 {
		postsJSON, err := json.Marshal(posts)
		if err == nil {
			r.redis.Store(ctx, cacheKey, string(postsJSON), 1*time.Minute)
		}
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
