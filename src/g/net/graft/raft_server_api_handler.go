// 返回格式统一：
// {result:1, message:"", data:""}

package graft

import (
    "strings"
    "g/net/ghttp"
    "g/encoding/gjson"
    "fmt"
)


// K-V API管理
func (n *Node) kvApiHandler(r *ghttp.Request, w *ghttp.Response) {
    method := strings.ToUpper(r.Method)
    switch method {
        case "GET":
            k := r.GetRequestString("k")
            if k == "" {
                w.ResponseJson(1, "ok", *n.KVMap.Clone())
            } else {
                if n.KVMap.Contains(k) {
                    w.ResponseJson(1, "ok", n.KVMap.Get(k))
                } else {
                    w.ResponseJson(1, "ok", nil)
                }
            }

        case "PUT":
            fallthrough
        case "POST":
            fallthrough
        case "DELETE":
            data := r.GetRaw()
            if data == "" {
                w.ResponseJson(0, "invalid input", nil)
                return
            }
            var items interface{}
            err := gjson.DecodeTo(&data, &items)
            if err != nil {
                w.ResponseJson(0, "invalid data type: " + err.Error(), nil)
                return
            }
            // 只允许map[string]interface{}和[]interface{}两种数据类型
            isSSMap  := "map[string]interface {}" == fmt.Sprintf("%T", items)
            isSArray := "[]interface {}" == fmt.Sprintf("%T", items)
            if !(isSSMap && (method == "PUT" || method == "POST")) && !(isSArray && method == "DELETE")  {
                w.ResponseJson(0, "invalid data type for " + method, nil)
                return
            }
            // 请求到leader
            conn := n.getConn(n.getLeader(), gPORT_REPL)
            if conn == nil {
                w.ResponseJson(0, "could not connect to leader: " + n.getLeader(), nil)
                return
            }
            head := gMSG_HEAD_SET
            if method == "DELETE" {
                head = gMSG_HEAD_REMOVE
            }
            err   = n.sendMsg(conn, head, *gjson.Encode(items))
            if err != nil {
                w.ResponseJson(0, "sending request error: " + err.Error(), nil)
            } else {
                msg := n.receiveMsg(conn)
                if msg.Head != gMSG_HEAD_REPL_RESPONSE {
                    w.ResponseJson(0, "handling request error", nil)
                } else {
                    w.ResponseJson(1, "ok", nil)
                }
            }
            conn.Close()

        default:
            w.ResponseJson(0, "unsupported method " + method, nil)
    }
}

// 节点信息API管理
func (n *Node) nodeApiHandler(r *ghttp.Request, w *ghttp.Response) {
    method := strings.ToUpper(r.Method)
    switch method {
        case "GET":
            conn := n.getConn(n.getLeader(), gPORT_RAFT)
            if conn == nil {
                w.ResponseJson(0, "could not connect to leader: " + n.getLeader(), nil)
                return
            }
            err := n.sendMsg(conn, gMSG_HEAD_PEERS_INFO, "")
            if err != nil {
                w.ResponseJson(0, "sending request error: " + err.Error(), nil)
            } else {
                var data interface{}
                msg := n.receiveMsg(conn)
                err  = gjson.DecodeTo(&msg.Body, &data)
                if err != nil {
                    w.ResponseJson(0, "received error from leader: " + err.Error(), nil)
                } else {
                    w.ResponseJson(1, "ok", data)
                }
            }
            conn.Close()

        case "PUT":
            fallthrough
        case "POST":
            fallthrough
        case "DELETE":
            data := r.GetRaw()
            if data == "" {
                w.ResponseJson(0, "invalid input", nil)
                return
            }
            var items interface{}
            err := gjson.DecodeTo(&data, &items)
            if err != nil {
                w.ResponseJson(0, "invalid data type: " + err.Error(), nil)
                return
            }
            // 只允许[]interface{}数据类型
            if "[]interface {}" != fmt.Sprintf("%T", items)  {
                w.ResponseJson(0, "invalid data type", nil)
                return
            }
            // 请求到leader
            conn := n.getConn(n.getLeader(), gPORT_RAFT)
            if conn == nil {
                w.ResponseJson(0, "could not connect to leader: " + n.getLeader(), nil)
                return
            }
            head := gMSG_HEAD_PEERS_ADD
            if method == "DELETE" {
                head = gMSG_HEAD_PEERS_REMOVE
            }
            err   = n.sendMsg(conn, head, *gjson.Encode(items))
            if err != nil {
                w.ResponseJson(0, "sending request error: " + err.Error(), nil)
            } else {
                msg := n.receiveMsg(conn)
                if msg.Head != gMSG_HEAD_RAFT_RESPONSE {
                    w.ResponseJson(0, "handling request error", nil)
                } else {
                    w.ResponseJson(1, "ok", nil)
                }
            }
            conn.Close()

        default:
            w.ResponseJson(0, "unsupported method " + method, nil)
    }
}

