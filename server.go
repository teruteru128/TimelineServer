package main

import (
	"log"
	"runtime"

	"github.com/TinyKitten/TimelineServer/api"
	"github.com/TinyKitten/TimelineServer/config"
)

func main() {
	debugMode := config.GetAPIConfig().Debug

	if !debugMode {
		defer func() {
			err := recover()
			if err != nil {
				log.Println("panick recover. ", err)
			}
		}()
	}

	runtime.GOMAXPROCS(runtime.NumCPU())

	api.StartServer()
}
