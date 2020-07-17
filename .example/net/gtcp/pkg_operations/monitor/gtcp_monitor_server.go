package main

import (
	"encoding/json"

	"github.com/jin502437344/gf/.example/net/gtcp/pkg_operations/monitor/types"
	"github.com/jin502437344/gf/net/gtcp"
	"github.com/jin502437344/gf/os/glog"
)

func main() {
	// 服务端，接收客户端数据并格式化为指定数据结构，打印
	gtcp.NewServer("127.0.0.1:8999", func(conn *gtcp.Conn) {
		defer conn.Close()
		for {
			data, err := conn.RecvPkg()
			if err != nil {
				if err.Error() == "EOF" {
					glog.Println("client closed")
				}
				break
			}
			info := &types.NodeInfo{}
			if err := json.Unmarshal(data, info); err != nil {
				glog.Errorf("invalid package structure: %s", err.Error())
			} else {
				glog.Println(info)
				conn.SendPkg([]byte("ok"))
			}
		}
	}).Run()
}
