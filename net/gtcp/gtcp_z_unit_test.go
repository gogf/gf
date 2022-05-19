// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gtcp_test

import (
	"crypto/tls"
	"fmt"
	"github.com/gogf/gf/v2/debug/gdebug"
	"github.com/gogf/gf/v2/net/gtcp"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/test/gtest"
	"testing"
	"time"
)

var (
	simpleTimeout = time.Millisecond * 100
	crtFile       = gfile.Dir(gdebug.CallerFilePath()) + gfile.Separator + "testdata/server.crt"
	keyFile       = gfile.Dir(gdebug.CallerFilePath()) + gfile.Separator + "testdata/server.key"
)

func getFreePortAddr() string {
	addr := "127.0.0.1:%d"
	freePort, _ := gtcp.GetFreePort()
	return fmt.Sprintf(addr, freePort)
}

func startTCPServer(addr string) {
	s := gtcp.NewServer(addr, func(conn *gtcp.Conn) {
		data := []byte("gtcp Server received")
		conn.Send(data)
		time.Sleep(simpleTimeout)
		conn.Close()
	})
	go s.Run()
}

func startTCPPkgServer(addr string) {
	s := gtcp.NewServer(addr, func(conn *gtcp.Conn) {
		data := []byte("gtcp Server received")
		conn.SendPkg(data)
		time.Sleep(simpleTimeout)
		conn.Close()
	})
	go s.Run()
}

func startTCPTLSServer(addr string) {
	tlsConfig := &tls.Config{}
	s := gtcp.NewServerTLS(addr, tlsConfig, func(conn *gtcp.Conn) {
		data := []byte("gtcp tls Server received")
		conn.Send(data)
		time.Sleep(simpleTimeout)
		conn.Close()
	})
	go s.Run()
}

func startTCPKeyCrtServer(addr string) {
	s, _ := gtcp.NewServerKeyCrt(addr, crtFile, keyFile, func(conn *gtcp.Conn) {
		data := []byte("gtcp tls Server received")
		conn.Send(data)
		time.Sleep(simpleTimeout)
		conn.Close()
	})
	go s.Run()
}

func TestGetFreePorts(t *testing.T) {
	ports, _ := gtcp.GetFreePorts(2)
	gtest.C(t, func(t *gtest.T) {
		t.AssertGT(ports[0], 0)
		t.AssertGT(ports[1], 0)
	})

	addr := fmt.Sprintf("%s:%d", "127.0.0.1", ports[0])
	startTCPServer(addr)
	time.Sleep(simpleTimeout)

	gtest.C(t, func(t *gtest.T) {
		conn, err := gtcp.NewPoolConn(fmt.Sprintf("127.0.0.1:%d", ports[0]))
		t.AssertNil(err)
		defer conn.Close()
		data := []byte("9999")
		recv, err := conn.SendRecv(data, -1)
		t.AssertNil(err)
		t.AssertNE(recv, nil)
	})
	gtest.C(t, func(t *gtest.T) {
		conn, err := gtcp.NewPoolConn(fmt.Sprintf("127.0.0.1:%d", 80))
		t.AssertNE(err, nil)
		t.AssertNil(conn)
	})
}

func TestMustGetFreePort(t *testing.T) {
	port := gtcp.MustGetFreePort()
	addr := fmt.Sprintf("%s:%d", "127.0.0.1", port)
	startTCPServer(addr)
	time.Sleep(simpleTimeout)

	gtest.C(t, func(t *gtest.T) {
		recv, err := gtcp.SendRecv(addr, []byte("hello"), -1)
		t.AssertNil(err)
		t.AssertGT(len(recv), 0)
	})
}

func TestNewConn(t *testing.T) {
	addr := getFreePortAddr()

	gtest.C(t, func(t *gtest.T) {
		conn, err := gtcp.NewConn(addr, simpleTimeout)
		t.AssertNil(conn)
		t.AssertNE(err, nil)
	})

	gtest.C(t, func(t *gtest.T) {
		startTCPServer(addr)

		time.Sleep(simpleTimeout)

		conn, err := gtcp.NewConn(addr, simpleTimeout)
		t.AssertNil(err)
		t.AssertNE(conn, nil)
		defer conn.Close()
		data := []byte("9999")
		recv, err := conn.SendRecv(data, -1)
		t.AssertNil(err)
		t.AssertNE(recv, nil)
	})
}

//TODO
func TestNewConnTLS(t *testing.T) {
	addr := getFreePortAddr()

	gtest.C(t, func(t *gtest.T) {
		conn, err := gtcp.NewConnTLS(addr, &tls.Config{})
		t.AssertNil(conn)
		t.AssertNE(err, nil)
	})

	gtest.C(t, func(t *gtest.T) {
		startTCPTLSServer(addr)

		time.Sleep(simpleTimeout)

		conn, err := gtcp.NewConnTLS(addr, &tls.Config{})
		t.AssertNil(conn)
		t.AssertNE(err, nil)
	})
}

func TestNewConnKeyCrt(t *testing.T) {
	addr := getFreePortAddr()

	gtest.C(t, func(t *gtest.T) {
		conn, err := gtcp.NewConnKeyCrt(addr, crtFile, keyFile)
		t.AssertNil(conn)
		t.AssertNE(err, nil)
	})

	gtest.C(t, func(t *gtest.T) {
		startTCPKeyCrtServer(addr)

		time.Sleep(simpleTimeout)

		conn, err := gtcp.NewConnKeyCrt(addr, crtFile, keyFile)
		t.AssertNil(conn)
		t.AssertNE(err, nil)
	})
}

func TestConn_Send(t *testing.T) {
	addr := getFreePortAddr()

	startTCPServer(addr)

	time.Sleep(simpleTimeout)

	gtest.C(t, func(t *gtest.T) {
		conn, err := gtcp.NewConn(addr)
		t.AssertNil(err)
		t.AssertNE(conn, nil)
		err = conn.Send([]byte("hello"), gtcp.Retry{Count: 1})
		t.AssertNil(err)
		recv, err := conn.Recv(-1)
		t.AssertNil(err)
		t.AssertNE(recv, nil)
	})
}

func TestConn_SendWithTimeout(t *testing.T) {
	addr := getFreePortAddr()

	startTCPServer(addr)

	time.Sleep(simpleTimeout)

	gtest.C(t, func(t *gtest.T) {
		conn, err := gtcp.NewConn(addr)
		t.AssertNil(err)
		t.AssertNE(conn, nil)
		err = conn.SendWithTimeout([]byte("hello"), time.Second, gtcp.Retry{Count: 1})
		t.AssertNil(err)
		recv, err := conn.Recv(-1)
		t.AssertNil(err)
		t.AssertNE(recv, nil)
	})
}

func TestConn_SendRecv(t *testing.T) {
	addr := getFreePortAddr()

	startTCPServer(addr)

	time.Sleep(simpleTimeout)

	gtest.C(t, func(t *gtest.T) {
		conn, err := gtcp.NewConn(addr)
		t.AssertNil(err)
		t.AssertNE(conn, nil)
		recv, err := conn.SendRecv([]byte("hello"), -1, gtcp.Retry{Count: 1})
		t.AssertNil(err)
		t.AssertNE(recv, nil)
	})
}

func TestConn_SendRecvWithTimeout(t *testing.T) {
	addr := getFreePortAddr()

	startTCPServer(addr)

	time.Sleep(simpleTimeout)

	gtest.C(t, func(t *gtest.T) {
		conn, err := gtcp.NewConn(addr)
		t.AssertNil(err)
		t.AssertNE(conn, nil)
		recv, err := conn.SendRecvWithTimeout([]byte("hello"), -1, time.Second, gtcp.Retry{Count: 1})
		t.AssertNil(err)
		t.AssertNE(recv, nil)
	})
}

func TestConn_RecvWithTimeout(t *testing.T) {
	addr := getFreePortAddr()

	startTCPServer(addr)

	time.Sleep(simpleTimeout)

	gtest.C(t, func(t *gtest.T) {
		conn, err := gtcp.NewConn(addr)
		t.AssertNil(err)
		t.AssertNE(conn, nil)
		conn.Send([]byte("hello"))
		recv, err := conn.RecvWithTimeout(-1, time.Second, gtcp.Retry{Count: 1})
		t.AssertNil(err)
		t.AssertNE(recv, nil)
	})
}

func TestConn_RecvLine(t *testing.T) {
	addr := getFreePortAddr()

	s := gtcp.NewServer(addr, func(conn *gtcp.Conn) {
		data := []byte("gtcp Server received\n")
		conn.Send(data)
		time.Sleep(simpleTimeout)
		conn.Close()
	})
	go s.Run()

	time.Sleep(simpleTimeout)

	gtest.C(t, func(t *gtest.T) {
		conn, err := gtcp.NewConn(addr)
		t.AssertNil(err)
		t.AssertNE(conn, nil)
		conn.Send([]byte("hello"))
		recv, err := conn.RecvLine(gtcp.Retry{Count: 1})
		t.AssertNil(err)
		t.AssertNE(recv, nil)
	})
}

func TestConn_RecvTill(t *testing.T) {
	addr := getFreePortAddr()

	startTCPServer(addr)

	time.Sleep(simpleTimeout)

	gtest.C(t, func(t *gtest.T) {
		conn, err := gtcp.NewConn(addr)
		t.AssertNil(err)
		t.AssertNE(conn, nil)
		conn.Send([]byte("hello"))
		recv, err := conn.RecvTill([]byte("received"), gtcp.Retry{Count: 1})
		t.AssertNil(err)
		t.AssertNE(recv, nil)
	})
}

func TestConn_SetDeadline(t *testing.T) {
	addr := getFreePortAddr()

	startTCPServer(addr)

	time.Sleep(simpleTimeout)

	gtest.C(t, func(t *gtest.T) {
		conn, err := gtcp.NewConn(addr)
		t.AssertNil(err)
		t.AssertNE(conn, nil)
		conn.SetDeadline(time.Time{})
		err = conn.Send([]byte("hello"), gtcp.Retry{Count: 1})
		t.AssertNil(err)
		recv, err := conn.Recv(-1)
		t.AssertNil(err)
		t.AssertNE(recv, nil)
	})
}

func TestConn_SetReceiveBufferWait(t *testing.T) {
	addr := getFreePortAddr()

	startTCPServer(addr)

	time.Sleep(simpleTimeout)

	gtest.C(t, func(t *gtest.T) {
		conn, err := gtcp.NewConn(addr)
		t.AssertNil(err)
		t.AssertNE(conn, nil)
		conn.SetReceiveBufferWait(time.Millisecond * 100)
		err = conn.Send([]byte("hello"), gtcp.Retry{Count: 1})
		t.AssertNil(err)
		recv, err := conn.Recv(-1)
		t.AssertNil(err)
		t.AssertNE(recv, nil)
	})
}

func TestNewNetConnKeyCrt(t *testing.T) {
	addr := getFreePortAddr()

	startTCPKeyCrtServer(addr)

	time.Sleep(simpleTimeout)

	gtest.C(t, func(t *gtest.T) {
		conn, err := gtcp.NewNetConnKeyCrt(addr, crtFile, keyFile, time.Second)
		t.AssertNil(conn)
		t.AssertNE(err, nil)
	})
}

func TestSend(t *testing.T) {
	addr := getFreePortAddr()

	startTCPServer(addr)

	time.Sleep(simpleTimeout)

	gtest.C(t, func(t *gtest.T) {
		err := gtcp.Send(addr, []byte("hello"), gtcp.Retry{Count: 1})
		t.AssertNil(err)
	})
}

func TestSendRecv(t *testing.T) {
	addr := getFreePortAddr()

	startTCPServer(addr)

	time.Sleep(simpleTimeout)

	gtest.C(t, func(t *gtest.T) {
		recv, err := gtcp.SendRecv(addr, []byte("hello"), -1)
		t.AssertNil(err)
		t.AssertGT(len(recv), 0)
	})
}

func TestSendWithTimeout(t *testing.T) {
	addr := getFreePortAddr()

	startTCPServer(addr)

	time.Sleep(simpleTimeout)

	gtest.C(t, func(t *gtest.T) {
		err := gtcp.SendWithTimeout("127.0.0.1:80", []byte("hello"), time.Millisecond*500)
		t.AssertNE(err, nil)
		err = gtcp.SendWithTimeout(addr, []byte("hello"), time.Millisecond*500)
		t.AssertNil(err)
	})
}

func TestSendRecvWithTimeout(t *testing.T) {
	addr := getFreePortAddr()

	startTCPServer(addr)

	time.Sleep(simpleTimeout)

	gtest.C(t, func(t *gtest.T) {
		recv, err := gtcp.SendRecvWithTimeout("127.0.0.1:80", []byte("hello"), -1, time.Millisecond*500)
		t.AssertNil(recv)
		t.AssertNE(err, nil)
		recv, err = gtcp.SendRecvWithTimeout(addr, []byte("hello"), -1, time.Millisecond*500)
		t.AssertNil(err)
		t.AssertNE(recv, nil)
	})
}

func TestSendPkg(t *testing.T) {
	addr := getFreePortAddr()

	startTCPPkgServer(addr)

	time.Sleep(simpleTimeout)

	gtest.C(t, func(t *gtest.T) {
		err := gtcp.SendPkg(addr, []byte("hello"))
		t.AssertNil(err)
		err = gtcp.SendPkg("127.0.0.1:80", []byte("hello"))
		t.AssertNE(err, nil)
	})
}
