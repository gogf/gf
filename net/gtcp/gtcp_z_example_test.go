// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gtcp_test

import (
	"fmt"
	"github.com/gogf/gf/v2/debug/gdebug"
	"github.com/gogf/gf/v2/os/gfile"
	"time"

	"github.com/gogf/gf/v2/net/gtcp"
)

func ExampleGetFreePort() {
	fmt.Println(gtcp.GetFreePort())

	// May Output:
	// 57429 <nil>
}

func ExampleGetFreePorts() {
	fmt.Println(gtcp.GetFreePorts(2))

	// May Output:
	// [57743 57744] <nil>
}

func ExampleSendRecvWithTimeout() {
	addr := "127.0.0.1:%d"
	freePort, _ := gtcp.GetFreePort()
	addr = fmt.Sprintf(addr, freePort)

	s := gtcp.NewServer(addr, func(conn *gtcp.Conn) {
		conn.Send([]byte("Server Received"))
	})
	defer s.Close()
	go s.Run()

	time.Sleep(time.Millisecond * 500)

	_, err := gtcp.SendRecvWithTimeout("127.0.0.1:80", []byte("hello"), -1, time.Millisecond*500)
	fmt.Println(err != nil)

	_, err = gtcp.SendRecvWithTimeout(addr, []byte("hello"), -1, time.Millisecond*500, gtcp.Retry{Count: 1})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(err == nil)

	// Output:
	// true
	// true
}

func ExampleMustGetFreePort() {
	var (
		host = "127.0.0.1"
		port = gtcp.MustGetFreePort()
	)

	addr := fmt.Sprintf("%s:%d", host, port)

	s := gtcp.NewServer(addr, func(conn *gtcp.Conn) {
	})
	defer s.Close()
	go s.Run()

	time.Sleep(time.Millisecond * 500)

	err := gtcp.Send(addr, []byte("hello"), gtcp.Retry{Count: 1})
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(err == nil)

	// Output:
	// true
}

func ExampleSendPkg() {
	addr := "127.0.0.1:%d"
	freePort, _ := gtcp.GetFreePort()
	addr = fmt.Sprintf(addr, freePort)

	s := gtcp.NewServer(addr, func(conn *gtcp.Conn) {
	})
	defer s.Close()
	go s.Run()

	time.Sleep(time.Millisecond * 500)

	err := gtcp.SendPkg("127.0.0.1:80", []byte("hello"))
	fmt.Println(err != nil)

	err = gtcp.SendPkg(addr, []byte("hello"))
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(err == nil)

	// Output:
	// true
	// true
}

func ExampleSendRecvPkg() {
	addr := "127.0.0.1:%d"
	freePort, _ := gtcp.GetFreePort()
	addr = fmt.Sprintf(addr, freePort)

	s := gtcp.NewServer(addr, func(conn *gtcp.Conn) {
		conn.SendPkg([]byte("Server Received"), gtcp.PkgOption{Retry: gtcp.Retry{Count: 1}})
	})
	defer s.Close()
	go s.Run()

	time.Sleep(time.Millisecond * 500)

	_, err := gtcp.SendRecvPkg("127.0.0.1:80", []byte("hello"))
	fmt.Println(err != nil)

	_, err = gtcp.SendRecvPkg(addr, []byte("hello"))
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(err == nil)

	// Output:
	// true
	// true
}

func ExampleSendPkgWithTimeout() {
	addr := "127.0.0.1:%d"
	freePort, _ := gtcp.GetFreePort()
	addr = fmt.Sprintf(addr, freePort)

	s := gtcp.NewServer(addr, func(conn *gtcp.Conn) {
		conn.SendPkg([]byte("Server Received"))
	})
	defer s.Close()
	go s.Run()

	time.Sleep(time.Millisecond * 500)

	err := gtcp.SendPkgWithTimeout("127.0.0.1:80", []byte("hello"), time.Second)
	fmt.Println(err != nil)

	err = gtcp.SendPkgWithTimeout(addr, []byte("hello"), time.Second)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(err == nil)

	// Output:
	// true
	// true
}

func ExampleSendRecvPkgWithTimeout() {
	addr := "127.0.0.1:%d"
	freePort, _ := gtcp.GetFreePort()
	addr = fmt.Sprintf(addr, freePort)

	s := gtcp.NewServer(addr, func(conn *gtcp.Conn) {
		conn.SendPkg([]byte("Server Received"))
	})
	defer s.Close()
	go s.Run()

	time.Sleep(time.Millisecond * 500)

	_, err := gtcp.SendRecvPkgWithTimeout("127.0.0.1:80", []byte("hello"), time.Second)
	fmt.Println(err != nil)

	_, err = gtcp.SendRecvPkgWithTimeout(addr, []byte("hello"), time.Second)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(err == nil)

	// Output:
	// true
	// true
}

func ExampleNewServer() {
	addr := "127.0.0.1:%d"
	freePort, _ := gtcp.GetFreePort()
	addr = fmt.Sprintf(addr, freePort)

	s := gtcp.NewServer(addr, func(conn *gtcp.Conn) {
		conn.SendPkg([]byte("Server Received"))
	}, "NewServer")
	defer s.Close()
	go s.Run()

	time.Sleep(time.Millisecond * 500)

	err := gtcp.Send(addr, []byte("hello"), gtcp.Retry{Count: 1})
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(err == nil)

	// Output:
	// true
}

func ExampleGetServer() {
	addr := "127.0.0.1:%d"
	freePort, _ := gtcp.GetFreePort()
	addr = fmt.Sprintf(addr, freePort)

	s := gtcp.GetServer("GetServer")
	defer s.Close()
	go s.Run()

	fmt.Println(s != nil)

	// Output:
	// true
}

func ExampleSetAddress() {
	addr := "127.0.0.1:%d"
	freePort, _ := gtcp.GetFreePort()
	addr = fmt.Sprintf(addr, freePort)

	s := gtcp.NewServer("", func(conn *gtcp.Conn) {
	})
	s.SetAddress(addr)
	defer s.Close()
	go s.Run()

	time.Sleep(time.Millisecond * 500)

	err := gtcp.Send(addr, []byte("hello"), gtcp.Retry{Count: 1})
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(err == nil)

	// Output:
	// true
}

func ExampleSetHandler() {
	addr := "127.0.0.1:%d"
	freePort, _ := gtcp.GetFreePort()
	addr = fmt.Sprintf(addr, freePort)

	s := gtcp.NewServer(addr, nil)
	s.SetHandler(func(conn *gtcp.Conn) {
	})
	defer s.Close()
	go s.Run()

	time.Sleep(time.Millisecond * 500)

	err := gtcp.Send(addr, []byte("hello"), gtcp.Retry{Count: 1})
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(err == nil)

	// Output:
	// true
}

func ExampleRun_NilHandle() {
	addr := "127.0.0.1:%d"
	freePort, _ := gtcp.GetFreePort()
	addr = fmt.Sprintf(addr, freePort)

	s := gtcp.NewServer(addr, nil)
	defer s.Close()
	go func() {
		err := s.Run()
		fmt.Println(err != nil)
	}()

	time.Sleep(time.Millisecond * 100)

	// Output:
	// true
}

func ExampleNewServerKeyCrt() {
	var (
		crtFile = gfile.Dir(gdebug.CallerFilePath()) + gfile.Separator + "testdata/server.crt"
		keyFile = gfile.Dir(gdebug.CallerFilePath()) + gfile.Separator + "testdata/server.key"
	)
	addr := "127.0.0.1:%d"
	freePort, _ := gtcp.GetFreePort()
	addr = fmt.Sprintf(addr, freePort)

	s, err := gtcp.NewServerKeyCrt(addr, crtFile, keyFile, func(conn *gtcp.Conn) {
	})
	if err != nil {
		fmt.Println(s == nil)
		return
	}
	defer s.Close()
	go s.Run()

	fmt.Println(s != nil)

	// Output:
	// true
}
