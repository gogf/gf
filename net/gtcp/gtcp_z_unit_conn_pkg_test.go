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
	"github.com/gogf/gf/v2/util/gconv"
)

func Test_Package_Basic(t *testing.T) {
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
	// SendPkg
	gtest.C(t, func(t *gtest.T) {
		conn, err := gtcp.NewConn(fmt.Sprintf("127.0.0.1:%d", p))
		t.AssertNil(err)
		defer conn.Close()
		for i := 0; i < 100; i++ {
			err := conn.SendPkg([]byte(gconv.String(i)))
			t.AssertNil(err)
		}
		for i := 0; i < 100; i++ {
			err := conn.SendPkgWithTimeout([]byte(gconv.String(i)), time.Second)
			t.AssertNil(err)
		}
	})
	// SendPkg with big data - failure.
	gtest.C(t, func(t *gtest.T) {
		conn, err := gtcp.NewConn(fmt.Sprintf("127.0.0.1:%d", p))
		t.AssertNil(err)
		defer conn.Close()
		data := make([]byte, 65536)
		err = conn.SendPkg(data)
		t.AssertNE(err, nil)
	})
	// SendRecvPkg
	gtest.C(t, func(t *gtest.T) {
		conn, err := gtcp.NewConn(fmt.Sprintf("127.0.0.1:%d", p))
		t.AssertNil(err)
		defer conn.Close()
		for i := 100; i < 200; i++ {
			data := []byte(gconv.String(i))
			result, err := conn.SendRecvPkg(data)
			t.AssertNil(err)
			t.Assert(result, data)
		}
		for i := 100; i < 200; i++ {
			data := []byte(gconv.String(i))
			result, err := conn.SendRecvPkgWithTimeout(data, time.Second)
			t.AssertNil(err)
			t.Assert(result, data)
		}
	})
	// SendRecvPkg with big data - failure.
	gtest.C(t, func(t *gtest.T) {
		conn, err := gtcp.NewConn(fmt.Sprintf("127.0.0.1:%d", p))
		t.AssertNil(err)
		defer conn.Close()
		data := make([]byte, 65536)
		result, err := conn.SendRecvPkg(data)
		t.AssertNE(err, nil)
		t.Assert(result, nil)
	})
	// SendRecvPkg with big data - success.
	gtest.C(t, func(t *gtest.T) {
		conn, err := gtcp.NewConn(fmt.Sprintf("127.0.0.1:%d", p))
		t.AssertNil(err)
		defer conn.Close()
		data := make([]byte, 65500)
		data[100] = byte(65)
		data[65400] = byte(85)
		result, err := conn.SendRecvPkg(data)
		t.AssertNil(err)
		t.Assert(result, data)
	})
}

func Test_Package_Timeout(t *testing.T) {
	p, _ := gtcp.GetFreePort()
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
	gtest.C(t, func(t *gtest.T) {
		conn, err := gtcp.NewConn(fmt.Sprintf("127.0.0.1:%d", p))
		t.AssertNil(err)
		defer conn.Close()
		data := []byte("10000")
		result, err := conn.SendRecvPkgWithTimeout(data, time.Millisecond*500)
		t.AssertNE(err, nil)
		t.Assert(result, nil)
	})
	gtest.C(t, func(t *gtest.T) {
		conn, err := gtcp.NewConn(fmt.Sprintf("127.0.0.1:%d", p))
		t.AssertNil(err)
		defer conn.Close()
		data := []byte("10000")
		result, err := conn.SendRecvPkgWithTimeout(data, time.Second*2)
		t.AssertNil(err)
		t.Assert(result, data)
	})
}

func Test_Package_Option(t *testing.T) {
	p, _ := gtcp.GetFreePort()
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
	gtest.C(t, func(t *gtest.T) {
		conn, err := gtcp.NewConn(fmt.Sprintf("127.0.0.1:%d", p))
		t.AssertNil(err)
		defer conn.Close()
		data := make([]byte, 0xFF+1)
		result, err := conn.SendRecvPkg(data, gtcp.PkgOption{HeaderSize: 1})
		t.AssertNE(err, nil)
		t.Assert(result, nil)
	})
	// SendRecvPkg with big data - success.
	gtest.C(t, func(t *gtest.T) {
		conn, err := gtcp.NewConn(fmt.Sprintf("127.0.0.1:%d", p))
		t.AssertNil(err)
		defer conn.Close()
		data := make([]byte, 0xFF)
		data[100] = byte(65)
		data[200] = byte(85)
		result, err := conn.SendRecvPkg(data, gtcp.PkgOption{HeaderSize: 1})
		t.AssertNil(err)
		t.Assert(result, data)
	})
}

func Test_Package_Option_HeadSize4(t *testing.T) {
	p, _ := gtcp.GetFreePort()
	s := gtcp.NewServer(fmt.Sprintf(`:%d`, p), func(conn *gtcp.Conn) {
		defer conn.Close()
		option := gtcp.PkgOption{HeaderSize: 4}
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
	gtest.C(t, func(t *gtest.T) {
		conn, err := gtcp.NewConn(fmt.Sprintf("127.0.0.1:%d", p))
		t.AssertNil(err)
		defer conn.Close()
		data := make([]byte, 0xFFFF+1)
		_, err = conn.SendRecvPkg(data, gtcp.PkgOption{HeaderSize: 4})
		t.Assert(err, nil)
	})
	// SendRecvPkg with big data - success.
	gtest.C(t, func(t *gtest.T) {
		conn, err := gtcp.NewConn(fmt.Sprintf("127.0.0.1:%d", p))
		t.AssertNil(err)
		defer conn.Close()
		data := make([]byte, 0xFF)
		data[100] = byte(65)
		data[200] = byte(85)
		result, err := conn.SendRecvPkg(data, gtcp.PkgOption{HeaderSize: 4})
		t.AssertNil(err)
		t.Assert(result, data)
	})
	// pkgOption.HeaderSize oversize
	gtest.C(t, func(t *gtest.T) {
		conn, err := gtcp.NewConn(fmt.Sprintf("127.0.0.1:%d", p))
		t.AssertNil(err)
		defer conn.Close()
		data := make([]byte, 0xFF)
		data[100] = byte(65)
		data[200] = byte(85)
		_, err = conn.SendRecvPkg(data, gtcp.PkgOption{HeaderSize: 5})
		t.AssertNE(err, nil)
	})
}
