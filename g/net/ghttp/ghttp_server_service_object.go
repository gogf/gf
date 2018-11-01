// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.
// 服务注册.

package ghttp

import (
    "errors"
    "gitee.com/johng/gf/g/os/glog"
    "strings"
    "reflect"
    "fmt"
    "gitee.com/johng/gf/g/util/gstr"
    "gitee.com/johng/gf/g/os/gfile"
)

// 绑定对象到URI请求处理中，会自动识别方法名称，并附加到对应的URI地址后面
// 第三个参数methods用以指定需要注册的方法，支持多个方法名称，多个方法以英文“,”号分隔，区分大小写
func (s *Server)BindObject(pattern string, obj interface{}, methods...string) error {
    methodMap := (map[string]bool)(nil)
    if len(methods) > 0 {
        methodMap = make(map[string]bool)
        for _, v := range strings.Split(methods[0], ",") {
            methodMap[strings.TrimSpace(v)] = true
        }
    }
    m     := make(handlerMap)
    v     := reflect.ValueOf(obj)
    t     := v.Type()
    sname := t.Elem().Name()
    finit := (func(*Request))(nil)
    fshut := (func(*Request))(nil)
    if v.MethodByName("Init").IsValid() {
        finit = v.MethodByName("Init").Interface().(func(*Request))
    }
    if v.MethodByName("Shut").IsValid() {
        fshut = v.MethodByName("Shut").Interface().(func(*Request))
    }
    pkgPath := t.Elem().PkgPath()
    pkgName := gfile.Basename(pkgPath)
    for i := 0; i < v.NumMethod(); i++ {
        mname := t.Method(i).Name
        if methodMap != nil && !methodMap[mname] {
            continue
        }
        if mname == "Init" || mname == "Shut" {
            continue
        }
        faddr, ok := v.Method(i).Interface().(func(*Request))
        if !ok {
            if methodMap != nil {
                s := fmt.Sprintf(`invalid medthod definition "%s", while "func(*Request))" is required`, v.Method(i).Type().String())
                glog.Error(s)
                return errors.New(s)
            }
            continue
        }
        objName := gstr.Replace(t.String(), fmt.Sprintf(`%s.`, pkgName), "")
        if objName[0] == '*' {
            objName = fmt.Sprintf(`(%s)`, objName)
        }
        key    := s.mergeBuildInNameToPattern(pattern, sname, mname, true)
        m[key]  = &handlerItem {
            name  : fmt.Sprintf(`%s.%s.%s`, pkgPath, objName, mname),
            rtype : gROUTE_REGISTER_OBJECT,
            ctype : nil,
            fname : "",
            faddr : faddr,
            finit : finit,
            fshut : fshut,
        }
        // 如果方法中带有Index方法，那么额外自动增加一个路由规则匹配主URI
        if strings.EqualFold(mname, "Index") {
            p := key
            if strings.EqualFold(p[len(p) - 6:], "/index") {
                p = p[0 : len(p) - 6]
                if len(p) == 0 {
                    p = "/"
                }
            }
            m[p] = &handlerItem {
                name  : fmt.Sprintf(`%s.%s.%s`, pkgPath, objName, mname),
                rtype : gROUTE_REGISTER_OBJECT,
                ctype : nil,
                fname : "",
                faddr : faddr,
                finit : finit,
                fshut : fshut,
            }
        }
    }
    return s.bindHandlerByMap(m)
}

// 绑定对象到URI请求处理中，会自动识别方法名称，并附加到对应的URI地址后面
// 第三个参数methods支持多个方法注册，多个方法以英文“,”号分隔，区分大小写
func (s *Server)BindObjectMethod(pattern string, obj interface{}, method string) error {
    m     := make(handlerMap)
    v     := reflect.ValueOf(obj)
    t     := v.Type()
    sname := t.Elem().Name()
    mname := strings.TrimSpace(method)
    fval  := v.MethodByName(mname)
    if !fval.IsValid() {
        return errors.New("invalid method name:" + mname)
    }
    faddr, ok := fval.Interface().(func(*Request))
    if !ok {
        s := fmt.Sprintf(`invalid medthod definition "%s", while "func(*Request)" is required`, fval.Type().String())
        glog.Error(s)
        return errors.New(s)
    }
    finit := (func(*Request))(nil)
    fshut := (func(*Request))(nil)
    if v.MethodByName("Init").IsValid() {
        finit = v.MethodByName("Init").Interface().(func(*Request))
    }
    if v.MethodByName("Shut").IsValid() {
        fshut = v.MethodByName("Shut").Interface().(func(*Request))
    }
    pkgPath := t.Elem().PkgPath()
    pkgName := gfile.Basename(pkgPath)
    objName := gstr.Replace(t.String(), fmt.Sprintf(`%s.`, pkgName), "")
    if objName[0] == '*' {
        objName = fmt.Sprintf(`(%s)`, objName)
    }
    key   := s.mergeBuildInNameToPattern(pattern, sname, mname, false)
    m[key] = &handlerItem{
        name  : fmt.Sprintf(`%s.%s.%s`, pkgPath, objName, mname),
        rtype : gROUTE_REGISTER_OBJECT,
        ctype : nil,
        fname : "",
        faddr : faddr,
        finit : finit,
        fshut : fshut,
    }

    return s.bindHandlerByMap(m)
}

// 绑定对象到URI请求处理中，会自动识别方法名称，并附加到对应的URI地址后面
// 需要注意对象方法的定义必须按照ghttp.HandlerFunc来定义
func (s *Server)BindObjectRest(pattern string, obj interface{}) error {
    m     := make(handlerMap)
    v     := reflect.ValueOf(obj)
    t     := v.Type()
    finit := (func(*Request))(nil)
    fshut := (func(*Request))(nil)
    if v.MethodByName("Init").IsValid() {
        finit = v.MethodByName("Init").Interface().(func(*Request))
    }
    if v.MethodByName("Shut").IsValid() {
        fshut = v.MethodByName("Shut").Interface().(func(*Request))
    }
    pkgPath := t.Elem().PkgPath()
    for i := 0; i < v.NumMethod(); i++ {
        mname  := t.Method(i).Name
        method := strings.ToUpper(mname)
        if _, ok := s.methodsMap[method]; !ok {
            continue
        }
        faddr, ok := v.Method(i).Interface().(func(*Request))
        if !ok {
            s := fmt.Sprintf(`invalid medthod definition "%s", while "func()" is required`, v.Method(i).Type().String())
            glog.Error(s)
            return errors.New(s)
        }
        pkgName := gfile.Basename(pkgPath)
        objName := gstr.Replace(t.String(), fmt.Sprintf(`%s.`, pkgName), "")
        if objName[0] == '*' {
            objName = fmt.Sprintf(`(%s)`, objName)
        }
        key   := mname + ":" + pattern
        m[key] = &handlerItem {
            name  : fmt.Sprintf(`%s.%s.%s`, pkgPath, objName, mname),
            rtype : gROUTE_REGISTER_OBJECT,
            ctype : nil,
            fname : "",
            faddr : faddr,
            finit : finit,
            fshut : fshut,
        }
    }
    return s.bindHandlerByMap(m)
}