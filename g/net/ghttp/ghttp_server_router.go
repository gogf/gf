// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.
// 路由控制.

package ghttp

import (
    "errors"
    "strings"
    "container/list"
    "gitee.com/johng/gf/g/util/gregex"
    "gitee.com/johng/gf/g/container/gset"
    "fmt"
)

// 查询请求处理方法.
// 内部带锁机制，可以并发读，但是不能并发写；并且有缓存机制，按照Host、Method、Path进行缓存.
func (s *Server) getHandlerWithCache(r *Request) *handlerRegisterItem {
    var cacheItem *handlerParsedItem
    cacheKey := s.handlerKey(r.Method, r.URL.Path, r.GetHost())
    if v := s.handlerCache.Get(cacheKey); v == nil {
        cacheItem = s.searchHandler(r.Method, r.URL.Path, r.GetHost())
        if cacheItem != nil {
            s.handlerCache.Set(cacheKey, cacheItem, 0)
        }
    } else {
        cacheItem = v.(*handlerParsedItem)
    }

    if cacheItem != nil {
        if r.Router == nil {
            for k, v := range cacheItem.values {
                r.routerVars[k] = v
            }
            r.Router = cacheItem.item.router
        }
        return cacheItem.item
    }
    return nil
}

// 解析pattern
func (s *Server)parsePattern(pattern string) (domain, method, uri string, err error) {
    uri    = pattern
    domain = gDEFAULT_DOMAIN
    method = gDEFAULT_METHOD
    if array, err := gregex.MatchString(`([a-zA-Z]+):(.+)`, pattern); len(array) > 1 && err == nil {
        method = array[1]
        uri    = array[2]
    }
    if array, err := gregex.MatchString(`(.+)@([\w\.\-]+)`, uri); len(array) > 1 && err == nil {
        uri     = array[1]
        domain  = array[2]
    }
    if uri == "" {
        err = errors.New("invalid pattern")
    }
    // 去掉末尾的"/"符号，与路由匹配时处理一直
    if uri != "/" {
        uri = strings.TrimRight(uri, "/")
    }
    return
}

// 路由注册处理方法。
// 如果带有hook参数，表示是回调注册方法，否则为普通路由执行方法。
func (s *Server) setHandler(pattern string, handler *handlerItem, hook ... string) error {
    // Web Server正字运行时无法动态注册路由方法
    if s.status == gSERVER_STATUS_RUNNING {
        return errors.New("cannnot bind handler while server running")
    }
    var hookName string
    if len(hook) > 0 {
        hookName = hook[0]
    }
    domain, method, uri, err := s.parsePattern(pattern)
    if err != nil {
        return errors.New("invalid pattern")
    }

    // 路由对象
    router := &Router {
        Uri      : uri,
        Domain   : domain,
        Method   : method,
        Priority : strings.Count(uri[1:], "/"),
    }
    router.RegRule, router.RegNames = s.patternToRegRule(uri)

    // 注册对象
    registerItem := &handlerRegisterItem {
        handler : handler,
        hooks   : make(map[string]*list.List),
        router  : router,
    }
    if len(hookName) > 0 {
        registerItem.handler         = nil
        registerItem.hooks[hookName] = list.New()
        registerItem.hooks[hookName].PushBack(handler)
    }

    // 动态注册，首先需要判断是否是动态注册，如果不是那么就没必要添加到动态注册记录变量中。
    // 非叶节点为哈希表检索节点，按照URI注册的层级进行高效检索，直至到叶子链表节点；
    // 叶子节点是链表，按照优先级进行排序，优先级高的排前面，按照遍历检索，按照哈希表层级检索后的叶子链表数据量不会很大，所以效率比较高；
    if _, ok := s.handlerTree[domain]; !ok {
        s.handlerTree[domain] = make(map[string]interface{})
    }
    // 用于遍历的指针
    p     := s.handlerTree[domain]
    // 当前节点的规则链表
    lists := make([]*list.List, 0)
    array := ([]string)(nil)
    if strings.EqualFold("/", uri) {
        array = []string{"/"}
    } else {
        array = strings.Split(uri[1:], "/")
    }
    // 键名"*fuzz"代表模糊匹配节点，其下会有一个链表；
    // 键名"*list"代表链表，叶子节点和模糊匹配节点都有该属性；
    for k, v := range array {
        if len(v) == 0 {
            continue
        }
        // 判断是否模糊匹配规则
        if gregex.IsMatchString(`^[:\*]|{[\w\.\-]+}`, v) {
            v = "*fuzz"
            // 由于是模糊规则，因此这里会有一个*list，用以将后续的路由规则加进来，
            // 检索会从叶子节点的链表往根节点按照优先级进行检索
            if v, ok := p.(map[string]interface{})["*list"]; !ok {
                p.(map[string]interface{})["*list"] = list.New()
                lists = append(lists, p.(map[string]interface{})["*list"].(*list.List))
            } else {
                lists = append(lists, v.(*list.List))
            }
        }
        // 属性层级数据写入
        if _, ok := p.(map[string]interface{})[v]; !ok {
            p.(map[string]interface{})[v] = make(map[string]interface{})
        }
        p = p.(map[string]interface{})[v]
        // 到达叶子节点，往list中增加匹配规则(条件 v != "*fuzz" 是因为模糊节点的话在前面已经添加了*list链表)
        if k == len(array) - 1 && v != "*fuzz" {
            if v, ok := p.(map[string]interface{})["*list"]; !ok {
                p.(map[string]interface{})["*list"] = list.New()
                lists = append(lists, p.(map[string]interface{})["*list"].(*list.List))
            } else {
                lists = append(lists, v.(*list.List))
            }
        }
    }
    // 得到的lists是该路由规则一路匹配下来相关的模糊匹配链表(注意不是这棵树所有的链表)，
    // 从头开始遍历每个节点的模糊匹配链表，将该路由项插入进去(按照优先级高的放在前面)
    item := (*handlerRegisterItem)(nil)
    // 用以标记 *handlerRegisterItem 指向的对象是否已经处理过，因为多个节点可能会关联同一个该对象
    pushedItemSet := gset.NewStringSet()
    if len(hookName) == 0 {
        // 普通方法路由注册，追加或者覆盖
        for _, l := range lists {
            pushed  := false
            address := ""
            for e := l.Front(); e != nil; e = e.Next() {
                item    = e.Value.(*handlerRegisterItem)
                address = fmt.Sprintf("%p", item)
                if pushedItemSet.Contains(address) {
                    pushed = true
                    break
                }
                // 判断是否已存在相同的路由注册项
                if strings.EqualFold(router.Domain, item.router.Domain) &&
                    strings.EqualFold(router.Method, item.router.Method) &&
                    strings.EqualFold(router.Uri, item.router.Uri) {
                    item.handler = handler
                    pushed = true
                    break
                }
                if s.compareRouterPriority(router, item.router) {
                    l.InsertBefore(registerItem, e)
                    pushed = true
                    break
                }
            }
            if pushed {
                if len(address) > 0 {
                    pushedItemSet.Add(address)
                }
            } else {
                l.PushBack(registerItem)
            }
        }
    } else {
        // 回调方法路由注册，将方法追加到链表末尾
        for _, l := range lists {
            pushed  := false
            address := ""
            for e := l.Front(); e != nil; e = e.Next() {
                item    = e.Value.(*handlerRegisterItem)
                address = fmt.Sprintf("%p", item)
                if pushedItemSet.Contains(address) {
                    pushed = true
                    break
                }
                // 判断是否已存在相同的路由注册项
                if strings.EqualFold(router.Domain, item.router.Domain) &&
                    strings.EqualFold(router.Method, item.router.Method) &&
                    strings.EqualFold(router.Uri, item.router.Uri) {
                    if _, ok := item.hooks[hookName]; !ok {
                        item.hooks[hookName] = list.New()
                    }
                    item.hooks[hookName].PushBack(handler)
                    pushed = true
                    break
                }
                if s.compareRouterPriority(router, item.router) {
                    l.InsertBefore(registerItem, e)
                    pushed = true
                    break
                }
            }
            if pushed {
                if len(address) > 0 {
                    pushedItemSet.Add(address)
                }
            } else {
                l.PushBack(registerItem)
            }
        }
    }
    //gutil.Dump(s.handlerTree)
    return nil
}

// 对比两个handlerItem的优先级，需要非常注意的是，注意新老对比项的参数先后顺序
// 优先级比较规则：
// 1、层级越深优先级越高(对比/数量)；
// 2、模糊规则优先级：{xxx} > :xxx > *xxx；
func (s *Server) compareRouterPriority(newRouter, oldRouter *Router) bool {
    if newRouter.Priority > oldRouter.Priority {
        return true
    }
    if newRouter.Priority < oldRouter.Priority {
        return false
    }
    // 例如：/{user}/{act} 比 /:user/:act 优先级高
    if strings.Count(newRouter.Uri, "{") > strings.Count(oldRouter.Uri, "{") {
        return true
    }
    // 例如: /:name/update 比 /:name/:action优先级高
    if strings.Count(newRouter.Uri, "/:") < strings.Count(oldRouter.Uri, "/:") {
        // 例如: /:name/:action 比 /:name/*any 优先级高
        if strings.Count(newRouter.Uri, "/*") < strings.Count(oldRouter.Uri, "/*") {
            return true
        }
        return false
    }
    return false
}

// 服务方法检索
func (s *Server) searchHandler(method, path, domain string) *handlerParsedItem {
    domains := []string{ gDEFAULT_DOMAIN }
    if !strings.EqualFold(gDEFAULT_DOMAIN, domain) {
        domains = append(domains, domain)
    }
    array := ([]string)(nil)
    if strings.EqualFold("/", path) {
        array = []string{"/"}
    } else {
        array = strings.Split(path[1:], "/")
    }
    for _, domain := range domains {
        p, ok := s.handlerTree[domain]
        if !ok {
            continue
        }
        // 多层链表(每个节点都有一个*list链表)的目的是当叶子节点未有任何规则匹配时，让父级模糊匹配规则继续处理
        lists := make([]*list.List, 0)
        for k, v := range array {
            if _, ok := p.(map[string]interface{})["*list"]; ok {
                lists = append(lists, p.(map[string]interface{})["*list"].(*list.List))
            }
            if _, ok := p.(map[string]interface{})[v]; ok {
                p = p.(map[string]interface{})[v]
                if k == len(array) - 1 {
                    if _, ok := p.(map[string]interface{})["*list"]; ok {
                        lists = append(lists, p.(map[string]interface{})["*list"].(*list.List))
                        break
                    }
                }
            } else {
                if _, ok := p.(map[string]interface{})["*fuzz"]; ok {
                    p = p.(map[string]interface{})["*fuzz"]
                }
            }
            // 如果是叶子节点，同时判断当前层级的"*fuzz"键名，解决例如：/user/*action 匹配 /user 的规则
            if k == len(array) - 1 {
                if _, ok := p.(map[string]interface{})["*fuzz"]; ok {
                    p = p.(map[string]interface{})["*fuzz"]
                }
                if _, ok := p.(map[string]interface{})["*list"]; ok {
                    lists = append(lists, p.(map[string]interface{})["*list"].(*list.List))
                }
            }
        }

        // 多层链表遍历检索，从数组末尾的链表开始遍历，末尾的深度高优先级也高
        for i := len(lists) - 1; i >= 0; i-- {
            for e := lists[i].Front(); e != nil; e = e.Next() {
                item := e.Value.(*handlerRegisterItem)
                // 动态匹配规则带有gDEFAULT_METHOD的情况，不会像静态规则那样直接解析为所有的HTTP METHOD存储
                if strings.EqualFold(item.router.Method, gDEFAULT_METHOD) || strings.EqualFold(item.router.Method, method) {
                    // 注意当不带任何动态路由规则时，len(match) == 1
                    if match, err := gregex.MatchString(item.router.RegRule, path); err == nil && len(match) > 0 {
                        //gutil.Dump(match)
                        //gutil.Dump(names)
                        handlerItem := &handlerParsedItem{item, nil}
                        // 如果需要query匹配，那么需要重新正则解析URL
                        if len(item.router.RegNames) > 0 {
                            if len(match) > len(item.router.RegNames) {
                                handlerItem.values = make(map[string][]string)
                                // 如果存在存在同名路由参数名称，那么执行数组追加
                                for i, name := range item.router.RegNames {
                                    if _, ok := handlerItem.values[name]; ok {
                                        handlerItem.values[name] = append(handlerItem.values[name], match[i + 1])
                                    } else {
                                        handlerItem.values[name] = []string{match[i + 1]}
                                    }
                                }
                            }
                        }
                        return handlerItem
                    }
                }
            }
        }
    }
    return nil
}

// 将pattern（不带method和domain）解析成正则表达式匹配以及对应的query字符串
func (s *Server) patternToRegRule(rule string) (regrule string, names []string) {
    if len(rule) < 2 {
        return rule, nil
    }
    regrule = "^"
    array  := strings.Split(rule[1:], "/")
    for _, v := range array {
        if len(v) == 0 {
            continue
        }
        switch v[0] {
            case ':':
                regrule += `/([\w\.\-]+)`
                names    = append(names, v[1:])
            case '*':
                regrule += `/{0,1}(.*)`
                names    = append(names, v[1:])
            default:
                s, _ := gregex.ReplaceStringFunc(`{[\w\.\-]+}`, v, func(s string) string {
                    names = append(names, s[1 : len(s) - 1])
                    return `([\w\.\-]+)`
                })
                if strings.EqualFold(s, v) {
                    regrule += "/" + v
                } else {
                    regrule += "/" + s
                }
        }
    }
    regrule += `$`
    return
}

// 生成回调方法查询的Key
func (s *Server) handlerKey(method, path, domain string) string {
    return strings.ToUpper(method) + ":" + path + "@" + strings.ToLower(domain)
}

