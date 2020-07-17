// Copyright 2017 gf Author(https://github.com/jin502437344/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/jin502437344/gf.

// go test *.go -bench=".*" -benchmem

package cmdenv

import (
	"os"
	"testing"

	"github.com/jin502437344/gf/test/gtest"
)

func Test_Get(t *testing.T) {
	os.Args = []string{"--gf.test.value1=111"}
	os.Setenv("GF_TEST_VALUE1", "222")
	os.Setenv("GF_TEST_VALUE2", "333")
	doInit()
	gtest.C(t, func(t *gtest.T) {
		t.Assert(Get("gf.test.value1").String(), "111")
		t.Assert(Get("gf.test.value2").String(), "333")
		t.Assert(Get("gf.test.value3").String(), "")
		t.Assert(Get("gf.test.value3", 1).String(), "1")
	})
}
