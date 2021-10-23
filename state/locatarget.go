package state

import (
	"errors"
	"reflect"

	"github.com/en-v/link/actor/action"
)

type HandlerReflection struct {
	Actions   map[string]*action.Action
	Interface interface{}
	Value     reflect.Value
	Type      reflect.Type
}

func (self *HandlerReflection) FindAction(actionName string) (*action.Action, error) {
	act, exists := self.Actions[actionName]
	if !exists {
		return nil, errors.New("Local target action not found, " + actionName)
	}
	return act, nil
}
