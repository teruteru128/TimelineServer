package api

import (
	"encoding/json"
	"net/http"

	"gopkg.in/mgo.v2/bson"

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
				if j, _ := h.checkFollow(claimID, post); j != nil {
					so.Emit(claimID, string(*j))
					h.logger.Debug(loggerTopic, zap.Any("Sent", j))
				}
			}

			go func(postChan chan models.Post) {
				// 投稿監視
				for post := range postChan {
					j, followers := h.checkFollow(claimID, post)
					if j != nil {
						so.Emit(claimID, string(*j))
						h.logger.Debug(loggerTopic, zap.Any("Sent", claimID))
						if followers != nil {
							for _, f := range followers {
								so.Emit(f.Hex(), string(*j))
								h.logger.Debug(loggerTopic, zap.Any("Sent", f.Hex()))
							}
						}
						h.logger.Debug(loggerTopic, zap.Any("Sent", j))
					}
				}
			}(postChan)
		})
	})
	return server
}

func (h *handler) checkFollow(claimID string, post models.Post) (*[]byte, []bson.ObjectId) {
	// フォローしている人か確認
	sender, err := h.db.FindUser(post.UserID)
	if err != nil {
		h.logger.Debug(loggerTopic, zap.String("Error", err.Error()))
		return nil, nil
	}
	respUser := models.UserToUserResponse(*sender)
	resp := StreamPostResp{
		post,
		respUser,
	}
	j, err := json.Marshal(resp)
	if err != nil {
		h.logger.Debug(loggerTopic, zap.String("Error", err.Error()))
		return nil, nil
	}

	// 自分がフォローしている
	for _, senderFollower := range sender.Followers {
		if senderFollower.Hex() == claimID {
			return &j, nil
		}
	}

	// 自分の投稿
	if sender.ID.Hex() == claimID {
		if len(sender.Followers) == 0 {
			return &j, nil
		}
		return &j, sender.Followers
	}

	return nil, nil
}
