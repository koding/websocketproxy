package common

import (
	"errors"

	"github.com/gorilla/websocket"
)

const (
	MessageTypeUnknown int = iota
	MessageTypeControl
	MessageTypeData
)

var (
	// DefaultUpgrader specifies the parameters for upgrading an HTTP
	// connection to a WebSocket connection.
	DefaultUpgrader = &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	// DefaultDialer is a dialer with all fields set to the default zero values.
	DefaultDialer = websocket.DefaultDialer

	// ErrConnectionNotEstablished connection not established
	ErrConnectionNotEstablished = errors.New("connection not established")
)
