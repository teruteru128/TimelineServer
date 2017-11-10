package api

import (
	"strconv"

	"github.com/TinyKitten/TimelineServer/config"
	"github.com/TinyKitten/TimelineServer/db"
	"github.com/TinyKitten/TimelineServer/logger"
	validator "gopkg.in/go-playground/validator.v9"

	"github.com/labstack/echo"
)

// StartServer APIサーバを起動する
func StartServer() {
	mongoIns := db.MongoInstance{Conf: config.GetDBConfig()}
	logger := logger.GetLogger()
	h := handler{db: &mongoIns, logger: logger}

	apiConfig := config.GetAPIConfig()
	port := strconv.Itoa(int(apiConfig.Port))

	host := apiConfig.Endpoint + ":" + port
	e := echo.New()
	e.Validator = &customValidator{validator: validator.New()}

	// API Version 1 Base
	v1 := e.Group("v1")

	// /v1/users Handlers
	v1.POST("/signup", h.signupHandler)
	v1.POST("/login", h.loginHandler)

	v1.GET("/:id", h.getUserHandler)
	v1.DELETE("/:id", h.userDeleteHandler)
	v1.POST("/:id/suspend", h.userSuspendHandler)

	e.Logger.Fatal(e.Start(host))
}
