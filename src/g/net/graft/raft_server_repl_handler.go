// 数据同步需要注意的是：
// leader只有在通知完所有follower更新完数据之后，自身才会进行数据更新
// 因此leader
package graft

import (
    "net"
    "g/encoding/gjson"
    "time"
    "log"
    "sync"
    "g/core/types/gmap"
    "g/os/gfile"
    "g/core/types/gset"
)

// 用以识别节点当前是否正在数据同步中
var isInReplication bool

// 集群数据同步接口回调函数
func (n *Node) replTcpHandler(conn net.Conn) {
    msg := n.receiveMsg(conn)
    if msg == nil {
        //log.Println("receive nil")
        conn.Close()
        return
    }
    switch msg.Head {
        case gREPLI_MSG_HEAD_SET:
            fallthrough
        case gREPLI_MSG_HEAD_REMOVE:
            n.SetStatusInReplication(true)
            if n.getRole() == gROLE_LEADER {
                var item LogRequest
                pitem := &item
                body  := msg.Body.(string)
                gjson.DecodeTo(&body, pitem)
                log.Println(body)
                log.Println(item)
                if pitem != nil {
                    var entry = LogEntry {
                        Id    : time.Now().UnixNano(),
                        Act   : msg.Head,
                        Key   : item.Key,
                        Value : item.Value,
                    }
                    n.LogList.PushFront(entry)
                    n.LogChan <- struct{}{}
                }
            } else {
                var entry LogEntry
                log.Println("receiving log entry", entry)
                body := msg.Body.(string)
                gjson.DecodeTo(&body, &entry)
                n.saveLogEntry(entry)
            }
            n.SetStatusInReplication(false)
            n.sendMsg(conn, gREPLI_MSG_HEAD_LOG_REPL_RESPONSE, nil)

        // 数据同步自动检测
        case gREPLI_MSG_HEAD_LOG_REPL_HEARTBEAT:
            n.sendMsg(conn, gREPLI_MSG_HEAD_LOG_REPL_HEARTBEAT, nil)
    }

    n.replTcpHandler(conn)
}

// leader到其他节点的数据同步监听
func (n *Node) logAutoReplicationHandler() {
    var wg sync.WaitGroup
    // 初始化数据同步心跳检测
    n.logAutoReplicationCheckHandler()

    for {
        select {
            case <- n.LogChan:
                n.SetStatusInReplication(true)
                conns := gmap.NewStringInterfaceMap()
                for {
                    item := n.LogList.PopBack()
                    if item == nil {
                        break;
                    }
                    entry := item.(LogEntry)
                    log.Println("sending log entry", entry)
                    for ip, status := range n.Peers.M {
                        if status != gSTATUS_ALIVE {
                            continue
                        }
                        conn := n.getConnFromPool(ip, gPORT_REPL, conns)
                        if conn != nil {
                            wg.Add(1)
                            go func(conn net.Conn, entry LogEntry) {
                                if err := n.sendMsg(conn, entry.Act, entry); err != nil {
                                    log.Println(err)
                                    conn.Close()
                                    wg.Done()
                                    return
                                }
                                n.receiveMsg(conn)
                                wg.Done()
                            }(conn, entry)
                        }
                    }
                    wg.Wait()
                    // 当所有节点的请求处理后，再保存数据到自身
                    // 以便leader与follower之间的数据同步判断
                    n.saveLogEntry(entry)
                }
                n.SetStatusInReplication(false)
                for _, c := range conns.M {
                    c.(net.Conn).Close()
                }
        }
    }
}

// 日志自动同步检查，类似心跳
func (n *Node) logAutoReplicationCheckHandler() {
    conns := gset.NewStringSet()
    for {
        if n.getRole() == gROLE_LEADER {
            for ip, _ := range n.Peers.M {
                if conns.Contains(ip) {
                    continue
                }
                conn := n.getConn(ip, gPORT_REPL)
                if conn == nil {
                    conns.Remove(ip)
                    continue
                }
                conns.Add(ip)
                go func(ip string, conn net.Conn) {
                    for {
                        // 如果当前正在数据同步操作中，那么等待
                        for n.GetStatusInReplication() {
                            time.Sleep(100 * time.Millisecond)
                        }
                        if n.getRole() != gROLE_LEADER {
                            conn.Close()
                            conns.Remove(ip)
                            return
                        }
                        if err := n.sendMsg(conn, gREPLI_MSG_HEAD_LOG_REPL_HEARTBEAT, nil); err != nil {
                            log.Println(err)
                            conn.Close()
                            conns.Remove(ip)
                            return
                        }
                        msg := n.receiveMsg(conn)
                        if msg == nil {
                            n.Peers.Set(ip, gSTATUS_DEAD)
                            conns.Remove(ip)
                            conn.Close()
                            return
                        } else {
                            switch msg.Head {
                                default:
                                    time.Sleep(gLOG_REPL_TIMEOUT_HEARTBEAT * time.Millisecond)
                            }
                        }
                    }
                }(ip, conn)
            }
        }
        time.Sleep(100 * time.Millisecond)
    }
}

// 保存日志数据
func (n *Node) saveLogEntry(entry LogEntry) {
    switch entry.Act {
        case gREPLI_MSG_HEAD_SET:
            log.Println("setting log entry", entry)
            n.KVMap.Set(entry.Key, entry.Value)

        case gREPLI_MSG_HEAD_REMOVE:
            log.Println("removing log entry", entry)
            n.KVMap.Remove(entry.Key)
    }
    n.setLastLogId(entry.Id)
    n.addLogCount()
}

// 日志自动保存处理
func (n *Node) logAutoSavingHandler() {
    // 初始化KVMap
    path := n.getDataFilePath()
    if gfile.Exists(path) {
        data := gfile.GetContents(path)
        if data != nil {
            log.Println("initializing kvmap from data file")
            content := string(data)
            gjson.DecodeTo(&content, n.KVMap)
        }
    }
    // 循环监听
    for {
        if n.getLastLogId() != n.getLastSavedLogId() {
            log.Println("saving data to file")
            data := gjson.Encode(n.KVMap)
            gfile.PutContents(n.DataPath, *data)
            n.setLastSavedLogId(n.getLastLogId())
        } else {
            time.Sleep(100 * time.Millisecond)
        }
    }
}

func (n *Node) getLastLogId() int64 {
    n.mutex.Lock()
    r := n.LastLogId
    n.mutex.Unlock()
    return r
}

func (n *Node) getLastSavedLogId() int64 {
    n.mutex.Lock()
    r := n.LastSavedLogId
    n.mutex.Unlock()
    return r
}

func (n *Node) GetStatusInReplication() bool {
    n.mutex.RLock()
    r := isInReplication
    n.mutex.RUnlock()
    return r
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

func (n *Node) SetStatusInReplication(status bool ) {
    n.mutex.Lock()
    isInReplication = status
    n.mutex.Unlock()
}

func (n *Node) addLogCount() {
    n.mutex.Lock()
    n.LogCount++
    n.mutex.Unlock()
}
