// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// WebSocketClient wraps the underlying websocket client connection
// and provides convenient functions.
type WebSocketClient struct {
	*websocket.Dialer
}

// NewWebSocketClient New creates and returns a new WebSocketClient object.
func NewWebSocketClient() *WebSocketClient {
	return &WebSocketClient{
		&websocket.Dialer{
			Proxy:            http.ProxyFromEnvironment,
			HandshakeTimeout: 45 * time.Second,
		},
	}
}
