// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gproc implements management and communication for processes.
package gproc

import (
	"os"
	"runtime"
	"time"

	"github.com/gogf/gf/v2/os/genv"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
)

const (
	envKeyPPid            = "GPROC_PPID"
	tracingInstrumentName = "github.com/gogf/gf/v2/os/gproc.Process"
)

var (
	processPid       = os.Getpid() // processPid is the pid of current process.
	processStartTime = time.Now()  // processStartTime is the start time of current process.
)

// Pid returns the pid of current process.
func Pid() int {
	return processPid
}

// PPid returns the custom parent pid if exists, or else it returns the system parent pid.
func PPid() int {
	if !IsChild() {
		return Pid()
	}
	ppidValue := os.Getenv(envKeyPPid)
	if ppidValue != "" && ppidValue != "0" {
		return gconv.Int(ppidValue)
	}
	return PPidOS()
}

// PPidOS returns the system parent pid of current process.
// Note that the difference between PPidOS and PPid function is that the PPidOS returns
// the system ppid, but the PPid functions may return the custom pid by gproc if the custom
// ppid exists.
func PPidOS() int {
	return os.Getppid()
}

// IsChild checks and returns whether current process is a child process.
// A child process is forked by another gproc process.
func IsChild() bool {
	ppidValue := os.Getenv(envKeyPPid)
	return ppidValue != "" && ppidValue != "0"
}

// SetPPid sets custom parent pid for current process.
func SetPPid(ppid int) error {
	if ppid > 0 {
		return os.Setenv(envKeyPPid, gconv.String(ppid))
	} else {
		return os.Unsetenv(envKeyPPid)
	}
}

// StartTime returns the start time of current process.
func StartTime() time.Time {
	return processStartTime
}

// Uptime returns the duration which current process has been running
func Uptime() time.Duration {
	return time.Now().Sub(processStartTime)
}

// SearchBinary searches the binary `file` in current working folder and PATH environment.
func SearchBinary(file string) string {
	// Check if it is absolute path of exists at current working directory.
	if gfile.Exists(file) {
		return file
	}
	return SearchBinaryPath(file)
}

// SearchBinaryPath searches the binary `file` in PATH environment.
func SearchBinaryPath(file string) string {
	array := ([]string)(nil)
	switch runtime.GOOS {
	case "windows":
		envPath := genv.Get("PATH", genv.Get("Path")).String()
		if gstr.Contains(envPath, ";") {
			array = gstr.SplitAndTrim(envPath, ";")
		} else if gstr.Contains(envPath, ":") {
			array = gstr.SplitAndTrim(envPath, ":")
		}
		if gfile.Ext(file) != ".exe" {
			file += ".exe"
		}

	default:
		array = gstr.SplitAndTrim(genv.Get("PATH").String(), ":")
	}
	if len(array) > 0 {
		path := ""
		for _, v := range array {
			path = v + gfile.Separator + file
			if gfile.Exists(path) && gfile.IsFile(path) {
				return path
			}
		}
	}
	return ""
}
