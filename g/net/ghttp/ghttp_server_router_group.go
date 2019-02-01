// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.
// 分组路由管理.

package ghttp

import (
    "gitee.com/johng/gf/g/os/glog"
    "gitee.com/johng/gf/g/util/gconv"
    "reflect"
    "strings"
)

// 分组路由对象
type RouterGroup struct {
    server *Server // Server
    domain *Domain // Domain
    prefix string  // URI前缀
}

// 分组路由批量绑定项
type GroupItem = []interface{}

// 获取分组路由对象
func (s *Server) Group(prefix...string) *RouterGroup {
    if len(prefix) > 0 {
        return &RouterGroup{
            server : s,
            prefix : prefix[0],
        }
    }
    return &RouterGroup{}
}

// 获取分组路由对象
func (d *Domain) Group(prefix...string) *RouterGroup {
    if len(prefix) > 0 {
        return &RouterGroup{
            domain : d,
            prefix : prefix[0],
        }
    }
    return &RouterGroup{}
}

// 执行分组路由批量绑定
func (g *RouterGroup) Bind(group string, items []GroupItem) {
    for _, item := range items {
        if len(item) < 3 {
            glog.Fatalfln("invalid router item: %s", item)
        }
        if strings.EqualFold(gconv.String(item[0]), "REST") {
            g.bind("REST", gconv.String(item[0]) + ":" + gconv.String(item[1]), item[2])
        } else {
            if len(item) > 3 {
                g.bind("HANDLER", gconv.String(item[0]) + ":" + gconv.String(item[1]), item[2], item[3])
            } else {
                g.bind("HANDLER", gconv.String(item[0]) + ":" + gconv.String(item[1]), item[2])
            }
        }
    }
}

// 绑定所有的HTTP Method请求方式
func (g *RouterGroup) ALL(pattern string, object interface{}, params...interface{}) {
    g.bind("HANDLER", gDEFAULT_METHOD + ":" + pattern, object, params...)
}

// 绑定常用方法: GET/PUT/POST/DELETE
func (g *RouterGroup) COMMON(pattern string, object interface{}, params...interface{}) {
    g.GET(pattern, object, params...)
    g.PUT(pattern, object, params...)
    g.POST(pattern, object, params...)
    g.DELETE(pattern, object, params...)
}

func (g *RouterGroup) GET(pattern string, object interface{}, params...interface{}) {
    g.bind("HANDLER", "GET:" + pattern, object, params...)
}

func (g *RouterGroup) PUT(pattern string, object interface{}, params...interface{}) {
    g.bind("HANDLER", "PUT:" + pattern, object, params...)
}

func (g *RouterGroup) POST(pattern string, object interface{}, params...interface{}) {
    g.bind("HANDLER", "POST:" + pattern, object, params...)
}

func (g *RouterGroup) DELETE(pattern string, object interface{}, params...interface{}) {
    g.bind("HANDLER", "DELETE:" + pattern, object, params...)
}

func (g *RouterGroup) PATCH(pattern string, object interface{}, params...interface{}) {
    g.bind("HANDLER", "PATCH:" + pattern, object, params...)
}

func (g *RouterGroup) HEAD(pattern string, object interface{}, params...interface{}) {
    g.bind("HANDLER", "HEAD:" + pattern, object, params...)
}

func (g *RouterGroup) CONNECT(pattern string, object interface{}, params...interface{}) {
    g.bind("HANDLER", "CONNECT:" + pattern, object, params...)
}

func (g *RouterGroup) OPTIONS(pattern string, object interface{}, params...interface{}) {
    g.bind("HANDLER", "OPTIONS:" + pattern, object, params...)
}

func (g *RouterGroup) TRACE(pattern string, object interface{}, params...interface{}) {
    g.bind("HANDLER", "TRACE:" + pattern, object, params...)
}

// REST路由注册
func (g *RouterGroup) REST(pattern string, object interface{}) {
    g.bind("REST", pattern, object)
}

// 执行路由绑定
func (g *RouterGroup) bind(bindType string, pattern string, object interface{}, params...interface{}) {
    // 注册路由处理
    if len(g.prefix) > 0 {
        domain, method, path, err := g.server.parsePattern(pattern)
        if err != nil {
            glog.Fatalfln("invalid pattern: %s", pattern)
        }
        if bindType == "HANDLER" {
            pattern = g.server.serveHandlerKey(method, g.prefix + "/" + strings.TrimLeft(path, "/"), domain)
        } else {
            pattern = g.prefix + "/" + strings.TrimLeft(path, "/")
        }
    }
    methods := gconv.Strings(params)
    // 判断是否事件回调注册
    if _, ok := object.(HandlerFunc); ok && len(methods) > 0 {
        bindType = "HOOK"
    }
    switch bindType {
        case "HANDLER":
            if h, ok := object.(HandlerFunc); ok {
                if g.server != nil {
                    g.server.BindHandler(pattern, h)
                } else {
                    g.domain.BindHandler(pattern, h)
                }
            } else if g.isController(object) {
                if len(methods) > 0 {
                    if g.server != nil {
                        g.server.BindControllerMethod(pattern, object.(Controller), methods[0])
                    } else {
                        g.domain.BindControllerMethod(pattern, object.(Controller), methods[0])
                    }
                } else {
                    if g.server != nil {
                        g.server.BindController(pattern, object.(Controller))
                    } else {
                        g.domain.BindController(pattern, object.(Controller))
                    }
                }
            } else {
                if len(methods) > 0 {
                    if g.server != nil {
                        g.server.BindObjectMethod(pattern, object, methods[0])
                    } else {
                        g.domain.BindObjectMethod(pattern, object, methods[0])
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
                    g.server.BindHookHandler(pattern, methods[0], h)
                } else {
                    g.domain.BindHookHandler(pattern, methods[0], h)
                }
            } else {
                glog.Fatalfln("invalid hook handler for pattern:%s", pattern)
            }
    }
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
