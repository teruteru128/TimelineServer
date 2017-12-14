package api

import (
	"os"

	"github.com/TinyKitten/TimelineServer/api/v1"
	"github.com/TinyKitten/TimelineServer/config"
	"github.com/labstack/echo/middleware"
	"golang.org/x/crypto/acme/autocert"
)

// StartServer APIサーバを起動する
func StartServer() {
	r := v1.NewV1Router()

	r.Use(middleware.Logger())
	r.Use(middleware.Recover())
	r.Pre(middleware.RemoveTrailingSlash())
	r.Use(middleware.CORS())

	apiConfig := config.GetAPIConfig()
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	host := ":" + port

	if apiConfig.Secure {
		r.AutoTLSManager.Cache = autocert.DirCache(".cache")
		r.Logger.Fatal(r.StartAutoTLS(host))
	} else {
		r.Logger.Fatal(r.Start(host))
	}
}
