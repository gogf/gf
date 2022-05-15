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
	"github.com/gogf/gf/v2/os/gfile"
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
		})
		go s.Run()
	}

	time.Sleep(time.Millisecond * 500)

	err1 := gtcp.Send(fmt.Sprintf("%s:%d", host, ports[0]), []byte("hello"), gtcp.Retry{Count: 1})
	err2 := gtcp.Send(fmt.Sprintf("%s:%d", host, ports[1]), []byte("hello"), gtcp.Retry{Count: 1})

	fmt.Println(err1 == nil && err2 == nil)

	// Output:
	// true
}

func ExampleNewConn() {
	addr := "127.0.0.1:%d"
	freePort, _ := gtcp.GetFreePort()
	addr = fmt.Sprintf(addr, freePort)

	conn, _ := gtcp.NewConn(addr, time.Second)
	fmt.Println(conn)

	//
	s := gtcp.NewServer(addr, func(conn *gtcp.Conn) {
		conn.Send([]byte("Server Received"))
	})
	go s.Run()
	defer s.Close()

	gtcp.NewConn(addr, time.Second)

	// Output:
	// <nil>
}

func ExampleNewConnTLS() {
	var (
		tlsConfig = &tls.Config{}
	)
	addr := "127.0.0.1:%d"
	freePort, _ := gtcp.GetFreePort()
	addr = fmt.Sprintf(addr, freePort)

	conn, _ := gtcp.NewConnTLS(addr, tlsConfig)
	fmt.Println(conn)

	//
	s := gtcp.NewServerTLS(addr, tlsConfig, func(conn *gtcp.Conn) {
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
		tlsConfig = &tls.Config{}
		crtFile   = gfile.Dir(gdebug.CallerFilePath()) + gfile.Separator + "testdata/crtFile"
		keyFile   = gfile.Dir(gdebug.CallerFilePath()) + gfile.Separator + "testdata/keyFile"
	)

	addr := "127.0.0.1:%d"
	freePort, _ := gtcp.GetFreePort()
	addr = fmt.Sprintf(addr, freePort)

	conn, _ := gtcp.NewConnKeyCrt(addr, crtFile, keyFile)
	fmt.Println(conn)

	//
	s := gtcp.NewServerTLS(addr, tlsConfig, func(conn *gtcp.Conn) {
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
	addr := "127.0.0.1:%d"
	freePort, _ := gtcp.GetFreePort()
	addr = fmt.Sprintf(addr, freePort)

	s := gtcp.NewServer(addr, func(conn *gtcp.Conn) {
		time.Sleep(time.Second * 2)
		conn.Close()
	})
	defer s.Close()
	go s.Run()

	time.Sleep(time.Millisecond * 500)

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
	addr := "127.0.0.1:%d"
	freePort, _ := gtcp.GetFreePort()
	addr = fmt.Sprintf(addr, freePort)

	s := gtcp.NewServer(addr, func(conn *gtcp.Conn) {
		time.Sleep(time.Second * 2)
		conn.Close()
	})
	defer s.Close()
	go s.Run()

	time.Sleep(time.Millisecond * 500)

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
	addr := "127.0.0.1:%d"
	freePort, _ := gtcp.GetFreePort()
	addr = fmt.Sprintf(addr, freePort)

	s := gtcp.NewServer(addr, func(conn *gtcp.Conn) {
		for {
			if _, err := conn.Recv(-1); err == nil {
				conn.Send([]byte("Server Received"))
			}
		}
	})
	defer s.Close()
	go s.Run()

	time.Sleep(time.Millisecond * 500)

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
	addr := "127.0.0.1:%d"
	freePort, _ := gtcp.GetFreePort()
	addr = fmt.Sprintf(addr, freePort)

	s := gtcp.NewServer(addr, func(conn *gtcp.Conn) {
		for {
			if _, err := conn.Recv(-1); err == nil {
				conn.Send([]byte("Server Received"))
			}
		}
	})
	defer s.Close()
	go s.Run()

	time.Sleep(time.Millisecond * 500)

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
		sendContent = make([]byte, 512)
	)

	addr := "127.0.0.1:%d"
	freePort, _ := gtcp.GetFreePort()
	addr = fmt.Sprintf(addr, freePort)

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

	time.Sleep(time.Millisecond * 500)

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

func ExampleConn_Recv_Once() {
	addr := "127.0.0.1:%d"
	freePort, _ := gtcp.GetFreePort()
	addr = fmt.Sprintf(addr, freePort)

	s := gtcp.NewServer(addr, func(conn *gtcp.Conn) {
		if _, err := conn.Recv(0); err == nil {
			conn.Close()
		} else {
			fmt.Println(err)
		}
	})
	defer s.Close()
	go s.Run()

	time.Sleep(time.Millisecond * 500)

	var (
		err error
	)
	conn, _ := gtcp.NewConn(addr)
	if conn != nil {
		for {
			_, err = conn.SendRecv([]byte("hello"), -1)
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

func ExampleConn_RecvWithTimeout() {
	var (
		sendContent = make([]byte, 512)
	)

	addr := "127.0.0.1:%d"
	freePort, _ := gtcp.GetFreePort()
	addr = fmt.Sprintf(addr, freePort)

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

	time.Sleep(time.Millisecond * 500)

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
	addr := "127.0.0.1:%d"
	freePort, _ := gtcp.GetFreePort()
	addr = fmt.Sprintf(addr, freePort)

	s := gtcp.NewServer(addr, func(conn *gtcp.Conn) {
		for {
			if _, err := conn.RecvLine(gtcp.Retry{Count: 1}); err != nil {
			}
		}
	})
	defer s.Close()
	go s.Run()

	time.Sleep(time.Millisecond * 500)

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
	addr := "127.0.0.1:%d"
	freePort, _ := gtcp.GetFreePort()
	addr = fmt.Sprintf(addr, freePort)

	s := gtcp.NewServer(addr, func(conn *gtcp.Conn) {
		for {
			if _, err := conn.RecvTill([]byte("finish"), gtcp.Retry{Count: 1}); err != nil {

			}
		}
	})
	defer s.Close()
	go s.Run()

	time.Sleep(time.Millisecond * 500)

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
	addr := "127.0.0.1:%d"
	freePort, _ := gtcp.GetFreePort()
	addr = fmt.Sprintf(addr, freePort)

	s := gtcp.NewServer(addr, func(conn *gtcp.Conn) {
		time.Sleep(time.Second * 2)
		conn.Close()
	})
	defer s.Close()
	go s.Run()

	time.Sleep(time.Millisecond * 500)

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
	addr := "127.0.0.1:%d"
	freePort, _ := gtcp.GetFreePort()
	addr = fmt.Sprintf(addr, freePort)

	s := gtcp.NewServer(addr, func(conn *gtcp.Conn) {
		time.Sleep(time.Second * 2)
		conn.Close()
	})
	defer s.Close()
	go s.Run()

	time.Sleep(time.Millisecond * 500)

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

func ExampleNewNetConnKeyCrt() {
	var (
		crtFile = gfile.Dir(gdebug.CallerFilePath()) + gfile.Separator + "testdata/server.crt"
		keyFile = gfile.Dir(gdebug.CallerFilePath()) + gfile.Separator + "testdata/server.key"
	)

	addr := "127.0.0.1:%d"
	freePort, _ := gtcp.GetFreePort()
	addr = fmt.Sprintf(addr, freePort)

	s := gtcp.NewServer(addr, func(conn *gtcp.Conn) {
	})
	defer s.Close()
	go s.Run()

	time.Sleep(time.Millisecond * 500)

	conn, _ := gtcp.NewNetConnKeyCrt(addr, crtFile, keyFile, time.Second)
	fmt.Println(conn != nil)

	// Output:
	// false
}

func ExampleSend() {
	addr := "127.0.0.1:%d"
	freePort, _ := gtcp.GetFreePort()
	addr = fmt.Sprintf(addr, freePort)

	s := gtcp.NewServer(addr, func(conn *gtcp.Conn) {
	})
	defer s.Close()
	go s.Run()

	time.Sleep(time.Millisecond * 500)

	err := gtcp.Send("127.0.0.1:80", []byte("hello"))

	fmt.Println(err != nil)

	err = gtcp.Send(addr, []byte("hello"), gtcp.Retry{Count: 1})
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(err == nil)

	// Output:
	// true
	// true
}

func ExampleSendRecv() {
	addr := "127.0.0.1:%d"
	freePort, _ := gtcp.GetFreePort()
	addr = fmt.Sprintf(addr, freePort)

	s := gtcp.NewServer(addr, func(conn *gtcp.Conn) {
		conn.Send([]byte("Server Received"))
	})
	defer s.Close()
	go s.Run()

	time.Sleep(time.Millisecond * 500)

	_, err := gtcp.SendRecv("127.0.0.1:80", []byte("hello"), -1)
	fmt.Println(err != nil)

	_, err = gtcp.SendRecv(addr, []byte("hello"), -1, gtcp.Retry{Count: 1})
	fmt.Println(err == nil)

	// Output:
	// true
	// true
}

func ExampleSendWithTimeout() {
	addr := "127.0.0.1:%d"
	freePort, _ := gtcp.GetFreePort()
	addr = fmt.Sprintf(addr, freePort)

	s := gtcp.NewServer(addr, func(conn *gtcp.Conn) {
		conn.Send([]byte("Server Received"))
	})
	defer s.Close()
	go s.Run()

	time.Sleep(time.Millisecond * 500)

	err := gtcp.SendWithTimeout("127.0.0.1:80", []byte("hello"), time.Millisecond*500)
	fmt.Println(err != nil)

	err = gtcp.SendWithTimeout(addr, []byte("hello"), time.Millisecond*500, gtcp.Retry{Count: 1})
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(err == nil)

	// Output:
	// true
	// true
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
		crtFile = gfile.Dir(gdebug.CallerFilePath()) + gfile.Separator + "testdata/crtFile"
		keyFile = gfile.Dir(gdebug.CallerFilePath()) + gfile.Separator + "testdata/keyFile"
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
	s.Run()

	fmt.Println(s != nil)

	// Output:
	// true
}
