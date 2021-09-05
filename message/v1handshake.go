package message

import (
	"github.com/en-v/link/core"
	"github.com/pkg/errors"
)

func VerifyHandshakeResponse(response *MsgV1) error {
	if response.Type != RESPONSE || response.Action != ACTION_HANDSHAKE {
		return errors.New("Handshake response has wrong direction or wrong action name")
	}

	if response.RecvrId == "" || len(response.RecvrId) != core.IDENTIFIERS_BYTES_COUNT*2 {
		return errors.New("Handshake response server ID is undefined")
	}

	if response.Status != STATUS_OK {
		errstr, casterr := response.Payload.GetString(ERROR_PAYLOAD_FIELD)

		if casterr != nil {
			return errors.New("Handshake response status is wrong, error: " + errstr)
		} else {
			return errors.New("Handshake response status is wrong")
		}
	}

	actions, err := response.Payload.Get(HANDSHAKE_ACTIONS)
	if err != nil {
		return errors.Wrap(err, "VerifyHandshakeResponse")
	}

	if response.Payload == nil || actions == nil {
		return errors.New("Handshake response doesnt contain result/actions data")
	}
	return nil
}

func VerifyHandshakeRequest(request *MsgV1) error {

	if request.Type != REQUEST || request.Action != ACTION_HANDSHAKE {
		return errors.New("Received message is not a handshake request")
	}

	if request.SenderId == "" || len(request.SenderId) != core.IDENTIFIERS_BYTES_COUNT*2 {
		return errors.New("Handshake request client ID is wrong or undefined")
	}

	actions, err := request.Payload.Get(HANDSHAKE_ACTIONS)
	if err != nil {
		return errors.Wrap(err, "VerifyHandshakeRequest")
	}

	token, err := request.Payload.Get(HANDSHAKE_TOKEN)
	if err != nil {
		return errors.Wrap(err, "VerifyHandshakeRequest")
	}

	if request.Payload == nil || token == nil || actions == nil {
		return errors.New("Handshake request doesnt contain params/token or params/actions data")
	}
	return nil
}
