// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gtcp_test

import (
	"fmt"
	"github.com/gogf/gf/v2/text/gstr"
	"testing"
	"time"

	"github.com/gogf/gf/v2/net/gtcp"
	"github.com/gogf/gf/v2/test/gtest"
)

func Test_Pool_Basic1(t *testing.T) {
	p, _ := gtcp.GetFreePort()
	s := gtcp.NewServer(fmt.Sprintf(`:%d`, p), func(conn *gtcp.Conn) {
		defer conn.Close()
		for {
			data, err := conn.RecvPkg()
			if err != nil {
				break
			}
			conn.SendPkg(data)
		}
	})
	go s.Run()
	defer s.Close()
	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		conn, err := gtcp.NewPoolConn(fmt.Sprintf("127.0.0.1:%d", p))
		t.AssertNil(err)
		defer conn.Close()
		data := []byte("9999")
		err = conn.SendPkg(data)
		t.AssertNil(err)
		err = conn.SendPkgWithTimeout(data, time.Second)
		t.AssertNil(err)
	})

	gtest.C(t, func(t *gtest.T) {
		_, err := gtcp.NewPoolConn(fmt.Sprintf("127.0.0.1:80"))
		t.AssertNE(err, nil)
	})
}

func Test_Pool_Basic2(t *testing.T) {
	p, _ := gtcp.GetFreePort()
	s := gtcp.NewServer(fmt.Sprintf(`:%d`, p), func(conn *gtcp.Conn) {
		conn.Close()
	})
	go s.Run()
	defer s.Close()
	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		conn, err := gtcp.NewPoolConn(fmt.Sprintf("127.0.0.1:%d", p))
		t.AssertNil(err)
		defer conn.Close()
		data := []byte("9999")
		err = conn.SendPkg(data)
		t.AssertNil(err)
		//err = conn.SendPkgWithTimeout(data, time.Second)
		//t.AssertNil(err)

		_, err = conn.SendRecv(data, -1)
		t.AssertNE(err, nil)
	})
}

func Test_Pool_Send(t *testing.T) {
	p, _ := gtcp.GetFreePort()
	s := gtcp.NewServer(fmt.Sprintf(`:%d`, p), func(conn *gtcp.Conn) {
		for {
			data, err := conn.Recv(-1)
			if err != nil {
				break
			}
			conn.Send(data)
		}
	})
	go s.Run()
	defer s.Close()
	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		conn, err := gtcp.NewPoolConn(fmt.Sprintf("127.0.0.1:%d", p))
		t.AssertNil(err)
		defer conn.Close()
		data := []byte("9999")
		err = conn.Send(data)
		t.AssertNil(err)
		recv, err := conn.Recv(-1)
		t.AssertNil(err)
		t.Assert(recv, data)
	})
}

func Test_Pool_Recv(t *testing.T) {
	p, _ := gtcp.GetFreePort()
	s := gtcp.NewServer(fmt.Sprintf(`:%d`, p), func(conn *gtcp.Conn) {
		for {
			data, err := conn.Recv(-1)
			if err != nil {
				break
			}
			conn.Send(data)
		}
	})
	go s.Run()
	defer s.Close()
	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		conn, err := gtcp.NewPoolConn(fmt.Sprintf("127.0.0.1:%d", p))
		t.AssertNil(err)
		defer conn.Close()
		data := []byte("9999")
		err = conn.Send(data)
		t.AssertNil(err)
		time.Sleep(100 * time.Millisecond)
		recv, err := conn.Recv(-1)
		t.AssertNil(err)
		t.Assert(recv, data)
	})
}

func Test_Pool_RecvLine(t *testing.T) {
	p, _ := gtcp.GetFreePort()
	s := gtcp.NewServer(fmt.Sprintf(`:%d`, p), func(conn *gtcp.Conn) {
		for {
			data, err := conn.Recv(-1)
			if err != nil {
				break
			}
			conn.Send(data)
		}
	})
	go s.Run()
	defer s.Close()
	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		conn, err := gtcp.NewPoolConn(fmt.Sprintf("127.0.0.1:%d", p))
		t.AssertNil(err)
		defer conn.Close()
		data := []byte("9999\n")
		err = conn.Send(data)
		t.AssertNil(err)
		time.Sleep(100 * time.Millisecond)
		recv, err := conn.RecvLine()
		t.AssertNil(err)
		splitData := gstr.Split(string(data), "\n")
		t.Assert(recv, splitData[0])
	})
}

func Test_Pool_RecvTill(t *testing.T) {
	p, _ := gtcp.GetFreePort()
	s := gtcp.NewServer(fmt.Sprintf(`:%d`, p), func(conn *gtcp.Conn) {
		for {
			data, err := conn.Recv(-1)
			if err != nil {
				break
			}
			conn.Send(data)
		}
	})
	go s.Run()
	defer s.Close()
	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		conn, err := gtcp.NewPoolConn(fmt.Sprintf("127.0.0.1:%d", p))
		t.AssertNil(err)
		defer conn.Close()
		data := []byte("9999\n")
		err = conn.Send(data)
		t.AssertNil(err)
		time.Sleep(100 * time.Millisecond)
		recv, err := conn.RecvTill([]byte("\n"))
		t.AssertNil(err)
		t.Assert(recv, data)
	})
}

func Test_Pool_RecvWithTimeout(t *testing.T) {
	p, _ := gtcp.GetFreePort()
	s := gtcp.NewServer(fmt.Sprintf(`:%d`, p), func(conn *gtcp.Conn) {
		for {
			data, err := conn.Recv(-1)
			if err != nil {
				break
			}
			conn.Send(data)
		}
	})
	go s.Run()
	defer s.Close()
	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		conn, err := gtcp.NewPoolConn(fmt.Sprintf("127.0.0.1:%d", p))
		t.AssertNil(err)
		defer conn.Close()
		data := []byte("9999")
		err = conn.Send(data)
		t.AssertNil(err)
		time.Sleep(100 * time.Millisecond)
		recv, err := conn.RecvWithTimeout(-1, time.Millisecond*500)
		t.AssertNil(err)
		t.Assert(data, recv)
	})
}

func Test_Pool_SendWithTimeout(t *testing.T) {
	p, _ := gtcp.GetFreePort()
	s := gtcp.NewServer(fmt.Sprintf(`:%d`, p), func(conn *gtcp.Conn) {
		for {
			data, err := conn.Recv(-1)
			if err != nil {
				break
			}
			conn.Send(data)
		}
	})
	go s.Run()
	defer s.Close()
	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		conn, err := gtcp.NewPoolConn(fmt.Sprintf("127.0.0.1:%d", p))
		t.AssertNil(err)
		defer conn.Close()
		data := []byte("9999")
		err = conn.SendWithTimeout(data, time.Millisecond*500)
		t.AssertNil(err)
		recv, err := conn.Recv(-1)
		t.AssertNil(err)
		t.Assert(data, recv)
	})
}

func Test_Pool_SendRecvWithTimeout(t *testing.T) {
	p, _ := gtcp.GetFreePort()
	s := gtcp.NewServer(fmt.Sprintf(`:%d`, p), func(conn *gtcp.Conn) {
		for {
			data, err := conn.Recv(-1)
			if err != nil {
				break
			}
			conn.Send(data)
		}
	})
	go s.Run()
	defer s.Close()
	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		conn, err := gtcp.NewPoolConn(fmt.Sprintf("127.0.0.1:%d", p))
		t.AssertNil(err)
		defer conn.Close()
		data := []byte("9999")
		recv, err := conn.SendRecvWithTimeout(data, -1, time.Millisecond*500)
		t.AssertNil(err)
		t.Assert(data, recv)
	})
}
