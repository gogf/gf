// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package utils_test

import (
	"regexp"
	"testing"

	"github.com/gogf/gf/v2/internal/utils"
)

var (
	replaceCharReg, _ = regexp.Compile(`[\-\.\_\s]+`)
)

func Benchmark_RemoveSymbols(b *testing.B) {
	for i := 0; i < b.N; i++ {
		utils.RemoveSymbols(`-a-b._a c1!@#$%^&*()_+:";'.,'01`)
	}
}

func Benchmark_RegularReplaceChars(b *testing.B) {
	for i := 0; i < b.N; i++ {
		replaceCharReg.ReplaceAllString(`-a-b._a c1!@#$%^&*()_+:";'.,'01`, "")
	}
}
