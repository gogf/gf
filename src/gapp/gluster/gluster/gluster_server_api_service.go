// 目前支持两种类型服务的健康检查：mysql, web
// 支持对每个服务节点使用自定义的check配置健康检查，check的格式参考示例。
// check配置是非必需的，默认情况下gluster会采用默认的配置进行健康检查。
// 注意：
// 服务在gluster中也是使用kv存储，使用配置中的name作为键名，因此，每个服务的名称不能重复。
//
// 1、MySQL数据库服务：
// {
//     "name" : "Site Database",
//     "type" : "mysql",
//     "list" : [
//         {"host":"192.168.2.102", "port":"3306", "user":"root", "pass":"123456", "database":"test"},
//         {"host":"192.168.2.124", "port":"3306", "user":"root", "pass":"123456", "database":"test"}
//     ]
// }
// 2、WEB服务：
// {
//     "name" : "Home Site",
//     "type" : "web",
//     "list" : [
//         {"url":"http://192.168.2.102", "check":"http://192.168.2.102/health"},
//         {"url":"http://192.168.2.124", }
//     ]
// }
// 返回格式统一：
// {result:1, message:"", data:""}

package gluster

import (
    "g/net/ghttp"
    "g/encoding/gjson"
)


// service 查询
func (this *NodeApiService) GET(r *ghttp.ClientRequest, w *ghttp.ServerResponse) {
    name := r.GetRequestString("name")
    if name == "" {
        w.ResponseJson(1, "ok", *this.node.ServiceForApi.Clone())
    } else {
        if this.node.Service.Contains(name) {
            w.ResponseJson(1, "ok", this.node.ServiceForApi.Get(name))
        } else {
            w.ResponseJson(1, "ok", nil)
        }
    }
}

// service 新增
func (this *NodeApiService) PUT(r *ghttp.ClientRequest, w *ghttp.ServerResponse) {
    this.POST(r, w)
}

// service 修改
func (this *NodeApiService) POST(r *ghttp.ClientRequest, w *ghttp.ServerResponse) {
    list := make([]ServiceStruct, 0)
    err  := gjson.DecodeTo(r.GetRaw(), &list)
    if err != nil {
        w.ResponseJson(0, "invalid data type: " + err.Error(), nil)
        return
    }
    for _, v := range list {
        err  = this.node.SendToLeader(gMSG_API_SERVICE_SET, gPORT_REPL, gjson.Encode(v))
        if err != nil {
            w.ResponseJson(0, err.Error(), nil)
            return
        }
    }
    w.ResponseJson(1, "ok", nil)
}

// service 删除
func (this *NodeApiService) DELETE(r *ghttp.ClientRequest, w *ghttp.ServerResponse) {
    list := make([]string, 0)
    err  := gjson.DecodeTo(r.GetRaw(), &list)
    if err != nil {
        w.ResponseJson(0, "invalid data type: " + err.Error(), nil)
        return
    }
    err  = this.node.SendToLeader(gMSG_API_SERVICE_REMOVE, gPORT_REPL, gjson.Encode(list))
    if err != nil {
        w.ResponseJson(0, err.Error(), nil)
    } else {
        w.ResponseJson(1, "ok", nil)
    }
}
