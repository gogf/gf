package utils

import (
	"fmt"

	"github.com/gogf/gf/cmd/gf/v2/internal/utility/mlog"
	"github.com/gogf/gf/v2/os/gproc"
)

var (
	gofmtPath     = gproc.SearchBinaryPath("gofmt")     // gofmtPath is the binary path of command `gofmt`.
	goimportsPath = gproc.SearchBinaryPath("goimports") // gofmtPath is the binary path of command `goimports`.
)

// GoFmt formats the source file using command `gofmt -w -s PATH`.
func GoFmt(path string) {
	if gofmtPath == "" {
		mlog.Fatal(`command "gofmt" not found`)
	}
	var command = fmt.Sprintf(`%s -w %s`, gofmtPath, path)
	result, err := gproc.ShellExec(command)
	if err != nil {
		mlog.Fatalf(`error executing command "%s": %s`, command, result)
	}
}

// GoImports formats the source file using command `goimports -w PATH`.
func GoImports(path string) {
	if goimportsPath == "" {
		mlog.Fatal(`command "goimports" not found`)
	}
	var command = fmt.Sprintf(`%s -w %s`, goimportsPath, path)
	result, err := gproc.ShellExec(command)
	if err != nil {
		mlog.Fatalf(`error executing command "%s": %s`, command, result)
	}
}
