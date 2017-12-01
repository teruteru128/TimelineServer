package v1

import (
	"net/http"
	"unicode/utf8"

	"go.uber.org/zap"
	"gopkg.in/mgo.v2/bson"

	"github.com/TinyKitten/TimelineServer/models"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/TinyKitten/TimelineServer/config"
)

type (
	PostReq struct {
		Status            string `json:"status" validate:"required"`
		InReplyToStatusID string `json:"in_reply_to_status_id"`
	}
	PostResponse struct {
		Favorited           bool                `json:"favorited"`
		CreatedAt           string              `json:"created_at"`
		ID                  string              `json:"id"`
		Entities            models.PostEntity   `json:"entities"`
		InReplyToUserID     string              `json:"in_reply_to_user_id"`
		Text                string              `json:"text"`
		Shared              bool                `json:"shared"`
		SharedCount         int                 `json:"shared_count"`
		User                models.UserResponse `json:"user"`
		InReplyToScreenName string              `json:"in_reply_to_screen_name"`
	}
)

func (h *APIHandler) UpdateStatus(c echo.Context) error {
	jwtUser := c.Get("user").(*jwt.Token)
	claims := jwtUser.Claims.(jwt.MapClaims)
	idStr := claims["id"].(string)
	id := bson.ObjectIdHex(idStr)

	req := new(PostReq)
	if err := c.Bind(req); err != nil {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: ErrParamsRequired}
	}
	if err := c.Validate(req); err != nil {
		h.logger.Debug("API Error", zap.String("Error", err.Error()))
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: ErrBadFormat}
	}

	if utf8.RuneCountInString(req.Status) > 140 {
		return &echo.HTTPError{Code: http.StatusRequestEntityTooLarge, Message: ErrTooLong}
	}

	u, err := h.db.FindUserByOID(id, true)
	if err != nil {
		return handleMgoError(err)
	}

	newPost := models.NewPost(u.ID, bson.ObjectId(req.InReplyToStatusID), req.Status)

	err = h.db.UpdatePost(*newPost)
	if err != nil {
		h.logger.Debug("API Error", zap.String("Error", err.Error()))
		return handleMgoError(err)
	}

	go func(post models.Post, postChan chan models.PostResponse) {
		resp := models.PostToPostResponse(post, *u)
		postChan <- resp
	}(*newPost, postChan)

	return c.JSON(http.StatusOK, &messageResponse{Message: "ok"})
}

func (h *APIHandler) GetUserPosts(c echo.Context) error {
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

	screenName := c.QueryParam("screen_name")
	userID := c.QueryParam("user_id")

	if screenName != "" {
		user, err := h.db.FindUser(screenName, true)
		if err != nil {
			return handleMgoError(err)
		}
		posts, err := h.db.GetPostsByOIDArray(user.Posts)
		if err != nil {
			return handleMgoError(err)
		}

		resp := models.PostsToPostResponseArray(posts, []models.User{*user}, true)

		return c.JSON(http.StatusOK, &resp)
	}

	if userID != "" {
		user, err := h.db.FindUserByOID(bson.ObjectIdHex(userID), true)
		if err != nil {
			return handleMgoError(err)
		}
		posts, err := h.db.GetPostsByOIDArray(user.Posts)
		if err != nil {
			return handleMgoError(err)
		}

		return c.JSON(http.StatusOK, &posts)

	}

	return c.JSON(http.StatusOK, &models.Post{})
}