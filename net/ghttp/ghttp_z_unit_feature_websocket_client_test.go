// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gorilla/websocket"
)

func Test_WebSocketClient(t *testing.T) {
	p, _ := ports.PopRand()
	s := g.Server(p)
	s.BindHandler("/ws", func(r *ghttp.Request) {
		ws, err := r.WebSocket()
		if err != nil {
			r.Exit()
		}
		for {
			msgType, msg, err := ws.ReadMessage()
			if err != nil {
				return
			}
			if err = ws.WriteMessage(msgType, msg); err != nil {
				return
			}
		}
	})
	s.SetPort(p)
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		client := ghttp.NewWebSocketClient()
		client.Proxy = http.ProxyFromEnvironment
		client.HandshakeTimeout = time.Minute

		conn, _, err := client.Dial(fmt.Sprintf("ws://127.0.0.1:%d/ws", p), nil)
		t.Assert(err, nil)
		defer conn.Close()

		msg := []byte("hello")
		err = conn.WriteMessage(websocket.TextMessage, msg)
		t.Assert(err, nil)

		mt, data, err := conn.ReadMessage()
		t.Assert(err, nil)
		t.Assert(mt, websocket.TextMessage)
		t.Assert(data, msg)
	})
}
