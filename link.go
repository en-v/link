package link

import (
	"net"

	"github.com/en-v/link/core"
	"github.com/en-v/link/link"
	"github.com/en-v/link/types"
)

/* LINK */

type Link interface {
	// For all modes
	Shutdown()
	SetAlias(string)
	GetClientToken(string) (string, error)
	GetClientIP(string) (*net.TCPAddr, error)

	// Used in the client mode
	Reconnect(gateUrl string, token string) error
	Connect(gateUrl string, token string) error
	Invoke(action string, params *types.Payload) (*types.Payload, error)

	// Used in the gate mode
	Open(port int) error
	OpenLocal(port int) error
	OpenSecure(port int, certFile string, keyFile string) error

	InvokeOn(clientId string, action string, params *types.Payload) (*types.Payload, error)
	Broadcast(action string, params *types.Payload) (map[string]*types.Payload, map[string]error)
	Kickout(clientId string) error
}

//Payload - create new types.Payload object.
//It needs as parameters for Invoke and InvokeOn methods.
func Payload() *types.Payload {
	return types.Newp()
}

//New - create new Link instance.
//linkProxy - linked object which methods the Link will provide
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
