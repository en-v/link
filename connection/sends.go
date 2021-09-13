package connection

import (
	"time"

	"github.com/en-v/link/core"
	"github.com/en-v/link/message"
	"github.com/en-v/log"
	"github.com/pkg/errors"
)

func (self *Connection) WriteJSON(msg *message.MsgV1) error {

	self.mutex.Lock()
	err := self.socket.WriteJSON(msg)
	self.mutex.Unlock()

	if err != nil {
		self.lastact = time.Now()
	}
	return err
}

func (self *Connection) SendRequst(req *message.MsgV1) error {

	if !self.enabled {
		return errors.New("Current connection is disabled)")
	}

	if req.Type != message.REQUEST {
		return errors.New("Current message is not request)")
	}

	err := self.WriteJSON(req)
	if err != nil {
		return errors.Wrap(err, "SendRequest.WriteJSON")
	}

	if req.Action == message.ACTION_KEEPALIVE {
		return nil
	}

	err = self.SetWaiter(req.MsgId)
	if err != nil {
		return errors.Wrap(err, "SendRequest.SetResponseWaiter")
	}

	return nil
}

func (self *Connection) WaitResponse(messageId string) (*message.MsgV1, error) {

	waiter, err := self.GetWaiter(messageId)
	if err != nil {
		return nil, errors.Wrap(err, "WaitResponse.WaitResponse")
	}

	select {
	case resp := <-waiter.Channel: // waiting for response message for current request
		self.DropWaiter(waiter.MsgId)
		return resp, nil

	case <-time.After(core.RESPONSE_WAITING_EXPIRE_INTERVAL): // ... or drop the waiting if time expired
		self.DropWaiter(waiter.MsgId)
		return nil, errors.New("Waiting for response time expired")
	}
}

func (self *Connection) SendResponse(resp *message.MsgV1) error {

	if core.DEBUG {
		log.Debugw("Send", "Response", resp.ToJson())
	}

	if resp.Type != message.RESPONSE {
		return errors.New("Current message is not response")
	}

	err := self.WriteJSON(resp)
	if err != nil {
		return errors.Wrap(err, "SendResponse.WriteJSON")
	}

	return nil
}
