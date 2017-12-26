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
    d      := &Domain{
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

// 绑定方法
func (d *Domain) BindControllerMethod(pattern string, c Controller, method string) error {
    return d.s.BindControllerMethod(pattern, c, method)
}

// 绑定控制器
func (d *Domain) BindController(uri string, c Controller) error {
    for domain, _ := range d.m {
        if err := d.s.BindController(uri + "@" + domain, c); err != nil {
            return err
        }
    }
    return nil
}