// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// 随机数管理
package grand

import (
    "time"
    "math/rand"
)
var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
var digits  = []rune("0123456789")

// 修改全局随机数seed
func Seed (seed...int64) {
    if len(seed) > 0 {
        rand.Seed(seed[0])
    } else {
        rand.Seed(time.Now().UnixNano())
    }
}

// 获得一个 min, max 之间的随机数(min <= x <= max)
func Rand (min, max int) int {
    if min >= max {
        return min
    }
    n := rand.Intn(max + 1)
    if n < min {
        return Rand(min, max)
    }
    return n
}

// 获得指定长度的随机字符串(可能包含数字和字母)
func RandStr(n int) string {
    b := make([]rune, n)
    for i := range b {
        if rand.Intn(2) == 1 {
            b[i] = digits[rand.Intn(10)]
        } else {
            b[i] = letters[rand.Intn(52)]
        }
    }
    return string(b)
}

// 获得指定长度的随机数字字符串
func RandDigits(n int) string {
    b := make([]rune, n)
    for i := range b {
        b[i] = digits[rand.Intn(10)]
    }
    return string(b)
}

// 获得指定长度的随机字母字符串
func RandLetters(n int) string {
    b := make([]rune, n)
    for i := range b {
        b[i] = letters[rand.Intn(52)]
    }
    return string(b)
}
