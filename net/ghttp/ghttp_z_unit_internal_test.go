// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import (
	"context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gsession"
	"github.com/gogf/gf/v2/os/gstructs"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gtag"
	"github.com/gogf/gf/v2/util/guid"
)

type testRequestTagsRes struct {
	Name string
}

type testRequestTagsNoTag struct {
	Name string
	Age  int
}

type testRequestTagsSub struct {
	X string `v:"required"`
}

type testRequestTagsNested struct {
	Sub testRequestTagsSub
}

type testRequestTagsNestedPointer struct {
	Sub *testRequestTagsSub
}

type testRequestTagsSlice struct {
	Items []testRequestTagsSub
}

// TestRequestTagsEmbeddedSub is an exported embedded struct type for tests.
type TestRequestTagsEmbeddedSub struct {
	X string `v:"required"`
}

type testRequestTagsEmbeddedExported struct {
	TestRequestTagsEmbeddedSub
}

// The attributes of an unexported embedded struct are not validated,
// please see Test_RequestTags_EmbeddedUnexported.
type testRequestTagsEmbeddedUnexported struct {
	testRequestTagsSub
}

type testRequestTagsIn struct {
	Name string `in:"header" p:"name"`
}

type testRequestTagsDefault struct {
	Id int `d:"100"`
}

type testRequestTagsSelfRef struct {
	Name string
	Next *testRequestTagsSelfRef
}

// testRequestTagsCyclicPointer is a cyclic pointer type whose element type is itself.
type testRequestTagsCyclicPointer *testRequestTagsCyclicPointer

type testRequestTagsCyclicPointerReq struct {
	Name string
	Data testRequestTagsCyclicPointer
	Next *testRequestTagsCyclicPointerReq
}

// The default value comes from a tag variable, please see Test_RequestTags_TagVariable.
type testRequestTagsTagVariable struct {
	Page int `d:"{testRequestTagsPageSize}"`
}

// The source of the `in` tag comes from a tag variable, please see Test_RequestTags_TagVariableIn.
type testRequestTagsTagVariableIn struct {
	Name string `in:"{testRequestTagsInSource}" p:"name"`
}

type testRequestTagsMap struct {
	Items map[string]testRequestTagsSub
}

// The nested named struct is not a request struct field that the parameter
// merging handles, so its `in` tag does not disable the optimization.
type testRequestTagsNestedIn struct {
	Sub testRequestTagsNestedInSub
}

type testRequestTagsNestedInSub struct {
	Name string `in:"header" p:"name"`
}

// The empty tag values take no effect, please see Test_RequestTags_EmptyValues.
type testRequestTagsEmptyValues struct {
	EmptyDefault string `d:""`
	EmptyIn      string `in:"" p:"empty_in"`
	EmptyValid   string `v:""`
}

// The space tag value takes effect as it is not empty.
type testRequestTagsSpaceValue struct {
	Name string `d:" "`
}

// The unexported attributes share their names with the exported ones,
// for which the parameter merging of the `in` and `default` tags still
// takes effect as the converted values are matched by the exported ones.
type testRequestTagsUnexportedIn struct {
	Name string
	name string `in:"header"`
}

type testRequestTagsUnexportedDefault struct {
	Id int
	id int `d:"100"`
}

type testRequestTagsSelfRefTagged struct {
	Sub *testRequestTagsSelfRefTagged `v:"required"`
}

// newTestHandler registers given standard handler on `/test` and returns the
// server and the handler information, as the registration procedure does.
func newTestHandler(t testing.TB, handler any) (*Server, handlerFuncInfo) {
	var s = GetServer("test-request-tags-" + guid.S())
	s.sessionManager = gsession.New(time.Hour, gsession.NewStorageMemory())
	funcInfo, err := s.checkAndCreateFuncInfo(handler, "", "", "")
	if err != nil {
		t.Fatal(err)
	}
	s.doBindHandler(context.Background(), doBindHandlerInput{
		Pattern:  "ALL:/test",
		FuncInfo: funcInfo,
	})
	return s, funcInfo
}

// serveTestRequest serves one JSON request for the handler of given server.
func serveTestRequest(t testing.TB, s *Server, funcInfo handlerFuncInfo, body string, headers ...map[string]string) *Request {
	var req = httptest.NewRequest(http.MethodPost, "/test", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	if len(headers) > 0 {
		for k, v := range headers[0] {
			req.Header.Set(k, v)
		}
	}
	var r = newRequest(s, req, httptest.NewRecorder())
	r.handlers, r.serveHandler, r.hasHookHandler, r.hasServeHandler = s.getHandlersWithCache(r)
	funcInfo.Func(r)
	return r
}

// newTestServingRequest registers given standard handler on `/test` and serves
// one JSON request for it, as the real serving procedure does.
func newTestServingRequest(t testing.TB, handler any, body string, headers ...map[string]string) *Request {
	var s, funcInfo = newTestHandler(t, handler)
	return serveTestRequest(t, s, funcInfo, body, headers...)
}

// Test_RequestTags_Scan tests the prechecking of the struct tag usages at
// handler registration.
func Test_RequestTags_Scan(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var cases = []struct {
			pointer    any
			hasValid   bool
			hasIn      bool
			hasDefault bool
		}{
			{&testRequestTagsNoTag{}, false, false, false},
			{&testRequestTagsSub{}, true, false, false},
			{&testRequestTagsNested{}, true, false, false},
			{&testRequestTagsNestedPointer{}, true, false, false},
			{&testRequestTagsSlice{}, true, false, false},
			{&testRequestTagsMap{}, true, false, false},
			{&testRequestTagsEmbeddedExported{}, true, false, false},
			// The attributes of the unexported embedded struct are included in the
			// request struct fields which the parameter merging handles.
			{&testRequestTagsEmbeddedUnexported{}, true, false, false},
			{&testRequestTagsIn{}, false, true, false},
			{&testRequestTagsNestedIn{}, false, false, false},
			{&testRequestTagsDefault{}, false, false, true},
			{&testRequestTagsUnexportedIn{}, false, true, false},
			{&testRequestTagsUnexportedDefault{}, false, false, true},
			// The self reference types should not cause infinite recursion.
			{&testRequestTagsSelfRef{}, false, false, false},
			{&testRequestTagsSelfRefTagged{}, true, false, false},
			// The cyclic pointer type should not cause infinite dereference.
			{&testRequestTagsCyclicPointerReq{}, false, false, false},
		}
		for _, c := range cases {
			fields, err := gstructs.Fields(gstructs.FieldsInput{
				Pointer:         c.pointer,
				RecursiveOption: gstructs.RecursiveOptionEmbedded,
			})
			t.AssertNil(err)
			var tags = scanRequestStructTags(reflect.TypeOf(c.pointer), fields)
			t.Assert(tags.HasValid, c.hasValid)
			t.Assert(tags.HasIn, c.hasIn)
			t.Assert(tags.HasDefault, c.hasDefault)
		}
	})
}

// Test_RequestTags_NoValidTag tests the serving procedure for a request struct
// that has no validation tag, for which the validation is skipped while the
// parameter converting still works.
func Test_RequestTags_NoValidTag(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var captured *testRequestTagsNoTag
		var r = newTestServingRequest(t, func(ctx context.Context, req *testRequestTagsNoTag) (*testRequestTagsRes, error) {
			captured = req
			return nil, nil
		}, `{"name":"john","age":18}`)
		t.AssertNil(r.error)
		t.Assert(captured.Name, "john")
		t.Assert(captured.Age, 18)
		// The validation is prechecked unnecessary at registration.
		t.Assert(r.serveHandler.Handler.Info.ReqStructTags.HasValid, false)
	})
}

// Test_RequestTags_NestedValidTag tests that the validation tags inside nested
// attribute types are detected at registration, so the validation is still
// performed for them.
func Test_RequestTags_NestedValidTag(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var r = newTestServingRequest(t, func(ctx context.Context, req *testRequestTagsNested) (*testRequestTagsRes, error) {
			return nil, nil
		}, `{}`)
		t.AssertNE(r.error, nil)
		t.Assert(gerror.Code(r.error), gcode.CodeValidationFailed)
		// The validation is prechecked necessary at registration.
		t.Assert(r.serveHandler.Handler.Info.ReqStructTags.HasValid, true)
	})
}

// Test_RequestTags_InTag tests that the `in` tag still works for the serving procedure.
func Test_RequestTags_InTag(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var captured *testRequestTagsIn
		var r = newTestServingRequest(t, func(ctx context.Context, req *testRequestTagsIn) (*testRequestTagsRes, error) {
			captured = req
			return nil, nil
		}, `{}`, map[string]string{"name": "john"})
		t.AssertNil(r.error)
		t.Assert(captured.Name, "john")
	})
}

// Test_RequestTags_DefaultTag tests that the `default` tag still works for the serving procedure.
func Test_RequestTags_DefaultTag(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var captured *testRequestTagsDefault
		var r = newTestServingRequest(t, func(ctx context.Context, req *testRequestTagsDefault) (*testRequestTagsRes, error) {
			captured = req
			return nil, nil
		}, `{}`)
		t.AssertNil(r.error)
		t.Assert(captured.Id, 100)
	})
}

// Test_RequestTags_UnexportedTag tests that the tags of unexported attributes
// are also prechecked, as the parameter merging handles all the request struct
// fields, not only the exported ones.
func Test_RequestTags_UnexportedTag(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var capturedIn *testRequestTagsUnexportedIn
		var r = newTestServingRequest(t, func(ctx context.Context, req *testRequestTagsUnexportedIn) (*testRequestTagsRes, error) {
			capturedIn = req
			return nil, nil
		}, `{}`, map[string]string{"name": "john"})
		t.AssertNil(r.error)
		t.Assert(capturedIn.Name, "john")
		t.Assert(r.serveHandler.Handler.Info.ReqStructTags.HasIn, true)

		var capturedDefault *testRequestTagsUnexportedDefault
		r = newTestServingRequest(t, func(ctx context.Context, req *testRequestTagsUnexportedDefault) (*testRequestTagsRes, error) {
			capturedDefault = req
			return nil, nil
		}, `{}`)
		t.AssertNil(r.error)
		t.Assert(capturedDefault.Id, 100)
		t.Assert(r.serveHandler.Handler.Info.ReqStructTags.HasDefault, true)
	})
}

// Test_RequestTags_CyclicPointer tests the handler registration for a request
// struct whose attribute is a cyclic pointer type, whose element type is itself.
//
// Note that the registration should not hang for such types. Note that the
// serving of such a request struct hangs in package gconv, which dereferences
// the attribute type in loop as well (see genFieldConvertFunc of package
// structcache), and it is not related to this prechecking.
func Test_RequestTags_CyclicPointer(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		_, funcInfo := newTestHandler(t, func(ctx context.Context, req *testRequestTagsCyclicPointerReq) (*testRequestTagsRes, error) {
			return nil, nil
		})
		var tags = funcInfo.ReqStructTags
		t.Assert(tags.HasValid, false)
		t.Assert(tags.HasIn, false)
		t.Assert(tags.HasDefault, false)
	})
}

// Test_RequestTags_TagVariable tests the tag usage prechecking for a tag value
// that comes from a tag variable, which might be registered or updated after
// the handler registration, see package gtag.
func Test_RequestTags_TagVariable(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// The tag variable is registered as empty, and is filled with the real
		// value after the handler registration but before the serving.
		gtag.SetOver("testRequestTagsPageSize", "")
		defer gtag.SetOver("testRequestTagsPageSize", "")

		var captured *testRequestTagsTagVariable
		var s, funcInfo = newTestHandler(t, func(ctx context.Context, req *testRequestTagsTagVariable) (*testRequestTagsRes, error) {
			captured = req
			return nil, nil
		})
		gtag.SetOver("testRequestTagsPageSize", "100")

		var r = serveTestRequest(t, s, funcInfo, `{}`)
		t.AssertNil(r.error)
		t.Assert(captured.Page, 100)
		t.Assert(r.serveHandler.Handler.Info.ReqStructTags.HasDefault, true)
	})
}

// Test_RequestTags_TagVariableIn tests the tag usage prechecking for an `in`
// tag whose value comes from a tag variable, which might be registered or
// updated after the handler registration, see package gtag.
func Test_RequestTags_TagVariableIn(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// The tag variable is registered as empty, and is filled with the real
		// value after the handler registration but before the serving.
		gtag.SetOver("testRequestTagsInSource", "")
		defer gtag.SetOver("testRequestTagsInSource", "")

		var captured *testRequestTagsTagVariableIn
		var s, funcInfo = newTestHandler(t, func(ctx context.Context, req *testRequestTagsTagVariableIn) (*testRequestTagsRes, error) {
			captured = req
			return nil, nil
		})
		gtag.SetOver("testRequestTagsInSource", "header")

		var r = serveTestRequest(t, s, funcInfo, `{}`, map[string]string{"name": "john"})
		t.AssertNil(r.error)
		t.Assert(captured.Name, "john")
		t.Assert(r.serveHandler.Handler.Info.ReqStructTags.HasIn, true)
	})
}

// Test_RequestTags_EmptyValues tests the tag prechecking for the empty tag
// values, which take no effect for the parameter merging and the validation,
// so they are prechecked as unused.
func Test_RequestTags_EmptyValues(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var captured *testRequestTagsEmptyValues
		var r = newTestServingRequest(t, func(ctx context.Context, req *testRequestTagsEmptyValues) (*testRequestTagsRes, error) {
			captured = req
			return nil, nil
		}, `{}`, map[string]string{"empty_in": "from-header"})
		t.AssertNil(r.error)
		t.Assert(captured.EmptyDefault, "")
		t.Assert(captured.EmptyIn, "")
		t.Assert(captured.EmptyValid, "")
		var tags = r.serveHandler.Handler.Info.ReqStructTags
		t.Assert(tags.HasValid, false)
		t.Assert(tags.HasIn, false)
		t.Assert(tags.HasDefault, false)

		// The space tag value is not empty, so it takes effect.
		var capturedSpace *testRequestTagsSpaceValue
		r = newTestServingRequest(t, func(ctx context.Context, req *testRequestTagsSpaceValue) (*testRequestTagsRes, error) {
			capturedSpace = req
			return nil, nil
		}, `{}`)
		t.AssertNil(r.error)
		t.Assert(capturedSpace.Name, " ")
		t.Assert(r.serveHandler.Handler.Info.ReqStructTags.HasDefault, true)
	})
}

// Test_RequestTags_SameMap tests the same map object checking, which is used by
// the request parameter merging to avoid the duplicated merging of the form map
// and the body map, as they are the same map for JSON content requests.
func Test_RequestTags_SameMap(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var m1 = map[string]any{"name": "john"}
		var m2 = map[string]any{"name": "john"}
		t.Assert(sameMap(m1, m1), true)
		t.Assert(sameMap(m1, m2), false)
		t.Assert(sameMap(m1, nil), false)
		t.Assert(sameMap(nil, nil), false)
	})
}

// Test_RequestTags_SelfRef tests the serving procedure for a self reference
// request struct that has no validation tag.
func Test_RequestTags_SelfRef(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var r = newTestServingRequest(t, func(ctx context.Context, req *testRequestTagsSelfRef) (*testRequestTagsRes, error) {
			return nil, nil
		}, `{"name":"john"}`)
		t.AssertNil(r.error)
	})
}

// Test_RequestTags_EmbeddedUnexported tests the serving procedure for a request
// struct that embeds an unexported struct which has validation tags.
//
// Note that the attributes of an unexported embedded struct are not validated
// by package gvalid, but they are included in the request struct fields which
// the parameter merging handles, so the tag is prechecked as used.
func Test_RequestTags_EmbeddedUnexported(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var r = newTestServingRequest(t, func(ctx context.Context, req *testRequestTagsEmbeddedUnexported) (*testRequestTagsRes, error) {
			return nil, nil
		}, `{}`)
		t.AssertNil(r.error)
		t.Assert(r.serveHandler.Handler.Info.ReqStructTags.HasValid, true)
	})
}
