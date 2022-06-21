// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gudp_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/gogf/gf/v2/net/gudp"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gconv"
)

var (
	simpleTimeout = time.Millisecond * 100
	sendData      = []byte("hello")
	invalidAddr   = "127.0.0.1:99999"
)

func startUDPServer(addr string) {
	s := gudp.NewServer(addr, func(conn *gudp.Conn) {
		defer conn.Close()
		for {
			data, err := conn.Recv(-1)
			if err != nil {
				break
			}
			conn.Send(data)
		}
	})
	go s.Run()
	time.Sleep(simpleTimeout)
}

func Test_Basic(t *testing.T) {
	var (
		ctx = context.TODO()
	)
	p, _ := gudp.GetFreePort()
	s := gudp.NewServer(fmt.Sprintf("127.0.0.1:%d", p), func(conn *gudp.Conn) {
		defer conn.Close()
		for {
			data, err := conn.Recv(-1)
			if len(data) > 0 {
				if err := conn.Send(append([]byte("> "), data...)); err != nil {
					glog.Error(ctx, err)
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
	gtest.C(t, func(t *gtest.T) {
		for i := 0; i < 100; i++ {
			conn, err := gudp.NewConn(fmt.Sprintf("127.0.0.1:%d", p))
			t.AssertNil(err)
			t.Assert(conn.Send([]byte(gconv.String(i))), nil)
			t.AssertNil(conn.RemoteAddr())
			result, err := conn.Recv(-1)
			t.AssertNil(err)
			t.AssertNE(conn.RemoteAddr(), nil)
			t.Assert(string(result), fmt.Sprintf(`> %d`, i))
			conn.Close()
		}
	})
	// gudp.Conn.SendRecv
	gtest.C(t, func(t *gtest.T) {
		for i := 0; i < 100; i++ {
			conn, err := gudp.NewConn(fmt.Sprintf("127.0.0.1:%d", p))
			t.AssertNil(err)
			_, err = conn.SendRecv([]byte(gconv.String(i)), -1)
			t.AssertNil(err)
			//t.Assert(string(result), fmt.Sprintf(`> %d`, i))
			conn.Close()
		}
	})
	// gudp.Conn.SendWithTimeout
	gtest.C(t, func(t *gtest.T) {
		for i := 0; i < 100; i++ {
			conn, err := gudp.NewConn(fmt.Sprintf("127.0.0.1:%d", p))
			t.AssertNil(err)
			err = conn.SendWithTimeout([]byte(gconv.String(i)), time.Second)
			t.AssertNil(err)
			conn.Close()
		}
	})
	// gudp.Conn.RecvWithTimeout
	gtest.C(t, func(t *gtest.T) {
		for i := 0; i < 100; i++ {
			conn, err := gudp.NewConn(fmt.Sprintf("127.0.0.1:%d", p))
			t.AssertNil(err)
			err = conn.Send([]byte(gconv.String(i)))
			t.AssertNil(err)
			conn.SetRecvBufferWait(time.Millisecond * 100)
			result, err := conn.RecvWithTimeout(-1, time.Second)
			t.AssertNil(err)
			t.Assert(string(result), fmt.Sprintf(`> %d`, i))
			conn.Close()
		}
	})
	// gudp.Conn.SendRecvWithTimeout
	gtest.C(t, func(t *gtest.T) {
		for i := 0; i < 100; i++ {
			conn, err := gudp.NewConn(fmt.Sprintf("127.0.0.1:%d", p))
			t.AssertNil(err)
			result, err := conn.SendRecvWithTimeout([]byte(gconv.String(i)), -1, time.Second)
			t.AssertNil(err)
			t.Assert(string(result), fmt.Sprintf(`> %d`, i))
			conn.Close()
		}
	})
	// gudp.Send
	gtest.C(t, func(t *gtest.T) {
		for i := 0; i < 100; i++ {
			err := gudp.Send(fmt.Sprintf("127.0.0.1:%d", p), []byte(gconv.String(i)))
			t.AssertNil(err)
		}
	})
	// gudp.SendRecv
	gtest.C(t, func(t *gtest.T) {
		for i := 0; i < 100; i++ {
			result, err := gudp.SendRecv(fmt.Sprintf("127.0.0.1:%d", p), []byte(gconv.String(i)), -1)
			t.AssertNil(err)
			t.Assert(string(result), fmt.Sprintf(`> %d`, i))
		}
	})
}

// If the read buffer size is less than the sent package size,
// the rest data would be dropped.
func Test_Buffer(t *testing.T) {
	var (
		ctx = context.TODO()
	)
	p, _ := gudp.GetFreePort()
	s := gudp.NewServer(fmt.Sprintf("127.0.0.1:%d", p), func(conn *gudp.Conn) {
		defer conn.Close()
		for {
			data, err := conn.Recv(1)
			if len(data) > 0 {
				if err := conn.Send(data); err != nil {
					glog.Error(ctx, err)
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
	gtest.C(t, func(t *gtest.T) {
		result, err := gudp.SendRecv(fmt.Sprintf("127.0.0.1:%d", p), []byte("123"), -1)
		t.AssertNil(err)
		t.Assert(string(result), "1")
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := gudp.SendRecv(fmt.Sprintf("127.0.0.1:%d", p), []byte("456"), -1)
		t.AssertNil(err)
		t.Assert(string(result), "4")
	})
}

func Test_NewConn(t *testing.T) {
	var (
		port, _ = gudp.GetFreePort()
	)

	startUDPServer(fmt.Sprintf("127.0.0.1:%d", port))

	gtest.C(t, func(t *gtest.T) {
		conn, err := gudp.NewConn(fmt.Sprintf("127.0.0.1:%d", port), fmt.Sprintf("127.0.0.1:%d", port+1))
		t.AssertNil(err)
		conn.SetDeadline(time.Now().Add(time.Second))
		t.Assert(conn.Send(sendData), nil)
		conn.Close()
	})

	gtest.C(t, func(t *gtest.T) {
		conn, err := gudp.NewConn(fmt.Sprintf("127.0.0.1:%d", port), fmt.Sprintf("127.0.0.1:%d", 99999))
		t.AssertNil(conn)
		t.AssertNE(err, nil)
	})

	gtest.C(t, func(t *gtest.T) {
		conn, err := gudp.NewConn(fmt.Sprintf("127.0.0.1:%d", 99999))
		t.AssertNil(conn)
		t.AssertNE(err, nil)
	})
}

func Test_GetFreePorts(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ports, err := gudp.GetFreePorts(2)
		t.AssertNil(err)
		t.AssertEQ(len(ports), 2)
	})
}

func Test_Server(t *testing.T) {
	p, _ := gudp.GetFreePort()
	gudp.NewServer(fmt.Sprintf("127.0.0.1:%d", p), func(conn *gudp.Conn) {
		defer conn.Close()
		for {
			data, err := conn.Recv(1)
			if len(data) > 0 {
				conn.Send(data)
			}
			if err != nil {
				break
			}
		}
	}, "GoFrameUDPServer")

	gtest.C(t, func(t *gtest.T) {
		server := gudp.GetServer("GoFrameUDPServer")
		t.AssertNE(server, nil)
		server = gudp.GetServer("TestUDPServer")
		t.AssertNE(server, nil)
		server.SetAddress("127.0.0.1:8888")
		server.SetHandler(func(conn *gudp.Conn) {
			defer conn.Close()
			for {
				conn.Send([]byte("OtherHandle"))
			}
		})
	})
}
