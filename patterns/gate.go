package patterns

import (
	"github.com/en-v/link"
	"github.com/en-v/link/types"
	"github.com/pkg/errors"
)

type GateTarget struct {
	id   string
	link link.Link
}

func NewGateTarget() (*GateTarget, error) {
	gateTarget := &GateTarget{}
	gateLink, err := link.New(gateTarget)
	if err != nil {
		return nil, errors.Wrap(err, "NodesGate.New.NewLink")
	}
	gateTarget.link = gateLink
	return gateTarget, nil
}

func (self *GateTarget) GetLinkHandlers() *types.Hooks {
	return &types.Hooks{
		LocalId:    self.id,
		Verify:  self.AuthCaller,
		CheckIn:   self.RegCaller,
		CheckOut: self.UnregCaller,
	}
}

func (self *GateTarget) RegCaller(callerId string) error {
	return nil
}

func (self *GateTarget) UnregCaller(callerId string) error {
	return nil
}

func (self *GateTarget) AuthCaller(callerId string, token string) error {
	return nil
}
