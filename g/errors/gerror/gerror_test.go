// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gerror_test

import (
	"testing"

	"github.com/gogf/gf/g/internal/errors"
	"github.com/gogf/gf/g/test/gtest"
)

func interfaceNil() interface{} {
	return nil
}

func nilError() error {
	return nil
}

func Test_Nil(t *testing.T) {
	gtest.Case(t, func() {
		gtest.Assert(errors.New(interfaceNil()), nil)
		gtest.Assert(errors.Wrap(nilError(), "test"), nil)
	})
}

func Test_Wrap(t *testing.T) {
	gtest.Case(t, func() {
		err := errors.New("1")
		err = errors.Wrap(err, "func2 error")
		err = errors.Wrap(err, "func3 error")
		gtest.AssertNE(err, nil)
		gtest.Assert(err.Error(), "func3 error: func2 error: 1")
	})
}
