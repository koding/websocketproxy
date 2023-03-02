package halfduplexproxy

import (
	websocket "github.com/dvonthenen/websocket"

	common "github.com/dvonthenen/websocketproxy/pkg/common"
)

// HalfDuplexWebsocketProxy
type HalfDuplexWebsocketProxy struct {
	*common.WebsocketProxy

	ToBackend *websocket.Conn
	ToClient  *websocket.Conn
}
