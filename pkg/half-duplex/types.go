package halfduplexproxy

import (
	websocket "github.com/gorilla/websocket"
	common "github.com/koding/websocketproxy/pkg/common"
)

// HalfDuplexWebsocketProxy
type HalfDuplexWebsocketProxy struct {
	*common.WebsocketProxy

	ToBackend *websocket.Conn
	ToClient  *websocket.Conn
}
