package patterns

import (
	"github.com/en-v/link"
	"github.com/en-v/link/types"
	"github.com/pkg/errors"
)

type Client struct {
	Id   string
	Link link.Link
}

func (c *Client) GetId() string {
	return c.Id
}

func NewClient() (*Client, error) {
	client := &Client{}
	Link, err := link.New(client)
	if err != nil {
		return nil, errors.Wrap(err, "Make Client")
	}
	client.Link = Link
	return client, nil
}

func (c *Client) Hooks() *types.Hooks {
	return &types.Hooks{
		LocalId: c.GetId,
	}
}
