// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

type (
	IBase interface {
		// sBase Init
		Init()
		// sBase Destory
		Destory()
		// baseDestory BeforeDestory
		BeforeDestory()
	}
)

var (
	localBase IBase
)

func Base() IBase {
	if localBase == nil {
		panic("implement not found for interface IBase, forgot register?")
	}
	return localBase
}

func RegisterBase(i IBase) {
	localBase = i
}
