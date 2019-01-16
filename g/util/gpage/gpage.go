// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// Package gpage provides useful paging functionality for web pages.
// 
// 分页管理.
package gpage

import (
    "fmt"
    "math"
    url2 "net/url"
    "gitee.com/johng/gf/g/util/gconv"
    "gitee.com/johng/gf/g/net/ghttp"
    "gitee.com/johng/gf/g/util/gregex"
    "gitee.com/johng/gf/g/util/gstr"
    "strings"
)

// 分页对象
type Page struct {
    Url            *url2.URL      // 当前页面的URL对象
    Router         *ghttp.Router  // 当前页面的路由对象(与gf框架耦合，在静态分页下有效)
    UrlTemplate    string         // URL生成规则，内部可使用{.page}变量指定页码
    TotalSize      int            // 总共数据条数
    TotalPage      int            // 总页数
    CurrentPage    int            // 当前页码
    PageName       string         // 分页参数名称(GET参数)
    NextPageTag    string         // 下一页标签
    PrevPageTag    string         // 上一页标签
    FirstPageTag   string         // 首页标签
    LastPageTag    string         // 尾页标签
    PrevBar        string         // 上一分页条
    NextBar        string         // 下一分页条
    PageBarNum     int            // 控制分页条的数量
    AjaxActionName string         // AJAX方法名，当该属性有值时，表示使用AJAX分页
}

// 创建一个分页对象，输入参数分别为：
// 总数量、每页数量、当前页码、当前的URL(URI+QUERY)、(可选)路由规则(例如: /user/list/:page、/order/list/*page、/order/list/{page}.html)
func New(TotalSize, perPage int,  CurrentPage interface{}, url string, router...*ghttp.Router) *Page {
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
    if len(router) > 0 {
        page.Router = router[0]
    }
    return page
}

// 启用AJAX分页
func (page *Page) EnableAjax(actionName string) {
    page.AjaxActionName = actionName
}

// 设置URL生成规则模板，模板中可使用{.page}变量指定页码位置
func (page *Page) SetUrlTemplate(template string) {
    page.UrlTemplate = template
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
    // 复制一个URL对象
    url := *page.Url
    if len(page.UrlTemplate) == 0  && page.Router != nil {
        page.UrlTemplate = page.makeUrlTemplate(url.Path, page.Router)
    }
    if len(page.UrlTemplate) > 0 {
        // 指定URL生成模板
        url.Path = gstr.Replace(page.UrlTemplate, "{.page}", gconv.String(pageNo))
        return url.String()
    }

    values := page.Url.Query()
    values.Set(page.PageName, gconv.String(pageNo))
    url.RawQuery = values.Encode()
    return url.String()
}

// 根据当前URL以及注册路由信息计算出对应的URL模板
func (page *Page) makeUrlTemplate(url string, router *ghttp.Router) (tpl string) {
    if page.Router != nil && len(router.RegNames) > 0 {
        if match, err := gregex.MatchString(router.RegRule, url); err == nil && len(match) > 0 {
            if len(match) > len(router.RegNames) {
                tpl          = router.Uri
                hasPageName := false
                for i, name := range router.RegNames {
                    rule := fmt.Sprintf(`[:\*]%s|\{%s\}`, name, name)
                    if !hasPageName && strings.Compare(name, page.PageName) == 0 {
                        hasPageName = true
                        tpl, _ = gregex.ReplaceString(rule, `{.page}`, tpl)
                    } else {
                        tpl, _ = gregex.ReplaceString(rule, match[i + 1], tpl)
                    }
                }
                if !hasPageName {
                    tpl = ""
                }
            }
        }
    }
    return
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

