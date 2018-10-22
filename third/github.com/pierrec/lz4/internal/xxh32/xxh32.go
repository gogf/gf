// Package xxh32 implements the very fast XXH hashing algorithm (32 bits version).
// (https://github.com/Cyan4973/XXH/)
package xxh32

import (
	"encoding/binary"
	"hash"
)

const (
	prime32_1 uint32 = 2654435761
	prime32_2 uint32 = 2246822519
	prime32_3 uint32 = 3266489917
	prime32_4 uint32 = 668265263
	prime32_5 uint32 = 374761393

	prime32_1plus2 uint32 = 606290984
	prime32_minus1 uint32 = 1640531535
)

var _ hash.Hash32 = (*XXH)(nil)

// XXH represents an xxhash32 object.
type XXH struct {
	seed     uint32
	v1       uint32
	v2       uint32
	v3       uint32
	v4       uint32
	totalLen uint64
	buf      [16]byte
	bufused  int
}

// Init sets the seed and Reset().
func (xxh *XXH) Init(seed uint32) {
	xxh.seed = seed
	xxh.Reset()
}

// Sum appends the current hash to b and returns the resulting slice.
// It does not change the underlying hash state.
func (xxh XXH) Sum(b []byte) []byte {
	h32 := xxh.Sum32()
	return append(b, byte(h32), byte(h32>>8), byte(h32>>16), byte(h32>>24))
}

// Reset resets the Hash to its initial state.
func (xxh *XXH) Reset() {
	seed := xxh.seed
	xxh.v1 = seed + prime32_1plus2
	xxh.v2 = seed + prime32_2
	xxh.v3 = seed
	xxh.v4 = seed - prime32_1
	xxh.totalLen = 0
	xxh.bufused = 0
}

// Size returns the number of bytes returned by Sum().
func (xxh *XXH) Size() int {
	return 4
}

// BlockSize gives the minimum number of bytes accepted by Write().
func (xxh *XXH) BlockSize() int {
	return 1
}

// Write adds input bytes to the Hash.
// It never returns an error.
func (xxh *XXH) Write(input []byte) (int, error) {
	n := len(input)
	m := xxh.bufused

	xxh.totalLen += uint64(n)

	r := len(xxh.buf) - m
	if n < r {
		copy(xxh.buf[m:], input)
		xxh.bufused += len(input)
		return n, nil
	}

	// Causes compiler to work directly from registers instead of stack:
	v1, v2, v3, v4 := xxh.v1, xxh.v2, xxh.v3, xxh.v4
	p := 0
	if m > 0 {
		// some data left from previous update
		copy(xxh.buf[xxh.bufused:], input[:r])
		xxh.bufused += len(input) - r

		// fast rotl(13)
		buf := xxh.buf[:16] // BCE hint.
		xxh.v1 = rol13(v1+binary.LittleEndian.Uint32(buf[:])*prime32_2) * prime32_1
		xxh.v2 = rol13(v2+binary.LittleEndian.Uint32(buf[4:])*prime32_2) * prime32_1
		xxh.v3 = rol13(v3+binary.LittleEndian.Uint32(buf[8:])*prime32_2) * prime32_1
		xxh.v4 = rol13(v4+binary.LittleEndian.Uint32(buf[12:])*prime32_2) * prime32_1
		p = r
		xxh.bufused = 0
	}

	for n := n - 16; p <= n; p += 16 {
		sub := input[p:][:16] //BCE hint for compiler
		v1 = rol13(v1+binary.LittleEndian.Uint32(sub[:])*prime32_2) * prime32_1
		v2 = rol13(v2+binary.LittleEndian.Uint32(sub[4:])*prime32_2) * prime32_1
		v3 = rol13(v3+binary.LittleEndian.Uint32(sub[8:])*prime32_2) * prime32_1
		v4 = rol13(v4+binary.LittleEndian.Uint32(sub[12:])*prime32_2) * prime32_1
	}
	xxh.v1, xxh.v2, xxh.v3, xxh.v4 = v1, v2, v3, v4

	copy(xxh.buf[xxh.bufused:], input[p:])
	xxh.bufused += len(input) - p

	return n, nil
}

// Sum32 returns the 32 bits Hash value.
func (xxh *XXH) Sum32() uint32 {
	h32 := uint32(xxh.totalLen)
	if h32 >= 16 {
		h32 += rol1(xxh.v1) + rol7(xxh.v2) + rol12(xxh.v3) + rol18(xxh.v4)
	} else {
		h32 += xxh.seed + prime32_5
	}

	p := 0
	n := xxh.bufused
	buf := xxh.buf
	for n := n - 4; p <= n; p += 4 {
		h32 += binary.LittleEndian.Uint32(buf[p:p+4]) * prime32_3
		h32 = rol17(h32) * prime32_4
	}
	for ; p < n; p++ {
		h32 += uint32(buf[p]) * prime32_5
		h32 = rol11(h32) * prime32_1
	}

	h32 ^= h32 >> 15
	h32 *= prime32_2
	h32 ^= h32 >> 13
	h32 *= prime32_3
	h32 ^= h32 >> 16

	return h32
}

// Checksum returns the 32bits Hash value.
func Checksum(input []byte, seed uint32) uint32 {
	n := len(input)
	h32 := uint32(n)

	if n < 16 {
		h32 += seed + prime32_5
	} else {
		v1 := seed + prime32_1 + prime32_2
		v2 := seed + prime32_2
		v3 := seed
		v4 := seed - prime32_1
		p := 0
		for n := n - 16; p <= n; p += 16 {
			sub := input[p:][:16] //BCE hint for compiler
			v1 = rol13(v1+binary.LittleEndian.Uint32(sub[:])*prime32_2) * prime32_1
			v2 = rol13(v2+binary.LittleEndian.Uint32(sub[4:])*prime32_2) * prime32_1
			v3 = rol13(v3+binary.LittleEndian.Uint32(sub[8:])*prime32_2) * prime32_1
			v4 = rol13(v4+binary.LittleEndian.Uint32(sub[12:])*prime32_2) * prime32_1
		}
		input = input[p:]
		n -= p
		h32 += rol1(v1) + rol7(v2) + rol12(v3) + rol18(v4)
	}

	p := 0
	for n := n - 4; p <= n; p += 4 {
		h32 += binary.LittleEndian.Uint32(input[p:p+4]) * prime32_3
		h32 = rol17(h32) * prime32_4
	}
	for p < n {
		h32 += uint32(input[p]) * prime32_5
		h32 = rol11(h32) * prime32_1
		p++
	}

	h32 ^= h32 >> 15
	h32 *= prime32_2
	h32 ^= h32 >> 13
	h32 *= prime32_3
	h32 ^= h32 >> 16

	return h32
}

func rol1(u uint32) uint32 {
	return u<<1 | u>>31
}

func rol7(u uint32) uint32 {
	return u<<7 | u>>25
}

func rol11(u uint32) uint32 {
	return u<<11 | u>>21
}

func rol12(u uint32) uint32 {
	return u<<12 | u>>20
}

func rol13(u uint32) uint32 {
	return u<<13 | u>>19
}

func rol17(u uint32) uint32 {
	return u<<17 | u>>15
}

func rol18(u uint32) uint32 {
	return u<<18 | u>>14
}

// Uint32 hashes x with the given seed.
func Uint32(x, seed uint32) uint32 {
	h := seed + prime32_5 + 4 + x*prime32_3
	h = rol17(h) * prime32_4
	h ^= h >> 15
	h *= prime32_2
	h ^= h >> 13
	h *= prime32_3
	h ^= h >> 16
	return h
}
