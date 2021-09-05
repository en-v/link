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
	LocalId          string
	CallerToken      string
	Mode             core.LinkMode
	ClientsAuthFunc  types.ClientsAuthFunc
	ClientsRegFunc   types.ClientsRegFunc
	ClientsUnregFunc types.ClientsUnregFunc
	LocalTarget      *LocalTargetReflection
	Net              *NetParams
	Connections      []*connection.Connection
	Alias            string

	StopFallen chan byte
}

type NetParams struct {
	RemoteGateUrl string
	LocalGatePort int
	CertFile      string
	KeyFile       string
}

func New(localTarget types.LinkProxy) (*State, error) {

	if localTarget == nil {
		return nil, errors.New("Target is empty (NIL)")
	}

	levers := localTarget.GetLinkHandlers()
	if levers.ID == "" || len(levers.ID) < core.IDENTIFIERS_BYTES_COUNT {
		return nil, errors.New("ID is empty or too short")
	}

	localTargetReflection := &LocalTargetReflection{
		Actions:   make(map[string]*action.Action),
		Interface: localTarget,
		Value:     reflect.ValueOf(localTarget),
		Type:      reflect.TypeOf(localTarget),
	}

	return &State{
		LocalId:          levers.ID,
		ClientsAuthFunc:  levers.Auth,
		ClientsRegFunc:   levers.Reg,
		ClientsUnregFunc: levers.Unreg,
		LocalTarget:      localTargetReflection,
		Net:              &NetParams{},
		StopFallen:       make(chan byte),
	}, nil
}
