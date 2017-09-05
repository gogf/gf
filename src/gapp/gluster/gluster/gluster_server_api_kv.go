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
        if this.node.DataMap.Size() > 100 {
            w.ResponseJson(0, "too large data size, need a key to search", nil)
        } else {
            w.ResponseJson(1, "ok", *this.node.DataMap.Clone())
        }
    } else {
        if this.node.DataMap.Contains(k) {
            w.ResponseJson(1, "ok", this.node.DataMap.Get(k))
        } else {
            w.ResponseJson(0, "ok", nil)
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
    err  = this.node.SendToLeader(gMSG_REPL_DATA_SET, gPORT_REPL, gjson.Encode(items))
    if err != nil {
        w.ResponseJson(0, err.Error(), nil)
    } else {
        // 为保证客户端能够及时相应（例如在写入请求的下一次获取请求将一定能够获取到最新的数据），
        // 因此，请求端应当在leader返回成功后，同时将该数据写入到本地
        if this.node.Id != this.node.getLeader().Id {
            this.node.DataMap.BatchSet(items)
        }
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
    err  = this.node.SendToLeader(gMSG_REPL_DATA_REMOVE, gPORT_REPL, gjson.Encode(list))
    if err != nil {
        w.ResponseJson(0, err.Error(), nil)
    } else {
        if this.node.Id != this.node.getLeader().Id {
            this.node.DataMap.BatchRemove(list)
        }
        w.ResponseJson(1, "ok", nil)
    }
}
