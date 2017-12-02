package v1

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/TinyKitten/TimelineServer/config"
	"github.com/TinyKitten/TimelineServer/models"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/googollee/go-socket.io"
	"go.uber.org/zap"
)

const (
	loggerTopic = "Socket.io"
)

var postChan chan models.PostResponse
var chanClosed chan struct{}

func (h *APIHandler) SocketIO() http.Handler {
	chanClosed = make(chan struct{})

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
			postChan = make(chan models.PostResponse)

			so.On("disconnection", func() {
				h.logger.Debug(loggerTopic, zap.String("connection", "On disconnect"))
				close(chanClosed)
			})

			apiConfig := config.GetAPIConfig()
			token, err := jwt.Parse(tokenReq, func(token *jwt.Token) (interface{}, error) {
				return []byte(apiConfig.Jwt), nil
			})
			if err != nil {
				h.logger.Debug(loggerTopic, zap.String("Error", ErrInvalidJwt))
				so.Disconnect()
				close(chanClosed)
			}

			claims := token.Claims.(jwt.MapClaims)

			claimID := claims["id"].(string)

			err = so.Emit("authenticated")
			if err != nil {
				h.logger.Error(loggerTopic, zap.Error(err))
			}
			// UNION Timeline
			// 投稿監視
			for {
				select {
				case post := <-postChan:
					j, err := json.Marshal(post)
					if err != nil {
						h.logger.Error(loggerTopic, zap.Error(err))
						continue
					}

					// UNION timeline
					err = so.Emit("union", string(j))
					if err != nil {
						h.logger.Error(loggerTopic, zap.Error(err))
					} else {
						h.logger.Debug(loggerTopic, zap.Any("Sent", "UNION"), zap.String("Data", string(j)))
					}

					// 自分の投稿
					if post.User.ID == claimID {
						err = so.Emit("home", string(j))
						if err != nil {
							h.logger.Error(loggerTopic, zap.Error(err))
						} else {
							h.logger.Debug(loggerTopic, zap.Any("Sent", claimID), zap.String("Data", string(j)))
						}
						continue
					} else {
						// 自分がフォローしている人の投稿
						if len(post.User.Followers) != 0 {
							for _, follower := range post.User.Followers {
								if claimID == follower.Hex() {
									err = so.Emit("home", string(j))
									if err != nil {
										h.logger.Error(loggerTopic, zap.Error(err))
									} else {
										h.logger.Debug(loggerTopic, zap.Any("Sent broadcast", claimID), zap.String("Data", string(j)))
									}
								}
							}
						}
					}
				case <-chanClosed:
					return
				}
			}
		})
	})
	return server
}
