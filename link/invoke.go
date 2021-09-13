package link

import (
	"github.com/en-v/link/core"
	"github.com/en-v/link/types"
	"github.com/pkg/errors"
)

// ActOn - call method on remote server
func (self *Link) Invoke(actionName string, payload *types.Payload) (*types.Payload, error) {

	if self.state.Mode != core.CLIENT_MODE {
		return nil, errors.New("Link is not in client mode, you need to use ActOnRemote for server mode or Connect before")
	}

	if self.state.Connections[0] == nil {
		return nil, errors.New("Remote connection is not found, try to Connect before Act calls")
	}
	remoteId := self.state.Connections[0].RemId

	res, err := self.actor.InvokeOnRemote(remoteId, actionName, payload)
	if err != nil {
		return nil, errors.Wrap(err, "Link.ActOnRemote.Act")
	}

	return res, nil
}

// InvokeOn - call method on remote client
func (self *Link) InvokeOn(callerId string, actionName string, payload *types.Payload) (*types.Payload, error) {

	if self.state.Mode != core.GATE_MODE {
		return nil, errors.New("Link is not in server mode, you need to use Act method")
	}

	res, err := self.actor.InvokeOnRemote(callerId, actionName, payload)
	if err != nil {
		return nil, errors.Wrap(err, "Link.ActOnRemote.ActOnRemote")
	}

	return res, nil
}

// InvokeOnCallers - call action on all remote clients
func (self *Link) Broadcast(actionName string, payload *types.Payload) (map[string]*types.Payload, map[string]error) {

	errs := make(map[string]error)

	if self.state.Mode != core.GATE_MODE {
		errs[self.state.LocalId] = errors.New("Link is not in server mode, you need to use Act method")
		return nil, errs
	}

	results := make(map[string]*types.Payload)

	for _, conn := range self.state.Connections {
		if conn != nil {
			res, err := self.actor.InvokeOnRemote(conn.RemId, actionName, payload)
			if err != nil {
				errs[conn.RemId] = errors.Wrap(err, "Link.ActOnRemote.ActOnRemote")
			} else {
				results[conn.RemId] = res
			}
		}
	}

	return results, errs
}
