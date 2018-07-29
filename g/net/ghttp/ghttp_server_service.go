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
    "gitee.com/johng/gf/g/util/gstr"
)

// 绑定URI到操作函数/方法
// pattern的格式形如：/user/list, put:/user, delete:/user, post:/user@johng.cn
// 支持RESTful的请求格式，具体业务逻辑由绑定的处理方法来执行
func (s *Server)bindHandlerItem(pattern string, item *handlerItem) error {
    if s.status == 1 {
        return errors.New("server handlers cannot be changed while running")
    }
    return s.setHandler(pattern, item)
}

// 通过映射数组绑定URI到操作函数/方法
func (s *Server)bindHandlerByMap(m handlerMap) error {
    for p, h := range m {
        if err := s.bindHandlerItem(p, h); err != nil {
            return err
        }
    }
    return nil
}

// 将内置的名称按照设定的规则合并到pattern中，内置名称按照{.xxx}规则命名。
// 规则1：pattern中的URI包含{.struct}关键字，则替换该关键字为结构体名称；
// 规则1：pattern中的URI包含{.method}关键字，则替换该关键字为方法名称；
// 规则2：如果不满足规则1，那么直接将防发明附加到pattern中的URI后面；
func (s *Server) mergeBuildInNameToPattern(pattern string, structName, methodName string) string {
    structName = s.nameToUrlPart(structName)
    methodName = s.nameToUrlPart(methodName)
    pattern    = strings.Replace(pattern, "{.struct}", structName, -1)
    if strings.Index(pattern, "{.method}") != -1 {
        return strings.Replace(pattern, "{.method}", methodName, -1)
    }
    // 检测域名后缀
    array := strings.Split(pattern, "@")
    // 分离URI(其实可能包含HTTP Method)
    uri := array[0]
    uri  = strings.TrimRight(uri, "/") + "/" + methodName
    // 加上指定域名后缀
    if len(array) > 1 {
        return uri + "@" + array[1]
    }
    return uri
}

// 将给定的名称转换为URL规范格式。
// 规范1: 全部转换为小写；
// 规范2: 方法名中间存在大写字母，转换为小写URI地址以“-”号链接每个单词；
func (s *Server) nameToUrlPart(name string) string {
    part := ""
    for i := 0; i < len(name); i++ {
        if i > 0 && gstr.IsLetterUpper(name[i]) {
            part += "-"
        }
        part += string(name[i])
    }
    return strings.ToLower(part)
}

// 注意该方法是直接绑定函数的内存地址，执行的时候直接执行该方法，不会存在初始化新的控制器逻辑
func (s *Server)BindHandler(pattern string, handler HandlerFunc) error {
    return s.bindHandlerItem(pattern, &handlerItem{
        ctype : nil,
        fname : "",
        faddr : handler,
    })
}

// 绑定对象到URI请求处理中，会自动识别方法名称，并附加到对应的URI地址后面
// 需要注意对象方法的定义必须按照ghttp.HandlerFunc来定义
func (s *Server)BindObject(pattern string, obj interface{}) error {
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
            m[pattern] = &handlerItem {
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
            m[pattern] = &handlerItem {
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

// 绑定控制器，控制器需要实现gmvc.Controller接口
// 这种方式绑定的控制器每一次请求都会初始化一个新的控制器对象进行处理，对应不同的请求会话
func (s *Server)BindController(pattern string, c Controller) error {
    // 遍历控制器，获取方法列表，并构造成uri
    m     := make(handlerMap)
    v     := reflect.ValueOf(c)
    t     := v.Type()
    sname := t.Elem().Name()
    for i := 0; i < v.NumMethod(); i++ {
        mname := t.Method(i).Name
        if mname == "Init" || mname == "Shut" || mname == "Exit"  {
            continue
        }
        key   := s.mergeBuildInNameToPattern(pattern, sname, mname)
        m[key] = &handlerItem {
            ctype : v.Elem().Type(),
            fname : mname,
            faddr : nil,
        }
        // 如果方法中带有Index方法，那么额外自动增加一个路由规则匹配主URI
        if strings.EqualFold(mname, "Index") {
            m[pattern] = &handlerItem {
                ctype : v.Elem().Type(),
                fname : mname,
                faddr : nil,
            }
        }
    }
    return s.bindHandlerByMap(m)
}

// 这种方式绑定的控制器每一次请求都会初始化一个新的控制器对象进行处理，对应不同的请求会话
// 第三个参数methods支持多个方法注册，多个方法以英文“,”号分隔，不区分大小写
func (s *Server)BindControllerMethod(pattern string, c Controller, methods string) error {
    m     := make(handlerMap)
    v     := reflect.ValueOf(c)
    e     := v.Type().Elem()
    t     := v.Elem().Type()
    sname := e.Name()
    for _, method := range strings.Split(methods, ",") {
        mname := strings.TrimSpace(method)
        if !v.MethodByName(mname).IsValid() {
            return errors.New("invalid method name:" + mname)
        }
        key    := s.mergeBuildInNameToPattern(pattern, sname, mname)
        m[key]  = &handlerItem {
            ctype : t,
            fname : mname,
            faddr : nil,
        }
        // 如果方法中带有Index方法，那么额外自动增加一个路由规则匹配主URI
        if strings.EqualFold(mname, "Index") {
            m[pattern] = &handlerItem {
                ctype : t,
                fname : mname,
                faddr : nil,
            }
        }
    }
    return s.bindHandlerByMap(m)
}

// 绑定控制器(RESTFul)，控制器需要实现gmvc.Controller接口
// 方法会识别HTTP方法，并做REST绑定处理，例如：Post方法会绑定到HTTP POST的方法请求处理，Delete方法会绑定到HTTP DELETE的方法请求处理
// 因此只会绑定HTTP Method对应的方法，其他方法不会自动注册绑定
// 这种方式绑定的控制器每一次请求都会初始化一个新的控制器对象进行处理，对应不同的请求会话
func (s *Server)BindControllerRest(pattern string, c Controller) error {
    // 遍历控制器，获取方法列表，并构造成uri
    m := make(handlerMap)
    v := reflect.ValueOf(c)
    t := v.Type()
    // 如果存在与HttpMethod对应名字的方法，那么绑定这些方法
    for i := 0; i < v.NumMethod(); i++ {
        name   := t.Method(i).Name
        method := strings.ToUpper(name)
        if _, ok := s.methodsMap[method]; !ok {
            continue
        }
        key   := name + ":" + pattern
        m[key] = &handlerItem {
            ctype : v.Elem().Type(),
            fname : name,
            faddr : nil,
        }
    }
    return s.bindHandlerByMap(m)
}
