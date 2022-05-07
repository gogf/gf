// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"
)

// BindObject registers object to server routes with a given pattern.
//
// The optional parameter `method` is used to specify the method to be registered, which
// supports multiple method names; multiple methods are separated by char ',', case-sensitive.
func (s *Server) BindObject(pattern string, object interface{}, method ...string) {
	bindMethod := ""
	if len(method) > 0 {
		bindMethod = method[0]
	}
	s.doBindObject(context.TODO(), doBindObjectInput{
		Prefix:     "",
		Pattern:    pattern,
		Object:     object,
		Method:     bindMethod,
		Middleware: nil,
		Source:     "",
	})
}

// BindObjectMethod registers specified method of the object to server routes with a given pattern.
//
// The optional parameter `method` is used to specify the method to be registered, which
// does not support multiple method names but only one, case-sensitive.
func (s *Server) BindObjectMethod(pattern string, object interface{}, method string) {
	s.doBindObjectMethod(context.TODO(), doBindObjectMethodInput{
		Prefix:     "",
		Pattern:    pattern,
		Object:     object,
		Method:     method,
		Middleware: nil,
		Source:     "",
	})
}

// BindObjectRest registers object in REST API styles to server with a specified pattern.
func (s *Server) BindObjectRest(pattern string, object interface{}) {
	s.doBindObjectRest(context.TODO(), doBindObjectInput{
		Prefix:     "",
		Pattern:    pattern,
		Object:     object,
		Method:     "",
		Middleware: nil,
		Source:     "",
	})
}

type doBindObjectInput struct {
	Prefix     string
	Pattern    string
	Object     interface{}
	Method     string
	Middleware []HandlerFunc
	Source     string
}

func (s *Server) doBindObject(ctx context.Context, in doBindObjectInput) {
	// Convert input method to map for convenience and high performance searching purpose.
	var methodMap map[string]bool
	if len(in.Method) > 0 {
		methodMap = make(map[string]bool)
		for _, v := range strings.Split(in.Method, ",") {
			methodMap[strings.TrimSpace(v)] = true
		}
	}
	// If the `method` in `pattern` is `defaultMethod`,
	// it removes for convenience for next statement control.
	domain, method, path, err := s.parsePattern(in.Pattern)
	if err != nil {
		s.Logger().Fatalf(ctx, `%+v`, err)
		return
	}
	if strings.EqualFold(method, defaultMethod) {
		in.Pattern = s.serveHandlerKey("", path, domain)
	}
	var (
		handlerMap   = make(map[string]*HandlerItem)
		reflectValue = reflect.ValueOf(in.Object)
		reflectType  = reflectValue.Type()
		initFunc     func(*Request)
		shutFunc     func(*Request)
	)
	// If given `object` is not pointer, it then creates a temporary one,
	// of which the value is `reflectValue`.
	// It then can retrieve all the methods both of struct/*struct.
	if reflectValue.Kind() == reflect.Struct {
		newValue := reflect.New(reflectType)
		newValue.Elem().Set(reflectValue)
		reflectValue = newValue
		reflectType = reflectValue.Type()
	}
	structName := reflectType.Elem().Name()
	if reflectValue.MethodByName(specialMethodNameInit).IsValid() {
		initFunc = reflectValue.MethodByName(specialMethodNameInit).Interface().(func(*Request))
	}
	if reflectValue.MethodByName(specialMethodNameShut).IsValid() {
		shutFunc = reflectValue.MethodByName(specialMethodNameShut).Interface().(func(*Request))
	}
	pkgPath := reflectType.Elem().PkgPath()
	pkgName := gfile.Basename(pkgPath)
	for i := 0; i < reflectValue.NumMethod(); i++ {
		methodName := reflectType.Method(i).Name
		if methodMap != nil && !methodMap[methodName] {
			continue
		}
		if methodName == specialMethodNameInit || methodName == specialMethodNameShut {
			continue
		}
		objName := gstr.Replace(reflectType.String(), fmt.Sprintf(`%s.`, pkgName), "")
		if objName[0] == '*' {
			objName = fmt.Sprintf(`(%s)`, objName)
		}

		funcInfo, err := s.checkAndCreateFuncInfo(reflectValue.Method(i).Interface(), pkgPath, objName, methodName)
		if err != nil {
			s.Logger().Fatalf(ctx, `%+v`, err)
		}

		key := s.mergeBuildInNameToPattern(in.Pattern, structName, methodName, true)
		handlerMap[key] = &HandlerItem{
			Name:       fmt.Sprintf(`%s.%s.%s`, pkgPath, objName, methodName),
			Type:       HandlerTypeObject,
			Info:       funcInfo,
			InitFunc:   initFunc,
			ShutFunc:   shutFunc,
			Middleware: in.Middleware,
			Source:     in.Source,
		}
		// If there's "Index" method, then an additional route is automatically added
		// to match the main URI, for example:
		// If pattern is "/user", then "/user" and "/user/index" are both automatically
		// registered.
		//
		// Note that if there's built-in variables in pattern, this route will not be added
		// automatically.
		var (
			isIndexMethod = strings.EqualFold(methodName, specialMethodNameIndex)
			hasBuildInVar = gregex.IsMatchString(`\{\.\w+\}`, in.Pattern)
			hashTwoParams = funcInfo.Type.NumIn() == 2
		)
		if isIndexMethod && !hasBuildInVar && !hashTwoParams {
			p := gstr.PosRI(key, "/index")
			k := key[0:p] + key[p+6:]
			if len(k) == 0 || k[0] == '@' {
				k = "/" + k
			}
			handlerMap[k] = &HandlerItem{
				Name:       fmt.Sprintf(`%s.%s.%s`, pkgPath, objName, methodName),
				Type:       HandlerTypeObject,
				Info:       funcInfo,
				InitFunc:   initFunc,
				ShutFunc:   shutFunc,
				Middleware: in.Middleware,
				Source:     in.Source,
			}
		}
	}
	s.bindHandlerByMap(ctx, in.Prefix, handlerMap)
}

type doBindObjectMethodInput struct {
	Prefix     string
	Pattern    string
	Object     interface{}
	Method     string
	Middleware []HandlerFunc
	Source     string
}

func (s *Server) doBindObjectMethod(ctx context.Context, in doBindObjectMethodInput) {
	var (
		handlerMap   = make(map[string]*HandlerItem)
		reflectValue = reflect.ValueOf(in.Object)
		reflectType  = reflectValue.Type()
		initFunc     func(*Request)
		shutFunc     func(*Request)
	)
	// If given `object` is not pointer, it then creates a temporary one,
	// of which the value is `v`.
	if reflectValue.Kind() == reflect.Struct {
		newValue := reflect.New(reflectType)
		newValue.Elem().Set(reflectValue)
		reflectValue = newValue
		reflectType = reflectValue.Type()
	}
	var (
		structName  = reflectType.Elem().Name()
		methodName  = strings.TrimSpace(in.Method)
		methodValue = reflectValue.MethodByName(methodName)
	)
	if !methodValue.IsValid() {
		s.Logger().Fatalf(ctx, "invalid method name: %s", methodName)
		return
	}
	if reflectValue.MethodByName(specialMethodNameInit).IsValid() {
		initFunc = reflectValue.MethodByName(specialMethodNameInit).Interface().(func(*Request))
	}
	if reflectValue.MethodByName(specialMethodNameShut).IsValid() {
		shutFunc = reflectValue.MethodByName(specialMethodNameShut).Interface().(func(*Request))
	}
	var (
		pkgPath = reflectType.Elem().PkgPath()
		pkgName = gfile.Basename(pkgPath)
		objName = gstr.Replace(reflectType.String(), fmt.Sprintf(`%s.`, pkgName), "")
	)
	if objName[0] == '*' {
		objName = fmt.Sprintf(`(%s)`, objName)
	}

	funcInfo, err := s.checkAndCreateFuncInfo(methodValue.Interface(), pkgPath, objName, methodName)
	if err != nil {
		s.Logger().Fatalf(ctx, `%+v`, err)
	}

	key := s.mergeBuildInNameToPattern(in.Pattern, structName, methodName, false)
	handlerMap[key] = &HandlerItem{
		Name:       fmt.Sprintf(`%s.%s.%s`, pkgPath, objName, methodName),
		Type:       HandlerTypeObject,
		Info:       funcInfo,
		InitFunc:   initFunc,
		ShutFunc:   shutFunc,
		Middleware: in.Middleware,
		Source:     in.Source,
	}

	s.bindHandlerByMap(ctx, in.Prefix, handlerMap)
}

func (s *Server) doBindObjectRest(ctx context.Context, in doBindObjectInput) {
	var (
		handlerMap   = make(map[string]*HandlerItem)
		reflectValue = reflect.ValueOf(in.Object)
		reflectType  = reflectValue.Type()
		initFunc     func(*Request)
		shutFunc     func(*Request)
	)
	// If given `object` is not pointer, it then creates a temporary one,
	// of which the value is `v`.
	if reflectValue.Kind() == reflect.Struct {
		newValue := reflect.New(reflectType)
		newValue.Elem().Set(reflectValue)
		reflectValue = newValue
		reflectType = reflectValue.Type()
	}
	structName := reflectType.Elem().Name()
	if reflectValue.MethodByName(specialMethodNameInit).IsValid() {
		initFunc = reflectValue.MethodByName(specialMethodNameInit).Interface().(func(*Request))
	}
	if reflectValue.MethodByName(specialMethodNameShut).IsValid() {
		shutFunc = reflectValue.MethodByName(specialMethodNameShut).Interface().(func(*Request))
	}
	pkgPath := reflectType.Elem().PkgPath()
	for i := 0; i < reflectValue.NumMethod(); i++ {
		methodName := reflectType.Method(i).Name
		if _, ok := methodsMap[strings.ToUpper(methodName)]; !ok {
			continue
		}
		pkgName := gfile.Basename(pkgPath)
		objName := gstr.Replace(reflectType.String(), fmt.Sprintf(`%s.`, pkgName), "")
		if objName[0] == '*' {
			objName = fmt.Sprintf(`(%s)`, objName)
		}

		funcInfo, err := s.checkAndCreateFuncInfo(
			reflectValue.Method(i).Interface(),
			pkgPath,
			objName,
			methodName,
		)
		if err != nil {
			s.Logger().Fatalf(ctx, `%+v`, err)
		}

		key := s.mergeBuildInNameToPattern(methodName+":"+in.Pattern, structName, methodName, false)
		handlerMap[key] = &HandlerItem{
			Name:       fmt.Sprintf(`%s.%s.%s`, pkgPath, objName, methodName),
			Type:       HandlerTypeObject,
			Info:       funcInfo,
			InitFunc:   initFunc,
			ShutFunc:   shutFunc,
			Middleware: in.Middleware,
			Source:     in.Source,
		}
	}
	s.bindHandlerByMap(ctx, in.Prefix, handlerMap)
}
