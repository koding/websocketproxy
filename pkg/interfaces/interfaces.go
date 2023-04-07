// Copyright 2022 The dvonthenen WebSocketProxy Authors. All Rights Reserved.
// Use of this source code is governed by an Apache-2.0
// license that can be found in the LICENSE file.
// SPDX-License-Identifier: Apache-2.0

// Package websocketproxy is a reverse proxy for WebSocket connections.
package interfaces

import (
	"net/http"
)

// ManageCallback is a callback to manage connections
type ManageCallback interface {
	RemoveConnection(uniqueId string)
}

// MessageCallback is a callback to view messages as they passthrough the proxy
type MessageCallback interface {
	HandleMessage(byMsg []byte) error
}

// DirectorCallback is a callback to modify the header before they passthrough the proxy
type DirectorCallback interface {
	AdjustHeaders(incoming *http.Request, out http.Header)
}
