package api

import (
	"strconv"

	"github.com/TinyKitten/TimelineServer/config"
	"github.com/TinyKitten/TimelineServer/db"
	"github.com/TinyKitten/TimelineServer/logger"
	"github.com/rs/cors"
	"go.uber.org/zap"
	"golang.org/x/crypto/acme/autocert"
	validator "gopkg.in/go-playground/validator.v9"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// StartServer APIサーバを起動する
func StartServer() {
	logger := logger.GetLogger()
	mongoIns, err := db.NewMongoInstance()
	if err != nil {
		logger.Panic("Failed to connect database.", zap.Skip())
	}
	h := handler{
		db:     mongoIns,
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
	v1 := e.Group(apiConfig.Version)

	account := v1.Group("/account")
	account.POST("/create.json", h.signupHandler)
	account.POST("/login.json", h.loginHandler)

	// JWT RESTRICTED
	v1j := v1.Group("")
	v1j.Use(middleware.JWT([]byte(apiConfig.Jwt)))

	account.GET("/settings.json", h.getSettingsHandler)
	accountj := v1j.Group("/account")
	accountj.POST("/settings.json", h.setSettingsHandler)
	accountj.POST("/update_profile_image.json", h.updateProfileImageHandler)

	v1.GET("/posts/public", h.getPublicPostsHandler)

	// Not restricted /users
	users := v1.Group("/users")
	users.GET("/show.json", h.getUserHandler)

	// Administrator
	superj := v1j.Group("/super")
	superj.POST("/update_suspend.json", h.userSuspendHandler)
	superj.POST("/update_official.json", h.setOfficalFlagHandler)
	superj.DELETE("destroy_user.json/:id", h.userDeleteHandler)

	// Static
	v1.Static("/static", "static")

	// Friendship
	friendshipj := v1j.Group("friendships")
	friendshipj.POST("/create.json", h.followHandler)
	friendshipj.POST("/destroy.json", h.unfollowHandler)
	friends := v1.Group("/friends")
	friends.GET("/ids.json", h.friendsIdsHandler)
	friends.GET("/list.json", h.friendsListHandler)
	followers := v1.Group("/followers")
	followers.GET("/ids.json", h.followerIdsHandler)
	followers.GET("/list.json", h.followerListHandler)

	statusesj := v1j.Group("/statuses")
	statusesj.POST("/update.json", h.updateStatusHandler)

	// Socket.io
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
	})
	sioHandler := c.Handler(h.socketIOHandler())
	e.GET("/socket.io", echo.WrapHandler(sioHandler))
	e.POST("/socket.io", echo.WrapHandler(sioHandler))

	e.AutoTLSManager.HostPolicy = autocert.HostWhitelist(apiConfig.Endpoint)
	e.AutoTLSManager.Cache = autocert.DirCache(".cache")
	e.Logger.Fatal(e.StartTLS(host, "cert.pem", "key.pem"))
}
