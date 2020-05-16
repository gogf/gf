// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package grand

import (
	"crypto/rand"
)

const (
	// Buffer size for uint32 random number.
	gBUFFER_SIZE = 10000
)

var (
	// bufferChan is the buffer for random bytes,
	// every item storing 4 bytes.
	bufferChan = make(chan []byte, gBUFFER_SIZE)
)

func init() {
	go asyncProducingRandomBufferBytesLoop()
}

// asyncProducingRandomBufferBytes is a named goroutine, which uses a asynchronous goroutine
// to produce the random bytes, and a buffer chan to store the random bytes.
// So it has high performance to generate random numbers.
func asyncProducingRandomBufferBytesLoop() {
	var (
		step   = 0
		buffer = make([]byte, 1024)
	)
	for {
		if n, err := rand.Read(buffer); err != nil {
			panic(err)
		} else {
			for i := 0; i < n-4; {
				bufferChan <- buffer[i : i+4]
				i++
			}
			// Reuse the rand buffer.
			for i := 0; i < n; i++ {
				step = int(buffer[0]) % 10
				if step != 0 {
					break
				}
			}
			// The step cannot be 0,
			// as it will produce the same random number as previous.
			if step == 0 {
				step = 2
			}
			for i := 0; i < n-4; {
				bufferChan <- buffer[i : i+4]
				i += step
			}
		}
	}
}
