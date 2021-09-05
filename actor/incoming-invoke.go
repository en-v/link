package actor

import (
	"reflect"

	"github.com/en-v/link/message"
	"github.com/en-v/link/types"
	"github.com/pkg/errors"
)

func (this *Actor) InvokeIncomingAction(req *message.MsgV1) (*types.Payload, error) {
	action, err := this.state.LocalTarget.FindAction(req.Action)
	if err != nil {
		return nil, errors.Wrap(err, "ActOnLocal.FindActionByName")
	}

	arg := []reflect.Value{this.state.LocalTarget.Value, reflect.ValueOf(req.Payload)}
	res := action.Func.Call(arg)

	if res[1].Interface() != nil {
		err = res[1].Interface().(error)
		if err != nil {
			return nil, errors.Wrap(err, "Actor.InvokeOnLocal.Error")
		}
	}
	results := res[0].Interface().(*types.Payload)
	return results, nil
}
