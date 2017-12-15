package main

import (
	"log"
	"os"
	"runtime"

	"github.com/TinyKitten/TimelineServer/api"
	"github.com/spf13/viper"
)

func initializeConfig() {
	viper.SetDefault("Version", "v1")
	viper.SetDefault("Debug", false)
	viper.SetDefault("Endpoint", "localhost")
	viper.SetDefault("JWT_TOKEN", "DEFAULT_JWT_TOKEN_CHANGE_ME")
	viper.SetDefault("Secure", false)

	viper.SetDefault("DB_HOST", "mongo")
	viper.SetDefault("DB_USERNAME", "")
	viper.SetDefault("DB_PASSWORD", "")
	viper.SetDefault("DB_PORT", 27017)
	viper.SetDefault("DB_NAME", "timeline_dev")

	viper.SetDefault("REDIS_HOST", "localhost")
	viper.SetDefault("REDIS_PASSWORD", "")
	viper.SetDefault("REDIS_PORT", 6379)

	viper.SetDefault("IMAGE_UPLOAD_PATH", "/uploads/img")
}

func main() {
	initializeConfig()

	f, _ := os.Create("./server.log")
	defer f.Close()
	log.SetOutput(f)

	defer func() {
		err := recover()
		if err != nil {
			log.Println("panic recover. ", err)
		}
	}()

	runtime.GOMAXPROCS(runtime.NumCPU())

	api.StartServer()
}
