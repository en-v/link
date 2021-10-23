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
		log.Debugw("CreateActor", "Type", stt.Handler.Type)
	}

	var plt = reflect.TypeOf(types.Newp())

	for i := 0; i < stt.Handler.Type.NumMethod(); i++ {

		reflectedMethod := stt.Handler.Type.Method(i)
		if reflectedMethod.Type.NumIn() == 2 && reflectedMethod.Type.NumOut() == 2 {
			if reflectedMethod.Type.In(1) == plt && reflectedMethod.Type.Out(0) == plt {
				stt.Handler.Actions[reflectedMethod.Name], err = action.New(reflectedMethod)
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

func (self *Actor) GetLocalActions() []string {

	actions := make([]string, len(self.state.Handler.Actions))
	i := 0

	for _, action := range self.state.Handler.Actions {
		actions[i] = action.Name
		i++
	}

	return actions
}

func (self *Actor) GetLocalId() string {
	return self.state.LocalId
}

func (self *Actor) DeleteConnection(remoteId string) {
	for i := range self.state.Connections {
		if self.state.Connections[i] != nil && self.state.Connections[i].RemId == remoteId {
			self.state.Connections[i] = nil
			self.state.CheckOut(remoteId)
			if core.DEBUG {
				log.Debug("Connection was deleted on server side")
			}
			return
		}
	}
}
