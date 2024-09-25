// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

type DefaultModelInterfaceImpl struct {
	*Model
}

func (m DefaultModelInterfaceImpl) setModel(model *Model) {
	m.Model = model
}

var (
	registerModelInterface = func(model *Model) ModelInterface {
		return DefaultModelInterfaceImpl{
			Model: model,
		}
	}
)

func RegisterModelInterface(fn func(model *Model) ModelInterface) {
	if fn == nil {
		return
	}
	registerModelInterface = fn
}
