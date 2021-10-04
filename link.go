package link

import (
	"net"

	"github.com/en-v/link/core"
	"github.com/en-v/link/link"
	"github.com/en-v/link/types"
	"github.com/pkg/errors"
)

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

//Payload - create a new types.Payload object.
//It is needed as parameters for Invoke and InvokeOn methods.
func Payload() *types.Payload {
	return types.Newp()
}

//Results - create a complete Payload with one named field of data
func Results(field string, data interface{}) *types.Payload {
	return types.Newp().Sert(field, data)
}

//Result - create a complete Payload with data as payload
func Result(data interface{}) (*types.Payload, error) {
	res := types.Newp()
	err := res.Put(data)
	if err != nil {
		return nil, errors.Wrap(err, "Result")
	}
	return res, nil
}

//New - create a new Link instance, linkProxy is linked object which methods the Link will provide
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
