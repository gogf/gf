package main

import (
	"github.com/jin502437344/gf/frame/g"
	"github.com/jin502437344/gf/net/ghttp"
)

type RegisterReq struct {
	Name  string
	Pass  string `p:"password1"`
	Pass2 string `p:"password2"`
}

type RegisterRes struct {
	Code  int         `json:"code"`
	Error string      `json:"error"`
	Data  interface{} `json:"data"`
}

func main() {
	s := g.Server()
	s.BindHandler("/register", func(r *ghttp.Request) {
		var req *RegisterReq
		if err := r.Parse(&req); err != nil {
			r.Response.WriteJsonExit(RegisterRes{
				Code:  1,
				Error: err.Error(),
			})
		}
		// ...
		r.Response.WriteJsonExit(RegisterRes{
			Data: req,
		})
	})
	s.SetPort(8199)
	s.Run()
}
