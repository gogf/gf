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

// handlerCacheItem is an item for router searching cache.
type handlerCacheItem struct {
	parsedItems []*handlerParsedItem
	hasHook     bool
	hasServe    bool
}

// getHandlersWithCache searches the router item with cache feature for given request.
func (s *Server) getHandlersWithCache(r *Request) (parsedItems []*handlerParsedItem, hasHook, hasServe bool) {
	method := r.Method
	// Special http method OPTIONS handling.
	// It searches the handler with the request method instead of OPTIONS method.
	if method == "OPTIONS" {
		if v := r.Request.Header.Get("Access-Control-Request-Method"); v != "" {
			method = v
		}
	}
	// Search and cache the router handlers.
	value := s.serveCache.GetOrSetFunc(s.serveHandlerKey(method, r.URL.Path, r.GetHost()), func() interface{} {
		parsedItems, hasHook, hasServe = s.searchHandlers(method, r.URL.Path, r.GetHost())
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

// searchHandlers retrieves and returns the routers with given parameters.
// Note that the returned routers contain serving handler, middleware handlers and hook handlers.
func (s *Server) searchHandlers(method, path, domain string) (parsedItems []*handlerParsedItem, hasHook, hasServe bool) {
	if len(path) == 0 {
		return nil, false, false
	}
	// Split the URL.path to separate parts.
	var array []string
	if strings.EqualFold("/", path) {
		array = []string{"/"}
	} else {
		array = strings.Split(path[1:], "/")
	}
	parsedItemList := glist.New()
	lastMiddlewareElem := (*glist.Element)(nil)
	repeatHandlerCheckMap := make(map[int]struct{}, 16)
	// Default domain has the most priority when iteration.
	for _, domain := range []string{gDEFAULT_DOMAIN, domain} {
		p, ok := s.serveTree[domain]
		if !ok {
			continue
		}
		// Make a list array with capacity of 16.
		lists := make([]*glist.List, 0, 16)
		for i, part := range array {
			// In case of double '/' URI, eg: /user//index
			if part == "" {
				continue
			}
			if v, ok := p.(map[string]interface{})["*list"]; ok {
				lists = append(lists, v.(*glist.List))
			}
			if v, ok := p.(map[string]interface{})[part]; ok {
				p = v
				if i == len(array)-1 {
					if v, ok := p.(map[string]interface{})["*list"]; ok {
						lists = append(lists, v.(*glist.List))
						break
					}
				}
			} else {
				if v, ok := p.(map[string]interface{})["*fuzz"]; ok {
					p = v
				}
			}
			// 如果是叶子节点，同时判断当前层级的"*fuzz"键名，解决例如：/user/*action 匹配 /user 的规则
			if i == len(array)-1 {
				if v, ok := p.(map[string]interface{})["*fuzz"]; ok {
					p = v
				}
				if v, ok := p.(map[string]interface{})["*list"]; ok {
					lists = append(lists, v.(*glist.List))
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
				if item.router.Method == gDEFAULT_METHOD || item.router.Method == method {
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

// serveHandlerKey creates and returns a cache key for router.
func (s *Server) serveHandlerKey(method, path, domain string) string {
	if len(domain) > 0 {
		domain = "@" + domain
	}
	if method == "" {
		return path + strings.ToLower(domain)
	}
	return strings.ToUpper(method) + ":" + path + strings.ToLower(domain)
}
