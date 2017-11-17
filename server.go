package main

import (
	"runtime"

	"github.com/TinyKitten/TimelineServer/api"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	api.StartServer()
}
