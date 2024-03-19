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
// 检测ptr是不是*struct 或者 **struct []*struct *[]*struct
// 如果满足类型，就直接解引用返回
// 后续可以通过返回值来设置字段之类的操作
// 不管是一级还是二级
// 使用以下的方法调用都会导致空指针
// var t *T
// 1.checkValidStructPtr(t)  ===> panic
// var t **T
// 2.checkValidStructPtr(t)  ===> panic
// 只有以下几种情况能调用成功
// var t T ==>checkValidStructPtr(&t)   OK
// var t *T ==>checkValidStructPtr(&t)   OK
func checkValidStructPtr(ptr any) (reflect.Value, reflect.Kind, error) {
	srcVal, ok := ptr.(reflect.Value)
	if ok {
		// 用于标准路由
		srcTyp := srcVal.Type()
		// XXXReq(ctx,&LoginReq{})
		elem, err := setStruct(srcTyp, srcVal, 1)
		return elem, reflect.Struct, err
	} else {
		// 其他路由
		srcTyp := reflect.TypeOf(ptr)

		if srcTyp.Kind() != reflect.Ptr {
			return srcVal, 0, fmt.Errorf("传入的值不是指针类型==%v", srcTyp)
		}
		srcValof := reflect.ValueOf(ptr)
		// 判断是不是nil指针
		if srcValof.IsNil() {
			return srcVal, 0, fmt.Errorf("不能传入空指针，ptr=%v", srcValof)

		}

		elemTyp := srcTyp.Elem()
		// 解引用次数
		derefCount := 1
		// 判断是不是二级指针
		if elemTyp.Kind() == reflect.Ptr {
			elemTyp = elemTyp.Elem()
			derefCount++
		}
		// 走到这里类型没错，也不是空指针
		switch elemTyp.Kind() {
		case reflect.Slice, reflect.Array:
			elem, err := setStructSlice(elemTyp, srcValof, derefCount)
			return elem, elemTyp.Kind(), err

		case reflect.Struct:
			elem, err := setStruct(elemTyp, srcValof, derefCount)
			return elem, reflect.Struct, err

		default:

			return reflect.Value{}, 0, fmt.Errorf("不支持的类型%s, 目前只支持*struct或者**struct *[]struct *[]*struct", srcTyp)
		}
	}

}

func setStructSlice(elemTyp reflect.Type, srcValof reflect.Value, derefCount int) (reflect.Value, error) {
	// [] =>  *struct / struct
	// 切片解引用之后 可能是*struct / struct
	elemTyp = elemTyp.Elem()
	if elemTyp.Kind() == reflect.Ptr {
		elemTyp = elemTyp.Elem()
	}

	if elemTyp.Kind() != reflect.Struct {
		return srcValof, fmt.Errorf("不支持的类型%s, 目前只支持*struct或者**struct *[]struct *[]*struct", elemTyp)
	}

	return srcValof, nil
}

// derefCount 等于几级指针，一级指针就是1，二级指针就是2，需要做解引用
func setStruct(elemTyp reflect.Type, srcValof reflect.Value, derefCount int) (reflect.Value, error) {
	// 走到这里类型没错
	var elemPtr reflect.Value

	// 一级指针
	if derefCount == 1 {
		// 解引用
		elemPtr = srcValof.Elem()

	} else if derefCount == 2 {

		// 检测二级指针是不是nil的，如果是，则使用new赋值
		srcValof = srcValof.Elem()
		if srcValof.IsNil() {
			newTyp := reflect.New(elemTyp)

			srcValof.Set(newTyp)
		}
		elemPtr = srcValof.Elem()

	} else {
		return srcValof, fmt.Errorf("不支持的类型%s, 目前只支持*struct或者**struct *[]struct *[]*struct", elemTyp)
	}
	// 走到这里 指针已经解引用，可以设置值了

	return elemPtr, nil
}
