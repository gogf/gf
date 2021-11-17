// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gbinary

// NOTE: THIS IS AN EXPERIMENTAL FEATURE!

// Bit Binary bit (0 | 1)
type Bit int8

// EncodeBits does encode bits return bits Default coding
func EncodeBits(bits []Bit, i int, l int) []Bit {
	return EncodeBitsWithUint(bits, uint(i), l)
}

// EncodeBitsWithUint . Merge ui bitwise into the bits array and occupy the length bits
// (Note: binary 0 | 1 digits are stored in the uis array)
func EncodeBitsWithUint(bits []Bit, ui uint, l int) []Bit {
	a := make([]Bit, l)
	for i := l - 1; i >= 0; i-- {
		a[i] = Bit(ui & 1)
		ui >>= 1
	}
	if bits != nil {
		return append(bits, a...)
	}
	return a
}

// EncodeBitsToBytes . does encode bits to bytes
// Convert bits to [] byte, encode from left to right, and add less than 1 byte from 0 to the end.
func EncodeBitsToBytes(bits []Bit) []byte {
	if len(bits)%8 != 0 {
		for i := 0; i < len(bits)%8; i++ {
			bits = append(bits, 0)
		}
	}
	b := make([]byte, 0)
	for i := 0; i < len(bits); i += 8 {
		b = append(b, byte(DecodeBitsToUint(bits[i:i+8])))
	}
	return b
}

// DecodeBits .does decode bits to int
// Resolve to int
func DecodeBits(bits []Bit) int {
	v := 0
	for _, i := range bits {
		v = v<<1 | int(i)
	}
	return v
}

// DecodeBitsToUint .Resolve to uint
func DecodeBitsToUint(bits []Bit) uint {
	v := uint(0)
	for _, i := range bits {
		v = v<<1 | uint(i)
	}
	return v
}

// DecodeBytesToBits .Parsing [] byte into character array [] uint8
func DecodeBytesToBits(bs []byte) []Bit {
	bits := make([]Bit, 0)
	for _, b := range bs {
		bits = EncodeBitsWithUint(bits, uint(b), 8)
	}
	return bits
}
