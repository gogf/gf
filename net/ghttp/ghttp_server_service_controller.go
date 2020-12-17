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

// BindController registers controller to server routes with specified pattern. The controller
// needs to implement the gmvc.Controller interface. Each request of the controller bound in
// this way will initialize a new controller object for processing, corresponding to different
// request sessions.
//
// The optional parameter <method> is used to specify the method to be registered, which
// supports multiple method names, multiple methods are separated by char ',', case sensitive.
func (s *Server) BindController(pattern string, controller Controller, method ...string) {
	bindMethod := ""
	if len(method) > 0 {
		bindMethod = method[0]
	}
	s.doBindController(pattern, controller, bindMethod, nil, "")
}

// BindControllerMethod registers specified method to server routes with specified pattern.
//
// The optional parameter <method> is used to specify the method to be registered, which
// does not supports multiple method names but only one, case sensitive.
func (s *Server) BindControllerMethod(pattern string, controller Controller, method string) {
	s.doBindControllerMethod(pattern, controller, method, nil, "")
}

// BindControllerRest registers controller in REST API style to server with specified pattern.
// The controller needs to implement the gmvc.Controller interface. Each request of the controller
// bound in this way will initialize a new controller object for processing, corresponding to
// different request sessions.
// The method will recognize the HTTP method and do REST binding, for example:
// The method "Post" of controller will be bound to the HTTP POST method request processing,
// and the method "Delete" will be bound to the HTTP DELETE method request processing.
// Therefore, only the method corresponding to the HTTP Method will be bound, other methods will
// not automatically register the binding.
func (s *Server) BindControllerRest(pattern string, controller Controller) {
	s.doBindControllerRest(pattern, controller, nil, "")
}

func (s *Server) doBindController(
	pattern string, controller Controller, method string,
	middleware []HandlerFunc, source string,
) {
	// Convert input method to map for convenience and high performance searching.
	var methodMap map[string]bool
	if len(method) > 0 {
		methodMap = make(map[string]bool)
		for _, v := range strings.Split(method, ",") {
			methodMap[strings.TrimSpace(v)] = true
		}
	}
	domain, method, path, err := s.parsePattern(pattern)
	if err != nil {
		s.Logger().Fatal(err)
		return
	}
	if strings.EqualFold(method, defaultMethod) {
		pattern = s.serveHandlerKey("", path, domain)
	}
	// Retrieve a list of methods, create construct corresponding URI.
	m := make(map[string]*handlerItem)
	v := reflect.ValueOf(controller)
	t := v.Type()
	pkgPath := t.Elem().PkgPath()
	pkgName := gfile.Basename(pkgPath)
	structName := t.Elem().Name()
	for i := 0; i < v.NumMethod(); i++ {
		methodName := t.Method(i).Name
		if methodMap != nil && !methodMap[methodName] {
			continue
		}
		if methodName == "Init" || methodName == "Shut" || methodName == "Exit" {
			continue
		}
		ctlName := gstr.Replace(t.String(), fmt.Sprintf(`%s.`, pkgName), "")
		if ctlName[0] == '*' {
			ctlName = fmt.Sprintf(`(%s)`, ctlName)
		}
		if _, ok := v.Method(i).Interface().(func()); !ok {
			if len(methodMap) > 0 {
				// If registering with specified method, print error.
				s.Logger().Errorf(
					`invalid route method: %s.%s.%s defined as "%s", but "func()" is required for controller registry`,
					pkgPath, ctlName, methodName, v.Method(i).Type().String(),
				)
			} else {
				// Else, just print debug information.
				s.Logger().Debugf(
					`ignore route method: %s.%s.%s defined as "%s", no match "func()" for controller registry`,
					pkgPath, ctlName, methodName, v.Method(i).Type().String(),
				)
			}
			continue
		}
		key := s.mergeBuildInNameToPattern(pattern, structName, methodName, true)
		m[key] = &handlerItem{
			itemName: fmt.Sprintf(`%s.%s.%s`, pkgPath, ctlName, methodName),
			itemType: handlerTypeController,
			ctrlInfo: &handlerController{
				name:    methodName,
				reflect: v.Elem().Type(),
			},
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
				itemName: fmt.Sprintf(`%s.%s.%s`, pkgPath, ctlName, methodName),
				itemType: handlerTypeController,
				ctrlInfo: &handlerController{
					name:    methodName,
					reflect: v.Elem().Type(),
				},
				middleware: middleware,
				source:     source,
			}
		}
	}
	s.bindHandlerByMap(m)
}

func (s *Server) doBindControllerMethod(
	pattern string,
	controller Controller,
	method string,
	middleware []HandlerFunc,
	source string,
) {
	m := make(map[string]*handlerItem)
	v := reflect.ValueOf(controller)
	t := v.Type()
	structName := t.Elem().Name()
	methodName := strings.TrimSpace(method)
	methodValue := v.MethodByName(methodName)
	if !methodValue.IsValid() {
		s.Logger().Fatal("invalid method name: " + methodName)
		return
	}
	pkgPath := t.Elem().PkgPath()
	pkgName := gfile.Basename(pkgPath)
	ctlName := gstr.Replace(t.String(), fmt.Sprintf(`%s.`, pkgName), "")
	if ctlName[0] == '*' {
		ctlName = fmt.Sprintf(`(%s)`, ctlName)
	}
	if _, ok := methodValue.Interface().(func()); !ok {
		s.Logger().Errorf(
			`invalid route method: %s.%s.%s defined as "%s", but "func()" is required for controller registry`,
			pkgPath, ctlName, methodName, methodValue.Type().String(),
		)
		return
	}
	key := s.mergeBuildInNameToPattern(pattern, structName, methodName, false)
	m[key] = &handlerItem{
		itemName: fmt.Sprintf(`%s.%s.%s`, pkgPath, ctlName, methodName),
		itemType: handlerTypeController,
		ctrlInfo: &handlerController{
			name:    methodName,
			reflect: v.Elem().Type(),
		},
		middleware: middleware,
		source:     source,
	}
	s.bindHandlerByMap(m)
}

func (s *Server) doBindControllerRest(
	pattern string, controller Controller,
	middleware []HandlerFunc, source string,
) {
	m := make(map[string]*handlerItem)
	v := reflect.ValueOf(controller)
	t := v.Type()
	pkgPath := t.Elem().PkgPath()
	structName := t.Elem().Name()
	for i := 0; i < v.NumMethod(); i++ {
		methodName := t.Method(i).Name
		if _, ok := methodsMap[strings.ToUpper(methodName)]; !ok {
			continue
		}
		pkgName := gfile.Basename(pkgPath)
		ctlName := gstr.Replace(t.String(), fmt.Sprintf(`%s.`, pkgName), "")
		if ctlName[0] == '*' {
			ctlName = fmt.Sprintf(`(%s)`, ctlName)
		}
		if _, ok := v.Method(i).Interface().(func()); !ok {
			s.Logger().Errorf(
				`invalid route method: %s.%s.%s defined as "%s", but "func()" is required for controller registry`,
				pkgPath, ctlName, methodName, v.Method(i).Type().String(),
			)
			return
		}
		key := s.mergeBuildInNameToPattern(methodName+":"+pattern, structName, methodName, false)
		m[key] = &handlerItem{
			itemName: fmt.Sprintf(`%s.%s.%s`, pkgPath, ctlName, methodName),
			itemType: handlerTypeController,
			ctrlInfo: &handlerController{
				name:    methodName,
				reflect: v.Elem().Type(),
			},
			middleware: middleware,
			source:     source,
		}
	}
	s.bindHandlerByMap(m)
}
