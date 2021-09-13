package connection

import (
	"errors"
	"strconv"

	"github.com/en-v/link/connection/waiter"
)

func (self *Connection) GetWaiter(msgId string) (*waiter.Waiter, error) {
	for _, w := range self.waiters {
		if w != nil && w.MsgId == msgId {
			return w, nil
		}
	}
	return nil, errors.New("Waiter not found")
}

func (self *Connection) SetWaiter(msgId string) error {
	for i, w := range self.waiters {
		if w == nil {
			self.waiters[i] = waiter.New(msgId)
			return nil
		}
	}
	return errors.New("Empty slots not found, total " + strconv.Itoa(len(self.waiters)))
}

func (self *Connection) DropWaiter(msgId string) {
	for i, w := range self.waiters {
		if w != nil && w.MsgId == msgId {
			w.Drop()
			self.waiters[i] = nil
			return
		}
	}
}
