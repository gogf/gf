// 返回格式统一：
// {result:1, message:"", data:""}

package gluster

import (
    "strings"
    "g/net/ghttp"
    "g/encoding/gjson"
    "reflect"
)

// 节点信息API管理
func (this *NodeApiNode) GET(r *ghttp.Request, w *ghttp.ServerResponse) {
    conn := this.node.getConn(this.node.getLeader(), gPORT_RAFT)
    if conn == nil {
        w.ResponseJson(0, "could not connect to leader: " + this.node.getLeader(), nil)
        return
    }
    err := this.node.sendMsg(conn, gMSG_API_PEERS_INFO, "")
    if err != nil {
        w.ResponseJson(0, "sending request error: " + err.Error(), nil)
    } else {
        var data interface{}
        msg := this.node.receiveMsg(conn)
        err  = gjson.DecodeTo(&msg.Body, &data)
        if err != nil {
            w.ResponseJson(0, "received error from leader: " + err.Error(), nil)
        } else {
            w.ResponseJson(1, "ok", data)
        }
    }
    conn.Close()
}

func (this *NodeApiNode) PUT(r *ghttp.Request, w *ghttp.ServerResponse) {
    this.DELETE(r, w)
}

func (this *NodeApiNode) POST(r *ghttp.Request, w *ghttp.ServerResponse) {
    this.DELETE(r, w)
}

func (this *NodeApiNode) DELETE(r *ghttp.Request, w *ghttp.ServerResponse) {
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
    // 只允许[]interface{}数据类型
    if "[]interface {}" != reflect.TypeOf(items).String()  {
        w.ResponseJson(0, "invalid data type", nil)
        return
    }
    // 请求到leader
    conn := this.node.getConn(this.node.getLeader(), gPORT_RAFT)
    if conn == nil {
        w.ResponseJson(0, "could not connect to leader: " + this.node.getLeader(), nil)
        return
    }
    head := gMSG_API_PEERS_ADD
    if method == "DELETE" {
        head = gMSG_API_PEERS_REMOVE
    }
    err   = this.node.sendMsg(conn, head, *gjson.Encode(items))
    if err != nil {
        w.ResponseJson(0, "sending request error: " + err.Error(), nil)
    } else {
        msg := this.node.receiveMsg(conn)
        if msg.Head != gMSG_RAFT_RESPONSE {
            w.ResponseJson(0, "handling request error", nil)
        } else {
            w.ResponseJson(1, "ok", nil)
        }
    }
    conn.Close()
}

