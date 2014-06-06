package websocketproxy

import (
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

func ProxyFunc(w http.ResponseWriter, r *http.Request) {
	u, _ := url.Parse("ws://127.0.0.1:8888")
	ProxyHandler(u).ServeHTTP(w, r)
}

func TestProxy(t *testing.T) {
	// websocket proxy
	mux := http.NewServeMux()
	mux.HandleFunc("/proxy", ProxyFunc)
	go func() {
		if err := http.ListenAndServe(":7777", mux); err != nil {
			t.Fatal("ListenAndServe: ", err)
		}
	}()

	time.Sleep(time.Millisecond * 100)

	// backend echo server
	go func() {
		mux2 := http.NewServeMux()
		mux2.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			conn, err := DefaultUpgrader.Upgrade(w, r, nil)
			if err != nil {
				log.Println(err)
				return
			}

			messageType, p, err := conn.ReadMessage()
			if err != nil {
				return
			}

			if err = conn.WriteMessage(messageType, p); err != nil {
				return
			}
		})

		err := http.ListenAndServe(":8888", mux2)
		if err != nil {
			t.Fatal("ListenAndServe: ", err)
		}
	}()

	time.Sleep(time.Millisecond * 100)

	// frontend server, dial now our proxy, which will reverse proxy our
	// message to the backend websocket server.
	conn, _, err := websocket.DefaultDialer.Dial(serverURL+"/proxy", nil)
	if err != nil {
		t.Fatal(err)
	}

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
}
