package graft

import (
    "time"
    "g/net/gip"
    "log"
    "g/net/gtcp"
    "net"
    "fmt"
    "g/net/gscanner"
    "encoding/json"
    "io"
)

// 集群协议通信接口回调函数
func (n *Node) raftTcpHandler(conn net.Conn) {
    msg := n.recieve(conn)
    if msg == nil {
        return
    }
    fromip, _ := gip.ParseAddress(conn.RemoteAddr().String())
    switch msg.Act {
        // 上线通知
        case "hi":
            n.Peers[fromip] = msg.From.RaftInfo.Role
            n.send(conn, "hi2", nil)
            //fmt.Println("add peer:", fromip, "to", n.Ip, ", remote", conn.RemoteAddr(), ", local", conn.LocalAddr())
    }
}

// 集群数据同步接口回调函数
func (n *Node) repliTcpHandler(conn net.Conn) {

}

// 局域网扫描回调函数，类似广播消息
func (n *Node) scannerRaftCallback(conn net.Conn) {
    fromip, _ := gip.ParseAddress(conn.RemoteAddr().String())
    if fromip == n.Ip {
        //fmt.Println(fromip, "==", n.Ip)
        return
    }
    err := n.send(conn, "hi", nil)
    if err != nil {
        log.Println(err)
        return
    }

    msg := n.recieve(conn)
    switch msg.Act {
        // 收到通知，找到组织
        case "hi2":
            n.Peers[fromip] = msg.From.RaftInfo.Role
            //fmt.Println("add peer from scan:", fromip, "to", n.Ip)
    }
}

// 获取数据
// @todo 数据量接收超过1024 byte时的处理
func (n *Node) recieve(conn net.Conn) *Msg {
    buffersize := 1024
    buffer     := make([]byte, buffersize)
    count, err := conn.Read(buffer)
    if err != nil {
        if err != io.EOF {
            log.Println("conn.Read", err)
        }
        return nil
    }
    if count > 0 {
        var msg Msg
        err = json.Unmarshal(buffer[0:count], &msg)
        if err != nil {
            log.Println(err)
            return nil
        }
        return &msg
    }
    return nil
}

// 发送数据
func (n *Node) send(conn net.Conn, act string, data interface{}) error {
    var msg = Msg{act, data, *n.getMsgFromInfo()}
    s, err := json.Marshal(msg)
    if err != nil {
        return err
    }
    _, err = conn.Write([]byte(s))
    if err != nil {
        return err
    }
    return nil
}

// 获得当前节点进行数据通信时的来源信息结构
func (n *Node) getMsgFromInfo() *MsgFrom {
    var from = MsgFrom {
        Name : n.Name,
        Role : n.Role,
    }
    from.RaftInfo.Term = n.RaftInfo.Term
    from.RaftInfo.Role = n.RaftInfo.Role
    return &from
}

// 运行节点
func (n *Node) Run() {
    // 创建接口监听
    gtcp.NewServer(fmt.Sprintf("%s:%d", n.Ip, gCLUSTER_PORT_RAFT),  n.raftTcpHandler).Run()
    gtcp.NewServer(fmt.Sprintf("%s:%d", n.Ip, gCLUSTER_PORT_REPLI), n.repliTcpHandler).Run()
    // 通知上线
    n.sayHiToAll()
    // 测试
    n.showPeers()
}

// 测试使用，展示当前节点通信的主机列表
func (n *Node) showPeers() {
    time.AfterFunc(2 * time.Second, func(){
        fmt.Println(n.Ip, n.Peers)
        n.showPeers()
    })
}

// 向局域网内其他主机通知上线
func (n *Node) sayHiToAll() {
    segment := gip.GetSegment(n.Ip)
    if segment != "" {
        startIp := fmt.Sprintf("%s.1",   segment)
        endIp   := fmt.Sprintf("%s.255", segment)
        gscanner.New().SetTimeout(6*time.Second).ScanIp(startIp, endIp, gCLUSTER_PORT_RAFT, n.scannerRaftCallback)
    }
}