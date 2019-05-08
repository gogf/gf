package main

import (
	"github.com/gogf/gf/g/net/gtcp"
	"github.com/gogf/gf/g/os/glog"
	"github.com/gogf/gf/g/os/gtimer"
	"github.com/gogf/gf/geg/net/gtcp/pkg_operations/common/funcs"
	"github.com/gogf/gf/geg/net/gtcp/pkg_operations/common/types"
	"time"
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
			case "hello":     onServerHello(conn, msg)
			case "doexit":    onServerDoExit(conn, msg)
			case "heartbeat": onServerHeartBeat(conn, msg)
			default:
				glog.Errorfln("invalid message: %v", msg)
				break
		}
	}
}

func onServerHello(conn *gtcp.Conn, msg *types.Msg) {
	glog.Printfln("hello response message from [%s]: %s", conn.RemoteAddr().String(), msg.Data)
}

func onServerHeartBeat(conn *gtcp.Conn, msg *types.Msg) {
	glog.Printfln("heartbeat from [%s]", conn.RemoteAddr().String())
}

func onServerDoExit(conn *gtcp.Conn, msg *types.Msg) {
	glog.Printfln("exit command from [%s]", conn.RemoteAddr().String())
	conn.Close()
}