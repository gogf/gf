package stats

import (
    "gitee.com/johng/gf/g"
    "gitee.com/johng/gf/g/net/ghttp"
)

var (
    total int
)

func init() {
    g.Server().BindHandler("/stats/total", showTotal)
}

func showTotal(r *ghttp.Request) {
    total++
    r.Response.Write("total:", total)
}
