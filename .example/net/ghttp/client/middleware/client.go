package main

import (
	"bytes"
	"fmt"
	"github.com/gogf/gf/container/garray"
	"github.com/gogf/gf/crypto/gmd5"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/internal/json"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/util/gconv"
	"github.com/gogf/gf/util/guid"
	"github.com/gogf/gf/util/gutil"
	"io/ioutil"
	"net/http"
)

const (
	appId     = "123"
	appSecret = "456"
)

// 注入统一的接口签名参数
func injectSignature(jsonContent []byte) []byte {
	var m map[string]interface{}
	_ = json.Unmarshal(jsonContent, &m)
	if len(m) > 0 {
		m["appid"] = appId
		m["nonce"] = guid.S()
		m["timestamp"] = gtime.Timestamp()
		var (
			keyArray   = garray.NewSortedStrArrayFrom(gutil.Keys(m))
			sigContent string
		)
		keyArray.Iterator(func(k int, v string) bool {
			sigContent += v
			sigContent += gconv.String(m[v])
			return true
		})
		m["signature"] = gmd5.MustEncryptString(gmd5.MustEncryptString(sigContent) + appSecret)
		jsonContent, _ = json.Marshal(m)
	}
	return jsonContent
}

func main() {
	c := g.Client()
	c.Use(func(c *ghttp.Client, r *http.Request) (resp *ghttp.ClientResponse, err error) {
		bodyBytes, _ := ioutil.ReadAll(r.Body)
		if len(bodyBytes) > 0 {
			// 注入签名相关参数，修改Request原有的提交参数
			bodyBytes = injectSignature(bodyBytes)
			r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
			r.ContentLength = int64(len(bodyBytes))
		}
		return c.Next(r)
	})
	content := c.ContentJson().PostContent("http://127.0.0.1:8199/", g.Map{
		"name": "goframe",
		"site": "https://goframe.org",
	})
	fmt.Println(content)
}
