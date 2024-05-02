// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
	"go/ast"

	gdbalias "github.com/gogf/gf/v2/database/gdb"
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
		A1o2(ctx context.Context, str string, a *ast.GoStmt, b *ast.GoStmt) error
		B_2(ctx context.Context, db gdbalias.Raw) (err error)
		// T1 random comment
		T1(ctx context.Context, id uint, id2 uint) (gdb gdbas.Model, err error)
		// T3
		/**
		 * random comment @*4213hHY1&%##%><<Y
		 * @param b
		 * @return c, d
		 * @return err
		 * @author oldme
		 */
		T3(ctx context.Context, b *gdbas.Model) (c *gdbas.Model, d *gdbas.Model, err error)
		// func (s *sArticle) T4(i interface{}) interface{}
		// # $ % ^ & * ( ) _ + - = { } | [ ] \ : " ; ' < > ? , . /
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
