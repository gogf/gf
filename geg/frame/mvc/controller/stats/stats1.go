package stats

import (
    "gitee.com/johng/gf/g"
    "gitee.com/johng/gf/g/net/ghttp"
    "gitee.com/johng/gf/g/util/gconv"
)

var (
    total1 int
)

func init() {
    g.Server().BindHandler("/stats/total1", showTotal1)
}

func showTotal1(r *ghttp.Request) {
    total1++
    r.Response.Write("total:", gconv.String(total1))
}
