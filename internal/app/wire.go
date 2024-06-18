//go:build wireinject
// +build wireinject

package app

import (
	"github.com/google/wire"
	"github.com/labstack/echo/v4"
)

func InitHTTPServer() *echo.Echo {
	wire.Build(allSet)
	return &echo.Echo{}
}
