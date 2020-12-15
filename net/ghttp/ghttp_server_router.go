// Copyright GoFrame Author(https://github.com/gogf/gf). All Rights Reserved.
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

// routerMapKey creates and returns an unique router key for given parameters.
// This key is used for Server.routerMap attribute, which is mainly for checks for
// repeated router registering.
func (s *Server) routerMapKey(hook, method, path, domain string) string {
	return hook + "%" + s.serveHandlerKey(method, path, domain)
}

// parsePattern parses the given pattern to domain, method and path variable.
func (s *Server) parsePattern(pattern string) (domain, method, path string, err error) {
	path = strings.TrimSpace(pattern)
	domain = defaultDomainName
	method = defaultMethod
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
	if handler.source == "" {
		_, file, line := gdebug.CallerWithFilter(gFILTER_KEY)
		handler.source = fmt.Sprintf(`%s:%d`, file, line)
	}
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
	routerKey := s.routerMapKey(handler.hookName, method, uri, domain)
	if !s.config.RouteOverWrite {
		switch handler.itemType {
		case handlerTypeHandler, handlerTypeObject, handlerTypeController:
			if item, ok := s.routesMap[routerKey]; ok {
				s.Logger().Fatalf(
					`duplicated route registry "%s" at %s , already registered at %s`,
					pattern, handler.source, item[0].source,
				)
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
	handler.router.RegRule, handler.router.RegNames = s.patternToRegular(uri)

	if _, ok := s.serveTree[domain]; !ok {
		s.serveTree[domain] = make(map[string]interface{})
	}
	// List array, very important for router registering.
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

	routeItem := registeredRouteItem{
		source:  handler.source,
		handler: handler,
	}
	switch handler.itemType {
	case handlerTypeHandler, handlerTypeObject, handlerTypeController:
		// Overwrite the route.
		s.routesMap[routerKey] = []registeredRouteItem{routeItem}
	default:
		// Append the route.
		s.routesMap[routerKey] = append(s.routesMap[routerKey], routeItem)
	}
}

// compareRouterPriority compares the priority between <newItem> and <oldItem>. It returns true
// if <newItem>'s priority is higher than <oldItem>, else it returns false. The higher priority
// item will be insert into the router list before the other one.
//
// Comparison rules:
// 1. The middleware has the most high priority.
// 2. URI: The deeper the higher (simply check the count of char '/' in the URI).
// 3. Route type: {xxx} > :xxx > *xxx.
func (s *Server) compareRouterPriority(newItem *handlerItem, oldItem *handlerItem) bool {
	// If they're all type of middleware, the priority is according their registered sequence.
	if newItem.itemType == handlerTypeMiddleware && oldItem.itemType == handlerTypeMiddleware {
		return false
	}
	// The middleware has the most high priority.
	if newItem.itemType == handlerTypeMiddleware && oldItem.itemType != handlerTypeMiddleware {
		return true
	}
	// URI: The deeper the higher (simply check the count of char '/' in the URI).
	if newItem.router.Priority > oldItem.router.Priority {
		return true
	}
	if newItem.router.Priority < oldItem.router.Priority {
		return false
	}

	// Compare the length of their URI,
	// but the fuzzy and named parts of the URI are not calculated to the result.

	// Eg:
	// /admin-goods-{page} > /admin-{page}
	// /{hash}.{type}      > /{hash}
	var uriNew, uriOld string
	uriNew, _ = gregex.ReplaceString(`\{[^/]+?\}`, "", newItem.router.Uri)
	uriOld, _ = gregex.ReplaceString(`\{[^/]+?\}`, "", oldItem.router.Uri)
	uriNew, _ = gregex.ReplaceString(`:[^/]+?`, "", uriNew)
	uriOld, _ = gregex.ReplaceString(`:[^/]+?`, "", uriOld)
	uriNew, _ = gregex.ReplaceString(`\*[^/]*`, "", uriNew) // Replace "/*" and "/*any".
	uriOld, _ = gregex.ReplaceString(`\*[^/]*`, "", uriOld) // Replace "/*" and "/*any".
	if len(uriNew) > len(uriOld) {
		return true
	}
	if len(uriNew) < len(uriOld) {
		return false
	}

	// Route type checks: {xxx} > :xxx > *xxx.
	// Eg:
	// /name/act > /{name}/:act
	var (
		fuzzyCountFieldNew int
		fuzzyCountFieldOld int
		fuzzyCountNameNew  int
		fuzzyCountNameOld  int
		fuzzyCountAnyNew   int
		fuzzyCountAnyOld   int
		fuzzyCountTotalNew int
		fuzzyCountTotalOld int
	)
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

	// If the counts of their fuzzy rules equal.

	// Eg: /name/{act} > /name/:act
	if fuzzyCountFieldNew > fuzzyCountFieldOld {
		return true
	}
	if fuzzyCountFieldNew < fuzzyCountFieldOld {
		return false
	}
	// Eg: /name/:act > /name/*act
	if fuzzyCountNameNew > fuzzyCountNameOld {
		return true
	}
	if fuzzyCountNameNew < fuzzyCountNameOld {
		return false
	}

	// It then compares the accuracy of their http method,
	// the more accurate the more priority.
	if newItem.router.Method != defaultMethod {
		return true
	}
	if oldItem.router.Method != defaultMethod {
		return true
	}

	// If they have different router type,
	// the new router item has more priority than the other one.
	if newItem.itemType == handlerTypeHandler ||
		newItem.itemType == handlerTypeObject ||
		newItem.itemType == handlerTypeController {
		return true
	}

	// Other situations, like HOOK items,
	// the old router item has more priority than the other one.
	return false
}

// patternToRegular converts route rule to according regular expression.
func (s *Server) patternToRegular(rule string) (regular string, names []string) {
	if len(rule) < 2 {
		return rule, nil
	}
	regular = "^"
	array := strings.Split(rule[1:], "/")
	for _, v := range array {
		if len(v) == 0 {
			continue
		}
		switch v[0] {
		case ':':
			if len(v) > 1 {
				regular += `/([^/]+)`
				names = append(names, v[1:])
			} else {
				regular += `/[^/]+`
			}
		case '*':
			if len(v) > 1 {
				regular += `/{0,1}(.*)`
				names = append(names, v[1:])
			} else {
				regular += `/{0,1}.*`
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
				regular += "/" + v
			} else {
				regular += "/" + s
			}
		}
	}
	regular += `$`
	return
}
