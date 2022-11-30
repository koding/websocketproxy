// Package websocketproxy is a reverse proxy for WebSocket connections.
package interfaces

import (
	"net/http"
)

// MessageCallback is a callback to view messages as they passthrough the proxy
type MessageCallback interface {
	HandleMessage(byMsg []byte) error
}

// DirectorCallback is a callback to modify the header before they passthrough the proxy
type DirectorCallback interface {
	AdjustHeaders(incoming *http.Request, out http.Header)
}
