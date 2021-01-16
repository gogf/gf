// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gerror_test

import (
	"errors"
	"fmt"
	"github.com/gogf/gf/internal/json"
	"testing"

	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/test/gtest"
)

func nilError() error {
	return nil
}

func Test_Nil(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gerror.New(""), nil)
		t.Assert(gerror.Wrap(nilError(), "test"), nil)
	})
}

func Test_New(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		err := gerror.New("1")
		t.AssertNE(err, nil)
		t.Assert(err.Error(), "1")
	})
	gtest.C(t, func(t *gtest.T) {
		err := gerror.Newf("%d", 1)
		t.AssertNE(err, nil)
		t.Assert(err.Error(), "1")
	})
	gtest.C(t, func(t *gtest.T) {
		err := gerror.NewSkip(1, "1")
		t.AssertNE(err, nil)
		t.Assert(err.Error(), "1")
	})
	gtest.C(t, func(t *gtest.T) {
		err := gerror.NewSkipf(1, "%d", 1)
		t.AssertNE(err, nil)
		t.Assert(err.Error(), "1")
	})
}

func Test_Wrap(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		err := errors.New("1")
		err = gerror.Wrap(err, "2")
		err = gerror.Wrap(err, "3")
		t.AssertNE(err, nil)
		t.Assert(err.Error(), "3: 2: 1")
	})
	gtest.C(t, func(t *gtest.T) {
		err := gerror.New("1")
		err = gerror.Wrap(err, "2")
		err = gerror.Wrap(err, "3")
		t.AssertNE(err, nil)
		t.Assert(err.Error(), "3: 2: 1")
	})
	gtest.C(t, func(t *gtest.T) {
		err := gerror.New("1")
		err = gerror.Wrap(err, "")
		t.AssertNE(err, nil)
		t.Assert(err.Error(), "1")
	})
}

func Test_Wrapf(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		err := errors.New("1")
		err = gerror.Wrapf(err, "%d", 2)
		err = gerror.Wrapf(err, "%d", 3)
		t.AssertNE(err, nil)
		t.Assert(err.Error(), "3: 2: 1")
	})
	gtest.C(t, func(t *gtest.T) {
		err := gerror.New("1")
		err = gerror.Wrapf(err, "%d", 2)
		err = gerror.Wrapf(err, "%d", 3)
		t.AssertNE(err, nil)
		t.Assert(err.Error(), "3: 2: 1")
	})
	gtest.C(t, func(t *gtest.T) {
		err := gerror.New("1")
		err = gerror.Wrapf(err, "")
		t.AssertNE(err, nil)
		t.Assert(err.Error(), "1")
	})
}

func Test_WrapSkip(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		err := errors.New("1")
		err = gerror.WrapSkip(1, err, "2")
		err = gerror.WrapSkip(1, err, "3")
		t.AssertNE(err, nil)
		t.Assert(err.Error(), "3: 2: 1")
	})
	gtest.C(t, func(t *gtest.T) {
		err := gerror.New("1")
		err = gerror.WrapSkip(1, err, "2")
		err = gerror.WrapSkip(1, err, "3")
		t.AssertNE(err, nil)
		t.Assert(err.Error(), "3: 2: 1")
	})
	gtest.C(t, func(t *gtest.T) {
		err := gerror.New("1")
		err = gerror.WrapSkip(1, err, "")
		t.AssertNE(err, nil)
		t.Assert(err.Error(), "1")
	})
}

func Test_WrapSkipf(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		err := errors.New("1")
		err = gerror.WrapSkipf(1, err, "2")
		err = gerror.WrapSkipf(1, err, "3")
		t.AssertNE(err, nil)
		t.Assert(err.Error(), "3: 2: 1")
	})
	gtest.C(t, func(t *gtest.T) {
		err := gerror.New("1")
		err = gerror.WrapSkipf(1, err, "2")
		err = gerror.WrapSkipf(1, err, "3")
		t.AssertNE(err, nil)
		t.Assert(err.Error(), "3: 2: 1")
	})
	gtest.C(t, func(t *gtest.T) {
		err := gerror.New("1")
		err = gerror.WrapSkipf(1, err, "")
		t.AssertNE(err, nil)
		t.Assert(err.Error(), "1")
	})
}

func Test_Cause(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		err := errors.New("1")
		t.Assert(gerror.Cause(err), err)
	})

	gtest.C(t, func(t *gtest.T) {
		err := errors.New("1")
		err = gerror.Wrap(err, "2")
		err = gerror.Wrap(err, "3")
		t.Assert(gerror.Cause(err), "1")
	})

	gtest.C(t, func(t *gtest.T) {
		err := gerror.New("1")
		t.Assert(gerror.Cause(err), "1")
	})

	gtest.C(t, func(t *gtest.T) {
		err := gerror.New("1")
		err = gerror.Wrap(err, "2")
		err = gerror.Wrap(err, "3")
		t.Assert(gerror.Cause(err), "1")
	})
}

func Test_Format(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		err := errors.New("1")
		err = gerror.Wrap(err, "2")
		err = gerror.Wrap(err, "3")
		t.AssertNE(err, nil)
		t.Assert(fmt.Sprintf("%s", err), "3: 2: 1")
		t.Assert(fmt.Sprintf("%v", err), "3: 2: 1")
	})

	gtest.C(t, func(t *gtest.T) {
		err := gerror.New("1")
		err = gerror.Wrap(err, "2")
		err = gerror.Wrap(err, "3")
		t.AssertNE(err, nil)
		t.Assert(fmt.Sprintf("%s", err), "3: 2: 1")
		t.Assert(fmt.Sprintf("%v", err), "3: 2: 1")
	})

	gtest.C(t, func(t *gtest.T) {
		err := gerror.New("1")
		err = gerror.Wrap(err, "2")
		err = gerror.Wrap(err, "3")
		t.AssertNE(err, nil)
		t.Assert(fmt.Sprintf("%-s", err), "3")
		t.Assert(fmt.Sprintf("%-v", err), "3")
	})
}

func Test_Stack(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		err := errors.New("1")
		t.Assert(fmt.Sprintf("%+v", err), "1")
	})

	gtest.C(t, func(t *gtest.T) {
		err := errors.New("1")
		err = gerror.Wrap(err, "2")
		err = gerror.Wrap(err, "3")
		t.AssertNE(err, nil)
		//fmt.Printf("%+v", err)
	})

	gtest.C(t, func(t *gtest.T) {
		err := gerror.New("1")
		t.AssertNE(fmt.Sprintf("%+v", err), "1")
		//fmt.Printf("%+v", err)
	})

	gtest.C(t, func(t *gtest.T) {
		err := gerror.New("1")
		err = gerror.Wrap(err, "2")
		err = gerror.Wrap(err, "3")
		t.AssertNE(err, nil)
		//fmt.Printf("%+v", err)
	})
}

func Test_Current(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		err := errors.New("1")
		err = gerror.Wrap(err, "2")
		err = gerror.Wrap(err, "3")
		t.Assert(err.Error(), "3: 2: 1")
		t.Assert(gerror.Current(err).Error(), "3")
	})
}

func Test_Next(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		err := errors.New("1")
		err = gerror.Wrap(err, "2")
		err = gerror.Wrap(err, "3")
		t.Assert(err.Error(), "3: 2: 1")

		err = gerror.Next(err)
		t.Assert(err.Error(), "2: 1")

		err = gerror.Next(err)
		t.Assert(err.Error(), "1")

		err = gerror.Next(err)
		t.Assert(err, nil)
	})
}

func Test_Code(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		err := errors.New("123")
		t.Assert(gerror.Code(err), -1)
		t.Assert(err.Error(), "123")
	})
	gtest.C(t, func(t *gtest.T) {
		err := gerror.NewCode(1, "123")
		t.Assert(gerror.Code(err), 1)
		t.Assert(err.Error(), "123")
	})
	gtest.C(t, func(t *gtest.T) {
		err := gerror.NewCodef(1, "%s", "123")
		t.Assert(gerror.Code(err), 1)
		t.Assert(err.Error(), "123")
	})
	gtest.C(t, func(t *gtest.T) {
		err := gerror.NewCodeSkip(1, 0, "123")
		t.Assert(gerror.Code(err), 1)
		t.Assert(err.Error(), "123")
	})
	gtest.C(t, func(t *gtest.T) {
		err := gerror.NewCodeSkipf(1, 0, "%s", "123")
		t.Assert(gerror.Code(err), 1)
		t.Assert(err.Error(), "123")
	})
	gtest.C(t, func(t *gtest.T) {
		err := errors.New("1")
		err = gerror.Wrap(err, "2")
		err = gerror.WrapCode(1, err, "3")
		t.Assert(gerror.Code(err), 1)
		t.Assert(err.Error(), "3: 2: 1")
	})
	gtest.C(t, func(t *gtest.T) {
		err := errors.New("1")
		err = gerror.Wrap(err, "2")
		err = gerror.WrapCodef(1, err, "%s", "3")
		t.Assert(gerror.Code(err), 1)
		t.Assert(err.Error(), "3: 2: 1")
	})
	gtest.C(t, func(t *gtest.T) {
		err := errors.New("1")
		err = gerror.Wrap(err, "2")
		err = gerror.WrapCodeSkip(1, 100, err, "3")
		t.Assert(gerror.Code(err), 1)
		t.Assert(err.Error(), "3: 2: 1")
	})
	gtest.C(t, func(t *gtest.T) {
		err := errors.New("1")
		err = gerror.Wrap(err, "2")
		err = gerror.WrapCodeSkipf(1, 100, err, "%s", "3")
		t.Assert(gerror.Code(err), 1)
		t.Assert(err.Error(), "3: 2: 1")
	})
}

func Test_Json(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		err := gerror.Wrap(gerror.New("1"), "2")
		b, e := json.Marshal(err)
		t.Assert(e, nil)
		t.Assert(string(b), `"2: 1"`)
	})
}
