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
	// 用于服务函数的ID生成变量
	handlerIdGenerator = gtype.NewInt()
)

// 解析pattern
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
	// 去掉末尾的"/"符号，与路由匹配时处理一致
	if path != "/" {
		path = strings.TrimRight(path, "/")
	}
	return
}

// 路由注册处理方法。
// 非叶节点为哈希表检索节点，按照URI注册的层级进行高效检索，直至到叶子链表节点；
// 叶子节点是链表，按照优先级进行排序，优先级高的排前面，按照遍历检索，按照哈希表层级检索后的叶子链表数据量不会很大，所以效率比较高；
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
	// 注册地址记录及重复注册判断
	regKey := s.handlerKey(handler.hookName, method, uri, domain)
	if !s.config.RouteOverWrite {
		switch handler.itemType {
		case gHANDLER_TYPE_HANDLER, gHANDLER_TYPE_OBJECT, gHANDLER_TYPE_CONTROLLER:
			if item, ok := s.routesMap[regKey]; ok {
				s.Logger().Fatalf(`duplicated route registry "%s", already registered at %s`, pattern, item[0].file)
				return
			}
		}
	}
	// 注册的路由信息对象
	handler.router = &Router{
		Uri:      uri,
		Domain:   domain,
		Method:   method,
		Priority: strings.Count(uri[1:], "/"),
	}
	handler.router.RegRule, handler.router.RegNames = s.patternToRegRule(uri)

	if _, ok := s.serveTree[domain]; !ok {
		s.serveTree[domain] = make(map[string]interface{})
	}
	// 当前节点的规则链表
	lists := make([]*glist.List, 0)
	array := ([]string)(nil)
	if strings.EqualFold("/", uri) {
		array = []string{"/"}
	} else {
		array = strings.Split(uri[1:], "/")
	}
	// 键名"*fuzz"代表当前节点为模糊匹配节点，该节点也会有一个*list链表；
	// 键名"*list"代表链表，叶子节点和模糊匹配节点都有该属性，优先级越高越排前；
	p := s.serveTree[domain]
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
				p.(map[string]interface{})["*list"] = glist.New()
				lists = append(lists, p.(map[string]interface{})["*list"].(*glist.List))
			} else {
				lists = append(lists, v.(*glist.List))
			}
		}
		// 属性层级数据写入
		if _, ok := p.(map[string]interface{})[v]; !ok {
			p.(map[string]interface{})[v] = make(map[string]interface{})
		}
		p = p.(map[string]interface{})[v]
		// 到达叶子节点，往list中增加匹配规则(条件 v != "*fuzz" 是因为模糊节点的话在前面已经添加了*list链表)
		if k == len(array)-1 && v != "*fuzz" {
			if v, ok := p.(map[string]interface{})["*list"]; !ok {
				p.(map[string]interface{})["*list"] = glist.New()
				lists = append(lists, p.(map[string]interface{})["*list"].(*glist.List))
			} else {
				lists = append(lists, v.(*glist.List))
			}
		}
	}

	// 上面循环后得到的lists是该路由规则一路匹配下来相关的模糊匹配链表(注意不是这棵树所有的链表)。
	// 下面从头开始遍历每个节点的模糊匹配链表，将该路由项插入进去(按照优先级高的放在lists链表的前面)
	item := (*handlerItem)(nil)
	for _, l := range lists {
		pushed := false
		for e := l.Front(); e != nil; e = e.Next() {
			item = e.Value.(*handlerItem)
			// 判断优先级，决定当前注册项的插入顺序
			if s.compareRouterPriority(handler, item) {
				l.InsertBefore(e, handler)
				pushed = true
				goto end
			}
		}
	end:
		if !pushed {
			l.PushBack(handler)
		}
	}
	//gutil.Dump(s.serveTree)
	if _, ok := s.routesMap[regKey]; !ok {
		s.routesMap[regKey] = make([]registeredRouteItem, 0)
	}
	_, file, line := gdebug.CallerWithFilter(gFILTER_KEY)
	s.routesMap[regKey] = append(s.routesMap[regKey], registeredRouteItem{
		file:    fmt.Sprintf(`%s:%d`, file, line),
		handler: handler,
	})
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
			// 特殊字符替换
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
