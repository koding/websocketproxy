module github.com/koding/websocketproxy

go 1.18

require (
	github.com/gorilla/websocket v1.5.0
	k8s.io/klog/v2 v2.90.0
)

require github.com/go-logr/logr v1.2.0 // indirect

replace github.com/gorilla/websocket => github.com/dvonthenen/websocket v1.5.1-0.20230208185225-642cd054e185

// replace github.com/gorilla/websocket => ../../gorilla/websocket
