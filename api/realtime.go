package api

import (
	"encoding/json"
	"net/http"

	"github.com/TinyKitten/TimelineServer/models"

	"github.com/TinyKitten/TimelineServer/config"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/googollee/go-socket.io"
	"github.com/labstack/gommon/log"
	"go.uber.org/zap"
)

var postChan = make(chan models.Post)

const loggerTopic = "Socket.io"

func (h *handler) socketIOHandler() http.Handler {
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
				h.logger.Debug(loggerTopic, zap.String("Error", err.Error()))
				so.Emit("unauthorized", ErrInvalidJwt)
				so.Disconnect()
				return
			}

			claims := token.Claims.(jwt.MapClaims)

			claimID := claims["id"].(string)

			so.Emit("authenticated")

			so.Join(claimID)

			// 初回送信
			posts, err := h.db.GetAllPosts()
			if err != nil {
				h.logger.Debug(loggerTopic, zap.String("Error", err.Error()))
			}
			for _, post := range *posts {
				me, err := h.db.FindUser(claimID)
				if err != nil {
					handleMgoError(err)
				}
				resp, err := h.newPostResponse(post, "")
				if err != nil {
					handleMgoError(err)
				}
				j, err := json.Marshal(resp)
				if err != nil {
					h.logger.Debug(loggerTopic, zap.Any("Error", err.Error()))
					return
				}

				if post.UserID.Hex() == claimID {
					so.Emit(claimID, string(j))
					h.logger.Debug(loggerTopic, zap.Any("Sent", claimID))
				}

				for _, follow := range me.Following {
					if follow == post.UserID {
						so.Emit(claimID, string(j))
						h.logger.Debug(loggerTopic, zap.Any("Sent", claimID))
					}
				}
			}

			go func(postChan chan models.Post) {
				// 投稿監視
				for post := range postChan {
					sender, err := h.db.FindUser(post.UserID.Hex())
					if err != nil {
						handleMgoError(err)
					}

					resp, err := h.newPostResponse(post, "")
					if err != nil {
						handleMgoError(err)
					}

					j, err := json.Marshal(resp)
					if err != nil {
						h.logger.Debug(loggerTopic, zap.Any("Error", err.Error()))
						continue
					}
					so.Emit(claimID, string(j))
					h.logger.Debug(loggerTopic, zap.Any("Sent", claimID))

					if len(sender.Followers) == 0 {
						continue
					}

					if post.UserID.Hex() == claimID {
						so.Emit(claimID, string(j))
						h.logger.Debug(loggerTopic, zap.Any("Sent", claimID))
					}

					for _, follower := range sender.Followers {
						so.Emit(follower.Hex(), string(j))
						h.logger.Debug(loggerTopic, zap.Any("Sent", follower.Hex()))
					}
				}
			}(postChan)
		})
	})
	return server
}
