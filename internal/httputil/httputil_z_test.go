// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package httputil

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildParams(t *testing.T) {
	// normal && special cases
	params := map[string]string{
		"code":  "x",
		"code1": "x&a=1", // for fix
		"code2": "x&a=111",
		"id":    "1+- ", // for fix
		"f":     "1#a=+- ",
		"v":     "",
		"n":     "null",
	}

	res1 := BuildParams(params)
	vs, err := url.ParseQuery(res1)
	assert.Truef(t, err == nil && len(params) == len(vs), "expected-len vs actual-len: [%d] vs [%d]", len(params), len(vs))
	for k := range vs {
		vv := vs.Get(k)
		_, ok := params[k]
		assert.Truef(t, ok && params[k] == vv, "[%s] expected vs actual: [%v] vs [%v]", k, params[k], vv)
	}
}
