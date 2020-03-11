// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gudp_test

import (
	"fmt"
	"github.com/gogf/gf/net/gudp"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/test/gtest"
	"github.com/gogf/gf/util/gconv"
	"testing"
	"time"
)

func Test_Basic(t *testing.T) {
	p := ports.PopRand()
	s := gudp.NewServer(fmt.Sprintf("127.0.0.1:%d", p), func(conn *gudp.Conn) {
		defer conn.Close()
		for {
			data, err := conn.Recv(-1)
			if len(data) > 0 {
				if err := conn.Send(append([]byte("> "), data...)); err != nil {
					glog.Error(err)
				}
			}
			if err != nil {
				break
			}
		}
	})
	go s.Run()
	defer s.Close()
	time.Sleep(100 * time.Millisecond)
	// gudp.Conn.Send
	gtest.Case(t, func() {
		for i := 0; i < 100; i++ {
			conn, err := gudp.NewConn(fmt.Sprintf("127.0.0.1:%d", p))
			gtest.Assert(err, nil)
			gtest.Assert(conn.Send([]byte(gconv.String(i))), nil)
			conn.Close()
		}
	})
	// gudp.Conn.SendRecv
	gtest.Case(t, func() {
		for i := 0; i < 100; i++ {
			conn, err := gudp.NewConn(fmt.Sprintf("127.0.0.1:%d", p))
			gtest.Assert(err, nil)
			result, err := conn.SendRecv([]byte(gconv.String(i)), -1)
			gtest.Assert(err, nil)
			gtest.Assert(string(result), fmt.Sprintf(`> %d`, i))
			conn.Close()
		}
	})
	// gudp.Send
	gtest.Case(t, func() {
		for i := 0; i < 100; i++ {
			err := gudp.Send(fmt.Sprintf("127.0.0.1:%d", p), []byte(gconv.String(i)))
			gtest.Assert(err, nil)
		}
	})
	// gudp.SendRecv
	gtest.Case(t, func() {
		for i := 0; i < 100; i++ {
			result, err := gudp.SendRecv(fmt.Sprintf("127.0.0.1:%d", p), []byte(gconv.String(i)), -1)
			gtest.Assert(err, nil)
			gtest.Assert(string(result), fmt.Sprintf(`> %d`, i))
		}
	})
}

// If the read buffer size is less than the sent package size,
// the rest data would be dropped.
func Test_Buffer(t *testing.T) {
	p := ports.PopRand()
	s := gudp.NewServer(fmt.Sprintf("127.0.0.1:%d", p), func(conn *gudp.Conn) {
		defer conn.Close()
		for {
			data, err := conn.Recv(1)
			if len(data) > 0 {
				if err := conn.Send(data); err != nil {
					glog.Error(err)
				}
			}
			if err != nil {
				break
			}
		}
	})
	go s.Run()
	defer s.Close()
	time.Sleep(100 * time.Millisecond)
	gtest.Case(t, func() {
		result, err := gudp.SendRecv(fmt.Sprintf("127.0.0.1:%d", p), []byte("123"), -1)
		gtest.Assert(err, nil)
		gtest.Assert(string(result), "1")
	})
	gtest.Case(t, func() {
		result, err := gudp.SendRecv(fmt.Sprintf("127.0.0.1:%d", p), []byte("456"), -1)
		gtest.Assert(err, nil)
		gtest.Assert(string(result), "4")
	})
}
