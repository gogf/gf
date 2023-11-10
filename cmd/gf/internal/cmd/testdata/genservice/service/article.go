// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
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
