// Copyright 2018-2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gvar provides an universal variable type, like generics.
package gvar

import (
	"github.com/gogf/gf/internal/intstore"
	"github.com/gogf/gf/os/gtime"
	"reflect"
	"time"
)

// Var is a universal variable interface, like generics.
type Var interface {
	IsEmpty() bool
	IsNil() bool

	Val() interface{}
	Set(value interface{}) (old interface{})
	Interface() interface{}
	String() string

	Bool() bool

	Int() int
	Int16() int16
	Int32() int32
	Int64() int64
	Int8() int8

	Uint() uint
	Uint16() uint16
	Uint32() uint32
	Uint64() uint64
	Uint8() uint8

	Float32() float32
	Float64() float64

	Ints() []int
	Uints() []uint
	Vars() []Var
	Bytes() []byte
	Floats() []float64
	Array() []interface{}
	Slice() []interface{}
	Strings() []string
	Interfaces() []interface{}

	Map(tags ...string) map[string]interface{}
	MapDeep(tags ...string) map[string]interface{}
	MapStrStr(tags ...string) map[string]string
	MapStrStrDeep(tags ...string) map[string]string
	MapStrVar(tags ...string) map[string]Var
	MapStrVarDeep(tags ...string) map[string]Var
	MapToMap(pointer interface{}, mapping ...map[string]string) (err error)
	MapToMapDeep(pointer interface{}, mapping ...map[string]string) (err error)
	MapToMaps(pointer interface{}, mapping ...map[string]string) (err error)
	MapToMapsDeep(pointer interface{}, mapping ...map[string]string) (err error)

	Scan(pointer interface{}, mapping ...map[string]string) error
	ScanDeep(pointer interface{}, mapping ...map[string]string) error
	Struct(pointer interface{}, mapping ...map[string]string) error
	StructDeep(pointer interface{}, mapping ...map[string]string) error
	Structs(pointer interface{}, mapping ...map[string]string) (err error)
	StructsDeep(pointer interface{}, mapping ...map[string]string) (err error)

	Time(format ...string) time.Time
	Duration() time.Duration
	GTime(format ...string) *gtime.Time

	MarshalJSON() ([]byte, error)
	UnmarshalJSON(b []byte) error
	UnmarshalValue(value interface{}) error
}

func init() {
	// Register the type of gvar.VarImp to local variable.
	intstore.ReflectTypeVarImp = reflect.TypeOf(VarImp{})
}
