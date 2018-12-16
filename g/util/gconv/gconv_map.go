// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package gconv

import (
    "reflect"
)

// 任意类型转换为 map[string]interface{} 类型,
// 如果给定的输入参数i不是map类型，那么转换会失败，返回nil.
// 当i为struct对象时，第二个参数noTagCheck表示不检测json标签，否则将会使用json tag作为map的键名。
func Map(i interface{}, noTagCheck...bool) map[string]interface{} {
    if i == nil {
        return nil
    }
    if r, ok := i.(map[string]interface{}); ok {
        return r
    } else {
        // 仅对常见的几种map组合进行断言，最后才会使用反射
        m := make(map[string]interface{})
        switch i.(type) {
            case map[interface{}]interface{}:
                for k, v := range i.(map[interface{}]interface{}) {
                    m[String(k)] = v
                }
            case map[interface{}]string:
                for k, v := range i.(map[interface{}]string) {
                    m[String(k)] = v
                }
            case map[interface{}]int:
                for k, v := range i.(map[interface{}]int) {
                    m[String(k)] = v
                }
            case map[interface{}]uint:
                for k, v := range i.(map[interface{}]uint) {
                    m[String(k)] = v
                }
            case map[interface{}]float32:
                for k, v := range i.(map[interface{}]float32) {
                    m[String(k)] = v
                }
            case map[interface{}]float64:
                for k, v := range i.(map[interface{}]float64) {
                    m[String(k)] = v
                }

            case map[string]bool:
                for k, v := range i.(map[string]bool) {
                    m[k] = v
                }
            case map[string]int:
                for k, v := range i.(map[string]int) {
                    m[k] = v
                }
            case map[string]uint:
                for k, v := range i.(map[string]uint) {
                    m[k] = v
                }
            case map[string]float32:
                for k, v := range i.(map[string]float32) {
                    m[k] = v
                }
            case map[string]float64:
                for k, v := range i.(map[string]float64) {
                    m[k] = v
                }

            case map[int]interface{}:
                for k, v := range i.(map[int]interface{}) {
                    m[String(k)] = v
                }
            case map[int]string:
                for k, v := range i.(map[int]string) {
                    m[String(k)] = v
                }
            case map[uint]string:
                for k, v := range i.(map[uint]string) {
                    m[String(k)] = v
                }
            // 不是常见类型，则使用反射
            default:
                rv   := reflect.ValueOf(i)
                kind := rv.Kind()
                // 如果是指针，那么需要转换到指针对应的数据项，以便识别真实的类型
                if kind == reflect.Ptr {
                    rv   = rv.Elem()
                    kind = rv.Kind()
                }
                if kind == reflect.Map {
                    ks := rv.MapKeys()
                    for _, k := range ks {
                        m[String(k.Interface())] = rv.MapIndex(k).Interface()
                    }
                } else if kind == reflect.Struct {
                    rt   := rv.Type()
                    name := ""
                    for i := 0; i < rv.NumField(); i++ {
                        // 检查json tag
                        if len(noTagCheck) == 0 || !noTagCheck[0] {
                            if name = rt.Field(i).Tag.Get("json"); name == "" {
                                name = rt.Field(i).Name
                            }
                        }
                        m[name] = rv.Field(i).Interface()
                    }
                } else {
                    return nil
                }
        }
        return m
    }
}
