package delivery

import (
	"go-solid/internal/usecase"
	"net/http"

	"github.com/labstack/echo/v4"
)

type PostHandler struct {
	uc usecase.PostUsecase
}

func NewPostHandler(uc usecase.PostUsecase) *PostHandler {
	return &PostHandler{uc: uc}
}

func (u *PostHandler) Sync(c echo.Context) (err error) {
	err = u.uc.SyncPosts(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Posts synced successfully"})
}

func (u *PostHandler) GetItems(c echo.Context) (err error) {
	posts, err := u.uc.GetItems(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, posts)
}
