package state

import (
	"errors"
	"reflect"

	"github.com/en-v/link/actor/action"
)

type LocalTargetReflection struct {
	Actions   map[string]*action.Action
	Interface interface{}
	Value     reflect.Value
	Type      reflect.Type
}

func (this *LocalTargetReflection) FindAction(actionName string) (*action.Action, error) {
	act, exists := this.Actions[actionName]
	if !exists {
		return nil, errors.New("Local target action not found, " + actionName)
	}
	return act, nil
}
