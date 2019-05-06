package main

import (
	"crypto/tls"
	"fmt"
	"github.com/gogf/gf/g/net/ghttp"
	"net/http"
)

func main() {
	c := ghttp.NewClient()
	c.Transport = &http.Transport{
		TLSClientConfig : &tls.Config{ InsecureSkipVerify: true},
	}
	r, e := c.Clone().Get("https://127.0.0.1:8199")
	fmt.Println(e)
	fmt.Println(r.StatusCode)
}
