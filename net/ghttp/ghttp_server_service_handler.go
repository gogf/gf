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
	"github.com/gogf/gf/v2/util/gtag"
)

// BindHandler registers a handler function to server with a given pattern.
//
// Note that the parameter `handler` can be type of:
// 1. func(*ghttp.Request)
// 2. func(context.Context, BizRequest)(BizResponse, error)
func (s *Server) BindHandler(pattern string, handler any) {
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
	f any, pkgPath, structName, methodName string,
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
		inputObjectPtr any
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

	if reflectType.In(1).Kind() != reflect.Pointer ||
		(reflectType.In(1).Kind() == reflect.Pointer && reflectType.In(1).Elem().Kind() != reflect.Struct) {
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
		if reflectType.Out(0).Kind() != reflect.Pointer ||
			(reflectType.Out(0).Kind() == reflect.Pointer && reflectType.Out(0).Elem().Kind() != reflect.Struct) {
			err = gerror.NewCodef(
				gcode.CodeInvalidParameter,
				`invalid handler: defined as "%s", but the first output parameter should be type of pointer to struct like "*BizRes"`,
				reflectType.String(),
			)
			return
		}
	*/

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
	funcInfo.ReqStructTags = scanRequestStructTags(funcInfo.Type.In(1), fields)
	funcInfo.Func = createRouterFunc(funcInfo)
	return
}

// requestStructTags marks the tag usages of a request struct, which is
// prechecked once at handler registration for the request handling to skip the
// tag-driven work that the struct does not use.
//
// Note that it never skips the work that might take effect: the raw tag values
// are checked instead of the parsed ones, as the tag variables can be
// registered or updated later in the boot procedure, see package gtag; and the
// tags are checked on exactly the scope that their handling covers, see
// scanRequestStructTags.
type requestStructTags struct {
	HasValid   bool // Whether any field uses the `valid`/`v` tag.
	HasIn      bool // Whether any field uses the `in` tag.
	HasDefault bool // Whether any field uses the `default`/`d` tag.
}

// scanRequestStructTags reports the tag usages of the given request struct,
// which is precomputed at handler registration for the request handling to skip
// the tag-driven work that the struct does not use.
//
// It is called once at handler registration.
//
// Note that the raw tag values are checked instead of the parsed ones, as the
// tag variables might be registered or updated later in the boot procedure, see
// package gtag. Keep the checked scopes in sync with their handling:
//   - the `in` and `default` tags are checked from the given request struct
//     fields, which are exactly the fields the parameter merging handles;
//   - the `valid` tag is checked from the whole type graph, as the validation
//     also walks into the nested attribute types, see the doCheckValueRecursively
//     of package gvalid.
func scanRequestStructTags(requestType reflect.Type, requestFields []gstructs.Field) (tags requestStructTags) {
	for _, field := range requestFields {
		if field.Field.Tag.Get(gtag.In) != "" {
			tags.HasIn = true
		}
		if field.Field.Tag.Get(gtag.Default) != "" || field.Field.Tag.Get(gtag.DefaultShort) != "" {
			tags.HasDefault = true
		}
	}
	tags.HasValid = hasValidTagInType(requestType, make(map[reflect.Type]struct{}))
	return
}

// hasValidTagInType checks and returns whether the `valid` tag is used by the
// given type or by any type reachable from its attributes.
func hasValidTagInType(structType reflect.Type, visited map[reflect.Type]struct{}) bool {
	// It walks the pointer types by dereferencing, and marks each type as visited,
	// which avoids both the infinite recursion for the self reference types and the
	// infinite dereference for the cyclic pointer types like `type P *P`, whose
	// element type is itself.
	for {
		if _, ok := visited[structType]; ok {
			return false
		}
		visited[structType] = struct{}{}
		if structType.Kind() != reflect.Pointer {
			break
		}
		structType = structType.Elem()
	}
	switch structType.Kind() {
	case reflect.Struct:
		for i := 0; i < structType.NumField(); i++ {
			var field = structType.Field(i)
			if field.Tag.Get(gtag.Valid) != "" || field.Tag.Get(gtag.ValidShort) != "" {
				return true
			}
			if hasValidTagInType(field.Type, visited) {
				return true
			}
		}

	case reflect.Slice, reflect.Array, reflect.Map:
		return hasValidTagInType(structType.Elem(), visited)
	}
	return false
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
			if funcInfo.Type.In(1).Kind() == reflect.Pointer {
				inputObject = reflect.New(funcInfo.Type.In(1).Elem())
				// The validation is prechecked unnecessary for the request struct
				// at handler registration, see handlerFuncInfo.ReqStructTags.
				// The flags are read from the handler item here, which is the same
				// source the parameter merging reads, instead of the copied funcInfo.
				r.error = r.doParse(inputObject.Interface(), parseTypeRequest,
					r.serveHandler.Handler.Info.ReqStructTags.HasValid)
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
