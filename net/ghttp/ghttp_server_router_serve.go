// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gogf/gf/container/glist"
	"github.com/gogf/gf/text/gregex"
)

// 缓存数据项
type handlerCacheItem struct {
	parsedItems []*handlerParsedItem
	hasHook     bool
	hasServe    bool
}

// 查询请求处理方法.
// 内部带锁机制，可以并发读，但是不能并发写；并且有缓存机制，按照Host、Method、Path进行缓存.
func (s *Server) getHandlersWithCache(r *Request) (parsedItems []*handlerParsedItem, hasHook, hasServe bool) {
	value := s.serveCache.GetOrSetFunc(s.serveHandlerKey(r.Method, r.URL.Path, r.GetHost()), func() interface{} {
		parsedItems, hasHook, hasServe = s.searchHandlers(r.Method, r.URL.Path, r.GetHost())
		if parsedItems != nil {
			return &handlerCacheItem{parsedItems, hasHook, hasServe}
		}
		return nil
	}, gROUTE_CACHE_DURATION)
	if value != nil {
		item := value.(*handlerCacheItem)
		return item.parsedItems, item.hasHook, item.hasServe
	}
	return
}

// 路由注册方法检索，返回所有该路由的注册函数，构造成数组返回
func (s *Server) searchHandlers(method, path, domain string) (parsedItems []*handlerParsedItem, hasHook, hasServe bool) {
	if len(path) == 0 {
		return nil, false, false
	}
	// 遍历检索的域名列表，优先遍历默认域名
	domains := []string{gDEFAULT_DOMAIN}
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
	parsedItemList := glist.New()
	lastMiddlewareElem := (*glist.Element)(nil)
	repeatHandlerCheckMap := make(map[int]struct{})
	for _, domain := range domains {
		p, ok := s.serveTree[domain]
		if !ok {
			continue
		}
		// 多层链表(每个节点都有一个*list链表)的目的是当叶子节点未有任何规则匹配时，让父级模糊匹配规则继续处理
		lists := make([]*glist.List, 0, 16)
		for k, v := range array {
			// In case of double '/' URI, eg: /user//index
			if v == "" {
				continue
			}
			if _, ok := p.(map[string]interface{})["*list"]; ok {
				lists = append(lists, p.(map[string]interface{})["*list"].(*glist.List))
			}
			if _, ok := p.(map[string]interface{})[v]; ok {
				p = p.(map[string]interface{})[v]
				if k == len(array)-1 {
					if _, ok := p.(map[string]interface{})["*list"]; ok {
						lists = append(lists, p.(map[string]interface{})["*list"].(*glist.List))
						break
					}
				}
			} else {
				if _, ok := p.(map[string]interface{})["*fuzz"]; ok {
					p = p.(map[string]interface{})["*fuzz"]
				}
			}
			// 如果是叶子节点，同时判断当前层级的"*fuzz"键名，解决例如：/user/*action 匹配 /user 的规则
			if k == len(array)-1 {
				if _, ok := p.(map[string]interface{})["*fuzz"]; ok {
					p = p.(map[string]interface{})["*fuzz"]
				}
				if _, ok := p.(map[string]interface{})["*list"]; ok {
					lists = append(lists, p.(map[string]interface{})["*list"].(*glist.List))
				}
			}
		}

		// 多层链表遍历检索，从数组末尾的链表开始遍历，末尾的深度高优先级也高
		for i := len(lists) - 1; i >= 0; i-- {
			for e := lists[i].Front(); e != nil; e = e.Next() {
				item := e.Value.(*handlerItem)
				// 主要是用于路由注册函数的重复添加判断(特别是中间件和钩子函数)
				if _, ok := repeatHandlerCheckMap[item.itemId]; ok {
					continue
				} else {
					repeatHandlerCheckMap[item.itemId] = struct{}{}
				}
				// 服务路由函数只能添加一次，将重复判断放在这里提高检索效率
				if hasServe {
					switch item.itemType {
					case gHANDLER_TYPE_HANDLER, gHANDLER_TYPE_OBJECT, gHANDLER_TYPE_CONTROLLER:
						continue
					}
				}
				// 动态匹配规则带有gDEFAULT_METHOD的情况，不会像静态规则那样直接解析为所有的HTTP METHOD存储
				if strings.EqualFold(item.router.Method, gDEFAULT_METHOD) || strings.EqualFold(item.router.Method, method) {
					// 注意当不带任何动态路由规则时，len(match) == 1
					if match, err := gregex.MatchString(item.router.RegRule, path); err == nil && len(match) > 0 {
						parsedItem := &handlerParsedItem{item, nil}
						// 如果需要路由规则中带有URI名称匹配，那么需要重新正则解析URL
						if len(item.router.RegNames) > 0 {
							if len(match) > len(item.router.RegNames) {
								parsedItem.values = make(map[string]string)
								// 如果存在存在同名路由参数名称，那么执行覆盖
								for i, name := range item.router.RegNames {
									parsedItem.values[name] = match[i+1]
								}
							}
						}
						switch item.itemType {
						// 服务路由函数只能添加一次
						case gHANDLER_TYPE_HANDLER, gHANDLER_TYPE_OBJECT, gHANDLER_TYPE_CONTROLLER:
							hasServe = true
							parsedItemList.PushBack(parsedItem)

						// 中间件需要排序在链表中服务函数之前，并且多个中间件按照顺序添加以便于后续执行
						case gHANDLER_TYPE_MIDDLEWARE:
							if lastMiddlewareElem == nil {
								lastMiddlewareElem = parsedItemList.PushFront(parsedItem)
							} else {
								lastMiddlewareElem = parsedItemList.InsertAfter(lastMiddlewareElem, parsedItem)
							}

						// 钩子函数存在性判断
						case gHANDLER_TYPE_HOOK:
							hasHook = true
							parsedItemList.PushBack(parsedItem)

						default:
							panic(fmt.Sprintf(`invalid handler type %d`, item.itemType))
						}
					}
				}
			}
		}
	}
	if parsedItemList.Len() > 0 {
		index := 0
		parsedItems = make([]*handlerParsedItem, parsedItemList.Len())
		for e := parsedItemList.Front(); e != nil; e = e.Next() {
			parsedItems[index] = e.Value.(*handlerParsedItem)
			index++
		}
	}
	return
}

// MarshalJSON implements the interface MarshalJSON for json.Marshal.
func (item *handlerItem) MarshalJSON() ([]byte, error) {
	switch item.itemType {
	case gHANDLER_TYPE_HOOK:
		return json.Marshal(
			fmt.Sprintf(
				`%s %s:%s (%s)`,
				item.router.Uri,
				item.router.Domain,
				item.router.Method,
				item.hookName,
			),
		)
	case gHANDLER_TYPE_MIDDLEWARE:
		return json.Marshal(
			fmt.Sprintf(
				`%s %s:%s (MIDDLEWARE)`,
				item.router.Uri,
				item.router.Domain,
				item.router.Method,
			),
		)
	default:
		return json.Marshal(
			fmt.Sprintf(
				`%s %s:%s`,
				item.router.Uri,
				item.router.Domain,
				item.router.Method,
			),
		)
	}
}

// MarshalJSON implements the interface MarshalJSON for json.Marshal.
func (item *handlerParsedItem) MarshalJSON() ([]byte, error) {
	return json.Marshal(item.handler)
}

// 生成回调方法查询的Key
func (s *Server) serveHandlerKey(method, path, domain string) string {
	if len(domain) > 0 {
		domain = "@" + domain
	}
	if method == "" {
		return path + strings.ToLower(domain)
	}
	return strings.ToUpper(method) + ":" + path + strings.ToLower(domain)
}
