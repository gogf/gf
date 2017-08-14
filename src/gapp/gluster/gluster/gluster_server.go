package gluster

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
    "g/util/grand"
    "g/net/ghttp"
)

// 局域网扫描回调函数，类似广播消息
func (n *Node) scannerRaftCallback(conn net.Conn) {
    fromip, _ := gip.ParseAddress(conn.RemoteAddr().String())
    if fromip == n.Ip {
        //log.Println(fromip, "==", n.Ip)
        return
    }
    err := n.sendMsg(conn, gMSG_RAFT_HI, "")
    if err != nil {
        log.Println(err)
        return
    }

    msg := n.receiveMsg(conn)
    if msg != nil && msg.Head == gMSG_RAFT_HI2 {
        n.updatePeerInfo(msg.Info)
        if msg.Info.Role == gROLE_LEADER {
            log.Println(n.Ip, "scanner: found leader", fromip)
            n.setLeader(fromip)
            n.setRole(gROLE_FOLLOWER)
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
    var msg = Msg { head, body, *n.getNodeInfo() }
    s, err := json.Marshal(msg)
    if err != nil {
        log.Println("send msg parse err:", err)
        return err
    }
    return n.send(conn, s)
}

// 获得TCP链接
func (n *Node) getConn(ip string, port int) net.Conn {
    conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", ip, port), 3 * time.Second)
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

// 运行节点
func (n *Node) Run() {
    // 初始化节点数据
    n.restoreDataFromFile()

    // 创建接口监听
    go gtcp.NewServer(fmt.Sprintf("%s:%d", n.Ip, gPORT_RAFT),  n.raftTcpHandler).Run()
    go gtcp.NewServer(fmt.Sprintf("%s:%d", n.Ip, gPORT_REPL),  n.replTcpHandler).Run()
    go func() {
        ips, _  := gip.IntranetIP()
        address := fmt.Sprintf("%s:%d", n.Ip, gPORT_API)
        if len(ips) == 1 {
            address = fmt.Sprintf(":%d", gPORT_API)
        }
        api := ghttp.NewServerByAddr(address)
        api.BindController("/kv",      &NodeApiKv{node: n})
        api.BindController("/node",    &NodeApiNode{node: n})
        api.BindController("/service", &NodeApiService{node: n})
        api.Run()
    }()

    // 通知上线（这里采用局域网扫描的方式进行广播通知）
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
    // 服务健康检查
    go n.serviceHealthCheckHandler()

    // 测试
    //go n.show()
}

// 测试使用，展示当前节点通信的主机列表
func (n *Node) show() {
    gtime.SetInterval(1 * time.Second, func() bool{
        //log.Println(n.Ip + ":", n.getScoreCount(), n.getScore(), n.getLeader(), n.getRole())
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

// 获取当前节点的信息
func (n *Node) getNodeInfo() *NodeInfo {
    return &NodeInfo {
        Name             : n.Name,
        Ip               : n.Ip,
        Status           : gSTATUS_ALIVE,
        Role             : n.getRole(),
        Score            : n.getScore(),
        ScoreCount       : n.getScoreCount(),
        LastLogId        : n.getLastLogId(),
        LastServiceLogId : n.getLastServiceLogId(),
        LastHeartbeat    : gtime.Millisecond(),
        Version          : gVERSION,
    }
}

func (n *Node) getLeader() string {
    n.mutex.RLock()
    r := n.Leader
    n.mutex.RUnlock()
    return r
}

func (n *Node) getRole() int {
    n.mutex.RLock()
    r := n.Role
    n.mutex.RUnlock()
    return r
}

func (n *Node) getScore() int64 {
    n.mutex.RLock()
    r := n.Score
    n.mutex.RUnlock()
    return r
}

func (n *Node) getScoreCount() int {
    n.mutex.RLock()
    r := n.ScoreCount
    n.mutex.RUnlock()
    return r
}

func (n *Node) getLastLogId() int64 {
    n.mutex.RLock()
    r := n.LastLogId
    n.mutex.RUnlock()
    return r
}

func (n *Node) getLastSavedLogId() int64 {
    n.mutex.Lock()
    r := n.LastSavedLogId
    n.mutex.Unlock()
    return r
}

func (n *Node) getLastServiceLogId() int64 {
    n.mutex.Lock()
    r := n.LastServiceLogId
    n.mutex.Unlock()
    return r
}

func (n *Node) getStatusInReplication() bool {
    n.mutex.RLock()
    r := isInReplication
    n.mutex.RUnlock()
    return r
}

func (n *Node) getElectionDeadline() int64 {
    n.mutex.RLock()
    r := n.ElectionDeadline
    n.mutex.RUnlock()
    return r
}

// 获取数据文件的绝对路径
func (n *Node) getDataFilePath() string {
    n.mutex.RLock()
    path := n.SavePath + gfile.Separator + n.Ip + "." + n.FileName
    n.mutex.RUnlock()
    return path
}

// 设置数据保存目录路径
func (n *Node) SetSavePath(path string) {
    n.mutex.Lock()
    n.SavePath = path
    n.mutex.Unlock()
}

// 添加比分节
func (n *Node) addScore(s int64) {
    n.mutex.Lock()
    n.Score += s
    n.mutex.Unlock()
}

// 添加比分节点数
func (n *Node) addScoreCount() {
    n.mutex.Lock()
    n.ScoreCount++
    n.mutex.Unlock()
}

// 重置为候选者，并初始化投票给自己
func (n *Node) resetAsCandidate() {
    n.mutex.Lock()
    n.Role       = gROLE_CANDIDATE
    n.Leader     = ""
    n.Score      = 0
    n.ScoreCount = 0
    n.mutex.Unlock()
}

// 重置为选民，并清空选票信息
func (n *Node) resetAsFollower() {
    n.mutex.Lock()
    n.Role      = gROLE_FOLLOWER
    n.Leader    = ""
    n.Score      = 0
    n.ScoreCount = 0
    n.mutex.Unlock()
}

func (n *Node) setRole(role int) {
    n.mutex.Lock()
    n.Role = role
    n.mutex.Unlock()
}

func (n *Node) setLeader(ip string) {
    n.mutex.Lock()
    n.Leader = ip
    n.mutex.Unlock()
}

func (n *Node) SetMonitor(ip string) {
    n.mutex.Lock()
    n.Monitor = ip
    n.mutex.Unlock()
}

func (n *Node) SetFileName(name string) {
    n.mutex.Lock()
    n.FileName = name
    n.mutex.Unlock()
}

func (n *Node) setLastLogId(id int64) {
    n.mutex.Lock()
    n.LastLogId = id
    n.mutex.Unlock()
}

func (n *Node) setLastSavedLogId(id int64) {
    n.mutex.Lock()
    n.LastSavedLogId = id
    n.mutex.Unlock()
}

func (n *Node) setLastServiceLogId(id int64) {
    n.mutex.Lock()
    n.LastServiceLogId = id
    n.mutex.Unlock()
}

func (n *Node) setStatusInReplication(status bool ) {
    n.mutex.Lock()
    isInReplication = status
    n.mutex.Unlock()
}

func (n *Node) setService(m *gmap.StringInterfaceMap) {
    if m == nil {
        return
    }
    n.mutex.Lock()
    n.Service = m
    n.mutex.Unlock()
}

func (n *Node) setPeers(m *gmap.StringInterfaceMap) {
    if m == nil {
        return
    }
    n.mutex.Lock()
    n.Peers = m
    n.mutex.Unlock()
}

func (n *Node) setKVMap(m *gmap.StringStringMap) {
    if m == nil {
        return
    }
    n.mutex.Lock()
    n.KVMap = m
    n.mutex.Unlock()
}

// 更新节点信息
func (n *Node) updatePeerInfo(info NodeInfo) {
    n.Peers.Set(info.Ip, info)
}

func (n *Node) updatePeerStatus(ip string, status int) {
    r := n.Peers.Get(ip)
    if r != nil {
        info       := r.(NodeInfo)
        info.Status = status
        if status == gSTATUS_ALIVE {
            info.LastHeartbeat = gtime.Millisecond()
        }
        n.Peers.Set(ip, info)
    }
    //if status == gSTATUS_DEAD {
    //    log.Println(ip, "was dead")
    //}
}

// 更新选举截止时间
func (n *Node) updateElectionDeadline() {
    n.mutex.Lock()
    n.ElectionDeadline = gtime.Millisecond() + int64(grand.Rand(gELECTION_TIMEOUT_MIN, gELECTION_TIMEOUT_MAX))
    n.mutex.Unlock()
}