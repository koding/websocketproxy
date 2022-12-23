package fullduplexproxy

import (
	websocket "github.com/gorilla/websocket"
	common "github.com/koding/websocketproxy/pkg/common"
)

// FullDuplexWebsocketProxy
type FullDuplexWebsocketProxy struct {
	*common.WebsocketProxy

	ToBackend *websocket.Conn
	ToClient  *websocket.Conn
}
