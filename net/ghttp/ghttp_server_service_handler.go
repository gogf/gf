// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import (
	"bytes"
	"github.com/gogf/gf/debug/gdebug"
	"strings"

	"github.com/gogf/gf/text/gstr"
)

// 注意该方法是直接绑定函数的内存地址，执行的时候直接执行该方法，不会存在初始化新的控制器逻辑
func (s *Server) BindHandler(pattern string, handler HandlerFunc) {
	s.doBindHandler(pattern, handler, nil, "")
}

// 绑定URI到操作函数/方法
// pattern的格式形如：/user/list, put:/user, delete:/user, post:/user@johng.cn
// 支持RESTful的请求格式，具体业务逻辑由绑定的处理方法来执行
func (s *Server) doBindHandler(
	pattern string, handler HandlerFunc,
	middleware []HandlerFunc, source string,
) {
	s.setHandler(pattern, &handlerItem{
		itemName:   gdebug.FuncPath(handler),
		itemType:   gHANDLER_TYPE_HANDLER,
		itemFunc:   handler,
		middleware: middleware,
		source:     source,
	})
}

// 通过映射map绑定URI到操作函数/方法
func (s *Server) bindHandlerByMap(m map[string]*handlerItem) {
	for p, h := range m {
		s.setHandler(p, h)
	}
}

// 将内置的名称按照设定的规则合并到pattern中，内置名称按照{.xxx}规则命名。
// 规则1：pattern中的URI包含{.struct}关键字，则替换该关键字为结构体名称；
// 规则2：pattern中的URI包含{.method}关键字，则替换该关键字为方法名称；
// 规则2：如果不满足规则1，那么直接将防发明附加到pattern中的URI后面；
func (s *Server) mergeBuildInNameToPattern(pattern string, structName, methodName string, allowAppend bool) string {
	structName = s.nameToUri(structName)
	methodName = s.nameToUri(methodName)
	pattern = strings.Replace(pattern, "{.struct}", structName, -1)
	if strings.Index(pattern, "{.method}") != -1 {
		return strings.Replace(pattern, "{.method}", methodName, -1)
	}
	// 不允许将方法名称append到路由末尾
	if !allowAppend {
		return pattern
	}
	// 检测域名后缀
	array := strings.Split(pattern, "@")
	// 分离URI(其实可能包含HTTP Method)
	uri := array[0]
	uri = strings.TrimRight(uri, "/") + "/" + methodName
	// 加上指定域名后缀
	if len(array) > 1 {
		return uri + "@" + array[1]
	}
	return uri
}

// 将给定的名称转换为URL规范格式。
// 规则0: 全部转换为小写，方法名中间存在大写字母，转换为小写URI地址以“-”号链接每个单词；
// 规则1: 不处理名称，以原有名称构建成URI
// 规则2: 仅转为小写，单词间不使用连接符号
// 规则3: 采用驼峰命名方式
func (s *Server) nameToUri(name string) string {
	switch s.config.NameToUriType {
	case URI_TYPE_FULLNAME:
		return name

	case URI_TYPE_ALLLOWER:
		return strings.ToLower(name)

	case URI_TYPE_CAMEL:
		part := bytes.NewBuffer(nil)
		if gstr.IsLetterUpper(name[0]) {
			part.WriteByte(name[0] + 32)
		} else {
			part.WriteByte(name[0])
		}
		part.WriteString(name[1:])
		return part.String()

	case URI_TYPE_DEFAULT:
		fallthrough
	default:
		part := bytes.NewBuffer(nil)
		for i := 0; i < len(name); i++ {
			if i > 0 && gstr.IsLetterUpper(name[i]) {
				part.WriteByte('-')
			}
			if gstr.IsLetterUpper(name[i]) {
				part.WriteByte(name[i] + 32)
			} else {
				part.WriteByte(name[i])
			}
		}
		return part.String()
	}
}
