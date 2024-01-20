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
	"github.com/gogf/gf/v2/os/gstructs"
	"github.com/gogf/gf/v2/text/gstr"
)

// BindHandler registers a handler function to server with a given pattern.
//
// Note that the parameter `handler` can be type of:
// 1. func(*ghttp.Request)
// 2. func(context.Context, BizRequest)(BizResponse, error)
func (s *Server) BindHandler(pattern string, handler interface{}) {
	var ctx = context.TODO()
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
// /user/list, put:/user, delete:/user, post:/user@goframe.org
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
	pattern = strings.ReplaceAll(pattern, "{.struct}", structName)
	if strings.Contains(pattern, "{.method}") {
		return strings.ReplaceAll(pattern, "{.method}", methodName)
	}
	if !allowAppend {
		return pattern
	}
	// Check domain parameter.
	var (
		array = strings.Split(pattern, "@")
		uri   = strings.TrimRight(array[0], "/") + "/" + methodName
	)
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

func (s *Server) checkAndCreateFuncInfo(
	f interface{}, pkgPath, structName, methodName string,
) (funcInfo handlerFuncInfo, err error) {
	funcInfo = handlerFuncInfo{
		Type:  reflect.TypeOf(f),
		Value: reflect.ValueOf(f),
	}
	if handlerFunc, ok := f.(HandlerFunc); ok {
		funcInfo.Func = handlerFunc
		return
	}

	var (
		reflectType    = funcInfo.Type
		inputObject    reflect.Value
		inputObjectPtr interface{}
	)
	if reflectType.NumIn() != 2 || reflectType.NumOut() != 2 {
		if pkgPath != "" {
			err = gerror.NewCodef(
				gcode.CodeInvalidParameter,
				`invalid handler: %s.%s.%s defined as "%s", but "func(*ghttp.Request)" or "func(context.Context, *BizReq)(*BizRes, error)" is required`,
				pkgPath, structName, methodName, reflectType.String(),
			)
		} else {
			err = gerror.NewCodef(
				gcode.CodeInvalidParameter,
				`invalid handler: defined as "%s", but "func(*ghttp.Request)" or "func(context.Context, *BizReq)(*BizRes, error)" is required`,
				reflectType.String(),
			)
		}
		return
	}

	if !reflectType.In(0).Implements(reflect.TypeOf((*context.Context)(nil)).Elem()) {
		err = gerror.NewCodef(
			gcode.CodeInvalidParameter,
			`invalid handler: defined as "%s", but the first input parameter should be type of "context.Context"`,
			reflectType.String(),
		)
		return
	}

	if !reflectType.Out(1).Implements(reflect.TypeOf((*error)(nil)).Elem()) {
		err = gerror.NewCodef(
			gcode.CodeInvalidParameter,
			`invalid handler: defined as "%s", but the last output parameter should be type of "error"`,
			reflectType.String(),
		)
		return
	}

	if reflectType.In(1).Kind() != reflect.Ptr ||
		(reflectType.In(1).Kind() == reflect.Ptr && reflectType.In(1).Elem().Kind() != reflect.Struct) {
		err = gerror.NewCodef(
			gcode.CodeInvalidParameter,
			`invalid handler: defined as "%s", but the second input parameter should be type of pointer to struct like "*BizReq"`,
			reflectType.String(),
		)
		return
	}

	// Do not enable this logic, as many users are already using none struct pointer type
	// as the first output parameter.
	/*
		if reflectType.Out(0).Kind() != reflect.Ptr ||
			(reflectType.Out(0).Kind() == reflect.Ptr && reflectType.Out(0).Elem().Kind() != reflect.Struct) {
			err = gerror.NewCodef(
				gcode.CodeInvalidParameter,
				`invalid handler: defined as "%s", but the first output parameter should be type of pointer to struct like "*BizRes"`,
				reflectType.String(),
			)
			return
		}
	*/

	// The request struct should be named as `xxxReq`.
	reqStructName := trimGeneric(reflectType.In(1).String())
	if !gstr.HasSuffix(reqStructName, `Req`) {
		err = gerror.NewCodef(
			gcode.CodeInvalidParameter,
			`invalid struct naming for request: defined as "%s", but it should be named with "Req" suffix like "XxxReq"`,
			reqStructName,
		)
		return
	}

	// The response struct should be named as `xxxRes`.
	resStructName := trimGeneric(reflectType.Out(0).String())
	if !gstr.HasSuffix(resStructName, `Res`) {
		err = gerror.NewCodef(
			gcode.CodeInvalidParameter,
			`invalid struct naming for response: defined as "%s", but it should be named with "Res" suffix like "XxxRes"`,
			resStructName,
		)
		return
	}

	funcInfo.IsStrictRoute = true

	inputObject = reflect.New(funcInfo.Type.In(1).Elem())
	inputObjectPtr = inputObject.Interface()

	// It retrieves and returns the request struct fields.
	fields, err := gstructs.Fields(gstructs.FieldsInput{
		Pointer:         inputObjectPtr,
		RecursiveOption: gstructs.RecursiveOptionEmbedded,
	})
	if err != nil {
		return funcInfo, err
	}
	funcInfo.ReqStructFields = fields
	funcInfo.Func = createRouterFunc(funcInfo)
	return
}

func createRouterFunc(funcInfo handlerFuncInfo) func(r *Request) {
	return func(r *Request) {
		var (
			ok          bool
			err         error
			inputValues = []reflect.Value{
				reflect.ValueOf(r.Context()),
			}
		)
		if funcInfo.Type.NumIn() == 2 {
			var inputObject reflect.Value
			if funcInfo.Type.In(1).Kind() == reflect.Ptr {
				inputObject = reflect.New(funcInfo.Type.In(1).Elem())
				r.error = r.Parse(inputObject.Interface())
			} else {
				inputObject = reflect.New(funcInfo.Type.In(1).Elem()).Elem()
				r.error = r.Parse(inputObject.Addr().Interface())
			}
			if r.error != nil {
				return
			}
			inputValues = append(inputValues, inputObject)
		}
		// Call handler with dynamic created parameter values.
		results := funcInfo.Value.Call(inputValues)
		switch len(results) {
		case 1:
			if !results[0].IsNil() {
				if err, ok = results[0].Interface().(error); ok {
					r.error = err
				}
			}

		case 2:
			r.handlerResponse = results[0].Interface()
			if !results[1].IsNil() {
				if err, ok = results[1].Interface().(error); ok {
					r.error = err
				}
			}
		}
	}
}

// trimGeneric removes type definitions string from response type name if generic
func trimGeneric(structName string) string {
	var (
		leftBraceIndex  = strings.LastIndex(structName, "[") // for generic, it is faster to start at the end than at the beginning
		rightBraceIndex = strings.LastIndex(structName, "]")
	)
	if leftBraceIndex == -1 || rightBraceIndex == -1 {
		// not found '[' or ']'
		return structName
	} else if leftBraceIndex+1 == rightBraceIndex {
		// may be a slice, because generic is '[X]', not '[]'
		// to be compatible with bad return parameter type: []XxxRes
		return structName
	}
	return structName[:leftBraceIndex]
}
