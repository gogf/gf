// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package grand

import (
	"crypto/rand"
	"encoding/binary"
)

const (
	// Buffer size for uint32 random number.
	gBUFFER_SIZE = 10000
)

var (
	// Buffer chan.
	bufferChan = make(chan uint32, gBUFFER_SIZE)
)

// It uses a asynchronous goroutine to produce the random number,
// and a buffer chan to store the random number. So it has high performance
// to generate random number.
func init() {
	step := 0
	buffer := make([]byte, 1024)
	go func() {
		for {
			if n, err := rand.Read(buffer); err != nil {
				panic(err)
			} else {
				for i := 0; i < n-4; {
					bufferChan <- binary.LittleEndian.Uint32(buffer[i : i+4])
					i++
				}
				// Reuse the rand buffer.
				for i := 0; i < n; i++ {
					step = int(buffer[0]) % 10
					if step != 0 {
						break
					}
				}
				if step == 0 {
					step = 2
				}
				for i := 0; i < n-4; {
					bufferChan <- binary.BigEndian.Uint32(buffer[i : i+4])
					i += step
				}
			}
		}
	}()
}

// Intn returns a int number which is between 0 and max - [0, max).
//
// Note:
// 1. The <max> can only be geater than 0, or else it return <max> directly;
// 2. The result is greater than or equal to 0, but less than <max>;
// 3. The result number is 32bit and less than math.MaxUint32.
func Intn(max int) int {
	if max <= 0 {
		return max
	}
	n := int(<-bufferChan) % max
	if (max > 0 && n < 0) || (max < 0 && n > 0) {
		return -n
	}
	return n
}
