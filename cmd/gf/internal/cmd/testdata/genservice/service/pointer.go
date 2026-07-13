// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
)

type (
	IPointer interface {
		TestPointer(ctx context.Context, opts *PointerOptions) error
	}
)

var (
	localPointer IPointer
)

func Pointer() IPointer {
	if localPointer == nil {
		panic("implement not found for interface IPointer, forgot register?")
	}
	return localPointer
}

func RegisterPointer(i IPointer) {
	localPointer = i
}
