// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package app

import (
	"github.com/labstack/echo/v4"
	"go-solid/internal/delivery/http"
	"go-solid/internal/repository/user"
	user2 "go-solid/internal/usecase/user"
)

// Injectors from wire.go:

func InitHTTPServer() *echo.Echo {
	userRepository := user.NewRepository()
	userUsecase := user2.NewUsecase(userRepository)
	userHandler := http.NewUserHandler(userUsecase)
	echoEcho := provideHTTPServer(userHandler)
	return echoEcho
}