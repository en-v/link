package message

type Status int

const STATUS_UND Status = 0
const STATUS_OK Status = 1
const STATUS_ERR Status = 2

type Type int

const REQUEST Type = 1
const RESPONSE Type = 2

const ACTION_CLOSE_CONNECTION = "$CLOSE"
const ACTION_KEEPALIVE = "$KEEPALIVE"
const ACTION_HANDSHAKE = "$HANDSHAKE"

const HANDSHAKE_ACTIONS = "actions"
const HANDSHAKE_TOKEN = "token"
const ERROR_PAYLOAD_FIELD = "error"
