package gluster

import (
    "net"
)

// monitor通信接口回调函数
func (n *Node) monitorTcpHandler(conn net.Conn) {
    msg := n.receiveMsg(conn)
    if msg == nil || msg.Info.Group != n.Group {
        conn.Close()
        return
    }

    //这里不用自动关闭链接，由于链接有读取超时，当一段时间没有数据时会自动关闭
    n.raftTcpHandler(conn)
}
