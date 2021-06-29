package utils

import (
	"fmt"
	"github.com/gogf/gf/os/gproc"
)

var (
	// gofmtPath is the binary path of command `gofmt`.
	gofmtPath = gproc.SearchBinaryPath("gofmt")
)

// GoFmt formats the source file using command `gofmt -w -s PATH`.
func GoFmt(path string) {
	if gofmtPath != "" {
		gproc.ShellExec(fmt.Sprintf(`%s -w -s %s`, gofmtPath, path))
	}
}
