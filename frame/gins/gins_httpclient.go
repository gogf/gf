// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gins

import (
	"fmt"

	"github.com/gogf/gf/v2/internal/instance"
	"github.com/gogf/gf/v2/net/gclient"
)

// HttpClient returns an instance of http client with specified name.
func HttpClient(name ...interface{}) *gclient.Client {
	var instanceKey = fmt.Sprintf("%s.%v", frameCoreComponentNameHttpClient, name)
	return instance.GetOrSetFuncLock(instanceKey, func() interface{} {
		return gclient.New()
	}).(*gclient.Client)
}
