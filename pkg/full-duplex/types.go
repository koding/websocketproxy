// Copyright 2022 The dvonthenen WebSocketProxy Authors. All Rights Reserved.
// Use of this source code is governed by an Apache-2.0
// license that can be found in the LICENSE file.
// SPDX-License-Identifier: Apache-2.0

package fullduplexproxy

import (
	websocket "github.com/dvonthenen/websocket"

	common "github.com/dvonthenen/websocketproxy/pkg/common"
)

// FullDuplexWebsocketProxy
type FullDuplexWebsocketProxy struct {
	*common.WebsocketProxy

	ToBackend *websocket.Conn
	ToClient  *websocket.Conn
}
