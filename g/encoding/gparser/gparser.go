// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gp.

// Package gparser provides convenient API for accessing/converting variable and JSON/XML/YAML/TOML.
package gparser

import (
    "github.com/gogf/gf/g/encoding/gjson"
    "time"
)

type Parser struct {
    json *gjson.Json
}

// New creates a Parser object with any variable type of <data>,
// but <data> should be a map or slice for data access reason,
// or it will make no sense.
// The <unsafe> param specifies whether using this Parser object
// in un-concurrent-safe context, which is false in default.
func New(value interface{}, unsafe...bool) *Parser {
    return &Parser{gjson.New(value, unsafe...)}
}

// NewUnsafe creates a un-concurrent-safe Parser object.
func NewUnsafe (value...interface{}) *Parser {
    if len(value) > 0 {
        return &Parser{gjson.New(value[0], false)}
    }
    return &Parser{gjson.New(nil, false)}
}

// Load loads content from specified file <path>,
// and creates a Parser object from its content.
func Load (path string, unsafe...bool) (*Parser, error) {
    if j, e := gjson.Load(path, unsafe...); e == nil {
        return &Parser{j}, nil
    } else {
        return nil, e
    }
}

// LoadContent creates a Parser object from given content,
// it checks the data type of <content> automatically,
// supporting JSON, XML, YAML and TOML types of data.
func LoadContent (data []byte, unsafe...bool) (*Parser, error) {
    if j, e := gjson.LoadContent(data, unsafe...); e == nil {
        return &Parser{j}, nil
    } else {
        return nil, e
    }
}

// SetSplitChar sets the separator char for hierarchical data access.
func (p *Parser) SetSplitChar(char byte) {
    p.json.SetSplitChar(char)
}

// SetViolenceCheck enables/disables violence check for hierarchical data access.
func (p *Parser) SetViolenceCheck(check bool) {
    p.json.SetViolenceCheck(check)
}

// GetToVar gets the value by specified <pattern>,
// and converts it to specified golang variable <v>.
// The <v> should be a pointer type.
func (p *Parser) GetToVar(pattern string, v interface{}) error {
    return p.json.GetToVar(pattern, v)
}

// GetMap gets the value by specified <pattern>,
// and converts it to map[string]interface{}.
func (p *Parser) GetMap(pattern string) map[string]interface{} {
    return p.json.GetMap(pattern)
}

// GetArray gets the value by specified <pattern>,
// and converts it to a slice of []interface{}.
func (p *Parser) GetArray(pattern string) []interface{} {
    return p.json.GetArray(pattern)
}

// GetString gets the value by specified <pattern>,
// and converts it to string.
func (p *Parser) GetString(pattern string) string {
    return p.json.GetString(pattern)
}

// GetStrings gets the value by specified <pattern>,
// and converts it to a slice of []string.
func (p *Parser) GetStrings(pattern string) []string {
    return p.json.GetStrings(pattern)
}

func (p *Parser) GetInterfaces(pattern string) []interface{} {
    return p.json.GetInterfaces(pattern)
}

func (p *Parser) GetTime(pattern string, format ... string) time.Time {
    return p.json.GetTime(pattern, format...)
}

func (p *Parser) GetTimeDuration(pattern string) time.Duration {
    return p.json.GetTimeDuration(pattern)
}

// GetBool gets the value by specified <pattern>,
// and converts it to bool.
// It returns false when value is: "", 0, false, off, nil;
// or returns true instead.
func (p *Parser) GetBool(pattern string) bool {
    return p.json.GetBool(pattern)
}

func (p *Parser) GetInt(pattern string) int {
    return p.json.GetInt(pattern)
}

func (p *Parser) GetInt8(pattern string) int8 {
    return p.json.GetInt8(pattern)
}

func (p *Parser) GetInt16(pattern string) int16 {
    return p.json.GetInt16(pattern)
}

func (p *Parser) GetInt32(pattern string) int32 {
    return p.json.GetInt32(pattern)
}

func (p *Parser) GetInt64(pattern string) int64 {
    return p.json.GetInt64(pattern)
}

func (p *Parser) GetInts(pattern string) []int {
    return p.json.GetInts(pattern)
}

func (p *Parser) GetUint(pattern string) uint {
    return p.json.GetUint(pattern)
}

func (p *Parser) GetUint8(pattern string) uint8 {
    return p.json.GetUint8(pattern)
}

func (p *Parser) GetUint16(pattern string) uint16 {
    return p.json.GetUint16(pattern)
}

func (p *Parser) GetUint32(pattern string) uint32 {
    return p.json.GetUint32(pattern)
}

func (p *Parser) GetUint64(pattern string) uint64 {
    return p.json.GetUint64(pattern)
}

func (p *Parser) GetFloat32(pattern string) float32 {
    return p.json.GetFloat32(pattern)
}

func (p *Parser) GetFloat64(pattern string) float64 {
    return p.json.GetFloat64(pattern)
}

func (p *Parser) GetFloats(pattern string) []float64 {
    return p.json.GetFloats(pattern)
}

// GetToStruct gets the value by specified <pattern>,
// and converts it to specified object <objPointer>.
// The <objPointer> should be the pointer to an object.
func (p *Parser) GetToStruct(pattern string, objPointer interface{}) error {
    return p.json.GetToStruct(pattern, objPointer)
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

// Get returns value by specified <pattern>.
// It returns all values of current Json object, if <pattern> is empty or not specified.
// It returns nil if no value found by <pattern>.
//
// We can also access slice item by its index number in <pattern>,
// eg: "items.name.first", "list.10".
func (p *Parser) Get(pattern...string) interface{} {
    return p.json.Get(pattern...)
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

func (p *Parser) ToXml(rootTag...string) ([]byte, error) {
    return p.json.ToXml(rootTag...)
}

func (p *Parser) ToXmlIndent(rootTag...string) ([]byte, error) {
    return p.json.ToXmlIndent(rootTag...)
}

func (p *Parser) ToJson() ([]byte, error) {
    return p.json.ToJson()
}

func (p *Parser) ToJsonIndent() ([]byte, error) {
    return p.json.ToJsonIndent()
}

func (p *Parser) ToYaml() ([]byte, error) {
    return p.json.ToYaml()
}

func (p *Parser) ToToml() ([]byte, error) {
    return p.json.ToToml()
}

// Dump prints current Json object with more manually readable.
func (p *Parser) Dump() error {
    return p.json.Dump()
}

// ToStruct converts current Json object to specified object.
// The <objPointer> should be a pointer type.
func (p *Parser) ToStruct(o interface{}) error {
    return p.json.ToStruct(o)
}

func VarToXml(value interface{}, rootTag...string) ([]byte, error) {
    return New(value).ToXml(rootTag...)
}

func VarToXmlIndent(value interface{}, rootTag...string) ([]byte, error) {
    return New(value).ToXmlIndent(rootTag...)
}

func VarToJson(value interface{}) ([]byte, error) {
    return New(value).ToJson()
}

func VarToJsonIndent(value interface{}) ([]byte, error) {
    return New(value).ToJsonIndent()
}

func VarToYaml(value interface{}) ([]byte, error) {
    return New(value).ToYaml()
}

func VarToToml(value interface{}) ([]byte, error) {
    return New(value).ToToml()
}

func VarToStruct(value interface{}, obj interface{}) error {
    return New(value).ToStruct(obj)
}

