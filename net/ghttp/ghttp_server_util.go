// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
)

// WrapF is a helper function for wrapping http.HandlerFunc and returns a ghttp.HandlerFunc.
func WrapF(f http.HandlerFunc) HandlerFunc {
	return func(r *Request) {
		f(r.Response.Writer, r.Request)
	}
}

// WrapH is a helper function for wrapping http.Handler and returns a ghttp.HandlerFunc.
func WrapH(h http.Handler) HandlerFunc {
	return func(r *Request) {
		h.ServeHTTP(r.Response.Writer, r.Request)
	}
}

// *struct/ **struct
// Check whether ptr is *struct or **struct []*struct *[]*struct
// If the type is satisfied, directly dereference and return
// Subsequent operations such as setting fields can be performed by returning values.
// Whether it is first level or second level
// Using the following method calls will result in a null pointer
// var t *T
// 1.checkValidRequestParams(t) ===> panic
// var t **T
// 2.checkValidRequestParams(t) ===> panic
// Only the following situations can call successfully
// var t T ==>checkValidRequestParams(&t) OK
// var t *T ==>checkValidRequestParams(&t) OK
func checkValidRequestParams(ptr any) (reflect.Value, reflect.Kind, error) {

	srcVal, ok := ptr.(reflect.Value)
	if ok {
		// used for standard routing
		srcTyp := srcVal.Type()
		// XXXReq(ctx,&LoginReq{})
		elem, err := checkValidStruct(srcTyp, srcVal, 1)
		return elem, reflect.Struct, err
	} else {
		//Routes registered by other rules
		srcTyp := reflect.TypeOf(ptr)
		if srcTyp.Kind() != reflect.Ptr {
			return srcVal, 0, gerror.NewCodef(
				gcode.CodeInvalidParameter,
				`invalid parameter type "%v", of which kind should be of *struct/**struct/*[]struct/*[]*struct, but got: "%v"`,
				srcTyp,
				srcTyp.Kind(),
			)

		}
		srcValof := reflect.ValueOf(ptr)
		// Determine whether it is a nil pointer
		if srcValof.IsNil() {
			return srcVal, 0, gerror.NewCodef(
				gcode.CodeInvalidParameter, `Cannot pass in a null pointer`)

		}

		elemTyp := srcTyp.Elem()
		//Number of dereferences
		derefCount := 1
		// Determine whether it is a secondary pointer
		if elemTyp.Kind() == reflect.Ptr {
			elemTyp = elemTyp.Elem()
			derefCount++
		}
		// The type is correct when we get here, and it is not a null pointer.
		switch elemTyp.Kind() {
		case reflect.Slice, reflect.Array:
			elem, err := checkValidStructSlice(elemTyp, srcValof, derefCount)
			fmt.Println("elemTyp==", elemTyp)
			return elem, elemTyp.Kind(), err

		case reflect.Struct:
			elem, err := checkValidStruct(elemTyp, srcValof, derefCount)
			return elem, reflect.Struct, err

		default:

			return reflect.Value{}, 0, gerror.NewCodef(
				gcode.CodeInvalidParameter,
				`invalid parameter type "%v", of which kind should be of *struct/**struct/*[]struct/*[]*struct, but got: "%v"`,
				srcTyp,
				srcTyp.Kind(),
			)
		}
	}

}

func checkValidStructSlice(elemTyp reflect.Type, srcValof reflect.Value, derefCount int) (reflect.Value, error) {
	// []slice =>  *struct / struct
	// After slice dereference it may be *struct/struct
	elemTyp = elemTyp.Elem()
	if elemTyp.Kind() == reflect.Ptr {
		elemTyp = elemTyp.Elem()
	}

	if elemTyp.Kind() != reflect.Struct {
		return srcValof, gerror.NewCodef(
			gcode.CodeInvalidParameter,
			`invalid parameter type "%v", of which kind should be of *struct/**struct/*[]struct/*[]*struct, but got: "%v"`,
			elemTyp,
			elemTyp.Kind(),
		)
	}
	// type ok
	return srcValof.Elem(), nil
}

// derefCount is equal to several levels of pointers.
// The first-level pointer is 1 and the second-level pointer is 2. Dereference is required.
func checkValidStruct(elemTyp reflect.Type, srcValof reflect.Value, derefCount int) (reflect.Value, error) {
	// The type is correct when we get here
	var elemPtr reflect.Value
	if derefCount == 1 {
		elemPtr = srcValof.Elem()
	} else if derefCount == 2 {
		// Check whether the secondary pointer is nil, if so, use new assignment
		srcValof = srcValof.Elem()
		if srcValof.IsNil() {
			newTyp := reflect.New(elemTyp)
			srcValof.Set(newTyp)
		}
		elemPtr = srcValof.Elem()

	} else {
		return srcValof, gerror.NewCodef(
			gcode.CodeInvalidParameter,
			`invalid parameter type "%v", of which kind should be of *struct/**struct/*[]struct/*[]*struct, but got: "%v"`,
			elemTyp,
			elemTyp.Kind(),
		)
	}
	// Come here, the pointer has been dereferenced and the value can be set.
	return elemPtr, nil
}
