package actor

import (
	"reflect"

	"github.com/en-v/link/actor/action"
	"github.com/en-v/link/core"
	"github.com/en-v/link/state"
	"github.com/en-v/link/types"
	"github.com/en-v/log"
)

type Actor struct {
	state *state.State
}

func New(stt *state.State) (*Actor, error) {
	var err error

	if core.DEBUG {
		log.Debugw("CreateActor", "Type", stt.LocalTarget.Type)
	}

	var plt = reflect.TypeOf(types.Newp())

	for i := 0; i < stt.LocalTarget.Type.NumMethod(); i++ {

		reflectedMethod := stt.LocalTarget.Type.Method(i)
		if reflectedMethod.Type.NumIn() == 2 && reflectedMethod.Type.NumOut() == 2 {
			if reflectedMethod.Type.In(1) == plt && reflectedMethod.Type.Out(0) == plt {
				stt.LocalTarget.Actions[reflectedMethod.Name], err = action.New(reflectedMethod)
				if err != nil {
					return nil, err
				}
				if core.DEBUG {
					log.Debugw("CreateActor.RegisterActions", "Action", reflectedMethod.Name)
				}
			}
		}

	}

	return &Actor{
		state: stt,
	}, nil
}

func (this *Actor) GetLocalActions() []string {

	actions := make([]string, len(this.state.LocalTarget.Actions))
	i := 0

	for _, action := range this.state.LocalTarget.Actions {
		actions[i] = action.Name
		i++
	}

	return actions
}

func (this *Actor) GetLocalId() string {
	return this.state.LocalId
}

func (this *Actor) DeleteConnection(remoteId string) {
	for i := range this.state.Connections {
		if this.state.Connections[i] != nil && this.state.Connections[i].RemId == remoteId {
			this.state.Connections[i] = nil
			this.state.ClientsUnregFunc(remoteId)
			if core.DEBUG {
				log.Debug("Connection was deleted on server side")
			}
			return
		}
	}
}
