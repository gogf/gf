package graft

import (
    "time"
    "g/net/gip"
    "log"
    "g/net/gtcp"
    "net"
    "fmt"
    "g/net/gscanner"
    "g/encoding/gjson"
)

const (
    // 集群端口定义
    CLUSTER_PORT_RAFT    = 4166 // 集群协议通信接口
    CLUSTER_PORT_REPLI   = 4167 // 集群数据同步接口
    // 集群角色
    CLUSTER_ROLE_SERVER  = 0
    CLUSTER_ROLE_CLIENT  = 1
    CLUSTER_ROLE_MONITOR = 2
    // raft 角色
    RAFT_ROLE_FOLLOWER   = 0
    RAFT_ROLE_CANDIDATE  = 1
    RAFT_ROLE_LEADER     = 2
    // 超时时间设置
    ELECTION_TIMEOUT_MIN = 500  * time.Millisecond // 官方推荐 150ms - 300ms
    ELECTION_TIMEOUT_MAX = 1000 * time.Millisecond // 官方推荐 150ms - 300ms
)

// 通信消息结构体
type Msg struct {
    Act  string
    Data interface{}
}

// 节点结构体
type Node struct {
    Name     string            // 节点名称
    Role     int               // 集群角色
    Peers    map[string]Node   // 集群所有的节点(ip->信息)
    RaftInfo struct {
        Role       int             // raft角色
        Term       int             // 时间阶段
        Vote       Node            // 投票的节点
        Leader     Node            // Leader节点
        VoteCount  int             // 获得的选票数量
        TotalCount int             // 总共节点数
    }
}

// 集群协议通信接口回调函数
func (n *Node) raftTcpHandler(conn net.Conn) {
    buffer     := make([]byte, 1024)
    count, err := conn.Read(buffer)
    if err != nil {
        log.Println(err)
        return
    }
    data, err := gjson.Decode(&string(buffer[0:count]))
    if err != nil {
        log.Println(err)
        return
    }
    act := data.GetString("act")
    switch act {
    case "hi":

        
    }
}

// 集群数据同步接口回调函数
func (n *Node) repliTcpHandler(conn net.Conn) {

}

// 局域网扫描回调函数
func (n *Node) scannerRaftCallback(conn net.Conn) {
    data, err := gjson.Encode(Msg{ "hi", n })
    if err != nil {
        log.Println(err)
        return
    }
    _, err = conn.Write([]byte(data))
    if err != nil {
        log.Println(err)
        return
    }
}

// 创建一个节点对象
func New() *Node {
    return &Node{}
}

// 运行节点
func (n *Node) Run() {
    // 创建接口监听
    gtcp.NewServer(fmt.Sprintf(":%d", CLUSTER_PORT_RAFT),  n.raftTcpHandler).Run()
    gtcp.NewServer(fmt.Sprintf(":%d", CLUSTER_PORT_REPLI), n.repliTcpHandler).Run()
    // 通知上线
    n.sayHello()
}

// 向集群内其他主机通知上线
func (n *Node) sayHello() {
    ips, err := gip.IntranetIP()
    if err != nil {
        log.Println(err)
        return
    }
    if len(ips) < 1 {
        log.Println("empty lan ips")
        return
    }
    for _, ip := range ips {
        segment := gip.GetSegment(ip)
        if segment != "" {
            startIp := fmt.Sprintf("%s.1",   segment)
            endIp   := fmt.Sprintf("%s.255", segment)
            gscanner.New().SetTimeout(6*time.Second).ScanIp(startIp, endIp, CLUSTER_PORT_RAFT, n.scannerRaftCallback)
        }
    }
}