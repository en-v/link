package types

type LinkProxy interface {
	GetLinkHandlers() *Handlers
}

type Handlers struct {
	ID    string
	Auth  ClientsAuthFunc
	Reg   ClientsRegFunc
	Unreg ClientsUnregFunc
}

type ClientsRegFunc func(ClientId string) error
type ClientsUnregFunc func(ClientId string) error
type ClientsAuthFunc func(ClientId string, Token string) error

/*
type Payload map[string]interface{}

func Pl() Payload {
	return Payload{}
}

func (self Payload) Set(field string, data interface{}) Payload {
	self[field] = data
	return self
}

const LINK_REM_ID_PROP_NAME = "$#6@8*%*8@6#$"

func (self Payload) SetRemoteId(remId string) {
	self[LINK_REM_ID_PROP_NAME] = remId
}

func (self Payload) GetRemoteId() string {
	remId, ok := self[LINK_REM_ID_PROP_NAME]
	if ok {
		return remId.(string)
	}
	return ""
}
*/
