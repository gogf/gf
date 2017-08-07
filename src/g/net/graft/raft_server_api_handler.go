// 返回格式统一：
// {result:1, message:"", data:""}

package graft

import (
    "strings"
    "g/net/ghttp"
    "g/encoding/gjson"
)


// K-V API管理
func (n *Node) kvApiHandler(r *ghttp.Request, w *ghttp.Response) {
    method := strings.ToUpper(r.Method)
    switch method {
        case "GET":
            k := r.GetRequestString("k")
            if k == "" {
                w.ResponseJson(1, "", *n.KVMap.Clone())
            } else {
                w.ResponseJson(1, "", n.KVMap.Get(k))
            }

        // @todo 类型检查
        case "PUT":
            fallthrough
        case "POST":
            fallthrough
        case "DELETE":
            data := r.GetRaw()
            if data == "" {
                w.ResponseJson(0, "invalid k-v input", nil)
                return
            }
            var items interface{}
            err := gjson.DecodeTo(&data, &items)
            if err != nil {
                w.ResponseJson(0, "invalid data type: " + err.Error(), nil)
                return
            }
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
                if msg.Head != gMSG_HEAD_LOG_REPL_RESPONSE {
                    w.ResponseJson(0, "handling request error", nil)
                } else {
                    w.ResponseJson(1, "", nil)
                }
            }
            conn.Close()
    }

}

// 节点信息API管理
func (n *Node) nodeApiHandler(r *ghttp.Request, w *ghttp.Response) {
    //io.WriteString(w, "hello\n")
}

