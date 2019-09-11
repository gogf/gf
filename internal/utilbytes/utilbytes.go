// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package utibytes provides some bytes functions for internal usage.
package utilbytes

import (
	"bytes"
	"fmt"
)

func Export(b []byte) string {
	buffer := bytes.NewBuffer(nil)
	buffer.WriteString(`[]byte("`)
	for _, v := range b {
		fmt.Fprintf(buffer, `\x%02x`, v)
	}
	buffer.WriteString(`")`)
	return buffer.String()
}
