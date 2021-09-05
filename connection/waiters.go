package connection

import (
	"errors"
	"strconv"

	"github.com/en-v/link/connection/waiter"
)

func (this *Connection) GetWaiter(msgId string) (*waiter.Waiter, error) {
	for _, w := range this.waiters {
		if w != nil && w.MsgId == msgId {
			return w, nil
		}
	}
	return nil, errors.New("Waiter not found")
}

func (this *Connection) SetWaiter(msgId string) error {
	for i, w := range this.waiters {
		if w == nil {
			this.waiters[i] = waiter.New(msgId)
			return nil
		}
	}
	return errors.New("Empty slots not found, total " + strconv.Itoa(len(this.waiters)))
}

func (this *Connection) DropWaiter(msgId string) {
	for i, w := range this.waiters {
		if w != nil && w.MsgId == msgId {
			w.Drop()
			this.waiters[i] = nil
			return
		}
	}
}
