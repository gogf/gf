// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import (
	"bytes"
	"context"
	"reflect"
	"strings"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/text/gstr"
)

// BindHandler registers a handler function to server with a given pattern.
//
// Note that the parameter `handler` can be type of:
// 1. func(*ghttp.Request)
// 2. func(context.Context, BizRequest)(BizResponse, error).
func (s *Server) BindHandler(pattern string, handler interface{}) {
	ctx := context.TODO()
	funcInfo, err := s.checkAndCreateFuncInfo(handler, "", "", "")
	if err != nil {
		s.Logger().Fatalf(ctx, `%+v`, err)
	}
	s.doBindHandler(ctx, doBindHandlerInput{
		Prefix:     "",
		Pattern:    pattern,
		FuncInfo:   funcInfo,
		Middleware: nil,
		Source:     "",
	})
}

type doBindHandlerInput struct {
	Prefix     string
	Pattern    string
	FuncInfo   handlerFuncInfo
	Middleware []HandlerFunc
	Source     string
}

// doBindHandler registers a handler function to server with given pattern.
//
// The parameter `pattern` is like:
// /user/list, put:/user, delete:/user, post:/user@goframe.org.
func (s *Server) doBindHandler(ctx context.Context, in doBindHandlerInput) {
	s.setHandler(ctx, setHandlerInput{
		Prefix:  in.Prefix,
		Pattern: in.Pattern,
		HandlerItem: &HandlerItem{
			Type:       HandlerTypeHandler,
			Info:       in.FuncInfo,
			Middleware: in.Middleware,
			Source:     in.Source,
		},
	})
}

// bindHandlerByMap registers handlers to server using map.
func (s *Server) bindHandlerByMap(ctx context.Context, prefix string, m map[string]*HandlerItem) {
	for pattern, handler := range m {
		s.setHandler(ctx, setHandlerInput{
			Prefix:      prefix,
			Pattern:     pattern,
			HandlerItem: handler,
		})
	}
}

// mergeBuildInNameToPattern merges build-in names into the pattern according to the following
// rules, and the built-in names are named like "{.xxx}".
// Rule 1: The URI in pattern contains the {.struct} keyword, it then replaces the keyword with the struct name;
// Rule 2: The URI in pattern contains the {.method} keyword, it then replaces the keyword with the method name;
// Rule 2: If Rule 1 is not met, it then adds the method name directly to the URI in the pattern;
//
// The parameter `allowAppend` specifies whether allowing appending method name to the tail of pattern.
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

// nameToUri converts the given name to the URL format using the following rules:
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

func (s *Server) checkAndCreateFuncInfo(f interface{}, pkgPath, structName, methodName string) (info handlerFuncInfo, err error) {
	handlerFunc, ok := f.(HandlerFunc)
	if !ok {
		reflectType := reflect.TypeOf(f)
		if reflectType.NumIn() != 2 || reflectType.NumOut() != 2 {
			if pkgPath != "" {
				err = gerror.NewCodef(
					gcode.CodeInvalidParameter,
					`invalid handler: %s.%s.%s defined as "%s", but "func(*ghttp.Request)" or "func(context.Context, *BizRequest)(*BizResponse, error)" is required`,
					pkgPath, structName, methodName, reflect.TypeOf(f).String(),
				)
			} else {
				err = gerror.NewCodef(
					gcode.CodeInvalidParameter,
					`invalid handler: defined as "%s", but "func(*ghttp.Request)" or "func(context.Context, *BizRequest)(*BizResponse, error)" is required`,
					reflect.TypeOf(f).String(),
				)
			}
			return
		}

		if reflectType.In(0).String() != "context.Context" {
			err = gerror.NewCodef(
				gcode.CodeInvalidParameter,
				`invalid handler: defined as "%s", but the first input parameter should be type of "context.Context"`,
				reflect.TypeOf(f).String(),
			)
			return
		}

		if reflectType.Out(1).String() != "error" {
			err = gerror.NewCodef(
				gcode.CodeInvalidParameter,
				`invalid handler: defined as "%s", but the last output parameter should be type of "error"`,
				reflect.TypeOf(f).String(),
			)
			return
		}

		// The request struct should be named as `xxxReq`.
		if !gstr.HasSuffix(reflectType.In(1).String(), `Req`) {
			err = gerror.NewCodef(
				gcode.CodeInvalidParameter,
				`invalid struct naming for request: defined as "%s", but it should be named with "Req" suffix like "XxxReq"`,
				reflectType.In(1).String(),
			)
			return
		}

		// The response struct should be named as `xxxRes`.
		if !gstr.HasSuffix(reflectType.Out(0).String(), `Res`) {
			err = gerror.NewCodef(
				gcode.CodeInvalidParameter,
				`invalid struct naming for response: defined as "%s", but it should be named with "Res" suffix like "XxxRes"`,
				reflectType.Out(0).String(),
			)
			return
		}
	}
	info.Func = handlerFunc
	info.Type = reflect.TypeOf(f)
	info.Value = reflect.ValueOf(f)
	return
}
