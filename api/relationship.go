package api

import (
	"net/http"

	"github.com/TinyKitten/TimelineServer/models"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"gopkg.in/mgo.v2/bson"
)

type (
	usersResponse struct {
		Users []models.User `json:"users"`
	}
)

func (h *handler) followHandler(c echo.Context) error {
	jwtUser := c.Get("user").(*jwt.Token)
	claims := jwtUser.Claims.(jwt.MapClaims)
	idStr := claims["id"].(string)
	objID := bson.ObjectId(bson.ObjectIdHex(idStr))

	displayName := c.Param("id")
	followUser, err := h.db.FindUser(displayName)
	if err != nil {
		return handleMgoError(err)
	}

	err = h.db.FollowUser(objID, followUser.ID)
	if err != nil {
		return handleMgoError(err)
	}

	return c.JSON(http.StatusOK, &messageResponse{Message: RespFollowed})
}

func (h *handler) unfollowHandler(c echo.Context) error {
	jwtUser := c.Get("user").(*jwt.Token)
	claims := jwtUser.Claims.(jwt.MapClaims)
	idStr := claims["id"].(string)
	objID := bson.ObjectId(bson.ObjectIdHex(idStr))

	displayName := c.Param("id")
	followUser, err := h.db.FindUser(displayName)
	if err != nil {
		return handleMgoError(err)
	}

	err = h.db.UnfollowUser(objID, followUser.ID)
	if err != nil {
		return handleMgoError(err)
	}

	return c.JSON(http.StatusOK, &messageResponse{Message: RespUnfollowed})
}

func (h *handler) followingListHandler(c echo.Context) error {
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

	return c.JSON(http.StatusOK, &usersResponse{Users: users})
}

func (h *handler) followerListHandler(c echo.Context) error {
	id := c.Param("id")
	user, err := h.db.FindUser(id)
	if err != nil {
		return handleMgoError(err)
	}

	if len(user.Followers) == 0 {
		return c.JSON(http.StatusOK, &usersResponse{})
	}

	users, err := h.db.FindUserByOIDArray(user.Following)
	if err != nil {
		return handleMgoError(err)
	}

	return c.JSON(http.StatusOK, &usersResponse{Users: users})
}
