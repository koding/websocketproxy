module github.com/koding/websocketproxy

go 1.18

require (
	github.com/dvonthenen/websocket v1.5.1-dyv.2
	k8s.io/klog/v2 v2.90.0
)

require github.com/go-logr/logr v1.2.0 // indirect

// replace github.com/gorilla/websocket => ../../gorilla/websocket
