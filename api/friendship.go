package api

import (
	"net/http"

	"github.com/TinyKitten/TimelineServer/config"
	"github.com/TinyKitten/TimelineServer/models"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"gopkg.in/mgo.v2/bson"
)

type (
	usersResponse struct {
		Users []models.UserResponse `json:"users"`
	}
	basicRequest struct {
		DisplayName string `json:"displayName"`
		UserID      string `json:"userId"`
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

	followUser, err := h.db.FindUser(req.DisplayName)
	if err != nil {
		return handleMgoError(err)
	}

	err = h.db.FollowUser(bson.ObjectId(idStr), followUser.ID)
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
	objID := bson.ObjectId(bson.ObjectIdHex(idStr))

	req := new(basicRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusInternalServerError, &messageResponse{Message: ErrUnknown})
	}

	removeUser, err := h.db.FindUser(req.DisplayName)
	if err != nil {
		return handleMgoError(err)
	}

	err = h.db.UnfollowUser(objID, removeUser.ID)
	if err != nil {
		return handleMgoError(err)
	}

	resp := models.UserToUserResponse(*removeUser)

	return c.JSON(http.StatusOK, &resp)
}

func (h *handler) followingListHandler(c echo.Context) error {
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

	id := c.Param("id")
	user, err := h.db.FindUser(id)
	if err != nil {
		return handleMgoError(err)
	}

	if len(user.Following) == 0 {
		return c.JSON(http.StatusOK, &usersResponse{})
	}

	users, err := h.db.FindUserByOIDArray(user.Following)
	if err != nil {
		return handleMgoError(err)
	}
	usersResp := models.UsersToUserResponseArray(users)
	return c.JSON(http.StatusOK, &usersResponse{Users: usersResp})
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

	id := c.Param("id")
	user, err := h.db.FindUser(id)
	if err != nil {
		return handleMgoError(err)
	}

	if len(user.Followers) == 0 {
		return c.JSON(http.StatusOK, &usersResponse{})
	}

	users, err := h.db.FindUserByOIDArray(user.Followers)
	if err != nil {
		return handleMgoError(err)
	}
	usersResp := models.UsersToUserResponseArray(users)

	return c.JSON(http.StatusOK, &usersResponse{Users: usersResp})
}
