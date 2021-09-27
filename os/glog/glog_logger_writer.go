// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package glog

import (
	"bytes"
	"context"
)

// Write implements the io.Writer interface.
// It just prints the content using Print.
func (l *Logger) Write(p []byte) (n int, err error) {
	l.Header(false).Print(context.TODO(), string(bytes.TrimRight(p, "\r\n")))
	return len(p), nil
}
