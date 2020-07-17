package main

import (
	"time"

	"github.com/jin502437344/gf/.example/net/gtcp/pkg_operations/common/funcs"
	"github.com/jin502437344/gf/.example/net/gtcp/pkg_operations/common/types"
	"github.com/jin502437344/gf/net/gtcp"
	"github.com/jin502437344/gf/os/glog"
	"github.com/jin502437344/gf/os/gtimer"
)

func main() {
	conn, err := gtcp.NewConn("127.0.0.1:8999")
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	// 心跳消息
	gtimer.SetInterval(time.Second, func() {
		if err := funcs.SendPkg(conn, "heartbeat"); err != nil {
			panic(err)
		}
	})
	// 测试消息, 3秒后向服务端发送hello消息
	gtimer.SetTimeout(3*time.Second, func() {
		if err := funcs.SendPkg(conn, "hello", "My name's John!"); err != nil {
			panic(err)
		}
	})
	for {
		msg, err := funcs.RecvPkg(conn)
		if err != nil {
			if err.Error() == "EOF" {
				glog.Println("server closed")
			}
			break
		}
		switch msg.Act {
		case "hello":
			onServerHello(conn, msg)
		case "doexit":
			onServerDoExit(conn, msg)
		case "heartbeat":
			onServerHeartBeat(conn, msg)
		default:
			glog.Errorf("invalid message: %v", msg)
			break
		}
	}
}

func onServerHello(conn *gtcp.Conn, msg *types.Msg) {
	glog.Printf("hello response message from [%s]: %s", conn.RemoteAddr().String(), msg.Data)
}

func onServerHeartBeat(conn *gtcp.Conn, msg *types.Msg) {
	glog.Printf("heartbeat from [%s]", conn.RemoteAddr().String())
}

func onServerDoExit(conn *gtcp.Conn, msg *types.Msg) {
	glog.Printf("exit command from [%s]", conn.RemoteAddr().String())
	conn.Close()
}
