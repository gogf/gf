// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package guid provides simple and high performance unique id generation functionality.
package guid

import (
	"os"
	"strconv"
	"time"

	"github.com/gogf/gf/v2/container/gtype"
	"github.com/gogf/gf/v2/encoding/ghash"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/gipv4"
	"github.com/gogf/gf/v2/util/grand"
)

const (
	sequenceMax   = uint32(46655)                          // Sequence max("zzz").
	randomStrBase = "0123456789abcdefghijklmnopqrstuvwxyz" // Random chars string(36 bytes).
)

var (
	sequence     gtype.Uint32 // Sequence for unique purpose of the current process.
	macAddrStr   = "0000000"  // Hash result of MAC addresses in 7 bytes.
	processIdStr = "0000"     // Process id in 4 bytes.
)

// init initializes several fixed local variable.
func init() {
	// Hash result of MAC addresses in 7 bytes.
	macs, _ := gipv4.GetMacArray()
	if len(macs) > 0 {
		var macAddrBytes []byte
		for _, mac := range macs {
			macAddrBytes = append(macAddrBytes, []byte(mac)...)
		}
		b := []byte{'0', '0', '0', '0', '0', '0', '0'}
		s := strconv.FormatUint(uint64(ghash.SDBM(macAddrBytes)), 36)
		copy(b, s)
		macAddrStr = string(b)
	}
	// Process id in 4 bytes.
	{
		b := []byte{'0', '0', '0', '0'}
		s := strconv.FormatInt(int64(os.Getpid()), 36)
		copy(b, s)
		processIdStr = string(b)
	}
}

// S creates and returns a global unique string in 32 bytes that meets most common
// usages without strict UUID algorithm. It returns a unique string using default
// unique algorithm if no `data` is given.
//
// The specified `data` can be no more than 2 parts. No matter how long each of the
// `data` size is, each of them will be hashed into 7 bytes as part of the result.
// If given `data` parts is less than 2, the leftover size of the result bytes will
// be token by random string.
//
// The returned string is composed with:
// 1. Default:    MACHash(7) + PID(4) + TimestampNano(12) + Sequence(3) + RandomString(6)
// 2. CustomData: DataHash(7/14) + TimestampNano(12) + Sequence(3) + RandomString(3/10)
//
// Note that：
//  1. The returned length is fixed to 32 bytes for performance purpose.
//  2. The custom parameter `data` composed should have unique attribute in your
//     business scenario.
func S(data ...[]byte) string {
	var (
		b       = make([]byte, 32)
		nanoStr = strconv.FormatInt(time.Now().UnixNano(), 36)
	)
	if len(data) == 0 {
		copy(b, macAddrStr)
		copy(b[7:], processIdStr)
		copy(b[11:], nanoStr)
		copy(b[23:], getSequence())
		copy(b[26:], getRandomStr(6))
	} else if len(data) <= 2 {
		n := 0
		for i, v := range data {
			// Ignore empty data item bytes.
			if len(v) > 0 {
				copy(b[i*7:], getDataHashStr(v))
				n += 7
			}
		}
		copy(b[n:], nanoStr)
		copy(b[n+12:], getSequence())
		copy(b[n+12+3:], getRandomStr(32-n-12-3))
	} else {
		panic(gerror.NewCode(
			gcode.CodeInvalidParameter,
			"too many data parts, it should be no more than 2 parts",
		))
	}
	return string(b)
}

// getSequence increases and returns the sequence string in 3 bytes.
// The sequence is less than "zzz"(46655).
func getSequence() []byte {
	b := []byte{'0', '0', '0'}
	s := strconv.FormatUint(uint64(sequence.Add(1)%sequenceMax), 36)
	copy(b, s)
	return b
}

// getRandomStr randomly picks and returns `n` count of chars from randomStrBase.
func getRandomStr(n int) []byte {
	if n <= 0 {
		return []byte{}
	}
	var (
		b           = make([]byte, n)
		numberBytes = grand.B(n)
	)
	for i := range b {
		b[i] = randomStrBase[numberBytes[i]%36]
	}
	return b
}

// getDataHashStr creates and returns hash bytes in 7 bytes with given data bytes.
func getDataHashStr(data []byte) []byte {
	b := []byte{'0', '0', '0', '0', '0', '0', '0'}
	s := strconv.FormatUint(uint64(ghash.SDBM(data)), 36)
	copy(b, s)
	return b
}
