package connection

import (
	"github.com/en-v/link/core"
	"github.com/en-v/link/message"
	"github.com/en-v/link/types"
	"github.com/en-v/log"
)

func (this *Connection) handleRequst(req *message.MsgV1) {

	if req.Payload == nil {
		req.Payload = types.Newp()
	}
	req.Payload.SetSenderId(req.SenderId)

	if this.handleServiceRequst(req) {
		return
	}

	if core.DEBUG {
		log.Debugw("Request handler", "Req", req.ToJson())
	}

	res, err := this.actor.InvokeIncomingAction(req)
	if err != nil {
		errmsg := req.ToError(this.actor.GetLocalId(), err)
		err = this.SendResponse(errmsg)
		if err != nil {
			log.Error(err)
		}
		return
	}

	resp := req.ToResponse(this.actor.GetLocalId(), res)
	err = this.SendResponse(resp)
	if err != nil {
		log.Error(err)
	}
}

func (this *Connection) handleServiceRequst(req *message.MsgV1) bool {

	switch req.Action {

	case message.ACTION_KEEPALIVE:
		if core.DEBUG {
			log.Debug("Keep alive message", "From", req.SenderId, "LastAct", this.lastact)
		}
		go this.sendKeepAliveMessage()
		return true

	case message.ACTION_CLOSE_CONNECTION:
		this.Close(false)
		this.actor.DeleteConnection(req.SenderId)
		if core.DEBUG {
			log.Debug("Connection was closed on remote side")
		}
		return true

	default:
		return false
	}
}

func (this *Connection) handleResponse(resp *message.MsgV1) {

	rrw, err := this.GetWaiter(resp.MsgId)
	if err != nil {
		log.Error(err)
		return
	}

	if resp.Payload == nil {
		resp.Payload = types.Newp()
	}

	rrw.Channel <- resp
	this.DropWaiter(resp.MsgId)
}
