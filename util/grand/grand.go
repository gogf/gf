// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package grand provides high performance random bytes/number/string generation functionality.
package grand

import (
	"encoding/binary"
	"time"
)

var (
	letters    = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ" // 52
	symbols    = "!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~"                   // 32
	digits     = "0123456789"                                           // 10
	characters = letters + digits + symbols                             // 94
)

// Intn returns an int number which is between 0 and max: [0, max).
//
// Note that:
// 1. The `max` can only be greater than 0, or else it returns `max` directly;
// 2. The result is greater than or equal to 0, but less than `max`;
// 3. The result number is 32bit and less than math.MaxUint32.
func Intn(max int) int {
	if max <= 0 {
		return max
	}
	n := int(binary.LittleEndian.Uint32(<-bufferChan)) % max
	if (max > 0 && n < 0) || (max < 0 && n > 0) {
		return -n
	}
	return n
}

// B retrieves and returns random bytes of given length `n`.
func B(n int) []byte {
	if n <= 0 {
		return nil
	}
	i := 0
	b := make([]byte, n)
	for {
		copy(b[i:], <-bufferChan)
		i += 4
		if i >= n {
			break
		}
	}
	return b
}

// N returns a random int between min and max: [min, max].
// The `min` and `max` also support negative numbers.
func N(min, max int) int {
	if min >= max {
		return min
	}
	if min >= 0 {
		// Because Intn dose not support negative number,
		// so we should first shift the value to left,
		// then call Intn to produce the random number,
		// and finally shift the result back to right.
		return Intn(max-(min-0)+1) + (min - 0)
	}
	if min < 0 {
		// Because Intn dose not support negative number,
		// so we should first shift the value to right,
		// then call Intn to produce the random number,
		// and finally shift the result back to left.
		return Intn(max+(0-min)+1) - (0 - min)
	}
	return 0
}

// S returns a random string which contains digits and letters, and its length is `n`.
// The optional parameter `symbols` specifies whether the result could contain symbols,
// which is false in default.
func S(n int, symbols ...bool) string {
	if n <= 0 {
		return ""
	}
	var (
		b           = make([]byte, n)
		numberBytes = B(n)
	)
	for i := range b {
		if len(symbols) > 0 && symbols[0] {
			b[i] = characters[numberBytes[i]%94]
		} else {
			b[i] = characters[numberBytes[i]%62]
		}
	}
	return string(b)
}

// D returns a random time.Duration between min and max: [min, max].
func D(min, max time.Duration) time.Duration {
	multiple := int64(1)
	if min != 0 {
		for min%10 == 0 {
			multiple *= 10
			min /= 10
			max /= 10
		}
	}
	n := int64(N(int(min), int(max)))
	return time.Duration(n * multiple)
}

// Str randomly picks and returns `n` count of chars from given string `s`.
// It also supports unicode string like Chinese/Russian/Japanese, etc.
func Str(s string, n int) string {
	if n <= 0 {
		return ""
	}
	var (
		b     = make([]rune, n)
		runes = []rune(s)
	)
	if len(runes) <= 255 {
		numberBytes := B(n)
		for i := range b {
			b[i] = runes[int(numberBytes[i])%len(runes)]
		}
	} else {
		for i := range b {
			b[i] = runes[Intn(len(runes))]
		}
	}
	return string(b)
}

// Digits returns a random string which contains only digits, and its length is `n`.
func Digits(n int) string {
	if n <= 0 {
		return ""
	}
	var (
		b           = make([]byte, n)
		numberBytes = B(n)
	)
	for i := range b {
		b[i] = digits[numberBytes[i]%10]
	}
	return string(b)
}

// Letters returns a random string which contains only letters, and its length is `n`.
func Letters(n int) string {
	if n <= 0 {
		return ""
	}
	var (
		b           = make([]byte, n)
		numberBytes = B(n)
	)
	for i := range b {
		b[i] = letters[numberBytes[i]%52]
	}
	return string(b)
}

// Symbols returns a random string which contains only symbols, and its length is `n`.
func Symbols(n int) string {
	if n <= 0 {
		return ""
	}
	var (
		b           = make([]byte, n)
		numberBytes = B(n)
	)
	for i := range b {
		b[i] = symbols[numberBytes[i]%32]
	}
	return string(b)
}

// Perm returns, as a slice of n int numbers, a pseudo-random permutation of the integers [0,n).
// TODO performance improving for large slice producing.
func Perm(n int) []int {
	m := make([]int, n)
	for i := 0; i < n; i++ {
		j := Intn(i + 1)
		m[i] = m[j]
		m[j] = i
	}
	return m
}

// Meet randomly calculate whether the given probability `num`/`total` is met.
func Meet(num, total int) bool {
	return Intn(total) < num
}

// MeetProb randomly calculate whether the given probability is met.
func MeetProb(prob float32) bool {
	return Intn(1e7) < int(prob*1e7)
}
