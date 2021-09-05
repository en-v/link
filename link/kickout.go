package link

import "github.com/pkg/errors"

func (this *Link) Kickout(clientId string) error {
	err := this.state.CloseClientConnection(clientId)
	if err != nil {
		return errors.Wrap(err, "Kickout")
	}
	return nil
}
