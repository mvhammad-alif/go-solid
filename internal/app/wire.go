//go:build wireinject
// +build wireinject

package app

import (
	"github.com/google/wire"
	"github.com/labstack/echo/v4"
	"go-solid/internal/tools"
)

func InitHTTPServer() (*echo.Echo, error) {
	wire.Build(httpServerSet)
	return &echo.Echo{}, nil
}

func InitCronService() (*tools.CronService, error) {
	wire.Build(cronServiceSet)
	return &tools.CronService{}, nil
}
