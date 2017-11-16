package api

import (
	"net/http"
	"strconv"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"gopkg.in/mgo.v2/bson"
)

// 管理者API　ObjectIDで処理
func (h *handler) userSuspendHandler(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	// 管理者チェック
	if claims["admin"].(bool) == false {
		return &echo.HTTPError{Code: http.StatusForbidden, Message: ErrAdminOnly}
	}

	oid := bson.ObjectIdHex(c.QueryParam("oid"))

	if err := h.db.SuspendUser(oid, true); err != nil {
		return handleMgoError(err)
	}

	return c.JSON(http.StatusOK, &messageResponse{Message: "ok"})
}

func (h *handler) setOfficalFlagHandler(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	// 管理者チェック
	if claims["admin"].(bool) == false {
		return &echo.HTTPError{Code: http.StatusForbidden, Message: ErrAdminOnly}
	}

	oid := bson.ObjectIdHex(c.QueryParam("oid"))
	flag, err := strconv.ParseBool(c.QueryParam("flag"))
	if err != nil {
		return &echo.HTTPError{Code: http.StatusForbidden, Message: ErrBadFormat}
	}

	if err := h.db.SetOfficial(oid, flag); err != nil {
		return handleMgoError(err)
	}

	return c.JSON(http.StatusOK, &messageResponse{Message: "ok"})
}
