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
	var (
		host = "127.0.0.1"
	)

	ports, _ := gtcp.GetFreePorts(2)

	for _, port := range ports {
		addr := fmt.Sprintf("%s:%d", host, port)

		s := gtcp.NewServer(addr, func(conn *gtcp.Conn) {
			time.Sleep(time.Second * 2)
			conn.Close()
		})
		defer s.Close()
		go s.Run()
	}

	time.Sleep(time.Millisecond * 10)

	err1 := gtcp.Send(fmt.Sprintf("%s:%d", host, ports[0]), []byte("hello"), gtcp.Retry{Count: 1})
	err2 := gtcp.Send(fmt.Sprintf("%s:%d", host, ports[1]), []byte("hello"), gtcp.Retry{Count: 1})

	fmt.Println(err1 == nil && err2 == nil)

	// Output:
	// true
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
		gtcp.NewNetConnTLS(addr, &tls.Config{}, time.Second)
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

func ExampleConn_SendWithTimeout() {
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
			err = conn.SendWithTimeout([]byte("hello"), time.Second, gtcp.Retry{Count: 1})
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

func ExampleConn_SendRecv() {
	var (
		addr = "127.0.0.1:80"
	)

	s := gtcp.NewServer(addr, func(conn *gtcp.Conn) {
		for {
			if _, err := conn.Recv(-1); err == nil {
				conn.Send([]byte("Server Received"))
			}
		}
	})
	defer s.Close()
	go s.Run()

	time.Sleep(time.Millisecond * 10)

	var (
		err error
	)
	conn, _ := gtcp.NewConn(addr)
	if conn != nil {
		sendCount := 0
		for {
			if sendCount > 3 {
				break
			}

			if _, err = conn.SendRecv([]byte("hello server"), -1); err != nil {
				break
			}

			sendCount++
			time.Sleep(time.Millisecond * 200)
		}
	}

	fmt.Println(err != nil)

	// Output:
	// false
}

func ExampleConn_SendRecvWithTimeout() {
	var (
		addr = "127.0.0.1:80"
	)

	s := gtcp.NewServer(addr, func(conn *gtcp.Conn) {
		for {
			if _, err := conn.Recv(-1); err == nil {
				conn.Send([]byte("Server Received"))
			}
		}
	})
	defer s.Close()
	go s.Run()

	time.Sleep(time.Millisecond * 10)

	var (
		err error
	)
	conn, _ := gtcp.NewConn(addr)
	if conn != nil {
		sendCount := 0
		for {
			if sendCount > 3 {
				break
			}

			if _, err = conn.SendRecvWithTimeout([]byte("hello server"), -1, time.Millisecond*500); err != nil {
				break
			}

			sendCount++
			time.Sleep(time.Millisecond * 200)
		}
	}

	fmt.Println(err != nil)

	// Output:
	// false
}

func ExampleConn_Recv() {
	var (
		addr        = "127.0.0.1:80"
		sendContent = make([]byte, 512)
	)
	for i := 0; i < 512; i++ {
		sendContent[i] = 'a'
	}

	s := gtcp.NewServer(addr, func(conn *gtcp.Conn) {
		for {
			if _, err := conn.Recv(-1); err != nil {
				conn.Close()
				return
			}
		}
	})
	defer s.Close()
	go s.Run()

	time.Sleep(time.Millisecond * 10)

	var (
		err error
	)
	conn, _ := gtcp.NewConn(addr)
	if conn != nil {
		sendCount := 0
		for {
			if sendCount > 5 {
				break
			}
			err = conn.Send(sendContent[0:64*sendCount], gtcp.Retry{Count: 1})
			if err != nil {
				break
			}
			time.Sleep(time.Millisecond * 200)
			sendCount++
		}
	}

	fmt.Println(err != nil)

	// Output:
	// false
}

func ExampleConn_RecvWithTimeout() {
	var (
		addr        = "127.0.0.1:80"
		sendContent = make([]byte, 512)
	)
	for i := 0; i < 512; i++ {
		sendContent[i] = 'a'
	}

	s := gtcp.NewServer(addr, func(conn *gtcp.Conn) {
		for {
			if _, err := conn.RecvWithTimeout(-1, time.Millisecond*500, gtcp.Retry{Count: 1}); err != nil {
			}
		}
	})
	defer s.Close()
	go s.Run()

	time.Sleep(time.Millisecond * 10)

	var (
		err error
	)
	conn, _ := gtcp.NewConn(addr)
	if conn != nil {
		sendCount := 0
		for {
			sendCount++
			if sendCount > 5 {
				break
			}
			err = conn.Send(sendContent[0:64*sendCount], gtcp.Retry{Count: 1})
			if err != nil {
				break
			}
			time.Sleep(time.Millisecond * 200 * time.Duration(sendCount))
		}
	}

	fmt.Println(err != nil)

	// Output:
	// false
}

func ExampleConn_RecvLine() {
	var (
		addr = "127.0.0.1:80"
	)

	s := gtcp.NewServer(addr, func(conn *gtcp.Conn) {
		for {
			if c, err := conn.RecvLine(gtcp.Retry{Count: 1}); err != nil {
				fmt.Println(c)
			}
		}
	})
	defer s.Close()
	go s.Run()

	time.Sleep(time.Millisecond * 10)

	var (
		err         error
		sendContent = []string{"hello", "\n", "world"}
	)
	conn, _ := gtcp.NewConn(addr)
	if conn != nil {
		sendCount := 0
		for {
			sendCount++
			if sendCount >= len(sendContent) {
				break
			}
			err = conn.Send([]byte(sendContent[sendCount]))
			if err != nil {
				break
			}
			time.Sleep(time.Millisecond * 200)
		}
	}

	fmt.Println(err != nil)

	// Output:
	// false
}

func ExampleConn_RecvTill() {
	var (
		addr = "127.0.0.1:80"
	)

	s := gtcp.NewServer(addr, func(conn *gtcp.Conn) {
		for {
			if c, err := conn.RecvTill([]byte("finish"), gtcp.Retry{Count: 1}); err != nil {
				fmt.Println(c)
			}
		}
	})
	defer s.Close()
	go s.Run()

	time.Sleep(time.Millisecond * 10)

	var (
		err         error
		sendContent = []string{"hello", "world", "finish"}
	)
	conn, _ := gtcp.NewConn(addr)
	if conn != nil {
		sendCount := 0
		for {
			sendCount++
			if sendCount >= len(sendContent) {
				break
			}
			err = conn.Send([]byte(sendContent[sendCount]))
			if err != nil {
				break
			}
			time.Sleep(time.Millisecond * 200)
		}
	}

	fmt.Println(err != nil)

	// Output:
	// false
}

func ExampleConn_SetDeadline() {
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
		conn.SetDeadline(time.Time{})
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

func ExampleConn_SetReceiveBufferWait() {
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
		conn.SetReceiveBufferWait(time.Millisecond * 100)
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

func ExampleSend() {
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

	err := gtcp.Send(addr, []byte("hello"), gtcp.Retry{Count: 1})

	fmt.Println(err == nil)

	// Output:
	// true
}

func ExampleSendRecv() {
	var (
		addr = "127.0.0.1:80"
	)

	s := gtcp.NewServer(addr, func(conn *gtcp.Conn) {
		conn.Send([]byte("Server Received"))
		time.Sleep(time.Second * 2)
		conn.Close()
	})
	defer s.Close()
	go s.Run()

	time.Sleep(time.Millisecond * 10)

	_, err := gtcp.SendRecv(addr, []byte("hello"), -1, gtcp.Retry{Count: 1})

	fmt.Println(err == nil)

	// Output:
	// true
}

func ExampleSendWithTimeout() {
	var (
		addr = "127.0.0.1:80"
	)

	s := gtcp.NewServer(addr, func(conn *gtcp.Conn) {
		conn.Send([]byte("Server Received"))
		time.Sleep(time.Second * 2)
		conn.Close()
	})
	defer s.Close()
	go s.Run()

	time.Sleep(time.Millisecond * 10)

	err := gtcp.SendWithTimeout(addr, []byte("hello"), time.Millisecond*500, gtcp.Retry{Count: 1})

	fmt.Println(err == nil)

	// Output:
	// true
}

func ExampleSendRecvWithTimeout() {
	var (
		addr = "127.0.0.1:80"
	)

	s := gtcp.NewServer(addr, func(conn *gtcp.Conn) {
		conn.Send([]byte("Server Received"))
		time.Sleep(time.Second * 2)
		conn.Close()
	})
	defer s.Close()
	go s.Run()

	time.Sleep(time.Millisecond * 10)

	_, err := gtcp.SendRecvWithTimeout(addr, []byte("hello"), -1, time.Millisecond*500, gtcp.Retry{Count: 1})

	fmt.Println(err == nil)

	// Output:
	// true
}

func ExampleMustGetFreePort() {
	var (
		host = "127.0.0.1"
		port = gtcp.MustGetFreePort()
	)

	addr := fmt.Sprintf("%s:%d", host, port)

	s := gtcp.NewServer(addr, func(conn *gtcp.Conn) {
		time.Sleep(time.Second * 2)
		conn.Close()
	})
	defer s.Close()
	go s.Run()

	time.Sleep(time.Millisecond * 10)

	err := gtcp.Send(addr, []byte("hello"), gtcp.Retry{Count: 1})

	fmt.Println(err == nil)

	// Output:
	// true
}

func ExampleSendPkg() {
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

	err := gtcp.SendPkg(addr, []byte("hello"))

	fmt.Println(err == nil)

	// Output:
	// true
}

func ExampleSendRecvPkg() {
	var (
		addr = "127.0.0.1:80"
	)

	s := gtcp.NewServer(addr, func(conn *gtcp.Conn) {
		conn.SendPkg([]byte("Server Received"))
		time.Sleep(time.Second * 2)
		conn.Close()
	})
	defer s.Close()
	go s.Run()

	time.Sleep(time.Millisecond * 10)

	_, err := gtcp.SendRecvPkg(addr, []byte("hello"))

	fmt.Println(err == nil)

	// Output:
	// true
}

func ExampleSendPkgWithTimeout() {
	var (
		addr = "127.0.0.1:80"
	)

	s := gtcp.NewServer(addr, func(conn *gtcp.Conn) {
		conn.SendPkg([]byte("Server Received"))
		time.Sleep(time.Second * 2)
		conn.Close()
	})
	defer s.Close()
	go s.Run()

	time.Sleep(time.Millisecond * 10)

	err := gtcp.SendPkgWithTimeout(addr, []byte("hello"), time.Second)

	fmt.Println(err == nil)

	// Output:
	// true
}

func ExampleSendRecvPkgWithTimeout() {
	var (
		addr = "127.0.0.1:80"
	)

	s := gtcp.NewServer(addr, func(conn *gtcp.Conn) {
		conn.SendPkg([]byte("Server Received"))
		time.Sleep(time.Second * 2)
		conn.Close()
	})
	defer s.Close()
	go s.Run()

	time.Sleep(time.Millisecond * 10)

	_, err := gtcp.SendRecvPkgWithTimeout(addr, []byte("hello"), time.Second)

	fmt.Println(err == nil)

	// Output:
	// true
}
