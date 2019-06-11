<<<<<<< HEAD
// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.
// 路由控制.
=======
// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.
// 路由控制基本方法.
>>>>>>> upstream/master

package ghttp

import (
<<<<<<< HEAD
    "errors"
    "strings"
    "container/list"
    "gitee.com/johng/gf/g/util/gregx"
)

// handler缓存项，根据URL.Path进行缓存，因此对象中带有缓存参数
type handlerCacheItem struct {
    item    *HandlerItem        // 准确的执行方法内存地址
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
            r.values[k] = v
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
    if array, err := gregx.MatchString(`([a-zA-Z]+):(.+)`, pattern); len(array) > 1 && err == nil {
        method = array[1]
        uri    = array[2]
    }
    if array, err := gregx.MatchString(`(.+)@([\w\.\-]+)`, uri); len(array) > 1 && err == nil {
        uri     = array[1]
        domain  = array[2]
    }
    if uri == "" {
        err = errors.New("invalid pattern")
    }
    // 去掉末尾的"/"符号，与路由匹配时处理一直
    if uri != "/" {
        uri = strings.TrimRight(uri, "/")
=======
    "container/list"
    "errors"
    "fmt"
    "github.com/gogf/gf/g/os/glog"
    "github.com/gogf/gf/g/text/gregex"
    "github.com/gogf/gf/g/text/gstr"
    "runtime"
    "strings"
)


// 解析pattern
func (s *Server)parsePattern(pattern string) (domain, method, path string, err error) {
    path   = strings.TrimSpace(pattern)
    domain = gDEFAULT_DOMAIN
    method = gDEFAULT_METHOD
    if array, err := gregex.MatchString(`([a-zA-Z]+):(.+)`, pattern); len(array) > 1 && err == nil {
        path = strings.TrimSpace(array[2])
        if v := strings.TrimSpace(array[1]); v != "" {
            method = v
        }
    }
    if array, err := gregex.MatchString(`(.+)@([\w\.\-]+)`, path); len(array) > 1 && err == nil {
        path = strings.TrimSpace(array[1])
        if v := strings.TrimSpace(array[2]); v != "" {
            domain = v
        }
    }
    if path == "" {
        err = errors.New("invalid pattern: URI should not be empty")
    }
    // 去掉末尾的"/"符号，与路由匹配时处理一致
    if path != "/" {
        path = strings.TrimRight(path, "/")
>>>>>>> upstream/master
    }
    return
}

<<<<<<< HEAD
// 注册服务处理方法
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
    if s.isUriHasRule(uri) {
        // 动态注册，首先需要判断是否是动态注册，如果不是那么就没必要添加到动态注册记录变量中
        // 非叶节点为哈希表检索节点，按照URI注册的层级进行高效检索，直至到叶子链表节点；
        // 叶子节点是链表，按照优先级进行排序，优先级高的排前面，按照遍历检索，按照哈希表层级检索后的叶子链表一般数据量不大，所以效率比较高；
        if _, ok := s.handlerTree[domain]; !ok {
            s.handlerTree[domain] = make(map[string]interface{})
        }
        p     := s.handlerTree[domain]
        lists := make([]*list.List, 0)
        array := strings.Split(uri[1:], "/")
        item.router.Priority = len(array)
        for k, v := range array {
            if len(v) == 0 {
                continue
            }
            switch v[0] {
                case ':':
                    fallthrough
                case '*':
                    v = "/"
                    if v, ok := p.(map[string]interface{})["*list"]; !ok {
                        p.(map[string]interface{})["*list"] = list.New()
                        lists = append(lists, p.(map[string]interface{})["*list"].(*list.List))
                    } else {
                        lists = append(lists, v.(*list.List))
                    }
                    fallthrough
                default:
                    if _, ok := p.(map[string]interface{})[v]; !ok {
                        p.(map[string]interface{})[v] = make(map[string]interface{})
                    }
                    p = p.(map[string]interface{})[v]
                    // 到达叶子节点，往list中增加匹配规则
                    if v != "/" && k == len(array) - 1 {
                        if v, ok := p.(map[string]interface{})["*list"]; !ok {
                            p.(map[string]interface{})["*list"] = list.New()
                            lists = append(lists, p.(map[string]interface{})["*list"].(*list.List))
                        } else {
                            lists = append(lists, v.(*list.List))
                        }
                    }
            }
        }
        // 从头开始遍历链表，优先级高的放在前面
        for _, l := range lists {
            for e := l.Front(); e != nil; e = e.Next() {
                if s.compareHandlerItemPriority(item, e.Value.(*HandlerItem)) {
                    l.InsertBefore(item, e)
                    return nil
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
    return nil
}

// 对比两个HandlerItem的优先级，需要非常注意的是，注意新老对比项的参数先后顺序
func (s *Server) compareHandlerItemPriority(newItem, oldItem *HandlerItem) bool {
    if newItem.router.Priority > oldItem.router.Priority {
        return true
    }
    if newItem.router.Priority < oldItem.router.Priority {
        return false
    }
    if strings.Count(newItem.router.Uri, "/:") > strings.Count(oldItem.router.Uri, "/:") {
        return true
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
        // 多层链表的目的是当叶子节点未有任何规则匹配时，让父级模糊匹配规则继续处理
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
            }
            // 如果是叶子节点，同时判断当前层级的"/"键名，解决例如：/user/*action 匹配 /user 的规则
            if k == len(array) - 1 {
                if _, ok := p.(map[string]interface{})["/"]; ok {
                    p = p.(map[string]interface{})["/"]
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
                        handlerItem := &handlerCacheItem{item, nil}
                        // 如果需要query匹配，那么需要重新解析URL
                        if len(names) > 0 {
                            if match, err := gregx.MatchString(regrule, r.URL.Path); err == nil {
                                array := strings.Split(names, ",")
                                if len(match) > len(array) {
                                    handlerItem.values = make(map[string][]string)
                                    for index, name := range array {
                                        handlerItem.values[name] = []string{match[index + 1]}
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
func (s *Server) patternToRegRule(rule string) (regrule string, names string) {
    if len(rule) < 2 {
        return rule, ""
=======
// 获得服务注册的文件地址信息
func (s *Server) getHandlerRegisterCallerLine(handler *handlerItem) string {
    skip := 5
    if handler.rtype == gROUTE_REGISTER_HANDLER {
       skip = 4
    }
    if _, cfile, cline, ok := runtime.Caller(skip); ok {
        return fmt.Sprintf("%s:%d", cfile, cline)
    }
    return ""
}

// 路由注册处理方法。
// 如果带有hook参数，表示是回调注册方法; 否则为普通路由执行方法。
func (s *Server) setHandler(pattern string, handler *handlerItem, hook ... string) {
    // Web Server正常运行时无法动态注册路由方法
    if s.Status() == SERVER_STATUS_RUNNING {
        glog.Error("cannot bind handler while server running")
        return
    }
    var hookName string
    if len(hook) > 0 {
        hookName = hook[0]
    }
    domain, method, uri, err := s.parsePattern(pattern)
    if err != nil {
        glog.Error("invalid pattern:", pattern, err)
        return
    }
    if len(uri) == 0 || uri[0] != '/' {
        glog.Error("invalid pattern:", pattern, "URI should lead with '/'")
        return
    }
    // 注册地址记录及重复注册判断
    regkey := s.handlerKey(hookName, method, uri, domain)
    caller := s.getHandlerRegisterCallerLine(handler)
    if len(hook) == 0 {
        if item, ok := s.routesMap[regkey]; ok {
            glog.Errorf(`duplicated route registry "%s", already registered at %s`, pattern, item[0].file)
            return
        }
    }

    // 路由对象
    handler.router = &Router {
        Uri      : uri,
        Domain   : domain,
        Method   : method,
        Priority : strings.Count(uri[1:], "/"),
    }
    handler.router.RegRule, handler.router.RegNames = s.patternToRegRule(uri)

    // 动态注册，首先需要判断是否是动态注册，如果不是那么就没必要添加到动态注册记录变量中。
    // 非叶节点为哈希表检索节点，按照URI注册的层级进行高效检索，直至到叶子链表节点；
    // 叶子节点是链表，按照优先级进行排序，优先级高的排前面，按照遍历检索，按照哈希表层级检索后的叶子链表数据量不会很大，所以效率比较高；
    tree := (map[string]interface{})(nil)
    if len(hookName) == 0 {
        tree = s.serveTree
    } else {
        tree = s.hooksTree
    }
    if _, ok := tree[domain]; !ok {
        tree[domain] = make(map[string]interface{})
    }
    // 用于遍历的指针
    p := tree[domain]
    if len(hookName) > 0 {
        if _, ok := p.(map[string]interface{})[hookName]; !ok {
            p.(map[string]interface{})[hookName] = make(map[string]interface{})
        }
        p = p.(map[string]interface{})[hookName]
    }
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
        if gregex.IsMatchString(`^[:\*]|\{[\w\.\-]+\}|\*`, v) {
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
    // 上面循环后得到的lists是该路由规则一路匹配下来相关的模糊匹配链表(注意不是这棵树所有的链表)。
    // 下面从头开始遍历每个节点的模糊匹配链表，将该路由项插入进去(按照优先级高的放在lists链表的前面)
    item := (*handlerItem)(nil)
    for _, l := range lists {
        pushed  := false
        for e := l.Front(); e != nil; e = e.Next() {
            item = e.Value.(*handlerItem)
            // 判断是否已存在相同的路由注册项，(如果不是hook注册)是则进行替换
            if len(hookName) == 0 {
                if strings.EqualFold(handler.router.Domain, item.router.Domain) &&
                    strings.EqualFold(handler.router.Method, item.router.Method) &&
                    strings.EqualFold(handler.router.Uri, item.router.Uri) {
                    e.Value = handler
                    pushed  = true
                    break
                }
            }
            // 如果路由注册项不相等，那么判断优先级，决定插入顺序
            if s.compareRouterPriority(handler.router, item.router) {
                l.InsertBefore(handler, e)
                pushed = true
                break
            }
        }
        if !pushed {
            l.PushBack(handler)
        }
    }
    // gutil.Dump(s.serveTree)
    // gutil.Dump(s.hooksTree)
    if _, ok := s.routesMap[regkey]; !ok {
        s.routesMap[regkey] = make([]registeredRouteItem, 0)
    }
    s.routesMap[regkey] = append(s.routesMap[regkey], registeredRouteItem {
        file    : caller,
        handler : handler,
    })
}

// 对比两个handlerItem的优先级，需要非常注意的是，注意新老对比项的参数先后顺序。
// 返回值true表示newRouter优先级比oldRouter高，会被添加链表中oldRouter的前面；否则后面。
// 优先级比较规则：
// 1、层级越深优先级越高(对比/数量)；
// 2、模糊规则优先级：{xxx} > :xxx > *xxx；
func (s *Server) compareRouterPriority(newRouter, oldRouter *Router) bool {
    // 优先比较层级，层级越深优先级越高
    if newRouter.Priority > oldRouter.Priority {
        return true
    }
    if newRouter.Priority < oldRouter.Priority {
        return false
    }
    // 精准匹配比模糊匹配规则优先级高，例如：/name/act 比 /{name}/:act 优先级高
    var fuzzyCountFieldNew, fuzzyCountFieldOld int
    var fuzzyCountNameNew,  fuzzyCountNameOld  int
    var fuzzyCountAnyNew,   fuzzyCountAnyOld   int
    var fuzzyCountTotalNew, fuzzyCountTotalOld int
    for _, v := range newRouter.Uri {
        switch v {
            case '{':
                fuzzyCountFieldNew++
            case ':':
                fuzzyCountNameNew++
            case '*':
                fuzzyCountAnyNew++
        }
    }
    for _, v := range oldRouter.Uri {
        switch v {
            case '{':
                fuzzyCountFieldOld++
            case ':':
                fuzzyCountNameOld++
            case '*':
                fuzzyCountAnyOld++
        }
    }
    fuzzyCountTotalNew = fuzzyCountFieldNew + fuzzyCountNameNew + fuzzyCountAnyNew
    fuzzyCountTotalOld = fuzzyCountFieldOld + fuzzyCountNameOld + fuzzyCountAnyOld
    if fuzzyCountTotalNew < fuzzyCountTotalOld {
        return true
    }
    if fuzzyCountTotalNew > fuzzyCountTotalOld {
        return false
    }

    /** 如果模糊规则数量相等，那么执行分别的数量判断 **/

    // 例如：/name/{act} 比 /name/:act 优先级高
    if fuzzyCountFieldNew > fuzzyCountFieldOld {
        return true
    }
    if fuzzyCountFieldNew < fuzzyCountFieldOld {
        return false
    }
    // 例如: /name/:act 比 /name/*act 优先级高
    if fuzzyCountNameNew > fuzzyCountNameOld {
        return true
    }
    if fuzzyCountNameNew < fuzzyCountNameOld {
        return false
    }

    /* 模糊规则数量相等，后续不用再判断*规则的数量比较了 */

    // 比较HTTP METHOD，更精准的优先级更高
    if newRouter.Method != gDEFAULT_METHOD {
        return true
    }
    if oldRouter.Method != gDEFAULT_METHOD {
        return true
    }

    // 最后新的规则比旧的规则优先级低
    return false
}

// 将pattern（不带method和domain）解析成正则表达式匹配以及对应的query字符串
func (s *Server) patternToRegRule(rule string) (regrule string, names []string) {
    if len(rule) < 2 {
        return rule, nil
>>>>>>> upstream/master
    }
    regrule = "^"
    array  := strings.Split(rule[1:], "/")
    for _, v := range array {
        if len(v) == 0 {
            continue
        }
        switch v[0] {
            case ':':
<<<<<<< HEAD
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
                regrule += "/" + v
=======
                if len(v) > 1 {
                    regrule += `/([^/]+)`
                    names    = append(names, v[1:])
                    break
                } else {
                    regrule += `/[^/]+`
                    break
                }
            case '*':
                if len(v) > 1 {
                    regrule += `/{0,1}(.*)`
                    names    = append(names, v[1:])
                    break
                } else {
                    regrule += `/{0,1}.*`
                    break
                }
            default:
                // 特殊字符替换
                v = gstr.ReplaceByMap(v, map[string]string{
                    `.` : `\.`,
                    `+` : `\+`,
                    `*` : `.*`,
                })
                s, _ := gregex.ReplaceStringFunc(`\{[\w\.\-]+\}`, v, func(s string) string {
                    names = append(names, s[1 : len(s) - 1])
                    return `([^/]+)`
                })
                if strings.EqualFold(s, v) {
                    regrule += "/" + v
                } else {
                    regrule += "/" + s
                }
>>>>>>> upstream/master
        }
    }
    regrule += `$`
    return
}

<<<<<<< HEAD
// 判断URI中是否包含动态注册规则
func (s *Server) isUriHasRule(uri string) bool {
    if len(uri) > 1 && (strings.Index(uri, "/:") != -1 || strings.Index(uri, "/*") != -1) {
        return true
    }
    return false
}

// 生成回调方法查询的Key
func (s *Server) handlerKey(domain, method, uri string) string {
    return strings.ToUpper(method) + ":" + uri + "@" + strings.ToLower(domain)
}

=======
>>>>>>> upstream/master
