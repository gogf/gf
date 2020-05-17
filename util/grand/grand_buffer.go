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
	buffer := make([]byte, 1024)
	for {
		if n, err := rand.Read(buffer); err != nil {
			panic(err)
		} else {
			for i := 0; i < n-4; i += 4 {
				b := make([]byte, 4)
				copy(b, buffer[i:i+4])
				bufferChan <- b
			}
		}
	}
}
