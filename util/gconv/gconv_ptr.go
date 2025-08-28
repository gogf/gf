// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv

// PtrAny creates and returns an any pointer variable to this value.
func PtrAny(any any) *any {
	return &any
}

// PtrString creates and returns a string pointer variable to this value.
func PtrString(any any) *string {
	v := String(any)
	return &v
}

// PtrBool creates and returns a bool pointer variable to this value.
func PtrBool(any any) *bool {
	v := Bool(any)
	return &v
}

// PtrInt creates and returns an int pointer variable to this value.
func PtrInt(any any) *int {
	v := Int(any)
	return &v
}

// PtrInt8 creates and returns an int8 pointer variable to this value.
func PtrInt8(any any) *int8 {
	v := Int8(any)
	return &v
}

// PtrInt16 creates and returns an int16 pointer variable to this value.
func PtrInt16(any any) *int16 {
	v := Int16(any)
	return &v
}

// PtrInt32 creates and returns an int32 pointer variable to this value.
func PtrInt32(any any) *int32 {
	v := Int32(any)
	return &v
}

// PtrInt64 creates and returns an int64 pointer variable to this value.
func PtrInt64(any any) *int64 {
	v := Int64(any)
	return &v
}

// PtrUint creates and returns an uint pointer variable to this value.
func PtrUint(any any) *uint {
	v := Uint(any)
	return &v
}

// PtrUint8 creates and returns an uint8 pointer variable to this value.
func PtrUint8(any any) *uint8 {
	v := Uint8(any)
	return &v
}

// PtrUint16 creates and returns an uint16 pointer variable to this value.
func PtrUint16(any any) *uint16 {
	v := Uint16(any)
	return &v
}

// PtrUint32 creates and returns an uint32 pointer variable to this value.
func PtrUint32(any any) *uint32 {
	v := Uint32(any)
	return &v
}

// PtrUint64 creates and returns an uint64 pointer variable to this value.
func PtrUint64(any any) *uint64 {
	v := Uint64(any)
	return &v
}

// PtrFloat32 creates and returns a float32 pointer variable to this value.
func PtrFloat32(any any) *float32 {
	v := Float32(any)
	return &v
}

// PtrFloat64 creates and returns a float64 pointer variable to this value.
func PtrFloat64(any any) *float64 {
	v := Float64(any)
	return &v
}
