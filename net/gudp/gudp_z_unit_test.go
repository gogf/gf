// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gudp_test

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/net/gudp"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gconv"
)

var (
	ports = garray.NewIntArray(true)
)

func init() {
	for i := 9000; i <= 10000; i++ {
		ports.Append(i)
	}
}

func Test_Basic(t *testing.T) {
	var (
		ctx = context.TODO()
	)
	p, _ := ports.PopRand()
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
			t.Assert(err, nil)
			t.Assert(conn.Send([]byte(gconv.String(i))), nil)
			conn.Close()
		}
	})
	// gudp.Conn.SendRecv
	gtest.C(t, func(t *gtest.T) {
		for i := 0; i < 100; i++ {
			conn, err := gudp.NewConn(fmt.Sprintf("127.0.0.1:%d", p))
			t.Assert(err, nil)
			_, err = conn.SendRecv([]byte(gconv.String(i)), -1)
			t.Assert(err, nil)
			//t.Assert(string(result), fmt.Sprintf(`> %d`, i))
			conn.Close()
		}
	})
	// gudp.Send
	gtest.C(t, func(t *gtest.T) {
		for i := 0; i < 100; i++ {
			err := gudp.Send(fmt.Sprintf("127.0.0.1:%d", p), []byte(gconv.String(i)))
			t.Assert(err, nil)
		}
	})
	// gudp.SendRecv
	gtest.C(t, func(t *gtest.T) {
		for i := 0; i < 100; i++ {
			result, err := gudp.SendRecv(fmt.Sprintf("127.0.0.1:%d", p), []byte(gconv.String(i)), -1)
			t.Assert(err, nil)
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
	p, _ := ports.PopRand()
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
		t.Assert(err, nil)
		t.Assert(string(result), "1")
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := gudp.SendRecv(fmt.Sprintf("127.0.0.1:%d", p), []byte("456"), -1)
		t.Assert(err, nil)
		t.Assert(string(result), "4")
	})
}

func Test_Conn(t *testing.T) {
	var (
		ctx = context.TODO()
	)
	p, _ := ports.PopRand()
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

	gtest.C(t, func(t *gtest.T) {
		for i := 0; i < 100; i++ {
			targetAddr := net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: p}
			uc, err := net.DialUDP("udp", nil, &targetAddr)
			t.Assert(err, nil)
			conn := gudp.NewConnByNetConn(uc)

			err = conn.SendWithTimeout([]byte(gconv.String(i)), time.Minute)
			t.Assert(err, nil)

			res1, err := conn.RecvWithTimeout(-1, time.Minute)
			t.Assert(err, nil)
			t.Assert(string(res1), fmt.Sprintf(`> %d`, i))

			r := conn.RemoteAddr()
			t.Assert(r.String(), fmt.Sprintf("127.0.0.1:%d", p))

			conn.Close()
		}
	})

	gtest.C(t, func(t *gtest.T) {
		for i := 0; i < 100; i++ {
			targetAddr := net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: p}
			uc, err := net.DialUDP("udp", nil, &targetAddr)
			t.Assert(err, nil)
			conn := gudp.NewConnByNetConn(uc)

			err = conn.SetDeadline(time.Now().Add(time.Minute))
			t.Assert(err, nil)
			err = conn.SetRecvDeadline(time.Now().Add(time.Minute))
			t.Assert(err, nil)
			err = conn.SetSendDeadline(time.Now().Add(time.Minute))
			t.Assert(err, nil)
			conn.SetRecvBufferWait(2 * time.Millisecond)

			err = conn.SendWithTimeout([]byte(gconv.String(i)), time.Minute)
			t.Assert(err, nil)

			res1, err := conn.RecvWithTimeout(-1, time.Minute)
			t.Assert(err, nil)
			t.Assert(string(res1), fmt.Sprintf(`> %d`, i))

			conn.Close()
		}
	})

	gtest.C(t, func(t *gtest.T) {
		for i := 0; i < 100; i++ {
			targetAddr := net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: p}
			uc, err := net.DialUDP("udp", nil, &targetAddr)
			t.Assert(err, nil)
			conn := gudp.NewConnByNetConn(uc)

			res1, err := conn.SendRecvWithTimeout([]byte(gconv.String(i)), -1, time.Minute)
			t.Assert(err, nil)
			t.Assert(string(res1), fmt.Sprintf(`> %d`, i))

			conn.Close()
		}
	})

}

func Test_Server(t *testing.T) {
	var (
		ctx = context.TODO()
	)
	p, _ := ports.PopRand()
	ss := gudp.NewServer("", nil, "jeff")
	ss.SetAddress(fmt.Sprintf("127.0.0.1:%d", p))
	ss.SetHandler(func(conn *gudp.Conn) {
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
	s := gudp.GetServer("jeff")
	gtest.C(t, func(t *gtest.T) {
		t.Assert(ss, s)
	})

	go s.Run()
	defer s.Close()
	time.Sleep(100 * time.Millisecond)

	gtest.C(t, func(t *gtest.T) {
		for i := 0; i < 100; i++ {
			targetAddr := net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: p}
			uc, err := net.DialUDP("udp", nil, &targetAddr)
			t.Assert(err, nil)
			conn := gudp.NewConnByNetConn(uc)

			res1, err := conn.SendRecv([]byte(gconv.String(i)), -1)
			t.Assert(err, nil)
			t.Assert(string(res1), fmt.Sprintf(`> %d`, i))

			conn.Close()
		}
	})

}

func Test_GetServer(t *testing.T) {
	var (
		ctx = context.TODO()
	)
	p, _ := ports.PopRand()
	ss := gudp.GetServer("")
	ss.SetAddress(fmt.Sprintf("127.0.0.1:%d", p))
	ss.SetHandler(func(conn *gudp.Conn) {
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

	s := gudp.GetServer("jeff")
	gtest.C(t, func(t *gtest.T) {
		t.Assert(ss, s)
	})
}

func Test_Server_Addr(t *testing.T) {
	p, _ := ports.PopRand()
	ss := gudp.GetServer("")

	ss.SetAddress(fmt.Sprintf("0.0.0.1:%d", p))
	e := ss.Run()
	time.Sleep(100 * time.Millisecond)

	gtest.C(t, func(t *gtest.T) {
		t.AssertNE(e, nil)
	})
}

func Test_Server_Handler(t *testing.T) {
	s := gudp.GetServer("")
	s.SetHandler(nil)
	e := s.Run()
	time.Sleep(100 * time.Millisecond)

	gtest.C(t, func(t *gtest.T) {
		t.AssertNE(e, nil)
	})
}

func Test_Conn_Addr(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		conn, err := gudp.NewConn("")
		t.AssertNE(err, nil)
		_ = conn
	})

	gtest.C(t, func(t *gtest.T) {
		conn, err := gudp.NewConn("nil")
		t.AssertNE(err, nil)
		_ = conn
	})

	gtest.C(t, func(t *gtest.T) {
		conn, err := gudp.NewConn("127.0.0.1:8000", "nil")
		t.AssertNE(err, nil)
		_ = conn
	})

	gtest.C(t, func(t *gtest.T) {
		conn, err := gudp.NewConn("127.0.0.1:8000", "127.0.0.0:8000")
		t.AssertNE(err, nil)
		_ = conn
	})
}

func Test_Send_Addr(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		err := gudp.Send("", []byte("hello"))
		t.AssertNE(err, nil)
	})

	gtest.C(t, func(t *gtest.T) {
		_, err := gudp.SendRecv("", []byte("hello"), -1)
		t.AssertNE(err, nil)
	})
}

func Test_Conn_Para(t *testing.T) {
	p, _ := ports.PopRand()

	gtest.C(t, func(t *gtest.T) {
		conn, _ := gudp.NewConn(fmt.Sprintf("127.0.0.1:%d", p))
		conn.Close()

		t.AssertNE(conn.SetSendDeadline(time.Now()), nil)
		t.AssertNE(conn.SetDeadline(time.Now().Add(time.Minute)), nil)
		t.AssertNE(conn.SetRecvDeadline(time.Now().Add(time.Minute)), nil)
		t.AssertNE(conn.SetSendDeadline(time.Now().Add(time.Minute)), nil)
		t.AssertNE(conn.Close(), nil)
	})
}

func Test_Conn_Send(t *testing.T) {
	p, _ := ports.PopRand()

	gtest.C(t, func(t *gtest.T) {
		conn, err := gudp.NewConn(fmt.Sprintf("127.0.0.1:%d", p))
		conn.Close()

		t.AssertNE(conn.Send([]byte("hello")), nil)
		t.AssertNE(conn.SendWithTimeout([]byte("hello"), time.Second), nil)

		_, err = conn.Recv(-1)
		t.AssertNE(err, nil)
		_, err = conn.RecvWithTimeout(-1, time.Second)
		t.AssertNE(err, nil)

		_, err = conn.SendRecv([]byte("hello"), -1)
		t.AssertNE(err, nil)

		t.AssertNE(conn.SendWithTimeout([]byte("hello"), time.Minute), nil)
		_, err = conn.RecvWithTimeout(-1, time.Minute)
		t.AssertNE(err, nil)
	})
}

func Test_Retry(t *testing.T) {
	p, _ := ports.PopRand()

	gtest.C(t, func(t *gtest.T) {
		conn, err := gudp.NewConn(fmt.Sprintf("127.0.0.1:%d", p))

		_, err = conn.RecvWithTimeout(10, time.Microsecond, gudp.Retry{1, time.Millisecond})
		t.AssertNE(err, nil)

		buf := make([]byte, 1024*1024*8, 1024*1024*8)
		err = conn.SendWithTimeout(buf, time.Microsecond, gudp.Retry{1, time.Millisecond})
		t.AssertNE(err, nil)

		conn.Close()
	})
}

func Test_Retry_2(t *testing.T) {
	var (
		ctx = context.TODO()
	)
	p, _ := ports.PopRand()

	s := gudp.NewServer(fmt.Sprintf("127.0.0.1:%d", p), func(conn *gudp.Conn) {
		defer conn.Close()
		for {
			data, err := conn.Recv(-1)
			if len(data) > 0 {
				time.Sleep(100 * time.Millisecond)
				if err := conn.Send(append([]byte(""), data...)); err != nil {
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
		conn, _ := gudp.NewConn(fmt.Sprintf("127.0.0.1:%d", p))
		res1, _ := conn.SendRecvWithTimeout([]byte(gconv.String(1)), 1, time.Second,
			gudp.Retry{3, 50 * time.Millisecond})
		t.Assert(string(res1), fmt.Sprintf(`%d`, 1))
		conn.Close()
	})
}
