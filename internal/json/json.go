// Copyright 2020 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package json provides json operations wrapping ignoring stdlib or third-party lib json.
package json

import (
	json2 "encoding/json"
	"github.com/json-iterator/go"
	"io"
)

// ConfigCompatibleWithStandardLibrary tries to be 50% compatible
// with standard library behavior.
// 50% - -!
var json = jsoniter.ConfigCompatibleWithStandardLibrary

// Marshal adapts to json/encoding Marshal API.
//
// Marshal returns the JSON encoding of v, adapts to json/encoding Marshal API
// Refer to https://godoc.org/encoding/json#Marshal for more information.
func Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

// MarshalToString convenient method to write as string instead of []byte.
func MarshalToString(v interface{}) (string, error) {
	return json.MarshalToString(v)
}

// MarshalIndent same as json.MarshalIndent. Prefix is not supported.
func MarshalIndent(v interface{}, prefix, indent string) ([]byte, error) {
	return json2.MarshalIndent(v, prefix, indent)
}

// UnmarshalFromString is a convenient method to read from string instead of []byte.
func UnmarshalFromString(str string, v interface{}) error {
	return json.UnmarshalFromString(str, v)
}

// Unmarshal adapts to json/encoding Unmarshal API
//
// Unmarshal parses the JSON-encoded data and stores the result in the value pointed to by v.
// Refer to https://godoc.org/encoding/json#Unmarshal for more information.
func Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

// NewEncoder same as json.NewEncoder
func NewEncoder(writer io.Writer) *json2.Encoder {
	return json2.NewEncoder(writer)
}

// NewDecoder adapts to json/stream NewDecoder API.
//
// NewDecoder returns a new decoder that reads from r.
//
// Instead of a json/encoding Decoder, an Decoder is returned
// Refer to https://godoc.org/encoding/json#NewDecoder for more information.
func NewDecoder(reader io.Reader) *json2.Decoder {
	return json2.NewDecoder(reader)
}

// Valid reports whether data is a valid JSON encoding.
func Valid(data []byte) bool {
	return json2.Valid(data)
}
