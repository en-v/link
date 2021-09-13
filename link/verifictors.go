package link

import (
	"errors"

	"github.com/en-v/link/core"
)

func (self *Link) tlsFilesExist(certFile string, keyFile string) (bool, error) {
	if certFile != "" && keyFile != "" {

		if !FilexExist(certFile) || !FilexExist(keyFile) {
			return false, nil
		}

		return true, nil
	}
	return false, nil
}

func (self *Link) verifyLinkAsServer() error {

	if self.state.Mode == core.CLIENT_MODE {
		return errors.New("Link is alredy in client mode")
	}

	for i := range self.state.Connections {
		if self.state.Connections[i] != nil && self.state.Connections[i].IsActive() {
			return errors.New("More than 1 active or not nil connections found")
		}
	}

	return nil
}

func (self *Link) verifyLinkAsCaller() error {
	if self.state.Mode == core.GATE_MODE {
		return errors.New("Link is alredy in server mode")
	}

	if len(self.state.Connections) > 1 {
		return errors.New("More than 1 active connections found")
	}

	if len(self.state.Connections) == 1 && self.state.Connections[0] != nil && self.state.Connections[0].IsActive() {
		return errors.New("Local active client connections exists")
	}

	return nil
}
