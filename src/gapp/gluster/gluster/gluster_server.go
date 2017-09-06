package gluster

import (
    "time"
    "g/net/gip"
    "g/net/gtcp"
    "net"
    "fmt"
    "encoding/json"
    "g/util/gtime"
    "g/core/types/gmap"
    "g/os/gfile"
    "g/net/ghttp"
    "g/os/gconsole"
    "g/encoding/gjson"
    "g/os/glog"
    "strings"
    "os"
    "errors"
)

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
    ip, _  := gip.ParseAddress(conn.LocalAddr().String())
    info   := n.getNodeInfo()
    info.Ip = ip
    s, err := json.Marshal(Msg { head, body, *info })
    if err != nil {
        glog.Println("send msg parse err:", err)
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

// 运行节点
func (n *Node) Run() {
    // 命令行操作
    if gconsole.Value.Get(1) != "" {
        gconsole.AutoRun()
        os.Exit(0)
        return
    }
    // 显示当前节点信息
    glog.Println("Host Id  :", n.Id)
    glog.Println("Host Name:", n.Name)

    // 读取配置文件
    n.initFromCfg()

    // 初始化节点数据
    n.restoreDataFromFile()

    // 创建接口监听
    go gtcp.NewServer(fmt.Sprintf(":%d", gPORT_RAFT),  n.raftTcpHandler).Run()
    go gtcp.NewServer(fmt.Sprintf(":%d", gPORT_REPL),  n.replTcpHandler).Run()
    go func() {
        api := ghttp.NewServerByAddr(fmt.Sprintf(":%d", gPORT_API))
        api.BindController("/kv",      &NodeApiKv{node: n})
        api.BindController("/node",    &NodeApiNode{node: n})
        api.BindController("/service", &NodeApiService{node: n})
        api.BindController("/balance", &NodeApiBalance{node: n})
        api.Run()
    }()

    // 通知上线（这里采用局域网扫描的方式进行广播通知）
    //go n.sayHiToAll()
    //time.Sleep(2 * time.Second)
    // 配置同步
    go n.replicateConfigToLeader()
    // 选举超时检查
    go n.electionHandler()
    // 心跳保持及存活性检查
    go n.heartbeatHandler()
    // 日志同步处理
    go n.replicationHandler()
    // 本地日志存储处理
    go n.logAutoSavingHandler()
    // 服务健康检查
    go n.serviceHealthCheckHandler()

    // 测试
    //go n.show()
}

// 读取配置文件内容
func (n *Node) initFromCfg() {
    // 获取命令行指定的配置文件路径，如果不存在，那么使用默认路径的配置文件
    // 默认路径为gcluster执行文件的同一目录下的gluster.json文件
    cfgpath := gconsole.Option.Get("cfg")
    if cfgpath == "" {
        cfgpath = gfile.SelfDir() + gfile.Separator + "gluster.json"
    } else {
        if !gfile.Exists(cfgpath) {
            glog.Error(cfgpath, "does not exist")
            return
        }
    }
    if !gfile.Exists(cfgpath) {
        return
    }
    n.CfgFilePath = cfgpath

    j := gjson.DecodeToJson(string(gfile.GetContents(cfgpath)))
    if j == nil {
        glog.Fatalln("config file decoding failed(surely a json format?), exit")
    }
    glog.Println("initializing from", cfgpath)
    // 集群名称
    n.Group = j.GetString("Group")
    // 集群角色
    n.Role  = j.GetInt("Role")
    if n.Role < 0 || n.Role > 2 {
        glog.Fatalln("invalid role setting, exit")
    }
    // 数据保存路径(请保证运行gcluster的用户有权限写入)
    savepath := j.GetString("SavePath")
    if savepath != "" {
        if !gfile.Exists(savepath) {
            gfile.Mkdir(savepath)
        }
        if !gfile.IsWritable(savepath) {
            glog.Fatalln(savepath, "is not writable for saving data")
        }
        n.SetSavePath(strings.TrimRight(savepath, gfile.Separator))
    }
    // 日志保存路径
    logpath := j.GetString("LogPath")
    if logpath != "" {
        if !gfile.Exists(logpath) {
            gfile.Mkdir(logpath)
        }
        if !gfile.IsWritable(logpath) {
            glog.Fatalln(logpath, "is not writable for saving log")
        }
        glog.SetLogPath(logpath)
    }
    // (可选)初始化节点列表，包含自定义的所需添加的服务器IP或者域名列表
    peers := j.GetArray("Peers")
    if peers != nil {
        for _, v := range peers {
            ip := v.(string)
            if ip == n.Ip {
                continue
            }
            go func(ip string) {
                if !n.sayHi(ip) {
                    n.updatePeerInfo(NodeInfo{Id: ip, Ip: ip})
                }
            }(ip)
        }
    }
}

// 将本地配置信息同步到leader
func (n *Node) replicateConfigToLeader() {
    for !n.CfgReplicated {
        if n.getRaftRole() != gROLE_RAFT_LEADER {
            if n.getLeader() != nil {
                if gfile.Exists(n.CfgFilePath) {
                    glog.Println("replicate config to leader")
                    err := n.SendToLeader(gMSG_REPL_CONFIG_FROM_FOLLOWER, gPORT_REPL, gfile.GetContents(n.CfgFilePath))
                    if err == nil {
                        n.CfgReplicated = true
                        glog.Println("replicate config to leader, done")
                    } else {
                        glog.Error(err)
                    }
                } else {
                    n.CfgReplicated = true
                }
            }
        } else {
            n.CfgReplicated = true
        }
        time.Sleep(100 * time.Millisecond)
    }
}

// (测试使用)展示当前节点通信的主机列表
func (n *Node) show() {
    gtime.SetInterval(1 * time.Second, func() bool{
        //glog.Println(n.Ip + ":", n.getScoreCount(), n.getScore(), n.getLeader(), n.getRaftRole())
        glog.Println(n.getLastLogId(), n.getLastServiceLogId())
        return true
    })
}

// 获取当前节点的信息
func (n *Node) getNodeInfo() *NodeInfo {
    return &NodeInfo {
        Group            : n.Group,
        Id               : n.Id,
        Ip               : n.Ip,
        Name             : n.Name,
        Status           : gSTATUS_ALIVE,
        Role             : n.Role,
        RaftRole         : n.getRaftRole(),
        Score            : n.getScore(),
        ScoreCount       : n.getScoreCount(),
        LastLogId        : n.getLastLogId(),
        LogCount         : n.getLogCount(),
        LastActiveTime   : gtime.Millisecond(),
        LastServiceLogId : n.getLastServiceLogId(),
        Version          : gVERSION,
    }
}

// 向leader发送操作请求
func (n *Node) SendToLeader(head int, port int, body string) error {
    leader := n.getLeader()
    if leader == nil {
        return errors.New(fmt.Sprintf("leader not found, please try again after leader election done, request head: %d", head))
    }
    conn := n.getConn(leader.Ip, port)
    if conn == nil {
        return errors.New("could not connect to leader: " + leader.Ip)
    }
    defer conn.Close()
    err := n.sendMsg(conn, head, body)
    if err != nil {
        return errors.New("sending request error: " + err.Error())
    } else {
        msg := n.receiveMsg(conn)
        if msg != nil && ((port == gPORT_RAFT && msg.Head != gMSG_RAFT_RESPONSE) || (port == gPORT_REPL && msg.Head != gMSG_REPL_RESPONSE)) {
            return errors.New("handling request error")
        }
    }
    return nil
}

// 通过IP向一个节点发送消息并建立双方联系
func (n *Node) sayHi(ip string) bool {
    if ip == n.Ip {
        return false
    }
    conn := n.getConn(ip, gPORT_RAFT)
    if conn == nil {
        return false
    }
    defer conn.Close()
    // 如果是本地同一节点通信，那么移除掉
    if n.checkConnInLocalNode(conn) {
        n.Peers.Remove(ip)
        return false
    }
    err := n.sendMsg(conn, gMSG_RAFT_HI, "")
    if err != nil {
        return false
    }
    msg := n.receiveMsg(conn)
    if msg != nil && msg.Head == gMSG_RAFT_HI2 {
        n.updatePeerInfo(msg.Info)
        if msg.Info.RaftRole == gROLE_RAFT_LEADER && n.Leader == nil && n.getRaftRole() != gROLE_RAFT_LEADER {
            n.setLeader(&msg.Info)
            n.setRaftRole(gROLE_RAFT_FOLLOWER)
        }
    }
    return true
}

// 向局域网内其他主机通知上线
func (n *Node) sayHiToLocalLan() {
    segment := gip.GetSegment(n.Ip)
    if segment == "" {
        glog.Fatalln("invalid listening ip given")
        return
    }
    for i := 1; i < 256; i++ {
        go n.sayHi(fmt.Sprintf("%s.%d", segment, i))
    }
}

// 检查链接是否属于本地的一个链接(即：自己链接自己)
func (n *Node) checkConnInLocalNode(conn net.Conn) bool {
    localip,  _ := gip.ParseAddress(conn.LocalAddr().String())
    remoteip, _ := gip.ParseAddress(conn.RemoteAddr().String())
    return localip == remoteip
}

// 获得Peers节点信息(包含自身)
func (n *Node) getAllPeers() *[]NodeInfo{
    list := make([]NodeInfo, 0)
    list  = append(list, *n.getNodeInfo())
    for _, v := range n.Peers.Values() {
        list = append(list, v.(NodeInfo))
    }
    return &list
}

func (n *Node) getIp() string {
    n.mutex.RLock()
    r := n.Ip
    n.mutex.RUnlock()
    return r
}

func (n *Node) getLeader() *NodeInfo {
    n.mutex.RLock()
    r := n.Leader
    n.mutex.RUnlock()
    return r
}

func (n *Node) getRaftRole() int {
    n.mutex.RLock()
    r := n.RaftRole
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

func (n *Node) getLogCount() int {
    n.mutex.RLock()
    r := n.LogCount
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
    r := n.isInDataReplication
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
    //path := n.SavePath + gfile.Separator + n.Id + "." + n.FileName
    path := n.SavePath + gfile.Separator + n.FileName
    n.mutex.RUnlock()
    return path
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

// 添加日志总数
func (n *Node) addLogCount() {
    n.mutex.Lock()
    n.ScoreCount++
    n.mutex.Unlock()
}

// 重置为候选者，并初始化投票给自己
func (n *Node) resetAsCandidate() {
    n.mutex.Lock()
    n.RaftRole   = gROLE_RAFT_CANDIDATE
    n.Leader     = nil
    n.Score      = 0
    n.ScoreCount = 0
    n.mutex.Unlock()
}

// 重置为选民，并清空选票信息
func (n *Node) resetAsFollower() {
    n.mutex.Lock()
    n.RaftRole   = gROLE_RAFT_FOLLOWER
    n.Leader     = nil
    n.Score      = 0
    n.ScoreCount = 0
    n.mutex.Unlock()
}

func (n *Node) setIp(ip string) {
    n.mutex.Lock()
    n.Ip = ip
    n.mutex.Unlock()
}

func (n *Node) setRaftRole(role int) {
    n.mutex.Lock()
    if n.RaftRole != role {
        glog.Printf("role changed from %s to %s\n", n.raftRoleName(n.RaftRole), n.raftRoleName(role))
    }
    n.RaftRole = role
    n.mutex.Unlock()

}

func (n *Node) setLeader(info *NodeInfo) {
    n.mutex.Lock()
    if n.Leader != nil {
        glog.Printf("leader changed from %s to %s\n", n.Leader.Name, info.Name)
    } else {
        glog.Println("set leader:", info.Name)
    }
    n.Leader = info
    n.mutex.Unlock()

}

// 设置数据保存目录路径
func (n *Node) SetSavePath(path string) {
    n.mutex.Lock()
    n.SavePath = path
    n.mutex.Unlock()
}

func (n *Node) setLastLogId(id int64) {
    n.mutex.Lock()
    n.LastLogId = id
    n.mutex.Unlock()
}

func (n *Node) setLogCount(count int) {
    n.mutex.Lock()
    n.LogCount = count
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
    n.isInDataReplication = status
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

func (n *Node) setServiceForApi(m *gmap.StringInterfaceMap) {
    if m == nil {
        return
    }
    n.mutex.Lock()
    n.ServiceForApi = m
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

func (n *Node) setDataMap(m *gmap.StringStringMap) {
    if m == nil {
        return
    }
    n.mutex.Lock()
    n.DataMap = m
    n.mutex.Unlock()
}

// 更新节点信息
func (n *Node) updatePeerInfo(info NodeInfo) {
    n.Peers.Set(info.Id, info)
    // 去掉初始化时写入的IP键名记录
    if n.Peers.Contains(info.Ip) {
        if info.Id != info.Ip {
            n.Peers.Remove(info.Ip)
        }
    }
}

func (n *Node) updatePeerStatus(Id string, status int) {
    r := n.Peers.Get(Id)
    if r != nil {
        info       := r.(NodeInfo)
        info.Status = status
        if info.LastActiveTime == 0 || status == gSTATUS_ALIVE {
            info.LastActiveTime = gtime.Millisecond()
        }
        n.updatePeerInfo(info)
    }
}

// 更新选举截止时间
// 改进：固定时间进行比分，看谁的比分更多
func (n *Node) updateElectionDeadline() {
    n.mutex.Lock()
    n.ElectionDeadline = gtime.Millisecond() + gELECTION_TIMEOUT
    n.mutex.Unlock()
}

// 将RAFT角色字段转换为可读的字符串
func (n *Node) raftRoleName(role int) string {
    switch role {
        case gROLE_RAFT_FOLLOWER:  return "follower"
        case gROLE_RAFT_CANDIDATE: return "candidate"
        case gROLE_RAFT_LEADER:    return "leader"
    }
    return "unknown"
}

