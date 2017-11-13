package api

import (
	"strconv"

	"github.com/TinyKitten/TimelineServer/config"
	"github.com/TinyKitten/TimelineServer/db"
	"github.com/TinyKitten/TimelineServer/logger"
	"github.com/gorilla/websocket"
	validator "gopkg.in/go-playground/validator.v9"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var (
	upgrader = websocket.Upgrader{}
)

// StartServer APIサーバを起動する
func StartServer() {
	mongoIns := db.MongoInstance{Conf: config.GetDBConfig()}
	logger := logger.GetLogger()
	h := handler{
		db:     &mongoIns,
		logger: logger,
	}

	apiConfig := config.GetAPIConfig()
	port := strconv.Itoa(int(apiConfig.Port))

	host := apiConfig.Endpoint + ":" + port
	e := echo.New()
	e.Use(middleware.CORS())
	e.Validator = &customValidator{validator: validator.New()}

	// API Version 1 Base
	v1 := e.Group("v1")

	// /v1/signup Handler
	v1.POST("/signup", h.signupHandler)
	// /v1/login Handler
	v1.POST("/login", h.loginHandler)

	// JWT RESTRICTED
	v1j := v1.Group("")
	v1j.Use(middleware.JWT([]byte(apiConfig.Jwt)))

	// /v1/posts Handlers(Restricted)
	posts := v1j.Group("/posts")
	posts.POST("/", h.postHandler)
	// posts.GET("/public_stream", h.getPublicStreamHandler)
	posts.GET("/sample_stream", h.getSampleStreamHandler)
	// /v1/posts/public Handlers(Restricted)
	pubPosts := posts.Group("/public")
	pubPosts.GET("/", h.getPublicPostsHandler)

	// Not restricted /users
	users := v1.Group("/users")
	users.GET("/:id", h.getUserHandler)

	// Restricted /users
	usersj := v1j.Group("/users")
	usersj.DELETE("/:id", h.userDeleteHandler)
	usersj.POST("/:id/suspend", h.userSuspendHandler)

	e.Logger.Fatal(e.Start(host))
}
