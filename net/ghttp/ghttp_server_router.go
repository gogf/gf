// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import (
	"errors"
	"fmt"
	"github.com/gogf/gf/container/gtype"
	"strings"

	"github.com/gogf/gf/debug/gdebug"

	"github.com/gogf/gf/container/glist"
	"github.com/gogf/gf/text/gregex"
	"github.com/gogf/gf/text/gstr"
)

const (
	gFILTER_KEY = "/net/ghttp/ghttp"
)

var (
	// handlerIdGenerator is handler item id generator.
	handlerIdGenerator = gtype.NewInt()
)

// handlerKey creates and returns an unique router key for given parameters.
func (s *Server) handlerKey(hook, method, path, domain string) string {
	return hook + "%" + s.serveHandlerKey(method, path, domain)
}

// parsePattern parses the given pattern to domain, method and path variable.
func (s *Server) parsePattern(pattern string) (domain, method, path string, err error) {
	path = strings.TrimSpace(pattern)
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
	if path != "/" {
		path = strings.TrimRight(path, "/")
	}
	return
}

// setHandler creates router item with given handler and pattern and registers the handler to the router tree.
// The router tree can be treated as a multilayer hash table, please refer to the comment in following codes.
// This function is called during server starts up, which cares little about the performance. What really cares
// is the well designed router storage structure for router searching when the request is under serving.
func (s *Server) setHandler(pattern string, handler *handlerItem) {
	handler.itemId = handlerIdGenerator.Add(1)
	domain, method, uri, err := s.parsePattern(pattern)
	if err != nil {
		s.Logger().Fatal("invalid pattern:", pattern, err)
		return
	}
	if len(uri) == 0 || uri[0] != '/' {
		s.Logger().Fatal("invalid pattern:", pattern, "URI should lead with '/'")
		return
	}

	// Repeated router checks, this feature can be disabled by server configuration.
	routerKey := s.handlerKey(handler.hookName, method, uri, domain)
	if !s.config.RouteOverWrite {
		switch handler.itemType {
		case gHANDLER_TYPE_HANDLER, gHANDLER_TYPE_OBJECT, gHANDLER_TYPE_CONTROLLER:
			if item, ok := s.routesMap[routerKey]; ok {
				s.Logger().Fatalf(`duplicated route registry "%s", already registered at %s`, pattern, item[0].file)
				return
			}
		}
	}
	// Create a new router by given parameter.
	handler.router = &Router{
		Uri:      uri,
		Domain:   domain,
		Method:   strings.ToUpper(method),
		Priority: strings.Count(uri[1:], "/"),
	}
	handler.router.RegRule, handler.router.RegNames = s.patternToRegRule(uri)

	if _, ok := s.serveTree[domain]; !ok {
		s.serveTree[domain] = make(map[string]interface{})
	}
	// List array, very important for router register.
	// There may be multiple lists adding into this array when searching from root to leaf.
	lists := make([]*glist.List, 0)
	array := ([]string)(nil)
	if strings.EqualFold("/", uri) {
		array = []string{"/"}
	} else {
		array = strings.Split(uri[1:], "/")
	}
	// Multilayer hash table:
	// 1. Each node of the table is separated by URI path which is split by char '/'.
	// 2. The key "*fuzz" specifies this node is a fuzzy node, which has no certain name.
	// 3. The key "*list" is the list item of the node, MOST OF THE NODES HAVE THIS ITEM,
	//    especially the fuzzy node. NOTE THAT the fuzzy node must have the "*list" item,
	//    and the leaf node also has "*list" item. If the node is not a fuzzy node either
	//    a leaf, it neither has "*list" item.
	// 2. The "*list" item is a list containing registered router items ordered by their
	//    priorities from high to low.
	// 3. There may be repeated router items in the router lists. The lists' priorities
	//    from root to leaf are from low to high.
	p := s.serveTree[domain]
	for i, part := range array {
		// Ignore empty URI part, like: /user//index
		if part == "" {
			continue
		}
		// Check if it's a fuzzy node.
		if gregex.IsMatchString(`^[:\*]|\{[\w\.\-]+\}|\*`, part) {
			part = "*fuzz"
			// If it's a fuzzy node, it creates a "*list" item - which is a list - in the hash map.
			// All the sub router items from this fuzzy node will also be added to its "*list" item.
			if v, ok := p.(map[string]interface{})["*list"]; !ok {
				newListForFuzzy := glist.New()
				p.(map[string]interface{})["*list"] = newListForFuzzy
				lists = append(lists, newListForFuzzy)
			} else {
				lists = append(lists, v.(*glist.List))
			}
		}
		// Make a new bucket for current node.
		if _, ok := p.(map[string]interface{})[part]; !ok {
			p.(map[string]interface{})[part] = make(map[string]interface{})
		}
		// Loop to next bucket.
		p = p.(map[string]interface{})[part]
		// The leaf is a hash map and must have an item named "*list", which contains the router item.
		// The leaf can be furthermore extended by adding more ket-value pairs into its map.
		// Note that the `v != "*fuzz"` comparison is required as the list might be added in the former
		// fuzzy checks.
		if i == len(array)-1 && part != "*fuzz" {
			if v, ok := p.(map[string]interface{})["*list"]; !ok {
				leafList := glist.New()
				p.(map[string]interface{})["*list"] = leafList
				lists = append(lists, leafList)
			} else {
				lists = append(lists, v.(*glist.List))
			}
		}
	}
	// It iterates the list array of <lists>, compares priorities and inserts the new router item in
	// the proper position of each list. The priority of the list is ordered from high to low.
	item := (*handlerItem)(nil)
	for _, l := range lists {
		pushed := false
		for e := l.Front(); e != nil; e = e.Next() {
			item = e.Value.(*handlerItem)
			// Checks the priority whether inserting the route item before current item,
			// which means it has more higher priority.
			if s.compareRouterPriority(handler, item) {
				l.InsertBefore(e, handler)
				pushed = true
				goto end
			}
		}
	end:
		// Just push back in default.
		if !pushed {
			l.PushBack(handler)
		}
	}
	// Initialize the route map item.
	if _, ok := s.routesMap[routerKey]; !ok {
		s.routesMap[routerKey] = make([]registeredRouteItem, 0)
	}
	_, file, line := gdebug.CallerWithFilter(gFILTER_KEY)
	routeItem := registeredRouteItem{
		file:    fmt.Sprintf(`%s:%d`, file, line),
		handler: handler,
	}
	switch handler.itemType {
	case gHANDLER_TYPE_HANDLER, gHANDLER_TYPE_OBJECT, gHANDLER_TYPE_CONTROLLER:
		// Overwrite the route.
		s.routesMap[routerKey] = []registeredRouteItem{routeItem}
	default:
		// Append the route.
		s.routesMap[routerKey] = append(s.routesMap[routerKey], routeItem)
	}
	//gutil.Dump(s.serveTree)
}

// 对比两个handlerItem的优先级，需要非常注意的是，注意新老对比项的参数先后顺序。
// 返回值true表示newItem优先级比oldItem高，会被添加链表中oldRouter的前面；否则后面。
// 优先级比较规则：
// 1、中间件优先级最高，按照添加顺序优先级执行；
// 2、其他路由注册类型，层级越深优先级越高(对比/数量)；
// 3、模糊规则优先级：{xxx} > :xxx > *xxx；
func (s *Server) compareRouterPriority(newItem *handlerItem, oldItem *handlerItem) bool {
	// 中间件优先级最高，按照添加顺序优先级执行
	if newItem.itemType == gHANDLER_TYPE_MIDDLEWARE && oldItem.itemType == gHANDLER_TYPE_MIDDLEWARE {
		return false
	}
	if newItem.itemType == gHANDLER_TYPE_MIDDLEWARE && oldItem.itemType != gHANDLER_TYPE_MIDDLEWARE {
		return true
	}
	// 优先比较层级，层级越深优先级越高
	if newItem.router.Priority > oldItem.router.Priority {
		return true
	}
	if newItem.router.Priority < oldItem.router.Priority {
		return false
	}
	// 精准匹配比模糊匹配规则优先级高，例如：/name/act 比 /{name}/:act 优先级高
	var fuzzyCountFieldNew, fuzzyCountFieldOld int
	var fuzzyCountNameNew, fuzzyCountNameOld int
	var fuzzyCountAnyNew, fuzzyCountAnyOld int
	var fuzzyCountTotalNew, fuzzyCountTotalOld int
	for _, v := range newItem.router.Uri {
		switch v {
		case '{':
			fuzzyCountFieldNew++
		case ':':
			fuzzyCountNameNew++
		case '*':
			fuzzyCountAnyNew++
		}
	}
	for _, v := range oldItem.router.Uri {
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

	/** 比较路由规则长度，越长的规则优先级越高，模糊/命名规则不算长度 **/

	// 例如：/admin-goods-{page} 比 /admin-{page} 优先级高
	var uriNew, uriOld string
	uriNew, _ = gregex.ReplaceString(`\{[^/]+\}`, "", newItem.router.Uri)
	uriNew, _ = gregex.ReplaceString(`:[^/]+`, "", uriNew)
	uriNew, _ = gregex.ReplaceString(`\*[^/]+`, "", uriNew)
	uriOld, _ = gregex.ReplaceString(`\{[^/]+\}`, "", oldItem.router.Uri)
	uriOld, _ = gregex.ReplaceString(`:[^/]+`, "", uriOld)
	uriOld, _ = gregex.ReplaceString(`\*[^/]+`, "", uriOld)
	if len(uriNew) > len(uriOld) {
		return true
	}
	if len(uriNew) < len(uriOld) {
		return false
	}

	/* 模糊规则数量相等，后续不用再判断*规则的数量比较了 */

	// 比较HTTP METHOD，更精准的优先级更高
	if newItem.router.Method != gDEFAULT_METHOD {
		return true
	}
	if oldItem.router.Method != gDEFAULT_METHOD {
		return true
	}

	// 如果是服务路由，那么新的规则比旧的规则优先级高(路由覆盖)
	if newItem.itemType == gHANDLER_TYPE_HANDLER ||
		newItem.itemType == gHANDLER_TYPE_OBJECT ||
		newItem.itemType == gHANDLER_TYPE_CONTROLLER {
		return true
	}

	// 如果是其他路由(HOOK/中间件)，那么新的规则比旧的规则优先级低，使得注册相同路由则顺序执行
	return false
}

// 将pattern（不带method和domain）解析成正则表达式匹配以及对应的query字符串
func (s *Server) patternToRegRule(rule string) (regrule string, names []string) {
	if len(rule) < 2 {
		return rule, nil
	}
	regrule = "^"
	array := strings.Split(rule[1:], "/")
	for _, v := range array {
		if len(v) == 0 {
			continue
		}
		switch v[0] {
		case ':':
			if len(v) > 1 {
				regrule += `/([^/]+)`
				names = append(names, v[1:])
			} else {
				regrule += `/[^/]+`
			}
		case '*':
			if len(v) > 1 {
				regrule += `/{0,1}(.*)`
				names = append(names, v[1:])
			} else {
				regrule += `/{0,1}.*`
			}
		default:
			// Special chars replacement.
			v = gstr.ReplaceByMap(v, map[string]string{
				`.`: `\.`,
				`+`: `\+`,
				`*`: `.*`,
			})
			s, _ := gregex.ReplaceStringFunc(`\{[\w\.\-]+\}`, v, func(s string) string {
				names = append(names, s[1:len(s)-1])
				return `([^/]+)`
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
