package action

import (
	"fmt"
	"reflect"

	"github.com/pkg/errors"
)

type Action struct {
	Name        string
	Func        reflect.Value
	Type        reflect.Type
	ParamType   reflect.Type
	ResultType  reflect.Type
	ResultLen   int
	Ins         []reflect.Value
	targetKnown bool
}

func New(m reflect.Method) (*Action, error) {

	if m.Type.NumIn() != 2 {
		return nil, errors.New("Count of input parameters is not 2 (must be an owner and a param), " + fmt.Sprint(m.Type) + " " + m.Name)
	}

	if m.Type.NumOut() > 2 {
		return nil, errors.New("Count of output results is not 1 or 2 (must be a result and an error or only the error)")
	}

	if m.Type.Out(m.Type.NumOut()-1).Name() != "error" {
		return nil, errors.New("One of output result is not an error")
	}

	return &Action{
		Name:        m.Name,
		Func:        m.Func,
		Type:        m.Type.In(0),
		ParamType:   m.Type.In(1),
		ResultType:  m.Type.Out(0),
		ResultLen:   m.Type.NumOut(),
		Ins:         make([]reflect.Value, 2),
		targetKnown: false,
	}, nil
}

func (this *Action) act(target interface{}, param interface{}) (interface{}, error) {

	if reflect.TypeOf(target) != this.Type {
		return nil, errors.New("Wrong owner data type")
	}

	if reflect.TypeOf(param) != this.ParamType {
		return nil, errors.New("Wrong parameter data type")
	}

	if !this.targetKnown {
		this.Ins[0] = reflect.ValueOf(target)
		this.targetKnown = true
	}
	this.Ins[1] = reflect.ValueOf(param)

	res := this.Func.Call(this.Ins)

	if this.ResultLen == 2 && res[1].Interface() != nil {
		return res[0].Interface(), res[1].Interface().(error)
	}
	return res[0].Interface(), nil

}
