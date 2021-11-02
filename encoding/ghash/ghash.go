// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package ghash provides some classic hash functions(uint32/uint64) in go.
package ghash

// BKDRHash implements the classic BKDR hash algorithm for 32 bits.
func BKDRHash(str []byte) uint32 {
	var seed uint32 = 131 // 31 131 1313 13131 131313 etc..
	var hash uint32 = 0
	for i := 0; i < len(str); i++ {
		hash = hash*seed + uint32(str[i])
	}
	return hash
}

// BKDRHash64 implements the classic BKDR hash algorithm for 64 bits.
func BKDRHash64(str []byte) uint64 {
	var seed uint64 = 131 // 31 131 1313 13131 131313 etc..
	var hash uint64 = 0
	for i := 0; i < len(str); i++ {
		hash = hash*seed + uint64(str[i])
	}
	return hash
}

// SDBMHash implements the classic SDBM hash algorithm for 32 bits.
func SDBMHash(str []byte) uint32 {
	var hash uint32 = 0
	for i := 0; i < len(str); i++ {
		// equivalent to: hash = 65599*hash + uint32(str[i]);
		hash = uint32(str[i]) + (hash << 6) + (hash << 16) - hash
	}
	return hash
}

// SDBMHash64 implements the classic SDBM hash algorithm for 64 bits.
func SDBMHash64(str []byte) uint64 {
	var hash uint64 = 0
	for i := 0; i < len(str); i++ {
		// equivalent to: hash = 65599*hash + uint32(str[i])
		hash = uint64(str[i]) + (hash << 6) + (hash << 16) - hash
	}
	return hash
}

// RSHash implements the classic RS hash algorithm for 32 bits.
func RSHash(str []byte) uint32 {
	var b uint32 = 378551
	var a uint32 = 63689
	var hash uint32 = 0
	for i := 0; i < len(str); i++ {
		hash = hash*a + uint32(str[i])
		a *= b
	}
	return hash
}

// RSHash64 implements the classic RS hash algorithm for 64 bits.
func RSHash64(str []byte) uint64 {
	var b uint64 = 378551
	var a uint64 = 63689
	var hash uint64 = 0
	for i := 0; i < len(str); i++ {
		hash = hash*a + uint64(str[i])
		a *= b
	}
	return hash
}

// JSHash implements the classic JS hash algorithm for 32 bits.
func JSHash(str []byte) uint32 {
	var hash uint32 = 1315423911
	for i := 0; i < len(str); i++ {
		hash ^= (hash << 5) + uint32(str[i]) + (hash >> 2)
	}
	return hash
}

// JSHash64 implements the classic JS hash algorithm for 64 bits.
func JSHash64(str []byte) uint64 {
	var hash uint64 = 1315423911
	for i := 0; i < len(str); i++ {
		hash ^= (hash << 5) + uint64(str[i]) + (hash >> 2)
	}
	return hash
}

// PJWHash implements the classic PJW hash algorithm for 32 bits.
func PJWHash(str []byte) uint32 {
	var BitsInUnignedInt uint32 = 4 * 8
	var ThreeQuarters uint32 = (BitsInUnignedInt * 3) / 4
	var OneEighth uint32 = BitsInUnignedInt / 8
	var HighBits uint32 = (0xFFFFFFFF) << (BitsInUnignedInt - OneEighth)
	var hash uint32 = 0
	var test uint32 = 0
	for i := 0; i < len(str); i++ {
		hash = (hash << OneEighth) + uint32(str[i])
		if test = hash & HighBits; test != 0 {
			hash = (hash ^ (test >> ThreeQuarters)) & (^HighBits + 1)
		}
	}
	return hash
}

// PJWHash64 implements the classic PJW hash algorithm for 64 bits.
func PJWHash64(str []byte) uint64 {
	var BitsInUnignedInt uint64 = 4 * 8
	var ThreeQuarters uint64 = (BitsInUnignedInt * 3) / 4
	var OneEighth uint64 = BitsInUnignedInt / 8
	var HighBits uint64 = (0xFFFFFFFFFFFFFFFF) << (BitsInUnignedInt - OneEighth)
	var hash uint64 = 0
	var test uint64 = 0
	for i := 0; i < len(str); i++ {
		hash = (hash << OneEighth) + uint64(str[i])
		if test = hash & HighBits; test != 0 {
			hash = (hash ^ (test >> ThreeQuarters)) & (^HighBits + 1)
		}
	}
	return hash
}

// ELFHash implements the classic ELF hash algorithm for 32 bits.
func ELFHash(str []byte) uint32 {
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

// ELFHash64 implements the classic ELF hash algorithm for 64 bits.
func ELFHash64(str []byte) uint64 {
	var hash uint64 = 0
	var x uint64 = 0
	for i := 0; i < len(str); i++ {
		hash = (hash << 4) + uint64(str[i])
		if x = hash & 0xF000000000000000; x != 0 {
			hash ^= x >> 24
			hash &= ^x + 1
		}
	}
	return hash
}

// DJBHash implements the classic DJB hash algorithm for 32 bits.
func DJBHash(str []byte) uint32 {
	var hash uint32 = 5381
	for i := 0; i < len(str); i++ {
		hash += (hash << 5) + uint32(str[i])
	}
	return hash
}

// DJBHash64 implements the classic DJB hash algorithm for 64 bits.
func DJBHash64(str []byte) uint64 {
	var hash uint64 = 5381
	for i := 0; i < len(str); i++ {
		hash += (hash << 5) + uint64(str[i])
	}
	return hash
}

// APHash implements the classic AP hash algorithm for 32 bits.
func APHash(str []byte) uint32 {
	var hash uint32 = 0
	for i := 0; i < len(str); i++ {
		if (i & 1) == 0 {
			hash ^= (hash << 7) ^ uint32(str[i]) ^ (hash >> 3)
		} else {
			hash ^= ^((hash << 11) ^ uint32(str[i]) ^ (hash >> 5)) + 1
		}
	}
	return hash
}

// APHash64 implements the classic AP hash algorithm for 64 bits.
func APHash64(str []byte) uint64 {
	var hash uint64 = 0
	for i := 0; i < len(str); i++ {
		if (i & 1) == 0 {
			hash ^= (hash << 7) ^ uint64(str[i]) ^ (hash >> 3)
		} else {
			hash ^= ^((hash << 11) ^ uint64(str[i]) ^ (hash >> 5)) + 1
		}
	}
	return hash
}
