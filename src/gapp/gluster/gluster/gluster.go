/*
    使用raft算法处理集群的一致性
    已解决split brains造成的数据一致性问题：
        经典分区多节点脑裂问题：出现两个及以上的分区网络，网络之间无法相互连接
        复杂非分区三点脑裂问题：A-B-C，AC之间无法相互连接（B是双网卡），这样会造成A、C为leader，B为follower
        以上需要解决的是数据一致性问题，解决方案：检测集群必需为环形网络，剔除掉非环形网络的节点
 */

package gluster

import (
    "os"
    "g/core/types/gmap"
    "sync"
    "g/core/types/glist"
    "g/net/ghttp"
    "g/os/glog"
    "g/os/gfile"
    "net"
    "encoding/json"
    "time"
    "io"
    "g/util/gutil"
    "sort"
    "g/net/gip"
    "strings"
    "g/encoding/gmd5"
    "fmt"
    "g/encoding/gcompress"
    "g/os/gconsole"
)

const (
    gVERSION                        = "0.6"   // 当前版本
    gCOMPRESS_COMMUNICATION         = true    // 是否在通信时进行内容压缩
    gCOMPRESS_SAVING                = false   // 是否在存储时压缩内容
    // 集群端口定义
    gPORT_RAFT                      = 4166    // 集群协议通信接口
    gPORT_REPL                      = 4167    // 集群数据同步接口
    gPORT_API                       = 4168    // 服务器对外API接口
    gPORT_MONITOR                   = 4169    // 监控服务接口
    gPORT_WEBUI                     = 4170    // WEB管理界面

    // 节点状态
    gSTATUS_DEAD                    = 0
    gSTATUS_ALIVE                   = 1

    // 集群角色
    gROLE_SERVER                    = 0
    gROLE_CLIENT                    = 1
    gROLE_MONITOR                   = 2

    // RAFT角色
    gROLE_RAFT_FOLLOWER             = 0
    gROLE_RAFT_CANDIDATE            = 1
    gROLE_RAFT_LEADER               = 2

    // 超时时间设置
    gTCP_RETRY_COUNT                = 3       // TCP请求失败时的重试次数
    gTCP_READ_TIMEOUT               = 3000    // (毫秒)TCP链接读取超时
    gTCP_WRITE_TIMEOUT              = 3000    // (毫秒)TCP链接写入超时
    gELECTION_TIMEOUT               = 1000    // (毫秒)RAFT选举超时时间
    gELECTION_TIMEOUT_HEARTBEAT     = 500     // (毫秒)RAFT Leader统治维持心跳间隔
    gLOG_REPL_TIMEOUT_HEARTBEAT     = 1000    // (毫秒)数据同步检测心跳间隔(数据包括kv数据及service数据)
    gLOG_REPL_AUTOSAVE_INTERVAL     = 1000    // (毫秒)数据自动物理化保存的间隔
    gLOG_REPL_LOGCLEAN_INTERVAL     = 5000    // (毫秒)数据同步时的日志清理间隔
    gLOG_REPL_PEERS_INTERVAL        = 3000    // (毫秒)Peers节点信息同步(非完整同步)
    gSERVICE_HEALTH_CHECK_INTERVAL  = 2000    // (毫秒)健康检查默认间隔

    // RAFT操作
    gMSG_RAFT_HI                    = 110
    gMSG_RAFT_HI2                   = 120
    gMSG_RAFT_RESPONSE              = 130
    gMSG_RAFT_HEARTBEAT             = 140
    gMSG_RAFT_I_AM_LEADER           = 150
    gMSG_RAFT_SPLIT_BRAINS_CHECK    = 160
    gMSG_RAFT_SPLIT_BRAINS_UNSET    = 170
    gMSG_RAFT_SCORE_REQUEST         = 180
    gMSG_RAFT_SCORE_COMPARE_REQUEST = 190
    gMSG_RAFT_SCORE_COMPARE_FAILURE = 200
    gMSG_RAFT_SCORE_COMPARE_SUCCESS = 210

    // 数据同步操作
    gMSG_REPL_DATA_SET                      = 310
    gMSG_REPL_DATA_REMOVE                   = 320
    gMSG_REPL_INCREMENTAL_UPDATE            = 330
    gMSG_REPL_COMPLETELY_UPDATE             = 340
    gMSG_REPL_HEARTBEAT                     = 350
    gMSG_REPL_RESPONSE                      = 360
    gMSG_REPL_PEERS_UPDATE                  = 370
    gMSG_REPL_NEED_UPDATE_LEADER            = 375
    gMSG_REPL_NEED_UPDATE_FOLLOWER          = 380
    gMSG_REPL_CONFIG_FROM_FOLLOWER          = 383
    gMSG_REPL_SERVICE_COMPLETELY_UPDATE     = 385
    gMSG_REPL_SERVICE_NEED_UPDATE_LEADER    = 390
    gMSG_REPL_SERVICE_NEED_UPDATE_FOLLOWER  = 400

    // API相关
    gMSG_API_PEERS_ADD                      = 520
    gMSG_API_PEERS_REMOVE                   = 530
    gMSG_API_SERVICE_SET                    = 540
    gMSG_API_SERVICE_REMOVE                 = 550
)

// 消息
type Msg struct {
    Head int
    Body string
    Info NodeInfo
}

// 服务器节点信息
type Node struct {
    mutex               sync.RWMutex             // 通用锁，可以考虑不同的变量使用不同的锁以提高读写效率

    Group               string                   // 集群名称
    Id                  string                   // 节点ID(根据算法自动生成的集群唯一名称)
    Name                string                   // 节点主机名称
    Ip                  string                   // 主机节点的ip，由通信的时候进行填充，
                                                 // 一个节点可能会有多个IP，这里保存最近通信的那个，节点唯一性识别使用的是Name字段
    CfgFilePath         string                   // 配置文件绝对路径
    CfgReplicated       bool                     // 本地配置对象是否已同步到leader(配置同步需要注意覆盖问题)
    Peers               *gmap.StringInterfaceMap // 集群所有的节点信息(ip->节点信息)，不包含自身
    Role                int                      // 集群角色
    RaftRole            int                      // RAFT角色
    MinNode             int                      // 组成集群的最小节点数量
    Leader              *NodeInfo                // Leader节点信息
    Score               int64                    // 选举比分
    ScoreCount          int                      // 选举比分的节点数
    ElectionDeadline    int64                    // 选举超时时间点
    isInDataReplication bool                     // 是否正在数据同步过程中

    LastLogId           int64                    // 最后一次保存log的id，用以数据一致性判断
    LogCount            int                      // 物理化保存的日志总数量，用于数据一致性判断
    LastSavedLogId      int64                    // 最后一次物理化log的id，用以物理化保存识别
    LastServiceLogId    int64                    // 最后一次保存的service id号，用以识别service信息同步
    LogList             *glist.SafeList          // leader日志列表，用以数据同步
    SavePath            string                   // 物理存储的本地数据目录绝对路径
    FileName            string                   // 数据文件名称(包含后缀)
    Service             *gmap.StringInterfaceMap // 存储的服务配置表
    ServiceForApi       *gmap.StringInterfaceMap // 用于提高Service API响应的冗余map变量，内容与Service成员变量相同，但结构不同
    DataMap             *gmap.StringStringMap    // 存储的K-V哈希表
}

// 服务对象
type Service struct {
    Name  string
    Type  string
    Node  *gmap.StringInterfaceMap
}

// 用以可直接json化处理的Service数据结构
type ServiceStruct struct {
    Name  string                 `json:"name"`
    Type  string                 `json:"type"`
    Node  map[string]interface{} `json:"node"`
}

// 用于KV API接口的对象
type NodeApiKv struct {
    ghttp.Controller
    node *Node
}

// 用于Node API接口的对象
type NodeApiNode struct {
    ghttp.Controller
    node *Node
}

// 用于Service API接口的对象
type NodeApiService struct {
    ghttp.Controller
    node *Node
}

// 用于Service 负载均衡API接口的对象
type NodeApiBalance struct {
    ghttp.Controller
    node *Node
}

// 用于Monitor WebUI对象
type MonitorWebUI struct {
    ghttp.Controller
    node *Node
}

// 节点信息
type NodeInfo struct {
    Group            string
    Id               string
    Name             string
    Ip               string
    Status           int
    Role             int
    RaftRole         int
    Score            int64
    ScoreCount       int
    LastLogId        int64
    LogCount         int
    LastServiceLogId int64
    LastActiveTime   int64  // 上一次活跃的时间毫秒(活跃包含：新增、心跳)，该数据用于Peer数据表中
    Version          string // 节点的版本
}

// 数据保存结构体
type SaveInfo struct {
    LastLogId        int64
    LogCount         int
    LogList          []LogEntry
    LastServiceLogId int64
    Service          map[string]ServiceStruct
    Peers            map[string]interface{}
    DataMap          map[string]string
}

// 日志记录项
type LogEntry struct {
    Id               int64                  // 唯一ID
    Act              int
    Items            interface{}            // map[string]string或[]string
}

// 绑定本地IP并创建一个服务节点
func NewServer() *Node {
    // 主机名称
    hostname, err := os.Hostname()
    if err != nil {
        glog.Fatalln("getting local hostname failed:", err)
        return nil
    }
    node := Node {
        Id                  : nodeId(),
        Ip                  : "127.0.0.1",
        Name                : hostname,
        Role                : gROLE_SERVER,
        RaftRole            : gROLE_RAFT_FOLLOWER,
        MinNode             : 1,
        Leader              : nil,
        Peers               : gmap.NewStringInterfaceMap(),
        SavePath            : gfile.SelfDir(),
        FileName            : "gluster.db",
        LogList             : glist.NewSafeList(),
        Service             : gmap.NewStringInterfaceMap(),
        ServiceForApi       : gmap.NewStringInterfaceMap(),
        DataMap               : gmap.NewStringStringMap(),
        isInDataReplication : false,
    }
    ips, err := gip.IntranetIP()
    if err == nil && len(ips) == 1 {
        node.Ip = ips[0]
    }
    // 命令行操作绑定
    gconsole.BindHandle("getnode",    cmd_getnode)
    gconsole.BindHandle("addnode",    cmd_addnode)
    gconsole.BindHandle("delnode",    cmd_delnode)
    gconsole.BindHandle("getkv",      cmd_getkv)
    gconsole.BindHandle("addkv",      cmd_addkv)
    gconsole.BindHandle("delkv",      cmd_delkv)
    gconsole.BindHandle("getservice", cmd_getservice)
    gconsole.BindHandle("addservice", cmd_addservice)
    gconsole.BindHandle("delservice", cmd_delservice)

    return &node
}

// 生成节点的唯一ID(hostname+ips)
func nodeId() string {
    hostname, err := os.Hostname()
    if err != nil {
        glog.Fatalln("getting local hostname failed:", err)
    }
    ips, err      := gip.IntranetIP()
    if err != nil {
        glog.Fatalln("getting local ips:", err)
    }
    // 如果有多个IP，那么将IP升序排序
    sort.Slice(ips, func(i, j int) bool { return ips[i] < ips[j] })
    return strings.ToUpper(gmd5.EncodeString(fmt.Sprintf("%s/%s", hostname, strings.Join(ips, ","))))
}

// 获取数据
func Receive(conn net.Conn) []byte {
    conn.SetReadDeadline(time.Now().Add(gTCP_READ_TIMEOUT * time.Millisecond))
    retry      := 0
    buffersize := 1024
    data       := make([]byte, 0)
    for {
        buffer      := make([]byte, buffersize)
        length, err := conn.Read(buffer)
        if err != nil {
            if retry > gTCP_RETRY_COUNT - 1 {
                break;
            }
            if err != io.EOF {
                //glog.Println("receive err:", err, "retry:", retry)
            }
            retry ++
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
    if gCOMPRESS_COMMUNICATION {
        return gcompress.UnZlib(data)
    }
    return data
}

// 发送数据
func Send(conn net.Conn, data []byte) error {
    conn.SetReadDeadline(time.Now().Add(gTCP_WRITE_TIMEOUT * time.Millisecond))
    retry := 0
    for {
        if gCOMPRESS_COMMUNICATION {
            data = gcompress.Zlib(data)
        }
        _, err := conn.Write(data)
        if err != nil {
            if retry > gTCP_RETRY_COUNT - 1 {
                return err
            }
            //glog.Println("data send:", err, "try:", retry)
            retry ++
            time.Sleep(100 * time.Millisecond)
        } else {
            return nil
        }
    }
}

// 获取Msg
func RecieveMsg(conn net.Conn) *Msg {
    data := Receive(conn)
    //glog.Println(string(data))
    if data != nil && len(data) > 0 {
        var msg Msg
        err := json.Unmarshal(data, &msg)
        if err != nil {
            glog.Println("receive msg parse err:", err)
            return nil
        }
        ip, _      := gip.ParseAddress(conn.RemoteAddr().String())
        msg.Info.Ip = ip
        return &msg
    }
    return nil
}

// 发送Msg
func SendMsg(conn net.Conn, head int, body string) error {
    var msg = Msg{
        Head : head,
        Body : body,
    }
    s, err := json.Marshal(msg)
    if err != nil {
        glog.Println("send msg parse err:", err)
        return err
    }
    return Send(conn, s)
}
