// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.
// 服务注册路由控制.

package ghttp

import (
    "strings"
    "container/list"
    "gitee.com/johng/gf/g/util/gregex"
)

// 查询请求处理方法.
// 内部带锁机制，可以并发读，但是不能并发写；并且有缓存机制，按照Host、Method、Path进行缓存.
func (s *Server) getServeHandlerWithCache(r *Request) *handlerParsedItem {
    cacheItem := (*handlerParsedItem)(nil)
    cacheKey  := s.serveHandlerKey(r.Method, r.URL.Path, r.GetHost())
    if v := s.serveCache.Get(cacheKey); v == nil {
        cacheItem = s.searchServeHandler(r.Method, r.URL.Path, r.GetHost())
        if cacheItem != nil {
            s.serveCache.Set(cacheKey, cacheItem, s.config.RouterCacheExpire*1000)
        }
    } else {
        cacheItem = v.(*handlerParsedItem)
    }
    return cacheItem
}

// 服务方法检索
func (s *Server) searchServeHandler(method, path, domain string) *handlerParsedItem {
    if len(path) == 0 {
        return nil
    }
    // 遍历检索的域名列表
    domains := []string{ gDEFAULT_DOMAIN }
    if !strings.EqualFold(gDEFAULT_DOMAIN, domain) {
        domains = append(domains, domain)
    }
    // URL.Path层级拆分
    array := ([]string)(nil)
    if strings.EqualFold("/", path) {
        array = []string{"/"}
    } else {
        array = strings.Split(path[1:], "/")
    }
    for _, domain := range domains {
        p, ok := s.serveTree[domain]
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
                item := e.Value.(*handlerItem)
                // 动态匹配规则带有gDEFAULT_METHOD的情况，不会像静态规则那样直接解析为所有的HTTP METHOD存储
                if strings.EqualFold(item.router.Method, gDEFAULT_METHOD) || strings.EqualFold(item.router.Method, method) {
                    // 注意当不带任何动态路由规则时，len(match) == 1
                    if match, err := gregex.MatchString(item.router.RegRule, path); err == nil && len(match) > 0 {
                        //gutil.Dump(match)
                        //gutil.Dump(names)
                        parsedItem := &handlerParsedItem{item, nil}
                        // 如果需要query匹配，那么需要重新正则解析URL
                        if len(item.router.RegNames) > 0 {
                            if len(match) > len(item.router.RegNames) {
                                parsedItem.values = make(map[string][]string)
                                // 如果存在存在同名路由参数名称，那么执行数组追加
                                for i, name := range item.router.RegNames {
                                    if _, ok := parsedItem.values[name]; ok {
                                        parsedItem.values[name] = append(parsedItem.values[name], match[i + 1])
                                    } else {
                                        parsedItem.values[name] = []string{match[i + 1]}
                                    }
                                }
                            }
                        }
                        return parsedItem
                    }
                }
            }
        }
    }
    return nil
}

// 生成回调方法查询的Key
func (s *Server) serveHandlerKey(method, path, domain string) string {
    return strings.ToUpper(method) + ":" + path + "@" + strings.ToLower(domain)
}

