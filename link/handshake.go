package link

import (
	"encoding/json"

	"github.com/pkg/errors"

	"github.com/en-v/link/connection"
	"github.com/en-v/link/core"
	"github.com/en-v/link/message"
	"github.com/en-v/link/types"
	"github.com/en-v/log"
	"github.com/gorilla/websocket"
)

func (self *Link) handshakeWithGate(socket *websocket.Conn) (*connection.Connection, error) {

	payload := types.Newp()
	payload.Set(message.HANDSHAKE_ACTIONS, self.actor.GetLocalActions())
	payload.Set(message.HANDSHAKE_TOKEN, self.state.CallerToken)

	requestToGate := message.Request(self.state.LocalId, message.ACTION_HANDSHAKE, payload)

	err := socket.WriteJSON(requestToGate)
	if err != nil {
		return nil, errors.Wrap(err, "handshakeWithGate.WriteJSON")
	}

	responseFromGate := message.EmptyV1()
	err = socket.ReadJSON(responseFromGate)
	if err != nil {
		return nil, errors.Wrap(err, "handshakeWithGate.ReadJSON")
	}

	err = message.VerifyHandshakeResponse(responseFromGate)
	if err != nil {
		return nil, errors.Wrap(err, "handshakeWithGate.VerifyHandshakeResponse")
	}

	conn, err := connection.New(socket, responseFromGate, self.actor)
	if err != nil {
		return nil, errors.Wrap(err, "handshakeWithGate.ConnectionNew")
	}

	if core.DEBUG {
		log.Debug("Client handshake is successful")
	}

	return conn, nil
}

func (self *Link) handshakeFromCaller(socket *websocket.Conn) (*connection.Connection, error) {

	requestFromCaller := message.EmptyV1()
	_, data, err := socket.ReadMessage()
	if err != nil {
		return nil, errors.Wrap(err, "handshakeFromCaller.ReadMessage")
	}

	if core.DEBUG {
		log.Debug(data)
	}

	err = json.Unmarshal(data, requestFromCaller)
	if err != nil {
		return nil, errors.Wrap(err, "handshakeFromCaller.Unmarshal")
	}

	err = message.VerifyHandshakeRequest(requestFromCaller)
	if err != nil {
		return nil, err
	}

	token := ""
	if self.state.ClientsAuthFunc != nil {
		token, err = requestFromCaller.Payload.GetString("token")
		if err != nil {
			return nil, errors.Wrap(err, "handshakeFromCaller.Get.Token")
		}
		err = self.state.ClientsAuthFunc(requestFromCaller.SenderId, token)
		if err != nil {
			return nil, errors.New("Acees denied, the token is wrong, " + err.Error())
		}
	}

	//

	payload := types.Newp()
	payload.Set(message.HANDSHAKE_ACTIONS, self.actor.GetLocalActions())
	responseToCaller := requestFromCaller.ToResponse(self.state.LocalId, payload)

	err = socket.WriteJSON(responseToCaller)
	if err != nil {
		return nil, err
	}

	conn, err := connection.New(socket, requestFromCaller, self.actor)
	conn.Token = token
	if err != nil {
		return nil, err
	}

	if core.DEBUG {
		log.Debug("Gate's handshake is succefull")
	}

	return conn, nil
}
