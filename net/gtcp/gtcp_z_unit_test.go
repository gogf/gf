// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gtcp_test

import (
	"crypto/tls"
	"fmt"
	"testing"
	"time"

	"github.com/gogf/gf/v2/debug/gdebug"
	"github.com/gogf/gf/v2/net/gtcp"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/text/gstr"
)

var (
	simpleTimeout = time.Millisecond * 100
	sendData      = []byte("hello")
	invalidAddr   = "127.0.0.1:99999"
	crtFile       = gfile.Dir(gdebug.CallerFilePath()) + gfile.Separator + "testdata/server.crt"
	keyFile       = gfile.Dir(gdebug.CallerFilePath()) + gfile.Separator + "testdata/server.key"
)

func startTCPServer(addr string) *gtcp.Server {
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
	return s
}

func startTCPPkgServer(addr string) *gtcp.Server {
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
	return s
}

func startTCPTLSServer(addr string) *gtcp.Server {
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
	return s
}

func startTCPKeyCrtServer(addr string) *gtcp.Server {
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
	return s
}

func TestGetFreePorts(t *testing.T) {
	ports, _ := gtcp.GetFreePorts(2)
	gtest.C(t, func(t *gtest.T) {
		t.AssertGT(ports[0], 0)
		t.AssertGT(ports[1], 0)
	})

	startTCPServer(fmt.Sprintf("%s:%d", "127.0.0.1", ports[0]))

	gtest.C(t, func(t *gtest.T) {
		conn, err := gtcp.NewPoolConn(fmt.Sprintf("127.0.0.1:%d", ports[0]))
		t.AssertNil(err)
		defer conn.Close()
		result, err := conn.SendRecv(sendData, -1)
		t.AssertNil(err)
		t.Assert(result, sendData)
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
		result, err := gtcp.SendRecv(addr, sendData, -1)
		t.AssertNil(err)
		t.Assert(sendData, result)
	})
}

func TestNewConn(t *testing.T) {
	addr := gtcp.FreePortAddress

	gtest.C(t, func(t *gtest.T) {
		conn, err := gtcp.NewConn(addr, simpleTimeout)
		t.AssertNil(conn)
		t.AssertNE(err, nil)
	})

	gtest.C(t, func(t *gtest.T) {
		s := startTCPServer(gtcp.FreePortAddress)

		conn, err := gtcp.NewConn(s.GetListenedAddress(), simpleTimeout)
		t.AssertNil(err)
		t.AssertNE(conn, nil)
		defer conn.Close()
		result, err := conn.SendRecv(sendData, -1)
		t.AssertNil(err)
		t.Assert(result, sendData)
	})
}

// TODO
func TestNewConnTLS(t *testing.T) {
	addr := gtcp.FreePortAddress

	gtest.C(t, func(t *gtest.T) {
		conn, err := gtcp.NewConnTLS(addr, &tls.Config{})
		t.AssertNil(conn)
		t.AssertNE(err, nil)
	})

	gtest.C(t, func(t *gtest.T) {
		s := startTCPTLSServer(addr)

		conn, err := gtcp.NewConnTLS(s.GetListenedAddress(), &tls.Config{
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
	addr := gtcp.FreePortAddress

	gtest.C(t, func(t *gtest.T) {
		conn, err := gtcp.NewConnKeyCrt(addr, crtFile, keyFile)
		t.AssertNil(conn)
		t.AssertNE(err, nil)
	})

	gtest.C(t, func(t *gtest.T) {
		s := startTCPKeyCrtServer(addr)

		conn, err := gtcp.NewConnKeyCrt(s.GetListenedAddress(), crtFile, keyFile)
		t.AssertNil(conn)
		t.AssertNE(err, nil)
	})
}

func TestConn_Send(t *testing.T) {
	s := startTCPServer(gtcp.FreePortAddress)

	gtest.C(t, func(t *gtest.T) {
		conn, err := gtcp.NewConn(s.GetListenedAddress())
		t.AssertNil(err)
		t.AssertNE(conn, nil)
		err = conn.Send(sendData, gtcp.Retry{Count: 1})
		t.AssertNil(err)
		result, err := conn.Recv(-1)
		t.AssertNil(err)
		t.Assert(result, sendData)
	})
}

func TestConn_SendWithTimeout(t *testing.T) {
	s := startTCPServer(gtcp.FreePortAddress)

	gtest.C(t, func(t *gtest.T) {
		conn, err := gtcp.NewConn(s.GetListenedAddress())
		t.AssertNil(err)
		t.AssertNE(conn, nil)
		err = conn.SendWithTimeout(sendData, time.Second, gtcp.Retry{Count: 1})
		t.AssertNil(err)
		result, err := conn.Recv(-1)
		t.AssertNil(err)
		t.Assert(result, sendData)
	})
}

func TestConn_SendRecv(t *testing.T) {
	s := startTCPServer(gtcp.FreePortAddress)

	gtest.C(t, func(t *gtest.T) {
		conn, err := gtcp.NewConn(s.GetListenedAddress())
		t.AssertNil(err)
		t.AssertNE(conn, nil)
		result, err := conn.SendRecv(sendData, -1, gtcp.Retry{Count: 1})
		t.AssertNil(err)
		t.Assert(result, sendData)
	})
}

func TestConn_SendRecvWithTimeout(t *testing.T) {
	s := startTCPServer(gtcp.FreePortAddress)

	gtest.C(t, func(t *gtest.T) {
		conn, err := gtcp.NewConn(s.GetListenedAddress())
		t.AssertNil(err)
		t.AssertNE(conn, nil)
		result, err := conn.SendRecvWithTimeout(sendData, -1, time.Second, gtcp.Retry{Count: 1})
		t.AssertNil(err)
		t.Assert(result, sendData)
	})
}

func TestConn_RecvWithTimeout(t *testing.T) {
	s := startTCPServer(gtcp.FreePortAddress)

	gtest.C(t, func(t *gtest.T) {
		conn, err := gtcp.NewConn(s.GetListenedAddress())
		t.AssertNil(err)
		t.AssertNE(conn, nil)
		conn.Send(sendData)
		result, err := conn.RecvWithTimeout(-1, time.Second, gtcp.Retry{Count: 1})
		t.AssertNil(err)
		t.Assert(result, sendData)
	})
}

func TestConn_RecvLine(t *testing.T) {
	s := startTCPServer(gtcp.FreePortAddress)

	gtest.C(t, func(t *gtest.T) {
		conn, err := gtcp.NewConn(s.GetListenedAddress())
		t.AssertNil(err)
		t.AssertNE(conn, nil)
		data := []byte("hello\n")
		conn.Send(data)
		result, err := conn.RecvLine(gtcp.Retry{Count: 1})
		t.AssertNil(err)
		splitData := gstr.Split(string(data), "\n")
		t.Assert(result, splitData[0])
	})
}

func TestConn_RecvTill(t *testing.T) {
	s := startTCPServer(gtcp.FreePortAddress)

	gtest.C(t, func(t *gtest.T) {
		conn, err := gtcp.NewConn(s.GetListenedAddress())
		t.AssertNil(err)
		t.AssertNE(conn, nil)
		conn.Send(sendData)
		result, err := conn.RecvTill([]byte("hello"), gtcp.Retry{Count: 1})
		t.AssertNil(err)
		t.Assert(result, sendData)
	})
}

func TestConn_SetDeadline(t *testing.T) {
	s := startTCPServer(gtcp.FreePortAddress)

	gtest.C(t, func(t *gtest.T) {
		conn, err := gtcp.NewConn(s.GetListenedAddress())
		t.AssertNil(err)
		t.AssertNE(conn, nil)
		conn.SetDeadline(time.Time{})
		err = conn.Send(sendData, gtcp.Retry{Count: 1})
		t.AssertNil(err)
		result, err := conn.Recv(-1)
		t.AssertNil(err)
		t.Assert(result, sendData)
	})
}

func TestConn_SetReceiveBufferWait(t *testing.T) {
	s := startTCPServer(gtcp.FreePortAddress)

	gtest.C(t, func(t *gtest.T) {
		conn, err := gtcp.NewConn(s.GetListenedAddress())
		t.AssertNil(err)
		t.AssertNE(conn, nil)
		conn.SetReceiveBufferWait(time.Millisecond * 100)
		err = conn.Send(sendData, gtcp.Retry{Count: 1})
		t.AssertNil(err)
		result, err := conn.Recv(-1)
		t.AssertNil(err)
		t.Assert(result, sendData)
	})
}

func TestNewNetConnKeyCrt(t *testing.T) {
	addr := gtcp.FreePortAddress

	startTCPKeyCrtServer(addr)

	gtest.C(t, func(t *gtest.T) {
		conn, err := gtcp.NewNetConnKeyCrt(addr, "crtFile", keyFile, time.Second)
		t.AssertNil(conn)
		t.AssertNE(err, nil)
	})

	gtest.C(t, func(t *gtest.T) {
		conn, err := gtcp.NewNetConnKeyCrt(addr, crtFile, keyFile, time.Second)
		t.AssertNil(conn)
		t.AssertNE(err, nil)
	})
}

func TestSend(t *testing.T) {
	s := startTCPServer(gtcp.FreePortAddress)

	gtest.C(t, func(t *gtest.T) {
		err := gtcp.Send(invalidAddr, sendData, gtcp.Retry{Count: 1})
		t.AssertNE(err, nil)
	})

	gtest.C(t, func(t *gtest.T) {
		err := gtcp.Send(s.GetListenedAddress(), sendData, gtcp.Retry{Count: 1})
		t.AssertNil(err)
	})
}

func TestSendRecv(t *testing.T) {
	s := startTCPServer(gtcp.FreePortAddress)

	gtest.C(t, func(t *gtest.T) {
		result, err := gtcp.SendRecv(invalidAddr, sendData, -1)
		t.AssertNE(err, nil)
		t.Assert(result, nil)
	})

	gtest.C(t, func(t *gtest.T) {
		result, err := gtcp.SendRecv(s.GetListenedAddress(), sendData, -1)
		t.AssertNil(err)
		t.Assert(result, sendData)
	})
}

func TestSendWithTimeout(t *testing.T) {
	s := startTCPServer(gtcp.FreePortAddress)

	gtest.C(t, func(t *gtest.T) {
		err := gtcp.SendWithTimeout(invalidAddr, sendData, time.Millisecond*500)
		t.AssertNE(err, nil)
		err = gtcp.SendWithTimeout(s.GetListenedAddress(), sendData, time.Millisecond*500)
		t.AssertNil(err)
	})
}

func TestSendRecvWithTimeout(t *testing.T) {
	s := startTCPServer(gtcp.FreePortAddress)

	gtest.C(t, func(t *gtest.T) {
		result, err := gtcp.SendRecvWithTimeout(invalidAddr, sendData, -1, time.Millisecond*500)
		t.AssertNil(result)
		t.AssertNE(err, nil)
		result, err = gtcp.SendRecvWithTimeout(s.GetListenedAddress(), sendData, -1, time.Millisecond*500)
		t.AssertNil(err)
		t.Assert(result, sendData)
	})
}

func TestSendPkg(t *testing.T) {
	s := startTCPPkgServer(gtcp.FreePortAddress)

	gtest.C(t, func(t *gtest.T) {
		err := gtcp.SendPkg(s.GetListenedAddress(), sendData)
		t.AssertNil(err)
		err = gtcp.SendPkg(invalidAddr, sendData)
		t.AssertNE(err, nil)
	})

	gtest.C(t, func(t *gtest.T) {
		err := gtcp.SendPkg(s.GetListenedAddress(), sendData, gtcp.PkgOption{Retry: gtcp.Retry{Count: 3}})
		t.AssertNil(err)
		err = gtcp.SendPkg(s.GetListenedAddress(), sendData)
		t.AssertNil(err)
	})
}

func TestSendRecvPkg(t *testing.T) {
	s := startTCPPkgServer(gtcp.FreePortAddress)

	gtest.C(t, func(t *gtest.T) {
		err := gtcp.SendPkg(s.GetListenedAddress(), sendData)
		t.AssertNil(err)
		_, err = gtcp.SendRecvPkg(invalidAddr, sendData)
		t.AssertNE(err, nil)
	})

	gtest.C(t, func(t *gtest.T) {
		err := gtcp.SendPkg(s.GetListenedAddress(), sendData)
		t.AssertNil(err)
		result, err := gtcp.SendRecvPkg(s.GetListenedAddress(), sendData)
		t.AssertNil(err)
		t.Assert(result, sendData)
	})
}

func TestSendPkgWithTimeout(t *testing.T) {
	s := startTCPPkgServer(gtcp.FreePortAddress)

	gtest.C(t, func(t *gtest.T) {
		err := gtcp.SendPkg(s.GetListenedAddress(), sendData)
		t.AssertNil(err)
		err = gtcp.SendPkgWithTimeout(invalidAddr, sendData, time.Second)
		t.AssertNE(err, nil)
	})

	gtest.C(t, func(t *gtest.T) {
		err := gtcp.SendPkg(s.GetListenedAddress(), sendData)
		t.AssertNil(err)
		err = gtcp.SendPkgWithTimeout(s.GetListenedAddress(), sendData, time.Second)
		t.AssertNil(err)
	})
}

func TestSendRecvPkgWithTimeout(t *testing.T) {
	s := startTCPPkgServer(gtcp.FreePortAddress)

	gtest.C(t, func(t *gtest.T) {
		err := gtcp.SendPkg(s.GetListenedAddress(), sendData)
		t.AssertNil(err)
		_, err = gtcp.SendRecvPkgWithTimeout(invalidAddr, sendData, time.Second)
		t.AssertNE(err, nil)
	})

	gtest.C(t, func(t *gtest.T) {
		err := gtcp.SendPkg(s.GetListenedAddress(), sendData)
		t.AssertNil(err)
		result, err := gtcp.SendRecvPkgWithTimeout(s.GetListenedAddress(), sendData, time.Second)
		t.AssertNil(err)
		t.Assert(result, sendData)
	})
}

func TestNewServer(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := gtcp.NewServer(gtcp.FreePortAddress, func(conn *gtcp.Conn) {
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

		time.Sleep(simpleTimeout)

		result, err := gtcp.SendRecv(s.GetListenedAddress(), sendData, -1)
		t.AssertNil(err)
		t.Assert(result, sendData)
	})
}

func TestGetServer(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := gtcp.GetServer("GetServer")
		defer s.Close()
		go s.Run()

		t.Assert(s.GetAddress(), "")
	})

	gtest.C(t, func(t *gtest.T) {
		gtcp.NewServer(gtcp.FreePortAddress, func(conn *gtcp.Conn) {
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
		go s.Run()

		time.Sleep(simpleTimeout)

		result, err := gtcp.SendRecv(s.GetListenedAddress(), sendData, -1)
		t.AssertNil(err)
		t.Assert(result, sendData)
	})
}

func TestServer_SetAddress(t *testing.T) {
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
		s.SetAddress(gtcp.FreePortAddress)
		go s.Run()

		time.Sleep(simpleTimeout)

		result, err := gtcp.SendRecv(s.GetListenedAddress(), sendData, -1)
		t.AssertNil(err)
		t.Assert(result, sendData)
	})
}

func TestServer_SetHandler(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := gtcp.NewServer(gtcp.FreePortAddress, nil)
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

		result, err := gtcp.SendRecv(s.GetListenedAddress(), sendData, -1)
		t.AssertNil(err)
		t.Assert(result, sendData)
	})
}

func TestServer_Run(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := gtcp.NewServer(gtcp.FreePortAddress, func(conn *gtcp.Conn) {
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

		result, err := gtcp.SendRecv(s.GetListenedAddress(), sendData, -1)
		t.AssertNil(err)
		t.Assert(result, sendData)
	})

	gtest.C(t, func(t *gtest.T) {
		s := gtcp.NewServer(gtcp.FreePortAddress, nil)
		defer s.Close()
		go func() {
			err := s.Run()
			t.AssertNE(err, nil)
		}()
	})
}
