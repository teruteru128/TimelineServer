package v1

import (
	"github.com/TinyKitten/TimelineServer/config"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	validator "gopkg.in/go-playground/validator.v9"
)

func NewV1Router() *echo.Echo {
	h := NewHandler()

	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}

	apiConfig := config.GetAPIConfig()
	v1 := e.Group(apiConfig.Version)

	account := v1.Group("/account")
	account.POST("/create.json", h.AccountCreate)
	account.POST("/login.json", h.Login)
	account.GET("/settings.json", h.GetAccountSettings)

	account = v1.Group("/account")
	account.Use(middleware.JWT([]byte(apiConfig.Jwt)))
	account.POST("/settings.json", h.SetAccountSettings)
	account.POST("/update_profile_image.json", h.UpdateAccountProfileImage)

	users := v1.Group("/users")
	users.GET("/show.json", h.GetUser)

	// Administrator
	super := v1.Group("/super")
	super.Use(middleware.JWT([]byte(apiConfig.Jwt)))
	super.POST("/update_suspend.json", h.AUserSuspendHandler)
	super.POST("/update_official.json", h.ASetOfficialFlag)

	// Static
	v1.Static("/uploads", "uploads")

	// Friendship
	friendship := v1.Group("/friendships")
	friendship.Use(middleware.JWT([]byte(apiConfig.Jwt)))
	friendship.POST("/create.json", h.Follow)
	friendship.POST("/destroy.json", h.Unfollow)

	like := v1.Group("/like")
	like.Use(middleware.JWT([]byte(apiConfig.Jwt)))
	like.POST("/create.json", h.CreateLike)
	like.POST("/destroy.json", h.DestroyLike)

	friends := v1.Group("/friends")
	friends.GET("/ids.json", h.GetFriendsID)
	friends.GET("/list.json", h.GetFriendsList)

	followers := v1.Group("/followers")
	followers.GET("/ids.json", h.GetFollowersID)
	followers.GET("/list.json", h.GetFollowerList)

	statuses := v1.Group("/statuses")
	statuses.GET("/realtime.json", h.RealtimeHandler)
	statuses.GET("/union.json", h.UnionHandler)
	statuses.GET("/list.json", h.GetUserPosts)
	statuses.GET("/home.json", h.GetHomePosts)

	statuses.Use(middleware.JWT([]byte(apiConfig.Jwt)))
	statuses.POST("/update.json", h.UpdateStatus)

	search := v1.Group("/search")
	search.GET("/user.json", h.SearchUserHandler)

	return e
}
