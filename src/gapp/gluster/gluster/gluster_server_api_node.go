// 返回格式统一：
// {result:1, message:"", data:""}

package gluster

import (
    "g/net/ghttp"
    "g/encoding/gjson"
    "reflect"
)

// 节点信息API管理
func (this *NodeApiNode) GET(r *ghttp.ClientRequest, w *ghttp.ServerResponse) {
    conn := this.node.getConn(this.node.getLeader().Ip, gPORT_RAFT)
    if conn == nil {
        w.ResponseJson(0, "could not connect to leader: " + this.node.getLeader().Ip, nil)
        return
    }
    err := this.node.sendMsg(conn, gMSG_API_PEERS_INFO, "")
    if err != nil {
        w.ResponseJson(0, "sending request error: " + err.Error(), nil)
    } else {
        msg  := this.node.receiveMsg(conn)
        data := gjson.Decode(&msg.Body)
        if data == nil {
            w.ResponseJson(0, "error data type from leader", nil)
        } else {
            w.ResponseJson(1, "ok", data)
        }
    }
    conn.Close()
}

func (this *NodeApiNode) PUT(r *ghttp.ClientRequest, w *ghttp.ServerResponse) {
    this.POST(r, w)
}

func (this *NodeApiNode) POST(r *ghttp.ClientRequest, w *ghttp.ServerResponse) {
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
    // 只允许[]interface{}数据类型
    if "[]interface {}" != reflect.TypeOf(items).String()  {
        w.ResponseJson(0, "invalid data type", nil)
        return
    }
    // 请求到leader
    conn := this.node.getConn(this.node.getLeader().Ip, gPORT_RAFT)
    if conn == nil {
        w.ResponseJson(0, "could not connect to leader: " + this.node.getLeader().Ip, nil)
        return
    }
    err := this.node.sendMsg(conn, gMSG_API_PEERS_ADD, *gjson.Encode(items))
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

func (this *NodeApiNode) DELETE(r *ghttp.ClientRequest, w *ghttp.ServerResponse) {
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
    // 只允许[]interface{}数据类型
    if "[]interface {}" != reflect.TypeOf(items).String()  {
        w.ResponseJson(0, "invalid data type", nil)
        return
    }
    // 请求到leader
    conn := this.node.getConn(this.node.getLeader().Ip, gPORT_RAFT)
    if conn == nil {
        w.ResponseJson(0, "could not connect to leader: " + this.node.getLeader().Ip, nil)
        return
    }
    err := this.node.sendMsg(conn, gMSG_API_PEERS_REMOVE, *gjson.Encode(items))
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

