// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// 其他工具包
package gutil

import (
    "reflect"
    "gitee.com/johng/gf/g/util/gstr"
    "gitee.com/johng/gf/g/util/gconv"
)

// 将map键值对映射到对应的struct对象属性上，需要注意：
// 1、第二个参数为struct对象指针；
// 2、struct对象的公开属性才能被映射赋值；
// 3、map中的键名可以为小写，映射转换时会自动将键名首字母转为大写做匹配映射，如果无法匹配则忽略；
func MapToStruct(m map[string]interface{}, o interface{}) error {
    for k, v := range m {
        _MapToStructSetField(o, k, v)
    }
    return nil
}
func _MapToStructSetField(obj interface{}, name string, value interface{}) {
    structValue      := reflect.ValueOf(obj).Elem()
    structFieldValue := structValue.FieldByName(gstr.UcFirst(name))
    // 键名与对象属性匹配检测
    if !structFieldValue.IsValid() {
        //return fmt.Errorf("No such field: %s in obj", name)
        return
    }
    // CanSet的属性必须为公开属性(首字母大写)
    if !structFieldValue.CanSet() {
        //return fmt.Errorf("Cannot set %s field value", name)
        return
    }
    // 必须将value转换为struct属性的数据类型，这里必须用到gconv包
    structFieldValue.Set(reflect.ValueOf(gconv.Convert(value, structFieldValue.Type().String())))
}