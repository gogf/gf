package stats

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
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
