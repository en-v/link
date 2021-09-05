package link

import (
	"errors"

	"github.com/en-v/link/core"
)

func (this *Link) tlsFilesExist(certFile string, keyFile string) (bool, error) {
	if certFile != "" && keyFile != "" {

		if !FilexExist(certFile) || !FilexExist(keyFile) {
			return false, nil
		}

		return true, nil
	}
	return false, nil
}

func (this *Link) verifyLinkAsServer() error {

	if this.state.Mode == core.CLIENT_MODE {
		return errors.New("Link is alredy in client mode")
	}

	for i := range this.state.Connections {
		if this.state.Connections[i] != nil && this.state.Connections[i].IsActive() {
			return errors.New("More than 1 active or not nil connections found")
		}
	}

	return nil
}

func (this *Link) verifyLinkAsCaller() error {
	if this.state.Mode == core.GATE_MODE {
		return errors.New("Link is alredy in server mode")
	}

	if len(this.state.Connections) > 1 {
		return errors.New("More than 1 active connections found")
	}

	if len(this.state.Connections) == 1 && this.state.Connections[0] != nil && this.state.Connections[0].IsActive() {
		return errors.New("Local active client connections exists")
	}

	return nil
}
