// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package utils

import (
	"io"
)

// ReadCloser implements the io.ReadCloser interface
// which is used for reading request body content multiple times.
//
// Note that it cannot be closed.
type ReadCloser struct {
	index      int    // Current read position.
	content    []byte // Content.
	repeatable bool   // Mark the content can be repeatable read.
}

// NewReadCloser creates and returns a RepeatReadCloser object.
func NewReadCloser(content []byte, repeatable bool) io.ReadCloser {
	return &ReadCloser{
		content:    content,
		repeatable: repeatable,
	}
}

// Read implements the io.ReadCloser interface.
func (b *ReadCloser) Read(p []byte) (n int, err error) {
	// Make it repeatable reading.
	if b.index >= len(b.content) && b.repeatable {
		b.index = 0
	}
	n = copy(p, b.content[b.index:])
	b.index += n
	if b.index >= len(b.content) {
		return n, io.EOF
	}
	return n, nil
}

// Close implements the io.ReadCloser interface.
func (b *ReadCloser) Close() error {
	return nil
}
