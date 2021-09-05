package connection

import (
	"time"

	"github.com/en-v/link/core"
	"github.com/en-v/link/message"
	"github.com/en-v/log"
)

func (this *Connection) Listen(mode core.LinkMode) {

	var err error
	this.enabled = true

	// keepalive mesaaging initialization
	if mode == core.CLIENT_MODE {
		go this.sendKeepAliveMessage()
	}

	if core.DEBUG {
		log.Debugw("Connection is listening...", "Mode", "", "RemoteID", this.RemId)
	}

	for this.enabled {

		msg := message.EmptyV1()
		err = this.socket.ReadJSON(msg)
		this.lastact = time.Now()

		if err != nil {
			log.Error(err)
			break
		}

		if msg.Type == message.REQUEST {
			go this.handleRequst(msg)
			continue
		}

		if msg.Type == message.RESPONSE {
			go this.handleResponse(msg)
			continue
		}

		log.Error("Unknown message direction (not request and not response)", msg.Type)
	}

	this.Close(false)
	if core.DEBUG {
		log.Debugw("Connection is stoped", "RemoteID", this.RemId, "Error", err, "IsServer", mode == core.GATE_MODE)
	}
}

func (this *Connection) sendKeepAliveMessage() {
	if this.enabled {

		time.Sleep(core.KEEPALIVE_INTERVAL)
		if this.enabled {

			req := message.Request(this.actor.GetLocalId(), message.ACTION_KEEPALIVE, nil)
			err := this.SendRequst(req)
			if err != nil {
				log.Error(err)
			}

			if core.DEBUG {
				log.Debug("Keep alive message", "To", this.RemId)
			}
		}
	}
}
