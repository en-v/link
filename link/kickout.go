package link

import "github.com/pkg/errors"

func (self *Link) Kickout(clientId string) error {
	err := self.state.CloseClientConnection(clientId)
	if err != nil {
		return errors.Wrap(err, "Kickout")
	}
	return nil
}
