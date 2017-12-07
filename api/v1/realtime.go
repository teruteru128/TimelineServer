package v1

import (
	"encoding/json"
	"net/http"

	"github.com/TinyKitten/TimelineServer/config"
	"github.com/TinyKitten/TimelineServer/models"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
	"go.uber.org/zap"
)

var (
	postChan chan models.PostResponse
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

const (
	loggerTopic = "Realtime Stream"
)

func (h *APIHandler) RealtimeHandler(c echo.Context) error {
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
	claimID := claims["id"].(string)

	postChan = make(chan models.PostResponse)

	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()
	defer close(postChan)

	go func(postChan chan models.PostResponse) {
		for post := range postChan {
			bytes, err := json.Marshal(post)
			if err != nil {
				c.Logger().Error(err)
			}
			if post.User.ID == claimID {
				err := ws.WriteMessage(websocket.TextMessage, bytes)
				if err != nil {
					c.Logger().Error(err)
				}
			} else {
				// 自分がフォローしている人の投稿
				if len(post.User.Followers) != 0 {
					for _, follower := range post.User.Followers {
						if claimID == follower.Hex() {
							err := ws.WriteMessage(websocket.TextMessage, bytes)
							if err != nil {
								c.Logger().Error(err)
							}
						}
					}
				}
			}
		}
	}(postChan)

	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			c.Logger().Error(err)
		}
	}

	return nil
}

func (h *APIHandler) UnionHandler(c echo.Context) error {
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

	postChan = make(chan models.PostResponse)

	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	go func(postChan chan models.PostResponse) {
		for post := range postChan {
			bytes, err := json.Marshal(post)
			if err != nil {
				c.Logger().Error(err)
			}
			// 無条件で送信
			err = ws.WriteMessage(websocket.TextMessage, bytes)
			if err != nil {
				c.Logger().Error(err)
			}
		}
	}(postChan)

	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			c.Logger().Error(err)
		}
	}

	return nil
}
