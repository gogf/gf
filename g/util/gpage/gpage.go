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
    Url            *url2.URL // 当前页面的URL对象
    Route          string    // 当前页面的路由规则(在静态分页下有效)
    TotalSize      int       // 总共数据条数
    TotalPage      int       // 总页数
    CurrentPage    int       // 当前页码
    PageName       string    // 分页参数名称(GET参数)
    NextPageTag    string    // 下一页标签
    PrevPageTag    string    // 上一页标签
    FirstPageTag   string    // 首页标签
    LastPageTag    string    // 尾页标签
    PrevBar        string    // 上一分页条
    NextBar        string    // 下一分页条
    PageBarNum     int       // 控制分页条的数量
    AjaxActionName string    // AJAX方法名，当该属性有值时，表示使用AJAX分页
}

// 创建一个分页对象，输入参数分别为：
// 总数量、每页数量、当前页码、当前的URL(可以只是URI+QUERY)、(可选)路由规则(例如: /user/list/:page、/order/list/*page)
func New(TotalSize, perPage int,  CurrentPage interface{}, url string, route...string) *Page {
    u, _ := url2.Parse(url)
    page := &Page {
        PageName     : "page",
        PrevPageTag  : "<",
        NextPageTag  : ">",
        FirstPageTag : "|<",
        LastPageTag  : ">|",
        PrevBar      : "<<",
        NextBar      : ">>",
        TotalSize    : TotalSize,
        TotalPage    : int(math.Ceil(float64(TotalSize)/float64(perPage))),
        CurrentPage  : 1,
        PageBarNum   : 10,
        Url          : u,
    }
    curPage := gconv.Int(CurrentPage)
    if curPage > 0 {
        page.CurrentPage = curPage
    }
    if len(route) > 0 {
        page.Route = route[0]
    }
    return page
}

// 启用AJAX分页
func (page *Page) EnableAjax(actionName string) {
    page.AjaxActionName = actionName
}

// 获取显示"下一页"的内容.
func (page *Page) NextPage(styles ... string) string {
    var curStyle, style string
    if len(styles) > 0 {
        curStyle = styles[0]
    }
    if len(styles) > 1 {
        style    = styles[0]
    }
    if page.CurrentPage < page.TotalPage {
        return page.GetLink(page.GetUrl(page.CurrentPage + 1), page.NextPageTag, "下一页", style)
    }
    return fmt.Sprintf(`<span class="%s">%s</span>`, curStyle, page.NextPageTag)
}

/// 获取显示“上一页”的内容
func (page *Page) PrevPage(styles ... string) string {
    var curStyle, style string
    if len(styles) > 0 {
        curStyle = styles[0]
    }
    if len(styles) > 1 {
        style    = styles[0]
    }
    if page.CurrentPage > 1 {
        return page.GetLink(page.GetUrl(page.CurrentPage - 1), page.PrevPageTag, "上一页", style)
    }
    return fmt.Sprintf(`<span class="%s">%s</span>`, curStyle, page.PrevPageTag)
}

/**
* 获取显示“首页”的代码
*
* @return string
*/
func (page *Page) FirstPage(styles ... string) string {
    var curStyle, style string
    if len(styles) > 0 {
        curStyle = styles[0]
    }
    if len(styles) > 1 {
        style    = styles[0]
    }
    if page.CurrentPage == 1 {
        return fmt.Sprintf(`<span class="%s">%s</span>`, curStyle, page.FirstPageTag)
    }
    return page.GetLink(page.GetUrl(1), page.FirstPageTag, "第一页", style)
}

// 获取显示“尾页”的内容
func (page *Page) LastPage(styles ... string) string {
    var curStyle, style string
    if len(styles) > 0 {
        curStyle = styles[0]
    }
    if len(styles) > 1 {
        style    = styles[0]
    }
    if page.CurrentPage == page.TotalPage {
        return fmt.Sprintf(`<span class="%s">%s</span>`, curStyle, page.LastPageTag)
    }
    return page.GetLink(page.GetUrl(page.TotalPage), page.LastPageTag, "最后页", style)
}

// 获得分页条列表内容
func (page *Page) PageBar(styles ... string) string {
    var curStyle, style string
    if len(styles) > 0 {
        curStyle = styles[0]
    }
    if len(styles) > 1 {
        style    = styles[0]
    }
    plus := int(math.Ceil(float64(page.PageBarNum / 2)))
    if page.PageBarNum - plus + page.CurrentPage > page.TotalPage {
        plus = page.PageBarNum - page.TotalPage + page.CurrentPage
    }
    begin := page.CurrentPage - plus + 1
    if begin < 1 {
        begin = 1
    }
    ret := ""
    for i := begin; i < begin + page.PageBarNum; i++ {
        if i <= page.TotalPage {
            if i != page.CurrentPage {
                ret += page.GetLink(page.GetUrl(i), gconv.String(i), style, "")
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
func (page *Page) SelectBar() string {
    ret := `<select name="gpage_select" onchange="window.location.href=this.value">`
    for i := 1; i <= page.TotalPage; i++ {
        if i == page.CurrentPage {
            ret += fmt.Sprintf(`<option value="%s" selected>%d</option>`, page.GetUrl(i), i)
        } else {
            ret += fmt.Sprintf(`<option value="%s">%d</option>`, page.GetUrl(i), i)
        }
    }
    ret += "</select>"
    return ret
}

// 预定义的分页显示风格内容
func (page *Page) GetContent(mode int) string {
    switch mode {
        case 1:
            page.NextPageTag = "下一页"
            page.PrevPageTag = "上一页"
            return fmt.Sprintf(
                `%s <span class="current">%d</span> %s`,
                page.PrevPage(),
                page.CurrentPage,
                page.NextPage(),
            )

        case 2:
            page.NextPageTag  = "下一页>>"
            page.PrevPageTag  = "<<上一页"
            page.FirstPageTag = "首页"
            page.LastPageTag  = "尾页"
            return fmt.Sprintf(
                `%s%s<span class="current">[第%d页]</span>%s%s第%s页`,
                page.FirstPage(),
                page.PrevPage(),
                page.CurrentPage,
                page.NextPage(),
                page.LastPage(),
                page.SelectBar(),
            )

        case 3:
            page.NextPageTag  = "下一页"
            page.PrevPageTag  = "上一页"
            page.FirstPageTag = "首页"
            page.LastPageTag  = "尾页"
            pageStr := page.FirstPage()
            pageStr += page.PrevPage()
            pageStr += page.PageBar("current")
            pageStr += page.NextPage()
            pageStr += page.LastPage()
            pageStr += fmt.Sprintf(
                `<span>当前页%d/%d</span> <span>共%d条</span>`,
                page.CurrentPage,
                page.TotalPage,
                page.TotalSize,
            )
            return pageStr

        case 4:
            page.NextPageTag  = "下一页"
            page.PrevPageTag  = "上一页"
            page.FirstPageTag = "首页"
            page.LastPageTag  = "尾页"
            pageStr := page.FirstPage()
            pageStr += page.PrevPage()
            pageStr += page.PageBar("current")
            pageStr += page.NextPage()
            pageStr += page.LastPage()
            return pageStr
    }
    return ""
}

// 为指定的页面返回地址值
func (page *Page) GetUrl(pageNo int) string {
    url := *page.Url
    if len(page.Route) > 0 {
        // 这里基于路由匹配的URL页码替换比较简单，但能满足绝大多数场景
        index := -1
        array := strings.Split(page.Route, "/")
        for k, v := range array {
            if strings.EqualFold(v, ":" + page.PageName) || strings.EqualFold(v, "*" + page.PageName) {
                index = k
                break
            }
        }
        // 替换url.Path中的分页码
        if index != -1 {
            pathArray := strings.Split(page.Url.Path, "/")
            for i := 0; i <= index - len(pathArray); i++ {
                pathArray = append(pathArray, "")
            }
            pathArray[index] = gconv.String(pageNo)
            url.Path         = strings.TrimRight(strings.Join(pathArray, "/"), "/")
            return url.String()
        }
    }
    values := page.Url.Query()
    values.Set(page.PageName, gconv.String(pageNo))
    url.RawQuery = values.Encode()
    return url.String()
}

// 获取链接地址
func (page *Page) GetLink(url, text, title, style string) string {
    if len(style) > 0 {
        style = fmt.Sprintf(`class="%s" `, style)
    }
    if len(page.AjaxActionName) > 0 {
        return fmt.Sprintf(`<a %shref='#' onclick="%s('%s')">%s</a>`, style, page.AjaxActionName, url, text)
    } else {
        return fmt.Sprintf(`<a %shref="%s" title="%s">%s</a>`, style, url, title, text)
    }
}

