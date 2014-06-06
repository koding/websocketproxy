// Package websocketproxy is a reverse proxy for WebSocket connections.
package websocketproxy

import (
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
)

var (
	// DefaultUpgrader specifies the paramaters for upgrading an HTTP connection to
	// a WebSocket connection.
	DefaultUpgrader = &websocket.Upgrader{
		ReadBufferSize:  4096,
		WriteBufferSize: 4096,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	// DefaultDialer is a dialer with all fields set to the default zero values.
	DefaultDialer = websocket.DefaultDialer
)

// WebsocketProxy is an HTTP Handler that takes an incoming websocket
// connection and proxies it to another server.
type WebsocketProxy struct {
	// Backend returns the backend URL which the proxy uses to reverse proxy
	// the incoming websocket connection.
	Backend func() *url.URL

	// Upgrader specifies the paramaters for upgrading an HTTP connection to a
	// WebSocket connection. If nil, DefaultUpgrader is used.
	Upgrader *websocket.Upgrader

	//  Dialer contains options for connecting to WebSocket server.
	Dialer *websocket.Dialer
}

// ProxyHandler returns a new http.Handler interface that reverse proxies the
// request to the given target.
func ProxyHandler(target *url.URL) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		NewProxy(target).ServerHTTP(rw, req)
	})
}

// NewProxy returns a new Websocket reverse proxy that rewrites the
// URL's to the scheme, host and base path provider in target.
func NewProxy(target *url.URL) *WebsocketProxy {
	backend := func() *url.URL { return target }
	return &WebsocketProxy{Backend: backend}
}

// ServerHTTP implements the http.Handler that proxies WebSocket connections.
func (w *WebsocketProxy) ServerHTTP(rw http.ResponseWriter, req *http.Request) {
	upgrader := w.Upgrader
	if w.Upgrader == nil {
		upgrader = DefaultUpgrader
	}

	connPub, err := upgrader.Upgrade(rw, req, nil)
	if err != nil {
		log.Println("websocketproxy: couldn't upgrade %s", err)
		return
	}
	defer connPub.Close()

	backendURL := w.Backend()

	dialer := w.Dialer
	if w.Dialer == nil {
		dialer = DefaultDialer
	}

	connKite, _, err := dialer.Dial(backendURL.String(), nil)
	if err != nil {
		log.Println("websocketproxy: couldn't dial to remote backend url %s", err)
		return
	}
	defer connKite.Close()

	errc := make(chan error, 2)
	cp := func(dst io.Writer, src io.Reader) {
		_, err := io.Copy(dst, src)
		errc <- err
	}

	go cp(connKite.UnderlyingConn(), connPub.UnderlyingConn())
	go cp(connPub.UnderlyingConn(), connKite.UnderlyingConn())
	<-errc
}
