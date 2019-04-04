// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package grand provides high performance API for random functionality.
// 
// 随机数管理.
package grand

var (
    letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
    digits  = []rune("0123456789")
)

// 随机计算是否满足给定的概率(分子/分母)
func Meet(num, total int) bool {
    return Intn(total) < num
}

// 随机计算是否满足给定的概率(float32)
func MeetProb(prob float32) bool {
    return Intn(1e7) < int(prob*1e7)
}

// Rand 别名, 返回: [min, max]
func N (min, max int) int {
   return Rand(min, max)
}

// 获得一个 min, max 之间的随机数: [min, max]
func Rand (min, max int) int {
    if min >= max {
        return min
    }
    if min >= 0 {
        // 数值往左平移，再使用底层随机方法获得随机数，随后将结果数值往右平移
        return Intn(max - (min - 0) + 1) + (min - 0)
    }
    if min < 0 {
        // 数值往右平移，再使用底层随机方法获得随机数，随后将结果数值往左平移
        return Intn(max + (0 - min) + 1) - (0 - min)
    }
    return 0
}

// RandStr 别名
func Str(n int) string {
    return RandStr(n)
}

// 获得指定长度的随机字符串(可能包含数字和字母)
func RandStr(n int) string {
    b := make([]rune, n)
    for i := range b {
        if Intn(2) == 1 {
            b[i] = digits[Intn(10)]
        } else {
            b[i] = letters[Intn(52)]
        }
    }
    return string(b)
}

// RandDigits 别名
func Digits(n int) string {
    return RandDigits(n)
}

// 获得指定长度的随机数字字符串
func RandDigits(n int) string {
    b := make([]rune, n)
    for i := range b {
        b[i] = digits[Intn(10)]
    }
    return string(b)
}

// RandLetters 别名
func Letters(n int) string {
    return RandLetters(n)
}

// 获得指定长度的随机字母字符串
func RandLetters(n int) string {
    b := make([]rune, n)
    for i := range b {
        b[i] = letters[Intn(52)]
    }
    return string(b)
}

// Perm returns, as a slice of n ints, a pseudo-random permutation of the integers [0,n).
//
// 返回[0, n)的随机数组成的slice。
func Perm(n int) []int {
    m := make([]int, n)
    for i := 0; i < n; i++ {
        j := Intn(i + 1)
        m[i] = m[j]
        m[j] = i
    }
    return m
}
