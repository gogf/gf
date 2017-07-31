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
    "g/util/gtime"
    "g/util/gutil"
)

// 局域网扫描回调函数，类似广播消息
func (n *Node) scannerRaftCallback(conn net.Conn) {
    fromip, _ := gip.ParseAddress(conn.RemoteAddr().String())
    if fromip == n.Ip {
        //log.Println(fromip, "==", n.Ip)
        return
    }
    err := n.sendMsg(conn, "hi", nil)
    if err != nil {
        log.Println(err)
        return
    }

    msg := n.recieveMsg(conn)
    if msg.Head == "hi2" {
        n.Peers.Set(fromip, msg.From.RaftInfo.Role)
        if msg.From.RaftInfo.Role == gRAFT_ROLE_LEADER {
            n.setRaftLeader(fromip)
        }
    }
}

// 获取数据
func (n *Node) recieve(conn net.Conn) []byte {
    //log.Println(conn.LocalAddr().String(), "recieve from", conn.RemoteAddr().String())
    try        := 0
    buffersize := 1024
    data       := make([]byte, 0)
    for {
        buffer      := make([]byte, buffersize)
        length, err := conn.Read(buffer)
        if err != nil {
            if err != io.EOF {
                log.Println("node recieve:", err, "try:", try)
            }
            if try > 2 {
                break;
            }
            try ++
            time.Sleep(100 * time.Millisecond)
        } else {
            if length == buffersize {
                data = gutil.MergeSlice(data, buffer)
            } else {
                data = gutil.MergeSlice(data, buffer[0:length])
                break;
            }
        }
    }
    return data
}

// 获取Msg
func (n *Node) recieveMsg(conn net.Conn) *Msg {
    response := n.recieve(conn)
    if response != nil && len(response) > 0 {
        var msg Msg
        err := json.Unmarshal(response, &msg)
        if err != nil {
            log.Println(err)
            return nil
        }
        return &msg
    }
    return nil
}

// 发送数据
func (n *Node) send(conn net.Conn, data []byte) error {
    //log.Println(conn.LocalAddr().String(), "send to", conn.RemoteAddr().String())
    try := 0
    for {
        _, err := conn.Write(data)
        if err != nil {
            log.Println("data send:", err, "try:", try)
            if try > 2 {
                return err
            }
            try ++
            time.Sleep(100 * time.Millisecond)
        } else {
            return nil
        }
    }
}

// 发送Msg
func (n *Node) sendMsg(conn net.Conn, head string, body interface{}) error {
    var msg = Msg{head, body, *n.getMsgFromInfo()}
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

    // 测试
    go n.show()
}

// 测试使用，展示当前节点通信的主机列表
func (n *Node) show() {
    gtime.SetInterval(4 * time.Second, func() bool{
        log.Println(n.Name, n.Ip, n.Peers.M, n.RaftInfo)
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