package app

import (
	"go-solid/internal/config"
	"go-solid/internal/delivery"
	postRepo "go-solid/internal/repository/post"
	userRepo "go-solid/internal/repository/user"
	postUC "go-solid/internal/usecase/post"
	userUC "go-solid/internal/usecase/user"
	"go-solid/internal/tools"

	"github.com/google/wire"
	"github.com/labstack/echo/v4"
)

// HTTP Server dependencies
var (
	configSet = wire.NewSet(
		config.Load,
	)

	handlerSet = wire.NewSet(
		delivery.NewUserHandler,
		delivery.NewPostHandler,
	)

	usecaseSet = wire.NewSet(
		userUC.NewUsecase,
		postUC.NewUsecase,
	)

	repositorySet = wire.NewSet(
		userRepo.NewRepository,
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

func provideCronJobs(cronHandler *delivery.CronHandler) []tools.CronJob {
	return []tools.CronJob{
		{
			Name:     "sync_posts",
			Schedule: "*/15 * * * *", // Every 15 minutes
			Job:      cronHandler.SyncPostsJob,
		},
	}
}
