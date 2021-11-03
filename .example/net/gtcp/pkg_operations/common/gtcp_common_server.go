package main

import (
	"time"

	"github.com/gogf/gf/v2/.example/net/gtcp/pkg_operations/common/funcs"
	"github.com/gogf/gf/v2/.example/net/gtcp/pkg_operations/common/types"
	"github.com/gogf/gf/v2/net/gtcp"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/os/gtimer"
)

func main() {
	gtcp.NewServer("127.0.0.1:8999", func(conn *gtcp.Conn) {
		defer conn.Close()
		// 测试消息, 10秒后让客户端主动退出
		gtimer.SetTimeout(10*time.Second, func() {
			funcs.SendPkg(conn, "doexit")
		})
		for {
			msg, err := funcs.RecvPkg(conn)
			if err != nil {
				if err.Error() == "EOF" {
					glog.Print("client closed")
				}
				break
			}
			switch msg.Act {
			case "hello":
				onClientHello(conn, msg)
			case "heartbeat":
				onClientHeartBeat(conn, msg)
			default:
				glog.Errorf("invalid message: %v", msg)
				break
			}
		}
	}).Run()
}

func onClientHello(conn *gtcp.Conn, msg *types.Msg) {
	glog.Printf("hello message from [%s]: %s", conn.RemoteAddr().String(), msg.Data)
	funcs.SendPkg(conn, msg.Act, "Nice to meet you!")
}

func onClientHeartBeat(conn *gtcp.Conn, msg *types.Msg) {
	glog.Printf("heartbeat from [%s]", conn.RemoteAddr().String())
}
