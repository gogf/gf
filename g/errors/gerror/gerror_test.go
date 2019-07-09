// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gerror_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/gogf/gf/g/errors/gerror"
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
		gtest.Assert(gerror.New(interfaceNil()), nil)
		gtest.Assert(gerror.Wrap(nilError(), "test"), nil)
	})
}

func Test_Wrap(t *testing.T) {
	gtest.Case(t, func() {
		err := errors.New("1")
		err = gerror.Wrap(err, "2")
		err = gerror.Wrap(err, "3")
		gtest.AssertNE(err, nil)
		gtest.Assert(err.Error(), "3: 2: 1")
	})

	gtest.Case(t, func() {
		err := gerror.New("1")
		err = gerror.Wrap(err, "2")
		err = gerror.Wrap(err, "3")
		gtest.AssertNE(err, nil)
		gtest.Assert(err.Error(), "3: 2: 1")
	})
}

func Test_Cause(t *testing.T) {
	gtest.Case(t, func() {
		err := errors.New("1")
		gtest.Assert(gerror.Cause(err), err)
	})

	gtest.Case(t, func() {
		err := errors.New("1")
		err = gerror.Wrap(err, "2")
		err = gerror.Wrap(err, "3")
		gtest.Assert(gerror.Cause(err), "1")
	})

	gtest.Case(t, func() {
		err := gerror.New("1")
		gtest.Assert(gerror.Cause(err), "1")
	})

	gtest.Case(t, func() {
		err := gerror.New("1")
		err = gerror.Wrap(err, "2")
		err = gerror.Wrap(err, "3")
		gtest.Assert(gerror.Cause(err), "1")
	})
}

func Test_Format(t *testing.T) {
	gtest.Case(t, func() {
		err := errors.New("1")
		err = gerror.Wrap(err, "2")
		err = gerror.Wrap(err, "3")
		gtest.AssertNE(err, nil)
		gtest.Assert(fmt.Sprintf("%s", err), "3: 2: 1")
		gtest.Assert(fmt.Sprintf("%v", err), "3: 2: 1")
	})

	gtest.Case(t, func() {
		err := gerror.New("1")
		err = gerror.Wrap(err, "2")
		err = gerror.Wrap(err, "3")
		gtest.AssertNE(err, nil)
		gtest.Assert(fmt.Sprintf("%s", err), "3: 2: 1")
		gtest.Assert(fmt.Sprintf("%v", err), "3: 2: 1")
	})

	gtest.Case(t, func() {
		err := gerror.New("1")
		err = gerror.Wrap(err, "2")
		err = gerror.Wrap(err, "3")
		gtest.AssertNE(err, nil)
		gtest.Assert(fmt.Sprintf("%-s", err), "3")
		gtest.Assert(fmt.Sprintf("%-v", err), "3")
	})
}

func Test_Stack(t *testing.T) {
	gtest.Case(t, func() {
		err := errors.New("1")
		gtest.Assert(fmt.Sprintf("%+v", err), "1")
	})

	gtest.Case(t, func() {
		err := errors.New("1")
		err = gerror.Wrap(err, "2")
		err = gerror.Wrap(err, "3")
		gtest.AssertNE(err, nil)
		//fmt.Printf("%+v", err)
	})

	gtest.Case(t, func() {
		err := gerror.New("1")
		gtest.AssertNE(fmt.Sprintf("%+v", err), "1")
		//fmt.Printf("%+v", err)
	})

	gtest.Case(t, func() {
		err := gerror.New("1")
		err = gerror.Wrap(err, "2")
		err = gerror.Wrap(err, "3")
		gtest.AssertNE(err, nil)
		//fmt.Printf("%+v", err)
	})
}
