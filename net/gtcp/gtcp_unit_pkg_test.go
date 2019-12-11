// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gtcp_test

import (
	"fmt"
	"github.com/gogf/gf/net/gtcp"
	"github.com/gogf/gf/test/gtest"
	"github.com/gogf/gf/util/gconv"
	"testing"
	"time"
)

func Test_Package_Basic(t *testing.T) {
	p := ports.PopRand()
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
	// SendPkg
	gtest.Case(t, func() {
		conn, err := gtcp.NewConn(fmt.Sprintf("127.0.0.1:%d", p))
		gtest.Assert(err, nil)
		defer conn.Close()
		for i := 0; i < 100; i++ {
			err := conn.SendPkg([]byte(gconv.String(i)))
			gtest.Assert(err, nil)
		}
		for i := 0; i < 100; i++ {
			err := conn.SendPkgWithTimeout([]byte(gconv.String(i)), time.Second)
			gtest.Assert(err, nil)
		}
	})
	// SendPkg with big data - failure.
	gtest.Case(t, func() {
		conn, err := gtcp.NewConn(fmt.Sprintf("127.0.0.1:%d", p))
		gtest.Assert(err, nil)
		defer conn.Close()
		data := make([]byte, 65536)
		err = conn.SendPkg(data)
		gtest.AssertNE(err, nil)
	})
	// SendRecvPkg
	gtest.Case(t, func() {
		conn, err := gtcp.NewConn(fmt.Sprintf("127.0.0.1:%d", p))
		gtest.Assert(err, nil)
		defer conn.Close()
		for i := 100; i < 200; i++ {
			data := []byte(gconv.String(i))
			result, err := conn.SendRecvPkg(data)
			gtest.Assert(err, nil)
			gtest.Assert(result, data)
		}
		for i := 100; i < 200; i++ {
			data := []byte(gconv.String(i))
			result, err := conn.SendRecvPkgWithTimeout(data, time.Second)
			gtest.Assert(err, nil)
			gtest.Assert(result, data)
		}
	})
	// SendRecvPkg with big data - failure.
	gtest.Case(t, func() {
		conn, err := gtcp.NewConn(fmt.Sprintf("127.0.0.1:%d", p))
		gtest.Assert(err, nil)
		defer conn.Close()
		data := make([]byte, 65536)
		result, err := conn.SendRecvPkg(data)
		gtest.AssertNE(err, nil)
		gtest.Assert(result, nil)
	})
	// SendRecvPkg with big data - success.
	gtest.Case(t, func() {
		conn, err := gtcp.NewConn(fmt.Sprintf("127.0.0.1:%d", p))
		gtest.Assert(err, nil)
		defer conn.Close()
		data := make([]byte, 65500)
		data[100] = byte(65)
		data[65400] = byte(85)
		result, err := conn.SendRecvPkg(data)
		gtest.Assert(err, nil)
		gtest.Assert(result, data)
	})
}

func Test_Package_Timeout(t *testing.T) {
	p := ports.PopRand()
	s := gtcp.NewServer(fmt.Sprintf(`:%d`, p), func(conn *gtcp.Conn) {
		defer conn.Close()
		for {
			data, err := conn.RecvPkg()
			if err != nil {
				break
			}
			time.Sleep(time.Second)
			gtest.Assert(conn.SendPkg(data), nil)
		}
	})
	go s.Run()
	defer s.Close()
	time.Sleep(100 * time.Millisecond)
	gtest.Case(t, func() {
		conn, err := gtcp.NewConn(fmt.Sprintf("127.0.0.1:%d", p))
		gtest.Assert(err, nil)
		defer conn.Close()
		data := []byte("10000")
		result, err := conn.SendRecvPkgWithTimeout(data, time.Millisecond*500)
		gtest.AssertNE(err, nil)
		gtest.Assert(result, nil)
	})
	gtest.Case(t, func() {
		conn, err := gtcp.NewConn(fmt.Sprintf("127.0.0.1:%d", p))
		gtest.Assert(err, nil)
		defer conn.Close()
		data := []byte("10000")
		result, err := conn.SendRecvPkgWithTimeout(data, time.Second*2)
		gtest.Assert(err, nil)
		gtest.Assert(result, data)
	})
}

func Test_Package_Option(t *testing.T) {
	p := ports.PopRand()
	s := gtcp.NewServer(fmt.Sprintf(`:%d`, p), func(conn *gtcp.Conn) {
		defer conn.Close()
		option := gtcp.PkgOption{HeaderSize: 1}
		for {
			data, err := conn.RecvPkg(option)
			if err != nil {
				break
			}
			gtest.Assert(conn.SendPkg(data, option), nil)
		}
	})
	go s.Run()
	defer s.Close()
	time.Sleep(100 * time.Millisecond)
	// SendRecvPkg with big data - failure.
	gtest.Case(t, func() {
		conn, err := gtcp.NewConn(fmt.Sprintf("127.0.0.1:%d", p))
		gtest.Assert(err, nil)
		defer conn.Close()
		data := make([]byte, 0xFF+1)
		result, err := conn.SendRecvPkg(data, gtcp.PkgOption{HeaderSize: 1})
		gtest.AssertNE(err, nil)
		gtest.Assert(result, nil)
	})
	// SendRecvPkg with big data - success.
	gtest.Case(t, func() {
		conn, err := gtcp.NewConn(fmt.Sprintf("127.0.0.1:%d", p))
		gtest.Assert(err, nil)
		defer conn.Close()
		data := make([]byte, 0xFF)
		data[100] = byte(65)
		data[200] = byte(85)
		result, err := conn.SendRecvPkg(data, gtcp.PkgOption{HeaderSize: 1})
		gtest.Assert(err, nil)
		gtest.Assert(result, data)
	})
}
