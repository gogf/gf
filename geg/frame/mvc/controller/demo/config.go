package demo

import "gitee.com/johng/gf/g/net/ghttp"

func init() {
    ghttp.GetServer().BindHandler("/config", func (r *ghttp.Request) {
        r.Response.WriteString("Apple")
    })
}
