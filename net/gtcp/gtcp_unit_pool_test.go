// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gtcp_test

import (
	"fmt"
	"github.com/gogf/gf/net/gtcp"
	"github.com/gogf/gf/test/gtest"
	"testing"
	"time"
)

func Test_Pool_Basic1(t *testing.T) {
	p, _ := ports.PopRand()
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
		t.Assert(err, nil)
		defer conn.Close()
		data := []byte("9999")
		err = conn.SendPkg(data)
		t.Assert(err, nil)
		err = conn.SendPkgWithTimeout(data, time.Second)
		t.Assert(err, nil)
	})
}

func Test_Pool_Basic2(t *testing.T) {
	p, _ := ports.PopRand()
	s := gtcp.NewServer(fmt.Sprintf(`:%d`, p), func(conn *gtcp.Conn) {
		conn.Close()
	})
	go s.Run()
	defer s.Close()
	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		conn, err := gtcp.NewPoolConn(fmt.Sprintf("127.0.0.1:%d", p))
		t.Assert(err, nil)
		defer conn.Close()
		data := []byte("9999")
		err = conn.SendPkg(data)
		t.Assert(err, nil)
		//err = conn.SendPkgWithTimeout(data, time.Second)
		//t.Assert(err, nil)

		_, err = conn.SendRecv(data, -1)
		t.AssertNE(err, nil)
	})
}
