// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv

import (
    "errors"
    "fmt"
    "github.com/gogf/gf/g/text/gstr"
    "github.com/gogf/gf/third/github.com/fatih/structs"
    "reflect"
    "strings"
)

// Struct maps the params key-value pairs to the corresponding struct object properties.
// The third parameter mapping is unnecessary, indicating the mapping between the custom name
// and the attribute name.
//
// Note:
// 1. The <params> can be any type of may/struct, usually a map;
// 2. The second parameter <objPointer> should be a pointer to the struct object;
// 3. Only the public attributes of struct object can be mapped;
// 4. If <params> is a map, the key of the map <params> can be lowercase.
//    It will automatically convert the first letter of the key to uppercase
//    in mapping procedure to do the matching.
//    If it does not match, ignore the key;
func Struct(params interface{}, objPointer interface{}, attrMapping...map[string]string) error {
    if params == nil {
        return errors.New("params cannot be nil")
    }
    if objPointer == nil {
        return errors.New("object pointer cannot be nil")
    }
    paramsMap := Map(params)
    if paramsMap == nil {
        return fmt.Errorf("invalid params: %v", params)
    }
    // struct的反射对象
    elem := reflect.Value{}
    if v, ok := objPointer.(reflect.Value); ok {
        elem = v
    } else {
        rv := reflect.ValueOf(objPointer)
        if kind := rv.Kind(); kind != reflect.Ptr {
            return fmt.Errorf("object pointer should be type of: %v", kind)
        }
        if !rv.IsValid() || rv.IsNil() {
            return errors.New("object pointer cannot be nil")
        }
        elem = rv.Elem()
    }
    // 已执行过转换的属性，只执行一次转换。
    // 或者是已经执行过转换检查的属性(即使不进行转换), 以便重复判断。
    doneMap := make(map[string]bool)
    // 首先按照传递的映射关系进行匹配
    if len(attrMapping) > 0 && len(attrMapping[0]) > 0 {
        for mapK, mapV := range attrMapping[0] {
            if v, ok := paramsMap[mapK]; ok {
                doneMap[mapV] = true
                if err := bindVarToStructAttr(elem, mapV, v); err != nil {
                    return err
                }
            }
        }
    }
    // 其次匹配对象定义时绑定的属性名称,
    // 标签映射关系map，如果有的话
    tagMap := getTagMapOfStruct(objPointer)
    for tagk, tagv := range tagMap {
        if _, ok := doneMap[tagv]; ok {
            continue
        }
        if v, ok := paramsMap[tagk]; ok {
            doneMap[tagv] = true
            if err := bindVarToStructAttr(elem, tagv, v); err != nil {
                return err
            }
        }
    }
    // 最后按照默认规则进行匹配
    attrMap  := make(map[string]struct{})
    elemType := elem.Type()
    for i := 0; i < elem.NumField(); i++ {
        // 只转换公开属性
        if !gstr.IsLetterUpper(elemType.Field(i).Name[0]) {
            continue
        }
        attrMap[elemType.Field(i).Name] = struct{}{}
    }
    for mapK, mapV := range paramsMap {
        name := ""
        for _, checkName := range []string {
            gstr.UcFirst(mapK),
            gstr.ReplaceByMap(mapK, map[string]string{
                "_" : "",
                "-" : "",
                " " : "",
            })}  {
            if _, ok := doneMap[checkName]; ok {
                continue
            }
            if _, ok := tagMap[checkName]; ok {
                continue
            }
            // 循环查找属性名称进行匹配
            for value, _ := range attrMap {
                if strings.EqualFold(checkName, value) {
                    name = value
                    break
                }
                if strings.EqualFold(checkName, gstr.Replace(value, "_", "")) {
                    name = value
                    break
                }
            }
            doneMap[checkName] = true
            if name != "" {
                break
            }
        }
        // 如果没有匹配到属性名称，放弃
        if name == "" {
            continue
        }
        if err := bindVarToStructAttr(elem, name, mapV); err != nil {
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
        tag := field.Tag("gconv")
        if tag == "" {
            tag = field.Tag("json")
        }
        if tag != "" {
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

