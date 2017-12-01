package v1

import (
	"encoding/json"
	"net/http"

	"github.com/TinyKitten/TimelineServer/models"

	"github.com/TinyKitten/TimelineServer/config"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/googollee/go-socket.io"
	"github.com/labstack/gommon/log"
	"go.uber.org/zap"
	"fmt"
)

var postChan = make(chan models.PostResponse)

const loggerTopic = "Socket.io"

func (h *APIHandler) SocketIO() http.Handler {
	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}

	server.On("connection", func(so socketio.Socket) {
		server.On("error", func(so socketio.Socket, err error) {
			h.logger.Debug(loggerTopic, zap.String("Error", err.Error()))
		})

		// JWT Authentication
		so.On("authenticate", func(tokenReq string) {

			so.On("disconnection", func() {
				h.logger.Debug(loggerTopic, zap.String("connection", "On disconnect"))
			})

			apiConfig := config.GetAPIConfig()
			token, err := jwt.Parse(tokenReq, func(token *jwt.Token) (interface{}, error) {
				return []byte(apiConfig.Jwt), nil
			})
			if err != nil {
				h.logger.Debug(loggerTopic, zap.String("Error", ErrInvalidJwt))
				so.Disconnect()
				return
			}

			claims := token.Claims.(jwt.MapClaims)

			claimID := claims["id"].(string)

			err = so.Emit("authenticated")
			if err != nil {
				h.logger.Error(loggerTopic, zap.Error(err))
			}

			err = so.Join(claimID)
			if err != nil {
				h.logger.Error(loggerTopic, zap.Error(err))
			}

			go func(postChan chan models.PostResponse) {
				// 投稿監視
				for post := range postChan {
					j, err := json.Marshal(post)
					if err != nil {
						h.logger.Error(loggerTopic, zap.Error(err))
						continue
					}

					fmt.Println(string(j))

					// UNION timeline
					err = so.Emit("union", string(j))
					if err != nil {
						h.logger.Error(loggerTopic, zap.Error(err))
					} else {
						h.logger.Debug(loggerTopic, zap.Any("Sent", "UNION"), zap.String("Data", string(j)))
					}

					err = so.Emit(claimID, string(j))
					if err != nil {
						h.logger.Error(loggerTopic, zap.Error(err))
					} else {
						h.logger.Debug(loggerTopic, zap.Any("Sent", claimID), zap.String("Data", string(j)))
					}

					if len(post.User.Followers) == 0 {
						continue
					}

					if post.User.UserID == claimID {
						err = so.Emit(claimID, string(j))
						if err != nil {
							h.logger.Error(loggerTopic, zap.Error(err))
						} else {
							h.logger.Debug(loggerTopic, zap.Any("Sent", claimID), zap.String("Data", string(j)))
						}
					}

					for _, follower := range post.User.Followers {
						err = so.Emit(follower.Hex(), string(j))
						if err != nil {
							h.logger.Error(loggerTopic, zap.Error(err))
						} else {
							h.logger.Debug(loggerTopic, zap.Any("Sent", claimID), zap.String("Data", string(j)))
						}
					}
				}
			}(postChan)
		})
	})
	return server
}
