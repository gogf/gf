// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gpage provides useful paging functionality for web pages.
package gpage

import (
	"fmt"
	"math"
	"net/url"
	"strings"

	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/text/gregex"
	"github.com/gogf/gf/text/gstr"
	"github.com/gogf/gf/util/gconv"
)

// Page is the pagination implementer.
type Page struct {
	UrlTemplate    string // Custom url template for page url producing.
	TotalSize      int    // Total size.
	TotalPage      int    // Total page, which is automatically calculated.
	CurrentPage    int    // Current page number >= 1.
	PageName       string // Page variable name. It's "page" in default.
	NextPageTag    string // Tag name for next p.
	PrevPageTag    string // Tag name for prev p.
	FirstPageTag   string // Tag name for first p.
	LastPageTag    string // Tag name for last p.
	PrevBar        string // Tag string for prev bar.
	NextBar        string // Tag string for next bar.
	PageBarNum     int    // Page bar number for displaying.
	AjaxActionName string // Ajax function name. Ajax is enabled if this attribute is not empty.
}

// 创建一个分页对象，输入参数分别为：
// 总数量、每页数量、当前页码、当前的URL(URI+QUERY)、(可选)路由规则(例如: /user/list/:page、/order/list/*page、/order/list/{page}.html)
func New(totalSize, pageSize, currentPage int, urlTemplate string) *Page {
	p := &Page{
		PageName:     "page",
		PrevPageTag:  "<",
		NextPageTag:  ">",
		FirstPageTag: "|<",
		LastPageTag:  ">|",
		PrevBar:      "<<",
		NextBar:      ">>",
		TotalSize:    totalSize,
		TotalPage:    int(math.Ceil(float64(totalSize) / float64(pageSize))),
		CurrentPage:  1,
		PageBarNum:   10,
		UrlTemplate:  urlTemplate,
	}
	if currentPage > 0 {
		p.CurrentPage = currentPage
	}
	return p
}

// 获取显示"下一页"的内容.
func (p *Page) NextPage(styles ...string) string {
	var curStyle, style string
	if len(styles) > 0 {
		curStyle = styles[0]
	}
	if len(styles) > 1 {
		style = styles[0]
	}
	if p.CurrentPage < p.TotalPage {
		return p.GetLink(p.GetUrl(p.CurrentPage+1), p.NextPageTag, "下一页", style)
	}
	return fmt.Sprintf(`<span class="%s">%s</span>`, curStyle, p.NextPageTag)
}

// 获取显示“上一页”的内容
func (p *Page) PrevPage(styles ...string) string {
	var curStyle, style string
	if len(styles) > 0 {
		curStyle = styles[0]
	}
	if len(styles) > 1 {
		style = styles[0]
	}
	if p.CurrentPage > 1 {
		return p.GetLink(p.GetUrl(p.CurrentPage-1), p.PrevPageTag, "上一页", style)
	}
	return fmt.Sprintf(`<span class="%s">%s</span>`, curStyle, p.PrevPageTag)
}

// 获取显示“首页”的代码
func (p *Page) FirstPage(styles ...string) string {
	var curStyle, style string
	if len(styles) > 0 {
		curStyle = styles[0]
	}
	if len(styles) > 1 {
		style = styles[0]
	}
	if p.CurrentPage == 1 {
		return fmt.Sprintf(`<span class="%s">%s</span>`, curStyle, p.FirstPageTag)
	}
	return p.GetLink(p.GetUrl(1), p.FirstPageTag, "第一页", style)
}

// 获取显示“尾页”的内容
func (p *Page) LastPage(styles ...string) string {
	var curStyle, style string
	if len(styles) > 0 {
		curStyle = styles[0]
	}
	if len(styles) > 1 {
		style = styles[0]
	}
	if p.CurrentPage == p.TotalPage {
		return fmt.Sprintf(`<span class="%s">%s</span>`, curStyle, p.LastPageTag)
	}
	return p.GetLink(p.GetUrl(p.TotalPage), p.LastPageTag, "最后页", style)
}

// 获得分页条列表内容
func (p *Page) PageBar(styles ...string) string {
	var curStyle, style string
	if len(styles) > 0 {
		curStyle = styles[0]
	}
	if len(styles) > 1 {
		style = styles[0]
	}
	plus := int(math.Ceil(float64(p.PageBarNum / 2)))
	if p.PageBarNum-plus+p.CurrentPage > p.TotalPage {
		plus = p.PageBarNum - p.TotalPage + p.CurrentPage
	}
	begin := p.CurrentPage - plus + 1
	if begin < 1 {
		begin = 1
	}
	ret := ""
	for i := begin; i < begin+p.PageBarNum; i++ {
		if i <= p.TotalPage {
			if i != p.CurrentPage {
				ret += p.GetLink(p.GetUrl(i), gconv.String(i), style, "")
			} else {
				ret += fmt.Sprintf(`<span class="%s">%d</span>`, curStyle, i)
			}
		} else {
			break
		}
	}
	return ret
}

// 获取基于select标签的显示跳转按钮的代码
func (p *Page) SelectBar() string {
	ret := `<select name="GPageSelect" onchange="window.location.href=this.value">`
	for i := 1; i <= p.TotalPage; i++ {
		if i == p.CurrentPage {
			ret += fmt.Sprintf(`<option value="%s" selected>%d</option>`, p.GetUrl(i), i)
		} else {
			ret += fmt.Sprintf(`<option value="%s">%d</option>`, p.GetUrl(i), i)
		}
	}
	ret += "</select>"
	return ret
}

// 预定义的分页显示风格内容
func (p *Page) GetContent(mode int) string {
	switch mode {
	case 1:
		p.NextPageTag = "下一页"
		p.PrevPageTag = "上一页"
		return fmt.Sprintf(
			`%s <span class="current">%d</span> %s`,
			p.PrevPage(),
			p.CurrentPage,
			p.NextPage(),
		)

	case 2:
		p.NextPageTag = "下一页>>"
		p.PrevPageTag = "<<上一页"
		p.FirstPageTag = "首页"
		p.LastPageTag = "尾页"
		return fmt.Sprintf(
			`%s%s<span class="current">[第%d页]</span>%s%s第%s页`,
			p.FirstPage(),
			p.PrevPage(),
			p.CurrentPage,
			p.NextPage(),
			p.LastPage(),
			p.SelectBar(),
		)

	case 3:
		p.NextPageTag = "下一页"
		p.PrevPageTag = "上一页"
		p.FirstPageTag = "首页"
		p.LastPageTag = "尾页"
		pageStr := p.FirstPage()
		pageStr += p.PrevPage()
		pageStr += p.PageBar("current")
		pageStr += p.NextPage()
		pageStr += p.LastPage()
		pageStr += fmt.Sprintf(
			`<span>当前页%d/%d</span> <span>共%d条</span>`,
			p.CurrentPage,
			p.TotalPage,
			p.TotalSize,
		)
		return pageStr

	case 4:
		p.NextPageTag = "下一页"
		p.PrevPageTag = "上一页"
		p.FirstPageTag = "首页"
		p.LastPageTag = "尾页"
		pageStr := p.FirstPage()
		pageStr += p.PrevPage()
		pageStr += p.PageBar("current")
		pageStr += p.NextPage()
		pageStr += p.LastPage()
		return pageStr
	}
	return ""
}

// 为指定的页面返回地址值
func (p *Page) GetUrl(pageNo int) string {
	pattern := fmt.Sprintf(`(:%s|\*%s|\.%s)`, p.PageName, p.PageName, p.PageName)
	result, _ := gregex.ReplaceString(pattern, pageNo, p.UrlTemplate)
		url.Path = gstr.Replace(p.UrlTemplate, "{.page}", gconv.String(pageNo))
		return url.String()
	}

	values := p.Url.Query()
	values.Set(p.PageName, gconv.String(pageNo))
	url.RawQuery = values.Encode()
	return url.String()
}

// 获取链接地址
func (p *Page) GetLink(url, text, title, style string) string {
	if len(style) > 0 {
		style = fmt.Sprintf(`class="%s" `, style)
	}
	if len(p.AjaxActionName) > 0 {
		return fmt.Sprintf(`<a %shref='#' onclick="%s('%s')">%s</a>`, style, p.AjaxActionName, url, text)
	} else {
		return fmt.Sprintf(`<a %shref="%s" title="%s">%s</a>`, style, url, title, text)
	}
}
