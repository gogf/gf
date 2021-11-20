// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import (
	"github.com/gogf/gf/v2/net/ghttp/internal/client"
)

type (
	Client            = client.Client
	ClientResponse    = client.Response
	ClientHandlerFunc = client.HandlerFunc
)

// NewClient creates and returns a new HTTP client object.
func NewClient() *Client {
	return client.New()
}
