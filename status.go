package junipero

import "errors"

//
type MessageType int

const (
	// TextMessage denotes a text data message. The text message payload is
	// interpreted as UTF-8 encoded text data.
	TextMessage MessageType = 1

	// BinaryMessage denotes a binary data message.
	BinaryMessage MessageType = 2

	// CloseMessage denotes a close control message. The optional message
	// payload contains a numeric code and text. Use the FormatCloseMessage
	// function to format a close message payload.
	CloseMessage MessageType = 8

	// PingMessage denotes a ping control message. The optional message payload
	// is UTF-8 encoded text.
	PingMessage MessageType = 9

	// PongMessage denotes a pong control message. The optional message payload
	// is UTF-8 encoded text.
	PongMessage MessageType = 10
)

// CloseStatus 是連線被關閉時的狀態代號。
type CloseStatus int

const (
	CloseNormalClosure           CloseStatus = 1000
	CloseGoingAway               CloseStatus = 1001
	CloseProtocolError           CloseStatus = 1002
	CloseUnsupportedData         CloseStatus = 1003
	CloseNoStatusReceived        CloseStatus = 1005
	CloseAbnormalClosure         CloseStatus = 1006
	CloseInvalidFramePayloadData CloseStatus = 1007
	ClosePolicyViolation         CloseStatus = 1008
	CloseMessageTooBig           CloseStatus = 1009
	CloseMandatoryExtension      CloseStatus = 1010
	CloseInternalServerErr       CloseStatus = 1011
	CloseServiceRestart          CloseStatus = 1012
	CloseTryAgainLater           CloseStatus = 1013
	CloseTLSHandshake            CloseStatus = 1015
)

var (
	ErrEngineClosed         = errors.New("junipero: upgrading connections when engine closed")
	ErrChannelClosed        = errors.New("junipero: interacting with a closed channel")
	ErrSessionTimedOut      = errors.New("junipero: interacting with a timed out session")
	ErrConnectionClosed     = errors.New("junipero: interacting with a disconnected connection")
	ErrSessionClosed        = errors.New("junipero: interacting with a closed session")
	ErrChannelNotFound      = errors.New("junipero: interacting with a undefined channel")
	ErrChannelSubscribed    = errors.New("junipero: subscribing to a subscribed channel")
	ErrChannelNotSubscribed = errors.New("junipero: unsubscribing a unsubscribed channel")
	ErrKeyNotFound          = errors.New("junipero: accessing a undefined key from the session store")
	ErrWriteTimedOut        = errors.New("junipero: write timed out")
)
