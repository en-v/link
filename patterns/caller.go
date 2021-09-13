package patterns

import (
	"github.com/en-v/link"
	"github.com/en-v/link/types"
	"github.com/pkg/errors"
)

type CallerTarget struct {
	id   string
	link link.Link
}

func NewCallerTraget() (*CallerTarget, error) {
	callerTarget := &CallerTarget{}
	callerLink, err := link.New(callerTarget)
	if err != nil {
		return nil, errors.Wrap(err, "NodesGate.New.NewLink")
	}
	callerTarget.link = callerLink
	return callerTarget, nil
}

func (self *CallerTarget) GetLinkHandlers() *types.Handlers {
	return &types.Handlers{
		ID: self.id,
	}
}
