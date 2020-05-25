package main

import (
	"fmt"
	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/util/guid"
)

func CreateSessionId(r *ghttp.Request) string {
	var (
		agent   = r.UserAgent()
		address = r.RemoteAddr
		cookie  = r.Header.Get("Cookie")
	)
	return guid.S([]byte(agent), []byte(address), []byte(cookie))
}
func main() {
	body := "{\"id\": 413231383385427875}"
	if dat, err := gjson.DecodeToJson(body); err == nil {
		fmt.Println(dat.MustToJsonString())
	}
}
