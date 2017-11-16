package api

import (
	"net/http"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/TinyKitten/TimelineServer/config"
	"github.com/TinyKitten/TimelineServer/models"
	"github.com/TinyKitten/TimelineServer/sentence"
	"github.com/TinyKitten/TimelineServer/token"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/googollee/go-socket.io"
	"github.com/labstack/gommon/log"
	"go.uber.org/zap"
)

type singleton struct{}

var instance *singleton

func (h *handler) socketIOHandler() http.Handler {
	loggerTopic := "Socket.io"
	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}
	server.On("connection", func(so socketio.Socket) {
		// JWT Authentication
		so.On("authenticate", func(tokenReq string) {
			apiConfig := config.GetAPIConfig()
			claim := token.JwtClaim{}
			_, err := jwt.ParseWithClaims(tokenReq, &claim, func(token *jwt.Token) (interface{}, error) {
				return []byte(apiConfig.Jwt), nil
			})
			if err != nil {
				h.logger.Debug(loggerTopic, zap.String("Error", err.Error()))
				so.Emit("unauthorized", err.Error())
				so.Disconnect()
			}

			_, err = h.db.FindUserByOID(bson.ObjectId(bson.ObjectIdHex(claim.ID)))
			if err != nil {
				h.logger.Debug(loggerTopic, zap.String("Error", err.Error()))
				so.Emit("unauthorized", err.Error())
				so.Disconnect()
			}

			h.logger.Debug(loggerTopic, zap.String("Success", ""))
			so.Emit("authenticated")

			quit := make(chan bool)

			if instance == nil {
				h.logger.Debug(loggerTopic, zap.String("connection", "On connection"))
				so.On("sample", func(msg string) {
					h.logger.Debug(loggerTopic, zap.String("sample", msg))
				})

				so.On("disconnection", func() {
					h.logger.Debug(loggerTopic, zap.String("connection", "On disconnect"))
					quit <- true
					instance = nil
				})

				server.On("error", func(so socketio.Socket, err error) {
					log.Error("error:", err)
				})

				go func() {
					for {
						select {
						case <-quit:
							return
						default:
							user, post := sentence.GetRandomPost()
							respUser := models.UserToUserResponse(*user)
							resp := StreamPostResp{
								*post,
								respUser,
							}
							so.Emit("sample", resp)
							h.logger.Debug(loggerTopic, zap.Any("sample_emit", resp))

							time.Sleep(1 * time.Second)
						}
					}
				}()
				instance = new(singleton)
			}
		})
	})
	return server
}
