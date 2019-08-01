package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/gtime"
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
		id := r.Cookie.GetSessionId()
		if id == "" {
			return
		}
		// 应用服务器一般是多个节点构成的集群，
		// 当请求中带有SESSION ID时，自动从Redis读取并恢复数据。
		value, err := g.Redis().DoVar("GET", id)
		if err != nil {
			panic(err)
		}
		if !value.IsNil() {
			if err := r.Session.Restore(value.Bytes()); err != nil {
				panic(err)
			}
		}
	}
}

// 请求结束时将SESSION数据存储到Redis中，或者在SESSION删除时也删除Redis中的数据。
func RedisHandlerSet(r *ghttp.Request) {
	if !r.IsFileRequest() {
		id := r.Cookie.GetSessionId()
		if id == "" {
			return
		}
		err := (error)(nil)
		value := ([]byte)(nil)
		if r.Session.Size() > 0 {
			if value, err = r.Session.Export(); err == nil {
				if len(value) == 0 {
					return
				} else if !r.Session.IsDirty() {
					// 更新过期时间
					_, err = g.Redis().Do("EXPIRE", id, r.Server.GetSessionMaxAge())
				} else {
					// 更新Redis数据
					_, err = g.Redis().Do("SETEX", id, r.Server.GetSessionMaxAge(), value)
				}
			}
		} else {
			// 清空SESSION后自动删除Redis数据
			_, err = g.Redis().Do("DEL", id)
		}
		if err != nil {
			panic(err)
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
