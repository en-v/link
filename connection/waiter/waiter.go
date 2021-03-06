package waiter

import (
	"time"

	"github.com/en-v/link/core"
	"github.com/en-v/link/message"
)

type Waiter struct {
	MsgId   string
	Created time.Time
	Channel chan *message.MsgV1
}

func New(id string) *Waiter {
	return &Waiter{
		MsgId:   id,
		Created: time.Now(),
		Channel: make(chan *message.MsgV1),
	}
}

func (self *Waiter) Drop() {
	close(self.Channel)
}

func (self *Waiter) IsExpired() bool {
	return time.Since(self.Created) > core.RESPONSE_WAITING_EXPIRE_INTERVAL
}
