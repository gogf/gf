package main

import (
    "fmt"
    "gitee.com/johng/gf/g/encoding/gxml"
    "gitee.com/johng/gf/g/os/gfile"
)

func main() {
    //json := gfile.GetBinContents("/home/john/Workspace/Go/GOPATH/src/gitee.com/johng/gf/geg/frame/config.json")
    //y, err := yaml.JSONToYAML(json)
    //fmt.Println(err)
    //fmt.Println(string(y))
    //
    //j, err := yaml.YAMLToJSON(y)
    //fmt.Println(err)
    //fmt.Println(string(j))

    x := `
<?xml version="1.0" encoding="UTF-8"?>
<rss kk-name="1">
<channel>
	<title>Comments for 碎言碎语</title>
	<atom:link href="http://johng.cn/comments/feed/" rel="self" type="application/rss+xml" />
	<link>http://johng.cn</link>
	<description></description>
	<lastBuildDate>Fri, 05 Jan 2018 02:56:11 +0000</lastBuildDate>
	<sy:updatePeriod>hourly</sy:updatePeriod>
	<sy:updateFrequency>1</sy:updateFrequency>
	<generator>https://wordpress.org/?v=4.7.3</generator>
	<item>
		<title>Comment on Go性能优化：string与[ ]byte转换 by John</title>
		<link>http://johng.cn/go-optimize-string-bytes/#comment-114</link>
		<dc:creator><![CDATA[John]]></dc:creator>
		<pubDate>Fri, 05 Jan 2018 02:56:11 +0000</pubDate>
		<guid isPermaLink="false">http://johng.cn/?p=3435#comment-114</guid>
		<description><![CDATA[这篇文章我转自雨痕，由于string和[]byte之间的转换比较常用，所以这篇文章的性能对比才比较惊人。其他基本类型与[]byte这件的转换性能差别不大，struct与[]byte的转换性能差别主要在数据结构设计上，比如socket通信的话，基本都是自己通过二进制转换来组织需要的数据结构。总得来说，在go开发的时候尽量多用二进制参数([]byte)是好的。]]></description>
		<content:encoded><![CDATA[<p>这篇文章我转自雨痕，由于string和[]byte之间的转换比较常用，所以这篇文章的性能对比才比较惊人。其他基本类型与[]byte这件的转换性能差别不大，struct与[]byte的转换性能差别主要在数据结构设计上，比如socket通信的话，基本都是自己通过二进制转换来组织需要的数据结构。总得来说，在go开发的时候尽量多用二进制参数([]byte)是好的。</p>
]]></content:encoded>
	</item>

	<item>
		<title>Comment on 层次不同的人，是很难沟通的 by buhuipao</title>
		<link>http://johng.cn/hard-for-communitication-between-different-levels/#comment-105</link>
		<dc:creator><![CDATA[buhuipao]]></dc:creator>
		<pubDate>Tue, 19 Sep 2017 09:10:08 +0000</pubDate>
		<guid isPermaLink="false">http://johng.cn/?p=2854#comment-105</guid>
		<description><![CDATA[博主说得很对，认知很重要。]]></description>
		<content:encoded><![CDATA[<p>博主说得很对，认知很重要。</p>
]]></content:encoded>
	</item>
</channel>
</rss>`

    j, _ := gxml.ToJson([]byte(x))
    fmt.Println(string(j))
}