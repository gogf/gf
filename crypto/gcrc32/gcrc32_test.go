// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go -bench=".*"

package gcrc32_test

import (
	"testing"

	"github.com/gogf/gf/v2/crypto/gcrc32"
	"github.com/gogf/gf/v2/crypto/gmd5"
	"github.com/gogf/gf/v2/test/gtest"
)

func TestEncrypt(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := "pibigstar"
		result := 693191136
		encrypt1 := gcrc32.Encrypt(s)
		encrypt2 := gcrc32.Encrypt([]byte(s))
		t.AssertEQ(int(encrypt1), result)
		t.AssertEQ(int(encrypt2), result)

		strmd5, _ := gmd5.Encrypt(s)
		test1 := gcrc32.Encrypt(strmd5)
		test2 := gcrc32.Encrypt([]byte(strmd5))
		t.AssertEQ(test2, test1)
	})
}
