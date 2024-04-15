// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
	"go/ast"

	gdbas "github.com/gogf/gf/v2/database/gdb"
)

type (
	IArticle interface {
		// Get article details
		Get(ctx context.Context, id uint) (info struct{}, err error)
		// Create
		/**
		 * create an article.
		 * @author oldme
		 */
		Create(ctx context.Context, info struct{}) (id uint, err error)
		A1o2(ctx context.Context, str string, a, b *ast.GoStmt) error
		T1(ctx context.Context, id, id2 uint) (gdb gdbas.Model, err error)
		T3(ctx context.Context, b *gdbas.Model) (c, d *gdbas.Model, err error)
		T4(i interface{}) interface{}
	}
)

var (
	localArticle IArticle
)

func Article() IArticle {
	if localArticle == nil {
		panic("implement not found for interface IArticle, forgot register?")
	}
	return localArticle
}

func RegisterArticle(i IArticle) {
	localArticle = i
}
