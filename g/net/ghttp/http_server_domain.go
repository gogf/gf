// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.
//
package ghttp

import (
    "strings"
    "gitee.com/johng/gf/g/container/gmap"
)

// 域名管理器对象
type Domain struct {
    s *Server         // 所属Server
    m map[string]bool // 多域名
}

// 域名对象表，用以存储和检索域名(支持多域名)与域名对象之间的关联关系
var domains = gmap.NewStringInterfaceMap()

// 生成一个域名对象
func (s *Server) Domain(domain string) *Domain {
    if r := domains.Get(domain); r != nil {
        return r.(*Domain)
    }
    d := &Domain{
        s : s,
        m : make(map[string]bool),
    }
    result := strings.Split(domain, ",")
    for _, v := range result {
        d.m[strings.TrimSpace(v)] = true
    }
    domains.Set(domain, d)
    return d
}

// 注意该方法是直接绑定方法的内存地址，执行的时候直接执行该方法，不会存在初始化新的控制器逻辑
func (d *Domain) BindHandler(pattern string, handler HandlerFunc) error {
    for domain, _ := range d.m {
        if err := d.s.bindHandlerItem(pattern + "@" + domain, HandlerItem{nil, "", handler}); err != nil {
            return err
        }
    }
    return nil
}

// 绑定对象到URI请求处理中，会自动识别方法名称，并附加到对应的URI地址后面
// 需要注意对象方法的定义必须按照ghttp.HandlerFunc来定义
func (d *Domain) BindObject(pattern string, obj interface{}) error {
    for domain, _ := range d.m {
        if err := d.s.BindObject(pattern + "@" + domain, obj); err != nil {
            return err
        }
    }
    return nil
}

// 绑定对象到URI请求处理中，会自动识别方法名称，并附加到对应的URI地址后面
// 需要注意对象方法的定义必须按照ghttp.HandlerFunc来定义
func (d *Domain) BindObjectRest(pattern string, obj interface{}) error {
    for domain, _ := range d.m {
        if err := d.s.BindObjectRest(pattern + "@" + domain, obj); err != nil {
            return err
        }
    }
    return nil
}

// 绑定控制器
func (d *Domain) BindController(pattern string, c Controller) error {
    for domain, _ := range d.m {
        if err := d.s.BindController(pattern + "@" + domain, c); err != nil {
            return err
        }
    }
    return nil
}

// 绑定控制器(RESTFul)
func (d *Domain) BindControllerRest(pattern string, c Controller) error {
    for domain, _ := range d.m {
        if err := d.s.BindControllerRest(pattern + "@" + domain, c); err != nil {
            return err
        }
    }
    return nil
}

// 绑定控制器方法
func (d *Domain) BindControllerMethod(pattern string, c Controller, method string) error {
    for domain, _ := range d.m {
        if err := d.s.BindControllerMethod(pattern + "@" + domain, c, method); err != nil {
            return err
        }
    }
    return nil
}
