// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gtcp_test

import (
	"fmt"
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
		err = conn.Send(data)
		t.AssertNil(err)
		conn.Close()
		conn.Conn.Close()
		time.Sleep(100 * time.Millisecond)
		err = conn.Send(data)
		t.AssertNil(err)
	})
}

func Test_Pool_Recv(t *testing.T) {
	p, _ := gtcp.GetFreePort()
	s := gtcp.NewServer(fmt.Sprintf(`:%d`, p), func(conn *gtcp.Conn) {
		conn.Send([]byte("Server Received"))
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
		_, err = conn.Recv(-1)
		t.AssertNil(err)
	})
}

func Test_Pool_RecvLine(t *testing.T) {
	p, _ := gtcp.GetFreePort()
	s := gtcp.NewServer(fmt.Sprintf(`:%d`, p), func(conn *gtcp.Conn) {
		conn.Send([]byte("Server Received\n"))
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
		_, err = conn.RecvLine()
		t.AssertNil(err)
		conn.Close()
		conn.Conn.Close()
		_, err = conn.RecvLine()
		t.AssertNE(err, nil)
	})
}

func Test_Pool_RecvTill(t *testing.T) {
	p, _ := gtcp.GetFreePort()
	s := gtcp.NewServer(fmt.Sprintf(`:%d`, p), func(conn *gtcp.Conn) {
		conn.Send([]byte("Server Received\n"))
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
		_, err = conn.RecvTill([]byte("\n"))
		t.AssertNil(err)
		conn.Close()
		conn.Conn.Close()
		_, err = conn.RecvTill([]byte("\n"))
		t.AssertNE(err, nil)
	})
}

func Test_Pool_RecvWithTimeout(t *testing.T) {
	p, _ := gtcp.GetFreePort()
	s := gtcp.NewServer(fmt.Sprintf(`:%d`, p), func(conn *gtcp.Conn) {
		conn.Send([]byte("Server Received\n"))
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
		_, err = conn.RecvWithTimeout(-1, time.Millisecond*500)
		t.AssertNil(err)
	})
}

func Test_Pool_SendWithTimeout(t *testing.T) {
	p, _ := gtcp.GetFreePort()
	s := gtcp.NewServer(fmt.Sprintf(`:%d`, p), func(conn *gtcp.Conn) {
		conn.Send([]byte("Server Received\n"))
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
	})
}

func Test_Pool_SendRecvWithTimeout(t *testing.T) {
	p, _ := gtcp.GetFreePort()
	s := gtcp.NewServer(fmt.Sprintf(`:%d`, p), func(conn *gtcp.Conn) {
		conn.Send([]byte("Server Received\n"))
	})
	go s.Run()
	defer s.Close()
	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		conn, err := gtcp.NewPoolConn(fmt.Sprintf("127.0.0.1:%d", p))
		t.AssertNil(err)
		defer conn.Close()
		data := []byte("9999")
		_, err = conn.SendRecvWithTimeout(data, -1, time.Millisecond*500)
		t.AssertNil(err)
	})
}
