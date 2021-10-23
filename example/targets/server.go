package targets

import (
	"errors"
	"sync"
	"time"

	"github.com/en-v/goid"
	"github.com/en-v/link"
	"github.com/en-v/link/types"
	"github.com/en-v/log"
)

func TestServer() {

	s := &Server{
		Clients: make(map[string]string),
		Id:      goid.New(),
	}

	link, err := link.New(s)
	if err != nil {
		panic(err)
	}

	err = link.Open(8040)
	if err != nil {
		panic(err)
	}

	if !SILENT_MODE {
		go func() {
			ticker := time.NewTicker(time.Second * 6).C
			for {
				<-ticker
				for remId := range s.Clients {
					res, err := link.InvokeOn(s.Clients[remId], "SomeCleintMethod", nil)
					log.Debugw("Method called on server", "RemId", remId, "Result", res, "Error", err)
				}
				<-ticker
				res, err := link.Broadcast("SomeCleintMethod", nil)
				log.Debugw("Method called on server", "Result", res, "Error", err)
			}
		}()
	}
}

type Server struct {
	Clients map[string]string
	Id      goid.BBID
	mu      sync.Mutex
}

func (s *Server) Hooks() *types.Hooks {
	return &types.Hooks{
		LocalId:  s.Id.String,
		Verify:   s.AuthClient,
		CheckIn:  s.RegisterClient,
		CheckOut: s.UnregisterClient,
	}
}

func (s *Server) AuthClient(ClientId string, Token string) error {
	if Token != "123" {
		return errors.New("Token is workng")
	}
	return nil
}

func (s *Server) ServerMethod(params *types.Payload) (*types.Payload, error) {
	log.Debug("Server method", "Params", params)
	return link.Results("mirrored_params", params), nil
}

func (s *Server) RegisterClient(remoteId string) error {
	s.mu.Lock()
	s.Clients[remoteId] = remoteId
	s.mu.Unlock()
	log.Infow("New connection", "ID", remoteId)
	return nil
}

func (s *Server) UnregisterClient(remoteId string) error {
	s.mu.Lock()
	delete(s.Clients, remoteId)
	s.mu.Unlock()
	log.Infow("Connection deleted", "ID", remoteId)
	return nil
}
