// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package garray

type apiSliceInterface interface {
	Slice() []interface{}
}

type apiSliceInt interface {
	Slice() []int
}

type apiSliceString interface {
	Slice() []string
}
