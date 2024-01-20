// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package article

import (
	"context"

	"github.com/gogf/gf/cmd/gf/v2/internal/cmd/testdata/genservice/service"
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
	return id, err
}
