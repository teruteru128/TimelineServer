package v1

import (
	"net/http"

	"github.com/TinyKitten/TimelineServer/config"
	"github.com/TinyKitten/TimelineServer/models"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

func (h *APIHandler) SearchUserHandler(c echo.Context) error {
	// Jwtチェック
	config := config.GetAPIConfig()
	tokenStr := c.QueryParam("token")
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Jwt), nil
	})
	if err != nil {
		return &echo.HTTPError{Code: http.StatusForbidden, Message: ErrInvalidJwt}
	}
	if !token.Valid {
		return &echo.HTTPError{Code: http.StatusForbidden, Message: ErrInvalidJwt}
	}

	query := c.QueryParam("query")

	users, err := h.db.SearchUser(query, 5)
	if err != nil {
		return handleMgoError(err)
	}

	if len(*users) == 0 {
		return c.JSON(http.StatusOK, &[]models.User{})
	}

	resp := models.UsersToUserResponseArray(*users)

	return c.JSON(http.StatusOK, resp)
}
