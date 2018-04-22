// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// 分页管理.
package gpage

import (
    "gitee.com/johng/gf/g/net/ghttp"
    "math"
    "fmt"
    "net/url"
    "gitee.com/johng/gf/g/util/gconv"
    "strings"
)

type Page struct {
    pageName       string // 分页参数名称
    nextPageTag    string // 下一页标签
    prevPageTag     string // 上一页标签
    firstPageTag   string // 首页标签
    lastPageTag    string // 尾页标签

    prevBar        string // 上一分页条
    nextBar        string // 下一分页条
    formatLeft     string
    formatRight    string
    isAjax         bool   // 是否支持AJAX分页模式
    totalSize      int
    pagebarNum     int    // 控制记录条的个数。
    totalPage      int    // 总页数
    ajaxActionName string // AJAX动作名
    currentPage    int    // 当前页
    url            string // url地址头
    offset         int
}

func New(totalSize, perPage, currentPage int, url string) *Page {
    page := &Page {
        pageName    : "page",
        totalSize   : totalSize,
        totalPage   : int(math.Ceil(float64(totalSize/perPage))),
        currentPage : currentPage,
        offset      : (currentPage - 1)*perPage,
        url         : url,
    }
    if strings.Index(url, "?") != -1 {
        page.url = url + "&"
    } else {
        page.url = url + "?"
    }
    page.url += page.pageName + "="
    return page
}

// 启用AJAX分页
func (page *Page)EnableAjax(actionName string) {
    page.isAjax          = true
    page.ajaxActionName = actionName
}

// 获取显示"下一页"的内容.
func (page *Page) nextPage(curStyle , style string) string {
    if page.currentPage < page.totalPage {
        return page._getLink(page._getUrl(page.currentPage + 1), page.nextPageTag, "下一页", style)
    }
    return fmt.Sprintf(`<span class="%s">%s</span>`, curStyle, page.nextPageTag)
}

/// 获取显示“上一页”的内容
func (page *Page) prevPage(curStyle , style string) string {
    if page.currentPage > 1 {
        return page._getLink(page._getUrl(page.currentPage - 1), page.prevPageTag, "上一页", style)
    }
    return fmt.Sprintf(`<span class="%s">%s</span>`, curStyle, page.prevPageTag)
}

/**
* 获取显示“首页”的代码
*
* @return string
*/
func (page *Page)firstPage(curStyle, style string) string {
    if page.currentPage == 1 {
        return fmt.Sprintf(`<span class="%s">%s</span>`, curStyle, page.firstPageTag)
    }
    return page._getLink(page._getUrl(1), page.firstPageTag, "第一页", style)
}

// 获取显示“尾页”的内容
func (page *Page)lastPage(curStyle, style string) string {
    if page.currentPage == page.totalPage {
        return fmt.Sprintf(`<span class="%s">%s</span>`, curStyle, page.lastPageTag)
    }
    return page._getLink(page._getUrl(page.totalPage), page.lastPageTag, "最后页", style)
}

// 获得分页条。
func (page *Page) nowBar(curStyle, style string) string {
    plus := int(math.Ceil(float64(page.pagebarNum / 2)))
    if page.pagebarNum - plus + page.currentPage > page.totalPage {
        plus = page.pagebarNum - page.totalPage + page.currentPage
    }
    begin := page.currentPage - plus + 1
    if begin < 1 {
        begin = 1
    }
    ret := ""
    for i := begin; i < begin + page.pagebarNum; i++ {
        if i <= page.totalPage {
            if i != page.currentPage {
                ret += page._getText(page._getLink(page._getUrl(i), gconv.String(i), style, ""))
            } else {
                ret += page._getText(fmt.Sprintf(`<span class="%s">%d</span>`, curStyle, i))
            }
        } else {
            break
        }
        ret += "\n"
    }
    return ret
}
/**
* 获取显示跳转按钮的代码
*
* @return string
*/
func (page *Page) selectBar() string {
    url := page._getUrl("' + this.value")
    ret := fmt.Sprintf(`<select name="gpage_select" onchange="window.location.href='%s'">`, url)
    for i := 1; i <= page.totalPage; i++ {
        if (i == page.currentPage) {
            ret += fmt.Sprintf(`<option value="%d" selected>%d</option>`, i, i)
        } else {
            ret += fmt.Sprintf(`<option value="%d">%d</option>`, i, i)
        }
    }
    ret += "</select>"
    return ret
}

/**
* 控制分页显示风格（你可以继承后增加相应的风格）
*
* @param int mode 显示风格分类。
* @return string
*/
func (page *Page)show(mode int) string {
    //switch (mode) {
    //case '1':
    //page.nextPage = '下一页'
    //page.prevPage  = '上一页'
    //return page.prevPage()."<span class=\"current\">{page.currentPage}</span>".page.nextPage()
    //break
    //
    //case '2':
    //page.nextPage  = '下一页>>'
    //page.prevPage   = '<<上一页'
    //page.firstPage = '首页'
    //page.lastPage  = '尾页'
    //return page.firstPage().page.prevPage().'<span class="current">[第'.page.currentPage.'页]</span>'.page.nextPage().page.lastPage().'第'.page.select().'页'
    //break
    //
    //case '3':
    //page.nextPage  = '下一页'
    //page.prevPage   = '上一页'
    //page.firstPage = '首页'
    //page.lastPage  = '尾页'
    //pageStr  = page.firstPage()." ".page.prevPage()
    //pageStr .= ' '.page.nowbar('current')
    //pageStr .= ' '.page.nextPage()." ".page.lastPage()
    //pageStr .= "<span>当前页{page.currentPage}/{page.totalPage}</span> <span>共{page.totalSize}条</span>"
    //return pageStr
    //break
    //
    //case '4':
    //page.nextPage  = '下一页'
    //page.prevPage   = '上一页'
    //page.firstPage = '首页'
    //page.lastPage  = '尾页'
    //pageStr  = page.firstPage()." ".page.prevPage()
    //pageStr .= ' '.page.nowbar('current')
    //pageStr .= ' '.page.nextPage()." ".page.lastPage()
    //return pageStr
    //break
    //}
    return ""
}

// 为指定的页面返回地址值
func (page *Page) _getUrl(pageNoStr string) string {
    return page.url + pageNoStr
}

// 获取分页显示文字，比如说默认情况下_getText('<a href="">1</a>')将返回[<a href="">1</a>]
func (page *Page)_getText(str string) string {
    return page.formatLeft + str + page.formatRight
}

// 获取链接地址
func (page *Page)_getLink(url, text, title, style string) string {
    if len(style) > 0 {
        style = fmt.Sprintf(`class="%s"`, style)
    }
    if (page.isAjax) {
        return fmt.Sprintf(`<a %s href='#' onclick="%s('%s')">%s</a>`, style, page.ajaxActionName, url, text)
    } else {
        return fmt.Sprintf(`"<a %s href="%s" title="%s">%s</a>"`, style, url, title, text)
    }
}

