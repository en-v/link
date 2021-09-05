package main

import (
	"time"

	link "github.com/en-v/link"
	"github.com/en-v/link/test/targets"
)

func main() {

	link.DebugOn()
	targets.TestServer()

	for i := 0; i < targets.CLIENTS; i++ {
		time.Sleep(time.Millisecond * 5)
		go targets.TestClient()
	}

	<-make(chan bool)
}
