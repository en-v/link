package targets

import (
	"errors"
	"time"

	"github.com/en-v/goid"

	"github.com/en-v/link"
	"github.com/en-v/link/types"
	"github.com/en-v/log"
)

func TestClient() {

	token := "123"
	client := &Client{
		ID: goid.New(),
	}

	clientLink, err := link.New(client)
	if err != nil {
		panic(err)
	}

	err = clientLink.Connect("ws://localhost:8040", token)
	if err != nil {
		panic(err)
	}

	params := link.Payload()
	params.Set("param1", client.ID.String())
	params.Set("param2", client.ID.String())

	if !SILENT_MODE {
		go func() {
			ticker := time.NewTicker(time.Second * 8).C
			for {
				<-ticker

				res, err := clientLink.Invoke("ServerMethod", params)
				_ = res
				//log.Debugw("Method called on client", "Result", res, "Error", err)

				if err != nil {
					log.Error("Connection with server is lost")
					return
				}
			}
		}()
	}

	go func() {
		time.Sleep(time.Second * 30)
		clientLink.Shutdown()
	}()
}

type Client struct {
	ID goid.BBID
}

func (s *Client) GetLinkHandlers() *types.Handlers {
	return &types.Handlers{
		ID:    s.ID.String(),
		Auth:  nil,
		Reg:   nil,
		Unreg: nil,
	}
}

func (s *Client) SomeCleintMethod(params *types.Payload) (*types.Payload, error) {
	//log.Debugw("Cleitn method:", params["some"])
	return nil, nil
}

func (s *Client) SomeCleintMethod2(params *types.Payload) (*types.Payload, error) {
	//log.Debugw("Cleitn method:", params["some"])
	return nil, nil
}

func (c *Client) ClientError(a interface{}) error {
	log.Debug("Client method with error:", a)
	return errors.New("Client error")
}
