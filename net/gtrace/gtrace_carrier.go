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

type Carrier struct {
	data map[string]interface{}
}

func NewCarrier(data ...map[string]interface{}) *Carrier {
	carrier := &Carrier{}
	if len(data) > 0 && data[0] != nil {
		carrier.data = data[0]
	} else {
		carrier.data = make(map[string]interface{})
	}
	return carrier
}

func (c *Carrier) Get(k string) string {
	return gconv.String(c.data[k])
}

func (c *Carrier) Set(k, v string) {
	c.data[k] = v
}

func (c *Carrier) MustMarshal() []byte {
	b, err := json.Marshal(c.data)
	if err != nil {
		panic(err)
	}
	return b
}

func (c *Carrier) String() string {
	return string(c.MustMarshal())
}

func (c *Carrier) UnmarshalJSON(b []byte) error {
	carrier := NewCarrier(nil)
	return json.Unmarshal(b, carrier.data)
}
