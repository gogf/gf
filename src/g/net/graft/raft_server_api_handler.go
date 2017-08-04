// 返回格式统一：
// {result:1, message:"", data:""}

package graft

import (
    "io"
    "strings"
    "g/net/ghttp"
)


// K-V API管理
func (n *Node) kvApiHandler(r *ghttp.Request, w *ghttp.Response) {
    method := strings.ToUpper(r.Method)
    switch method {
        case "GET":
            k := r.GetRequest("k")
            if k == nil {
                w.ResponseJson(1, "", *n.KVMap.Clone())
            } else {
                w.ResponseJson(1, "", n.KVMap.Get(k.(string)))
            }

        case "PUT":
        case "POST":
        case "DELETE":

    }

}

// 节点信息API管理
func (n *Node) nodeApiHandler(r *ghttp.Request, w *ghttp.Response) {
    io.WriteString(w, "hello\n")
}

