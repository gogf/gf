package demo

<<<<<<< HEAD
import "gitee.com/johng/gf/g/net/ghttp"

func init() {
    ghttp.GetServer().BindHandler("/", func(r *ghttp.Request){
        r.Response.Write("Hello World!")
    })
}
=======
import "github.com/gogf/gf/g/net/ghttp"

func init() {
	ghttp.GetServer().BindHandler("/", func(r *ghttp.Request) {
		r.Response.Write("Hello World!")
	})
}
>>>>>>> upstream/master
