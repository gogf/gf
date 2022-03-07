// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gtcp_test

import (
	"crypto/tls"
	"fmt"
	"github.com/gogf/gf/v2/util/gutil"
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

func ExampleNewConn() {
	var (
		addr = "127.0.0.1:80"
	)

	conn, _ := gtcp.NewConn(addr, time.Second)
	fmt.Println(conn)

	//
	s := gtcp.NewServer(addr, func(conn *gtcp.Conn) {
		defer conn.Close()
	})
	go s.Run()
	defer s.Close()

	conn, _ = gtcp.NewConn(addr, time.Second)
	fmt.Println(conn.RemoteAddr())

	// Output:
	// <nil>
	// 127.0.0.1:80
}

func ExampleNewConnTLS() {
	var (
		addr      = "127.0.0.1:80"
		tlsConfig = &tls.Config{}
	)

	conn, _ := gtcp.NewConnTLS(addr, tlsConfig)
	fmt.Println(conn)

	//
	s := gtcp.NewServerTLS(addr, tlsConfig, func(conn *gtcp.Conn) {
		defer conn.Close()
	})
	go s.Run()
	defer s.Close()

	gutil.TryCatch(func() {
		conn, _ = gtcp.NewConnTLS(addr, &tls.Config{})
	})

	// Output:
	// <nil>
}

func ExampleNewConnKeyCrt() {
	var (
		addr      = "127.0.0.1:80"
		tlsConfig = &tls.Config{}
		crtFile   = "crtFile"
		keyFile   = "keyFile"
	)

	conn, _ := gtcp.NewConnKeyCrt(addr, crtFile, keyFile)
	fmt.Println(conn)

	//
	s := gtcp.NewServerTLS(addr, tlsConfig, func(conn *gtcp.Conn) {
		defer conn.Close()
	})
	go s.Run()
	defer s.Close()

	gutil.TryCatch(func() {
		conn, _ = gtcp.NewConnKeyCrt(addr, crtFile, keyFile)
	})

	// Output:
	// <nil>
}

func ExampleConn_Send() {
	var (
		addr = "127.0.0.1:80"
	)

	s := gtcp.NewServer(addr, func(conn *gtcp.Conn) {
		time.Sleep(time.Second * 2)
		conn.Close()
	})
	defer s.Close()
	go s.Run()

	time.Sleep(time.Millisecond * 10)

	var (
		err error
	)
	conn, _ := gtcp.NewConn(addr)
	if conn != nil {
		for {
			err = conn.Send([]byte("hello"), gtcp.Retry{Count: 1})
			if err != nil {
				break
			}
			time.Sleep(time.Second)
		}
	}

	fmt.Println(err != nil)

	// Output:
	// true
}
