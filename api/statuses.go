package api

import (
	"fmt"
	"net/http"
	"strconv"
	"unicode/utf8"

	"go.uber.org/zap"
	"gopkg.in/mgo.v2/bson"

	"github.com/TinyKitten/TimelineServer/config"
	"github.com/TinyKitten/TimelineServer/models"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
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

func (h *handler) newPostResponse(post models.Post, replyToUserID string) (*PostResponse, error) {
	sender, err := h.db.FindUserByOID(post.UserID)
	if err != nil {
		return nil, err
	}

	senderResp := models.UserToUserResponse(*sender)

	if replyToUserID != "" {
		toReply, err := h.db.FindUserByOID(post.UserID)
		if err != nil {
			return nil, err
		}
		return &PostResponse{
			Favorited: false,
			CreatedAt: post.CreatedAt.String(),
			ID:        post.ID.Hex(),
			Entities: models.PostEntity{
				URLs:         post.URLs,
				Hashtags:     post.Hashtags,
				UserMentions: []models.Post{},
			},
			Text:                post.Text,
			Shared:              false,
			SharedCount:         len(post.Shared),
			User:                senderResp,
			InReplyToScreenName: toReply.DisplayName,
			InReplyToUserID:     toReply.UserID,
		}, nil
	}
	return &PostResponse{
		Favorited: false,
		CreatedAt: post.CreatedAt.String(),
		ID:        post.ID.Hex(),
		Entities: models.PostEntity{
			URLs:         post.URLs,
			Hashtags:     post.Hashtags,
			UserMentions: []models.Post{},
		},
		Text:        post.Text,
		Shared:      false,
		SharedCount: len(post.Shared),
		User:        senderResp,
	}, nil
}

func (h *handler) getPublicPostsHandler(c echo.Context) error {
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

	limit, _ := strconv.Atoi(c.QueryParam("limit"))

	posts, err := h.db.GetPosts(limit)
	if err != nil {
		return handleMgoError(err)
	}

	return c.JSON(http.StatusOK, posts)
}

func (h *handler) updateStatusHandler(c echo.Context) error {
	jwtUser := c.Get("user").(*jwt.Token)
	claims := jwtUser.Claims.(jwt.MapClaims)
	idStr := claims["id"].(string)
	objID := bson.ObjectId(bson.ObjectIdHex(idStr))

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

	u, err := h.db.FindUserByOID(objID)
	if err != nil {
		return handleMgoError(err)
	}

	newPost := models.NewPost(u.ID.Hex(), req.InReplyToStatusID, req.Status)
	fmt.Println(newPost.Text)

	err = h.db.Insert("posts", newPost)
	if err != nil {
		h.logger.Debug("API Error", zap.String("Error", err.Error()))
		return handleMgoError(err)
	}

	go func(post models.Post, postChan chan models.Post) {
		postChan <- post
	}(*newPost, postChan)

	return c.JSON(http.StatusOK, &messageResponse{Message: "ok"})
}
