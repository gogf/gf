// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gogf/gf/v2/container/glist"
	"github.com/gogf/gf/v2/encoding/gurl"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/internal/intlog"
	"github.com/gogf/gf/v2/internal/json"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/util/gmeta"
)

// handlerCacheItem is an item just for internal router searching cache.
type handlerCacheItem struct {
	parsedItems []*HandlerItemParsed
	serveItem   *HandlerItemParsed
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

// getHandlersWithCache searches the router item with cache feature for a given request.
func (s *Server) getHandlersWithCache(r *Request) (parsedItems []*HandlerItemParsed, serveItem *HandlerItemParsed, hasHook, hasServe bool) {
	var (
		ctx    = r.Context()
		method = r.Method
		path   = r.URL.Path
		host   = r.GetHost()
	)
	// In case of, eg:
	// Case 1:
	// 		GET /net/http
	// 		r.URL.Path    : /net/http
	// 		r.URL.RawPath : (empty string)
	// Case 2:
	// 		GET /net%2Fhttp
	// 		r.URL.Path    : /net/http
	// 		r.URL.RawPath : /net%2Fhttp
	if r.URL.RawPath != "" {
		path = r.URL.RawPath
	}
	// Special http method OPTIONS handling.
	// It searches the handler with the request method instead of OPTIONS method.
	if method == http.MethodOptions {
		if v := r.Request.Header.Get("Access-Control-Request-Method"); v != "" {
			method = v
		}
	}
	// Search and cache the router handlers.
	if xUrlPath := r.Header.Get(HeaderXUrlPath); xUrlPath != "" {
		path = xUrlPath
	}
	var handlerCacheKey = s.serveHandlerKey(method, path, host)
	value, err := s.serveCache.GetOrSetFunc(ctx, handlerCacheKey, func(ctx context.Context) (interface{}, error) {
		parsedItems, serveItem, hasHook, hasServe = s.searchHandlers(method, path, host)
		if parsedItems != nil {
			return &handlerCacheItem{parsedItems, serveItem, hasHook, hasServe}, nil
		}
		return nil, nil
	}, routeCacheDuration)
	if err != nil {
		intlog.Errorf(ctx, `%+v`, err)
	}
	if value != nil {
		item := value.Val().(*handlerCacheItem)
		return item.parsedItems, item.serveItem, item.hasHook, item.hasServe
	}
	return
}

// searchHandlers retrieve and returns the routers with given parameters.
// Note that the returned routers contain serving handler, middleware handlers and hook handlers.
func (s *Server) searchHandlers(method, path, domain string) (parsedItems []*HandlerItemParsed, serveItem *HandlerItemParsed, hasHook, hasServe bool) {
	if len(path) == 0 {
		return nil, nil, false, false
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

	// The default domain has the most priority when iteration.
	// Please see doSetHandler if you want to get known about the structure of serveTree.
	for _, domainItem := range []string{DefaultDomainName, domain} {
		p, ok := s.serveTree[domainItem]
		if !ok {
			continue
		}
		// Make a list array with a capacity of 16.
		lists := make([]*glist.List, 0, 16)
		for i, part := range array {
			// Add all lists of each node to the list array.
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
				item := e.Value.(*HandlerItem)
				// Filter repeated handler items, especially the middleware and hook handlers.
				// It is necessary, do not remove this checks logic unless you really know how it is necessary.
				//
				// The `repeatHandlerCheckMap` is used for repeat handler filtering during handler searching.
				// As there are fuzzy nodes, and the fuzzy nodes have both sub-nodes and sub-list nodes, there
				// may be repeated handler items in both sub-nodes and sub-list nodes. It here uses handler item id to
				// identify the same handler item that registered.
				//
				// The same handler item is the one that is registered in the same function doSetHandler.
				// Note that, one handler function(middleware or hook function) may be registered multiple times as
				// different handler items using function doSetHandler, and they have different handler item id.
				//
				// Note that twice, the handler function may be registered multiple times as different handler items.
				if _, isRepeatedHandler := repeatHandlerCheckMap[item.Id]; isRepeatedHandler {
					continue
				} else {
					repeatHandlerCheckMap[item.Id] = struct{}{}
				}
				// Serving handler can only be added to the handler array just once.
				// The first route item in the list has the most priority than the rest.
				// This ignoring can implement route overwritten feature.
				if hasServe {
					switch item.Type {
					case HandlerTypeHandler, HandlerTypeObject:
						continue
					}
				}
				if item.Router.Method == defaultMethod || item.Router.Method == method {
					// Note the rule having no fuzzy rules: len(match) == 1
					if match, err := gregex.MatchString(item.Router.RegRule, path); err == nil && len(match) > 0 {
						parsedItem := &HandlerItemParsed{item, nil}
						// If the rule contains fuzzy names,
						// it needs paring the URL to retrieve the values for the names.
						if len(item.Router.RegNames) > 0 {
							if len(match) > len(item.Router.RegNames) {
								parsedItem.Values = make(map[string]string)
								// It there repeated names, it just overwrites the same one.
								for i, name := range item.Router.RegNames {
									parsedItem.Values[name], _ = gurl.Decode(match[i+1])
								}
							}
						}
						switch item.Type {
						// The serving handler can be added just once.
						case HandlerTypeHandler, HandlerTypeObject:
							hasServe = true
							serveItem = parsedItem
							parsedItemList.PushBack(parsedItem)

						// The middleware is inserted before the serving handler.
						// If there are multiple middleware, they're inserted into the result list by their registering order.
						// The middleware is also executed by their registered order.
						case HandlerTypeMiddleware:
							if lastMiddlewareElem == nil {
								lastMiddlewareElem = parsedItemList.PushFront(parsedItem)
							} else {
								lastMiddlewareElem = parsedItemList.InsertAfter(lastMiddlewareElem, parsedItem)
							}

						// HOOK handler, just push it back to the list.
						case HandlerTypeHook:
							hasHook = true
							parsedItemList.PushBack(parsedItem)

						default:
							panic(gerror.Newf(`invalid handler type %s`, item.Type))
						}
					}
				}
			}
		}
	}
	if parsedItemList.Len() > 0 {
		var index = 0
		parsedItems = make([]*HandlerItemParsed, parsedItemList.Len())
		for e := parsedItemList.Front(); e != nil; e = e.Next() {
			parsedItems[index] = e.Value.(*HandlerItemParsed)
			index++
		}
	}
	return
}

// MarshalJSON implements the interface MarshalJSON for json.Marshal.
func (h *HandlerItem) MarshalJSON() ([]byte, error) {
	switch h.Type {
	case HandlerTypeHook:
		return json.Marshal(
			fmt.Sprintf(
				`%s %s:%s (%s)`,
				h.Router.Uri,
				h.Router.Domain,
				h.Router.Method,
				h.HookName,
			),
		)
	case HandlerTypeMiddleware:
		return json.Marshal(
			fmt.Sprintf(
				`%s %s:%s (MIDDLEWARE)`,
				h.Router.Uri,
				h.Router.Domain,
				h.Router.Method,
			),
		)
	default:
		return json.Marshal(
			fmt.Sprintf(
				`%s %s:%s`,
				h.Router.Uri,
				h.Router.Domain,
				h.Router.Method,
			),
		)
	}
}

// MarshalJSON implements the interface MarshalJSON for json.Marshal.
func (h *HandlerItemParsed) MarshalJSON() ([]byte, error) {
	return json.Marshal(h.Handler)
}

// GetMetaTag retrieves and returns the metadata value associated with the given key from the request struct.
// The meta value is from struct tags from g.Meta/gmeta.Meta type.
func (h *HandlerItem) GetMetaTag(key string) string {
	if h == nil {
		return ""
	}
	metaValue := gmeta.Get(h.Info.Type.In(1), key)
	if metaValue != nil {
		return metaValue.String()
	}
	return ""
}

// GetMetaTag retrieves and returns the metadata value associated with the given key from the request struct.
// The meta value is from struct tags from g.Meta/gmeta.Meta type.
// For example:
//
//	type GetMetaTagReq struct {
//	    g.Meta `path:"/test" method:"post" summary:"meta_tag" tags:"meta"`
//	    // ...
//	}
//
// r.GetServeHandler().GetMetaTag("summary") // returns "meta_tag"
// r.GetServeHandler().GetMetaTag("method")  // returns "post"
func (h *HandlerItemParsed) GetMetaTag(key string) string {
	if h == nil || h.Handler == nil {
		return ""
	}
	return h.Handler.GetMetaTag(key)
}
