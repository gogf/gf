// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// 并发安全的数组.
// 底层使用通用的interface{}类型，从性能上考虑，类似于gmap那样可以为每种类型都定义一个array.
package garray
