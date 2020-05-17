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
	sequenceMax   = uint32(1000000)                        // Sequence max.
	randomStrBase = "0123456789abcdefghijklmnopqrstuvwxyz" // 36
	macAddrStr    = "0000000"                              // MAC addresses hash result in 7 bytes.
	processIdStr  = "0000"                                 // Process id in 4 bytes.
)

// init initializes several fixed local variable.
func init() {
	// MAC addresses hash result in 7 bytes.
	macs, _ := gipv4.MacArray()
	if len(macs) > 0 {
		var macAddrBytes []byte
		for _, mac := range macs {
			macAddrBytes = append(macAddrBytes, []byte(mac)...)
		}
		b := []byte{'0', '0', '0', '0', '0', '0', '0'}
		s := strconv.FormatUint(uint64(ghash.DJBHash(macAddrBytes)), 36)
		copy(b[7-len(s):], s)
		macAddrStr = string(b)
	}
	// Process id in 4 bytes.
	{
		b := []byte{'0', '0', '0', '0'}
		s := strconv.FormatInt(int64(os.Getpid()), 36)
		copy(b[4-len(s):], s)
		processIdStr = string(b)
	}
}

// S creates and returns an unique string in 36 bytes that meets most common usages
// without strict UUID algorithm. It returns an unique string using default unique
// algorithm if no <data> is given.
//
// The specified <data> can be no more than 3 count. No matter how long each of the <data>
// size is, each of them will be hashed into 7 bytes as part of the result. If given
// <data> count is less than 3, the leftover size of the result bytes will be token by
// random string.
//
// The returned string is composed with:
// 1. Default:    MAC(7) + PID(4) + Sequence(4) + TimestampNano(12) + RandomString(9)
// 2. CustomData: Data...(7 - 21) + TimestampNano(12) + RandomString(3 - 17)
//
// Note that the returned length is fixed to 36 bytes for performance purpose.
func S(data ...[]byte) string {
	var (
		b       = make([]byte, 36)
		nanoStr = strconv.FormatInt(time.Now().UnixNano(), 36)
	)
	if len(data) == 0 {
		copy(b, macAddrStr)
		copy(b[7:], processIdStr)
		copy(b[11:], getSequence())
		copy(b[15:], nanoStr)
		copy(b[27:], getRandomStr(9))
	} else if len(data) <= 3 {
		n := 0
		for i, v := range data {
			copy(b[i*7:], getDataHashStr(v))
			n += 7
		}
		copy(b[n:], nanoStr)
		copy(b[n+12:], getRandomStr(36-n-12))
	} else {
		panic("data count too long, no more than 3")
	}
	return gconv.UnsafeBytesToStr(b)
}

// getSequence increases and returns the sequence string in 4 bytes.
// The sequence is less than 1000000.
func getSequence() string {
	b := []byte{'0', '0', '0', '0'}
	s := strconv.FormatUint(uint64(sequence.Add(1)%sequenceMax), 36)
	copy(b[4-len(s):], s)
	return gconv.UnsafeBytesToStr(b)
}

// getRandomStr randomly picks and returns <n> count of chars from randomStrBase.
func getRandomStr(n int) string {
	if n <= 0 {
		return ""
	}
	var (
		b           = make([]byte, n)
		numberBytes = grand.B(n)
	)
	for i := range b {
		b[i] = randomStrBase[numberBytes[i]%36]
	}
	return gconv.UnsafeBytesToStr(b)
}

// getDataHashStr creates and returns hash bytes in 7 bytes with given data bytes.
func getDataHashStr(data []byte) string {
	b := []byte{'0', '0', '0', '0', '0', '0', '0'}
	s := strconv.FormatUint(uint64(ghash.DJBHash(data)), 36)
	copy(b[7-len(s):], s)
	return gconv.UnsafeBytesToStr(b)
}
