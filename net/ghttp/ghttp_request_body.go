// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import (
	"bytes"
	"io/ioutil"
)

// BodyReadCloser implements the io.ReadCloser interface
// which is used for reading request body content multiple times.
type BodyReadCloser struct {
	*bytes.Reader
}

// RefillBody refills the request body object after read all of its content.
// It makes the request body reusable for next reading.
func (r *Request) RefillBody() {
	if r.bodyContent == nil {
		r.bodyContent, _ = ioutil.ReadAll(r.Body)
	}
	r.Body = &BodyReadCloser{
		bytes.NewReader(r.bodyContent),
	}
}

// Close implements the io.ReadCloser interface.
func (b *BodyReadCloser) Close() error {
	return nil
}
