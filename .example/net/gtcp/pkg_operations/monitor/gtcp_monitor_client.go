package main

import (
	"encoding/json"

	"github.com/gogf/gf/v2/.example/net/gtcp/pkg_operations/monitor/types"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gtcp"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/os/gtime"
)

func main() {
	// 数据上报客户端
	conn, err := gtcp.NewConn("127.0.0.1:8999")
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	// 使用JSON格式化数据字段
	info, err := json.Marshal(types.NodeInfo{
		Cpu:  float32(66.66),
		Host: "localhost",
		Ip: g.Map{
			"etho": "192.168.1.100",
			"eth1": "114.114.10.11",
		},
		MemUsed:  15560320,
		MemTotal: 16333788,
		Time:     int(gtime.Timestamp()),
	})
	if err != nil {
		panic(err)
	}
	// 使用 SendRecvPkg 发送消息包并接受返回
	if result, err := conn.SendRecvPkg(info); err != nil {
		if err.Error() == "EOF" {
			glog.Print("server closed")
		}
	} else {
		glog.Print(string(result))
	}
}
