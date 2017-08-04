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
    "g/core/types/gmap"
    "g/os/gfile"
    "g/net/ghttp"
)

// 局域网扫描回调函数，类似广播消息
func (n *Node) scannerRaftCallback(conn net.Conn) {
    fromip, _ := gip.ParseAddress(conn.RemoteAddr().String())
    if fromip == n.Ip {
        //log.Println(fromip, "==", n.Ip)
        return
    }
    err := n.sendMsg(conn, gMSG_HEAD_HI, "")
    if err != nil {
        log.Println(err)
        return
    }

    msg := n.receiveMsg(conn)
    if msg.Head == gMSG_HEAD_HI2 {
        n.Peers.Set(fromip, gSTATUS_ALIVE)
        if msg.From.Role == gROLE_LEADER {
            log.Println(n.Ip, "scanner: found leader", fromip)
            n.setLeader(fromip)
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
func (n *Node) sendMsg(conn net.Conn, head int, body string) error {
    var msg = Msg{ head, body, *n.getMsgFromInfo() }
    s, err := json.Marshal(msg)
    if err != nil {
        return err
    }
    return n.send(conn, s)
}

// 获得TCP链接
func (n *Node) getConn(ip string, port int) net.Conn {
    conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", ip, port), 3*time.Second)
    if err == nil {
        return conn
    }
    return nil
}

// 通过连接池获取tcp链接，连接池地址是传入的conns
func (n *Node) getConnFromPool(ip string, port int, conns *gmap.StringInterfaceMap) net.Conn {
    var conn net.Conn
    if result := conns.Get(ip); result != nil {
        conn = result.(net.Conn)
    } else {
        conn = n.getConn(ip, port)
        if conn != nil {
            conns.Set(ip, conn)
        } else {
            conns.Remove(ip)
            return nil
        }
    }
    return conn
}

// 获得当前节点进行数据通信时的来源信息结构
func (n *Node) getMsgFromInfo() *MsgFrom {
    n.mutex.RLock()
    var from = MsgFrom {
        Name      : n.Name,
        Role      : n.Role,
        LastLogId : n.LastLogId,
        LogCount  : n.LogCount,
    }
    n.mutex.RUnlock()
    return &from
}

// 获取数据文件的绝对路径
func (n *Node) getDataFilePath() string {
    n.mutex.RLock()
    path := n.DataPath + gfile.Separator + n.Ip + ".graft.db"
    n.mutex.RUnlock()
    return path
}

// 设置数据保存目录路径
func (n *Node) SetDataPath(path string) {
    n.mutex.Lock()
    n.DataPath = path
    n.mutex.Unlock()
}

// 运行节点
func (n *Node) Run() {
    // 创建接口监听
    gtcp.NewServer(fmt.Sprintf("%s:%d", n.Ip, gPORT_RAFT),  n.raftTcpHandler).Run()
    gtcp.NewServer(fmt.Sprintf("%s:%d", n.Ip, gPORT_REPL),  n.replTcpHandler).Run()
    ips, _  := gip.IntranetIP()
    address := fmt.Sprintf("%s:%d", n.Ip, gPORT_API)
    if len(ips) == 1 {
        address = fmt.Sprintf(":%d", gPORT_API)
    }
    api := ghttp.NewServerByAddr(address)
    api.BindHandle("/kv",   n.kvApiHandler)
    api.BindHandle("/node", n.nodeApiHandler)
    api.Run()

    // 初始化节点数据
    n.restoreData()
    // 通知上线
    go n.sayHiToAll()
    time.Sleep(2 * time.Second)
    // 选举超时检查
    go n.electionHandler()
    // 心跳保持及存活性检查
    go n.heartbeatHandler()
    // 日志同步处理
    go n.logAutoReplicationHandler()
    // 本地日志存储处理
    go n.logAutoSavingHandler()

    // 测试
    //go n.show()
}

// 测试使用，展示当前节点通信的主机列表
func (n *Node) show() {
    gtime.SetInterval(4 * time.Second, func() bool{
        //log.Println(n.Ip + ":", n.getScoreCount(), n.getScore(), n.getLeader(), *n.Peers.Clone(), n.LogList.Len(), n.KVMap.M)
        log.Println(n.Ip + ":", n.getLeader(), n.getLastLogId(), *n.Peers.Clone(), n.LogList.Len(), n.KVMap.M)
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
    //log.Println(n.Ip, "say hi to all")
    gscanner.New().SetTimeout(6 * time.Second).ScanIp(startIp, endIp, gPORT_RAFT, n.scannerRaftCallback)
    //log.Println(n.Ip, "say hi to all done")
}