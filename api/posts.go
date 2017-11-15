package api

import (
	"net/http"
	"strconv"
	"unicode/utf8"

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
		return handleMgoError(err)
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

	if utf8.RuneCountInString(req.Text) > 140 {
		return &echo.HTTPError{Code: http.StatusRequestEntityTooLarge, Message: ErrTooLong}
	}

	newPost := models.NewPost(objID, req.Text)

	err := h.db.Create("posts", newPost)
	if err != nil {
		return handleMgoError(err)
	}

	return nil
}
