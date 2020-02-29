// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import (
	"bytes"
	"github.com/gogf/gf/os/glog"
)

// errorLogger is the error logging logger for underlying net/http.Server.
type errorLogger struct {
	logger *glog.Logger
}

// Write implements the io.Writer interface.
func (l *errorLogger) Write(p []byte) (n int, err error) {
	l.logger.Skip(1).Error(string(bytes.TrimRight(p, "\r\n")))
	return len(p), nil
}
