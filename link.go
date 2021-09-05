package link

import (
	"net"

	"github.com/en-v/link/core"
	"github.com/en-v/link/link"
	"github.com/en-v/link/types"
)

/* LINK */

type Link interface {
	// for all modes
	Shutdown()
	SetAlias(string)
	GetClientToken(string) (string, error)
	GetClientIP(string) (*net.TCPAddr, error)
	// use in client mode
	Reconnect(gateUrl string, token string) error
	Connect(gateUrl string, token string) error
	Invoke(action string, params *types.Payload) (*types.Payload, error)

	// use in gate mode
	Open(port int) error
	OpenLocal(port int) error
	OpenSecure(port int, certFile string, keyFile string) error

	InvokeOn(clientId string, action string, params *types.Payload) (*types.Payload, error)
	Broadcast(action string, params *types.Payload) (map[string]*types.Payload, map[string]error)
	Kickout(clientId string) error
}

func Payload() *types.Payload {
	return types.Newp()
}

func New(linkProxy types.LinkProxy) (Link, error) {
	return link.New(linkProxy)
}

//DebugOn - enable debug logging
func DebugOn() {
	core.DEBUG = true
}

//DebugOff - disable debug logging
func DebugOff() {
	core.DEBUG = false
}
