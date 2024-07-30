// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import (
	"context"
	"reflect"
	"unsafe"
)

type iface struct {
	typ  unsafe.Pointer
	data unsafe.Pointer
}

func doAsmClosureCallStrictRoute(funcInfo *handlerFuncInfo, r *Request, req unsafe.Pointer) {
	switch funcInfo.Type.NumOut() {
	case 2:
		if funcInfo.Type.Out(0).Kind() == reflect.Slice {
			res, err := doAnyCallRequest_with_sliceRes_err(funcInfo.handlerFuncClosure, r.Context(), req)
			sliceValue := reflect.New(funcInfo.Type.Out(0)).Elem().Interface()
			out := (*iface)(unsafe.Pointer(&sliceValue))
			(*out).data = unsafe.Pointer(&res)
			r.handlerResponse = sliceValue
			r.error = err
		} else {
			res, err := doAnyCallRequest_with_res_err(funcInfo.handlerFuncClosure, r.Context(), req)
			outType := funcInfo.Type.Out(0).Elem()
			outValue := reflect.New(outType).Interface()
			out := (*iface)(unsafe.Pointer(&outValue))
			(*out).data = res
			r.handlerResponse = outValue
			r.error = err
		}
	case 1:
		err := doAnyCallRequest_with_err(funcInfo.handlerFuncClosure, r.Context(), req)
		r.error = err
	}
}

func doAnyCallRequest_with_err(fn any, ctx context.Context, req unsafe.Pointer) error

func doAnyCallRequest_with_res_err(fn any, ctx context.Context, req unsafe.Pointer) (unsafe.Pointer, error)

type _slice struct {
	ptr unsafe.Pointer
	len int
	cap int
}

func doAnyCallRequest_with_sliceRes_err(fn any, ctx context.Context, req unsafe.Pointer) (_slice, error)

func doAsmMethodCallStrictRoute(funcInfo *handlerFuncInfo, r *Request, req unsafe.Pointer) {
	switch funcInfo.Type.NumOut() {
	case 2:
		if funcInfo.Type.Out(0).Kind() == reflect.Slice {
			res, err := doMethodCallRequest_with_sliceRes_err(funcInfo.rawHandlerFuncCodePtr, funcInfo.objPointer, r.Context(), req)
			sliceValue := reflect.New(funcInfo.Type.Out(0)).Elem().Interface()
			out := (*iface)(unsafe.Pointer(&sliceValue))
			(*out).data = unsafe.Pointer(&res)
			r.handlerResponse = sliceValue
			r.error = err
		} else {
			res, err := doMethodCallRequest_with_res_err(funcInfo.rawHandlerFuncCodePtr, funcInfo.objPointer, r.Context(), req)
			outType := funcInfo.Type.Out(0).Elem()
			outValue := reflect.New(outType).Interface()
			out := (*iface)(unsafe.Pointer(&outValue))
			(*out).data = res
			r.handlerResponse = outValue
			r.error = err
		}
	case 1:
		err := doMethodCallRequest_with_err(funcInfo.rawHandlerFuncCodePtr, funcInfo.objPointer, r.Context(), req)
		r.error = err
	}
}

func doMethodCallRequest_with_err(fn, obj unsafe.Pointer, ctx context.Context, req unsafe.Pointer) error

func doMethodCallRequest_with_res_err(fn, obj unsafe.Pointer, ctx context.Context, req unsafe.Pointer) (unsafe.Pointer, error)

func doMethodCallRequest_with_sliceRes_err(fn, obj unsafe.Pointer, ctx context.Context, req unsafe.Pointer) (_slice, error)
