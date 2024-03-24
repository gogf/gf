// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gmetric

import (
	"bytes"
	"fmt"
	"os"

	"github.com/gogf/gf/v2/internal/json"
	"github.com/gogf/gf/v2/os/gfile"
)

// Attributes is a slice of Attribute.
type Attributes []Attribute

// Attribute is the key-value pair item for Metric.
type Attribute interface {
	Key() string // The key for this attribute.
	Value() any  // The value for this attribute.
}

// AttributeKey is the attribute key.
type AttributeKey string

// Option holds the option for perform a metric operation.
type Option struct {
	// Attributes holds the dynamic key-value pair metadata.
	Attributes Attributes
}

// localAttribute implements interface Attribute.
type localAttribute struct {
	key   string
	value any
}

var (
	hostname    string
	processPath string
)

func init() {
	hostname, _ = os.Hostname()
	processPath = gfile.SelfPath()
}

// CommonAttributes returns the common used attributes for an instrument.
func CommonAttributes() Attributes {
	return Attributes{
		NewAttribute(`os.host.name`, hostname),
		NewAttribute(`process.path`, processPath),
	}
}

// NewAttribute creates and returns an Attribute by given `key` and `value`.
func NewAttribute(key string, value any) Attribute {
	return &localAttribute{
		key:   key,
		value: value,
	}
}

// Key returns the key of the attribute.
func (l *localAttribute) Key() string {
	return l.key
}

// Value returns the value of the attribute.
func (l *localAttribute) Value() any {
	return l.value
}

// MarshalJSON implements the interface MarshalJSON for json.Marshal.
func (l *localAttribute) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`{"%s":%#v}`, l.key, l.value)), nil
}

// MarshalJSON implements the interface MarshalJSON for json.Marshal.
func (attrs Attributes) String() string {
	bs, _ := attrs.MarshalJSON()
	return string(bs)
}

// MarshalJSON implements the interface MarshalJSON for json.Marshal.
func (attrs Attributes) MarshalJSON() ([]byte, error) {
	var (
		bs     []byte
		err    error
		buffer = bytes.NewBuffer(nil)
	)
	buffer.WriteByte('[')
	for _, attr := range attrs {
		bs, err = json.Marshal(attr)
		if err != nil {
			return nil, err
		}
		if buffer.Len() > 1 {
			buffer.WriteByte(',')
		}
		buffer.Write(bs)
	}
	buffer.WriteByte(']')
	return buffer.Bytes(), nil
}
