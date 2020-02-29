// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/gorilla/websocket"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/test/gtest"
)

func Test_WebSocket(t *testing.T) {
	p := ports.PopRand()
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
	gtest.Case(t, func() {
		conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://127.0.0.1:%d/ws", p), nil)
		gtest.Assert(err, nil)
		defer conn.Close()

		msg := []byte("hello")
		err = conn.WriteMessage(websocket.TextMessage, msg)
		gtest.Assert(err, nil)

		t, data, err := conn.ReadMessage()
		gtest.Assert(err, nil)
		gtest.Assert(t, websocket.TextMessage)
		gtest.Assert(data, msg)
	})
}
