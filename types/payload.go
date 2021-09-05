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

func (this *Payload) Set(field string, data interface{}) {
	this.items[field] = data
}

func (this *Payload) Sert(field string, data interface{}) *Payload {
	this.items[field] = data
	return this
}

func (this *Payload) Get(field string) (interface{}, error) {
	v, e := this.items[field]
	if e {
		return v, nil
	}
	return "", errors.New("Field doesnt exist, " + field)
}

// MAP

func (this *Payload) SetMap(items map[string]interface{}) {
	this.items = items
}

func (this *Payload) GetAsMap() map[string]interface{} {
	return this.items
}

// TYPED

func (this *Payload) GetString(field string) (string, error) {
	v, e := this.Get(field)
	if e != nil {
		return "", e

	}
	return v.(string), nil
}

func (this *Payload) GetFloat(field string) (float64, error) {
	v, e := this.Get(field)
	if e != nil {
		return 0, e

	}
	r, success := v.(float64)
	if !success {
		return 0, errors.New("Type cast error, value is not float64, " + field)
	}
	return r, nil
}

func (this *Payload) GetInt32(field string) (int32, error) {
	v, e := this.Get(field)
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

func (this *Payload) SetSenderId(semderId string) {
	this.senderId = semderId
}

func (this *Payload) SenderId() string {
	return this.senderId
}

// JSON

func (this *Payload) UnmarshalJSON(data []byte) error {
	this.rawjson = data
	items := make(map[string]interface{})

	err := json.Unmarshal(data, &items)
	if err != nil {
		return err
	}

	this.items = items
	return nil
}

func (this *Payload) MarshalJSON() ([]byte, error) {
	return json.Marshal(this.items)
}

func (this *Payload) Raw() json.RawMessage {
	return this.rawjson
}

// PULL AND PUT
// setting to object by pointer from payload data, get
func (this *Payload) Pull(pointer interface{}) error {

	if this.rawjson == nil {
		return errors.New("Raw JSON doesnt exist")
	}

	err := json.Unmarshal(this.rawjson, pointer)
	if err != nil {
		return errors.Wrap(err, "Payload.UnmarshallTo")
	}

	return nil
}

// setting data to payload from object by pointer, set
func (this *Payload) Put(pointer interface{}) error {

	rawjson, err := json.Marshal(pointer)
	if err != nil {
		return errors.Wrap(err, "Payload.MarshallToRaw")
	}

	err = this.UnmarshalJSON(rawjson)
	if err != nil {
		return errors.Wrap(err, "Payload.UnmarshalJSON")
	}

	return nil
}
