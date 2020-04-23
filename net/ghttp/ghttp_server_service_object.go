// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
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

// 绑定对象到URI请求处理中，会自动识别方法名称，并附加到对应的URI地址后面
// 第三个参数methods用以指定需要注册的方法，支持多个方法名称，多个方法以英文“,”号分隔，区分大小写
func (s *Server) BindObject(pattern string, object interface{}, method ...string) {
	bindMethod := ""
	if len(method) > 0 {
		bindMethod = method[0]
	}
	s.doBindObject(pattern, object, bindMethod, nil, "")
}

// 绑定对象到URI请求处理中，会自动识别方法名称，并附加到对应的URI地址后面，
// 第三个参数method仅支持一个方法注册，不支持多个，并且区分大小写。
func (s *Server) BindObjectMethod(pattern string, object interface{}, method string) {
	s.doBindObjectMethod(pattern, object, method, nil, "")
}

// 绑定对象到URI请求处理中，会自动识别方法名称，并附加到对应的URI地址后面,
// 需要注意对象方法的定义必须按照 ghttp.HandlerFunc 来定义
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
	if strings.EqualFold(method, gDEFAULT_METHOD) {
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
				// 指定的方法名称注册，那么需要使用错误提示
				s.Logger().Errorf(
					`invalid route method: %s.%s.%s defined as "%s", but "func(*ghttp.Request)" is required for object registry`,
					pkgPath, objName, methodName, v.Method(i).Type().String(),
				)
			} else {
				// 否则只是Debug提示
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
			itemType:   gHANDLER_TYPE_OBJECT,
			itemFunc:   itemFunc,
			initFunc:   initFunc,
			shutFunc:   shutFunc,
			middleware: middleware,
			source:     source,
		}
		// 如果方法中带有Index方法，那么额外自动增加一个路由规则匹配主URI。
		// 注意，当pattern带有内置变量时，不会自动加该路由。
		if strings.EqualFold(methodName, "Index") && !gregex.IsMatchString(`\{\.\w+\}`, pattern) {
			p := gstr.PosRI(key, "/index")
			k := key[0:p] + key[p+6:]
			if len(k) == 0 || k[0] == '@' {
				k = "/" + k
			}
			m[k] = &handlerItem{
				itemName:   fmt.Sprintf(`%s.%s.%s`, pkgPath, objName, methodName),
				itemType:   gHANDLER_TYPE_OBJECT,
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

// 绑定对象到URI请求处理中，会自动识别方法名称，并附加到对应的URI地址后面，
// 第三个参数method仅支持一个方法注册，不支持多个，并且区分大小写。
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
		s.Logger().Errorf(`invalid route method: %s.%s.%s defined as "%s", but "func(*ghttp.Request)" is required for object registry`,
			pkgPath, objName, methodName, methodValue.Type().String())
		return
	}
	key := s.mergeBuildInNameToPattern(pattern, structName, methodName, false)
	m[key] = &handlerItem{
		itemName:   fmt.Sprintf(`%s.%s.%s`, pkgPath, objName, methodName),
		itemType:   gHANDLER_TYPE_OBJECT,
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
			s.Logger().Errorf(`invalid route method: %s.%s.%s defined as "%s", but "func(*ghttp.Request)" is required for object registry`,
				pkgPath, objName, methodName, v.Method(i).Type().String())
			continue
		}
		key := s.mergeBuildInNameToPattern(methodName+":"+pattern, structName, methodName, false)
		m[key] = &handlerItem{
			itemName:   fmt.Sprintf(`%s.%s.%s`, pkgPath, objName, methodName),
			itemType:   gHANDLER_TYPE_OBJECT,
			itemFunc:   itemFunc,
			initFunc:   initFunc,
			shutFunc:   shutFunc,
			middleware: middleware,
			source:     source,
		}
	}
	s.bindHandlerByMap(m)
}
