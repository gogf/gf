// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go -bench=".*"

package gcrc32_test

import (
	"testing"

	"github.com/gogf/gf/g/crypto/gcrc32"
	"github.com/gogf/gf/g/test/gtest"
)

func TestEncrypt(t *testing.T) {
	gtest.Case(t, func() {
		s := "pibigstar"
		result := 693191136
		encrypt1 := gcrc32.EncryptString(s)
		encrypt2 := gcrc32.EncryptBytes([]byte(s))
		gtest.AssertEQ(int(encrypt1), result)
		gtest.AssertEQ(int(encrypt2), result)
	})
}
