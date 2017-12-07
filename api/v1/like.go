package v1

import (
	"net/http"

	"github.com/TinyKitten/TimelineServer/models"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"go.uber.org/zap"
	"gopkg.in/mgo.v2/bson"
)

type (
	LikeRequest struct {
		PostID string `json:"id" validate:"required"`
	}
)

func (h *APIHandler) CreateLike(c echo.Context) error {
	jwtUser := c.Get("user").(*jwt.Token)
	claims := jwtUser.Claims.(jwt.MapClaims)
	idStr := claims["id"].(string)

	req := new(LikeRequest)
	if err := c.Bind(req); err != nil {
		h.logger.Debug("API Error", zap.String("Error", err.Error()))
		return c.JSON(http.StatusInternalServerError, &messageResponse{Message: ErrUnknown})
	}

	if bson.IsObjectIdHex(req.PostID) {
		err := h.db.CreateLike(bson.ObjectIdHex(req.PostID), bson.ObjectIdHex(idStr))
		if err != nil {
			return handleMgoError(err)
		}

		updated, err := h.db.FindPost(bson.ObjectIdHex(req.PostID), true)
		if err != nil {
			return handleMgoError(err)
		}
		sender, err := h.db.FindUserByOID(updated.UserID, true)
		if err != nil {
			return handleMgoError(err)
		}

		resp := models.PostToPostResponse(*updated, *sender)

		return c.JSON(http.StatusOK, &resp)
	}

	h.logger.Debug("API Error", zap.String("Error", ErrParamsRequired))
	return &echo.HTTPError{Code: http.StatusBadRequest, Message: ErrParamsRequired}
}

func (h *APIHandler) DestroyLike(c echo.Context) error {
	jwtUser := c.Get("user").(*jwt.Token)
	claims := jwtUser.Claims.(jwt.MapClaims)
	idStr := claims["id"].(string)

	req := new(LikeRequest)
	if err := c.Bind(req); err != nil {
		h.logger.Debug("API Error", zap.String("Error", err.Error()))
		return c.JSON(http.StatusInternalServerError, &messageResponse{Message: ErrUnknown})
	}

	if bson.IsObjectIdHex(req.PostID) {
		err := h.db.DestroyLike(bson.ObjectIdHex(req.PostID), bson.ObjectIdHex(idStr))
		if err != nil {
			return handleMgoError(err)
		}

		updated, err := h.db.FindPost(bson.ObjectIdHex(req.PostID), true)
		if err != nil {
			return handleMgoError(err)
		}
		sender, err := h.db.FindUserByOID(updated.UserID, true)
		if err != nil {
			return handleMgoError(err)
		}

		resp := models.PostToPostResponse(*updated, *sender)

		return c.JSON(http.StatusOK, &resp)
	}

	h.logger.Debug("API Error", zap.String("Error", ErrParamsRequired))
	return &echo.HTTPError{Code: http.StatusBadRequest, Message: ErrParamsRequired}
}
