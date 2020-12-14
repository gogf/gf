// Copyright GoFrame Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import (
	"fmt"
	"github.com/gogf/gf/internal/json"
	"strings"

	"github.com/gogf/gf/container/glist"
	"github.com/gogf/gf/text/gregex"
)

// handlerCacheItem is an item just for internal router searching cache.
type handlerCacheItem struct {
	parsedItems []*handlerParsedItem
	hasHook     bool
	hasServe    bool
}

// serveHandlerKey creates and returns a handler key for router.
func (s *Server) serveHandlerKey(method, path, domain string) string {
	if len(domain) > 0 {
		domain = "@" + domain
	}
	if method == "" {
		return path + strings.ToLower(domain)
	}
	return strings.ToUpper(method) + ":" + path + strings.ToLower(domain)
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
	value, _ := s.serveCache.GetOrSetFunc(
		s.serveHandlerKey(method, r.URL.Path, r.GetHost()),
		func() (interface{}, error) {
			parsedItems, hasHook, hasServe = s.searchHandlers(method, r.URL.Path, r.GetHost())
			if parsedItems != nil {
				return &handlerCacheItem{parsedItems, hasHook, hasServe}, nil
			}
			return nil, nil
		}, routeCacheDuration)
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
	// In case of double '/' URI, for example:
	// /user//index, //user/index, //user//index//
	var previousIsSep = false
	for i := 0; i < len(path); {
		if path[i] == '/' {
			if previousIsSep {
				path = path[:i] + path[i+1:]
				continue
			} else {
				previousIsSep = true
			}
		} else {
			previousIsSep = false
		}
		i++
	}
	// Split the URL.path to separate parts.
	var array []string
	if strings.EqualFold("/", path) {
		array = []string{"/"}
	} else {
		array = strings.Split(path[1:], "/")
	}
	var (
		lastMiddlewareElem    *glist.Element
		parsedItemList        = glist.New()
		repeatHandlerCheckMap = make(map[int]struct{}, 16)
	)

	// Default domain has the most priority when iteration.
	for _, domain := range []string{defaultDomainName, domain} {
		p, ok := s.serveTree[domain]
		if !ok {
			continue
		}
		// Make a list array with capacity of 16.
		lists := make([]*glist.List, 0, 16)
		for i, part := range array {
			// Add all list of each node to the list array.
			if v, ok := p.(map[string]interface{})["*list"]; ok {
				lists = append(lists, v.(*glist.List))
			}
			if v, ok := p.(map[string]interface{})[part]; ok {
				// Loop to the next node by certain key name.
				p = v
				if i == len(array)-1 {
					if v, ok := p.(map[string]interface{})["*list"]; ok {
						lists = append(lists, v.(*glist.List))
						break
					}
				}
			} else if v, ok := p.(map[string]interface{})["*fuzz"]; ok {
				// Loop to the next node by fuzzy node item.
				p = v
			}
			if i == len(array)-1 {
				// It here also checks the fuzzy item,
				// for rule case like: "/user/*action" matches to "/user".
				if v, ok := p.(map[string]interface{})["*fuzz"]; ok {
					p = v
				}
				// The leaf must have a list item. It adds the list to the list array.
				if v, ok := p.(map[string]interface{})["*list"]; ok {
					lists = append(lists, v.(*glist.List))
				}
			}
		}

		// OK, let's loop the result list array, adding the handler item to the result handler result array.
		// As the tail of the list array has the most priority, it iterates the list array from its tail to head.
		for i := len(lists) - 1; i >= 0; i-- {
			for e := lists[i].Front(); e != nil; e = e.Next() {
				item := e.Value.(*handlerItem)
				// Filter repeated handler item, especially the middleware and hook handlers.
				// It is necessary, do not remove this checks logic unless you really know how it is necessary.
				if _, ok := repeatHandlerCheckMap[item.itemId]; ok {
					continue
				} else {
					repeatHandlerCheckMap[item.itemId] = struct{}{}
				}
				// Serving handler can only be added to the handler array just once.
				if hasServe {
					switch item.itemType {
					case handlerTypeHandler, handlerTypeObject, handlerTypeController:
						continue
					}
				}
				if item.router.Method == defaultMethod || item.router.Method == method {
					// Note the rule having no fuzzy rules: len(match) == 1
					if match, err := gregex.MatchString(item.router.RegRule, path); err == nil && len(match) > 0 {
						parsedItem := &handlerParsedItem{item, nil}
						// If the rule contains fuzzy names,
						// it needs paring the URL to retrieve the values for the names.
						if len(item.router.RegNames) > 0 {
							if len(match) > len(item.router.RegNames) {
								parsedItem.values = make(map[string]string)
								// It there repeated names, it just overwrites the same one.
								for i, name := range item.router.RegNames {
									parsedItem.values[name] = match[i+1]
								}
							}
						}
						switch item.itemType {
						// The serving handler can be only added just once.
						case handlerTypeHandler, handlerTypeObject, handlerTypeController:
							hasServe = true
							parsedItemList.PushBack(parsedItem)

						// The middleware is inserted before the serving handler.
						// If there're multiple middleware, they're inserted into the result list by their registering order.
						// The middleware are also executed by their registered order.
						case handlerTypeMiddleware:
							if lastMiddlewareElem == nil {
								lastMiddlewareElem = parsedItemList.PushFront(parsedItem)
							} else {
								lastMiddlewareElem = parsedItemList.InsertAfter(lastMiddlewareElem, parsedItem)
							}

						// HOOK handler, just push it back to the list.
						case handlerTypeHook:
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
	case handlerTypeHook:
		return json.Marshal(
			fmt.Sprintf(
				`%s %s:%s (%s)`,
				item.router.Uri,
				item.router.Domain,
				item.router.Method,
				item.hookName,
			),
		)
	case handlerTypeMiddleware:
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
