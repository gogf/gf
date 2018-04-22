// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// 分页管理.
package gpage

import (
    "fmt"
    "math"
    "strings"
    url2 "net/url"
    "gitee.com/johng/gf/g/util/gconv"
)

// 分页对象
type Page struct {
    pageName       string    // 分页参数名称
    nextPageTag    string    // 下一页标签
    prevPageTag    string    // 上一页标签
    firstPageTag   string    // 首页标签
    lastPageTag    string    // 尾页标签
    prevBar        string    // 上一分页条
    nextBar        string    // 下一分页条
    totalSize      int       // 总共条数
    pageBarNum     int       // 控制记录条的个数
    totalPage      int       // 总页数
    currentPage    int       // 当前页
    offset         int       // 分页的offset条数
    url            *url2.URL // URL对象
    route          string    // 路由规则
    ajaxActionName string    // AJAX动作名，当该属性有值时，表示使用AJAX分页
}

// 创建一个分页对象，输入参数分别为：
// 总数量、每页数量、当前页码、当前的URL(可以只是URI+QUERY)、(可选)路由规则(例如: /user/list/:page、/order/list/*page)
func New(totalSize, perPage, currentPage int, url string, route...string) *Page {
    u, _ := url2.Parse(url)
    page := &Page {
        pageName     : "page",
        prevPageTag  : "<",
        nextPageTag  : ">",
        firstPageTag : "|<",
        lastPageTag  : ">|",
        prevBar      : "<<",
        nextBar      : ">>",
        totalSize    : totalSize,
        totalPage    : int(math.Ceil(float64(totalSize/perPage))),
        currentPage  : currentPage,
        offset       : (currentPage - 1)*perPage,
        pageBarNum   : 10,
        url          : u,
    }
    if len(route) > 0 {
        page.route = route[0]
    }
    return page
}

// 启用AJAX分页
func (page *Page)EnableAjax(actionName string) {
    page.ajaxActionName = actionName
}

// 获取显示"下一页"的内容.
func (page *Page) nextPage(styles ... string) string {
    var curStyle, style string
    if len(styles) > 0 {
        curStyle = styles[0]
    }
    if len(styles) > 1 {
        style    = styles[0]
    }
    if page.currentPage < page.totalPage {
        return page.getLink(page.getUrl(page.currentPage + 1), page.nextPageTag, "下一页", style)
    }
    return fmt.Sprintf(`<span class="%s">%s</span>`, curStyle, page.nextPageTag)
}

/// 获取显示“上一页”的内容
func (page *Page) prevPage(styles ... string) string {
    var curStyle, style string
    if len(styles) > 0 {
        curStyle = styles[0]
    }
    if len(styles) > 1 {
        style    = styles[0]
    }
    if page.currentPage > 1 {
        return page.getLink(page.getUrl(page.currentPage - 1), page.prevPageTag, "上一页", style)
    }
    return fmt.Sprintf(`<span class="%s">%s</span>`, curStyle, page.prevPageTag)
}

/**
* 获取显示“首页”的代码
*
* @return string
*/
func (page *Page)firstPage(styles ... string) string {
    var curStyle, style string
    if len(styles) > 0 {
        curStyle = styles[0]
    }
    if len(styles) > 1 {
        style    = styles[0]
    }
    if page.currentPage == 1 {
        return fmt.Sprintf(`<span class="%s">%s</span>`, curStyle, page.firstPageTag)
    }
    return page.getLink(page.getUrl(1), page.firstPageTag, "第一页", style)
}

// 获取显示“尾页”的内容
func (page *Page)lastPage(styles ... string) string {
    var curStyle, style string
    if len(styles) > 0 {
        curStyle = styles[0]
    }
    if len(styles) > 1 {
        style    = styles[0]
    }
    if page.currentPage == page.totalPage {
        return fmt.Sprintf(`<span class="%s">%s</span>`, curStyle, page.lastPageTag)
    }
    return page.getLink(page.getUrl(page.totalPage), page.lastPageTag, "最后页", style)
}

// 获得分页条。
func (page *Page) nowBar(styles ... string) string {
    var curStyle, style string
    if len(styles) > 0 {
        curStyle = styles[0]
    }
    if len(styles) > 1 {
        style    = styles[0]
    }
    plus := int(math.Ceil(float64(page.pageBarNum / 2)))
    if page.pageBarNum - plus + page.currentPage > page.totalPage {
        plus = page.pageBarNum - page.totalPage + page.currentPage
    }
    begin := page.currentPage - plus + 1
    if begin < 1 {
        begin = 1
    }
    ret := ""
    for i := begin; i < begin + page.pageBarNum; i++ {
        if i <= page.totalPage {
            if i != page.currentPage {
                ret += page.getLink(page.getUrl(i), gconv.String(i), style, "")
            } else {
                ret += fmt.Sprintf(`<span class="%s">%d</span>`, curStyle, i)
            }
        } else {
            break
        }
        if i != begin + page.pageBarNum - 1 {
            ret += "\n"
        }
    }
    return ret
}
/**
* 获取显示跳转按钮的代码
*
* @return string
*/
func (page *Page) selectBar() string {
    ret := fmt.Sprintf(`<select name="gpage_select" onchange="window.location.href='%sthis.value'">`, page.url)
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

// 预定义的分页显示风格内容
func (page *Page) GetContent(mode int) string {
    switch (mode) {
        case 1:
            page.nextPageTag = "下一页"
            page.prevPageTag = "上一页"
            return fmt.Sprintf(`%s <span class="current">%d</span> %s`, page.prevPage(), page.currentPage, page.nextPage())


        case 2:
            page.nextPageTag  = "下一页>>"
            page.prevPageTag  = "<<上一页"
            page.firstPageTag = "首页"
            page.lastPageTag  = "尾页"
            return fmt.Sprintf(`%s%s <span class="current">[第%d页]</span> %s%s 第%s页`,
                page.firstPage(), page.prevPage(), page.currentPage, page.nextPage(), page.lastPage(), page.selectBar())


        case 3:
            page.nextPageTag  = "下一页"
            page.prevPageTag  = "上一页"
            page.firstPageTag = "首页"
            page.lastPageTag  = "尾页"
            pageStr := page.firstPage() + "\n"
            pageStr += page.prevPage() + "\n"
            pageStr += page.nowBar("current") + "\n"
            pageStr += page.nextPage() + "\n"
            pageStr += page.lastPage() + "\n"
            pageStr += fmt.Sprintf(`<span>当前页%d/%d</span> <span>共%d条</span>`, page.currentPage, page.totalPage, page.totalSize)
            return pageStr


        case 4:
            page.nextPageTag  = "下一页"
            page.prevPageTag  = "上一页"
            page.firstPageTag = "首页"
            page.lastPageTag  = "尾页"
            pageStr := page.firstPage() + "\n"
            pageStr += page.prevPage() + "\n"
            pageStr += page.nowBar("current") + "\n"
            pageStr += page.nextPage() + "\n"
            pageStr += page.lastPage() + "\n"
            return pageStr
    }
    return ""
}

// 为指定的页面返回地址值
func (page *Page) getUrl(pageNo int) string {
    url := *page.url
    if len(page.url.RawQuery) > 0 && len(page.url.Query().Get(page.pageName)) > 0 {
        values := page.url.Query()
        values.Set(page.pageName, gconv.String(pageNo))
        url.RawQuery = values.Encode()
    } else {
        // 这里基于路由匹配的URL页码替换比较简单，但能满足绝大多数场景
        index := -1
        array := strings.Split(page.route, "/")
        for k, v := range array {
            if strings.EqualFold(v, ":" + page.pageName) || strings.EqualFold(v, "*" + page.pageName) {
                index = k
                break
            }
        }
        // 替换url.Path中的分页码
        if index != -1 {
            array       := strings.Split(page.url.Path, "/")
            array[index] = gconv.String(pageNo)
            url.Path     = strings.Join(array, "/")
        }
    }
    return url.String()
}

// 获取链接地址
func (page *Page) getLink(url, text, title, style string) string {
    if len(style) > 0 {
        style = fmt.Sprintf(`class="%s" `, style)
    }
    if len(page.ajaxActionName) > 0 {
        return fmt.Sprintf(`<a %shref='#' onclick="%s('%s')">%s</a>`, style, page.ajaxActionName, url, text)
    } else {
        return fmt.Sprintf(`<a %shref="%s" title="%s">%s</a>`, style, url, title, text)
    }
}

