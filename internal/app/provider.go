package app

import (
	"go-solid/internal/config"
	http "go-solid/internal/delivery"
	postRepo "go-solid/internal/repository/post"
	userRepo "go-solid/internal/repository/user"
	postUC "go-solid/internal/usecase/post"
	userUC "go-solid/internal/usecase/user"

	"github.com/google/wire"
	"github.com/labstack/echo/v4"
)

var (
	configSet = wire.NewSet(
		config.Load,
	)

	handlerSet = wire.NewSet(
		http.NewUserHandler,
		http.NewPostHandler,
	)

	usecaseSet = wire.NewSet(
		userUC.NewUsecase,
		postUC.NewUsecase,
	)

	repositorySet = wire.NewSet(
		userRepo.NewRepository,
		postRepo.NewRepository,
	)

	allSet = wire.NewSet(
		configSet,
		handlerSet,
		usecaseSet,
		repositorySet,

		provideHTTPServer,
	)
)

func provideHTTPServer(postHandler *http.PostHandler) *echo.Echo {
	e := echo.New()
	e.GET("/sync", postHandler.Sync)
	e.GET("/items", postHandler.GetItems)
	return e
}
