// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package mlog

import (
	"context"

	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/os/genv"
	"github.com/gogf/gf/v2/os/glog"
)

const (
	headerPrintEnvName = "GF_CLI_MLOG_HEADER"
)

var (
	ctx    = context.TODO()
	logger = glog.New()
)

func init() {
	if genv.Get(headerPrintEnvName).String() == "1" {
		logger.SetHeaderPrint(true)
	} else {
		logger.SetHeaderPrint(false)
	}

	if gcmd.GetOpt("debug") != nil || gcmd.GetOpt("gf.debug") != nil {
		logger.SetHeaderPrint(true)
		logger.SetStackSkip(4)
		logger.SetFlags(logger.GetFlags() | glog.F_FILE_LONG)
		logger.SetDebug(true)
	} else {
		logger.SetStack(false)
		logger.SetDebug(false)
	}
}

// SetHeaderPrint enables/disables header printing to stdout.
func SetHeaderPrint(enabled bool) {
	logger.SetHeaderPrint(enabled)
	if enabled {
		_ = genv.Set(headerPrintEnvName, "1")
	} else {
		_ = genv.Set(headerPrintEnvName, "0")
	}
}

func Print(v ...interface{}) {
	logger.Print(ctx, v...)
}

func Printf(format string, v ...interface{}) {
	logger.Printf(ctx, format, v...)
}

func Fatal(v ...interface{}) {
	logger.Fatal(ctx, v...)
}

func Fatalf(format string, v ...interface{}) {
	logger.Fatalf(ctx, format, v...)
}

func Debug(v ...interface{}) {
	logger.Debug(ctx, v...)
}

func Debugf(format string, v ...interface{}) {
	logger.Debugf(ctx, format, v...)
}
