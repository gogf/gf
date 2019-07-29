// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.
// 域名服务注册管理.

package ghttp

import (
	"strings"
)

// 域名管理器对象
type Domain struct {
	s *Server         // 所属Server
	m map[string]bool // 多域名
}

// 生成一个域名对象, 参数 domains 支持给定多个域名。
func (s *Server) Domain(domains string) *Domain {
	d := &Domain{
		s: s,
		m: make(map[string]bool),
	}
	for _, v := range strings.Split(domains, ",") {
		d.m[strings.TrimSpace(v)] = true
	}
	return d
}

// 注意该方法是直接绑定方法的内存地址，执行的时候直接执行该方法，不会存在初始化新的控制器逻辑
func (d *Domain) BindHandler(pattern string, handler HandlerFunc) {
	for domain, _ := range d.m {
		d.s.BindHandler(pattern+"@"+domain, handler)
	}
}

// 执行对象方法
func (d *Domain) BindObject(pattern string, obj interface{}, methods ...string) {
	for domain, _ := range d.m {
		d.s.BindObject(pattern+"@"+domain, obj, methods...)
	}
}

// 执行对象方法注册，methods参数不区分大小写
func (d *Domain) BindObjectMethod(pattern string, obj interface{}, method string) {
	for domain, _ := range d.m {
		d.s.BindObjectMethod(pattern+"@"+domain, obj, method)
	}
}

// RESTful执行对象注册
func (d *Domain) BindObjectRest(pattern string, obj interface{}) {
	for domain, _ := range d.m {
		d.s.BindObjectRest(pattern+"@"+domain, obj)
	}
}

// 控制器注册
func (d *Domain) BindController(pattern string, c Controller, methods ...string) {
	for domain, _ := range d.m {
		d.s.BindController(pattern+"@"+domain, c, methods...)
	}
}

// 控制器方法注册，methods参数区分大小写
func (d *Domain) BindControllerMethod(pattern string, c Controller, method string) {
	for domain, _ := range d.m {
		d.s.BindControllerMethod(pattern+"@"+domain, c, method)
	}
}

// RESTful控制器注册
func (d *Domain) BindControllerRest(pattern string, c Controller) {
	for domain, _ := range d.m {
		d.s.BindControllerRest(pattern+"@"+domain, c)
	}
}

// 绑定指定的hook回调函数, hook参数的值由ghttp server设定，参数不区分大小写
// 目前hook支持：Init/Shut
func (d *Domain) BindHookHandler(pattern string, hook string, handler HandlerFunc) {
	for domain, _ := range d.m {
		d.s.BindHookHandler(pattern+"@"+domain, hook, handler)
	}
}

// 通过map批量绑定回调函数
func (d *Domain) BindHookHandlerByMap(pattern string, hookmap map[string]HandlerFunc) {
	for domain, _ := range d.m {
		d.s.BindHookHandlerByMap(pattern+"@"+domain, hookmap)
	}
}

// 绑定指定的状态码回调函数
func (d *Domain) BindStatusHandler(status int, handler HandlerFunc) {
	for domain, _ := range d.m {
		d.s.setStatusHandler(d.s.statusHandlerKey(status, domain), handler)
	}
}

// 通过map批量绑定状态码回调函数
func (d *Domain) BindStatusHandlerByMap(handlerMap map[int]HandlerFunc) {
	for k, v := range handlerMap {
		d.BindStatusHandler(k, v)
	}
}
