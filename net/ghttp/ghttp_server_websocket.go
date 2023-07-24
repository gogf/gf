// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import "github.com/gorilla/websocket"

// WebSocket wraps the underlying websocket connection
// and provides convenient functions.
type WebSocket struct {
	*websocket.Conn
}

const (
	// WsMsgText TextMessage denotes a text data message.
	// The text message payload is interpreted as UTF-8 encoded text data.
	WsMsgText = websocket.TextMessage

	// WsMsgBinary BinaryMessage denotes a binary data message.
	WsMsgBinary = websocket.BinaryMessage

	// WsMsgClose CloseMessage denotes a close control message.
	// The optional message payload contains a numeric code and text.
	// Use the FormatCloseMessage function to format a close message payload.
	WsMsgClose = websocket.CloseMessage

	// WsMsgPing PingMessage denotes a ping control message.
	// The optional message payload is UTF-8 encoded text.
	WsMsgPing = websocket.PingMessage

	// WsMsgPong PongMessage denotes a pong control message.
	// The optional message payload is UTF-8 encoded text.
	WsMsgPong = websocket.PongMessage
)
