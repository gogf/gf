// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package gconv

import (
    "errors"
    "fmt"
    "gitee.com/johng/gf/g/container/gset"
    "gitee.com/johng/gf/g/util/gstr"
    "gitee.com/johng/gf/third/github.com/fatih/structs"
    "reflect"
    "strings"
)

// 将params键值对参数映射到对应的struct对象属性上，第三个参数mapping为非必需，表示自定义名称与属性名称的映射关系。
// 需要注意：
// 1、第二个参数应当为struct对象指针；
// 2、struct对象的**公开属性(首字母大写)**才能被映射赋值；
// 3、map中的键名可以为小写，映射转换时会自动将键名首字母转为大写做匹配映射，如果无法匹配则忽略；
func Struct(params interface{}, objPointer interface{}, attrMapping...map[string]string) error {
    if params == nil {
        return nil
    }
    isParamMap := true
    paramsMap  := (map[string]interface{})(nil)
    // 先将参数转为 map[string]interface{} 类型
    if m, ok := params.(map[string]interface{}); ok {
        paramsMap = m
    } else {
        paramsMap = make(map[string]interface{})
        if reflect.ValueOf(params).Kind() == reflect.Map {
            ks := reflect.ValueOf(params).MapKeys()
            vs := reflect.ValueOf(params)
            for _, k := range ks {
                paramsMap[String(k.Interface())] = vs.MapIndex(k).Interface()
            }
        } else {
            isParamMap = false
        }
    }
    // struct的反射对象
    elem := reflect.Value{}
    if v, ok := objPointer.(reflect.Value); ok {
        elem = v
    } else {
        elem = reflect.ValueOf(objPointer).Elem()
    }
    // 如果给定的参数不是map类型，那么直接将参数值映射到第一个属性上
    if !isParamMap {
        if err := bindVarToStructByIndex(elem, 0, params); err != nil {
            return err
        }
        return nil
    }
    // 已执行过转换的属性，只执行一次转换
    dmap := make(map[string]bool)
    // 首先按照传递的映射关系进行匹配
    if len(attrMapping) > 0 && len(attrMapping[0]) > 0 {
        for mappingk, mappingv := range attrMapping[0] {
            if v, ok := paramsMap[mappingk]; ok {
                dmap[mappingv] = true
                if err := bindVarToStructAttr(elem, mappingv, v); err != nil {
                    return err
                }
            }
        }
    }
    // 其次匹配对象定义时绑定的属性名称
    // 标签映射关系map，如果有的话
    tagmap := getTagMapOfStruct(objPointer)
    for tagk, tagv := range tagmap {
        if _, ok := dmap[tagv]; ok {
            continue
        }
        if v, ok := paramsMap[tagk]; ok {
            dmap[tagv] = true
            if err := bindVarToStructAttr(elem, tagv, v); err != nil {
                return err
            }
        }
    }
    // 最后按照默认规则进行匹配
    attrset  := gset.NewStringSet(true)
    elemtype := elem.Type()
    for i := 0; i < elem.NumField(); i++ {
        attrset.Add(elemtype.Field(i).Name)
    }
    for mapk, mapv := range paramsMap {
        name := ""
        for _, checkName := range []string {
            gstr.UcFirst(mapk),
            gstr.ReplaceByMap(mapk, map[string]string{
                "_" : "",
                "-" : "",
                " " : "",
            })} {
            if _, ok := dmap[checkName]; ok {
                continue
            }
            if _, ok := tagmap[checkName]; ok {
                continue
            }
            // 循环查找属性名称进行匹配
            attrset.Iterator(func(value string) bool {
                if strings.EqualFold(checkName, value) {
                    name = value
                    return false
                }
                if strings.EqualFold(checkName, gstr.Replace(value, "_", "")) {
                    name = value
                    return false
                }
                return true
            })
            if name != "" {
                break
            }
        }
        // 如果没有匹配到属性名称，放弃
        if name == "" {
            continue
        }
        if err := bindVarToStructAttr(elem, name, mapv); err != nil {
            return err
        }
    }
    return nil
}

// 解析指针对象的tag
func getTagMapOfStruct(objPointer interface{}) map[string]string {
    tagmap := make(map[string]string)
    // 反射类型判断
    fields := ([]*structs.Field)(nil)
    if v, ok := objPointer.(reflect.Value); ok {
        fields = structs.Fields(v.Interface())
    } else {
        fields = structs.Fields(objPointer)
    }
    // 将struct中定义的属性转换名称构建成tagmap
    for _, field := range fields {
        if tag := field.Tag("gconv"); tag != "" {
            for _, v := range strings.Split(tag, ",") {
                tagmap[strings.TrimSpace(v)] = field.Name()
            }
        }
    }
    return tagmap
}

// 将参数值绑定到对象指定名称的属性上
func bindVarToStructAttr(elem reflect.Value, name string, value interface{}) (err error) {
    structFieldValue := elem.FieldByName(name)
    // 键名与对象属性匹配检测，map中如果有struct不存在的属性，那么不做处理，直接return
    if !structFieldValue.IsValid() {
        return nil
    }
    // CanSet的属性必须为公开属性(首字母大写)
    if !structFieldValue.CanSet() {
        return nil
    }
    // 必须将value转换为struct属性的数据类型，这里必须用到gconv包
    defer func() {
        // 如果转换失败，那么可能是类型不匹配造成(例如属性包含自定义类型)，那么执行递归转换
        if recover() != nil {
            err = bindVarToReflectValue(structFieldValue, value)
        }
    }()
    structFieldValue.Set(reflect.ValueOf(Convert(value, structFieldValue.Type().String())))
    return nil
}

// 将参数值绑定到对象指定索引位置的属性上
func bindVarToStructByIndex(elem reflect.Value, index int, value interface{}) (err error) {
    structFieldValue := elem.FieldByIndex([]int{index})
    // 键名与对象属性匹配检测
    if !structFieldValue.IsValid() {
        return nil
    }
    // CanSet的属性必须为公开属性(首字母大写)
    if !structFieldValue.CanSet() {
        return nil
    }
    // 必须将value转换为struct属性的数据类型，这里必须用到gconv包
    defer func() {
        // 如果转换失败，那么可能是类型不匹配造成(例如属性包含自定义类型)，那么执行递归转换
        if recover() != nil {
            err = bindVarToReflectValue(structFieldValue, value)
        }
    }()
    structFieldValue.Set(reflect.ValueOf(Convert(value, structFieldValue.Type().String())))
    return nil
}

// 当默认的基本类型转换失败时，通过recover判断后执行反射类型转换(处理复杂类型)
func bindVarToReflectValue(structFieldValue reflect.Value, value interface{}) error {
    switch structFieldValue.Kind() {
        // 属性为结构体
        case reflect.Struct:
            Struct(value, structFieldValue)

        // 属性为数组类型
        case reflect.Slice: fallthrough
        case reflect.Array:
            a := reflect.Value{}
            v := reflect.ValueOf(value)
            if v.Kind() == reflect.Slice || v.Kind() == reflect.Array {
                if v.Len() > 0 {
                    a  = reflect.MakeSlice(structFieldValue.Type(), v.Len(), v.Len())
                    t := a.Index(0).Type()
                    for i := 0; i < v.Len(); i++ {
                        if t.Kind() == reflect.Ptr {
                            e := reflect.New(t.Elem()).Elem()
                            Struct(v.Index(i).Interface(), e)
                            a.Index(i).Set(e.Addr())
                        } else {
                            e := reflect.New(t).Elem()
                            Struct(v.Index(i).Interface(), e)
                            a.Index(i).Set(e)
                        }
                    }
                }
            } else {
                a  = reflect.MakeSlice(structFieldValue.Type(), 1, 1)
                t := a.Index(0).Type()
                if t.Kind() == reflect.Ptr {
                    e := reflect.New(t.Elem()).Elem()
                    Struct(value, e)
                    a.Index(0).Set(e.Addr())
                } else {
                    e := reflect.New(t).Elem()
                    Struct(value, e)
                    a.Index(0).Set(e)
                }
            }
            structFieldValue.Set(a)

        // 属性为指针类型
        case reflect.Ptr:
            e := reflect.New(structFieldValue.Type().Elem()).Elem()
            Struct(value, e)
            structFieldValue.Set(e.Addr())

        default:
            return errors.New(fmt.Sprintf(`cannot convert to type "%s"`, structFieldValue.Type().String()))
    }
    return nil
}

