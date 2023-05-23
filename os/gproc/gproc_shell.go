// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gproc

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"runtime"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"

	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/text/gstr"
)

// Shell executes command `cmd` synchronously with given input pipe `in` and output pipe `out`.
// The command `cmd` reads the input parameters from input pipe `in`, and writes its output automatically
// to output pipe `out`.
func Shell(ctx context.Context, cmd string, out io.Writer, in io.Reader) error {
	p := NewProcess(
		getShell(),
		append([]string{getShellOption()}, parseCommand(cmd)...),
	)
	p.Stdin = in
	p.Stdout = out
	return p.Run(ctx)
}

// ShellRun executes given command `cmd` synchronously and outputs the command result to the stdout.
func ShellRun(ctx context.Context, cmd string) error {
	p := NewProcess(
		getShell(),
		append([]string{getShellOption()}, parseCommand(cmd)...),
	)
	return p.Run(ctx)
}

// ShellExec executes given command `cmd` synchronously and returns the command result.
func ShellExec(ctx context.Context, cmd string, environment ...[]string) (result string, err error) {
	var (
		buf = bytes.NewBuffer(nil)
		p   = NewProcess(
			getShell(),
			append([]string{getShellOption()}, parseCommand(cmd)...),
			environment...,
		)
	)
	p.Stdout = buf
	p.Stderr = buf
	err = p.Run(ctx)
	result = buf.String()
	return
}

// parseCommand parses command `cmd` into slice arguments.
//
// Note that it just parses the `cmd` for "cmd.exe" binary in windows, but it is not necessary
// parsing the `cmd` for other systems using "bash"/"sh" binary.
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

// getShell returns the shell command depending on current working operating system.
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

// getShellOption returns the shell option depending on current working operating system.
// It returns "/c" for windows, and "-c" for others.
func getShellOption() string {
	switch runtime.GOOS {
	case "windows":
		return "/c"

	default:
		return "-c"
	}
}

// tracingEnvFromCtx converts OpenTelemetry propagation data as environment variables.
func tracingEnvFromCtx(ctx context.Context) []string {
	var (
		a = make([]string, 0)
		m = make(map[string]string)
	)
	otel.GetTextMapPropagator().Inject(ctx, propagation.MapCarrier(m))
	for k, v := range m {
		a = append(a, fmt.Sprintf(`%s=%s`, k, v))
	}
	return a
}
