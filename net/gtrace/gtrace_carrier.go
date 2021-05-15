// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gtrace

import (
	"github.com/gogf/gf/internal/json"
	"github.com/gogf/gf/util/gconv"
)

// Carrier is the storage medium used by a TextMapPropagator.
type Carrier map[string]interface{}

func NewCarrier(data ...map[string]interface{}) Carrier {
	if len(data) > 0 && data[0] != nil {
		return data[0]
	} else {
		return make(map[string]interface{})
	}
}

// Get returns the value associated with the passed key.
func (c Carrier) Get(k string) string {
	return gconv.String(c[k])
}

// Set stores the key-value pair.
func (c Carrier) Set(k, v string) {
	c[k] = v
}

// Keys lists the keys stored in this carrier.
func (c Carrier) Keys() []string {
	keys := make([]string, 0, len(c))
	for k := range c {
		keys = append(keys, k)
	}
	return keys
}

func (c Carrier) MustMarshal() []byte {
	b, err := json.Marshal(c)
	if err != nil {
		panic(err)
	}
	return b
}

func (c Carrier) String() string {
	return string(c.MustMarshal())
}

func (c Carrier) UnmarshalJSON(b []byte) error {
	carrier := NewCarrier(nil)
	return json.UnmarshalUseNumber(b, carrier)
}
