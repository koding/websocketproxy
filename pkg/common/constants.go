// Copyright 2022 The dvonthenen WebSocketProxy Authors. All Rights Reserved.
// Use of this source code is governed by an Apache-2.0
// license that can be found in the LICENSE file.
// SPDX-License-Identifier: Apache-2.0

package common

import (
	"errors"

	"github.com/dvonthenen/websocket"
)

const (
	MessageTypeUnknown int = iota
	MessageTypeControl
	MessageTypeData
)

var (
	// DefaultUpgrader specifies the parameters for upgrading an HTTP
	// connection to a WebSocket connection.
	DefaultUpgrader = &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	// DefaultDialer is a dialer with all fields set to the default zero values.
	DefaultDialer = websocket.DefaultDialer

	// ErrConnectionNotEstablished connection not established
	ErrConnectionNotEstablished = errors.New("connection not established")
)
