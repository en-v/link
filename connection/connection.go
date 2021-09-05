package connection

import (
	"errors"
	"net"
	"sync"
	"time"

	"github.com/en-v/link/actor/iactor"
	"github.com/en-v/link/connection/waiter"
	"github.com/en-v/link/core"
	"github.com/en-v/link/message"
	"github.com/en-v/link/types"
	"github.com/en-v/log"
	"github.com/gorilla/websocket"
)

type Connection struct {
	RemId      string
	IP         *net.TCPAddr
	remActions []string
	waiters    []*waiter.Waiter
	socket     *websocket.Conn
	lastact    time.Time
	enabled    bool
	Token      string
	actor      iactor.IActor
	mutex      sync.Mutex
	termSig    chan byte
}

func New(socket *websocket.Conn, handShakeMessage *message.MsgV1, iactor iactor.IActor) (*Connection, error) {

	remoteId := ""

	if handShakeMessage.Type == message.RESPONSE { // response from the server
		remoteId = handShakeMessage.RecvrId
	}

	if handShakeMessage.Type == message.REQUEST { // request from the server
		remoteId = handShakeMessage.SenderId
	}

	if remoteId == "" {
		return nil, errors.New("Unknown message  direction")
	}

	conn := &Connection{
		waiters: make([]*waiter.Waiter, core.MAX_WAITERS_PER_CONNECTION_COUNT),
		socket:  socket,
		lastact: time.Now(),
		enabled: false,
		actor:   iactor,
		RemId:   remoteId,
		termSig: make(chan byte),
	}

	err := conn.applyRemoteActions(handShakeMessage.Payload)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func (this *Connection) Close(sendRequest bool) {

	if sendRequest {
		req := message.Request(this.actor.GetLocalId(), message.ACTION_CLOSE_CONNECTION, nil)
		err := this.SendRequst(req)
		if err != nil {
			log.Error(err)
		}
	}

	this.enabled = false
	this.socket.Close()
}

func (this *Connection) IsActive() bool {
	return time.Since(this.lastact) <= core.KEEPALIVE_INTERVAL*4
}

func (this *Connection) applyRemoteActions(payload *types.Payload) error {
	data, err := payload.Get(message.HANDSHAKE_ACTIONS)
	if err != nil {
		return errors.New("Actions field is not an array")
	}

	actions, ok := data.([]interface{})
	if !ok {
		return errors.New("Actions field casting error")
	}

	this.remActions = make([]string, len(actions))
	for a := range actions {
		this.remActions[a] = actions[a].(string)
	}
	if core.DEBUG {
		log.Debugw("Remote Actions", "Actions", this.remActions)
	}
	return nil
}

func (this *Connection) RemoteActionExists(action string) bool {
	for _, rem := range this.remActions {
		if rem == action {
			return true
		}
	}
	return false
}

func (this *Connection) Enable() {
	this.enabled = true
}

func (this *Connection) Disable() {
	this.enabled = false
}
