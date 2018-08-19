// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.
// 域名服务注册管理.

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
var domainMap = gmap.NewStringInterfaceMap()

// 生成一个域名对象
func (s *Server) Domain(domains string) *Domain {
    if r := domainMap.Get(domains); r != nil {
        return r.(*Domain)
    }
    d := &Domain{
        s : s,
        m : make(map[string]bool),
    }
    result := strings.Split(domains, ",")
    for _, v := range result {
        d.m[strings.TrimSpace(v)] = true
    }
    domainMap.Set(domains, d)
    return d
}

// 注意该方法是直接绑定方法的内存地址，执行的时候直接执行该方法，不会存在初始化新的控制器逻辑
func (d *Domain) BindHandler(pattern string, handler HandlerFunc) error {
    for domain, _ := range d.m {
        if err := d.s.BindHandler(pattern + "@" + domain, handler); err != nil {
            return err
        }
    }
    return nil
}

// 执行对象方法
func (d *Domain) BindObject(pattern string, obj interface{}, methods...string) error {
    if len(methods) > 0 {
        return d.BindObjectMethod(pattern, obj, strings.Join(methods, ","))
    }
    for domain, _ := range d.m {
        if err := d.s.BindObject(pattern + "@" + domain, obj); err != nil {
            return err
        }
    }
    return nil
}

// 执行对象方法注册，methods参数不区分大小写
func (d *Domain) BindObjectMethod(pattern string, obj interface{}, method string) error {
    for domain, _ := range d.m {
        if err := d.s.BindObjectMethod(pattern + "@" + domain, obj, method); err != nil {
            return err
        }
    }
    return nil
}

// RESTful执行对象注册
func (d *Domain) BindObjectRest(pattern string, obj interface{}) error {
    for domain, _ := range d.m {
        if err := d.s.BindObjectRest(pattern + "@" + domain, obj); err != nil {
            return err
        }
    }
    return nil
}

// 控制器注册
func (d *Domain) BindController(pattern string, c Controller, methods...string) error {
    if len(methods) > 0 {
        return d.BindControllerMethod(pattern, c, strings.Join(methods, ","))
    }
    for domain, _ := range d.m {
        if err := d.s.BindController(pattern + "@" + domain, c); err != nil {
            return err
        }
    }
    return nil
}

// 控制器方法注册，methods参数区分大小写
func (d *Domain) BindControllerMethod(pattern string, c Controller, method string) error {
    for domain, _ := range d.m {
        if err := d.s.BindControllerMethod(pattern + "@" + domain, c, method); err != nil {
            return err
        }
    }
    return nil
}

// RESTful控制器注册
func (d *Domain) BindControllerRest(pattern string, c Controller) error {
    for domain, _ := range d.m {
        if err := d.s.BindControllerRest(pattern + "@" + domain, c); err != nil {
            return err
        }
    }
    return nil
}

// 绑定指定的hook回调函数, hook参数的值由ghttp server设定，参数不区分大小写
// 目前hook支持：Init/Shut
func (d *Domain)BindHookHandler(pattern string, hook string, handler HandlerFunc) error {
    for domain, _ := range d.m {
        if err := d.s.BindHookHandler(pattern + "@" + domain, hook, handler); err != nil {
            return err
        }
    }
    return nil
}

// 通过map批量绑定回调函数
func (d *Domain)BindHookHandlerByMap(pattern string, hookmap map[string]HandlerFunc) error {
    for domain, _ := range d.m {
        if err := d.s.BindHookHandlerByMap(pattern + "@" + domain, hookmap); err != nil {
            return err
        }
    }
    return nil
}

// 绑定指定的状态码回调函数
func (d *Domain)BindStatusHandler(status int, handler HandlerFunc) {
    for domain, _ := range d.m {
        d.s.setStatusHandler(d.s.statusHandlerKey(status, domain), handler)
    }
}

// 通过map批量绑定状态码回调函数
func (d *Domain)BindStatusHandlerByMap(handlerMap map[int]HandlerFunc) {
    for k, v := range handlerMap {
        d.BindStatusHandler(k, v)
    }
}