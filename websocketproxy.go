// Package websocketproxy is a reverse websocket proxy handler
package websocketproxy

import (
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
)

// WebsocketProxy is an HTTP Handler that takes an incoming websocket
// connection and proxies it to another server.
type WebsocketProxy struct {
	// Backend returns the backend URL which the proxy uses to reverse proxy
	// the incoming websocket connection.
	Backend func() *url.URL
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// ProxyHandler returns a new http.Handler interface that reverse proxies the
// request to the given target.
func ProxyHandler(target *url.URL) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		NewWebsocketProxy(target).ServerHTTP(rw, req)
	})
}

// NewWebsocketProxy returns a new Websocket ReverseProxy that rewrites the
// URL's to the scheme, host and base path provider in target.
func NewWebsocketProxy(target *url.URL) *WebsocketProxy {
	backend := func() *url.URL { return target }
	return &WebsocketProxy{Backend: backend}
}

func (w *WebsocketProxy) ServerHTTP(rw http.ResponseWriter, req *http.Request) {
	connPub, err := upgrader.Upgrade(rw, req, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer connPub.Close()

	backendURL := w.Backend()

	connKite, _, err := websocket.DefaultDialer.Dial(backendURL.String(), nil)
	if err != nil {
		log.Println("websocket.Dialer", err)
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
