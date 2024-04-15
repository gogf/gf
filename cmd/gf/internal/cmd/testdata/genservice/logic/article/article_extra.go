// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package article

import (
	"context"

	gdbas "github.com/gogf/gf/v2/database/gdb"
)

func (s *sArticle) T1(ctx context.Context, id, id2 uint) (gdb gdbas.Model, err error) {
	g := gdbas.Model{}
	return g, err
}

func (s *sArticle) t2(ctx context.Context) (err error) {
	return nil
}

func (s *sArticle) T3(ctx context.Context, b *gdbas.Model) (c, d *gdbas.Model, err error) {
	return nil, nil, nil
}

func (s *sArticle) T4(i interface{}) interface{} {
	return nil
}
