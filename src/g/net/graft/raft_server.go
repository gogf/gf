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
    "g/util/gtime"
    "g/encoding/gjson"
)

// 局域网扫描回调函数，类似广播消息
func (n *Node) scannerRaftCallback(conn net.Conn) {
    fromip, _ := gip.ParseAddress(conn.RemoteAddr().String())
    if fromip == n.Ip {
        //log.Println(fromip, "==", n.Ip)
        return
    }
    err := n.sendMsg(conn, gRAFT_MSG_HEAD_HI, nil)
    if err != nil {
        log.Println(err)
        return
    }

    msg := n.receiveMsg(conn)
    if msg.Head == gRAFT_MSG_HEAD_HI2 {
        n.Peers.Set(fromip, msg.From.RaftInfo.Role)
        if msg.From.RaftInfo.Role == gRAFT_ROLE_LEADER {
            n.setRaftLeader(fromip)
        }
    }
}

// 获取数据
func (n *Node) receive(conn net.Conn) []byte {
    return Receive(conn)
}

// 获取Msg
func (n *Node) receiveMsg(conn net.Conn) *Msg {
    return RecieveMsg(conn)
}

// 发送数据
func (n *Node) send(conn net.Conn, data []byte) error {
    return Send(conn, data)
}

// 发送Msg
func (n *Node) sendMsg(conn net.Conn, head int, body interface{}) error {
    var msg = Msg{ head, *gjson.Encode(body), *n.getMsgFromInfo() }
    s, err := json.Marshal(msg)
    if err != nil {
        return err
    }
    return n.send(conn, s)
}

// 获得TCP链接
func (n *Node) getConn(ip string, port int) net.Conn {
    conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", ip, port), 6*time.Second)
    if err == nil {
        return conn
    }
    return nil
}

// 获得当前节点进行数据通信时的来源信息结构
func (n *Node) getMsgFromInfo() *MsgFrom {
    n.mutex.RLock()
    var from = MsgFrom {
        Name : n.Name,
        Role : n.Role,
    }
    n.mutex.RUnlock()
    from.RaftInfo.Term = n.RaftInfo.Term
    from.RaftInfo.Role = n.RaftInfo.Role
    return &from
}

// 运行节点
func (n *Node) Run() {
    // 创建接口监听
    gtcp.NewServer(fmt.Sprintf("%s:%d", n.Ip, gCLUSTER_PORT_RAFT),  n.raftTcpHandler).Run()
    gtcp.NewServer(fmt.Sprintf("%s:%d", n.Ip, gCLUSTER_PORT_REPLI), n.repliTcpHandler).Run()
    //gtcp.NewServer(fmt.Sprintf("%s:%d", n.Ip, gCLUSTER_PORT_API),   n.apiTcpHandler).Run()
    // 通知上线
    n.sayHiToAll()
    // 选举超时检查
    go n.electionHandler()
    // 心跳保持及存活性检查
    go n.heartbeatHandler()
    // 日志同步处理
    go n.logAutoReplicationHandler()
    // 本地日志存储处理
    go n.logAutoSavingHandler()

    // 测试
    go n.show()
}

// 测试使用，展示当前节点通信的主机列表
func (n *Node) show() {
    gtime.SetInterval(4 * time.Second, func() bool{
        // log.Println(n.Name, n.Ip, n.Peers.M, n.RaftInfo)
        log.Println(n.Ip, ":", n.getRaftLeader(), n.Peers.M, n.LogList.Len(), n.KVMap.M)
        return true
    })
}

// 向局域网内其他主机通知上线
func (n *Node) sayHiToAll() {
    segment := gip.GetSegment(n.Ip)
    if segment == "" {
        log.Fatalln("invalid listening ip given")
        return
    }
    startIp := fmt.Sprintf("%s.1",   segment)
    endIp   := fmt.Sprintf("%s.255", segment)
    gscanner.New().SetTimeout(3*time.Second).ScanIp(startIp, endIp, gCLUSTER_PORT_RAFT, n.scannerRaftCallback)
}