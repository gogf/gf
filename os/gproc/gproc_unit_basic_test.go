// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gproc

import (
	"github.com/gogf/gf/test/gtest"
	"testing"
)

func Test_parseCommand(t *testing.T) {
	gtest.Case(t, func() {
		commandMap := map[string]interface{}{
			`cmd`:             []string{`cmd`},
			`cmd /c`:          []string{`cmd`, `/c`},
			`cmd /c go build`: []string{`cmd`, `/c`, `go`, `build`},
			`cmd /c go build -ldflags "-X 'a=123' -X 'b=456'" test.go`: []string{`cmd`, `/c`, `go`, `build`, `-ldflags`, `"-X 'a=123' -X 'b=456'"`, `test.go`},
		}
		for k, v := range commandMap {
			gtest.Assert(parseCommand(k), v)
		}
	})
}
