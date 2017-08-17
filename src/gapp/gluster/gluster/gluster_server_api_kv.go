// 返回格式统一：
// {result:1, message:"", data:""}

package gluster

import (
    "strings"
    "g/net/ghttp"
    "g/encoding/gjson"
    "reflect"
)


// K-V 查询
func (this *NodeApiKv) GET(r *ghttp.ClientRequest, w *ghttp.ServerResponse) {
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
func (this *NodeApiKv) PUT(r *ghttp.ClientRequest, w *ghttp.ServerResponse) {
    this.POST(r, w)
}

// K-V 修改
// @todo 需要使用同步请求机制保证成功处理，并返回真实的处理结果
func (this *NodeApiKv) POST(r *ghttp.ClientRequest, w *ghttp.ServerResponse) {
    data := r.GetRaw()
    if data == "" {
        w.ResponseJson(0, "invalid input", nil)
        return
    }
    items := gjson.Decode(&data)
    if items == nil {
        w.ResponseJson(0, "invalid data type: json decoding failed", nil)
        return
    }
    if "map[string]interface {}" != reflect.TypeOf(items).String() {
        w.ResponseJson(0, "invalid data type", nil)
        return
    }
    // 请求到leader
    conn := this.node.getConn(this.node.getLeader(), gPORT_REPL)
    if conn == nil {
        w.ResponseJson(0, "could not connect to leader: " + this.node.getLeader(), nil)
        return
    }
    err := this.node.sendMsg(conn, gMSG_REPL_SET, *gjson.Encode(items))
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

// K-V 删除
// @todo 需要使用同步请求机制保证成功处理，并返回真实的处理结果
func (this *NodeApiKv) DELETE(r *ghttp.ClientRequest, w *ghttp.ServerResponse) {
    method := strings.ToUpper(r.Method)
    data   := r.GetRaw()
    if data == "" {
        w.ResponseJson(0, "invalid input", nil)
        return
    }

    items := gjson.Decode(&data)
    if items == nil {
        w.ResponseJson(0, "invalid data type: json decoding failed", nil)
        return
    }

    if "[]interface {}" != reflect.TypeOf(items).String()  {
        w.ResponseJson(0, "invalid data type for " + method, nil)
        return
    }
    // 请求到leader
    conn := this.node.getConn(this.node.getLeader(), gPORT_REPL)
    if conn == nil {
        w.ResponseJson(0, "could not connect to leader: " + this.node.getLeader(), nil)
        return
    }
    err := this.node.sendMsg(conn, gMSG_REPL_REMOVE, *gjson.Encode(items))
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
