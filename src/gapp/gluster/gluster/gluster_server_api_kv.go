// 返回格式统一：
// {result:1, message:"", data:""}

package gluster

import (
    "g/net/ghttp"
    "g/encoding/gjson"
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
func (this *NodeApiKv) POST(r *ghttp.ClientRequest, w *ghttp.ServerResponse) {
    items := make(map[string]string)
    err   := gjson.DecodeTo(r.GetRaw(), &items)
    if err != nil {
        w.ResponseJson(0, "invalid data type: " + err.Error(), nil)
        return
    }
    err  = this.node.SendToLeader(gMSG_REPL_SET, gPORT_REPL, gjson.Encode(items))
    if err != nil {
        w.ResponseJson(0, err.Error(), nil)
    } else {
        w.ResponseJson(1, "ok", nil)
    }
}

// K-V 删除
func (this *NodeApiKv) DELETE(r *ghttp.ClientRequest, w *ghttp.ServerResponse) {
    list := make([]string, 0)
    err  := gjson.DecodeTo(r.GetRaw(), &list)
    if err != nil {
        w.ResponseJson(0, "invalid data type: " + err.Error(), nil)
        return
    }
    err  = this.node.SendToLeader(gMSG_REPL_REMOVE, gPORT_REPL, gjson.Encode(list))
    if err != nil {
        w.ResponseJson(0, err.Error(), nil)
    } else {
        w.ResponseJson(1, "ok", nil)
    }
}
