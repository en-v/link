package connection

import (
	"time"

	"github.com/en-v/link/core"
	"github.com/en-v/link/message"
	"github.com/en-v/log"
	"github.com/pkg/errors"
)

func (this *Connection) WriteJSON(msg *message.MsgV1) error {

	this.mutex.Lock()
	err := this.socket.WriteJSON(msg)
	this.mutex.Unlock()

	if err != nil {
		this.lastact = time.Now()
	}
	return err
}

func (this *Connection) SendRequst(req *message.MsgV1) error {

	if !this.enabled {
		return errors.New("Current connection is disabled)")
	}

	if req.Type != message.REQUEST {
		return errors.New("Current message is not request)")
	}

	err := this.WriteJSON(req)
	if err != nil {
		return errors.Wrap(err, "SendRequest.WriteJSON")
	}

	if req.Action == message.ACTION_KEEPALIVE {
		return nil
	}

	err = this.SetWaiter(req.MsgId)
	if err != nil {
		return errors.Wrap(err, "SendRequest.SetResponseWaiter")
	}

	return nil
}

func (this *Connection) WaitResponse(messageId string) (*message.MsgV1, error) {

	waiter, err := this.GetWaiter(messageId)
	if err != nil {
		return nil, errors.Wrap(err, "WaitResponse.WaitResponse")
	}

	select {
	case resp := <-waiter.Channel: // waiting for response message for current request
		this.DropWaiter(waiter.MsgId)
		return resp, nil

	case <-time.After(core.RESPONSE_WAITING_EXPIRE_INTERVAL): // ... or drop the waiting if time expired
		this.DropWaiter(waiter.MsgId)
		return nil, errors.New("Waiting for response time expired")
	}
}

func (this *Connection) SendResponse(resp *message.MsgV1) error {

	if core.DEBUG {
		log.Debugw("Send", "Response", resp.ToJson())
	}

	if resp.Type != message.RESPONSE {
		return errors.New("Current message is not response")
	}

	err := this.WriteJSON(resp)
	if err != nil {
		return errors.Wrap(err, "SendResponse.WriteJSON")
	}

	return nil
}
