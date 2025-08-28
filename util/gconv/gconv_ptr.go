// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv

// PtrAny creates and returns an any pointer variable to this value.
func PtrAny(anyInput any) *any {
	return &anyInput
}

// PtrString creates and returns a string pointer variable to this value.
func PtrString(anyInput any) *string {
	v := String(anyInput)
	return &v
}

// PtrBool creates and returns a bool pointer variable to this value.
func PtrBool(anyInput any) *bool {
	v := Bool(anyInput)
	return &v
}

// PtrInt creates and returns an int pointer variable to this value.
func PtrInt(anyInput any) *int {
	v := Int(anyInput)
	return &v
}

// PtrInt8 creates and returns an int8 pointer variable to this value.
func PtrInt8(anyInput any) *int8 {
	v := Int8(anyInput)
	return &v
}

// PtrInt16 creates and returns an int16 pointer variable to this value.
func PtrInt16(anyInput any) *int16 {
	v := Int16(anyInput)
	return &v
}

// PtrInt32 creates and returns an int32 pointer variable to this value.
func PtrInt32(anyInput any) *int32 {
	v := Int32(anyInput)
	return &v
}

// PtrInt64 creates and returns an int64 pointer variable to this value.
func PtrInt64(anyInput any) *int64 {
	v := Int64(anyInput)
	return &v
}

// PtrUint creates and returns an uint pointer variable to this value.
func PtrUint(anyInput any) *uint {
	v := Uint(anyInput)
	return &v
}

// PtrUint8 creates and returns an uint8 pointer variable to this value.
func PtrUint8(anyInput any) *uint8 {
	v := Uint8(anyInput)
	return &v
}

// PtrUint16 creates and returns an uint16 pointer variable to this value.
func PtrUint16(anyInput any) *uint16 {
	v := Uint16(anyInput)
	return &v
}

// PtrUint32 creates and returns an uint32 pointer variable to this value.
func PtrUint32(anyInput any) *uint32 {
	v := Uint32(anyInput)
	return &v
}

// PtrUint64 creates and returns an uint64 pointer variable to this value.
func PtrUint64(anyInput any) *uint64 {
	v := Uint64(anyInput)
	return &v
}

// PtrFloat32 creates and returns a float32 pointer variable to this value.
func PtrFloat32(anyInput any) *float32 {
	v := Float32(anyInput)
	return &v
}

// PtrFloat64 creates and returns a float64 pointer variable to this value.
func PtrFloat64(anyInput any) *float64 {
	v := Float64(anyInput)
	return &v
}
