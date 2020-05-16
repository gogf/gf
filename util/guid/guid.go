// Copyright 2020 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package guid provides simple and high performance unique id generation functionality.
//
// PLEASE VERY NOTE:
// This package only provides unique number generation for simple, convenient and most common
// usage purpose, but does not provide strict global unique number generation. Please refer
// to UUID algorithm for global unique number generation if necessary.
package guid

import (
	"github.com/gogf/gf/container/gtype"
	"github.com/gogf/gf/encoding/ghash"
	"github.com/gogf/gf/net/gipv4"
	"github.com/gogf/gf/util/gconv"
	"github.com/gogf/gf/util/grand"
	"os"
	"strconv"
	"time"
)

var (
	sequence      gtype.Uint32                             // Sequence for unique purpose of current process.
	macAddrStr    string                                   // MAC addresses hash result in 7 bytes.
	processIdStr  string                                   // Process id in 4 bytes.
	randomStrBase = "0123456789abcdefghijklmnopqrstuvwxyz" // 36
)

func init() {
	// MAC addresses hash result in 7 bytes.
	macs, _ := gipv4.MacArray()
	if len(macs) > 0 {
		var b []byte
		for _, mac := range macs {
			b = append(b, []byte(mac)...)
		}
		macAddrStr = strconv.FormatUint(uint64(ghash.DJBHash(b)), 36)
		if n := 7 - len(macAddrStr); n > 0 {
			for i := 0; i < n; i++ {
				macAddrStr = "0" + macAddrStr
			}
		}
		if n := len(macAddrStr) - 7; n > 0 {
			macAddrStr = macAddrStr[n:]
		}
	}
	// Process id in 4 bytes.
	processIdStr = strconv.FormatInt(int64(os.Getpid()), 36)
	if n := 4 - len(processIdStr); n > 0 {
		for i := 0; i < n; i++ {
			processIdStr = "0" + processIdStr
		}
	}
	if n := len(processIdStr) - 4; n > 0 {
		processIdStr = processIdStr[n:]
	}
}

// S creates and returns an unique string in 36 bytes that meets most
// common usages without strict UUID algorithm.
// The returned string is composed with:
// MAC(7) + PID(4) + Sequence(7) + TimestampNano(12) + RandomString(6)
func S() string {
	b := make([]byte, 36)
	copy(b, macAddrStr)
	copy(b[7:], processIdStr)
	copy(b[11:], getSequence())
	copy(b[18:], strconv.FormatInt(time.Now().UnixNano(), 36))
	copy(b[30:], randomStr(randomStrBase, 6))
	return gconv.UnsafeBytesToStr(b)
}

// New creates and returns an unique string in 36 bytes using custom data.
// The returned string is composed with:
// Data...(7 - 21) + TimestampNano(12) + RandomString(3 - 17)
//
// The specified <data> can be count of 1 to 3. No matter how long each of the <data>
// size is, each of them will be hashed into 7 bytes as part of the result. If given
// <data> count is less than 3, the leftover size of the result bytes will be token by
// random string.
//
// Note that the <data> can not be empty.
func New(data ...[]byte) string {
	if len(data) == 0 {
		panic("data cannot be empty")
	}
	if len(data) > 3 {
		panic("data count too long, no more than 3")
	}
	b := make([]byte, 36)
	n := 0
	for i, v := range data {
		copy(b[i*7:], getDataHashStr(v))
		n += 7
	}
	copy(b[n:], strconv.FormatInt(time.Now().UnixNano(), 36))
	copy(b[n+12:], randomStr(randomStrBase, 36-n-12))
	return gconv.UnsafeBytesToStr(b)
}

// getSequence increases and returns the sequence string in 7 bytes.
func getSequence() string {
	b := make([]byte, 7)
	copy(b, []byte{'0', '0', '0', '0', '0', '0', '0'})
	s := strconv.FormatUint(uint64(sequence.Add(1)), 36)
	copy(b[7-len(s):], s)
	return gconv.UnsafeBytesToStr(b)
}

// Str randomly picks and returns <n> count of chars from given string <s>.
// It also supports unicode string like Chinese/Russian/Japanese, etc.
func randomStr(s string, n int) string {
	var (
		b           = make([]byte, n)
		numberBytes = grand.B(n)
	)
	for i := range b {
		b[i] = s[numberBytes[i]%36]
	}
	return gconv.UnsafeBytesToStr(b)
}

// getDataHashStr creates and returns hash bytes in 7 bytes with given data bytes.
func getDataHashStr(data []byte) string {
	b := make([]byte, 7)
	copy(b, []byte{'0', '0', '0', '0', '0', '0', '0'})
	s := strconv.FormatUint(uint64(ghash.DJBHash(data)), 36)
	copy(b[7-len(s):], s)
	return gconv.UnsafeBytesToStr(b)
}
