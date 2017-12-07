package v1

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
	BasicRequest struct {
		DisplayName string `json:"screen_name"`
		UserID      string `json:"user_id"`
	}
	FollowerResponse struct {
		Ids []bson.ObjectId `json:"ids"`
	}
)

func (h *APIHandler) Follow(c echo.Context) error {
	jwtUser := c.Get("user").(*jwt.Token)
	claims := jwtUser.Claims.(jwt.MapClaims)
	idStr := claims["id"].(string)

	req := new(BasicRequest)
	if err := c.Bind(req); err != nil {
		h.logger.Debug("API Error", zap.String("Error", err.Error()))
		return c.JSON(http.StatusInternalServerError, &messageResponse{Message: ErrUnknown})
	}

	if req.DisplayName != "" {
		f, err := h.db.FindUser(req.DisplayName, true)
		if err != nil {
			return handleMgoError(err)
		}
		err = h.db.FollowUser(bson.ObjectIdHex(idStr), f.ID)
		if err != nil {
			return handleMgoError(err)
		}

		resp := models.UserToUserResponse(*f)

		return c.JSON(http.StatusOK, &resp)
	}

	if req.UserID != "" {
		f, err := h.db.FindUserByOID(bson.ObjectIdHex(req.UserID), true)
		if err != nil {
			return handleMgoError(err)
		}
		err = h.db.FollowUser(bson.ObjectIdHex(idStr), f.ID)
		if err != nil {
			return handleMgoError(err)
		}

		resp := models.UserToUserResponse(*f)

		return c.JSON(http.StatusOK, &resp)
	}

	h.logger.Debug("API Error", zap.String("Error", ErrParamsRequired))
	return &echo.HTTPError{Code: http.StatusBadRequest, Message: ErrParamsRequired}
}

func (h *APIHandler) Unfollow(c echo.Context) error {
	jwtUser := c.Get("user").(*jwt.Token)
	claims := jwtUser.Claims.(jwt.MapClaims)
	idStr := claims["id"].(string)

	req := new(BasicRequest)
	if err := c.Bind(req); err != nil {
		h.logger.Debug("API Error", zap.String("Error", err.Error()))
		return c.JSON(http.StatusInternalServerError, &messageResponse{Message: ErrUnknown})
	}

	if req.DisplayName != "" {
		f, err := h.db.FindUser(req.DisplayName, true)
		if err != nil {
			return handleMgoError(err)
		}
		err = h.db.UnfollowUser(bson.ObjectIdHex(idStr), f.ID)
		if err != nil {
			return handleMgoError(err)
		}

		resp := models.UserToUserResponse(*f)

		return c.JSON(http.StatusOK, &resp)
	}

	if req.UserID != "" {
		f, err := h.db.FindUserByOID(bson.ObjectIdHex(req.UserID), true)
		if err != nil {
			return handleMgoError(err)
		}
		err = h.db.UnfollowUser(bson.ObjectIdHex(idStr), f.ID)
		if err != nil {
			return handleMgoError(err)
		}

		resp := models.UserToUserResponse(*f)

		return c.JSON(http.StatusOK, &resp)
	}

	h.logger.Debug("API Error", zap.String("Error", ErrParamsRequired))
	return &echo.HTTPError{Code: http.StatusBadRequest, Message: ErrParamsRequired}

}

func (h *APIHandler) GetFriendsID(c echo.Context) error {
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

	user, err := h.db.FindUser(id, true)
	if err != nil {
		return handleMgoError(err)
	}

	if len(user.Following) == 0 {
		return c.JSON(http.StatusOK, &FollowerResponse{})
	}

	return c.JSON(http.StatusOK, &FollowerResponse{Ids: user.Following})
}

func (h *APIHandler) GetFollowersID(c echo.Context) error {
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
		user, err = h.db.FindUserByOID(bson.ObjectIdHex(id), true)
		if err != nil {
			return handleMgoError(err)
		}
	}

	if displayName != "" {
		user, err = h.db.FindUser(displayName, true)
		if err != nil {
			return handleMgoError(err)
		}
	}

	if len(user.Followers) == 0 {
		return c.JSON(http.StatusOK, &FollowerResponse{})
	}

	return c.JSON(http.StatusOK, &FollowerResponse{Ids: user.Followers})
}

func (h *APIHandler) GetFollowerList(c echo.Context) error {
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
		user, err = h.db.FindUserByOID(bson.ObjectIdHex(id), true)
		if err != nil {
			return handleMgoError(err)
		}
	}

	if displayName != "" {
		user, err = h.db.FindUser(displayName, true)
		if err != nil {
			return handleMgoError(err)
		}
	}

	if len(user.Followers) == 0 {
		return c.JSON(http.StatusOK, &[]models.UserResponse{})
	}

	users, err := h.db.FindUserByOIDArray(user.Followers, true)
	if err != nil {
		return handleMgoError(err)
	}

	resp := models.UsersToUserResponseArray(users)

	return c.JSON(http.StatusOK, resp)
}

func (h *APIHandler) GetFriendsList(c echo.Context) error {
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
		user, err = h.db.FindUserByOID(bson.ObjectId(id), true)
		if err != nil {
			return handleMgoError(err)
		}
	}

	if displayName != "" {
		user, err = h.db.FindUser(displayName, true)
		if err != nil {
			return handleMgoError(err)
		}
	}

	if len(user.Followers) == 0 {
		return c.JSON(http.StatusOK, &[]models.UserResponse{})
	}

	users, err := h.db.FindUserByOIDArray(user.Following, true)
	if err != nil {
		return handleMgoError(err)
	}

	resp := models.UsersToUserResponseArray(users)

	return c.JSON(http.StatusOK, resp)
}
