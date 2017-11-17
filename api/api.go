package api

import (
	"net/http"
	"strconv"

	"github.com/TinyKitten/TimelineServer/config"
	"github.com/TinyKitten/TimelineServer/db"
	"github.com/TinyKitten/TimelineServer/logger"
	"github.com/gorilla/websocket"
	"github.com/rs/cors"
	validator "gopkg.in/go-playground/validator.v9"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

// StartServer APIサーバを起動する
func StartServer() {
	logger := logger.GetLogger()
	mongoIns := db.MongoInstance{
		Conf: config.GetDBConfig(),
	}
	h := handler{
		db:     &mongoIns,
		logger: logger,
	}

	apiConfig := config.GetAPIConfig()
	port := strconv.Itoa(int(apiConfig.Port))

	host := apiConfig.Endpoint + ":" + port
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.CORS())
	e.Validator = &customValidator{validator: validator.New()}

	// API Version 1 Base
	v1 := e.Group("/v1")

	// /v1/signup Handler
	v1.POST("/signup", h.signupHandler)
	// /v1/login Handler
	v1.POST("/login", h.loginHandler)

	// JWT RESTRICTED
	v1j := v1.Group("")
	v1j.Use(middleware.JWT([]byte(apiConfig.Jwt)))

	// /v1/posts Handlers(Restricted)
	posts := v1j.Group("/posts")
	posts.POST("", h.postHandler)
	// /v1/posts/public Handlers(Restricted)
	v1.GET("/posts/public", h.getPublicPostsHandler)

	// Not restricted /users
	users := v1.Group("/users")
	users.GET("/:id", h.getUserHandler)

	// Administrator
	v1j.POST("/suspend", h.userSuspendHandler)
	v1j.POST("/official", h.setOfficalFlagHandler)

	// Relationships
	v1j.PUT("/follow/:id", h.followHandler)
	v1j.PUT("/unfollow/:id", h.unfollowHandler)
	// Relations
	v1.GET("/following/:id", h.followingListHandler)
	v1.GET("/follower/:id", h.followerListHandler)

	// Restricted /users
	usersj := v1j.Group("/users")
	usersj.DELETE("/:id", h.userDeleteHandler)

	// Socket.io
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
	})
	sioHandler := c.Handler(h.socketIOHandler())
	e.GET("/socket.io", echo.WrapHandler(sioHandler))
	e.POST("/socket.io", echo.WrapHandler(sioHandler))

	e.Logger.Fatal(e.Start(host))
}
