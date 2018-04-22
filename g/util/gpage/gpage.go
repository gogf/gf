// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// 分页管理.
package gpage

import (
    "gitee.com/johng/gf/g/net/ghttp"
)

type Page struct {
    request        *ghttp.Request
    pageName       string // page标签，用来控制url页。比如说xxx.php?PBPage=2中的PBPage
    nextPage       string // 下一页标签
    prePage        string // 上一页标签
    firstPage      string // 首页标签
    lastPage       string // 尾页标签
    pre_bar        string // 上一分页条
    next_bar       string // 下一分页条
    formatLeft     string
    formatRight    string
    isAjax         bool // 是否支持AJAX分页模式
    totalSize      int
    pagebarNum     int    // 控制记录条的个数。
    totalPage      int    // 总页数
    ajaxActionName string // AJAX动作名
    currentPage    int    // 当前页
    url            string // url地址头
    offset         int
}
/**
* constructor构造函数
*
* @param array array['total'], array['perpage'], array['currentPage'], array['url'], array['ajax']...
*/
func New(total, perpage int) {
    total       = intval(array['total']);
    perpage     = (array_key_exists('perpage',array))     ? intval(array['perpage'])     : 10;
    currentPage = (array_key_exists('currentPage',array)) ? intval(array['currentPage']) : '';
    url         = (array_key_exists('url',array))         ? array['url']                 : '';
    } else {
    total       = array;
    perpage     = 10;
    currentPage ='';
    url         = '';
}
if ((!is_int(total)) || (total < 0)) {
this->_error(__FUNCTION__, 'invalid total');
}
if ((!is_int(perpage)) || (perpage <= 0)) {
this->_error(__FUNCTION__, 'invalid perpage');
}
if (!empty(array['pageName'])) {
// 设置pagename
this->set('pageName', array['pageName']);
}
this->_setCurrentPage(currentPage); // 设置当前页
this->_setUrl(url);                 // 设置链接地址
this->totalSize  = total;
this->totalPage  = ceil(total/perpage);
this->offset     = (this->currentPage-1)*perpage;
if (!empty(array['ajax'])) {
this->openAjax(array['ajax']);//打开AJAX模式
}

}

/**
* 设定类中指定变量名的值，如果改变量不属于这个类，将throw一个exception
*
* @param string var
* @param string value
*/
public function set(var, value)
{
if (inArray(var, get_object_vars(this))) {
this->var = value;
} else {
this->_error(__FUNCTION__, var." does not belong to PB_Page!");
}

}

/**
 * 使用AJAX模式。
 *
 * @param string action 默认ajax触发的动作名称。
 *
 * @return void
 */
public function openAjax(action)
{
this->isAjax          = true;
this->ajaxActionName = action;
}

/**
* 获取显示"下一页"的代码.
*
* @param string style
* @return string
*/
public function nextPage(curStyle = '', style = '')
{
if (this->currentPage < this->totalPage) {
return this->_getLink(this->_getUrl(this->currentPage+1), this->nextPage, '下一页', style);
}
return '<span class="'.curStyle.'">'.this->nextPage.'</span>';
}

/**
* 获取显示“上一页”的代码
*
* @param string style
* @return string
*/
public function prePage(curStyle='', style='')
{
if (this->currentPage > 1) {
return this->_getLink(this->_getUrl(this->currentPage - 1), this->prePage, '上一页', style);
}
return '<span class="'.curStyle.'">'.this->prePage.'</span>';
}

/**
* 获取显示“首页”的代码
*
* @return string
*/
public function firstPage(curStyle = '', style = '')
{
if (this->currentPage == 1) {
return '<span class="'.curStyle.'">'.this->firstPage.'</span>';
}
return this->_getLink(this->_getUrl(1), this->firstPage, '第一页', style);
}

/**
* 获取显示“尾页”的代码
*
* @return string
*/
public function lastPage(curStyle='', style='')
{
if (this->currentPage == this->totalPage) {
return '<span class="'.curStyle.'">'.this->lastPage.'</span>';
}
return this->_getLink(this->_getUrl(this->totalPage), this->lastPage, '最后页', style);
}

/**
 * 获得分页条。
 *
 * @param 当前页码 curStyle
 * @param 连接CSS style
 * @return 分页条字符串
 */
public function nowbar(curStyle = '', style = '')
{
plus = ceil(this->pagebarNum / 2);
if (this->pagebarNum - plus + this->currentPage > this->totalPage) {
plus = (this->pagebarNum - this->totalPage + this->currentPage);
}
begin  = this->currentPage - plus + 1;
begin  = (begin>=1) ? begin : 1;
return = '';
for (i = begin; i < begin + this->pagebarNum; i++) {
if (i <= this->totalPage) {
if (i != this->currentPage) {
return .= this->_getText(this->_getLink(this->_getUrl(i), i, style));
} else {
return .= this->_getText('<span class="'.curStyle.'">'.i.'</span>');
}
} else {
break;
}
return .= "\n";
}
unset(begin);
return return;
}
/**
* 获取显示跳转按钮的代码
*
* @return string
*/
public function select()
{
url    = this->_getUrl("' + this.value");
return = "<select name=\"PB_Page_Select\" onchange=\"window.location.href='url\">";
for (i=1; i <= this->totalPage; i++) {
if (i==this->currentPage) {
return .= '<option value="'.i.'" selected>'.i.'</option>';
} else {
return .= '<option value="'.i.'">'.i.'</option>';
}
}
unset(i);
return .= '</select>';
return return;
}

/**
* 获取mysql 语句中limit需要的值
*
* @return string
*/
public function offset()
{
return this->offset;
}

/**
* 控制分页显示风格（你可以继承后增加相应的风格）
*
* @param int mode 显示风格分类。
* @return string
*/
public function show(mode = 1)
{
switch (mode) {
case '1':
this->nextPage = '下一页';
this->prePage  = '上一页';
return this->prePage()."<span class=\"current\">{this->currentPage}</span>".this->nextPage();
break;

case '2':
this->nextPage  = '下一页>>';
this->prePage   = '<<上一页';
this->firstPage = '首页';
this->lastPage  = '尾页';
return this->firstPage().this->prePage().'<span class="current">[第'.this->currentPage.'页]</span>'.this->nextPage().this->lastPage().'第'.this->select().'页';
break;

case '3':
this->nextPage  = '下一页';
this->prePage   = '上一页';
this->firstPage = '首页';
this->lastPage  = '尾页';
pageStr  = this->firstPage()." ".this->prePage();
pageStr .= ' '.this->nowbar('current');
pageStr .= ' '.this->nextPage()." ".this->lastPage();
pageStr .= "<span>当前页{this->currentPage}/{this->totalPage}</span> <span>共{this->totalSize}条</span>";
return pageStr;
break;

case '4':
this->nextPage  = '下一页';
this->prePage   = '上一页';
this->firstPage = '首页';
this->lastPage  = '尾页';
pageStr  = this->firstPage()." ".this->prePage();
pageStr .= ' '.this->nowbar('current');
pageStr .= ' '.this->nextPage()." ".this->lastPage();
return pageStr;
break;
}

}

/*----------------private function (私有方法)-----------------------------------------------------------*/
/**
* 设置url头地址
* @param: string url
* @return boolean
*/
private function _setUrl(url = "")
{
if (!empty(url)) {
//手动设置
this->url = url.((stristr(url,'?')) ? '&' : '?').this->pageName."=";
} else {
parse = parse_url(_SERVER['REQUEST_URI']);
query = array();
if (!empty(parse['query'])) {
parse_str(parse['query'], query);
if (!empty(query) && isset(query[this->pageName])) {
unset(query[this->pageName]);
}
}
array = explode('?', _SERVER['REQUEST_URI']);
if (!empty(query)) {
this->url = array[0].'?'.http_build_query(query)."&{this->pageName}=";
} else {
this->url = array[0]."?{this->pageName}=";
}
}
}

/**
* 设置当前页面
*/
private function _setCurrentPage(currentPage)
{
if(empty(currentPage)) {
// 系统获取
if(isset(_GET[this->pageName])) {
this->currentPage = intval(_GET[this->pageName]);
}
} else {
//手动设置
this->currentPage = intval(currentPage);
}
}

/**
* 为指定的页面返回地址值
*
* @param int pageNo
* @return string url
*/
private function _getUrl(pageNo=1)
{
return this->url.pageNo;
}

/**
* 获取分页显示文字，比如说默认情况下_getText('<a href="">1</a>')将返回[<a href="">1</a>]
*
* @param String str
* @return string url
*/
private function _getText(str)
{
return this->formatLeft.str.this->formatRight;
}

//获取链接地址
private function _getLink(url, text, title='', style='')
{
style = (empty(style)) ? '' : 'class="'.style.'"';
if (this->isAjax) {
//如果是使用AJAX模式
return "<a style href='#' onclick=\"{this->ajaxActionName}('url');\">text</a>";
} else {
return "<a style href='url' title='title'>text</a>";
}
}

//出错处理方式
/**
 * 展示错误病终止执行.
 *
 * @param string function 错误产生的函数名称.
 * @param string errormsg 错误信息.
 *
 * @return void
 */
private function _error(function, errormsg)
{
die('Error in file <b>'.__FILE__.'</b> ,Function <b>'.function.'()</b> :'.errormsg);
}
}








