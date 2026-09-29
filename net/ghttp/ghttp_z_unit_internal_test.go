// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import (
	"testing"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gstructs"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/text/gstr"
)

// reqBodyValidReq declares the JSON array request body by one slice field.
type reqBodyValidReq struct {
	Items []int `json:"items" in:"body"`
}

// reqBodyPointerReq declares the JSON array request body by one pointer to slice field.
type reqBodyPointerReq struct {
	Items *[]int `json:"items" in:"body"`
}

// reqBodyNoTagNameReq declares the `in:"body"` tag without any priority tag name, of which the
// field name itself is used as the tag name.
type reqBodyNoTagNameReq struct {
	Items []int `in:"body"`
}

// reqBodyNotSliceReq declares the `in:"body"` tag on a field of non slice type.
type reqBodyNotSliceReq struct {
	Item int `json:"item" in:"body"`
}

// reqBodyFixedArrayReq declares the `in:"body"` tag on a fixed-size array field, which should
// use a slice with a length validation rule instead.
type reqBodyFixedArrayReq struct {
	Items [2]int `json:"items" in:"body"`
}

// reqBodyStructData is the attribute type of reqBodyStructReq.
type reqBodyStructData struct {
	Name string `json:"name"`
}

// reqBodyStructReq declares the `in:"body"` tag on a struct field, which does not need the tag
// as an object request body is received by the ordinary request struct fields.
type reqBodyStructReq struct {
	Data reqBodyStructData `json:"data" in:"body"`
}

// reqBodyUploadFileReq declares the `in:"body"` tag on a single uploading file field, which is
// not able to receive the request body.
type reqBodyUploadFileReq struct {
	File *UploadFile `json:"file" in:"body"`
}

// reqBodyUploadFilesReq declares the `in:"body"` tag on an uploading file slice field, of which
// the type is a slice as well but not a bodyable one.
type reqBodyUploadFilesReq struct {
	Files []*UploadFile `json:"files" in:"body"`
}

// reqBodyUploadFilesAliasReq declares the `in:"body"` tag on the UploadFiles type, which is
// another form of the uploading file slice field.
type reqBodyUploadFilesAliasReq struct {
	Files UploadFiles `json:"files" in:"body"`
}

// reqBodyDuplicatedInTagReq declares the `in` tag more than once with different values on one
// field, of which the parameter location is ambiguous.
type reqBodyDuplicatedInTagReq struct {
	Item string `json:"item" in:"query" in:"header"`
}

// reqBodySameInTagReq declares the `in` tag more than once with the same value, which is
// unambiguous for the parameter location.
type reqBodySameInTagReq struct {
	Item string `json:"item" in:"query" in:"query"`
}

// reqBodyDuplicatedReq declares the `in:"body"` tag on multiple fields.
type reqBodyDuplicatedReq struct {
	Items []int `json:"items" in:"body"`
	More  []int `json:"more" in:"body"`
}

// reqBodyAbsentReq declares no `in:"body"` tag at all.
type reqBodyAbsentReq struct {
	Items []int `json:"items"`
}

func newReqBodyFuncInfo(t *gtest.T, pointer any) handlerFuncInfo {
	fields, err := gstructs.Fields(gstructs.FieldsInput{
		Pointer:         pointer,
		RecursiveOption: gstructs.RecursiveOptionEmbedded,
	})
	t.AssertNil(err)
	return handlerFuncInfo{ReqStructFields: fields}
}

func Test_HandlerFuncInfo_CheckAndCreateReqBodyField(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// No field tagged with `in:"body"`.
		var info = newReqBodyFuncInfo(t, &reqBodyAbsentReq{})
		t.AssertNil(info.checkAndCreateReqBodyField())
		t.Assert(info.ReqBodyFieldName, "")

		// One slice field tagged with `in:"body"`.
		info = newReqBodyFuncInfo(t, &reqBodyValidReq{})
		t.AssertNil(info.checkAndCreateReqBodyField())
		t.Assert(info.ReqBodyFieldName, "Items")
		t.Assert(info.ReqBodyFieldTagName, "items")

		// One pointer to slice field tagged with `in:"body"`.
		info = newReqBodyFuncInfo(t, &reqBodyPointerReq{})
		t.AssertNil(info.checkAndCreateReqBodyField())
		t.Assert(info.ReqBodyFieldName, "Items")
		t.Assert(info.ReqBodyFieldTagName, "items")

		// The `in:"body"` field has no priority tag name, which then uses its field name.
		info = newReqBodyFuncInfo(t, &reqBodyNoTagNameReq{})
		t.AssertNil(info.checkAndCreateReqBodyField())
		t.Assert(info.ReqBodyFieldName, "Items")
		t.Assert(info.ReqBodyFieldTagName, "Items")

		// The field tagged with `in:"body"` is not of slice type.
		info = newReqBodyFuncInfo(t, &reqBodyNotSliceReq{})
		err := info.checkAndCreateReqBodyField()
		t.AssertNE(err, nil)
		t.Assert(gerror.Code(err), gcode.CodeInvalidParameter)

		// The field tagged with `in:"body"` is a fixed-size array.
		info = newReqBodyFuncInfo(t, &reqBodyFixedArrayReq{})
		err = info.checkAndCreateReqBodyField()
		t.AssertNE(err, nil)
		t.Assert(gerror.Code(err), gcode.CodeInvalidParameter)
		t.Assert(gstr.Contains(err.Error(), "fixed-size array"), true)

		// The field tagged with `in:"body"` is a struct.
		info = newReqBodyFuncInfo(t, &reqBodyStructReq{})
		err = info.checkAndCreateReqBodyField()
		t.AssertNE(err, nil)
		t.Assert(gerror.Code(err), gcode.CodeInvalidParameter)
		t.Assert(gstr.Contains(err.Error(), "ordinary fields of the request struct"), true)

		// The field tagged with `in:"body"` is of a single uploading file type.
		info = newReqBodyFuncInfo(t, &reqBodyUploadFileReq{})
		err = info.checkAndCreateReqBodyField()
		t.AssertNE(err, nil)
		t.Assert(gerror.Code(err), gcode.CodeInvalidParameter)
		t.Assert(gstr.Contains(err.Error(), "file type"), true)

		// The fields tagged with `in:"body"` are of uploading file slice types, which are also
		// slices but are not able to receive the JSON array request body.
		for _, req := range []any{&reqBodyUploadFilesReq{}, &reqBodyUploadFilesAliasReq{}} {
			info = newReqBodyFuncInfo(t, req)
			err = info.checkAndCreateReqBodyField()
			t.AssertNE(err, nil)
			t.Assert(gerror.Code(err), gcode.CodeInvalidParameter)
			t.Assert(gstr.Contains(err.Error(), "file type"), true)
			t.Assert(gstr.Contains(err.Error(), "multipart/form-data"), true)
		}

		// Multiple fields tagged with `in:"body"`.
		info = newReqBodyFuncInfo(t, &reqBodyDuplicatedReq{})
		err = info.checkAndCreateReqBodyField()
		t.AssertNE(err, nil)
		t.Assert(gerror.Code(err), gcode.CodeInvalidParameter)
	})
}

func Test_HandlerFuncInfo_CheckDuplicatedInTag(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// No `in` tag at all.
		var info = newReqBodyFuncInfo(t, &reqBodyAbsentReq{})
		t.AssertNil(info.checkDuplicatedInTag())

		// One `in` tag on the body field.
		info = newReqBodyFuncInfo(t, &reqBodyValidReq{})
		t.AssertNil(info.checkDuplicatedInTag())

		// The same `in` value declared more than once is not ambiguous.
		info = newReqBodyFuncInfo(t, &reqBodySameInTagReq{})
		t.AssertNil(info.checkDuplicatedInTag())

		// More than one `in` tag with different values on one field.
		info = newReqBodyFuncInfo(t, &reqBodyDuplicatedInTagReq{})
		err := info.checkDuplicatedInTag()
		t.AssertNE(err, nil)
		t.Assert(gerror.Code(err), gcode.CodeInvalidParameter)
		t.Assert(gstr.Contains(err.Error(), `conflicting values "query" and "header"`), true)
	})
}
