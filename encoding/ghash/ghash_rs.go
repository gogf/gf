// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghash

// RS implements the classic RS hash algorithm for 32 bits.
func RS(str []byte) uint32 {
	var (
		b    uint32 = 378551
		a    uint32 = 63689
		hash uint32 = 0
	)
	for i := 0; i < len(str); i++ {
		hash = hash*a + uint32(str[i])
		a *= b
	}
	return hash
}

// RS64 implements the classic RS hash algorithm for 64 bits.
func RS64(str []byte) uint64 {
	var (
		b    uint64 = 378551
		a    uint64 = 63689
		hash uint64 = 0
	)
	for i := 0; i < len(str); i++ {
		hash = hash*a + uint64(str[i])
		a *= b
	}
	return hash
}
