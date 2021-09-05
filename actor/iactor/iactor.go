package iactor

import (
	"github.com/en-v/link/message"
	"github.com/en-v/link/types"
)

type IActor interface {
	InvokeIncomingAction(*message.MsgV1) (*types.Payload, error)
	GetLocalId() string
	DeleteConnection(string)
}
