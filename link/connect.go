package link

import (
	"crypto/tls"

	"github.com/en-v/link/core"
	"github.com/en-v/log"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

func (self *Link) Reconnect(gateUrl string, token string) error {
	self.Shutdown()
	return self.Connect(gateUrl, token)
}

func (self *Link) Connect(gateUrl string, token string) error {

	err := self.verifyLinkAsCaller()
	if err != nil {
		return errors.Wrap(err, "Connect.VerifyLinkAsClient")
	}

	dialer := websocket.DefaultDialer
	dialer.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	socket, _, err := dialer.Dial(gateUrl, nil)
	if err != nil {
		return errors.Wrap(err, "Connect.WebScoektDial")
	}

	self.state.CallerToken = token
	conn, err := self.handshakeWithGate(socket)
	if err != nil {
		return errors.Wrap(err, "Connect.HandShakeToServer")
	}

	self.state.SetClientMode()
	self.state.Connections[0] = conn
	go self.state.Connections[0].Listen(core.CLIENT_MODE)

	if core.DEBUG {
		log.Debugw("Link as a client is connectedto", "Alias", self.state.Alias, "URL", gateUrl)
	}

	return nil
}
