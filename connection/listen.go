package connection

import (
	"time"

	"github.com/en-v/link/core"
	"github.com/en-v/link/message"
	"github.com/en-v/log"
)

func (self *Connection) Listen(mode core.LinkMode) {

	var err error
	self.enabled = true

	// keepalive mesaaging initialization
	if mode == core.CLIENT_MODE {
		go self.sendKeepAliveMessage()
	}

	if core.DEBUG {
		log.Debugw("Connection is listening...", "Mode", "", "RemoteID", self.RemId)
	}

	for self.enabled {

		msg := message.EmptyV1()
		err = self.socket.ReadJSON(msg)
		self.lastact = time.Now()

		if err != nil {
			log.Error(err)
			break
		}

		if msg.Type == message.REQUEST {
			go self.handleRequst(msg)
			continue
		}

		if msg.Type == message.RESPONSE {
			go self.handleResponse(msg)
			continue
		}

		log.Error("Unknown message direction (not request and not response)", msg.Type)
	}

	self.Close(false)
	if core.DEBUG {
		log.Debugw("Connection is stoped", "RemoteID", self.RemId, "Error", err, "IsServer", mode == core.GATE_MODE)
	}
}

func (self *Connection) sendKeepAliveMessage() {
	if self.enabled {

		time.Sleep(core.KEEPALIVE_INTERVAL)
		if self.enabled {

			req := message.Request(self.actor.GetLocalId(), message.ACTION_KEEPALIVE, nil)
			err := self.SendRequst(req)
			if err != nil {
				log.Error(err)
			}

			if core.DEBUG {
				log.Debug("Keep alive message", "To", self.RemId)
			}
		}
	}
}
