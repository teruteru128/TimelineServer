package api

import (
	"strconv"

	"github.com/TinyKitten/Timeline/config"
	"github.com/TinyKitten/Timeline/db"
	"github.com/TinyKitten/Timeline/logger"
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

	e.Logger.Fatal(e.Start(host))
}
