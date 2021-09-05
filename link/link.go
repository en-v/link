package link

import (
	"context"
	"net"
	"net/http"

	"github.com/en-v/link/actor"
	"github.com/en-v/link/core"
	"github.com/en-v/link/state"
	"github.com/en-v/link/types"
	"github.com/en-v/log"
	"github.com/gorilla/websocket"
)

type Link struct {
	state    *state.State
	actor    *actor.Actor
	httpsrv  *http.Server
	upgrader websocket.Upgrader
}

func New(target types.LinkProxy) (*Link, error) {

	stt, err := state.New(target)
	if err != nil {
		return nil, err
	}

	actor, err := actor.New(stt)
	if err != nil {
		return nil, err
	}

	return &Link{
		state: stt,
		actor: actor,
	}, nil
}

func (this *Link) Shutdown() {

	this.state.StopFallen <- 1

	if this.state.Mode == core.GATE_MODE {
		err := this.httpsrv.Shutdown(context.TODO())
		if err != nil {
			log.Error(err)
		}
	}

	this.state.SetBlancMode()
}

func (this *Link) SetAlias(alias string) {
	this.state.Alias = alias
}

func (this *Link) GetClientToken(remoteId string) (string, error) {

	c, err := this.state.GetConnection(remoteId)
	if err != nil {
		return "", err
	}

	t := make([]byte, len(c.Token))
	copy(t, []byte(c.Token))

	return string(t), err
}

func (this *Link) GetClientIP(remoteId string) (*net.TCPAddr, error) {

	c, err := this.state.GetConnection(remoteId)
	if err != nil {
		return nil, err
	}
	return c.IP, err
}
