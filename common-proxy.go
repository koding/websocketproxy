// Package websocketproxy is a reverse proxy for WebSocket connections.
package websocketproxy

import (
	"io"
	"net/http"
	"net/url"
)

// ProxyHandler returns a new http.Handler interface that reverse proxies the
// request to the given target.
func ProxyHandler(options ProxyOptions) http.Handler { return NewProxy(options) }

// NewProxy returns a new Websocket reverse proxy that rewrites the
// URL's to the scheme, host and base path provider in target.
func NewProxy(options ProxyOptions) *WebsocketProxy {
	backend := func(r *http.Request) *url.URL {
		// Shallow copy
		u := options.Url
		u.Fragment = r.URL.Fragment
		u.Path = r.URL.Path
		u.RawQuery = r.URL.RawQuery
		return u
	}
	return &WebsocketProxy{
		Director: options.Director,
		Viewer:   options.Viewer,
		Upgrader: options.Upgrader,
		Dialer:   options.Dialer,
		backend:  backend,
		options:  options,
	}
}

// Stop websocket proxy on demand
func (w *WebsocketProxy) CloseProxy() {
	close(w.stopBackendChan)
	close(w.stopClientChan)
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

func copyResponse(rw http.ResponseWriter, resp *http.Response) error {
	copyHeader(rw.Header(), resp.Header)
	rw.WriteHeader(resp.StatusCode)
	defer resp.Body.Close()

	_, err := io.Copy(rw, resp.Body)
	return err
}
