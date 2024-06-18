package app

import (
	"go-solid/internal/delivery/http"
	userRepo "go-solid/internal/repository/user"
	userUC "go-solid/internal/usecase/user"

	"github.com/google/wire"
	"github.com/labstack/echo/v4"
)

var (
	handlerSet = wire.NewSet(
		http.NewUserHandler,
	)

	usecaseSet = wire.NewSet(
		userUC.NewUsecase,
	)

	repositorySet = wire.NewSet(
		userRepo.NewRepository,
	)

	allSet = wire.NewSet(
		handlerSet,
		usecaseSet,
		repositorySet,

		provideHTTPServer,
	)
)

func provideHTTPServer(userHandler *http.UserHandler) *echo.Echo {
	e := echo.New()
	e.GET("/users/:id", userHandler.GetUserDetail)
	return e
}
