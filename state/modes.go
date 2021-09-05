package state

import (
	"github.com/en-v/link/connection"
	"github.com/en-v/link/core"
)

func (this *State) SetBlancMode() {

	for i := range this.Connections {
		if this.Connections[i] != nil {
			this.Connections[i].Close(true)
			this.Connections[i] = nil
		}
	}

	this.Connections = nil
	this.Mode = core.BLANK_MODE
}

func (this *State) SetGateMode() {
	this.Mode = core.GATE_MODE
	this.Connections = make([]*connection.Connection, core.MAX_CONNECTIONS_PER_GATE_COUNT)
}

func (this *State) SetClientMode() {
	this.Mode = core.CLIENT_MODE
	this.Connections = make([]*connection.Connection, 1)
}
