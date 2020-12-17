// Copyright GoFrame Author(https://github.com/gogf/gf). All Rights Reserved.
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

// BindHandler registers a handler function to server with given pattern.
func (s *Server) BindHandler(pattern string, handler HandlerFunc) {
	s.doBindHandler(pattern, handler, nil, "")
}

// doBindHandler registers a handler function to server with given pattern.
// The parameter <pattern> is like:
// /user/list, put:/user, delete:/user, post:/user@goframe.org
func (s *Server) doBindHandler(
	pattern string, handler HandlerFunc,
	middleware []HandlerFunc, source string,
) {
	s.setHandler(pattern, &handlerItem{
		itemName:   gdebug.FuncPath(handler),
		itemType:   handlerTypeHandler,
		itemFunc:   handler,
		middleware: middleware,
		source:     source,
	})
}

// bindHandlerByMap registers handlers to server using map.
func (s *Server) bindHandlerByMap(m map[string]*handlerItem) {
	for p, h := range m {
		s.setHandler(p, h)
	}
}

// mergeBuildInNameToPattern merges build-in names into the pattern according to the following
// rules, and the built-in names are named like "{.xxx}".
// Rule 1: The URI in pattern contains the {.struct} keyword, it then replaces the keyword with the struct name;
// Rule 2: The URI in pattern contains the {.method} keyword, it then replaces the keyword with the method name;
// Rule 2: If Rule 1 is not met, it then adds the method name directly to the URI in the pattern;
//
// The parameter <allowAppend> specifies whether allowing appending method name to the tail of pattern.
func (s *Server) mergeBuildInNameToPattern(pattern string, structName, methodName string, allowAppend bool) string {
	structName = s.nameToUri(structName)
	methodName = s.nameToUri(methodName)
	pattern = strings.Replace(pattern, "{.struct}", structName, -1)
	if strings.Index(pattern, "{.method}") != -1 {
		return strings.Replace(pattern, "{.method}", methodName, -1)
	}
	if !allowAppend {
		return pattern
	}
	// Check domain parameter.
	array := strings.Split(pattern, "@")
	uri := array[0]
	uri = strings.TrimRight(uri, "/") + "/" + methodName
	// Append the domain parameter to URI.
	if len(array) > 1 {
		return uri + "@" + array[1]
	}
	return uri
}

// nameToUri converts the given name to URL format using following rules:
// Rule 0: Convert all method names to lowercase, add char '-' between words.
// Rule 1: Do not convert the method name, construct the URI with the original method name.
// Rule 2: Convert all method names to lowercase, no connecting symbols between words.
// Rule 3: Use camel case naming.
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
