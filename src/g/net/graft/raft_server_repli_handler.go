package graft

import (
    "net"
)

// 集群数据同步接口回调函数
func (n *Node) repliTcpHandler(conn net.Conn) {
    msg := n.recieveMsg(conn)
    if msg == nil {
        return
    }
    switch msg.Head {
        case "set":

        case "remove":
    }
}
