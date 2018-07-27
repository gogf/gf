// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.
// 事件回调注册.

package ghttp

import (
    "errors"
    "strings"
    "container/list"
    "gitee.com/johng/gf/g/util/gregex"
)

// hook缓存项，根据URL.Path进行缓存，因此对象中带有缓存参数
type hookCacheItem struct {
    faddr    HandlerFunc         // 准确的执行方法内存地址
    values   map[string][]string // GET解析参数
}

// 事件回调注册方法
// 因为有事件回调优先级的关系，叶子节点必须为一个链表，因此这里只有动态注册
func (s *Server) setHookHandler(pattern string, hook string, handler *HandlerItem) error {
    return s.setHandler(pattern, handler, hook)
}

// 事件回调 - 检索动态路由规则
// 并按照指定hook回调函数的优先级及注册顺序进行调用
func (s *Server) callHookHandler(r *Request, hook string) {
    // 如果没有注册事件回调，那么不做后续处理
    if len(s.hooksTree) == 0 {
        return
    }

    s.hhcmu.RLock()
    defer s.hhcmu.RUnlock()

    var hookItems []*hookCacheItem
    cacheKey := s.handlerHookKey(r.GetHost(), r.Method, r.URL.Path, hook)
    if v := s.hooksCache.Get(cacheKey); v == nil {
        hookItems = s.searchHookHandler(r, hook)
        if hookItems != nil {
            s.hooksCache.Set(cacheKey, hookItems, 0)
        }
    } else {
        hookItems = v.([]*hookCacheItem)
    }
    if hookItems != nil {
        for _, item := range hookItems {
            for k, v := range item.values {
                r.queries[k] = v
            }
            item.faddr(r)
        }
    }
}

func (s *Server) searchHookHandler(r *Request, hook string) []*hookCacheItem {

}

// 绑定指定的hook回调函数, pattern参数同BindHandler，支持命名路由；hook参数的值由ghttp server设定，参数不区分大小写
func (s *Server)BindHookHandler(pattern string, hook string, handler HandlerFunc) error {
    return s.setHookHandler(pattern, hook, &HandlerItem{
        ctype : nil,
        fname : "",
        faddr : handler,
    })
    return nil
}

// 通过map批量绑定回调函数
func (s *Server)BindHookHandlerByMap(pattern string, hookmap map[string]HandlerFunc) error {
    for k, v := range hookmap {
        if err := s.BindHookHandler(pattern, k, v); err != nil {
            return err
        }
    }
    return nil
}

// 构造用于hooksMap检索的键名
func (s *Server)handlerHookKey(domain, method, uri, hook string) string {
    return strings.ToUpper(hook) + "^" + s.handlerKey(domain, method, uri)
}
