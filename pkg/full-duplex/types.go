package fullduplexproxy

import (
	websocket "github.com/dvonthenen/websocket"

	common "github.com/dvonthenen/websocketproxy/pkg/common"
)

// FullDuplexWebsocketProxy
type FullDuplexWebsocketProxy struct {
	*common.WebsocketProxy

	ToBackend *websocket.Conn
	ToClient  *websocket.Conn
}
