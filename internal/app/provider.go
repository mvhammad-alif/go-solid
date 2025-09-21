package app

import (
	"context"
	"database/sql"
	"fmt"

	"go-solid/internal/config"
	"go-solid/internal/delivery"
	postRepo "go-solid/internal/repository/post"
	"go-solid/internal/tools"
	postUC "go-solid/internal/usecase/post"

	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/wire"
	"github.com/labstack/echo/v4"
)

// HTTP Server dependencies
var (
	configSet = wire.NewSet(
		config.Load,
	)

	handlerSet = wire.NewSet(
		delivery.NewPostHandler,
	)

	usecaseSet = wire.NewSet(
		postUC.NewUsecase,
	)

	repositorySet = wire.NewSet(
		provideDatabase,
		provideRedisClient,
		postRepo.NewRepository,
	)

	httpServerSet = wire.NewSet(
		configSet,
		handlerSet,
		usecaseSet,
		repositorySet,
		provideHTTPServer,
	)
)

// Cron service dependencies
var (
	cronConfigSet = wire.NewSet(
		config.Load,
	)

	cronUsecaseSet = wire.NewSet(
		postUC.NewUsecase,
	)

	cronRepositorySet = wire.NewSet(
		postRepo.NewRepository,
	)

	cronHandlerSet = wire.NewSet(
		delivery.NewCronHandler,
	)

	cronToolSet = wire.NewSet(
		tools.NewCronService,
	)

	cronServiceSet = wire.NewSet(
		cronConfigSet,
		cronUsecaseSet,
		provideDatabase,
		provideRedisClient,
		cronRepositorySet,
		cronHandlerSet,
		cronToolSet,
		provideCronJobs,
	)
)

func provideHTTPServer(postHandler *delivery.PostHandler) *echo.Echo {
	e := echo.New()
	e.GET("/sync", postHandler.Sync)
	e.GET("/items", postHandler.GetItems)
	return e
}

func provideDatabase(cfg *config.Config) (*sql.DB, error) {
	db, err := sql.Open("mysql", cfg.GetDatabaseDSN())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

func provideRedisClient(cfg *config.Config) (*tools.RedisClient, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.GetRedisAddr(),
	})

	// Test Redis connection
	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return tools.NewRedisClient(redisClient), nil
}

func provideCronJobs(cronHandler *delivery.CronHandler) []tools.CronJob {
	return []tools.CronJob{
		{
			Name:     "sync_posts",
			Schedule: "*/15 * * * *", // Every 15 minutes
			Job:      cronHandler.SyncPostsJob,
		},
	}
}
