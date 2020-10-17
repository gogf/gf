// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gproc implements management and communication for processes.
package gproc

import (
	"bytes"
	"github.com/gogf/gf/os/genv"
	"github.com/gogf/gf/text/gstr"
	"io"
	"os"
	"runtime"
	"time"

	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/util/gconv"
)

const (
	gPROC_ENV_KEY_PPID_KEY = "GPROC_PPID"
)

var (
	// processPid is the pid of current process.
	processPid = os.Getpid()
	// processStartTime is the start time of current process.
	processStartTime = time.Now()
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
	ppidValue := os.Getenv(gPROC_ENV_KEY_PPID_KEY)
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
	ppidValue := os.Getenv(gPROC_ENV_KEY_PPID_KEY)
	return ppidValue != "" && ppidValue != "0"
}

// SetPPid sets custom parent pid for current process.
func SetPPid(ppid int) error {
	if ppid > 0 {
		return os.Setenv(gPROC_ENV_KEY_PPID_KEY, gconv.String(ppid))
	} else {
		return os.Unsetenv(gPROC_ENV_KEY_PPID_KEY)
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

// Shell executes command <cmd> synchronizingly with given input pipe <in> and output pipe <out>.
// The command <cmd> reads the input parameters from input pipe <in>, and writes its output automatically
// to output pipe <out>.
func Shell(cmd string, out io.Writer, in io.Reader) error {
	p := NewProcess(getShell(), append([]string{getShellOption()}, parseCommand(cmd)...))
	p.Stdin = in
	p.Stdout = out
	return p.Run()
}

// ShellRun executes given command <cmd> synchronizingly and outputs the command result to the stdout.
func ShellRun(cmd string) error {
	p := NewProcess(getShell(), append([]string{getShellOption()}, parseCommand(cmd)...))
	return p.Run()
}

// ShellExec executes given command <cmd> synchronizingly and returns the command result.
func ShellExec(cmd string, environment ...[]string) (string, error) {
	buf := bytes.NewBuffer(nil)
	p := NewProcess(getShell(), append([]string{getShellOption()}, parseCommand(cmd)...), environment...)
	p.Stdout = buf
	p.Stderr = buf
	err := p.Run()
	return buf.String(), err
}

// parseCommand parses command <cmd> into slice arguments.
//
// Note that it just parses the <cmd> for "cmd.exe" binary in windows, but it is not necessary
// parsing the <cmd> for other systems using "bash"/"sh" binary.
func parseCommand(cmd string) (args []string) {
	if runtime.GOOS != "windows" {
		return []string{cmd}
	}
	// Just for "cmd.exe" in windows.
	var argStr string
	var firstChar, prevChar, lastChar1, lastChar2 byte
	array := gstr.SplitAndTrim(cmd, " ")
	for _, v := range array {
		if len(argStr) > 0 {
			argStr += " "
		}
		firstChar = v[0]
		lastChar1 = v[len(v)-1]
		lastChar2 = 0
		if len(v) > 1 {
			lastChar2 = v[len(v)-2]
		}
		if prevChar == 0 && (firstChar == '"' || firstChar == '\'') {
			// It should remove the first quote char.
			argStr += v[1:]
			prevChar = firstChar
		} else if prevChar != 0 && lastChar2 != '\\' && lastChar1 == prevChar {
			// It should remove the last quote char.
			argStr += v[:len(v)-1]
			args = append(args, argStr)
			argStr = ""
			prevChar = 0
		} else if len(argStr) > 0 {
			argStr += v
		} else {
			args = append(args, v)
		}
	}
	return
}

// getShell returns the shell command depending on current working operation system.
// It returns "cmd.exe" for windows, and "bash" or "sh" for others.
func getShell() string {
	switch runtime.GOOS {
	case "windows":
		return SearchBinary("cmd.exe")
	default:
		// Check the default binary storage path.
		if gfile.Exists("/bin/bash") {
			return "/bin/bash"
		}
		if gfile.Exists("/bin/sh") {
			return "/bin/sh"
		}
		// Else search the env PATH.
		path := SearchBinary("bash")
		if path == "" {
			path = SearchBinary("sh")
		}
		return path
	}
}

// getShellOption returns the shell option depending on current working operation system.
// It returns "/c" for windows, and "-c" for others.
func getShellOption() string {
	switch runtime.GOOS {
	case "windows":
		return "/c"
	default:
		return "-c"
	}
}

// SearchBinary searches the binary <file> in current working folder and PATH environment.
func SearchBinary(file string) string {
	// Check if it's absolute path of exists at current working directory.
	if gfile.Exists(file) {
		return file
	}
	return SearchBinaryPath(file)
}

// SearchBinaryPath searches the binary <file> in PATH environment.
func SearchBinaryPath(file string) string {
	array := ([]string)(nil)
	switch runtime.GOOS {
	case "windows":
		envPath := genv.Get("PATH", genv.Get("Path"))
		if gstr.Contains(envPath, ";") {
			array = gstr.SplitAndTrim(envPath, ";")
		} else if gstr.Contains(envPath, ":") {
			array = gstr.SplitAndTrim(envPath, ":")
		}
		if gfile.Ext(file) != ".exe" {
			file += ".exe"
		}
	default:
		array = gstr.SplitAndTrim(genv.Get("PATH"), ":")
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
