// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gp.

package gparser

import (
	"github.com/gogf/gf/g/container/gvar"
	"github.com/gogf/gf/g/os/gtime"
	"time"
)

// Val returns the value.
func (p *Parser) Value() interface{} {
	return p.json.Value()
}

// Get returns value by specified <pattern>.
// It returns all values of current Json object, if <pattern> is empty or not specified.
// It returns nil if no value found by <pattern>.
//
// We can also access slice item by its index number in <pattern>,
// eg: "items.name.first", "list.10".
//
// It returns a default value specified by <def> if value for <pattern> is not found.
func (p *Parser) Get(pattern string, def...interface{}) interface{} {
	return p.json.Get(pattern, def...)
}

// GetVar returns a *gvar.Var with value by given <pattern>.
func (p *Parser) GetVar(pattern string, def...interface{}) *gvar.Var {
	return p.json.GetVar(pattern, def...)
}

// GetMap gets the value by specified <pattern>,
// and converts it to map[string]interface{}.
func (p *Parser) GetMap(pattern string, def...interface{}) map[string]interface{} {
    return p.json.GetMap(pattern, def...)
}

// GetArray gets the value by specified <pattern>,
// and converts it to a slice of []interface{}.
func (p *Parser) GetArray(pattern string, def...interface{}) []interface{} {
    return p.json.GetArray(pattern, def...)
}

// GetString gets the value by specified <pattern>,
// and converts it to string.
func (p *Parser) GetString(pattern string, def...interface{}) string {
    return p.json.GetString(pattern, def...)
}

// GetBool gets the value by specified <pattern>,
// and converts it to bool.
// It returns false when value is: "", 0, false, off, nil;
// or returns true instead.
func (p *Parser) GetBool(pattern string, def...interface{}) bool {
    return p.json.GetBool(pattern, def...)
}

func (p *Parser) GetInt(pattern string, def...interface{}) int {
    return p.json.GetInt(pattern, def...)
}

func (p *Parser) GetInt8(pattern string, def...interface{}) int8 {
    return p.json.GetInt8(pattern, def...)
}

func (p *Parser) GetInt16(pattern string, def...interface{}) int16 {
    return p.json.GetInt16(pattern, def...)
}

func (p *Parser) GetInt32(pattern string, def...interface{}) int32 {
    return p.json.GetInt32(pattern, def...)
}

func (p *Parser) GetInt64(pattern string, def...interface{}) int64 {
    return p.json.GetInt64(pattern, def...)
}

func (p *Parser) GetInts(pattern string, def...interface{}) []int {
    return p.json.GetInts(pattern, def...)
}

func (p *Parser) GetUint(pattern string, def...interface{}) uint {
    return p.json.GetUint(pattern, def...)
}

func (p *Parser) GetUint8(pattern string, def...interface{}) uint8 {
    return p.json.GetUint8(pattern, def...)
}

func (p *Parser) GetUint16(pattern string, def...interface{}) uint16 {
    return p.json.GetUint16(pattern, def...)
}

func (p *Parser) GetUint32(pattern string, def...interface{}) uint32 {
    return p.json.GetUint32(pattern, def...)
}

func (p *Parser) GetUint64(pattern string, def...interface{}) uint64 {
    return p.json.GetUint64(pattern, def...)
}

func (p *Parser) GetFloat32(pattern string, def...interface{}) float32 {
    return p.json.GetFloat32(pattern, def...)
}

func (p *Parser) GetFloat64(pattern string, def...interface{}) float64 {
    return p.json.GetFloat64(pattern, def...)
}

func (p *Parser) GetFloats(pattern string, def...interface{}) []float64 {
    return p.json.GetFloats(pattern, def...)
}

// GetStrings gets the value by specified <pattern>,
// and converts it to a slice of []string.
func (p *Parser) GetStrings(pattern string, def...interface{}) []string {
	return p.json.GetStrings(pattern, def...)
}

func (p *Parser) GetInterfaces(pattern string, def...interface{}) []interface{} {
	return p.json.GetInterfaces(pattern, def...)
}

func (p *Parser) GetTime(pattern string, format...string) time.Time {
	return p.json.GetTime(pattern, format...)
}

func (p *Parser) GetDuration(pattern string, def...interface{}) time.Duration {
	return p.json.GetDuration(pattern, def...)
}

func (p *Parser) GetGTime(pattern string, format...string) *gtime.Time {
	return p.json.GetGTime(pattern, format...)
}

// GetToVar gets the value by specified <pattern>,
// and converts it to specified golang variable <v>.
// The <v> should be a pointer type.
func (p *Parser) GetToVar(pattern string, pointer interface{}) error {
	return p.json.GetToVar(pattern, pointer)
}

// GetToStruct gets the value by specified <pattern>,
// and converts it to specified object <pointer>.
// The <pointer> should be the pointer to a struct.
func (p *Parser) GetToStruct(pattern string, pointer interface{}) error {
    return p.json.GetToStruct(pattern, pointer)
}

// Set sets value with specified <pattern>.
// It supports hierarchical data access by char separator, which is '.' in default.
func (p *Parser) Set(pattern string, value interface{}) error {
    return p.json.Set(pattern, value)
}

// Len returns the length/size of the value by specified <pattern>.
// The target value by <pattern> should be type of slice or map.
// It returns -1 if the target value is not found, or its type is invalid.
func (p *Parser) Len(pattern string) int {
    return p.json.Len(pattern)
}

// Append appends value to the value by specified <pattern>.
// The target value by <pattern> should be type of slice.
func (p *Parser) Append(pattern string, value interface{}) error {
    return p.json.Append(pattern, value)
}

// Remove deletes value with specified <pattern>.
// It supports hierarchical data access by char separator, which is '.' in default.
func (p *Parser) Remove(pattern string) error {
    return p.json.Remove(pattern)
}

// ToMap converts current object values to map[string]interface{}.
// It returns nil if fails.
func (p *Parser) ToMap() map[string]interface{} {
    return p.json.ToMap()
}

// ToArray converts current object values to []interface{}.
// It returns nil if fails.
func (p *Parser) ToArray() []interface{} {
    return p.json.ToArray()
}

// ToStruct converts current Json object to specified object.
// The <objPointer> should be a pointer type.
func (p *Parser) ToStruct(pointer interface{}) error {
    return p.json.ToStruct(pointer)
}

// Dump prints current Json object with more manually readable.
func (p *Parser) Dump() error {
	return p.json.Dump()
}