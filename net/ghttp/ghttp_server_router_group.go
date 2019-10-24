// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import (
	"reflect"
	"strings"

	"github.com/gogf/gf/text/gstr"

	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/util/gconv"
)

// 分组路由对象
type RouterGroup struct {
	parent *RouterGroup // 父级分组路由
	server *Server      // Server
	domain *Domain      // Domain
	prefix string       // URI前缀
}

// 分组路由批量绑定项
type GroupItem = []interface{}

// 预绑定路由项结构
type groupPreBindItem struct {
	group    *RouterGroup
	bindType string
	pattern  string
	object   interface{}
	params   []interface{}
}

var (
	// 预处理路由项存储数组
	preBindItems = make([]groupPreBindItem, 0, 64)
)

// 处理预绑定路由项
func (s *Server) handlePreBindItems() {
	for _, item := range preBindItems {
		if item.group.server != nil && item.group.server != s {
			continue
		}
		if item.group.domain != nil && item.group.domain.s != s {
			continue
		}
		item.group.doBind(item.bindType, item.pattern, item.object, item.params...)
	}
}

// 获取分组路由对象
func (s *Server) Group(prefix string, groups ...func(g *RouterGroup)) *RouterGroup {
	// 自动识别并加上/前缀
	if prefix[0] != '/' {
		prefix = "/" + prefix
	}
	if prefix == "/" {
		prefix = ""
	}
	group := &RouterGroup{
		server: s,
		prefix: prefix,
	}
	if len(groups) > 0 {
		for _, v := range groups {
			v(group)
		}
	}
	return group
}

// 获取分组路由对象(绑定域名)
func (d *Domain) Group(prefix string, groups ...func(g *RouterGroup)) *RouterGroup {
	if prefix == "/" {
		prefix = ""
	}
	group := &RouterGroup{
		domain: d,
		prefix: prefix,
	}
	if len(groups) > 0 {
		for _, v := range groups {
			v(group)
		}
	}
	return group
}

// 层级递归创建分组路由注册项
func (g *RouterGroup) Group(prefix string, groups ...func(g *RouterGroup)) *RouterGroup {
	if prefix == "/" {
		prefix = ""
	}
	group := &RouterGroup{
		parent: g,
		server: g.server,
		domain: g.domain,
		prefix: prefix,
	}
	if len(groups) > 0 {
		for _, v := range groups {
			v(group)
		}
	}
	return group
}

func (g *RouterGroup) Clone() *RouterGroup {
	return &RouterGroup{
		parent: g.parent,
		server: g.server,
		domain: g.domain,
		prefix: g.prefix,
	}
}

// 执行分组路由批量绑定
func (g *RouterGroup) Bind(items []GroupItem) *RouterGroup {
	group := g.Clone()
	for _, item := range items {
		if len(item) < 3 {
			glog.Fatalf("invalid router item: %s", item)
		}
		bindType := gstr.ToUpper(gconv.String(item[0]))
		switch bindType {
		case "REST":
			group.preBind("REST", gconv.String(item[0])+":"+gconv.String(item[1]), item[2])
		case "MIDDLEWARE":
			group.preBind("MIDDLEWARE", gconv.String(item[0])+":"+gconv.String(item[1]), item[2])
		default:
			if strings.EqualFold(bindType, "ALL") {
				bindType = ""
			} else {
				bindType += ":"
			}
			if len(item) > 3 {
				group.preBind("HANDLER", bindType+gconv.String(item[1]), item[2], item[3])
			} else {
				group.preBind("HANDLER", bindType+gconv.String(item[1]), item[2])
			}
		}
	}
	return group
}

// 绑定所有的HTTP Method请求方式
func (g *RouterGroup) ALL(pattern string, object interface{}, params ...interface{}) *RouterGroup {
	return g.Clone().preBind("HANDLER", gDEFAULT_METHOD+":"+pattern, object, params...)
}

// 绑定常用方法: GET/PUT/POST/DELETE
func (g *RouterGroup) COMMON(pattern string, object interface{}, params ...interface{}) *RouterGroup {
	group := g.Clone()
	group.preBind("HANDLER", "GET:"+pattern, object, params...)
	group.preBind("HANDLER", "PUT:"+pattern, object, params...)
	group.preBind("HANDLER", "POST:"+pattern, object, params...)
	group.preBind("HANDLER", "DELETE:"+pattern, object, params...)
	return group
}

func (g *RouterGroup) GET(pattern string, object interface{}, params ...interface{}) *RouterGroup {
	return g.Clone().preBind("HANDLER", "GET:"+pattern, object, params...)
}

func (g *RouterGroup) PUT(pattern string, object interface{}, params ...interface{}) *RouterGroup {
	return g.Clone().preBind("HANDLER", "PUT:"+pattern, object, params...)
}

func (g *RouterGroup) POST(pattern string, object interface{}, params ...interface{}) *RouterGroup {
	return g.Clone().preBind("HANDLER", "POST:"+pattern, object, params...)
}

func (g *RouterGroup) DELETE(pattern string, object interface{}, params ...interface{}) *RouterGroup {
	return g.Clone().preBind("HANDLER", "DELETE:"+pattern, object, params...)
}

func (g *RouterGroup) PATCH(pattern string, object interface{}, params ...interface{}) *RouterGroup {
	return g.Clone().preBind("HANDLER", "PATCH:"+pattern, object, params...)
}

func (g *RouterGroup) HEAD(pattern string, object interface{}, params ...interface{}) *RouterGroup {
	return g.Clone().preBind("HANDLER", "HEAD:"+pattern, object, params...)
}

func (g *RouterGroup) CONNECT(pattern string, object interface{}, params ...interface{}) *RouterGroup {
	return g.Clone().preBind("HANDLER", "CONNECT:"+pattern, object, params...)
}

func (g *RouterGroup) OPTIONS(pattern string, object interface{}, params ...interface{}) *RouterGroup {
	return g.Clone().preBind("HANDLER", "OPTIONS:"+pattern, object, params...)
}

func (g *RouterGroup) TRACE(pattern string, object interface{}, params ...interface{}) *RouterGroup {
	return g.Clone().preBind("HANDLER", "TRACE:"+pattern, object, params...)
}

func (g *RouterGroup) REST(pattern string, object interface{}) *RouterGroup {
	return g.Clone().preBind("REST", pattern, object)
}

func (g *RouterGroup) Hook(pattern string, hook string, handler HandlerFunc) *RouterGroup {
	return g.Clone().preBind("HANDLER", pattern, handler, hook)
}

func (g *RouterGroup) Middleware(handlers ...HandlerFunc) *RouterGroup {
	group := g.Clone()
	for _, handler := range handlers {
		group.preBind("MIDDLEWARE", "/*", handler)
	}
	return group
}

func (g *RouterGroup) MiddlewarePattern(pattern string, handlers ...HandlerFunc) *RouterGroup {
	group := g.Clone()
	for _, handler := range handlers {
		group.preBind("MIDDLEWARE", pattern, handler)
	}
	return group
}

func (g *RouterGroup) preBind(bindType string, pattern string, object interface{}, params ...interface{}) *RouterGroup {
	preBindItems = append(preBindItems, groupPreBindItem{
		group:    g,
		bindType: bindType,
		pattern:  pattern,
		object:   object,
		params:   params,
	})
	return g
}

func (g *RouterGroup) getPrefix() string {
	prefix := g.prefix
	parent := g.parent
	for parent != nil {
		prefix = parent.prefix + prefix
		parent = parent.parent
	}
	return prefix
}

// 执行路由绑定
func (g *RouterGroup) doBind(bindType string, pattern string, object interface{}, params ...interface{}) *RouterGroup {
	prefix := g.getPrefix()
	// 注册路由处理
	if len(prefix) > 0 {
		domain, method, path, err := g.server.parsePattern(pattern)
		if err != nil {
			glog.Fatalf("invalid pattern: %s", pattern)
		}
		// If there'a already a domain, unset the domain field in the pattern.
		if g.domain != nil {
			domain = ""
		}
		if bindType == "REST" {
			pattern = prefix + "/" + strings.TrimLeft(path, "/")
		} else {
			pattern = g.server.serveHandlerKey(method, prefix+"/"+strings.TrimLeft(path, "/"), domain)
		}
	}
	// 去掉可能重复出现的'//'符号
	pattern = gstr.Replace(pattern, "//", "/")
	// 将附加参数转换为字符串
	extras := gconv.Strings(params)
	// 判断是否事件回调注册
	if _, ok := object.(HandlerFunc); ok && len(extras) > 0 {
		bindType = "HOOK"
	}
	switch bindType {
	case "MIDDLEWARE":
		if h, ok := object.(HandlerFunc); ok {
			if g.server != nil {
				g.server.BindMiddleware(pattern, h)
			} else {
				g.domain.BindMiddleware(pattern, h)
			}
		} else {
			glog.Fatalf("invalid middleware handler for pattern:%s", pattern)
		}
	case "HANDLER":
		if h, ok := object.(HandlerFunc); ok {
			if g.server != nil {
				g.server.BindHandler(pattern, h)
			} else {
				g.domain.BindHandler(pattern, h)
			}
		} else if g.isController(object) {
			if len(extras) > 0 {
				if g.server != nil {
					g.server.BindControllerMethod(pattern, object.(Controller), extras[0])
				} else {
					g.domain.BindControllerMethod(pattern, object.(Controller), extras[0])
				}
			} else {
				if g.server != nil {
					g.server.BindController(pattern, object.(Controller))
				} else {
					g.domain.BindController(pattern, object.(Controller))
				}
			}
		} else {
			if len(extras) > 0 {
				if g.server != nil {
					g.server.BindObjectMethod(pattern, object, extras[0])
				} else {
					g.domain.BindObjectMethod(pattern, object, extras[0])
				}
			} else {
				if g.server != nil {
					g.server.BindObject(pattern, object)
				} else {
					g.domain.BindObject(pattern, object)
				}
			}
		}
	case "REST":
		if g.isController(object) {
			if g.server != nil {
				g.server.BindControllerRest(pattern, object.(Controller))
			} else {
				g.domain.BindControllerRest(pattern, object.(Controller))
			}
		} else {
			if g.server != nil {
				g.server.BindObjectRest(pattern, object)
			} else {
				g.domain.BindObjectRest(pattern, object)
			}
		}
	case "HOOK":
		if h, ok := object.(HandlerFunc); ok {
			if g.server != nil {
				g.server.BindHookHandler(pattern, extras[0], h)
			} else {
				g.domain.BindHookHandler(pattern, extras[0], h)
			}
		} else {
			glog.Fatalf("invalid hook handler for pattern:%s", pattern)
		}
	}
	return g
}

// 判断给定对象是否控制器对象：
// 控制器必须包含以下公开的属性对象：Request/Response/Server/Cookie/Session/View.
func (g *RouterGroup) isController(value interface{}) bool {
	// 首先判断是否满足控制器接口定义
	if _, ok := value.(Controller); !ok {
		return false
	}
	// 其次检查控制器的必需属性
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.FieldByName("Request").IsValid() && v.FieldByName("Response").IsValid() &&
		v.FieldByName("Server").IsValid() && v.FieldByName("Cookie").IsValid() &&
		v.FieldByName("Session").IsValid() && v.FieldByName("View").IsValid() {
		return true
	}
	return false
}
