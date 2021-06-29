package mlog

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gcmd"
	"github.com/gogf/gf/os/genv"
	"github.com/gogf/gf/os/glog"
)

const (
	headerPrintEnvName = "GF_CLI_MLOG_HEADER"
)

var (
	logger = glog.New()
)

func init() {
	logger.SetStack(false)
	logger.SetDebug(false)
	if genv.Get(headerPrintEnvName) == "1" {
		logger.SetHeaderPrint(true)
	} else {
		logger.SetHeaderPrint(false)
	}
	if gcmd.ContainsOpt("debug") {
		logger.SetDebug(true)
	}
}

// SetHeaderPrint enables/disables header printing to stdout.
func SetHeaderPrint(enabled bool) {
	logger.SetHeaderPrint(enabled)
	if enabled {
		genv.Set(headerPrintEnvName, "1")
	} else {
		genv.Set(headerPrintEnvName, "0")
	}
}

func Print(v ...interface{}) {
	logger.Print(v...)
}

func Printf(format string, v ...interface{}) {
	logger.Printf(format, v...)
}

func Fatal(v ...interface{}) {
	logger.Fatal(append(g.Slice{"Error:"}, v...)...)
}

func Fatalf(format string, v ...interface{}) {
	logger.Fatalf("Error: "+format, v...)
}

func Debug(v ...interface{}) {
	logger.Debug(append(g.Slice{"Debug:"}, v...)...)
}

func Debugf(format string, v ...interface{}) {
	logger.Debugf("Debug: "+format, v...)
}
