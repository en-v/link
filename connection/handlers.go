package connection

import (
	"github.com/en-v/link/core"
	"github.com/en-v/link/message"
	"github.com/en-v/link/types"
	"github.com/en-v/log"
)

func (self *Connection) handleRequst(req *message.MsgV1) {

	if req.Payload == nil {
		req.Payload = types.Newp()
	}
	req.Payload.SetSenderId(req.SenderId)

	if self.handleServiceRequst(req) {
		return
	}

	if core.DEBUG {
		log.Debugw("Request handler", "Req", req.ToJson())
	}

	res, err := self.actor.InvokeIncomingAction(req)
	if err != nil {
		errmsg := req.ToError(self.actor.GetLocalId(), err)
		err = self.SendResponse(errmsg)
		if err != nil {
			log.Error(err)
		}
		return
	}

	resp := req.ToResponse(self.actor.GetLocalId(), res)
	err = self.SendResponse(resp)
	if err != nil {
		log.Error(err)
	}
}

func (self *Connection) handleServiceRequst(req *message.MsgV1) bool {

	switch req.Action {

	case message.ACTION_KEEPALIVE:
		if core.DEBUG {
			log.Debug("Keep alive message", "From", req.SenderId, "LastAct", self.lastact)
		}
		go self.sendKeepAliveMessage()
		return true

	case message.ACTION_CLOSE_CONNECTION:
		self.Close(false)
		self.actor.DeleteConnection(req.SenderId)
		if core.DEBUG {
			log.Debug("Connection was closed on remote side")
		}
		return true

	default:
		return false
	}
}

func (self *Connection) handleResponse(resp *message.MsgV1) {

	rrw, err := self.GetWaiter(resp.MsgId)
	if err != nil {
		log.Error(err)
		return
	}

	if resp.Payload == nil {
		resp.Payload = types.Newp()
	}

	rrw.Channel <- resp
	self.DropWaiter(resp.MsgId)
}
