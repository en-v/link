package types

type Handler interface {
	Hooks() *Hooks
}

type Hooks struct {
	LocalId  GetIdFunc
	Verify   VefifyFunc
	CheckIn  CheckInFunc
	CheckOut CheckOutFunc
}

type GetIdFunc func() string
type CheckInFunc func(remoteId string) error
type CheckOutFunc func(remoteId string) error
type VefifyFunc func(remoteId string, token string) error
