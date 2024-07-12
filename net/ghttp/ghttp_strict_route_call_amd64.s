// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// ctx It is an interface type that will pass two pointers in
// fn.data = &closure         It is a closure structure, similar to the following
// type closure struct{
//      Fn unsafe.Pointer     What is saved here is the address of the actual function
//      captured variables... This is the parameter field captured by the closure function, which may be of pointer or value type
//      // A int
//      // B *int
//      // others...
// }
// The convention for calling closure functions in Go language is usually to store&closure in the dx register
// func request(fn {typ, data unsafe.Pointer}, ctx {typ, data unsafe.Pointer},req *HelloReq) (*HelloRes,error)
TEXT  ·doAnyCallRequest_with_res_err(SB), $24
    MOVQ fn_data+8(FP),DX

    MOVQ ctx_type+16(FP), AX
    MOVQ ctx_data+24(FP), BX
    MOVQ req+32(FP), CX

    CALL 0(DX)

    MOVQ AX, res+40(FP)
    MOVQ BX, err_type+48(FP)
    MOVQ CX, err_data+56(FP)
    RET

// func request(fn {typ, data unsafe.Pointer}, ctx {typ, data unsafe.Pointer},req *HelloReq) error
TEXT  ·doAnyCallRequest_with_err(SB), $16
    MOVQ fn_data+8(FP),DX

    MOVQ ctx_type+16(FP), AX
    MOVQ ctx_data+24(FP), BX
    MOVQ req+32(FP), CX

    CALL 0(DX)

    MOVQ BX, err_type+40(FP)
    MOVQ CX, err_data+48(FP)
    RET

// type slice struct{
//      ptr unsafe.Pointer
//      len int
//      cap int
// }
// func request(fn {typ, data unsafe.Pointer}, ctx {typ, data unsafe.Pointer},req *HelloReq) ([]HelloRes,error)
TEXT  ·doAnyCallRequest_with_sliceRes_err(SB), $40
    // Save the closure object to dx
    MOVQ fn_data+8(FP),DX
    // Passing ctx parameters
    MOVQ ctx_type+16(FP), AX
    MOVQ ctx_data+24(FP), BX
    // Passing req parameters
    MOVQ req+32(FP), CX
    // Call the closure function
    CALL 0(DX)
    // Set return value
    MOVQ AX, slice_ptr+40(FP)
    MOVQ BX, slice_len+48(FP)
    MOVQ CX, slice_cap+56(FP)
    MOVQ DI, err_type+64(FP)
    MOVQ SI, err_data+72(FP)
    RET

