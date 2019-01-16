// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// Package gmap provides kinds of concurrent-safe(alternative) maps.
//
// 并发安全的哈希MAP.
package gmap

// 默认的Map对象其实就是InterfaceInterfaceMap的别名。
// 注意
// 1、这个Map是所有并发安全Map中效率最低的，如果对效率要求比较高的场合，请合理选择对应数据类型的Map；
// 2、这个Map的优点是使用简便，由于键值都是interface{}类型，因此对键值的数据类型要求不高；
// 3、底层实现比较类似于sync.Map；
type Map = InterfaceInterfaceMap

func New(safe...bool) *Map {
    return NewInterfaceInterfaceMap(safe...)
}