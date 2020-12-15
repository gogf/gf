// Copyright GoFrame Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/text/gregex"
	"github.com/gogf/gf/text/gstr"
)

// BindObject registers object to server routes with given pattern.
//
// The optional parameter <method> is used to specify the method to be registered, which
// supports multiple method names, multiple methods are separated by char ',', case sensitive.
//
// Note that the route method should be defined as ghttp.HandlerFunc.
func (s *Server) BindObject(pattern string, object interface{}, method ...string) {
	bindMethod := ""
	if len(method) > 0 {
		bindMethod = method[0]
	}
	s.doBindObject(pattern, object, bindMethod, nil, "")
}

// BindObjectMethod registers specified method of object to server routes with given pattern.
//
// The optional parameter <method> is used to specify the method to be registered, which
// does not supports multiple method names but only one, case sensitive.
//
// Note that the route method should be defined as ghttp.HandlerFunc.
func (s *Server) BindObjectMethod(pattern string, object interface{}, method string) {
	s.doBindObjectMethod(pattern, object, method, nil, "")
}

// BindObjectRest registers object in REST API style to server with specified pattern.
// Note that the route method should be defined as ghttp.HandlerFunc.
func (s *Server) BindObjectRest(pattern string, object interface{}) {
	s.doBindObjectRest(pattern, object, nil, "")
}

func (s *Server) doBindObject(
	pattern string, object interface{}, method string,
	middleware []HandlerFunc, source string,
) {
	// Convert input method to map for convenience and high performance searching purpose.
	var methodMap map[string]bool
	if len(method) > 0 {
		methodMap = make(map[string]bool)
		for _, v := range strings.Split(method, ",") {
			methodMap[strings.TrimSpace(v)] = true
		}
	}
	// 当pattern中的method为all时，去掉该method，以便于后续方法判断
	domain, method, path, err := s.parsePattern(pattern)
	if err != nil {
		s.Logger().Fatal(err)
		return
	}
	if strings.EqualFold(method, defaultMethod) {
		pattern = s.serveHandlerKey("", path, domain)
	}
	m := make(map[string]*handlerItem)
	v := reflect.ValueOf(object)
	t := v.Type()
	initFunc := (func(*Request))(nil)
	shutFunc := (func(*Request))(nil)
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
		itemFunc, ok := v.Method(i).Interface().(func(*Request))
		if !ok {
			if len(methodMap) > 0 {
				s.Logger().Errorf(
					`invalid route method: %s.%s.%s defined as "%s", but "func(*ghttp.Request)" is required for object registry`,
					pkgPath, objName, methodName, v.Method(i).Type().String(),
				)
			} else {
				s.Logger().Debugf(
					`ignore route method: %s.%s.%s defined as "%s", no match "func(*ghttp.Request)" for object registry`,
					pkgPath, objName, methodName, v.Method(i).Type().String(),
				)
			}
			continue
		}
		key := s.mergeBuildInNameToPattern(pattern, structName, methodName, true)
		m[key] = &handlerItem{
			itemName:   fmt.Sprintf(`%s.%s.%s`, pkgPath, objName, methodName),
			itemType:   handlerTypeObject,
			itemFunc:   itemFunc,
			initFunc:   initFunc,
			shutFunc:   shutFunc,
			middleware: middleware,
			source:     source,
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
				itemName:   fmt.Sprintf(`%s.%s.%s`, pkgPath, objName, methodName),
				itemType:   handlerTypeObject,
				itemFunc:   itemFunc,
				initFunc:   initFunc,
				shutFunc:   shutFunc,
				middleware: middleware,
				source:     source,
			}
		}
	}
	s.bindHandlerByMap(m)
}

func (s *Server) doBindObjectMethod(
	pattern string, object interface{}, method string,
	middleware []HandlerFunc, source string,
) {
	m := make(map[string]*handlerItem)
	v := reflect.ValueOf(object)
	t := v.Type()
	structName := t.Elem().Name()
	methodName := strings.TrimSpace(method)
	methodValue := v.MethodByName(methodName)
	if !methodValue.IsValid() {
		s.Logger().Fatal("invalid method name: " + methodName)
		return
	}
	initFunc := (func(*Request))(nil)
	shutFunc := (func(*Request))(nil)
	if v.MethodByName("Init").IsValid() {
		initFunc = v.MethodByName("Init").Interface().(func(*Request))
	}
	if v.MethodByName("Shut").IsValid() {
		shutFunc = v.MethodByName("Shut").Interface().(func(*Request))
	}
	pkgPath := t.Elem().PkgPath()
	pkgName := gfile.Basename(pkgPath)
	objName := gstr.Replace(t.String(), fmt.Sprintf(`%s.`, pkgName), "")
	if objName[0] == '*' {
		objName = fmt.Sprintf(`(%s)`, objName)
	}
	itemFunc, ok := methodValue.Interface().(func(*Request))
	if !ok {
		s.Logger().Errorf(
			`invalid route method: %s.%s.%s defined as "%s", but "func(*ghttp.Request)" is required for object registry`,
			pkgPath, objName, methodName, methodValue.Type().String(),
		)
		return
	}
	key := s.mergeBuildInNameToPattern(pattern, structName, methodName, false)
	m[key] = &handlerItem{
		itemName:   fmt.Sprintf(`%s.%s.%s`, pkgPath, objName, methodName),
		itemType:   handlerTypeObject,
		itemFunc:   itemFunc,
		initFunc:   initFunc,
		shutFunc:   shutFunc,
		middleware: middleware,
		source:     source,
	}

	s.bindHandlerByMap(m)
}

func (s *Server) doBindObjectRest(
	pattern string, object interface{},
	middleware []HandlerFunc, source string,
) {
	m := make(map[string]*handlerItem)
	v := reflect.ValueOf(object)
	t := v.Type()
	initFunc := (func(*Request))(nil)
	shutFunc := (func(*Request))(nil)
	structName := t.Elem().Name()
	if v.MethodByName("Init").IsValid() {
		initFunc = v.MethodByName("Init").Interface().(func(*Request))
	}
	if v.MethodByName("Shut").IsValid() {
		shutFunc = v.MethodByName("Shut").Interface().(func(*Request))
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
		itemFunc, ok := v.Method(i).Interface().(func(*Request))
		if !ok {
			s.Logger().Errorf(
				`invalid route method: %s.%s.%s defined as "%s", but "func(*ghttp.Request)" is required for object registry`,
				pkgPath, objName, methodName, v.Method(i).Type().String(),
			)
			continue
		}
		key := s.mergeBuildInNameToPattern(methodName+":"+pattern, structName, methodName, false)
		m[key] = &handlerItem{
			itemName:   fmt.Sprintf(`%s.%s.%s`, pkgPath, objName, methodName),
			itemType:   handlerTypeObject,
			itemFunc:   itemFunc,
			initFunc:   initFunc,
			shutFunc:   shutFunc,
			middleware: middleware,
			source:     source,
		}
	}
	s.bindHandlerByMap(m)
}
