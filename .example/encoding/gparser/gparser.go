package main

import (
	"fmt"

	"github.com/gogf/gf/encoding/gparser"
	"github.com/gogf/gf/os/glog"
)

func getWithPattern1() {
	data :=
		`{
            "users" : {
                    "count" : 100,
                    "list"  : [
                        {"name" : "Ming", "score" : 60},
                        {"name" : "John", "score" : 99.5}
                    ]
            }
        }`

	if p, e := gparser.LoadContent([]byte(data)); e != nil {
		glog.Error(e)
	} else {
		fmt.Println("John Score:", p.GetFloat32("users.list.1.score"))
	}
}

func getWithPattern2() {
	data :=
		`<?xml version="1.0" encoding="UTF-8"?>
         <note>
           <to>Tove</to>
           <from>Jani</from>
           <heading>Reminder</heading>
           <body>Don't forget me this weekend!</body>
         </note>`

	if p, e := gparser.LoadContent([]byte(data)); e != nil {
		glog.Error(e)
	} else {
		fmt.Println("Heading:", p.GetString("note.heading"))
	}
}

// 当键名存在"."号时，检索优先级：键名->层级，因此不会引起歧义
func multiDots1() {
	data :=
		`{
            "users" : {
                "count" : 100
            },
            "users.count" : 101
        }`
	if p, e := gparser.LoadContent([]byte(data)); e != nil {
		glog.Error(e)
	} else {
		fmt.Println("Users Count:", p.Get("users.count"))
	}
}

func multiDots2() {
	data :=
		`{
            "users" : {
                "count" : {
                    "type1" : 1,
                    "type2" : 2
                },
                "count.type1" : 100
            }
        }`
	if p, e := gparser.LoadContent([]byte(data)); e != nil {
		glog.Error(e)
	} else {
		fmt.Println("Users Count:", p.Get("users.count.type1"))
		fmt.Println("Users Count:", p.Get("users.count.type2"))
	}
}

// 设置数据
func set1() {
	data :=
		`<?xml version="1.0" encoding="UTF-8"?>
         <article>
           <count>10</count>
           <list><title>gf article1</title><content>gf content1</content></list>
           <list><title>gf article2</title><content>gf content2</content></list>
           <list><title>gf article3</title><content>gf content3</content></list>
         </article>`
	if p, e := gparser.LoadContent([]byte(data)); e != nil {
		glog.Error(e)
	} else {
		p.Set("article.list.0", nil)
		c, _ := p.ToJson()
		fmt.Println(string(c))
		// {"article":{"count":"10","list":[{"content":"gf content2","title":"gf article2"},{"content":"gf content3","title":"gf article3"}]}}
	}
}

func set2() {
	data :=
		`{
            "users" : {
                "count" : 100
            }
        }`
	if p, e := gparser.LoadContent([]byte(data)); e != nil {
		glog.Error(e)
	} else {
		p.Set("users.count", 1)
		p.Set("users.list", []string{"John", "小明"})
		c, _ := p.ToJson()
		fmt.Println(string(c))
	}
}

func makeXml1() {
	p := gparser.New(nil)
	p.Set("name", "john")
	p.Set("age", 18)
	p.Set("scores", map[string]int{
		"语文": 100,
		"数学": 100,
		"英语": 100,
	})
	c, _ := p.ToXmlIndent("simple-xml")
	fmt.Println(string(c))
}

func makeJson1() {
	type Order struct {
		Id    int     `json:"id"`
		Price float32 `json:"price"`
	}
	p := gparser.New(nil)
	p.Set("orders.list.0", Order{1, 100})
	p.Set("orders.list.1", Order{2, 666})
	p.Set("orders.list.2", Order{3, 999.99})
	fmt.Println("Order 1 Price:", p.Get("orders.list.1.price"))
	c, _ := p.ToJson()
	fmt.Println(string(c))
}

func makeJson2() {
	p := gparser.New(map[string]string{
		"k1": "v1",
		"k2": "v2",
	})
	p.Set("k1.1", []int{1, 2, 3})
	//p.Set("0.0.1", []int{1,2,3})
	c, _ := p.ToJson()
	fmt.Println(string(c))
}

func makeJson3() {
	p := gparser.New([]string{"a"})
	p.Set("0.0.0", []int{1, 2, 3})
	c, _ := p.ToJson()
	fmt.Println(string(c))
}

func toStruct1() {
	type Info struct {
		Name string
		Url  string
	}
	o := Info{}
	p := gparser.New(map[string]string{
		"Name": "gf",
		"Url":  "https://gitee.com/johng",
	})
	p.ToStruct(&o)
	fmt.Println("Name:", o.Name)
	fmt.Println("Url :", o.Url)
}

func convert() {
	p := gparser.New(map[string]string{
		"name": "gf",
		"site": "https://gitee.com/johng",
	})
	c, _ := p.ToJson()
	fmt.Println("JSON:")
	fmt.Println(string(c))
	fmt.Println("======================")

	fmt.Println("XML:")
	c, _ = p.ToXmlIndent()
	fmt.Println(string(c))
	fmt.Println("======================")

	fmt.Println("YAML:")
	c, _ = p.ToYaml()
	fmt.Println(string(c))
	fmt.Println("======================")

	fmt.Println("TOML:")
	c, _ = p.ToToml()
	fmt.Println(string(c))

}

func remove1() {
	p := gparser.New(map[string]string{
		"k1": "v1",
		"k2": "v2",
	})
	p.Set("0.0.0.0.0.0.0.0", []int{1, 2, 3})
	p.Remove("0.0")
	c, _ := p.ToJson()
	fmt.Println(string(c))
}

func main() {
	toStruct1()
}
