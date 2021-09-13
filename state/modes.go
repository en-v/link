package state

import (
	"github.com/en-v/link/connection"
	"github.com/en-v/link/core"
)

func (self *State) SetBlancMode() {

	for i := range self.Connections {
		if self.Connections[i] != nil {
			self.Connections[i].Close(true)
			self.Connections[i] = nil
		}
	}

	self.Connections = nil
	self.Mode = core.BLANK_MODE
}

func (self *State) SetGateMode() {
	self.Mode = core.GATE_MODE
	self.Connections = make([]*connection.Connection, core.MAX_CONNECTIONS_PER_GATE_COUNT)
}

func (self *State) SetClientMode() {
	self.Mode = core.CLIENT_MODE
	self.Connections = make([]*connection.Connection, 1)
}
