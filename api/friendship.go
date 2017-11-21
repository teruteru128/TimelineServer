package api

import (
	"net/http"

	"github.com/TinyKitten/TimelineServer/config"
	"github.com/TinyKitten/TimelineServer/models"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"go.uber.org/zap"
	"gopkg.in/mgo.v2/bson"
)

type (
	basicRequest struct {
		DisplayName string `json:"screen_name"`
		UserID      string `json:"user_id"`
	}
	FollowerResponse struct {
		Ids []bson.ObjectId `json:"ids"`
	}
)

func (h *handler) followHandler(c echo.Context) error {
	jwtUser := c.Get("user").(*jwt.Token)
	claims := jwtUser.Claims.(jwt.MapClaims)
	idStr := claims["id"].(string)

	req := new(basicRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusInternalServerError, &messageResponse{Message: ErrUnknown})
	}

	followUser := &models.User{}
	if req.DisplayName != "" {
		f, err := h.db.FindUser(req.DisplayName)
		if err != nil {
			return handleMgoError(err)
		}
		followUser = f
	}

	if req.UserID != "" {
		f, err := h.db.FindUserByOID(bson.ObjectId(req.UserID))
		if err != nil {
			return handleMgoError(err)
		}
		followUser = f
	}

	if req.UserID == "" && req.DisplayName == "" {
		h.logger.Debug("API Error", zap.String("Error", ErrParamsRequired))
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: ErrParamsRequired}
	}

	err := h.db.FollowUser(bson.ObjectId(idStr), followUser.ID)
	if err != nil {
		return handleMgoError(err)
	}

	resp := models.UserToUserResponse(*followUser)

	return c.JSON(http.StatusOK, &resp)
}

func (h *handler) unfollowHandler(c echo.Context) error {
	jwtUser := c.Get("user").(*jwt.Token)
	claims := jwtUser.Claims.(jwt.MapClaims)
	idStr := claims["id"].(string)

	req := new(basicRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusInternalServerError, &messageResponse{Message: ErrUnknown})
	}

	removeUser := &models.User{}
	if req.DisplayName != "" {
		f, err := h.db.FindUser(req.DisplayName)
		if err != nil {
			return handleMgoError(err)
		}
		removeUser = f
	}

	if req.UserID != "" {
		f, err := h.db.FindUserByOID(bson.ObjectId(req.UserID))
		if err != nil {
			return handleMgoError(err)
		}
		removeUser = f
	}

	if req.UserID == "" && req.DisplayName == "" {
		h.logger.Debug("API Error", zap.String("Error", ErrParamsRequired))
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: ErrParamsRequired}
	}

	err := h.db.UnfollowUser(bson.ObjectId(idStr), removeUser.ID)
	if err != nil {
		return handleMgoError(err)
	}

	resp := models.UserToUserResponse(*removeUser)

	return c.JSON(http.StatusOK, &resp)
}

func (h *handler) friendsIdsHandler(c echo.Context) error {
	config := config.GetAPIConfig()
	tokenStr := c.QueryParam("token")
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Jwt), nil
	})
	if err != nil {
		h.logger.Debug("API Error", zap.String("Error", err.Error()))
		return &echo.HTTPError{Code: http.StatusForbidden, Message: ErrInvalidJwt}
	}
	if !token.Valid {
		h.logger.Debug("API Error", zap.String("Error", err.Error()))
		return &echo.HTTPError{Code: http.StatusForbidden, Message: ErrInvalidJwt}
	}
	claims := token.Claims.(jwt.MapClaims)
	id := claims["id"].(string)

	user, err := h.db.FindUser(id)
	if err != nil {
		return handleMgoError(err)
	}

	if len(user.Following) == 0 {
		return c.JSON(http.StatusOK, &FollowerResponse{})
	}

	return c.JSON(http.StatusOK, &FollowerResponse{Ids: user.Following})
}

func (h *handler) followerIdsHandler(c echo.Context) error {
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

	id := c.QueryParam("user_id")
	displayName := c.QueryParam("screen_name")

	user := &models.User{}
	if id != "" {
		user, err = h.db.FindUserByOID(bson.ObjectId(id))
		if err != nil {
			return handleMgoError(err)
		}
	}

	if displayName != "" {
		user, err = h.db.FindUser(displayName)
		if err != nil {
			return handleMgoError(err)
		}
	}

	if len(user.Followers) == 0 {
		return c.JSON(http.StatusOK, &FollowerResponse{})
	}

	return c.JSON(http.StatusOK, &FollowerResponse{Ids: user.Followers})
}

func (h *handler) followerListHandler(c echo.Context) error {
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

	id := c.QueryParam("user_id")
	displayName := c.QueryParam("screen_name")

	user := &models.User{}
	if id != "" {
		user, err = h.db.FindUserByOID(bson.ObjectId(id))
		if err != nil {
			return handleMgoError(err)
		}
	}

	if displayName != "" {
		user, err = h.db.FindUser(displayName)
		if err != nil {
			return handleMgoError(err)
		}
	}

	if len(user.Followers) == 0 {
		return c.JSON(http.StatusOK, &[]models.UserResponse{})
	}

	users, err := h.db.FindUserByOIDArray(user.Followers)
	if err != nil {
		return handleMgoError(err)
	}

	resp := models.UsersToUserResponseArray(users)

	return c.JSON(http.StatusOK, resp)
}

func (h *handler) friendsListHandler(c echo.Context) error {
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

	id := c.QueryParam("user_id")
	displayName := c.QueryParam("screen_name")

	user := &models.User{}
	if id != "" {
		user, err = h.db.FindUserByOID(bson.ObjectId(id))
		if err != nil {
			return handleMgoError(err)
		}
	}

	if displayName != "" {
		user, err = h.db.FindUser(displayName)
		if err != nil {
			return handleMgoError(err)
		}
	}

	if len(user.Followers) == 0 {
		return c.JSON(http.StatusOK, &[]models.UserResponse{})
	}

	users, err := h.db.FindUserByOIDArray(user.Following)
	if err != nil {
		return handleMgoError(err)
	}

	resp := models.UsersToUserResponseArray(users)

	return c.JSON(http.StatusOK, resp)
}
