// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gerror_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/internal/json"
	"github.com/gogf/gf/v2/test/gtest"
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
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gerror.Wrapf(nil, ""), nil)
	})
}

func Test_WrapSkip(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gerror.WrapSkip(1, nil, "2"), nil)
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
		t.Assert(gerror.WrapSkipf(1, nil, "2"), nil)
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
		t.Assert(gerror.Cause(nil), nil)
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

	gtest.C(t, func(t *gtest.T) {
		t.Assert(gerror.Stack(nil), "")
		err := errors.New("1")
		t.Assert(gerror.Stack(err), err)
	})

	gtest.C(t, func(t *gtest.T) {
		var e *gerror.Error = nil
		t.Assert(e.Cause(), nil)
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
		// fmt.Printf("%+v", err)
	})

	gtest.C(t, func(t *gtest.T) {
		err := gerror.New("1")
		t.AssertNE(fmt.Sprintf("%+v", err), "1")
		// fmt.Printf("%+v", err)
	})

	gtest.C(t, func(t *gtest.T) {
		err := gerror.New("1")
		err = gerror.Wrap(err, "2")
		err = gerror.Wrap(err, "3")
		t.AssertNE(err, nil)
		// fmt.Printf("%+v", err)
	})
}

func Test_Current(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gerror.Current(nil), nil)
		err := errors.New("1")
		err = gerror.Wrap(err, "2")
		err = gerror.Wrap(err, "3")
		t.Assert(err.Error(), "3: 2: 1")
		t.Assert(gerror.Current(err).Error(), "3")
	})
	gtest.C(t, func(t *gtest.T) {
		var e *gerror.Error = nil
		t.Assert(e.Current(), nil)
	})
}

func Test_Unwrap(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gerror.Unwrap(nil), nil)
		err := errors.New("1")
		err = gerror.Wrap(err, "2")
		err = gerror.Wrap(err, "3")
		t.Assert(err.Error(), "3: 2: 1")

		err = gerror.Unwrap(err)
		t.Assert(err.Error(), "2: 1")

		err = gerror.Unwrap(err)
		t.Assert(err.Error(), "1")

		err = gerror.Unwrap(err)
		t.AssertNil(err)
	})
	gtest.C(t, func(t *gtest.T) {
		var e *gerror.Error = nil
		t.Assert(e.Unwrap(), nil)
	})
}

func Test_Code(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		err := errors.New("123")
		t.Assert(gerror.Code(err), -1)
		t.Assert(err.Error(), "123")
	})
	gtest.C(t, func(t *gtest.T) {
		err := gerror.NewCode(gcode.CodeUnknown, "123")
		t.Assert(gerror.Code(err), gcode.CodeUnknown)
		t.Assert(err.Error(), "123")
	})
	gtest.C(t, func(t *gtest.T) {
		err := gerror.NewCodef(gcode.New(1, "", nil), "%s", "123")
		t.Assert(gerror.Code(err).Code(), 1)
		t.Assert(err.Error(), "123")
	})
	gtest.C(t, func(t *gtest.T) {
		err := gerror.NewCodeSkip(gcode.New(1, "", nil), 0, "123")
		t.Assert(gerror.Code(err).Code(), 1)
		t.Assert(err.Error(), "123")
	})
	gtest.C(t, func(t *gtest.T) {
		err := gerror.NewCodeSkipf(gcode.New(1, "", nil), 0, "%s", "123")
		t.Assert(gerror.Code(err).Code(), 1)
		t.Assert(err.Error(), "123")
	})
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gerror.WrapCode(gcode.New(1, "", nil), nil, "3"), nil)
		err := errors.New("1")
		err = gerror.Wrap(err, "2")
		err = gerror.WrapCode(gcode.New(1, "", nil), err, "3")
		t.Assert(gerror.Code(err).Code(), 1)
		t.Assert(err.Error(), "3: 2: 1")
	})
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gerror.WrapCodef(gcode.New(1, "", nil), nil, "%s", "3"), nil)
		err := errors.New("1")
		err = gerror.Wrap(err, "2")
		err = gerror.WrapCodef(gcode.New(1, "", nil), err, "%s", "3")
		t.Assert(gerror.Code(err).Code(), 1)
		t.Assert(err.Error(), "3: 2: 1")
	})
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gerror.WrapCodeSkip(gcode.New(1, "", nil), 100, nil, "3"), nil)
		err := errors.New("1")
		err = gerror.Wrap(err, "2")
		err = gerror.WrapCodeSkip(gcode.New(1, "", nil), 100, err, "3")
		t.Assert(gerror.Code(err).Code(), 1)
		t.Assert(err.Error(), "3: 2: 1")
	})
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gerror.WrapCodeSkipf(gcode.New(1, "", nil), 100, nil, "%s", "3"), nil)
		err := errors.New("1")
		err = gerror.Wrap(err, "2")
		err = gerror.WrapCodeSkipf(gcode.New(1, "", nil), 100, err, "%s", "3")
		t.Assert(gerror.Code(err).Code(), 1)
		t.Assert(err.Error(), "3: 2: 1")
	})
}

func TestError_Error(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var e *gerror.Error = nil
		t.Assert(e.Error(), nil)
	})
}

func TestError_Code(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var e *gerror.Error = nil
		t.Assert(e.Code(), gcode.CodeNil)
	})
}

func Test_SetCode(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		err := gerror.New("123")
		t.Assert(gerror.Code(err), -1)
		t.Assert(err.Error(), "123")

		err.(*gerror.Error).SetCode(gcode.CodeValidationFailed)
		t.Assert(gerror.Code(err), gcode.CodeValidationFailed)
		t.Assert(err.Error(), "123")
	})
	gtest.C(t, func(t *gtest.T) {
		var err *gerror.Error = nil
		err.SetCode(gcode.CodeValidationFailed)
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

func Test_HasStack(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		err1 := errors.New("1")
		err2 := gerror.New("1")
		t.Assert(gerror.HasStack(err1), false)
		t.Assert(gerror.HasStack(err2), true)
	})
}

func Test_Equal(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		err1 := errors.New("1")
		err2 := errors.New("1")
		err3 := gerror.New("1")
		err4 := gerror.New("4")
		t.Assert(gerror.Equal(err1, err2), false)
		t.Assert(gerror.Equal(err1, err3), true)
		t.Assert(gerror.Equal(err2, err3), true)
		t.Assert(gerror.Equal(err3, err4), false)
		t.Assert(gerror.Equal(err1, err4), false)
	})
	gtest.C(t, func(t *gtest.T) {
		var e = new(gerror.Error)
		t.Assert(e.Equal(e), true)
	})
}

func Test_Is(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		err1 := errors.New("1")
		err2 := gerror.Wrap(err1, "2")
		err2 = gerror.Wrap(err2, "3")
		t.Assert(gerror.Is(err2, err1), true)
	})
}

func Test_HashError(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		err1 := errors.New("1")
		err2 := gerror.Wrap(err1, "2")
		err2 = gerror.Wrap(err2, "3")
		t.Assert(gerror.HasError(err2, err1), true)
	})
}

func Test_HashCode(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gerror.HasCode(nil, gcode.CodeNotAuthorized), false)
		err1 := errors.New("1")
		err2 := gerror.WrapCode(gcode.CodeNotAuthorized, err1, "2")
		err3 := gerror.Wrap(err2, "3")
		err4 := gerror.Wrap(err3, "4")
		t.Assert(gerror.HasCode(err1, gcode.CodeNotAuthorized), false)
		t.Assert(gerror.HasCode(err2, gcode.CodeNotAuthorized), true)
		t.Assert(gerror.HasCode(err3, gcode.CodeNotAuthorized), true)
		t.Assert(gerror.HasCode(err4, gcode.CodeNotAuthorized), true)
	})
}

func Test_NewOption(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.AssertNE(gerror.NewOption(gerror.Option{
			Error: errors.New("NewOptionError"),
			Stack: true,
			Text:  "Text",
			Code:  gcode.CodeNotAuthorized,
		}), gerror.New("NewOptionError"))
	})
}
