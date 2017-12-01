package api

import (
	"strconv"

	"github.com/TinyKitten/TimelineServer/api/v1"
	"github.com/TinyKitten/TimelineServer/config"
	"github.com/labstack/echo/middleware"
)

// StartServer APIサーバを起動する
func StartServer() {
	r := v1.NewV1Router()

	r.Use(middleware.Logger())
	r.Use(middleware.Recover())
	r.Pre(middleware.RemoveTrailingSlash())
	r.Use(middleware.CORS())

	apiConfig := config.GetAPIConfig()
	port := strconv.Itoa(apiConfig.Port)

	host := ":" + port

	r.Logger.Fatal(r.Start(host))
}