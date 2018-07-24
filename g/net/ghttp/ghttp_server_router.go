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
)

// handler缓存项，根据URL.Path进行缓存，因此对象中带有缓存参数
type handlerCacheItem struct {
    item    *HandlerItem         // 准确的执行方法内存地址
    values   map[string][]string // GET解析参数
}

// 查询请求处理方法
// 这里有个锁机制，可以并发读，但是不能并发写
func (s *Server) getHandler(r *Request) *HandlerItem {
    // 缓存清空时是直接修改属性，因此必须使用互斥锁
    s.hmcmu.RLock()
    defer s.hmcmu.RUnlock()

    var handlerItem *handlerCacheItem
    cacheKey := s.handlerKey(r.GetHost(), r.Method, r.URL.Path)
    if v := s.handlerCache.Get(cacheKey); v == nil {
        handlerItem = s.searchHandler(r)
        if handlerItem != nil {
            s.handlerCache.Set(cacheKey, handlerItem, 0)
        }
    } else {
        handlerItem = v.(*handlerCacheItem)
    }
    if handlerItem != nil {
        for k, v := range handlerItem.values {
            r.queries[k] = v
        }
        r.Router = handlerItem.item.router
        return handlerItem.item
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

// 路由注册处理方法
func (s *Server) setHandler(pattern string, item *HandlerItem) error {
    domain, method, uri, err := s.parsePattern(pattern)
    if err != nil {
        return errors.New("invalid pattern")
    }
    item.router = &Router {
        Uri    : uri,
        Domain : domain,
        Method : method,
    }
    s.hmmu.Lock()
    defer s.hmmu.Unlock()
    defer s.clearHandlerCache()
    if s.isPatternUriHasFuzzRule(uri) {
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
        array := strings.Split(uri[1:], "/")
        item.router.Priority = len(array)
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
        for _, l := range lists {
            for e := l.Front(); e != nil; e = e.Next() {
                if s.compareHandlerItemPriority(item, e.Value.(*HandlerItem)) {
                    l.InsertBefore(item, e)
                    break
                }
            }
            l.PushBack(item)
        }
    } else {
        // 静态注册
        if method == gDEFAULT_METHOD {
            for v, _ := range s.methodsMap {
                s.handlerMap[s.handlerKey(domain, v, uri)] = item
            }
        } else {
            s.handlerMap[s.handlerKey(domain, method, uri)] = item
        }
    }
    //b, _ := gparser.VarToJsonIndent(s.handlerTree)
    //fmt.Println(string(b))
    return nil
}

// 对比两个HandlerItem的优先级，需要非常注意的是，注意新老对比项的参数先后顺序
// 优先级比较规则：
// 1、层级越深优先级越高(对比/数量)；
// 2、模糊规则优先级：{xxx} > :xxx > *xxx；
func (s *Server) compareHandlerItemPriority(newItem, oldItem *HandlerItem) bool {
    if newItem.router.Priority > oldItem.router.Priority {
        return true
    }
    if newItem.router.Priority < oldItem.router.Priority {
        return false
    }
    // 例如：/{user}/{act} 比 /:user/:act 优先级高
    if strings.Count(newItem.router.Uri, "{") > strings.Count(oldItem.router.Uri, "{") {
        return true
    }
    // 例如: /:name/update 比 /:name/:action优先级高
    if strings.Count(newItem.router.Uri, "/:") < strings.Count(oldItem.router.Uri, "/:") {
        // 例如: /:name/:action 比 /:name/*any 优先级高
        if strings.Count(newItem.router.Uri, "/*") < strings.Count(oldItem.router.Uri, "/*") {
            return true
        }
        return false
    }
    return false
}

// 服务方法检索
func (s *Server) searchHandler(r *Request) *handlerCacheItem {
    item := s.searchHandlerStatic(r)
    if item == nil {
        item = s.searchHandlerDynamic(r)
    }
    return item
}

// 检索静态路由规则
func (s *Server) searchHandlerStatic(r *Request) *handlerCacheItem {
    s.hmmu.RLock()
    defer s.hmmu.RUnlock()
    domains := []string{r.GetHost(), gDEFAULT_DOMAIN}
    // 首先进行静态匹配
    for _, domain := range domains {
        if f, ok := s.handlerMap[s.handlerKey(domain, r.Method, r.URL.Path)]; ok {
            return &handlerCacheItem{f, nil}
        }
    }
    return nil
}

// 检索动态路由规则
func (s *Server) searchHandlerDynamic(r *Request) *handlerCacheItem {
    s.hmmu.RLock()
    defer s.hmmu.RUnlock()
    domains := []string{gDEFAULT_DOMAIN, r.GetHost()}
    array   := strings.Split(r.URL.Path[1:], "/")
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
                item := e.Value.(*HandlerItem)
                // 动态匹配规则带有gDEFAULT_METHOD的情况，不会像静态规则那样直接解析为所有的HTTP METHOD存储
                if strings.EqualFold(item.router.Method, gDEFAULT_METHOD) || strings.EqualFold(item.router.Method, r.Method) {
                    // 不管正则关键字符转义问题
                    //rule, names := s.patternToRegRule(gstr.ReplaceByMap(gregex.Quote(item.router.Uri), map[string]string {
                    //    `\{` : `{`,
                    //    `\}` : `}`,
                    //    `\*` : `*`,
                    //}))
                    rule, names := s.patternToRegRule(item.router.Uri)
                    if match, err := gregex.MatchString(rule, r.URL.Path); err == nil && len(match) > 1 {
                        //gutil.Dump(match)
                        //gutil.Dump(names)
                        handlerItem := &handlerCacheItem{item, nil}
                        // 如果需要query匹配，那么需要重新正则解析URL
                        if len(names) > 0 {
                            array := strings.Split(names, ",")
                            if len(match) > len(array) {
                                handlerItem.values = make(map[string][]string)
                                for index, name := range array {
                                    handlerItem.values[name] = []string{match[index + 1]}
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
func (s *Server) patternToRegRule(rule string) (regrule string, names string) {
    if len(rule) < 2 {
        return rule, ""
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
                if len(names) > 0 {
                    names += ","
                }
                names += v[1:]
            case '*':
                regrule += `/{0,1}(.*)`
                if len(names) > 0 {
                    names += ","
                }
                names += v[1:]
            default:
                s, _ := gregex.ReplaceStringFunc(`{[\w\.\-]+}`, v, func(s string) string {
                    if len(names) > 0 {
                        names += ","
                    }
                    names += s[1 : len(s) - 1]
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

// 判断URI中是否包含动态注册规则
func (s *Server) isPatternUriHasFuzzRule(uri string) bool {
    if len(uri) > 1 && gregex.IsMatchString(`^/[:\*]|{[\w\.\-]+}`, uri) {
        return true
    }
    return false
}

// 生成回调方法查询的Key
func (s *Server) handlerKey(domain, method, uri string) string {
    return strings.ToUpper(method) + ":" + uri + "@" + strings.ToLower(domain)
}

