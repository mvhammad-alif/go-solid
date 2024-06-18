package http

import (
	"go-solid/internal/usecase"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	uc usecase.UserUsecase
}

func NewUserHandler(uc usecase.UserUsecase) *UserHandler {
	return &UserHandler{uc: uc}
}

func (u *UserHandler) GetUserDetail(c echo.Context) (err error) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return
	}

	user, err := u.uc.GetUserDetail(c.Request().Context(), int64(userID))
	if err != nil {
		return
	}

	return c.JSON(http.StatusOK, user)
}
