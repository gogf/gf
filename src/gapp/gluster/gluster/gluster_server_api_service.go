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
    "reflect"
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

    if "map[string]interface {}" != reflect.TypeOf(items).String()  {
        w.ResponseJson(0, "invalid data type", nil)
        return
    }
    // 请求到leader
    conn := this.node.getConn(this.node.getLeader(), gPORT_REPL)
    if conn == nil {
        w.ResponseJson(0, "could not connect to leader: " + this.node.getLeader(), nil)
        return
    }
    err := this.node.sendMsg(conn, gMSG_API_SERVICE_SET, *gjson.Encode(items))
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

// service 删除
func (this *NodeApiService) DELETE(r *ghttp.ClientRequest, w *ghttp.ServerResponse) {
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
        w.ResponseJson(0, "invalid data type", nil)
        return
    }
    // 请求到leader
    conn := this.node.getConn(this.node.getLeader(), gPORT_REPL)
    if conn == nil {
        w.ResponseJson(0, "could not connect to leader: " + this.node.getLeader(), nil)
        return
    }
    err := this.node.sendMsg(conn, gMSG_API_SERVICE_REMOVE, *gjson.Encode(items))
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
