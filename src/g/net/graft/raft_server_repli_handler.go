package graft

import (
    "net"
    "g/encoding/gjson"
    "time"
    "log"
    "sync"
    "g/core/types/gmap"
)

// 集群数据同步接口回调函数
func (n *Node) repliTcpHandler(conn net.Conn) {
    msg := n.receiveMsg(conn)
    if msg == nil {
        log.Println("receive nil")
        conn.Close()
        return
    }

    if n.getRaftRole() == gRAFT_ROLE_LEADER {
        var item *LogRequest
        body := msg.Body.(string)
        gjson.DecodeTo(&body, item)
        if item != nil {
            var entry = LogEntry {
                Id    : time.Now().UnixNano(),
                Act   : msg.Head,
                Key   : item.Key,
                Value : item.Value,
            }
            n.LogChan <- struct{}{}
            n.LogList.PushFront(entry)
            n.setLastLogId(entry.Id)
            n.addLogCount()
        }
    } else {
        var entry *LogEntry
        log.Println("receiving log entry", *entry)
        body := msg.Body.(string)
        gjson.DecodeTo(&body, entry)
        n.LogList.PushFront(entry)
        n.setLastLogId(entry.Id)
        n.addLogCount()
    }

    n.repliTcpHandler(conn)
}

// 日志自动同步节点
// 注意每一条数据的同步都是非异步执行
// 这里并不会检测数据是否同步成功，因为我们有plan b的机制进一步保证数据的一致性
func (n *Node) logAutoReplicationHandler() {
    var wg sync.WaitGroup
    for {
        select {
            // 事件通知进行数据发送同步
            case <- n.LogChan:
                conns := gmap.NewStringInterfaceMap()
                for {
                    item := n.LogList.PopBack()
                    if item == nil {
                        break;
                    }
                    entry := item.(LogEntry)
                    log.Println("sending log entry", entry)
                    ips := n.Peers.Keys()
                    for _, ip := range ips {
                        wg.Add(1)
                        conn, ok := conns.Get(ip)
                        if !ok {
                            conn = n.getConn(ip, gCLUSTER_PORT_REPLI)
                            conns.Set(ip, conn)
                        }
                        if conn == nil {
                            conns.Remove(ip)
                            continue
                        }
                        go func(conn net.Conn, entry LogEntry) {
                            n.sendMsg(conn, gREPLI_MSG_HEAD_SET, entry)
                            wg.Done()
                        }(conn.(net.Conn), entry)
                    }
                    wg.Wait()
                }
                for _, c := range conns.M {
                    c.(net.Conn).Close()
                }
        }
    }
}

// 日志自动保存处理
func (n *Node) logAutoSavingHandler() {
    for {
        v := n.LogList.PopBack()
        if v == nil {
            time.Sleep(100 * time.Millisecond)
        } else {
            entry := v.(LogEntry)
            switch entry.Act {
                case gREPLI_MSG_HEAD_SET:
                    log.Println("setting log entry", entry)
                    n.KVMap.Set(entry.Key, entry.Value)
                case gREPLI_MSG_HEAD_REMOVE:
                    log.Println("removing log entry", entry)
                    n.KVMap.Remove(entry.Key)
            }
        }
    }
}

func (n *Node) setLastLogId(id int64) {
    n.mutex.Lock()
    n.LastLogId = id
    n.mutex.Unlock()
}

func (n *Node) addLogCount() {
    n.mutex.Lock()
    n.LogCount++
    n.mutex.Unlock()
}
