// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gerror

import (
	"encoding/json"
)

// MarshalJSON implements the interface json.Marshaler for Error.
// It serializes the error using its string representation.
func (err *Error) MarshalJSON() ([]byte, error) {
	return json.Marshal(err.Error())
}
