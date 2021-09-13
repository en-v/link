package types

import (
	"encoding/json"

	"github.com/pkg/errors"
)

type Payload struct {
	senderId string
	items    map[string]interface{}
	rawjson  json.RawMessage
}

func Newp() *Payload {
	return &Payload{
		items: make(map[string]interface{}),
	}
}

func NewPut(pointer interface{}) (*Payload, error) {
	n := Newp()
	err := n.Put(pointer)
	if err != nil {
		return nil, err
	}
	return n, nil
}

func (self *Payload) Set(field string, data interface{}) {
	self.items[field] = data
}

func (self *Payload) Sert(field string, data interface{}) *Payload {
	self.items[field] = data
	return self
}

func (self *Payload) Get(field string) (interface{}, error) {
	v, e := self.items[field]
	if e {
		return v, nil
	}
	return "", errors.New("Field doesnt exist, " + field)
}

// MAP

func (self *Payload) SetMap(items map[string]interface{}) {
	self.items = items
}

func (self *Payload) GetAsMap() map[string]interface{} {
	return self.items
}

// TYPED

func (self *Payload) GetString(field string) (string, error) {
	v, e := self.Get(field)
	if e != nil {
		return "", e

	}
	return v.(string), nil
}

func (self *Payload) GetFloat(field string) (float64, error) {
	v, e := self.Get(field)
	if e != nil {
		return 0, e

	}
	r, success := v.(float64)
	if !success {
		return 0, errors.New("Type cast error, value is not float64, " + field)
	}
	return r, nil
}

func (self *Payload) GetInt32(field string) (int32, error) {
	v, e := self.Get(field)
	if e != nil {
		return 0, e

	}
	r, success := v.(float64)
	if !success {
		return 0, errors.New("Type cast error, value is not float64, " + field)
	}
	return int32(r), nil
}

// SENDER ID

func (self *Payload) SetSenderId(semderId string) {
	self.senderId = semderId
}

func (self *Payload) SenderId() string {
	return self.senderId
}

// JSON

func (self *Payload) UnmarshalJSON(data []byte) error {
	self.rawjson = data
	items := make(map[string]interface{})

	err := json.Unmarshal(data, &items)
	if err != nil {
		return err
	}

	self.items = items
	return nil
}

func (self *Payload) MarshalJSON() ([]byte, error) {
	return json.Marshal(self.items)
}

func (self *Payload) Raw() json.RawMessage {
	return self.rawjson
}

// PULL AND PUT
// setting to object by pointer from payload data, get
func (self *Payload) Pull(pointer interface{}) error {

	if self.rawjson == nil {
		return errors.New("Raw JSON doesnt exist")
	}

	err := json.Unmarshal(self.rawjson, pointer)
	if err != nil {
		return errors.Wrap(err, "Payload.UnmarshallTo")
	}

	return nil
}

// setting data to payload from object by pointer, set
func (self *Payload) Put(pointer interface{}) error {

	rawjson, err := json.Marshal(pointer)
	if err != nil {
		return errors.Wrap(err, "Payload.MarshallToRaw")
	}

	err = self.UnmarshalJSON(rawjson)
	if err != nil {
		return errors.Wrap(err, "Payload.UnmarshalJSON")
	}

	return nil
}
