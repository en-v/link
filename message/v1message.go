package message

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"time"

	"github.com/en-v/link/core"
	"github.com/en-v/link/types"
)

const V1 = 1

type MsgV1 struct {
	Version   int            `json:"ver"` // link-protocol version, =1
	MsgId     string         `json:"mid"` // messages pair ID (request and response - both have one common MessageID)
	SenderId  string         `json:"sid"` // sender ID (16 symbold ID string - 8 random bytes in hex-view)
	RecvrId   string         `json:"rid"` // receiver ID (16 symbold ID string - 8 random bytes in hex-view)
	Type      Type           `json:"tpe"` // message type: request (=1) or  response (=2)
	Action    string         `json:"act"` // invocated action name
	Status    Status         `json:"sta"` // message status: OK (=1), ERR (=2), unknown
	Payload   *types.Payload `json:"pld"` // message payload: any map[string]interface{} data
	Requested time.Time      `json:"qtm"` // reqesting time: time of request message creation on the sender side
	Responsed time.Time      `json:"rtm"` // responding time: time of response message creation on the receiver side
}

func EmptyV1() *MsgV1 {
	return &MsgV1{
		Version: V1,
		Status:  STATUS_UND,
	}
}

func randId() string {
	rnd := make([]byte, core.IDENTIFIERS_BYTES_COUNT)
	rand.Read(rnd)
	return hex.EncodeToString(rnd)
}

func Request(senderId string, action string, payload *types.Payload) *MsgV1 {
	return &MsgV1{
		Version:   V1,
		MsgId:     randId(),
		Status:    STATUS_UND,
		SenderId:  senderId,
		Type:      REQUEST,
		Payload:   payload,
		Action:    action,
		Requested: time.Now(),
	}
}

func (this *MsgV1) ToResponse(receiverId string, payload *types.Payload) *MsgV1 {
	return &MsgV1{
		Version:   this.Version,
		MsgId:     this.MsgId,
		SenderId:  this.SenderId,
		RecvrId:   receiverId,
		Action:    this.Action,
		Type:      RESPONSE,
		Status:    STATUS_OK,
		Responsed: time.Now(),
		Requested: this.Requested,
		Payload:   payload,
	}
}

func (this *MsgV1) ToError(receiverId string, err error) *MsgV1 {
	response := this.ToResponse(receiverId, nil)
	response.Status = STATUS_ERR
	response.Payload = types.Newp()
	response.Payload.Set(ERROR_PAYLOAD_FIELD, err.Error())
	return response
}

func (this *MsgV1) ToJson() string {
	d, e := json.Marshal(this)
	if e != nil {
		return e.Error()
	}
	return string(d)
}
