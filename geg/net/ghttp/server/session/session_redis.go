package main

import (
	"encoding/json"

	"github.com/gogf/gf/g"
	"github.com/gogf/gf/g/net/ghttp"
	"github.com/gogf/gf/g/os/gtime"
)

// 测试，SESSION写入
func SessionSet(r *ghttp.Request) {
	r.Session.Set("time", gtime.Second())
	r.Response.WriteJson("ok")
}

// 测试，SESSION读取
func SessionGet(r *ghttp.Request) {
	r.Response.WriteJson(r.Session.Map())
}

// 请求处理之前将Redis中的数据读取出来并存储到SESSION对象中。
func RedisHandlerGet(r *ghttp.Request) {
	if !r.IsFileRequest() {
		sessionId := r.Cookie.GetSessionId()
		if sessionId == "" {
			return
		}
		value, err := g.Redis().DoVar("GET", sessionId)
		if err != nil {
			panic(err)
		}
		if !value.IsNil() {
			err := r.Session.RestoreFromJson(value.Bytes())
			if err != nil {
				panic(err)
			}
		}
	}
}

// 请求结束时将SESSION数据存储到Redis中，或者在SESSION删除时也删除Redis中的数据。
func RedisHandlerSet(r *ghttp.Request) {
	if !r.IsFileRequest() {
		sessionId := r.Cookie.GetSessionId()
		if sessionId == "" {
			return
		}
		err := (error)(nil)
		if r.Session.Size() > 0 {
			value, err := json.Marshal(r.Session)
			if err != nil {
				panic(err)
			}
			_, err = g.Redis().Do("SETEX", r.Cookie.GetSessionId(), r.Server.GetSessionMaxAge(), value)
			if err != nil {
				panic(err)
			}
		} else {
			_, err = g.Redis().Do("DEL", r.Cookie.GetSessionId())
			if err != nil {
				panic(err)
			}
		}
	}
}

func main() {
	s := g.Server()
	s.BindHandler("/set", SessionSet)
	s.BindHandler("/get", SessionGet)
	s.BindHookHandlerByMap("/*", map[string]ghttp.HandlerFunc{
		ghttp.HOOK_BEFORE_SERVE: RedisHandlerGet,
		ghttp.HOOK_AFTER_SERVE:  RedisHandlerSet,
	})
	s.SetPort(8199)
	s.Run()
}
