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

func (this *GateTarget) GetLinkHandlers() *types.Handlers {
	return &types.Handlers{
		ID:    this.id,
		Auth:  this.AuthCaller,
		Reg:   this.RegCaller,
		Unreg: this.UnregCaller,
	}
}

func (this *GateTarget) RegCaller(callerId string) error {
	return nil
}

func (this *GateTarget) UnregCaller(callerId string) error {
	return nil
}

func (this *GateTarget) AuthCaller(callerId string, token string) error {
	return nil
}
