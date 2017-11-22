package v1

import (
	"github.com/TinyKitten/TimelineServer/config"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/rs/cors"
)

func NewV1Router() *echo.Echo {
	h := NewHandler()

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.CORS())

	apiConfig := config.GetAPIConfig()
	v1 := e.Group(apiConfig.Version)

	account := v1.Group("/account")
	account.POST("/create.json", h.AccountCreate)
	account.POST("/login.json", h.Login)

	// JWT RESTRICTED
	v1j := v1.Group("")
	v1j.Use(middleware.JWT([]byte(apiConfig.Jwt)))

	account.GET("/settings.json", h.GetAccountSettings)
	accountj := v1j.Group("/account")
	accountj.POST("/settings.json", h.SetAccountSettings)
	accountj.POST("/update_profile_image.json", h.UpdateAccountProfileImage)

	users := v1.Group("/users")
	users.GET("/show.json", h.GetUser)

	// Administrator
	superj := v1j.Group("/super")
	superj.POST("/update_suspend.json", h.AUserSuspendHandler)
	superj.POST("/update_official.json", h.ASetOfficialFlag)

	// Static
	v1.Static("/static", "static")

	// Friendship
	friendshipj := v1j.Group("friendships")
	friendshipj.POST("/create.json", h.Follow)
	friendshipj.POST("/destroy.json", h.Unfollow)
	friends := v1.Group("/friends")
	friends.GET("/ids.json", h.GetFriendsID)
	friends.GET("/list.json", h.GetFriendsList)
	followers := v1.Group("/followers")
	followers.GET("/ids.json", h.GetFollowersID)
	followers.GET("/list.json", h.GetFollowerList)

	statusesj := v1j.Group("/statuses")
	statusesj.POST("/update.json", h.UpdateStatus)

	// Socket.io
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
	})
	sio := c.Handler(h.SocketIO())
	e.GET("/socket.io", echo.WrapHandler(sio))
	e.POST("/socket.io", echo.WrapHandler(sio))
	return e
}
