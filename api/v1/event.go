package v1

import (
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"gopkg.in/mgo.v2/bson"
)

// EventListHandler イベントの一覧を返す
func (h *APIHandler) EventListHandler(c echo.Context) error {
	// Jwtチェック
	jwtUser := c.Get("user").(*jwt.Token)
	claims := jwtUser.Claims.(jwt.MapClaims)
	idStr := claims["id"].(string)
	id := bson.ObjectIdHex(idStr)

	events, err := h.db.GetEvents(id)
	if err != nil {
		return handleMgoError(err)
	}

	return c.JSON(http.StatusOK, &events)
}
