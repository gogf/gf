// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gerror_test

import (
	"fmt"
	"testing"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/test/gtest"
)

func Test_ErrorStack(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		err := gfErrorsTest1()
		errStackStr := fmt.Sprintf("%+v", err)

		t.Assert(errStackStr, `
`)
	})
}

func gfErrorsTest1() error {
	return gfErrorsTest2()
}

func gfErrorsTest2() error {
	return gfErrorsTest3()
}

func gfErrorsTest3() error {
	return gfErrorsTest4()
}

func gfErrorsTest4() error {
	return gerror.Wrap(gfErrorsTest5(), "gerror/Wrap:test wrap")
}

func gfErrorsTest5() error {
	return gerror.New("gerror/New:test new")
}
