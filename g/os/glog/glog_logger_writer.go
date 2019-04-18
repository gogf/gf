// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package glog

type Writer struct {
	logger *Logger
}

// Write implements the io.Writer interface.
// It just prints the content with header or level.
func (w *Writer) Write(p []byte) (n int, err error) {
	w.logger.Header(false).Print(string(p))
	return len(p), nil
}