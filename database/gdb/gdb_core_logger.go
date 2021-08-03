// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"context"
	"github.com/gogf/gf/os/glog"
)

// LoggerImp is the default implementation of interface Logger for DB.
type LoggerImp struct {
	*glog.Logger
}

// Error implements function Error for interface Logger.
func (l LoggerImp) Error(ctx context.Context, s string) {
	l.Ctx(ctx).Error(s)
}

// Debug implements function Debug for interface Logger.
func (l LoggerImp) Debug(ctx context.Context, s string) {
	l.Ctx(ctx).Debug(s)
}
