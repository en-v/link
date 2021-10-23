package state

import (
	"errors"
	"time"

	"github.com/en-v/link/connection"
	"github.com/en-v/link/core"
	"github.com/en-v/log"
)

func (self *State) SetConnection(conn *connection.Connection) error {

	for i := range self.Connections {
		if self.Connections[i] != nil && self.Connections[i].RemId == conn.RemId {
			self.Connections[i].Close(false)
			self.Connections[i] = conn
			if core.DEBUG {
				log.Debugw("Set Connection", "Action", "Update exists", "Id", self.Connections[i].RemId, "Active", self.ActiveConnectionsCount())
			}
			return nil
		}
	}

	for i := range self.Connections {
		if self.Connections[i] == nil {
			self.Connections[i] = conn
			if core.DEBUG {
				log.Debugw("Set Connection", "Action", "Add new", "Id", self.Connections[i].RemId, "Active", self.ActiveConnectionsCount())
			}
			return nil
		}
	}
	return errors.New("Free slot connection not found (array is full)")
}

func (self *State) ActiveConnectionsCount() int {
	c := 0
	for i := range self.Connections {
		if self.Connections[i] != nil {
			c++
		}
	}
	return c
}

func (self *State) GetConnection(remoteId string) (*connection.Connection, error) {
	for _, c := range self.Connections {
		if c != nil && c.RemId == remoteId {
			return c, nil
		}
	}
	return nil, errors.New("Connection with current remoteId not found, " + remoteId)
}

func (self *State) CloseClientConnection(remoteId string) error {
	for num, c := range self.Connections {
		if c != nil && c.RemId == remoteId {
			self.dropConnection(num, c)
			return nil
		}
	}
	return errors.New("Connection with current remoteId not found, " + remoteId)
}

func (self *State) dropConnection(num int, c *connection.Connection) {
	c.Close(false)
	if self.CheckOut != nil {
		self.CheckOut(c.RemId)
	}
	self.Connections[num] = nil
}

func (self *State) DropFallenConnections() {

	alivers, fallen, available := 0, 0, 0
	ticker := time.NewTicker(core.ZOOMBIE_HUNTER_INERVAL).C

	for self.Mode == core.GATE_MODE {
		alivers, fallen, available = 0, 0, 0

		select {
		case <-ticker:
			for num, c := range self.Connections {
				if c != nil {
					if !c.IsActive() {
						log.Debugw("Drop connection", "Alias", self.Alias, "ID", c.RemId)
						self.dropConnection(num, c)
						fallen++
					} else {
						alivers++
					}
				} else {
					available++
				}
			}

			if core.DEBUG && fallen > 0 {
				log.Debugw("DropFallenConnections", "Alias", self.Alias, "Fallen", fallen, "Alive", alivers, "Available", available)
			}

			if available == 0 {
				panic("BLINK CONNECTIONS STACK IS FULL")
			}

		case <-self.StopFallen:
			log.Debug("Fallen connections watch-process stopped", "Alias", self.Alias)
			return
		}
	}
}
