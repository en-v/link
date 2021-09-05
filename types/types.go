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

func (this Payload) Set(field string, data interface{}) Payload {
	this[field] = data
	return this
}

const LINK_REM_ID_PROP_NAME = "$#6@8*%*8@6#$"

func (this Payload) SetRemoteId(remId string) {
	this[LINK_REM_ID_PROP_NAME] = remId
}

func (this Payload) GetRemoteId() string {
	remId, ok := this[LINK_REM_ID_PROP_NAME]
	if ok {
		return remId.(string)
	}
	return ""
}
*/
