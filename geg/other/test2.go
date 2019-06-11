package main

import (
<<<<<<< HEAD
    "fmt"
    "gitee.com/johng/gf/g/os/gpm"
    "os"
    "time"
    "gitee.com/johng/gf/g/os/glog"
)

func main() {
    m   := gproc.New()
    env := os.Environ()
    env  = append(env, "child=1")
    p   := m.NewProcess(os.Args[0], os.Args, env)
    if os.Getenv("child") != "" {
        time.Sleep(3*time.Second)
        glog.Error("error")
    } else {
        pid, err := p.Run()
        fmt.Println(pid)
        fmt.Println(err)
        fmt.Println(p.Wait())
    }
}
=======
	"github.com/gogf/gf/g"
	"github.com/gogf/gf/g/net/ghttp"
)

func main() {
	s := g.Server()
	s.BindHookHandler("/*any", ghttp.HOOK_BEFORE_SERVE, func(r *ghttp.Request) {
		r.Response.SetAllowCrossDomainRequest("*", "PUT,GET,POST,DELETE,OPTIONS")
		r.Response.Header().Set("Access-Control-Allow-Credentials", "true")
		r.Response.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, Content-Type, token")
	})
	s.Group("/v1").COMMON("*", func(r *ghttp.Request) {
		r.Response.WriteJson(g.Map{"name": "john"})
	})
	s.SetPort(6789)
	s.Run()
}
>>>>>>> upstream/master
