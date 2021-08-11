// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import (
	"bytes"
	"github.com/gogf/gf/debug/gdebug"
	"github.com/gogf/gf/errors/gerror"
	"reflect"
	"strings"

	"github.com/gogf/gf/text/gstr"
)

// BindHandler registers a handler function to server with given pattern.
// The parameter `handler` can be type of:
// func(*ghttp.Request)
// func(context.Context)
// func(context.Context,TypeRequest)
// func(context.Context,TypeRequest) error
// func(context.Context,TypeRequest)(TypeResponse,error)
func (s *Server) BindHandler(pattern string, handler interface{}) {
	funcInfo, err := s.checkAndCreateFuncInfo(handler, "", "", "")
	if err != nil {
		s.Logger().Error(err.Error())
		return
	}
	s.doBindHandler(pattern, funcInfo, nil, "")
}

// doBindHandler registers a handler function to server with given pattern.
// The parameter <pattern> is like:
// /user/list, put:/user, delete:/user, post:/user@goframe.org
func (s *Server) doBindHandler(pattern string, funcInfo handlerFuncInfo, middleware []HandlerFunc, source string) {
	s.setHandler(pattern, &handlerItem{
		Name:       gdebug.FuncPath(funcInfo.Func),
		Type:       handlerTypeHandler,
		Info:       funcInfo,
		Middleware: middleware,
		Source:     source,
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
	case UriTypeFullName:
		return name

	case UriTypeAllLower:
		return strings.ToLower(name)

	case UriTypeCamel:
		part := bytes.NewBuffer(nil)
		if gstr.IsLetterUpper(name[0]) {
			part.WriteByte(name[0] + 32)
		} else {
			part.WriteByte(name[0])
		}
		part.WriteString(name[1:])
		return part.String()

	case UriTypeDefault:
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

func (s *Server) checkAndCreateFuncInfo(f interface{}, pkgPath, objName, methodName string) (info handlerFuncInfo, err error) {
	handlerFunc, ok := f.(HandlerFunc)
	if !ok {
		reflectType := reflect.TypeOf(f)
		if reflectType.NumIn() == 0 || reflectType.NumIn() > 2 || reflectType.NumOut() > 2 {
			if pkgPath != "" {
				err = gerror.NewCodef(
					gerror.CodeInvalidParameter,
					`invalid handler: %s.%s.%s defined as "%s", but "func(*ghttp.Request)" or "func(context.Context)/func(context.Context,Request)/func(context.Context,Request) error/func(context.Context,Request)(Response,error)" is required`,
					pkgPath, objName, methodName, reflect.TypeOf(f).String(),
				)
			} else {
				err = gerror.NewCodef(
					gerror.CodeInvalidParameter,
					`invalid handler: defined as "%s", but "func(*ghttp.Request)" or "func(context.Context)/func(context.Context,Request)/func(context.Context,Request) error/func(context.Context,Request)(Response,error)" is required`,
					reflect.TypeOf(f).String(),
				)
			}
			return
		}

		if reflectType.In(0).String() != "context.Context" {
			err = gerror.NewCodef(
				gerror.CodeInvalidParameter,
				`invalid handler: defined as "%s", but the first input parameter should be type of "context.Context"`,
				reflect.TypeOf(f).String(),
			)
			return
		}

		if reflectType.NumOut() > 0 && reflectType.Out(reflectType.NumOut()-1).String() != "error" {
			err = gerror.NewCodef(
				gerror.CodeInvalidParameter,
				`invalid handler: defined as "%s", but the last output parameter should be type of "error"`,
				reflect.TypeOf(f).String(),
			)
			return
		}
	}
	info.Func = handlerFunc
	info.Type = reflect.TypeOf(f)
	info.Value = reflect.ValueOf(f)
	return
}
