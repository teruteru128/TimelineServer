package main

import (
	"log"
	"os"
	"runtime"

	"github.com/TinyKitten/TimelineServer/api"
)

func main() {
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
