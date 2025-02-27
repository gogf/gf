// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gclient

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// WebSocketClient wraps the underlying websocket client connection
// and provides convenient functions.
//
// Deprecated: please use third-party library for websocket client instead.
type WebSocketClient struct {
	*websocket.Dialer
}

// NewWebSocket creates and returns a new WebSocketClient object.
//
// Deprecated: please use third-party library for websocket client instead.
func NewWebSocket() *WebSocketClient {
	return &WebSocketClient{
		&websocket.Dialer{
			Proxy:            http.ProxyFromEnvironment,
			HandshakeTimeout: 45 * time.Second,
		},
	}
}
