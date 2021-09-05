package state

import (
	"errors"
	"time"

	"github.com/en-v/link/connection"
	"github.com/en-v/link/core"
	"github.com/en-v/log"
)

func (this *State) SetConnection(conn *connection.Connection) error {

	for i := range this.Connections {
		if this.Connections[i] != nil && this.Connections[i].RemId == conn.RemId {
			this.Connections[i].Close(false)
			this.Connections[i] = conn
			if core.DEBUG {
				log.Debugw("Set Connection", "Action", "Update exists", "Id", this.Connections[i].RemId, "Active", this.ActiveConnectionsCount())
			}
			return nil
		}
	}

	for i := range this.Connections {
		if this.Connections[i] == nil {
			this.Connections[i] = conn
			if core.DEBUG {
				log.Debugw("Set Connection", "Action", "Add new", "Id", this.Connections[i].RemId, "Active", this.ActiveConnectionsCount())
			}
			return nil
		}
	}
	return errors.New("Free slot connection not found (array is full)")
}

func (this *State) ActiveConnectionsCount() int {
	c := 0
	for i := range this.Connections {
		if this.Connections[i] != nil {
			c++
		}
	}
	return c
}

func (this *State) GetConnection(remoteId string) (*connection.Connection, error) {
	for _, c := range this.Connections {
		if c != nil && c.RemId == remoteId {
			return c, nil
		}
	}
	return nil, errors.New("Connection with current remoteId not found, " + remoteId)
}

func (this *State) CloseClientConnection(remoteId string) error {
	for num, c := range this.Connections {
		if c != nil && c.RemId == remoteId {
			this.dropConnection(num, c)
			return nil
		}
	}
	return errors.New("Connection with current remoteId not found, " + remoteId)
}

func (this *State) dropConnection(num int, c *connection.Connection) {
	c.Close(false)
	if this.ClientsUnregFunc != nil {
		this.ClientsUnregFunc(c.RemId)
	}
	this.Connections[num] = nil
}

func (this *State) DropFallenConnections() {

	alivers, fallen, available := 0, 0, 0
	ticker := time.NewTicker(core.ZOOMBIE_HUNTER_INERVAL).C

	for this.Mode == core.GATE_MODE {
		alivers, fallen, available = 0, 0, 0

		select {
		case <-ticker:
			for num, c := range this.Connections {
				if c != nil {
					if !c.IsActive() {
						log.Debugw("Drop connection", "Alias", this.Alias, "ID", c.RemId)
						this.dropConnection(num, c)
						fallen++
					} else {
						alivers++
					}
				} else {
					available++
				}
			}

			if core.DEBUG && fallen > 0 {
				log.Debugw("DropFallenConnections", "Alias", this.Alias, "Fallen", fallen, "Alive", alivers, "Available", available)
			}

			if available == 0 {
				panic("BLINK CONNECTIONS STACK IS FULL")
			}

		case <-this.StopFallen:
			log.Debug("Fallen connections watch-process stopped", "Alias", this.Alias)
			return
		}
	}
}
