// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghash

// ELF implements the classic ELF hash algorithm for 32 bits.
func ELF(str []byte) uint32 {
	var hash uint32 = 0
	var x uint32 = 0
	for i := 0; i < len(str); i++ {
		hash = (hash << 4) + uint32(str[i])
		if x = hash & 0xF0000000; x != 0 {
			hash ^= x >> 24
			hash &= ^x + 1
		}
	}
	return hash
}

// ELF64 implements the classic ELF hash algorithm for 64 bits.
func ELF64(str []byte) uint64 {
	var (
		hash uint64 = 0
		x    uint64 = 0
	)
	for i := 0; i < len(str); i++ {
		hash = (hash << 4) + uint64(str[i])
		if x = hash & 0xF000000000000000; x != 0 {
			hash ^= x >> 24
			hash &= ^x + 1
		}
	}
	return hash
}
