package main

import (
    "fmt"
    "gitee.com/johng/gf/g"
    "gitee.com/johng/gf/g/util/gregex"
)

func main() {
    //s := `1544180795 -- s_has_sess -- 41570504 -decryptSess- 41570504__iuVycRYg9qE3y7CsSgGZH1K2nxTdjPZN4fXot65zHIEmULO0Ow6LweJp5raWl8Ft -postSess- eyJpdiI6IkFwSWZ3eXFMcGxBZE5JcWF4aXh0M3c9PSIsInZhbHVlIjoiV3ZLeGduMnRoRkFZdmxHTzM5ZzdyU1JHWDMycmZlRERvNnFkaUR0SitlRjBrZnlYR1JvS2puTGZNUThSeFR0bWtlT3pza0l0elFqRk5mdXF6XC9FWWpWZnljVjdJbHd3dTRybEhldHZHTk5DQ015dlpYNHljNmxKMWJTRUVpY0E4IiwibWFjIjoiOTkxMzIxOTRhMGUxZWZiODM4NWZjNDZjYmVhNWY2NjhlZDZkNmVlNjY1MTE2N2VhZDAzYzY4NDJmZGFkMjY5YyJ9 -- 0 -- B8105CF2-1588-4753-9F86-9B8C36EB1842 -- iPhone 7 -- 12.1 -- 6.8.7 -- i -- 10.111.153.5 -- medlinker -- service -- unknown
//`
//   s := `[08-Dec-2018 13:35:03 Asia/Shanghai] Medlinker\Services\Message\MessageService|updateUserInfo|用户头像 URI 不能为空 in /var/www/med-d2d/app/Services/Message/RongCloudService.php on line 851`
    s := `[2018-12-01 13:35:03 Asia/Shanghai] 1544180795 Medlinker\Services\Message\MessageService|updateUserInfo|用户头像 URI 不能为空 in /var/www/med-d2d/app/Services/Message/RongCloudService.php on line 851`
    //m, e := gregex.MatchString(`/var/log/medlinker/[\w\-\_]+/(.+?)/{0,1}[\d\-\_]*\.log`, `/var/log/medlinker/med-questionnaire/nginx/error/access-20181206.log`)
    //m, e := gregex.MatchString(`/var/log/medlinker/[\w\-\_]+/(.+?)/{0,1}[\d\-\_]*\.log`, `/var/log/medlinker/med-questionnaire/storagelogs/events/sqlLog/2018-12-06.log`)
    m, e := gregex.MatchString(`(.*?((\d{4}[-/\.]\d{2}[-/\.]\d{2}|\d{1,2}[-/\.][A-Za-z]{3,}[-/\.]\d{4})[:\sT-]*\d{0,2}:{0,1}\d{0,2}:{0,1}\d{0,2}\.{0,1}\d{0,9}[\sZ]{0,1}[\+-]{0,1}[:\d]*|\d{10}).+)`, s)
    fmt.Println(e)
    g.Dump(m)
}
