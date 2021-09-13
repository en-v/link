package actor

import (
	"github.com/en-v/link/connection"
	"github.com/en-v/link/message"
	"github.com/en-v/link/types"
	"github.com/pkg/errors"
)

func (self *Actor) InvokeOnRemote(remoteId string, actionName string, payload *types.Payload) (*types.Payload, error) {
	req, conn, err := self.prepareRequstAndConnection(remoteId, actionName, payload)
	if err != nil {
		return nil, errors.Wrap(err, "InvokeOnRemote.prepareRequstAndConnection")
	}

	err = conn.SendRequst(req)
	if err != nil {
		return nil, errors.Wrap(err, "InvokeOnRemote.SendRequst")
	}

	resp, err := conn.WaitResponse(req.MsgId)
	if err != nil {
		return nil, errors.Wrap(err, "InvokeOnRemote.WaitResponse")
	}

	if resp.Status == message.STATUS_ERR {
		errstt, err := resp.Payload.GetString("error")
		if err != nil {
			return nil, errors.Wrap(err, "Error field casting error")
		}
		return nil, errors.New(errstt)
	}

	return resp.Payload, nil
}

func (self *Actor) prepareRequstAndConnection(remoteId string, actionName string, payload *types.Payload) (*message.MsgV1, *connection.Connection, error) {

	conn, err := self.state.GetConnection(remoteId)
	if err != nil {
		return nil, nil, errors.Wrap(err, "prepareRequstAndConnection.FindConnectionById")
	}

	if !conn.RemoteActionExists(actionName) {
		return nil, nil, errors.New("Action with current name doesnt exist, " + actionName)
	}

	return message.Request(self.state.LocalId, actionName, payload), conn, nil
}
