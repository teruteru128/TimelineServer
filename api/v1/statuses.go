package v1

import (
	"net/http"
	"strconv"
	"unicode/utf8"

	"go.uber.org/zap"
	"gopkg.in/mgo.v2/bson"

	"github.com/TinyKitten/TimelineServer/config"
	"github.com/TinyKitten/TimelineServer/models"
	"github.com/TinyKitten/TimelineServer/utils"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

type (
	PostReq struct {
		Status            string `json:"status" validate:"required"`
		InReplyToStatusID string `json:"in_reply_to_status_id"`
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

	go func(post models.Post, postChan chan<- models.PostResponse) {
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

	limitStr := c.QueryParam("limit")
	cursorStr := c.QueryParam("cursor")
	screenName := c.QueryParam("screen_name")
	userID := c.QueryParam("user_id")

	if screenName != "" {
		// ScreenNameでの検索
		user, err := h.db.FindUser(screenName, true)
		if err != nil {
			return handleMgoError(err)
		}
		posts, err := h.db.GetPostsByOIDArray(user.Posts)
		if err != nil {
			return handleMgoError(err)
		}

		if limitStr == "" && cursorStr == "" {
			// リミット指定・カーソル指定なし
			resp := models.PostsToPostResponseArray(posts, []models.User{*user}, true)
			return c.JSON(http.StatusOK, &resp)
		}

		if limitStr == "" && cursorStr != "" {
			// リミット指定あり
			cursor, err := strconv.Atoi(cursorStr)
			if err != nil {
				h.logger.Debug("Param Error", zap.String("Error", err.Error()))
				return &echo.HTTPError{Code: http.StatusBadRequest, Message: ErrBadFormat}
			}

			if cursor >= len(posts) {
				return c.JSON(http.StatusOK, &[]models.PostResponse{})
			}

			resp := models.PostsToPostResponseArray(posts[cursor:], []models.User{*user}, true)
			return c.JSON(http.StatusOK, &resp)

		}

		if limitStr != "" && cursorStr == "" {
			// カーソル指定のみあり
			limit, err := strconv.Atoi(limitStr)
			if err != nil {
				h.logger.Debug("Param Error", zap.String("Error", err.Error()))
				return &echo.HTTPError{Code: http.StatusBadRequest, Message: ErrBadFormat}
			}
			resp := models.PostsToPostResponseArray(posts[:limit], []models.User{*user}, true)
			return c.JSON(http.StatusOK, &resp)
		}

		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			h.logger.Debug("Param Error", zap.String("Error", err.Error()))
			return &echo.HTTPError{Code: http.StatusBadRequest, Message: ErrBadFormat}
		}
		cursor, err := strconv.Atoi(cursorStr)
		if err != nil {
			h.logger.Debug("Param Error", zap.String("Error", err.Error()))
			return &echo.HTTPError{Code: http.StatusBadRequest, Message: ErrBadFormat}
		}
		if cursor >= len(posts) {
			return c.JSON(http.StatusOK, &[]models.PostResponse{})
		}
		resp := models.PostsToPostResponseArray(posts[cursor:limit], []models.User{*user}, true)
		return c.JSON(http.StatusOK, &resp)
	}

	if userID != "" {
		// IDでの検索
		user, err := h.db.FindUserByOID(bson.ObjectIdHex(userID), true)
		if err != nil {
			return handleMgoError(err)
		}
		posts, err := h.db.GetPostsByOIDArray(user.Posts)
		if err != nil {
			return handleMgoError(err)
		}

		if limitStr == "" && cursorStr == "" {
			// リミット指定・カーソル指定なし
			resp := models.PostsToPostResponseArray(posts, []models.User{*user}, true)
			return c.JSON(http.StatusOK, &resp)
		}

		if limitStr == "" && cursorStr != "" {
			// リミット指定あり
			cursor, err := strconv.Atoi(cursorStr)
			if err != nil {
				h.logger.Debug("Param Error", zap.String("Error", err.Error()))
				return &echo.HTTPError{Code: http.StatusBadRequest, Message: ErrBadFormat}
			}

			if cursor >= len(posts) {
				return c.JSON(http.StatusOK, &[]models.PostResponse{})
			}

			resp := models.PostsToPostResponseArray(posts[cursor:], []models.User{*user}, true)
			return c.JSON(http.StatusOK, &resp)

		}

		if limitStr != "" && cursorStr == "" {
			// カーソル指定のみあり
			limit, err := strconv.Atoi(limitStr)
			if err != nil {
				h.logger.Debug("Param Error", zap.String("Error", err.Error()))
				return &echo.HTTPError{Code: http.StatusBadRequest, Message: ErrBadFormat}
			}
			resp := models.PostsToPostResponseArray(posts[:limit], []models.User{*user}, true)
			return c.JSON(http.StatusOK, &resp)
		}

		// リミット指定・カーソル指定あり

		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			h.logger.Debug("Param Error", zap.String("Error", err.Error()))
			return &echo.HTTPError{Code: http.StatusBadRequest, Message: ErrBadFormat}
		}
		cursor, err := strconv.Atoi(cursorStr)
		if err != nil {
			h.logger.Debug("Param Error", zap.String("Error", err.Error()))
			return &echo.HTTPError{Code: http.StatusBadRequest, Message: ErrBadFormat}
		}
		if cursor >= len(posts) {
			return c.JSON(http.StatusOK, &[]models.PostResponse{})
		}
		resp := models.PostsToPostResponseArray(posts[cursor:limit], []models.User{*user}, true)

		return c.JSON(http.StatusOK, &resp)

	}

	return c.JSON(http.StatusOK, &[]models.Post{})
}

func (h *APIHandler) GetHomePosts(c echo.Context) error {
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

	limitStr := c.QueryParam("limit")
	cursorStr := c.QueryParam("cursor")

	user, err := h.db.FindUserByOID(bson.ObjectIdHex(id), true)
	if err != nil {
		return handleMgoError(err)
	}

	// 自分の投稿
	posts, err := h.db.GetPostsByOIDArray(user.Posts)
	if err != nil {
		return handleMgoError(err)
	}

	// フォローしている人の投稿

	friends, err := h.db.FindUserByOIDArray(user.Following, true)
	if err != nil {
		return handleMgoError(err)
	}
	for _, friend := range friends {
		friendPosts, err := h.db.GetPostsByOIDArray(friend.Posts)
		if err != nil {
			return handleMgoError(err)
		}
		posts = append(posts, friendPosts...)
	}

	posts = utils.SortByPostDates(posts)
	var senders []models.User

	for _, post := range posts {
		s, err := h.db.FindUserByOID(post.UserID, true)
		if err != nil {
			return handleMgoError(err)
		}
		senders = append(senders, *s)
	}

	if limitStr == "" && cursorStr == "" {
		// リミット指定・カーソル指定なし
		resp := models.PostsToPostResponseArray(posts, senders, false)
		return c.JSON(http.StatusOK, &resp)
	}

	if limitStr == "" && cursorStr != "" {
		// リミット指定あり
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			h.logger.Debug("Param Error", zap.String("Error", err.Error()))
			return &echo.HTTPError{Code: http.StatusBadRequest, Message: ErrBadFormat}
		}
		resp := models.PostsToPostResponseArray(posts[:limit], senders[:limit], false)
		return c.JSON(http.StatusOK, &resp)
	}

	if limitStr != "" && cursorStr == "" {
		// カーソル指定のみあり
		cursor, err := strconv.Atoi(cursorStr)
		if err != nil {
			h.logger.Debug("Param Error", zap.String("Error", err.Error()))
			return &echo.HTTPError{Code: http.StatusBadRequest, Message: ErrBadFormat}
		}
		if cursor >= len(posts) {
			return c.JSON(http.StatusOK, &[]models.PostResponse{})
		}
		resp := models.PostsToPostResponseArray(posts[cursor:], senders[cursor:], false)
		return c.JSON(http.StatusOK, &resp)

	}

	// リミット指定・カーソル指定あり
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		h.logger.Debug("Param Error", zap.String("Error", err.Error()))
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: ErrBadFormat}
	}
	cursor, err := strconv.Atoi(cursorStr)
	if err != nil {
		h.logger.Debug("Param Error", zap.String("Error", err.Error()))
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: ErrBadFormat}
	}
	if cursor >= len(posts) {
		return c.JSON(http.StatusOK, &[]models.PostResponse{})
	}
	resp := models.PostsToPostResponseArray(posts[cursor:limit], senders[cursor:limit], false)
	return c.JSON(http.StatusOK, &resp)
}
