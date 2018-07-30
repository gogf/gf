package stats

import (
    "gitee.com/johng/gf/g"
    "gitee.com/johng/gf/g/net/ghttp"
    "gitee.com/johng/gf/g/util/gconv"
    "gitee.com/johng/gf/g/container/gtype"
)

var (
    total2 = gtype.NewInt()
)

func init() {
    g.Server().BindHandler("/stats/total2", showTotal2)
}

func showTotal2(r *ghttp.Request) {
    r.Response.Write("total:", gconv.String(total2.Add(1)))
}
