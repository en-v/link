package state

import (
	"errors"
	"reflect"

	"github.com/en-v/link/actor/action"
	"github.com/en-v/link/connection"
	"github.com/en-v/link/core"
	"github.com/en-v/link/types"
)

type State struct {
	LocalId     string
	CallerToken string
	Mode        core.LinkMode
	Vefify      types.VefifyFunc
	CheckIn     types.CheckInFunc
	CheckOut    types.CheckOutFunc
	Handler     *HandlerReflection
	Net         *NetParams
	Connections []*connection.Connection
	Alias       string

	StopFallen chan byte
}

type NetParams struct {
	RemoteGateUrl string
	LocalGatePort int
	CertFile      string
	KeyFile       string
}

func New(handler types.Handler) (*State, error) {

	if handler == nil {
		return nil, errors.New("Target is empty (NIL)")
	}

	hooks := handler.Hooks()
	if hooks.LocalId() == "" || len(hooks.LocalId()) < core.IDENTIFIERS_BYTES_COUNT {
		return nil, errors.New("ID is empty or too short")
	}

	handlerReflection := &HandlerReflection{
		Actions:   make(map[string]*action.Action),
		Interface: handler,
		Value:     reflect.ValueOf(handler),
		Type:      reflect.TypeOf(handler),
	}

	return &State{
		LocalId:    hooks.LocalId(),
		Vefify:     hooks.Verify,
		CheckIn:    hooks.CheckIn,
		CheckOut:   hooks.CheckOut,
		Handler:    handlerReflection,
		Net:        &NetParams{},
		StopFallen: make(chan byte),
	}, nil
}
