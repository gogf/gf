// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghash

// SDBM implements the classic SDBM hash algorithm for 32 bits.
func SDBM(str []byte) uint32 {
	var hash uint32 = 0
	for i := 0; i < len(str); i++ {
		// equivalent to: hash = 65599*hash + uint32(str[i]);
		hash = uint32(str[i]) + (hash << 6) + (hash << 16) - hash
	}
	return hash
}

// SDBM64 implements the classic SDBM hash algorithm for 64 bits.
func SDBM64(str []byte) uint64 {
	var hash uint64 = 0
	for i := 0; i < len(str); i++ {
		// equivalent to: hash = 65599*hash + uint32(str[i])
		hash = uint64(str[i]) + (hash << 6) + (hash << 16) - hash
	}
	return hash
}
