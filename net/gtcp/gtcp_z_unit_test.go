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
	for _, port := range ports {
		addr := fmt.Sprintf("%s:%d", "127.0.0.1:%d", port)

		startTCPServer(addr)
	}
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

func TestNewConn(t *testing.T) {
	addr := getFreePortAddr()

	gtest.C(t, func(t *gtest.T) {
		conn, err := gtcp.NewConn(addr, simpleTimeout)
		t.AssertNil(conn)
		t.AssertNE(err, nil)
	})

	gtest.C(t, func(t *gtest.T) {
		startTCPServer(addr)

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
		conn, err := gtcp.NewNetConnTLS(addr, &tls.Config{}, simpleTimeout)
		t.AssertNil(conn)
		t.AssertNE(err, nil)
	})

	gtest.C(t, func(t *gtest.T) {
		startTCPTLSServer(addr)

		conn, err := gtcp.NewNetConnTLS(addr, &tls.Config{}, simpleTimeout)
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

		conn, err := gtcp.NewConnKeyCrt(addr, crtFile, keyFile)
		t.AssertNil(conn)
		t.AssertNE(err, nil)
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
