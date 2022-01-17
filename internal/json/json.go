// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package json provides json operations wrapping ignoring stdlib or third-party lib json.
package json

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/gogf/gf/v2/errors/gerror"
)

// Marshal adapts to json/encoding Marshal API.
//
// Marshal returns the JSON encoding of v, adapts to json/encoding Marshal API
// Refer to https://godoc.org/encoding/json#Marshal for more information.
func Marshal(v interface{}) (marshaledBytes []byte, err error) {
	marshaledBytes, err = json.Marshal(v)
	if err != nil {
		err = gerror.Wrap(err, `json.Marshal failed`)
	}
	return
}

// MarshalIndent same as json.MarshalIndent.
func MarshalIndent(v interface{}, prefix, indent string) (marshaledBytes []byte, err error) {
	marshaledBytes, err = json.MarshalIndent(v, prefix, indent)
	if err != nil {
		err = gerror.Wrap(err, `json.MarshalIndent failed`)
	}
	return
}

// Unmarshal adapts to json/encoding Unmarshal API
//
// Unmarshal parses the JSON-encoded data and stores the result in the value pointed to by v.
// Refer to https://godoc.org/encoding/json#Unmarshal for more information.
func Unmarshal(data []byte, v interface{}) (err error) {
	err = json.Unmarshal(data, v)
	if err != nil {
		err = gerror.Wrap(err, `json.Unmarshal failed`)
	}
	return
}

// UnmarshalUseNumber decodes the json data bytes to target interface using number option.
func UnmarshalUseNumber(data []byte, v interface{}) (err error) {
	decoder := NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	err = decoder.Decode(v)
	if err != nil {
		err = gerror.Wrap(err, `json.UnmarshalUseNumber failed`)
	}
	return
}

// NewEncoder same as json.NewEncoder
func NewEncoder(writer io.Writer) *json.Encoder {
	return json.NewEncoder(writer)
}

// NewDecoder adapts to json/stream NewDecoder API.
//
// NewDecoder returns a new decoder that reads from r.
//
// Instead of a json/encoding Decoder, a Decoder is returned
// Refer to https://godoc.org/encoding/json#NewDecoder for more information.
func NewDecoder(reader io.Reader) *json.Decoder {
	return json.NewDecoder(reader)
}

// Valid reports whether data is a valid JSON encoding.
func Valid(data []byte) bool {
	return json.Valid(data)
}
