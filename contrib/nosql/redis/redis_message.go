// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package redis

// Message received as result of a PUBLISH command issued by another client.
type Message struct {
	Channel      string
	Pattern      string
	Payload      string
	PayloadSlice []string
}
