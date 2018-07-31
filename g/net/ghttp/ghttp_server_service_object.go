// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.
// 服务注册.

package ghttp

import (
    "errors"
    "strings"
    "reflect"
)

// 绑定对象到URI请求处理中，会自动识别方法名称，并附加到对应的URI地址后面
// 需要注意对象方法的定义必须按照ghttp.HandlerFunc来定义
func (s *Server)BindObject(pattern string, obj interface{}, methods...string) error {
    if len(methods) > 0 {
        return s.BindObjectMethod(pattern, obj, strings.Join(methods, ","))
    }
    m := make(handlerMap)
    v := reflect.ValueOf(obj)
    t := v.Type()
    sname := t.Elem().Name()
    for i := 0; i < v.NumMethod(); i++ {
        method := t.Method(i).Name
        key    := s.mergeBuildInNameToPattern(pattern, sname, method)
        m[key]  = &handlerItem {
            ctype : nil,
            fname : "",
            faddr : v.Method(i).Interface().(func(*Request)),
        }
        // 如果方法中带有Index方法，那么额外自动增加一个路由规则匹配主URI
        if strings.EqualFold(method, "Index") {
            p := key
            if strings.EqualFold(p[len(p) - 6:], "/index") {
                p = p[0 : len(p) - 6]
            }
            m[p] = &handlerItem {
                ctype : nil,
                fname : "",
                faddr : v.Method(i).Interface().(func(*Request)),
            }
        }
    }
    return s.bindHandlerByMap(m)
}

// 绑定对象到URI请求处理中，会自动识别方法名称，并附加到对应的URI地址后面
// 第三个参数methods支持多个方法注册，多个方法以英文“,”号分隔，区分大小写
func (s *Server)BindObjectMethod(pattern string, obj interface{}, methods string) error {
    m     := make(handlerMap)
    v     := reflect.ValueOf(obj)
    t     := v.Type()
    sname := t.Elem().Name()
    for _, method := range strings.Split(methods, ",") {
        mname  := strings.TrimSpace(method)
        fval   := v.MethodByName(mname)
        if !fval.IsValid() {
            return errors.New("invalid method name:" + mname)
        }
        key   := s.mergeBuildInNameToPattern(pattern, sname, mname)
        m[key] = &handlerItem{
            ctype : nil,
            fname : "",
            faddr : fval.Interface().(func(*Request)),
        }
        // 如果方法中带有Index方法，那么额外自动增加一个路由规则匹配主URI
        if strings.EqualFold(mname, "Index") {
            p := key
            if strings.EqualFold(p[len(p) - 6:], "/index") {
                p = p[0 : len(p) - 6]
            }
            m[p] = &handlerItem {
                ctype : nil,
                fname : "",
                faddr : fval.Interface().(func(*Request)),
            }
        }
    }
    return s.bindHandlerByMap(m)
}

// 绑定对象到URI请求处理中，会自动识别方法名称，并附加到对应的URI地址后面
// 需要注意对象方法的定义必须按照ghttp.HandlerFunc来定义
func (s *Server)BindObjectRest(pattern string, obj interface{}) error {
    m := make(handlerMap)
    v := reflect.ValueOf(obj)
    t := v.Type()
    for i := 0; i < v.NumMethod(); i++ {
        name   := t.Method(i).Name
        method := strings.ToUpper(name)
        if _, ok := s.methodsMap[method]; !ok {
            continue
        }
        key   := name + ":" + pattern
        m[key] = &handlerItem {
            ctype : nil,
            fname : "",
            faddr : v.Method(i).Interface().(func(*Request)),
        }
    }
    return s.bindHandlerByMap(m)
}