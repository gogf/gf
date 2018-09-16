// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package gconv

import (
    "gitee.com/johng/gf/g/util/gstr"
    "reflect"
    "github.com/fatih/structs"
    "strings"
)

// 将params键值对参数映射到对应的struct对象属性上，第三个参数mapping为非必需，表示自定义名称与属性名称的映射关系。
// 需要注意：
// 1、第二个参数为struct对象指针；
// 2、struct对象的**公开属性(首字母大写)**才能被映射赋值；
// 3、map中的键名可以为小写，映射转换时会自动将键名首字母转为大写做匹配映射，如果无法匹配则忽略；
func MapToStruct(params map[string]interface{}, object interface{}, mapping...map[string]string) error {
    tagmap := make(map[string]string)
    fields := structs.Fields(object)
    // 将struct中定义的属性转换名称构建称tagmap
    for _, field := range fields {
        if tag := field.Tag("gconv"); tag != "" {
            for _, v := range strings.Split(tag, ",") {
                tagmap[strings.TrimSpace(v)] = field.Name()
            }
        }
    }
    elem := reflect.ValueOf(object).Elem()
    dmap := make(map[string]bool)
    // 首先按照传递的映射关系进行匹配
    if len(mapping) > 0 && len(mapping[0]) > 0 {
        for mappingk, mappingv := range mapping[0] {
            if v, ok := params[mappingk]; ok {
                dmap[mappingv] = true
                bindVarToStruct(elem, mappingv, v)
            }
        }
    }
    // 其次匹配对象定义时绑定的属性名称
    for tagk, tagv := range tagmap {
        if _, ok := dmap[tagv]; ok {
            continue
        }
        if v, ok := params[tagk]; ok {
            dmap[tagv] = true
            bindVarToStruct(elem, tagv, v)
        }
    }
    // 最后按照默认规则进行匹配
    for mapk, mapv := range params {
        name := gstr.UcFirst(mapk)
        if _, ok := dmap[name]; ok {
            continue
        }
        // 后续tag逻辑中会处理的key(重复的键名)这里便不处理
        if _, ok := tagmap[mapk]; !ok {
            bindVarToStruct(elem, name, mapv)
        }
    }
    return nil
}

// 将参数值绑定到对象指定名称的属性上
func bindVarToStruct(elem reflect.Value, name string, value interface{}) {
    structFieldValue := elem.FieldByName(name)
    // 键名与对象属性匹配检测
    if !structFieldValue.IsValid() {
        return
    }
    // CanSet的属性必须为公开属性(首字母大写)
    if !structFieldValue.CanSet() {
        return
    }
    // 必须将value转换为struct属性的数据类型，这里必须用到gconv包
    structFieldValue.Set(reflect.ValueOf(Convert(value, structFieldValue.Type().String())))
}
