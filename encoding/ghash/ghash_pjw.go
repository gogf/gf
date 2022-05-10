// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghash

// PJW implements the classic PJW hash algorithm for 32 bits.
func PJW(str []byte) uint32 {
	var (
		BitsInUnsignedInt uint32 = 4 * 8
		ThreeQuarters     uint32 = (BitsInUnsignedInt * 3) / 4
		OneEighth         uint32 = BitsInUnsignedInt / 8
		HighBits          uint32 = (0xFFFFFFFF) << (BitsInUnsignedInt - OneEighth)
		hash              uint32 = 0
		test              uint32 = 0
	)
	for i := 0; i < len(str); i++ {
		hash = (hash << OneEighth) + uint32(str[i])
		if test = hash & HighBits; test != 0 {
			hash = (hash ^ (test >> ThreeQuarters)) & (^HighBits + 1)
		}
	}
	return hash
}

// PJW64 implements the classic PJW hash algorithm for 64 bits.
func PJW64(str []byte) uint64 {
	var (
		BitsInUnsignedInt uint64 = 4 * 8
		ThreeQuarters     uint64 = (BitsInUnsignedInt * 3) / 4
		OneEighth         uint64 = BitsInUnsignedInt / 8
		HighBits          uint64 = (0xFFFFFFFFFFFFFFFF) << (BitsInUnsignedInt - OneEighth)
		hash              uint64 = 0
		test              uint64 = 0
	)
	for i := 0; i < len(str); i++ {
		hash = (hash << OneEighth) + uint64(str[i])
		if test = hash & HighBits; test != 0 {
			hash = (hash ^ (test >> ThreeQuarters)) & (^HighBits + 1)
		}
	}
	return hash
}
