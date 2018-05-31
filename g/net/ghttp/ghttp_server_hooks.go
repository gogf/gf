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
    "gitee.com/johng/gf/g/util/gregx"
)

// hook缓存项，根据URL.Path进行缓存，因此对象中带有缓存参数
type hookCacheItem struct {
    faddr    HandlerFunc         // 准确的执行方法内存地址
    values   map[string][]string // GET解析参数
}

// 事件回调注册方法
// 因为有事件回调优先级的关系，叶子节点必须为一个链表，因此这里只有动态注册
func (s *Server) setHookHandler(pattern string, hook string, item *HandlerItem) error {
    domain, method, uri, err := s.parsePattern(pattern)
    if err != nil {
        return errors.New("invalid pattern")
    }
    item.router = &Router {
        Uri    : uri,
        Domain : domain,
        Method : method,
    }

    s.hhmu.Lock()
    defer s.hhmu.Unlock()
    defer s.clearHooksCache()

    if _, ok := s.hooksTree[domain]; !ok {
        s.hooksTree[domain] = make(map[string]interface{})
    }
    p := s.hooksTree[domain]
    if _, ok := p.(map[string]interface{})[hook]; !ok {
        p.(map[string]interface{})[hook] = make(map[string]interface{})
    }
    p = p.(map[string]interface{})[hook]

    array := strings.Split(uri[1:], "/")
    item.router.Priority = len(array)
    for _, v := range array {
        if len(v) == 0 {
            continue
        }
        switch v[0] {
            case ':':
                fallthrough
            case '*':
                v = "/"
                fallthrough
            default:
                if _, ok := p.(map[string]interface{})[v]; !ok {
                    p.(map[string]interface{})[v] = make(map[string]interface{})
                }
                p = p.(map[string]interface{})[v]
        }
    }
    // 到达叶子节点
    var l *list.List
    if v, ok := p.(map[string]interface{})["*list"]; !ok {
        l = list.New()
        p.(map[string]interface{})["*list"] = l
    } else {
        l = v.(*list.List)
    }
    // 从头开始遍历链表，优先级高的放在前面
    for e := l.Front(); e != nil; e = e.Next() {
        if s.compareHandlerItemPriority(item, e.Value.(*HandlerItem)) {
            l.InsertBefore(item, e)
            return nil
        }
    }
    l.PushBack(item)

    return nil
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
    s.hhmu.RLock()
    defer s.hhmu.RUnlock()
    hookItems := make([]*hookCacheItem, 0)
    domains   := []string{r.GetHost(), gDEFAULT_DOMAIN}
    array     := strings.Split(r.URL.Path[1:], "/")
    for _, domain := range domains {
        p, ok := s.hooksTree[domain]
        if !ok {
            continue
        }
        p, ok = p.(map[string]interface{})[hook]
        if !ok {
            continue
        }
        // 多层链表的目的是当叶子节点未有任何规则匹配时，让父级模糊匹配规则继续处理
        lists := make([]*list.List, 0)
        for k, v := range array {
            if _, ok := p.(map[string]interface{})["*list"]; ok {
                lists = append(lists, p.(map[string]interface{})["*list"].(*list.List))
            }
            if _, ok := p.(map[string]interface{})[v]; !ok {
                if _, ok := p.(map[string]interface{})["/"]; ok {
                    p = p.(map[string]interface{})["/"]
                    if k == len(array) - 1 {
                        if _, ok := p.(map[string]interface{})["*list"]; ok {
                            lists = append(lists, p.(map[string]interface{})["*list"].(*list.List))
                        }
                    }
                } else {
                    break
                }
            } else {
                p = p.(map[string]interface{})[v]
                if k == len(array) - 1 {
                    if _, ok := p.(map[string]interface{})["*list"]; ok {
                        lists = append(lists, p.(map[string]interface{})["*list"].(*list.List))
                    }
                }
            }
        }

        // 多层链表遍历检索，从数组末尾的链表开始遍历，末尾的深度高优先级也高
        for i := len(lists) - 1; i >= 0; i-- {
            for e := lists[i].Front(); e != nil; e = e.Next() {
                item := e.Value.(*HandlerItem)
                if strings.EqualFold(item.router.Method, gDEFAULT_METHOD) || strings.EqualFold(item.router.Method, r.Method) {
                    regrule, names := s.patternToRegRule(item.router.Uri)
                    if gregx.IsMatchString(regrule, r.URL.Path) {
                        hookItem := &hookCacheItem {item.faddr, nil}
                        // 如果需要query匹配，那么需要重新解析URL
                        if len(names) > 0 {
                            if match, err := gregx.MatchString(regrule, r.URL.Path); err == nil {
                                array := strings.Split(names, ",")
                                if len(match) > len(array) {
                                    hookItem.values = make(map[string][]string)
                                    // 这里需要注意的是，注册事件回调如果带有规则匹配，那么会修改Request对象传递参数的值
                                    // 这个应当在注册事件回调的时候注意
                                    for index, name := range array {
                                        hookItem.values[name] = []string{match[index + 1]}
                                    }
                                }
                            }
                        }
                        hookItems = append(hookItems, hookItem)
                    }
                }
            }
        }
    }
    return hookItems
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
