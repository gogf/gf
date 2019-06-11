package demo

<<<<<<< HEAD
import "gitee.com/johng/gf/g/net/ghttp"

func init() {
    ghttp.GetServer().BindHandler("/apple",     Apple)
    ghttp.GetServer().BindHandler("/pen",       Pen)
    ghttp.GetServer().BindHandler("/apple-pen", ApplePen)
}

func Apple(r *ghttp.Request) {
    r.Response.Write("Apple")
}

func Pen(r *ghttp.Request) {
    r.Response.Write("Pen")
}

func ApplePen(r *ghttp.Request) {
    r.Response.Write("Apple-Pen")
}
=======
import (
	"github.com/gogf/gf/g"
	"github.com/gogf/gf/g/net/ghttp"
)

func init() {
	s := g.Server()
	s.BindHandler("/apple", Apple)
	s.BindHandler("/pen", Pen)
	s.BindHandler("/apple-pen", ApplePen)
}

func Apple(r *ghttp.Request) {
	r.Response.Write("Apple")
}

func Pen(r *ghttp.Request) {
	r.Response.Write("Pen")
}

func ApplePen(r *ghttp.Request) {
	r.Response.Write("Apple-Pen")
}
>>>>>>> upstream/master
