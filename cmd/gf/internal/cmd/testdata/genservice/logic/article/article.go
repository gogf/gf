// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package article

import (
	"context"
	"go/ast"
	t "time"

	"github.com/gogf/gf/cmd/gf/v2/internal/cmd/testdata/genservice/service"
	gdbalias "github.com/gogf/gf/v2/database/gdb"
)

type sArticle struct {
}

func init() {
	service.RegisterArticle(&sArticle{})
}

// Get article details
func (s *sArticle) Get(ctx context.Context, id uint) (info struct{}, err error) {
	return struct{}{}, err
}

// Create
/**
 * create an article.
 * @author oldme
 */
func (s *sArticle) Create(ctx context.Context, info struct{}) (id uint, err error) {
	// Use time package to test alias import.
	t.Now()
	return id, err
}

func (s *sArticle) A1o2(ctx context.Context, str string, a, b *ast.GoStmt) error {
	return nil
}

func (s *sArticle) B_2(ctx context.Context, db gdbalias.Raw) (err error) {
	return nil
}
