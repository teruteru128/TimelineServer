package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/TinyKitten/TimelineServer/sentence"

	"github.com/TinyKitten/TimelineServer/models"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
)

type (
	PostReq struct {
		Text string `json:"text" validate:"required"`
	}
	StreamPostResp struct {
		models.Post
		models.UserResponse `json:"user"`
	}
)

func (h *handler) getSampleStreamHandler(c echo.Context) error {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	objID := claims["id"].(string)

	if h.checkSuspended(objID) {
		return &echo.HTTPError{Code: http.StatusForbidden, Message: ErrSuspended}
	}

	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	for {
		user, post := sentence.GetRandomPost()
		respUser := models.UserToUserResponse(*user)
		resp := StreamPostResp{
			*post,
			respUser,
		}

		// Write
		err := ws.WriteJSON(resp)
		if err != nil {
			c.Logger().Error(err)
		}

		time.Sleep(1 * time.Second)
	}
}

func (h *handler) getPublicPostsHandler(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	objID := claims["id"].(string)

	if h.checkSuspended(objID) {
		return &echo.HTTPError{Code: http.StatusForbidden, Message: ErrSuspended}
	}

	limit, _ := strconv.Atoi(c.QueryParam("limit"))

	posts, err := h.db.GetAllPosts(limit)
	if err != nil {
		return &echo.HTTPError{Code: http.StatusInternalServerError, Message: ErrUnknown}
	}

	return c.JSON(http.StatusOK, posts)
}

func (h *handler) postHandler(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	objID := claims["id"].(string)

	if h.checkSuspended(objID) {
		return &echo.HTTPError{Code: http.StatusForbidden, Message: ErrSuspended}
	}

	req := new(PostReq)
	if err := c.Bind(req); err != nil {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: ErrParamsRequired}
	}
	if err := c.Validate(req); err != nil {
		log.Error(err.Error())
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: ErrBadFormat}
	}

	newPost := models.NewPost(objID, req.Text)

	err := h.db.Create("posts", newPost)
	if err != nil {
		return &echo.HTTPError{Code: http.StatusInternalServerError, Message: ErrUnknown}
	}

	return nil
}
