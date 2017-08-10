// 返回格式统一：
// {result:1, message:"", data:""}

package graft

import (
    "strings"
    "g/net/ghttp"
    "g/encoding/gjson"
    "fmt"
)


// K-V 查询
func (this *NodeApiKv) GET(r *ghttp.Request, w *ghttp.ServerResponse) {
    k := r.GetRequestString("k")
    if k == "" {
        w.ResponseJson(1, "ok", *this.node.KVMap.Clone())
    } else {
        if this.node.KVMap.Contains(k) {
            w.ResponseJson(1, "ok", this.node.KVMap.Get(k))
        } else {
            w.ResponseJson(1, "ok", nil)
        }
    }
}

// K-V 新增
func (this *NodeApiKv) PUT(r *ghttp.Request, w *ghttp.ServerResponse) {
    this.DELETE(r, w)
}

// K-V 修改
func (this *NodeApiKv) POST(r *ghttp.Request, w *ghttp.ServerResponse) {
    this.DELETE(r, w)
}

// K-V 删除
func (this *NodeApiKv) DELETE(r *ghttp.Request, w *ghttp.ServerResponse) {
    method := strings.ToUpper(r.Method)
    data   := r.GetRaw()
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
    conn := this.node.getConn(this.node.getLeader(), gPORT_REPL)
    if conn == nil {
        w.ResponseJson(0, "could not connect to leader: " + this.node.getLeader(), nil)
        return
    }
    head := gMSG_REPL_SET
    if method == "DELETE" {
        head = gMSG_REPL_REMOVE
    }
    err   = this.node.sendMsg(conn, head, *gjson.Encode(items))
    if err != nil {
        w.ResponseJson(0, "sending request error: " + err.Error(), nil)
    } else {
        msg := this.node.receiveMsg(conn)
        if msg.Head != gMSG_REPL_RESPONSE {
            w.ResponseJson(0, "handling request error", nil)
        } else {
            w.ResponseJson(1, "ok", nil)
        }
    }
    conn.Close()
}
