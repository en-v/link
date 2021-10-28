package patterns

import (
	"github.com/en-v/link"
	"github.com/en-v/link/types"
	"github.com/pkg/errors"
)

type Gate struct {
	Id   string
	Link link.Link
}

func (g *Gate) GetId() string {
	return g.Id
}

func MakeGate() (*Gate, error) {
	gate := &Gate{}
	Link, err := link.New(gate)
	if err != nil {
		return nil, errors.Wrap(err, "Make Gate")
	}
	gate.Link = Link
	return gate, nil
}

func (g *Gate) Hooks() *types.Hooks {
	return &types.Hooks{
		LocalId:  g.GetId,
		Verify:   g.Verify,
		CheckIn:  g.CheckIn,
		CheckOut: g.CheckOut,
	}
}

func (g *Gate) CheckIn(remid string) error {
	return nil
}

func (g *Gate) CheckOut(remid string) error {
	return nil
}

func (g *Gate) Verify(remid, token string) error {
	return nil
}
