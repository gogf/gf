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

// BindObject registers object to server routes with given pattern.
//
// The optional parameter <method> is used to specify the method to be registered, which
// supports multiple method names, multiple methods are separated by char ',', case-sensitive.
//
// Note that the route method should be defined as ghttp.HandlerFunc.
func (s *Server) BindObject(pattern string, object interface{}, method ...string) {
	var (
		bindMethod = ""
	)
	if len(method) > 0 {
		bindMethod = method[0]
	}
	s.doBindObject(context.TODO(), pattern, object, bindMethod, nil, "")
}

// BindObjectMethod registers specified method of object to server routes with given pattern.
//
// The optional parameter <method> is used to specify the method to be registered, which
// does not supports multiple method names but only one, case-sensitive.
//
// Note that the route method should be defined as ghttp.HandlerFunc.
func (s *Server) BindObjectMethod(pattern string, object interface{}, method string) {
	s.doBindObjectMethod(context.TODO(), pattern, object, method, nil, "")
}

// BindObjectRest registers object in REST API styles to server with specified pattern.
// Note that the route method should be defined as ghttp.HandlerFunc.
func (s *Server) BindObjectRest(pattern string, object interface{}) {
	s.doBindObjectRest(context.TODO(), pattern, object, nil, "")
}

func (s *Server) doBindObject(ctx context.Context, pattern string, object interface{}, method string, middleware []HandlerFunc, source string) {
	// Convert input method to map for convenience and high performance searching purpose.
	var (
		methodMap map[string]bool
	)
	if len(method) > 0 {
		methodMap = make(map[string]bool)
		for _, v := range strings.Split(method, ",") {
			methodMap[strings.TrimSpace(v)] = true
		}
	}
	// If the `method` in `pattern` is `defaultMethod`,
	// it removes for convenience for next statement control.
	domain, method, path, err := s.parsePattern(pattern)
	if err != nil {
		s.Logger().Fatalf(ctx, `%+v`, err)
		return
	}
	if strings.EqualFold(method, defaultMethod) {
		pattern = s.serveHandlerKey("", path, domain)
	}
	var (
		m        = make(map[string]*handlerItem)
		v        = reflect.ValueOf(object)
		t        = v.Type()
		initFunc func(*Request)
		shutFunc func(*Request)
	)
	// If given `object` is not pointer, it then creates a temporary one,
	// of which the value is `v`.
	if v.Kind() == reflect.Struct {
		newValue := reflect.New(t)
		newValue.Elem().Set(v)
		v = newValue
		t = v.Type()
	}
	structName := t.Elem().Name()
	if v.MethodByName("Init").IsValid() {
		initFunc = v.MethodByName("Init").Interface().(func(*Request))
	}
	if v.MethodByName("Shut").IsValid() {
		shutFunc = v.MethodByName("Shut").Interface().(func(*Request))
	}
	pkgPath := t.Elem().PkgPath()
	pkgName := gfile.Basename(pkgPath)
	for i := 0; i < v.NumMethod(); i++ {
		methodName := t.Method(i).Name
		if methodMap != nil && !methodMap[methodName] {
			continue
		}
		if methodName == "Init" || methodName == "Shut" {
			continue
		}
		objName := gstr.Replace(t.String(), fmt.Sprintf(`%s.`, pkgName), "")
		if objName[0] == '*' {
			objName = fmt.Sprintf(`(%s)`, objName)
		}

		funcInfo, err := s.checkAndCreateFuncInfo(v.Method(i).Interface(), pkgPath, objName, methodName)
		if err != nil {
			s.Logger().Fatalf(ctx, `%+v`, err)
		}

		key := s.mergeBuildInNameToPattern(pattern, structName, methodName, true)
		m[key] = &handlerItem{
			Name:       fmt.Sprintf(`%s.%s.%s`, pkgPath, objName, methodName),
			Type:       HandlerTypeObject,
			Info:       funcInfo,
			InitFunc:   initFunc,
			ShutFunc:   shutFunc,
			Middleware: middleware,
			Source:     source,
		}
		// If there's "Index" method, then an additional route is automatically added
		// to match the main URI, for example:
		// If pattern is "/user", then "/user" and "/user/index" are both automatically
		// registered.
		//
		// Note that if there's built-in variables in pattern, this route will not be added
		// automatically.
		if strings.EqualFold(methodName, "Index") && !gregex.IsMatchString(`\{\.\w+\}`, pattern) {
			p := gstr.PosRI(key, "/index")
			k := key[0:p] + key[p+6:]
			if len(k) == 0 || k[0] == '@' {
				k = "/" + k
			}
			m[k] = &handlerItem{
				Name:       fmt.Sprintf(`%s.%s.%s`, pkgPath, objName, methodName),
				Type:       HandlerTypeObject,
				Info:       funcInfo,
				InitFunc:   initFunc,
				ShutFunc:   shutFunc,
				Middleware: middleware,
				Source:     source,
			}
		}
	}
	s.bindHandlerByMap(ctx, m)
}

func (s *Server) doBindObjectMethod(ctx context.Context, pattern string, object interface{}, method string, middleware []HandlerFunc, source string) {
	var (
		m        = make(map[string]*handlerItem)
		v        = reflect.ValueOf(object)
		t        = v.Type()
		initFunc func(*Request)
		shutFunc func(*Request)
	)
	// If given `object` is not pointer, it then creates a temporary one,
	// of which the value is `v`.
	if v.Kind() == reflect.Struct {
		newValue := reflect.New(t)
		newValue.Elem().Set(v)
		v = newValue
		t = v.Type()
	}
	var (
		structName  = t.Elem().Name()
		methodName  = strings.TrimSpace(method)
		methodValue = v.MethodByName(methodName)
	)
	if !methodValue.IsValid() {
		s.Logger().Fatalf(ctx, "invalid method name: %s", methodName)
		return
	}
	if v.MethodByName("Init").IsValid() {
		initFunc = v.MethodByName("Init").Interface().(func(*Request))
	}
	if v.MethodByName("Shut").IsValid() {
		shutFunc = v.MethodByName("Shut").Interface().(func(*Request))
	}
	var (
		pkgPath = t.Elem().PkgPath()
		pkgName = gfile.Basename(pkgPath)
		objName = gstr.Replace(t.String(), fmt.Sprintf(`%s.`, pkgName), "")
	)
	if objName[0] == '*' {
		objName = fmt.Sprintf(`(%s)`, objName)
	}

	funcInfo, err := s.checkAndCreateFuncInfo(methodValue.Interface(), pkgPath, objName, methodName)
	if err != nil {
		s.Logger().Fatalf(ctx, `%+v`, err)
	}

	key := s.mergeBuildInNameToPattern(pattern, structName, methodName, false)
	m[key] = &handlerItem{
		Name:       fmt.Sprintf(`%s.%s.%s`, pkgPath, objName, methodName),
		Type:       HandlerTypeObject,
		Info:       funcInfo,
		InitFunc:   initFunc,
		ShutFunc:   shutFunc,
		Middleware: middleware,
		Source:     source,
	}

	s.bindHandlerByMap(ctx, m)
}

func (s *Server) doBindObjectRest(ctx context.Context, pattern string, object interface{}, middleware []HandlerFunc, source string) {
	var (
		m        = make(map[string]*handlerItem)
		v        = reflect.ValueOf(object)
		t        = v.Type()
		initFunc func(*Request)
		shutFunc func(*Request)
	)
	// If given `object` is not pointer, it then creates a temporary one,
	// of which the value is `v`.
	if v.Kind() == reflect.Struct {
		newValue := reflect.New(t)
		newValue.Elem().Set(v)
		v = newValue
		t = v.Type()
	}
	structName := t.Elem().Name()
	if v.MethodByName(methodNameInit).IsValid() {
		initFunc = v.MethodByName(methodNameInit).Interface().(func(*Request))
	}
	if v.MethodByName(methodNameShut).IsValid() {
		shutFunc = v.MethodByName(methodNameShut).Interface().(func(*Request))
	}
	pkgPath := t.Elem().PkgPath()
	for i := 0; i < v.NumMethod(); i++ {
		methodName := t.Method(i).Name
		if _, ok := methodsMap[strings.ToUpper(methodName)]; !ok {
			continue
		}
		pkgName := gfile.Basename(pkgPath)
		objName := gstr.Replace(t.String(), fmt.Sprintf(`%s.`, pkgName), "")
		if objName[0] == '*' {
			objName = fmt.Sprintf(`(%s)`, objName)
		}

		funcInfo, err := s.checkAndCreateFuncInfo(v.Method(i).Interface(), pkgPath, objName, methodName)
		if err != nil {
			s.Logger().Fatalf(ctx, `%+v`, err)
		}

		key := s.mergeBuildInNameToPattern(methodName+":"+pattern, structName, methodName, false)
		m[key] = &handlerItem{
			Name:       fmt.Sprintf(`%s.%s.%s`, pkgPath, objName, methodName),
			Type:       HandlerTypeObject,
			Info:       funcInfo,
			InitFunc:   initFunc,
			ShutFunc:   shutFunc,
			Middleware: middleware,
			Source:     source,
		}
	}
	s.bindHandlerByMap(ctx, m)
}
