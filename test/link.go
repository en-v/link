package main

import (
	"time"

	bblink "github.com/en-v/link"
	"github.com/en-v/link/test/targets"
)

func main() {

	bblink.DebugOn()
	targets.TestServer()

	for i := 0; i < targets.CLIENTS; i++ {
		time.Sleep(time.Millisecond * 5)
		go targets.TestClient()
	}

	<-make(chan bool)
}
