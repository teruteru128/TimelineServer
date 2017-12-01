package v1

import (
	"net/http"
	"strconv"

	"github.com/TinyKitten/TimelineServer/models"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"go.uber.org/zap"
	"gopkg.in/mgo.v2/bson"
)

type ABasicRequest struct {
	UserID bson.ObjectId `json:"user_id"`
}

// 管理者API　ObjectIDで処理
func (h *APIHandler) AUserSuspendHandler(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	// 管理者チェック
	if claims["admin"].(bool) == false {
		return &echo.HTTPError{Code: http.StatusForbidden, Message: ErrAdminOnly}
	}
	req := new(ABasicRequest)
	if err := c.Bind(req); err != nil {
		h.logger.Debug("API Error", zap.String("Error", ErrParamsRequired))
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: ErrParamsRequired}
	}

	u, err := h.db.FindUserByOID(req.UserID, true)
	if err != nil {
		h.logger.Debug("API Error", zap.String("Error", ErrParamsRequired))
		return handleMgoError(err)
	}

	err = h.db.UpdateUser(req.UserID, "suspended", true)
	if err != nil {
		h.logger.Debug("API Error", zap.String("Error", ErrParamsRequired))
		return handleMgoError(err)
	}
	resp := models.UserToUserResponse(*u)

	return c.JSON(http.StatusOK, resp)
}

func (h *APIHandler) ASetOfficialFlag(c echo.Context) error {
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
