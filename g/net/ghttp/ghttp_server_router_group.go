// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.
// 分组路由管理.

package ghttp

// 分组路由对象
type RouterGroup struct {
    server *Server // Web Server
    prefix string  // URI前缀
}

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

// REST路由注册
func (g *RouterGroup) REST(pattern string, object interface{}) {

}

// 绑定所有的HTTP Method请求方式
func (g *RouterGroup) ALL(pattern string, params...interface{}) {

}

func (g *RouterGroup) GET(pattern string, params...interface{}) {

}

func (g *RouterGroup) PUT(pattern string, params...interface{}) {

}

func (g *RouterGroup) POST(pattern string, params...interface{}) {

}

func (g *RouterGroup) DELETE(pattern string, params...interface{}) {

}

func (g *RouterGroup) PATCH(pattern string, params...interface{}) {

}

func (g *RouterGroup) HEAD(pattern string, params...interface{}) {

}

func (g *RouterGroup) CONNECT(pattern string, params...interface{}) {

}

func (g *RouterGroup) OPTIONS(pattern string, params...interface{}) {

}

func (g *RouterGroup) TRTACE(pattern string, params...interface{}) {

}

// 执行路由绑定
func (g *RouterGroup) bind(method string, pattern string, params...interface{}) {

}

// 判断给定对象是否控制器对象：
// 控制器必须包含以下公开的属性对象：Request/Response/Server/Cookie/Session/View.
func (g *RouterGroup) isController(object interface{}) {

}