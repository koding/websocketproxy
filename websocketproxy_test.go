package websocketproxy

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

var (
	serverURL  = "ws://127.0.0.1:7777"
	backendURL = "ws://127.0.0.1:8888"
)

func TestProxy(t *testing.T) {
	// websocket proxy
	u, _ := url.Parse(backendURL)
	supportedSubProtocols := []string{"test-protocol"}
	upgrader := &websocket.Upgrader{
		ReadBufferSize:  4096,
		WriteBufferSize: 4096,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		Subprotocols: supportedSubProtocols,
	}

	proxy := NewProxy(u)
	proxy.Upgrader = upgrader

	go func() {
		mux := http.NewServeMux()
		mux.Handle("/proxy", proxy)
		if err := http.ListenAndServe(":7777", mux); err != nil {
			t.Fatal("ListenAndServe: ", err)
		}
	}()

	time.Sleep(time.Millisecond * 100)

	// backend echo server
	websocketMsgRcverCBackend := make(chan websocketMsg, 1)
	go func() {
		mux2 := http.NewServeMux()
		mux2.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			conn, err := upgrader.Upgrade(w, r, nil)
			if err != nil {
				log.Println(err)
				return
			}

			for {
				messageType, p, err := conn.ReadMessage()
				if err != nil {
					websocketMsgRcverCBackend <- websocketMsg{messageType, p, err}
					return
				}

				if err = conn.WriteMessage(messageType, p); err != nil {
					return
				}
			}
		})

		err := http.ListenAndServe(":8888", mux2)
		if err != nil {
			t.Fatal("ListenAndServe: ", err)
		}
	}()

	time.Sleep(time.Millisecond * 100)

	// define subprotocols for client, appending one not supported by the server
	clientSubProtocols := append(supportedSubProtocols, []string{"test-notsupported"}...)
	h := http.Header{}
	for _, subprot := range clientSubProtocols {
		h.Add("Sec-WebSocket-Protocol", subprot)
	}

	// frontend server, dial now our proxy, which will reverse proxy our
	// message to the backend websocket server.
	conn, resp, err := websocket.DefaultDialer.Dial(serverURL+"/proxy", h)
	if err != nil {
		t.Fatal(err)
	}

	// check if the server really accepted only the first one
	in := func(desired string) bool {
		for _, prot := range resp.Header[http.CanonicalHeaderKey("Sec-WebSocket-Protocol")] {
			if desired == prot {
				return true
			}
		}
		return false
	}

	if !in("test-protocol") {
		t.Error("test-protocol should be available")
	}

	if in("test-notsupported") {
		t.Error("test-notsupported should be not recevied from the server.")
	}

	// send msg to the backend server which goes through proxy
	msg := "hello kite"
	err = conn.WriteMessage(websocket.TextMessage, []byte(msg))
	if err != nil {
		t.Error(err)
	}

	messageType, p, err := conn.ReadMessage()
	if err != nil {
		t.Error(err)
	}

	if messageType != websocket.TextMessage {
		t.Error("incoming message type is not Text")
	}

	if msg != string(p) {
		t.Errorf("expecting: %s, got: %s", msg, string(p))
	}

	// shutdown procedure
	//
	proxy.Shutdown(context.Background())

	wsErrBackend := <-websocketMsgRcverCBackend
	e, ok := wsErrBackend.err.(*websocket.CloseError)
	if !ok {
		t.Fatal("backend error is not websocket.CloseError")
	}
	if e.Code != websocket.CloseGoingAway {
		t.Error("backend error code is not websocket.CloseGoingAway")
	}
	if e.Text != websocketProxyClosingMsg {
		t.Errorf("backend error test expecting: %s, got: %s", websocketProxyClosingMsg, e.Text)
	}
}
