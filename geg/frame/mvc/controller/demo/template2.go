package demo

import (
<<<<<<< HEAD
    "gitee.com/johng/gf/g/net/ghttp"
    "gitee.com/johng/gf/g"
)

func init() {
    ghttp.GetServer().BindHandler("/template2", func(r *ghttp.Request){
        content, _ := g.View().Parse("index.tpl", map[string]interface{}{
            "id"   : 123,
            "name" : "john",
        })
        r.Response.Write(content)
    })
}
=======
	"github.com/gogf/gf/g"
	"github.com/gogf/gf/g/net/ghttp"
)

func init() {
	ghttp.GetServer().BindHandler("/template2", func(r *ghttp.Request) {
		content, _ := g.View().Parse("index.tpl", map[string]interface{}{
			"id":   123,
			"name": "john",
		})
		r.Response.Write(content)
	})
}
>>>>>>> upstream/master
