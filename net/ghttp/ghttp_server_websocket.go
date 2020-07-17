// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
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
	// TextMessage denotes a text data message. The text message payload is
	// interpreted as UTF-8 encoded text data.
	WS_MSG_TEXT = websocket.TextMessage

	// BinaryMessage denotes a binary data message.
	WS_MSG_BINARY = websocket.BinaryMessage

	// CloseMessage denotes a close control message. The optional message
	// payload contains a numeric code and text. Use the FormatCloseMessage
	// function to format a close message payload.
	WS_MSG_CLOSE = websocket.CloseMessage

	// PingMessage denotes a ping control message. The optional message payload
	// is UTF-8 encoded text.
	WS_MSG_PING = websocket.PingMessage

	// PongMessage denotes a pong control message. The optional message payload
	// is UTF-8 encoded text.
	WS_MSG_PONG = websocket.PongMessage
)
