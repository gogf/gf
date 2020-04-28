package main

import (
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

// ListsInput ListsInput
type ListsInput struct {
	Page  int `v:"min:1#page不能小于1"`
	Limit int `v:"between:10,100#limit必须是10-100之间的整数"`
	Search
}

type ListRequest struct {
	ListsInput
}

// Search Search
type Search struct {
	GoodsName    string
	GoodsSource  string
	FreeShipping string
	MarketPrice  struct {
		Min       int
		Max       int
		IsInclude bool
	}
	GuidePrice struct {
		Min       int
		Max       int
		IsInclude bool
	}
	AgreementPrice struct {
		Min       int
		Max       int
		IsInclude bool
	}
	NormalProfitMargin struct {
		Min       int
		Max       int
		IsInclude bool
	}
	ActivityProfitMargin struct {
		Min       int
		Max       int
		IsInclude bool
	}
	Sales struct {
		Cycle    string
		Quantity struct {
			Min       int
			Max       int
			IsInclude bool
		}
	}
	SalesReturn struct {
		Cycle    string
		Quantity struct {
			Min       int
			Max       int
			IsInclude bool
		}
	}
	ChooseGoods struct {
		Cycle    string
		Quantity struct {
			Min       int
			Max       int
			IsInclude bool
		}
	}
}

// Lists Lists
func Lists(r *ghttp.Request) {
	var params *ListRequest
	if err := r.Parse(&params); err != nil {
		r.Response.WriteExit(gerror.Stack(err))
	}

	r.Response.Write(params)
}

func main() {
	s := g.Server()
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.POST("/test", Lists)
	})
	s.SetPort(8199)
	s.Run()
}
