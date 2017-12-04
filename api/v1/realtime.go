package v1

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/TinyKitten/TimelineServer/config"
	"github.com/TinyKitten/TimelineServer/models"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"go.uber.org/zap"
	"golang.org/x/net/websocket"
)

var postChan chan models.PostResponse

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

	websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()
		go func(postChan chan models.PostResponse) {
			for post := range postChan {
				bytes, err := json.Marshal(post)
				if err != nil {
					c.Logger().Error(err)
				}
				if post.User.ID == claimID {
					err = websocket.Message.Send(ws, string(bytes))
					if err != nil {
						c.Logger().Error(err)
					}
				} else {
					// 自分がフォローしている人の投稿
					if len(post.User.Followers) != 0 {
						for _, follower := range post.User.Followers {
							if claimID == follower.Hex() {
								err = websocket.Message.Send(ws, string(bytes))
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
			// Read
			msg := ""
			err = websocket.Message.Receive(ws, &msg)
			if err != nil {
				c.Logger().Error(err)
			}
			fmt.Printf("%s\n", msg)
		}
	}).ServeHTTP(c.Response(), c.Request())
	return nil
}
