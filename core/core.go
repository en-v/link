package core

import "time"

const VERSION = "0.1"

type LinkMode byte

const BLANK_MODE LinkMode = 0
const GATE_MODE LinkMode = 1
const CLIENT_MODE LinkMode = 2

const MAX_WAITERS_PER_CONNECTION_COUNT = 128
const MAX_CONNECTIONS_PER_GATE_COUNT = 1024
const IDENTIFIERS_BYTES_COUNT = 8

const RESPONSE_WAITING_EXPIRE_INTERVAL = time.Minute * 2
const KEEPALIVE_INTERVAL = time.Second * 10
const ZOOMBIE_HUNTER_INERVAL = KEEPALIVE_INTERVAL * 3

var DEBUG = false

func (this LinkMode) ToString() string {
	switch this {
	case GATE_MODE:
		return "GATE"

	case BLANK_MODE:
		return "BLANK"

	case CLIENT_MODE:
		return "CLIENT"

	default:
		return "UNKNOWN"
	}
}
