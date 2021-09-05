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

func (this *Link) handshakeWithGate(socket *websocket.Conn) (*connection.Connection, error) {

	payload := types.Newp()
	payload.Set(message.HANDSHAKE_ACTIONS, this.actor.GetLocalActions())
	payload.Set(message.HANDSHAKE_TOKEN, this.state.CallerToken)

	requestToGate := message.Request(this.state.LocalId, message.ACTION_HANDSHAKE, payload)

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

	conn, err := connection.New(socket, responseFromGate, this.actor)
	if err != nil {
		return nil, errors.Wrap(err, "handshakeWithGate.ConnectionNew")
	}

	if core.DEBUG {
		log.Debug("Client handshake is successful")
	}

	return conn, nil
}

func (this *Link) handshakeFromCaller(socket *websocket.Conn) (*connection.Connection, error) {

	requestFromCaller := message.EmptyV1()
	_, data, err := socket.ReadMessage()
	if err != nil {
		return nil, errors.Wrap(err, "handshakeFromCaller.ReadMessage")
	}

	if core.DEBUG {
		log.Debug(data)
	}
	//TODO
	err = json.Unmarshal(data, requestFromCaller)
	if err != nil {
		return nil, errors.Wrap(err, "handshakeFromCaller.Unmarshal")
	}

	err = message.VerifyHandshakeRequest(requestFromCaller)
	if err != nil {
		return nil, err
	}

	token := ""
	if this.state.ClientsAuthFunc != nil {
		token, err = requestFromCaller.Payload.GetString("token")
		if err != nil {
			return nil, errors.Wrap(err, "handshakeFromCaller.Get.Token")
		}
		err = this.state.ClientsAuthFunc(requestFromCaller.SenderId, token)
		if err != nil {
			return nil, errors.New("Acees denied, token is wrong, " + err.Error())
		}
	}

	//

	payload := types.Newp()
	payload.Set(message.HANDSHAKE_ACTIONS, this.actor.GetLocalActions())
	responseToCaller := requestFromCaller.ToResponse(this.state.LocalId, payload)

	err = socket.WriteJSON(responseToCaller)
	if err != nil {
		return nil, err
	}

	conn, err := connection.New(socket, requestFromCaller, this.actor)
	conn.Token = token
	if err != nil {
		return nil, err
	}

	if core.DEBUG {
		log.Debug("Gate handshake is succefull")
	}

	return conn, nil
}
