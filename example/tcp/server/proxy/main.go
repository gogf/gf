package main

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gtcp"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/os/gtimer"
)

const (
	AddressOfServer1 = ":8198"
	AddressOfServer2 = ":8199"
	UpStream         = "127.0.0.1:8198"
)

var (
	ctx = gctx.GetInitCtx()
)

// StartTCPServer1 starts Server1: A simple tcp server for demo.
// It reads the content from client connect and write it back to client.
func StartTCPServer1() {
	s := g.TCPServer(1)
	s.SetHandler(func(conn *gtcp.Conn) {
		defer conn.Close()
		for {
			data, err := conn.Recv(-1)
			if err != nil {
				g.Log().Errorf(ctx, `%+v`, err)
				break
			}
			if len(data) > 0 {
				err = conn.Send([]byte(fmt.Sprintf(`received: %s`, data)))
				if err != nil {
					g.Log().Errorf(ctx, `%+v`, err)
					break
				}
			}
		}
	})
	s.SetAddress(AddressOfServer1)
	s.Run()
}

// StartTCPServer2 starts Server2:
// All requests to Server2 are directly redirected to Server1.
func StartTCPServer2() {
	s := g.TCPServer(2)
	s.SetHandler(func(conn *gtcp.Conn) {
		defer conn.Close()
		// Each client connection associates an upstream connection.
		upstreamClient, err := gtcp.NewConn(UpStream)
		if err != nil {
			_, _ = conn.Write([]byte(fmt.Sprintf(
				`cannot connect to upstream "%s": %s`, UpStream, err.Error(),
			)))
			return
		}
		// Redirect the client connection reading and writing to upstream connection.
		for {
			go io.Copy(upstreamClient, conn)
			_, err = io.Copy(conn, upstreamClient)
			if err != nil {
				_, _ = conn.Write([]byte(fmt.Sprintf(
					`io.Copy to upstream "%s" failed: %s`, UpStream, err.Error(),
				)))
			}
		}
	})
	s.SetAddress(AddressOfServer2)
	s.Run()
}

func main() {
	go StartTCPServer1()
	go StartTCPServer2()
	time.Sleep(time.Second)
	gtimer.Add(ctx, time.Second, func(ctx context.Context) {
		address := fmt.Sprintf(`127.0.0.1%s`, AddressOfServer2)
		result, err := gtcp.SendRecv(address, []byte(gtime.Now().String()), -1)
		if err != nil {
			g.Log().Errorf(ctx, `send data failed: %+v`, err)
		}
		g.Log().Info(ctx, result)
	})
	g.Listen()
}
