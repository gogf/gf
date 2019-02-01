package main

<<<<<<< HEAD
func main() {
    //listener, err := reuseport.Listen("tcp", ":8881")
    //if err != nil {
    //    panic(err)
    //}
    //defer listener.Close()
    //
    //server := &http.Server{}
    //http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    //    fmt.Println(gproc.Pid())
    //    fmt.Fprintf(w, "%d\n", gproc.Pid())
    //})
    //
    //panic(server.Serve(listener))
=======
import (
    "gitee.com/johng/gf/g"
    "gitee.com/johng/gf/g/net/ghttp"
)

func main() {
    s := g.Server()
    s.BindHookHandler("/*any", ghttp.HOOK_BEFORE_SERVE, func(r *ghttp.Request) {
        r.Response.SetAllowCrossDomainRequest("*", "PUT,GET,POST,DELETE,OPTIONS")
        r.Response.Header().Set("Access-Control-Allow-Credentials", "true")
        r.Response.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, Content-Type, token")
    })
    s.Group("/v1").COMMON("*", func(r *ghttp.Request) {
        r.Response.WriteJson(g.Map{"name" : "john"})
    })
    s.SetPort(6789)
    s.Run()
>>>>>>> master
}