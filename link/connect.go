package link

import (
	"crypto/tls"

	"github.com/en-v/link/core"
	"github.com/en-v/log"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

func (this *Link) Reconnect(gateUrl string, token string) error {
	this.Shutdown()
	return this.Connect(gateUrl, token)
}

func (this *Link) Connect(gateUrl string, token string) error {

	err := this.verifyLinkAsCaller()
	if err != nil {
		return errors.Wrap(err, "Connect.VerifyLinkAsClient")
	}

	dialer := websocket.DefaultDialer
	dialer.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	socket, _, err := dialer.Dial(gateUrl, nil)
	if err != nil {
		return errors.Wrap(err, "Connect.WebScoektDial")
	}

	this.state.CallerToken = token
	conn, err := this.handshakeWithGate(socket)
	if err != nil {
		return errors.Wrap(err, "Connect.HandShakeToServer")
	}

	this.state.SetClientMode()
	this.state.Connections[0] = conn
	go this.state.Connections[0].Listen(core.CLIENT_MODE)

	if core.DEBUG {
		log.Debugw("Link as a client is connectedto", "Alias", this.state.Alias, "URL", gateUrl)
	}

	return nil
}
