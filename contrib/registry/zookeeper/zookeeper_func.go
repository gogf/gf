// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package zookeeper

import (
	"encoding/json"
)

func unmarshal(data []byte) (c *Content, err error) {
	err = json.Unmarshal(data, &c)
	return
}

func marshal(c *Content) ([]byte, error) {
	return json.Marshal(c)
}
