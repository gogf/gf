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
	"github.com/gogf/gf/v2/text/gstr"
	"testing"
	"time"
)

var (
	simpleTimeout = time.Millisecond * 100
	sendData      = []byte("hello")
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
		defer conn.Close()
		for {
			data, err := conn.Recv(-1)
			if err != nil {
				break
			}
			conn.Send(data)
		}
	})
	go s.Run()
	time.Sleep(simpleTimeout)
}

func startTCPPkgServer(addr string) {
	s := gtcp.NewServer(addr, func(conn *gtcp.Conn) {
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
	time.Sleep(simpleTimeout)
}

func startTCPTLSServer(addr string) {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		Certificates: []tls.Certificate{
			tls.Certificate{},
		},
	}
	s := gtcp.NewServerTLS(addr, tlsConfig, func(conn *gtcp.Conn) {
		defer conn.Close()
		for {
			data, err := conn.Recv(-1)
			if err != nil {
				break
			}
			conn.Send(data)
		}
	})
	go s.Run()
	time.Sleep(simpleTimeout)
}

func startTCPKeyCrtServer(addr string) {
	s, _ := gtcp.NewServerKeyCrt(addr, crtFile, keyFile, func(conn *gtcp.Conn) {
		defer conn.Close()
		for {
			data, err := conn.Recv(-1)
			if err != nil {
				break
			}
			conn.Send(data)
		}
	})
	go s.Run()
	time.Sleep(simpleTimeout)
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
		recv, err := conn.SendRecv(sendData, -1)
		t.AssertNil(err)
		t.Assert(recv, sendData)
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

	gtest.C(t, func(t *gtest.T) {
		recv, err := gtcp.SendRecv(addr, sendData, -1)
		t.AssertNil(err)
		t.Assert(sendData, recv)
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
		recv, err := conn.SendRecv(sendData, -1)
		t.AssertNil(err)
		t.Assert(recv, sendData)
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

		conn, err := gtcp.NewConnTLS(addr, &tls.Config{
			InsecureSkipVerify: true,
			Certificates: []tls.Certificate{
				tls.Certificate{},
			},
		})
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

	gtest.C(t, func(t *gtest.T) {
		conn, err := gtcp.NewConn(addr)
		t.AssertNil(err)
		t.AssertNE(conn, nil)
		err = conn.Send(sendData, gtcp.Retry{Count: 1})
		t.AssertNil(err)
		recv, err := conn.Recv(-1)
		t.AssertNil(err)
		t.Assert(recv, sendData)
	})
}

func TestConn_SendWithTimeout(t *testing.T) {
	addr := getFreePortAddr()

	startTCPServer(addr)

	gtest.C(t, func(t *gtest.T) {
		conn, err := gtcp.NewConn(addr)
		t.AssertNil(err)
		t.AssertNE(conn, nil)
		err = conn.SendWithTimeout(sendData, time.Second, gtcp.Retry{Count: 1})
		t.AssertNil(err)
		recv, err := conn.Recv(-1)
		t.AssertNil(err)
		t.Assert(recv, sendData)
	})
}

func TestConn_SendRecv(t *testing.T) {
	addr := getFreePortAddr()

	startTCPServer(addr)

	gtest.C(t, func(t *gtest.T) {
		conn, err := gtcp.NewConn(addr)
		t.AssertNil(err)
		t.AssertNE(conn, nil)
		recv, err := conn.SendRecv(sendData, -1, gtcp.Retry{Count: 1})
		t.AssertNil(err)
		t.Assert(recv, sendData)
	})
}

func TestConn_SendRecvWithTimeout(t *testing.T) {
	addr := getFreePortAddr()

	startTCPServer(addr)

	gtest.C(t, func(t *gtest.T) {
		conn, err := gtcp.NewConn(addr)
		t.AssertNil(err)
		t.AssertNE(conn, nil)
		recv, err := conn.SendRecvWithTimeout(sendData, -1, time.Second, gtcp.Retry{Count: 1})
		t.AssertNil(err)
		t.Assert(recv, sendData)
	})
}

func TestConn_RecvWithTimeout(t *testing.T) {
	addr := getFreePortAddr()

	startTCPServer(addr)

	gtest.C(t, func(t *gtest.T) {
		conn, err := gtcp.NewConn(addr)
		t.AssertNil(err)
		t.AssertNE(conn, nil)
		conn.Send(sendData)
		recv, err := conn.RecvWithTimeout(-1, time.Second, gtcp.Retry{Count: 1})
		t.AssertNil(err)
		t.Assert(recv, sendData)
	})
}

func TestConn_RecvLine(t *testing.T) {
	addr := getFreePortAddr()

	startTCPServer(addr)

	gtest.C(t, func(t *gtest.T) {
		conn, err := gtcp.NewConn(addr)
		t.AssertNil(err)
		t.AssertNE(conn, nil)
		data := []byte("hello\n")
		conn.Send(data)
		recv, err := conn.RecvLine(gtcp.Retry{Count: 1})
		t.AssertNil(err)
		splitData := gstr.Split(string(data), "\n")
		t.Assert(recv, splitData[0])
	})
}

func TestConn_RecvTill(t *testing.T) {
	addr := getFreePortAddr()

	startTCPServer(addr)

	gtest.C(t, func(t *gtest.T) {
		conn, err := gtcp.NewConn(addr)
		t.AssertNil(err)
		t.AssertNE(conn, nil)
		conn.Send(sendData)
		recv, err := conn.RecvTill([]byte("hello"), gtcp.Retry{Count: 1})
		t.AssertNil(err)
		t.Assert(recv, sendData)
	})
}

func TestConn_SetDeadline(t *testing.T) {
	addr := getFreePortAddr()

	startTCPServer(addr)

	gtest.C(t, func(t *gtest.T) {
		conn, err := gtcp.NewConn(addr)
		t.AssertNil(err)
		t.AssertNE(conn, nil)
		conn.SetDeadline(time.Time{})
		err = conn.Send(sendData, gtcp.Retry{Count: 1})
		t.AssertNil(err)
		recv, err := conn.Recv(-1)
		t.AssertNil(err)
		t.Assert(recv, sendData)
	})
}

func TestConn_SetReceiveBufferWait(t *testing.T) {
	addr := getFreePortAddr()

	startTCPServer(addr)

	gtest.C(t, func(t *gtest.T) {
		conn, err := gtcp.NewConn(addr)
		t.AssertNil(err)
		t.AssertNE(conn, nil)
		conn.SetReceiveBufferWait(time.Millisecond * 100)
		err = conn.Send(sendData, gtcp.Retry{Count: 1})
		t.AssertNil(err)
		recv, err := conn.Recv(-1)
		t.AssertNil(err)
		t.Assert(recv, sendData)
	})
}

func TestNewNetConnKeyCrt(t *testing.T) {
	addr := getFreePortAddr()

	startTCPKeyCrtServer(addr)

	gtest.C(t, func(t *gtest.T) {
		conn, err := gtcp.NewNetConnKeyCrt(addr, crtFile, keyFile, time.Second)
		t.AssertNil(conn)
		t.AssertNE(err, nil)
	})
}

func TestSend(t *testing.T) {
	addr := getFreePortAddr()

	startTCPServer(addr)

	gtest.C(t, func(t *gtest.T) {
		err := gtcp.Send(addr, sendData, gtcp.Retry{Count: 1})
		t.AssertNil(err)
	})
}

func TestSendRecv(t *testing.T) {
	addr := getFreePortAddr()

	startTCPServer(addr)

	gtest.C(t, func(t *gtest.T) {
		recv, err := gtcp.SendRecv(addr, sendData, -1)
		t.AssertNil(err)
		t.Assert(recv, sendData)
	})
}

func TestSendWithTimeout(t *testing.T) {
	addr := getFreePortAddr()

	startTCPServer(addr)

	time.Sleep(simpleTimeout)

	gtest.C(t, func(t *gtest.T) {
		err := gtcp.SendWithTimeout("127.0.0.1:99999", sendData, time.Millisecond*500)
		t.AssertNE(err, nil)
		err = gtcp.SendWithTimeout(addr, sendData, time.Millisecond*500)
		t.AssertNil(err)
	})
}

func TestSendRecvWithTimeout(t *testing.T) {
	addr := getFreePortAddr()

	startTCPServer(addr)

	time.Sleep(simpleTimeout)

	gtest.C(t, func(t *gtest.T) {
		recv, err := gtcp.SendRecvWithTimeout("127.0.0.1:99999", sendData, -1, time.Millisecond*500)
		t.AssertNil(recv)
		t.AssertNE(err, nil)
		recv, err = gtcp.SendRecvWithTimeout(addr, sendData, -1, time.Millisecond*500)
		t.AssertNil(err)
		t.Assert(recv, sendData)
	})
}

func TestSendPkg(t *testing.T) {
	addr := getFreePortAddr()

	startTCPPkgServer(addr)

	time.Sleep(simpleTimeout)

	gtest.C(t, func(t *gtest.T) {
		err := gtcp.SendPkg(addr, sendData)
		t.AssertNil(err)
		err = gtcp.SendPkg("127.0.0.1:99999", sendData)
		t.AssertNE(err, nil)
	})

	gtest.C(t, func(t *gtest.T) {
		err := gtcp.SendPkg(addr, sendData)
		t.AssertNil(err)
		err = gtcp.SendPkg(addr, sendData)
		t.AssertNil(err)
	})
}

func TestSendRecvPkg(t *testing.T) {
	addr := getFreePortAddr()

	startTCPPkgServer(addr)

	time.Sleep(simpleTimeout)

	gtest.C(t, func(t *gtest.T) {
		err := gtcp.SendPkg(addr, sendData)
		t.AssertNil(err)
		_, err = gtcp.SendRecvPkg("127.0.0.1:99999", sendData)
		t.AssertNE(err, nil)
	})

	gtest.C(t, func(t *gtest.T) {
		err := gtcp.SendPkg(addr, sendData)
		t.AssertNil(err)
		recv, err := gtcp.SendRecvPkg(addr, sendData)
		t.AssertNil(err)
		t.Assert(recv, sendData)
	})
}

func TestSendPkgWithTimeout(t *testing.T) {
	addr := getFreePortAddr()

	startTCPPkgServer(addr)

	time.Sleep(simpleTimeout)

	gtest.C(t, func(t *gtest.T) {
		err := gtcp.SendPkg(addr, sendData)
		t.AssertNil(err)
		err = gtcp.SendPkgWithTimeout("127.0.0.1:99999", sendData, time.Second)
		t.AssertNE(err, nil)
	})

	gtest.C(t, func(t *gtest.T) {
		err := gtcp.SendPkg(addr, sendData)
		t.AssertNil(err)
		err = gtcp.SendPkgWithTimeout(addr, sendData, time.Second)
		t.AssertNil(err)
	})
}

func TestSendRecvPkgWithTimeout(t *testing.T) {
	addr := getFreePortAddr()

	startTCPPkgServer(addr)

	time.Sleep(simpleTimeout)

	gtest.C(t, func(t *gtest.T) {
		err := gtcp.SendPkg(addr, sendData)
		t.AssertNil(err)
		_, err = gtcp.SendRecvPkgWithTimeout("127.0.0.1:99999", sendData, time.Second)
		t.AssertNE(err, nil)
	})

	gtest.C(t, func(t *gtest.T) {
		err := gtcp.SendPkg(addr, sendData)
		t.AssertNil(err)
		recv, err := gtcp.SendRecvPkgWithTimeout(addr, sendData, time.Second)
		t.AssertNil(err)
		t.Assert(recv, sendData)
	})
}

func TestNewServer(t *testing.T) {
	addr := getFreePortAddr()

	gtest.C(t, func(t *gtest.T) {
		s := gtcp.NewServer(addr, func(conn *gtcp.Conn) {
			defer conn.Close()
			for {
				data, err := conn.Recv(-1)
				if err != nil {
					break
				}
				conn.Send(data)
			}
		}, "NewServer")
		defer s.Close()
		go s.Run()
		t.AssertNE(s, nil)

		time.Sleep(simpleTimeout)

		recv, err := gtcp.SendRecv(addr, sendData, -1)
		t.AssertNil(err)
		t.Assert(recv, sendData)
	})
}

func TestGetServer(t *testing.T) {
	addr := getFreePortAddr()

	gtest.C(t, func(t *gtest.T) {
		s := gtcp.GetServer("GetServer")
		defer s.Close()
		go s.Run()

		t.AssertNE(s, nil)
		t.Assert(s.GetAddress(), "")
	})

	gtest.C(t, func(t *gtest.T) {
		gtcp.NewServer(addr, func(conn *gtcp.Conn) {
			defer conn.Close()
			for {
				data, err := conn.Recv(-1)
				if err != nil {
					break
				}
				conn.Send(data)
			}
		}, "NewServer")

		s := gtcp.GetServer("NewServer")
		defer s.Close()
		t.AssertNE(s, nil)
		go s.Run()

		time.Sleep(simpleTimeout)

		recv, err := gtcp.SendRecv(addr, sendData, -1)
		t.AssertNil(err)
		t.Assert(recv, sendData)
	})
}

func TestServer_SetAddress(t *testing.T) {
	addr := getFreePortAddr()

	gtest.C(t, func(t *gtest.T) {
		s := gtcp.NewServer("", func(conn *gtcp.Conn) {
			defer conn.Close()
			for {
				data, err := conn.Recv(-1)
				if err != nil {
					break
				}
				conn.Send(data)
			}
		})
		defer s.Close()
		t.Assert(s.GetAddress(), "")
		s.SetAddress(addr)
		go s.Run()

		time.Sleep(simpleTimeout)

		recv, err := gtcp.SendRecv(addr, sendData, -1)
		t.AssertNil(err)
		t.Assert(recv, sendData)
	})
}

func TestServer_SetHandler(t *testing.T) {
	addr := getFreePortAddr()

	gtest.C(t, func(t *gtest.T) {
		s := gtcp.NewServer(addr, nil)
		defer s.Close()
		s.SetHandler(func(conn *gtcp.Conn) {
			defer conn.Close()
			for {
				data, err := conn.Recv(-1)
				if err != nil {
					break
				}
				conn.Send(data)
			}
		})
		go s.Run()

		time.Sleep(simpleTimeout)

		recv, err := gtcp.SendRecv(addr, sendData, -1)
		t.AssertNil(err)
		t.Assert(recv, sendData)
	})
}

func TestServer_Run(t *testing.T) {
	addr := getFreePortAddr()

	gtest.C(t, func(t *gtest.T) {
		s := gtcp.NewServer(addr, func(conn *gtcp.Conn) {
			defer conn.Close()
			for {
				data, err := conn.Recv(-1)
				if err != nil {
					break
				}
				conn.Send(data)
			}
		})
		defer s.Close()
		go s.Run()

		time.Sleep(simpleTimeout)

		recv, err := gtcp.SendRecv(addr, sendData, -1)
		t.AssertNil(err)
		t.Assert(recv, sendData)
	})

	gtest.C(t, func(t *gtest.T) {
		s := gtcp.NewServer(addr, nil)
		defer s.Close()
		go func() {
			err := s.Run()
			t.AssertNE(err, nil)
		}()
	})
}
