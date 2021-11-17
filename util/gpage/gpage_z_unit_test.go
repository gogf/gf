// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go -bench=".*"

package gpage_test

import (
	"testing"

	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gpage"
)

func Test_New(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		page := gpage.New(9, 2, 1, `/user/list?page={.page}`)
		t.Assert(page.TotalSize, 9)
		t.Assert(page.TotalPage, 5)
		t.Assert(page.CurrentPage, 1)
	})
	gtest.C(t, func(t *gtest.T) {
		page := gpage.New(9, 2, 0, `/user/list?page={.page}`)
		t.Assert(page.TotalSize, 9)
		t.Assert(page.TotalPage, 5)
		t.Assert(page.CurrentPage, 1)
	})
}

func Test_Basic(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		page := gpage.New(9, 2, 1, `/user/list?page={.page}`)
		t.Assert(page.NextPage(), `<a class="GPageLink" href="/user/list?page=2" title="">></a>`)
		t.Assert(page.PrevPage(), `<span class="GPageSpan"><</span>`)
		t.Assert(page.FirstPage(), `<span class="GPageSpan">|<</span>`)
		t.Assert(page.LastPage(), `<a class="GPageLink" href="/user/list?page=5" title="">>|</a>`)
		t.Assert(page.PageBar(), `<span class="GPageSpan">1</span><a class="GPageLink" href="/user/list?page=2" title="2">2</a><a class="GPageLink" href="/user/list?page=3" title="3">3</a><a class="GPageLink" href="/user/list?page=4" title="4">4</a><a class="GPageLink" href="/user/list?page=5" title="5">5</a>`)
	})

	gtest.C(t, func(t *gtest.T) {
		page := gpage.New(9, 2, 3, `/user/list?page={.page}`)
		t.Assert(page.NextPage(), `<a class="GPageLink" href="/user/list?page=4" title="">></a>`)
		t.Assert(page.PrevPage(), `<a class="GPageLink" href="/user/list?page=2" title=""><</a>`)
		t.Assert(page.FirstPage(), `<a class="GPageLink" href="/user/list?page=1" title="">|<</a>`)
		t.Assert(page.LastPage(), `<a class="GPageLink" href="/user/list?page=5" title="">>|</a>`)
		t.Assert(page.PageBar(), `<a class="GPageLink" href="/user/list?page=1" title="1">1</a><a class="GPageLink" href="/user/list?page=2" title="2">2</a><span class="GPageSpan">3</span><a class="GPageLink" href="/user/list?page=4" title="4">4</a><a class="GPageLink" href="/user/list?page=5" title="5">5</a>`)
	})

	gtest.C(t, func(t *gtest.T) {
		page := gpage.New(9, 2, 5, `/user/list?page={.page}`)
		t.Assert(page.NextPage(), `<span class="GPageSpan">></span>`)
		t.Assert(page.PrevPage(), `<a class="GPageLink" href="/user/list?page=4" title=""><</a>`)
		t.Assert(page.FirstPage(), `<a class="GPageLink" href="/user/list?page=1" title="">|<</a>`)
		t.Assert(page.LastPage(), `<span class="GPageSpan">>|</span>`)
		t.Assert(page.PageBar(), `<a class="GPageLink" href="/user/list?page=1" title="1">1</a><a class="GPageLink" href="/user/list?page=2" title="2">2</a><a class="GPageLink" href="/user/list?page=3" title="3">3</a><a class="GPageLink" href="/user/list?page=4" title="4">4</a><span class="GPageSpan">5</span>`)
	})
}

func Test_CustomTag(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		page := gpage.New(5, 1, 2, `/user/list/{.page}`)
		page.PrevPageTag = "《"
		page.NextPageTag = "》"
		page.FirstPageTag = "|《"
		page.LastPageTag = "》|"
		page.PrevBarTag = "《《"
		page.NextBarTag = "》》"
		t.Assert(page.NextPage(), `<a class="GPageLink" href="/user/list/3" title="">》</a>`)
		t.Assert(page.PrevPage(), `<a class="GPageLink" href="/user/list/1" title="">《</a>`)
		t.Assert(page.FirstPage(), `<a class="GPageLink" href="/user/list/1" title="">|《</a>`)
		t.Assert(page.LastPage(), `<a class="GPageLink" href="/user/list/5" title="">》|</a>`)
		t.Assert(page.PageBar(), `<a class="GPageLink" href="/user/list/1" title="1">1</a><span class="GPageSpan">2</span><a class="GPageLink" href="/user/list/3" title="3">3</a><a class="GPageLink" href="/user/list/4" title="4">4</a><a class="GPageLink" href="/user/list/5" title="5">5</a>`)
	})
}

func Test_CustomStyle(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		page := gpage.New(5, 1, 2, `/user/list/{.page}`)
		page.LinkStyle = "MyPageLink"
		page.SpanStyle = "MyPageSpan"
		page.SelectStyle = "MyPageSelect"
		t.Assert(page.NextPage(), `<a class="MyPageLink" href="/user/list/3" title="">></a>`)
		t.Assert(page.PrevPage(), `<a class="MyPageLink" href="/user/list/1" title=""><</a>`)
		t.Assert(page.FirstPage(), `<a class="MyPageLink" href="/user/list/1" title="">|<</a>`)
		t.Assert(page.LastPage(), `<a class="MyPageLink" href="/user/list/5" title="">>|</a>`)
		t.Assert(page.PageBar(), `<a class="MyPageLink" href="/user/list/1" title="1">1</a><span class="MyPageSpan">2</span><a class="MyPageLink" href="/user/list/3" title="3">3</a><a class="MyPageLink" href="/user/list/4" title="4">4</a><a class="MyPageLink" href="/user/list/5" title="5">5</a>`)
		t.Assert(page.SelectBar(), `<select name="MyPageSelect" onchange="window.location.href=this.value"><option value="/user/list/1">1</option><option value="/user/list/2" selected>2</option><option value="/user/list/3">3</option><option value="/user/list/4">4</option><option value="/user/list/5">5</option></select>`)
	})
}

func Test_Ajax(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		page := gpage.New(5, 1, 2, `/user/list/{.page}`)
		page.AjaxActionName = "LoadPage"
		t.Assert(page.NextPage(), `<a class="GPageLink" href="javascript:LoadPage('/user/list/3')" title="">></a>`)
		t.Assert(page.PrevPage(), `<a class="GPageLink" href="javascript:LoadPage('/user/list/1')" title=""><</a>`)
		t.Assert(page.FirstPage(), `<a class="GPageLink" href="javascript:LoadPage('/user/list/1')" title="">|<</a>`)
		t.Assert(page.LastPage(), `<a class="GPageLink" href="javascript:LoadPage('/user/list/5')" title="">>|</a>`)
		t.Assert(page.PageBar(), `<a class="GPageLink" href="javascript:LoadPage('/user/list/1')" title="1">1</a><span class="GPageSpan">2</span><a class="GPageLink" href="javascript:LoadPage('/user/list/3')" title="3">3</a><a class="GPageLink" href="javascript:LoadPage('/user/list/4')" title="4">4</a><a class="GPageLink" href="javascript:LoadPage('/user/list/5')" title="5">5</a>`)
	})
}

func Test_PredefinedContent(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		page := gpage.New(5, 1, 2, `/user/list/{.page}`)
		page.AjaxActionName = "LoadPage"
		t.Assert(page.GetContent(1), `<a class="GPageLink" href="javascript:LoadPage('/user/list/1')" title="">上一页</a> <span class="current">2</span> <a class="GPageLink" href="javascript:LoadPage('/user/list/3')" title="">下一页</a>`)
		t.Assert(page.GetContent(2), `<a class="GPageLink" href="javascript:LoadPage('/user/list/1')" title="">首页</a><a class="GPageLink" href="javascript:LoadPage('/user/list/1')" title=""><<上一页</a><span class="current">[第2页]</span><a class="GPageLink" href="javascript:LoadPage('/user/list/3')" title="">下一页>></a><a class="GPageLink" href="javascript:LoadPage('/user/list/5')" title="">尾页</a>第<select name="GPageSelect" onchange="window.location.href=this.value"><option value="/user/list/1">1</option><option value="/user/list/2" selected>2</option><option value="/user/list/3">3</option><option value="/user/list/4">4</option><option value="/user/list/5">5</option></select>页`)
		t.Assert(page.GetContent(3), `<a class="GPageLink" href="javascript:LoadPage('/user/list/1')" title="">首页</a><a class="GPageLink" href="javascript:LoadPage('/user/list/1')" title="">上一页</a><a class="GPageLink" href="javascript:LoadPage('/user/list/1')" title="1">1</a><span class="GPageSpan">2</span><a class="GPageLink" href="javascript:LoadPage('/user/list/3')" title="3">3</a><a class="GPageLink" href="javascript:LoadPage('/user/list/4')" title="4">4</a><a class="GPageLink" href="javascript:LoadPage('/user/list/5')" title="5">5</a><a class="GPageLink" href="javascript:LoadPage('/user/list/3')" title="">下一页</a><a class="GPageLink" href="javascript:LoadPage('/user/list/5')" title="">尾页</a><span>当前页2/5</span> <span>共5条</span>`)
		t.Assert(page.GetContent(4), `<a class="GPageLink" href="javascript:LoadPage('/user/list/1')" title="">首页</a><a class="GPageLink" href="javascript:LoadPage('/user/list/1')" title="">上一页</a><a class="GPageLink" href="javascript:LoadPage('/user/list/1')" title="1">1</a><span class="GPageSpan">2</span><a class="GPageLink" href="javascript:LoadPage('/user/list/3')" title="3">3</a><a class="GPageLink" href="javascript:LoadPage('/user/list/4')" title="4">4</a><a class="GPageLink" href="javascript:LoadPage('/user/list/5')" title="5">5</a><a class="GPageLink" href="javascript:LoadPage('/user/list/3')" title="">下一页</a><a class="GPageLink" href="javascript:LoadPage('/user/list/5')" title="">尾页</a>`)
		t.Assert(page.GetContent(5), ``)
	})
}
