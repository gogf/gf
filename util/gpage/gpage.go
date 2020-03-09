// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gpage provides useful paging functionality for web pages.
package gpage

import (
	"fmt"
	"github.com/gogf/gf/text/gstr"
	"github.com/gogf/gf/util/gconv"
	"math"
)

// Page is the pagination implementer.
// All the attributes are public, you can change them when necessary.
type Page struct {
	TotalSize      int    // Total size.
	TotalPage      int    // Total page, which is automatically calculated.
	CurrentPage    int    // Current page number >= 1.
	UrlTemplate    string // Custom url template for page url producing.
	LinkStyle      string // CSS style name for HTML link tag <a>.
	SpanStyle      string // CSS style name for HTML span tag <span>, which is used for first, current and last page tag.
	SelectStyle    string // CSS style name for HTML select tag <select>.
	NextPageTag    string // Tag name for next p.
	PrevPageTag    string // Tag name for prev p.
	FirstPageTag   string // Tag name for first p.
	LastPageTag    string // Tag name for last p.
	PrevBarTag     string // Tag string for prev bar.
	NextBarTag     string // Tag string for next bar.
	PageBarNum     int    // Page bar number for displaying.
	AjaxActionName string // Ajax function name. Ajax is enabled if this attribute is not empty.
}

const (
	PAGE_NAME         = "page"    // PAGE_NAME defines the default page name.
	PAGE_PLACE_HOLDER = "{.page}" // PAGE_PLACE_HOLDER defines the place holder for the url template.
)

// New creates and returns a pagination manager.
// Note that the parameter <urlTemplate> specifies the URL producing template, like:
// /user/list/{.page}, /user/list/{.page}.html, /user/list?page={.page}&type=1, etc.
// The build-in variable in <urlTemplate> "{.page}" specifies the page number, which will be replaced by certain
// page number when producing.
func New(totalSize, pageSize, currentPage int, urlTemplate string) *Page {
	p := &Page{
		LinkStyle:    "GPageLink",
		SpanStyle:    "GPageSpan",
		SelectStyle:  "GPageSelect",
		PrevPageTag:  "<",
		NextPageTag:  ">",
		FirstPageTag: "|<",
		LastPageTag:  ">|",
		PrevBarTag:   "<<",
		NextBarTag:   ">>",
		TotalSize:    totalSize,
		TotalPage:    int(math.Ceil(float64(totalSize) / float64(pageSize))),
		CurrentPage:  currentPage,
		PageBarNum:   10,
		UrlTemplate:  urlTemplate,
	}
	if currentPage == 0 {
		p.CurrentPage = 1
	}
	return p
}

// NextPage returns the HTML content for the next page.
func (p *Page) NextPage() string {
	if p.CurrentPage < p.TotalPage {
		return p.GetLink(p.CurrentPage+1, p.NextPageTag, "")
	}
	return fmt.Sprintf(`<span class="%s">%s</span>`, p.SpanStyle, p.NextPageTag)
}

// PrevPage returns the HTML content for the previous page.
func (p *Page) PrevPage() string {
	if p.CurrentPage > 1 {
		return p.GetLink(p.CurrentPage-1, p.PrevPageTag, "")
	}
	return fmt.Sprintf(`<span class="%s">%s</span>`, p.SpanStyle, p.PrevPageTag)
}

// FirstPage returns the HTML content for the first page.
func (p *Page) FirstPage() string {
	if p.CurrentPage == 1 {
		return fmt.Sprintf(`<span class="%s">%s</span>`, p.SpanStyle, p.FirstPageTag)
	}
	return p.GetLink(1, p.FirstPageTag, "")
}

// LastPage returns the HTML content for the last page.
func (p *Page) LastPage() string {
	if p.CurrentPage == p.TotalPage {
		return fmt.Sprintf(`<span class="%s">%s</span>`, p.SpanStyle, p.LastPageTag)
	}
	return p.GetLink(p.TotalPage, p.LastPageTag, "")
}

// PageBar returns the HTML page bar content with link and span tags.
func (p *Page) PageBar() string {
	plus := int(math.Ceil(float64(p.PageBarNum / 2)))
	if p.PageBarNum-plus+p.CurrentPage > p.TotalPage {
		plus = p.PageBarNum - p.TotalPage + p.CurrentPage
	}
	begin := p.CurrentPage - plus + 1
	if begin < 1 {
		begin = 1
	}
	barContent := ""
	for i := begin; i < begin+p.PageBarNum; i++ {
		if i <= p.TotalPage {
			if i != p.CurrentPage {
				barText := gconv.String(i)
				barContent += p.GetLink(i, barText, barText)
			} else {
				barContent += fmt.Sprintf(`<span class="%s">%d</span>`, p.SpanStyle, i)
			}
		} else {
			break
		}
	}
	return barContent
}

// SelectBar returns the select HTML content for pagination.
func (p *Page) SelectBar() string {
	barContent := fmt.Sprintf(`<select name="%s" onchange="window.location.href=this.value">`, p.SelectStyle)
	for i := 1; i <= p.TotalPage; i++ {
		if i == p.CurrentPage {
			barContent += fmt.Sprintf(`<option value="%s" selected>%d</option>`, p.GetUrl(i), i)
		} else {
			barContent += fmt.Sprintf(`<option value="%s">%d</option>`, p.GetUrl(i), i)
		}
	}
	barContent += "</select>"
	return barContent
}

// GetContent returns the page content for predefined mode.
// These predefined contents are mainly for chinese localization purpose. You can defines your own
// page function retrieving the page content according to the implementation of this function.
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
		pageStr += p.PageBar()
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
		pageStr += p.PageBar()
		pageStr += p.NextPage()
		pageStr += p.LastPage()
		return pageStr
	}
	return ""
}

// GetUrl parses the UrlTemplate with given page number and returns the URL string.
// Note that the UrlTemplate attribute can be either an URL or a URI string with "{.page}"
// place holder specifying the page number position.
func (p *Page) GetUrl(page int) string {
	return gstr.Replace(p.UrlTemplate, PAGE_PLACE_HOLDER, gconv.String(page))
}

// GetLink returns the HTML link tag <a> content for given page number.
func (p *Page) GetLink(page int, text, title string) string {
	if len(p.AjaxActionName) > 0 {
		return fmt.Sprintf(
			`<a class="%s" href="javascript:%s('%s')" title="%s">%s</a>`,
			p.LinkStyle, p.AjaxActionName, p.GetUrl(page), title, text,
		)
	} else {
		return fmt.Sprintf(
			`<a class="%s" href="%s" title="%s">%s</a>`,
			p.LinkStyle, p.GetUrl(page), title, text,
		)
	}
}
