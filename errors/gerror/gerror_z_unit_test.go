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

// customError is used to test As function
type customError struct {
	Message string
}

func (e *customError) Error() string {
	return e.Message
}

// anotherError is used to test As function with different error type
type anotherError struct{}

func (e *anotherError) Error() string {
	return "another error"
}

// customCauseError implements ICause interface
type customCauseError struct {
	msg   string
	cause error
}

func (e *customCauseError) Error() string { return e.msg }
func (e *customCauseError) Cause() error  { return e.cause }

// customStackError implements IStack interface
type customStackError struct {
	msg   string
	stack string
}

func (e *customStackError) Error() string { return e.msg }
func (e *customStackError) Stack() string { return e.stack }

// customCurrentError implements ICurrent interface
type customCurrentError struct {
	msg     string
	current error
}

func (e *customCurrentError) Error() string  { return e.msg }
func (e *customCurrentError) Current() error { return e.current }

// customUnwrapError implements IUnwrap interface
type customUnwrapError struct {
	msg    string
	unwrap error
}

func (e *customUnwrapError) Error() string { return e.msg }
func (e *customUnwrapError) Unwrap() error { return e.unwrap }

// customEqualError implements IEqual interface
type customEqualError struct {
	msg string
}

func (e *customEqualError) Error() string           { return e.msg }
func (e *customEqualError) Equal(target error) bool { return e.msg == target.Error() }

// customCodeError implements ICode interface
type customCodeError struct {
	msg  string
	code gcode.Code
}

func (e *customCodeError) Error() string    { return e.msg }
func (e *customCodeError) Code() gcode.Code { return e.code }

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
	gtest.C(t, func(t *gtest.T) {
		errNormal := gerror.New("test")
		b, e := json.Marshal(errNormal)
		t.Assert(e, nil)
		t.Assert(string(b), `"test"`)
	})
	gtest.C(t, func(t *gtest.T) {
		// The string contains special characters.
		errWithSign := gerror.New(`test ""`)
		b, e := json.Marshal(errWithSign)
		t.Assert(e, nil)
		t.Assert(string(b), `"test \"\""`)
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

		var (
			errNotFound = errors.New("not found")
			gerror1     = gerror.Wrap(errNotFound, "wrapped")
			gerror2     = gerror.New("not found")
		)
		t.Assert(errors.Is(errNotFound, errNotFound), true)
		t.Assert(errors.Is(nil, errNotFound), false)
		t.Assert(errors.Is(nil, nil), true)

		t.Assert(gerror.Is(errNotFound, errNotFound), true)
		t.Assert(gerror.Is(nil, errNotFound), false)
		t.Assert(gerror.Is(nil, nil), true)

		t.Assert(errors.Is(gerror1, errNotFound), true)
		t.Assert(errors.Is(gerror2, errNotFound), false)
		t.Assert(gerror.Is(gerror1, errNotFound), true)
		t.Assert(gerror.Is(gerror2, errNotFound), false)
	})
}

func Test_HasError(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		err1 := errors.New("1")
		err2 := gerror.Wrap(err1, "2")
		err2 = gerror.Wrap(err2, "3")
		t.Assert(gerror.HasError(err2, err1), true)
	})
}

func Test_HasCode(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gerror.HasCode(nil, gcode.CodeNotAuthorized), false)
		err1 := errors.New("1")
		err2 := gerror.WrapCode(gcode.CodeNotAuthorized, err1, "2")
		err3 := gerror.Wrap(err2, "3")
		err4 := gerror.Wrap(err3, "4")
		err5 := gerror.WrapCode(gcode.CodeInvalidParameter, err4, "5")
		t.Assert(gerror.HasCode(err1, gcode.CodeNotAuthorized), false)
		t.Assert(gerror.HasCode(err2, gcode.CodeNotAuthorized), true)
		t.Assert(gerror.HasCode(err3, gcode.CodeNotAuthorized), true)
		t.Assert(gerror.HasCode(err4, gcode.CodeNotAuthorized), true)
		t.Assert(gerror.HasCode(err5, gcode.CodeNotAuthorized), true)
		t.Assert(gerror.HasCode(err5, gcode.CodeInvalidParameter), true)
		t.Assert(gerror.HasCode(err5, gcode.CodeInternalError), false)
	})
}

func Test_NewOption(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.AssertNE(gerror.NewWithOption(gerror.Option{
			Error: errors.New("NewOptionError"),
			Stack: true,
			Text:  "Text",
			Code:  gcode.CodeNotAuthorized,
		}), gerror.New("NewOptionError"))
	})
}

func Test_As(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var myerr = &customError{Message: "custom error"}

		// Test with nil error
		var targetErr *customError
		t.Assert(gerror.As(nil, &targetErr), false)
		t.Assert(targetErr, nil)

		// Test with standard error
		err1 := errors.New("standard error")
		t.Assert(gerror.As(err1, &targetErr), false)
		t.Assert(targetErr, nil)

		// Test with custom error type
		err2 := myerr
		t.Assert(gerror.As(err2, &targetErr), true)
		t.Assert(targetErr.Message, "custom error")

		// Test with wrapped error
		err3 := gerror.Wrap(myerr, "wrapped")
		targetErr = nil
		t.Assert(gerror.As(err3, &targetErr), true)
		t.Assert(targetErr.Message, "custom error")

		// Test with deeply wrapped error
		err4 := gerror.Wrap(gerror.Wrap(gerror.Wrap(myerr, "wrap3"), "wrap2"), "wrap1")
		targetErr = nil
		t.Assert(gerror.As(err4, &targetErr), true)
		t.Assert(targetErr.Message, "custom error")

		// Test with different error type
		var otherErr *anotherError
		t.Assert(gerror.As(err4, &otherErr), false)
		t.Assert(otherErr, nil)

		// Test with non-pointer target
		defer func() {
			t.Assert(recover() != nil, true)
		}()
		var nonPtr customError
		gerror.As(err4, nonPtr)
	})

	gtest.C(t, func(t *gtest.T) {
		// Test with nil target
		defer func() {
			t.Assert(recover() != nil, true)
		}()
		gerror.As(errors.New("error"), nil)
	})
}

func Test_NewOption_Deprecated(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Test deprecated NewOption function
		err := gerror.NewOption(gerror.Option{
			Error: errors.New("base error"),
			Stack: true,
			Text:  "option text",
			Code:  gcode.CodeInternalError,
		})
		t.AssertNE(err, nil)
		t.Assert(gerror.Code(err), gcode.CodeInternalError)
	})
}

func Test_Code_WithIUnwrap(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Test Code() with custom error that implements IUnwrap but not ICode
		innerErr := gerror.NewCode(gcode.CodeInternalError, "inner error")
		unwrapErr := &customUnwrapError{msg: "unwrap error", unwrap: innerErr}
		t.Assert(gerror.Code(unwrapErr), gcode.CodeInternalError)
	})
	gtest.C(t, func(t *gtest.T) {
		// Test Code() with nil
		t.Assert(gerror.Code(nil), gcode.CodeNil)
	})
	gtest.C(t, func(t *gtest.T) {
		// Test Code() with custom error that implements ICode
		codeErr := &customCodeError{msg: "code error", code: gcode.CodeNotFound}
		t.Assert(gerror.Code(codeErr), gcode.CodeNotFound)
	})
}

func Test_Cause_WithIUnwrap(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Test Cause() with custom error that implements IUnwrap but not ICause
		rootErr := errors.New("root error")
		unwrapErr := &customUnwrapError{msg: "unwrap error", unwrap: rootErr}
		t.Assert(gerror.Cause(unwrapErr), rootErr)
	})
}

func Test_Cause_WithICause(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Test Cause() with custom error that implements ICause
		rootErr := errors.New("root error")
		causeErr := &customCauseError{msg: "cause error", cause: rootErr}
		t.Assert(gerror.Cause(causeErr), rootErr)
	})
}

func Test_Stack_WithIStack(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Test Stack() with custom error that implements IStack
		stackErr := &customStackError{msg: "stack error", stack: "custom stack trace"}
		t.Assert(gerror.Stack(stackErr), "custom stack trace")
	})
}

func Test_Current_WithICurrent(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Test Current() with custom error that implements ICurrent
		currentErr := errors.New("current error")
		customErr := &customCurrentError{msg: "custom error", current: currentErr}
		t.Assert(gerror.Current(customErr), currentErr)
	})
	gtest.C(t, func(t *gtest.T) {
		// Test Current() with standard error (does not implement ICurrent)
		stdErr := errors.New("standard error")
		t.Assert(gerror.Current(stdErr), stdErr)
	})
}

func Test_Equal_WithIEqual(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Test Equal() when target implements IEqual
		err1 := errors.New("test error")
		err2 := &customEqualError{msg: "test error"}
		t.Assert(gerror.Equal(err1, err2), true)
	})
	gtest.C(t, func(t *gtest.T) {
		// Test Equal() when both are the same
		err := errors.New("test error")
		t.Assert(gerror.Equal(err, err), true)
	})
}

func Test_Error_Cause_WithICause(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Test Error.Cause() when inner error implements ICause
		rootErr := errors.New("root")
		causeErr := &customCauseError{msg: "cause", cause: rootErr}
		wrappedErr := gerror.Wrap(causeErr, "wrapped")
		t.Assert(gerror.Cause(wrappedErr), rootErr)
	})
}

func Test_Error_WithCodeMessage(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Test Error.Error() when text is empty but code has message
		err := gerror.NewCode(gcode.CodeInternalError)
		t.Assert(err.Error(), "Internal Error")
	})
	gtest.C(t, func(t *gtest.T) {
		// Test Error.Error() when text is empty and code has message, with wrapped error
		innerErr := errors.New("inner")
		err := gerror.WrapCode(gcode.CodeInternalError, innerErr)
		t.Assert(err.Error(), "Internal Error: inner")
	})
}

func Test_Format_PlusS(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Test %+s format (stack only)
		err := gerror.New("test error")
		stackStr := fmt.Sprintf("%+s", err)
		t.Assert(len(stackStr) > 0, true)
		t.AssertNE(stackStr, "test error")
	})
}

func Test_Format_MinusS_EmptyText(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Test %-s format when text is empty but code has message
		err := gerror.NewCode(gcode.CodeInternalError)
		result := fmt.Sprintf("%-s", err)
		t.Assert(result, "Internal Error")
	})
}

func Test_Stack_DeepNested(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Test deeply nested errors stack
		err := gerror.New("level1")
		for i := 2; i <= 5; i++ {
			err = gerror.Wrap(err, fmt.Sprintf("level%d", i))
		}
		stack := gerror.Stack(err)
		t.Assert(len(stack) > 0, true)
	})
}

func Test_Stack_NilError(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var err *gerror.Error = nil
		t.Assert(err.Stack(), "")
	})
}

func Test_Stack_WithStandardError(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Test stack with wrapped standard error
		stdErr := errors.New("standard error")
		err := gerror.Wrap(stdErr, "wrapped")
		stack := gerror.Stack(err)
		t.Assert(len(stack) > 0, true)
	})
}

func Test_NewCode_MultipleTexts(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Test NewCode with multiple text arguments
		err := gerror.NewCode(gcode.CodeInternalError, "text1", "text2", "text3")
		t.Assert(err.Error(), "text1, text2, text3")
	})
}

func Test_NewCodeSkip_MultipleTexts(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Test NewCodeSkip with multiple text arguments
		err := gerror.NewCodeSkip(gcode.CodeInternalError, 0, "text1", "text2")
		t.Assert(err.Error(), "text1, text2")
	})
}

func Test_WrapCode_MultipleTexts(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Test WrapCode with multiple text arguments
		innerErr := errors.New("inner")
		err := gerror.WrapCode(gcode.CodeInternalError, innerErr, "text1", "text2")
		t.Assert(err.Error(), "text1, text2: inner")
	})
}

func Test_WrapCodeSkip_MultipleTexts(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Test WrapCodeSkip with multiple text arguments
		innerErr := errors.New("inner")
		err := gerror.WrapCodeSkip(gcode.CodeInternalError, 0, innerErr, "text1", "text2")
		t.Assert(err.Error(), "text1, text2: inner")
	})
}
