package common

import (
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"

	"github.com/koding/websocketproxy/pkg/interfaces"
)

// ProxyOptions these are the available options for a Proxy
type ProxyOptions struct {
	Url           *url.URL
	NaturalTunnel bool
	Upgrader      *websocket.Upgrader
	Dialer        *websocket.Dialer

	Director interfaces.DirectorCallback
	Viewer   interfaces.MessageCallback
}

// WebsocketProxy is an HTTP Handler that takes an incoming WebSocket
// connection and proxies it to another server.
type WebsocketProxy struct {
	// Director, if non-nil, is a function that may copy additional request
	// headers from the incoming WebSocket connection into the output headers
	// which will be forwarded to another server.
	Director interfaces.DirectorCallback

	// Viewer, if non-nil, is a function that may view messages as they comeback
	// to the originating client
	Viewer interfaces.MessageCallback

	// Upgrader specifies the parameters for upgrading a incoming HTTP
	// connection to a WebSocket connection. If nil, DefaultUpgrader is used.
	Upgrader *websocket.Upgrader

	//  Dialer contains options for connecting to the backend WebSocket server.
	//  If nil, DefaultDialer is used.
	Dialer *websocket.Dialer

	// ProxyOptions describe how to initialize the Proxy
	Options ProxyOptions

	// Backend returns the backend URL which the proxy uses to reverse proxy
	// the incoming WebSocket connection. Request is the initial incoming and
	// unmodified request.
	Backend func(*http.Request) *url.URL

	// Stop channels to close the websocket on demand
	StopClientChan  chan struct{}
	StopBackendChan chan struct{}
}
