package main

import (
	"time"

	"github.com/en-v/link"
	"github.com/en-v/link/example/targets"
	"github.com/en-v/log"
)

func main() {
	log.Init(log.GREEN, "LINK")

	link.DebugOn()
	targets.TestServer()

	for i := 0; i < targets.CLIENTS; i++ {
		time.Sleep(time.Millisecond * 5)
		go targets.TestClient()
	}

	<-make(chan bool)
}
