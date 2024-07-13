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
TEXT  ·doAnyCallRequest_with_res_err(SB), $80
    MOVQ AX,48(SP)
    MOVQ BX,56(SP)
    MOVQ CX,64(SP)
    MOVQ DX,72(SP)

    MOVQ fn_data+8(FP),DX

    MOVQ ctx_type+16(FP), AX
    MOVQ ctx_data+24(FP), BX
    MOVQ req+32(FP), CX
    // Passing parameters
    MOVQ AX,0(SP)
    MOVQ BX,8(SP)
    MOVQ CX,16(SP)
    CALL 0(DX)
    // Set return value
    MOVQ AX, res+40(FP)
    MOVQ BX, err_type+48(FP)
    MOVQ CX, err_data+56(FP)

    MOVQ 48(SP),AX
    MOVQ 56(SP),BX
    MOVQ 64(SP),CX
    MOVQ 72(SP),DX
    RET

// func request(fn {typ, data unsafe.Pointer}, ctx {typ, data unsafe.Pointer},req *HelloReq) error
TEXT  ·doAnyCallRequest_with_err(SB), $72
    MOVQ AX,48(SP)
    MOVQ BX,56(SP)
    MOVQ CX,64(SP)
    MOVQ DX,72(SP)

    MOVQ fn_data+8(FP),DX

    MOVQ ctx_type+16(FP), AX
    MOVQ ctx_data+24(FP), BX
    MOVQ req+32(FP), CX
    // Passing parameters
    MOVQ AX,0(SP)
    MOVQ BX,8(SP)
    MOVQ CX,16(SP)
    CALL 0(DX)
    // Set return value
    MOVQ AX, err_type+40(FP)
    MOVQ BX, err_data+48(FP)

    MOVQ 48(SP),AX
    MOVQ 56(SP),BX
    MOVQ 64(SP),CX
    MOVQ 72(SP),DX
    RET

// type slice struct{
//      ptr unsafe.Pointer
//      len int
//      cap int
// }
// func request(fn {typ, data unsafe.Pointer}, ctx {typ, data unsafe.Pointer},req *HelloReq) ([]HelloRes,error)
TEXT  ·doAnyCallRequest_with_sliceRes_err(SB), $96
    MOVQ AX,64(SP)
    MOVQ BX,72(SP)
    MOVQ CX,80(SP)
    MOVQ DX,88(SP)

    // Save the closure object to dx
    MOVQ fn_data+8(FP),DX
    // Passing ctx parameters
    MOVQ ctx_type+16(FP), AX
    MOVQ ctx_data+24(FP), BX
    // Passing req parameters
    MOVQ req+32(FP), CX
    // Passing parameters
    MOVQ AX,0(SP)
    MOVQ BX,8(SP)
    MOVQ CX,16(SP)
    // Call the closure function
    CALL 0(DX)
    // Set return value
    MOVQ AX, slice_ptr+40(FP)
    MOVQ BX, slice_len+48(FP)
    MOVQ CX, slice_cap+56(FP)
    MOVQ DI, err_type+64(FP)
    MOVQ SI, err_data+72(FP)

    MOVQ 64(SP),AX
    MOVQ 72(SP),BX
    MOVQ 80(SP),CX
    MOVQ 88(SP),DX
    RET

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// func(fn,obj unsafe.Pointer, ctx context.Context, req unsafe.Pointer) (unsafe.Pointer, error)
TEXT  ·doMethodCallRequest_with_res_err(SB), $96
    MOVQ AX,56(SP)
    MOVQ BX,64(SP)
    MOVQ CX,72(SP)
    MOVQ DI,80(SP)
    MOVQ SI,88(SP)

    MOVQ fn+0(FP),SI

    MOVQ obj+8(FP),AX
    MOVQ ctx_typ+16(FP),BX
    MOVQ ctx_data+24(FP),CX
    MOVQ req+32(FP),DI
    // Passing parameters
    MOVQ AX, 0(SP)
    MOVQ BX, 8(SP)
    MOVQ CX, 16(SP)
    MOVQ DI, 24(SP)
    CALL SI
    // Set return value
    MOVQ AX,res+40(FP)
    MOVQ BX,err_typ+48(FP)
    MOVQ CX,err_data+56(FP)

    MOVQ 56(SP),AX
    MOVQ 64(SP),BX
    MOVQ 72(SP),CX
    MOVQ 80(SP),DI
    MOVQ 88(SP),SI
    RET

// func(fn,obj unsafe.Pointer, ctx context.Context, req unsafe.Pointer) error
TEXT  ·doMethodCallRequest_with_err(SB), $96
    MOVQ AX,56(SP)
    MOVQ BX,64(SP)
    MOVQ CX,72(SP)
    MOVQ DI,80(SP)
    MOVQ SI,88(SP)

     MOVQ fn+0(FP),SI

     MOVQ obj+8(FP),AX
     MOVQ ctx_typ+16(FP),BX
     MOVQ ctx_data+24(FP),CX
     MOVQ req+32(FP),DI
    // Passing parameters
     MOVQ AX, 0(SP)
     MOVQ BX, 8(SP)
     MOVQ CX, 16(SP)
     MOVQ DI, 24(SP)
     CALL SI
    // Set return value
     MOVQ AX,err_typ+40(FP)
     MOVQ BX,err_data+48(FP)

     MOVQ 56(SP),AX
     MOVQ 64(SP),BX
     MOVQ 72(SP),CX
     MOVQ 80(SP),DI
     MOVQ 88(SP),SI
     RET


// func(fn,obj unsafe.Pointer, ctx context.Context, req unsafe.Pointer) (_slice, error)
TEXT  ·doMethodCallRequest_with_sliceRes_err(SB), $112
    MOVQ AX,72(SP)
    MOVQ BX,80(SP)
    MOVQ CX,88(SP)
    MOVQ DI,96(SP)
    MOVQ SI,104(SP)

    MOVQ fn+0(FP),SI

    MOVQ obj+8(FP),AX
    MOVQ ctx_typ+16(FP),BX
    MOVQ ctx_data+24(FP),CX
    MOVQ req+32(FP),DI
    // Passing parameters
    MOVQ AX, 0(SP)
    MOVQ BX, 8(SP)
    MOVQ CX, 16(SP)
    MOVQ DI, 24(SP)
    CALL SI
    // Set return value
    MOVQ AX, slice_ptr+40(FP)
    MOVQ BX, slice_len+48(FP)
    MOVQ CX, slice_cap+56(FP)
    MOVQ DI, err_type+64(FP)
    MOVQ SI, err_data+72(FP)

    MOVQ 72(SP),AX
    MOVQ 80(SP),BX
    MOVQ 88(SP),CX
    MOVQ 96(SP),DI
    MOVQ 104(SP),SI
    RET
