// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
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
	bufferChanSize = 10000
)

var (
	// bufferChan is the buffer for random bytes,
	// every item storing 4 bytes.
	bufferChan = make(chan []byte, bufferChanSize)
)

func init() {
	go asyncProducingRandomBufferBytesLoop()
}

// asyncProducingRandomBufferBytes is a named goroutine, which uses a asynchronous goroutine
// to produce the random bytes, and a buffer chan to store the random bytes.
// So it has high performance to generate random numbers.
func asyncProducingRandomBufferBytesLoop() {
	var step int
	for {
		buffer := make([]byte, 1024)
		if n, err := rand.Read(buffer); err != nil {
			panic(err)
		} else {
			// The random buffer from system is very expensive,
			// so fully reuse the random buffer by changing
			// the step with a different number can
			// improve the performance a lot.
			// for _, step = range []int{4, 5, 6, 7} {
			for _, step = range []int{4} {
				for i := 0; i <= n-4; i += step {
					bufferChan <- buffer[i : i+4]
				}
			}
		}
	}
}
