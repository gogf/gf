// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.
// 事件回调注册.

package ghttp

// 事件回调注册方法
// 因为有事件回调优先级的关系，叶子节点必须为一个链表，因此这里只有动态注册
func (s *Server) setHookHandler(pattern string, hook string, handler *handlerItem) error {
    return s.setHandler(pattern, handler, hook)
}

// 检索事件回调方法
func (s *Server) searchHookHandler(r *Request, hook string) []*handlerItem {
    if item := s.getHandlerWithCache(r); item != nil {
        if l, ok := item.hooks[hook]; ok {
            items := make([]*handlerItem, 0)
            for e := l.Front(); e != nil; e = e.Next() {
                items = append(items, e.Value.(*handlerItem))
            }
            return items
        }
    }
    return nil
}

// 事件回调处理，内部使用了缓存处理.
// 并按照指定hook回调函数的优先级及注册顺序进行调用
func (s *Server) callHookHandler(r *Request, hook string) {
    var hookItems []*handlerItem
    cacheKey := s.hookHandlerKey(hook, r.Method, r.URL.Path, r.GetHost())
    if v := s.hooksCache.Get(cacheKey); v == nil {
        hookItems = s.searchHookHandler(r, hook)
        if hookItems != nil {
            s.hooksCache.Set(cacheKey, hookItems, 0)
        }
    } else {
        hookItems = v.([]*handlerItem)
    }
    if hookItems != nil {
        for _, item := range hookItems {
            item.faddr(r)
        }
    }
}

// 生成hook key
func (s *Server) hookHandlerKey(hook, method, path, domain string) string {
    return hook + "@" + s.handlerKey(method, path, domain)
}

// 绑定指定的hook回调函数, pattern参数同BindHandler，支持命名路由；hook参数的值由ghttp server设定，参数不区分大小写
func (s *Server)BindHookHandler(pattern string, hook string, handler HandlerFunc) error {
    return s.setHookHandler(pattern, hook, &handlerItem{
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
