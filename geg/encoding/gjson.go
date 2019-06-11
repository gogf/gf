package main

import (
<<<<<<< HEAD
    "fmt"
    "gitee.com/johng/gf/g/os/glog"
    "gitee.com/johng/gf/g/encoding/gjson"
    "gitee.com/johng/gf/g/os/gtime"
)

func getByPattern() {
    data :=
        `{
=======
	"fmt"
	"github.com/gogf/gf/g/encoding/gjson"
	"github.com/gogf/gf/g/os/glog"
	"github.com/gogf/gf/g/os/gtime"
)

func getByPattern() {
	data :=
		`{
>>>>>>> upstream/master
            "users" : {
                    "count" : 100,
                    "list"  : [
                        {"name" : "小明",  "score" : 60},
                        {"name" : "John", "score" : 99.5}
                    ]
            }
        }`
<<<<<<< HEAD
    j, err := gjson.DecodeToJson([]byte(data))
    if err != nil {
        glog.Error(err)
    } else {
        fmt.Println("John Score:", j.GetFloat32("users.list.1.score"))
    }
=======
	j, err := gjson.DecodeToJson([]byte(data))
	if err != nil {
		glog.Error(err)
	} else {
		fmt.Println("John Score:", j.GetFloat32("users.list.1.score"))
	}
>>>>>>> upstream/master
}

// 当键名存在"."号时，检索优先级：键名->层级，因此不会引起歧义
func testMultiDots() {
<<<<<<< HEAD
    data :=
        `{
=======
	data :=
		`{
>>>>>>> upstream/master
            "users" : {
                "count" : 100
            },
            "users.count" : 101
        }`
<<<<<<< HEAD
    j, err := gjson.DecodeToJson([]byte(data))
    if err != nil {
        glog.Error(err)
    } else {
        fmt.Println("Users Count:", j.GetInt("users.count"))
    }
=======
	j, err := gjson.DecodeToJson([]byte(data))
	if err != nil {
		glog.Error(err)
	} else {
		fmt.Println("Users Count:", j.GetInt("users.count"))
	}
>>>>>>> upstream/master
}

// 设置数据
func testSet() {
<<<<<<< HEAD
    data :=
        `{
=======
	data :=
		`{
>>>>>>> upstream/master
            "users" : {
                "count" : 100
            }
        }`
<<<<<<< HEAD
    j, err := gjson.DecodeToJson([]byte(data))
    if err != nil {
        glog.Error(err)
    } else {
        j.Set("users.count",  1)
        j.Set("users.list",  []string{"John", "小明"})
        c, _ := j.ToJson()
        fmt.Println(string(c))
    }
=======
	j, err := gjson.DecodeToJson([]byte(data))
	if err != nil {
		glog.Error(err)
	} else {
		j.Set("users.count", 1)
		j.Set("users.list", []string{"John", "小明"})
		c, _ := j.ToJson()
		fmt.Println(string(c))
	}
>>>>>>> upstream/master
}

// 将Json数据转换为其他数据格式
func testConvert() {
<<<<<<< HEAD
    data :=
        `{
=======
	data :=
		`{
>>>>>>> upstream/master
            "users" : {
                "count" : 100,
                "list"  : ["John", "小明"]
            }
        }`
<<<<<<< HEAD
    j, err := gjson.DecodeToJson([]byte(data))
    if err != nil {
        glog.Error(err)
    } else {
        c, _ := j.ToJson()
        fmt.Println("JSON:")
        fmt.Println(string(c))
        fmt.Println("======================")

        fmt.Println("XML:")
        c, _ = j.ToXmlIndent()
        fmt.Println(string(c))
        fmt.Println("======================")

        fmt.Println("YAML:")
        c, _ = j.ToYaml()
        fmt.Println(string(c))
        fmt.Println("======================")

        fmt.Println("TOML:")
        c, _ = j.ToToml()
        fmt.Println(string(c))
    }
}


func testSplitChar() {
    var v interface{}
    j := gjson.New(nil)
    t1 := gtime.Nanosecond()
    j.Set("a.b.c.d.e.f.g.h.i.j.k", 1)
    t2 := gtime.Nanosecond()
    fmt.Println(t2 - t1)

    t5 := gtime.Nanosecond()
    v = j.Get("a.b.c.d.e.f.g.h.i.j.k")
    t6 := gtime.Nanosecond()
    fmt.Println(v)
    fmt.Println(t6 - t5)

    j.SetSplitChar('#')



    t7 := gtime.Nanosecond()
    v = j.Get("a#b#c#d#e#f#g#h#i#j#k")
    t8 := gtime.Nanosecond()
    fmt.Println(v)
    fmt.Println(t8 - t7)
}


func testViolenceCheck() {
    j := gjson.New(nil)
    t1 := gtime.Nanosecond()
    j.Set("a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.a", 1)
    t2 := gtime.Nanosecond()
    fmt.Println(t2 - t1)

    t3 := gtime.Nanosecond()
    j.Set("a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.a", 1)
    t4 := gtime.Nanosecond()
    fmt.Println(t4 - t3)

    t5 := gtime.Nanosecond()
    j.Get("a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.a")
    t6 := gtime.Nanosecond()
    fmt.Println(t6 - t5)


    j.SetViolenceCheck(false)


    t7 := gtime.Nanosecond()
    j.Set("a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.a", 1)
    t8 := gtime.Nanosecond()
    fmt.Println(t8 - t7)

    t9 := gtime.Nanosecond()
    j.Get("a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.a")
    t10 := gtime.Nanosecond()
    fmt.Println(t10 - t9)
}

func main() {
    testViolenceCheck()
}
=======
	j, err := gjson.DecodeToJson([]byte(data))
	if err != nil {
		glog.Error(err)
	} else {
		c, _ := j.ToJson()
		fmt.Println("JSON:")
		fmt.Println(string(c))
		fmt.Println("======================")

		fmt.Println("XML:")
		c, _ = j.ToXmlIndent()
		fmt.Println(string(c))
		fmt.Println("======================")

		fmt.Println("YAML:")
		c, _ = j.ToYaml()
		fmt.Println(string(c))
		fmt.Println("======================")

		fmt.Println("TOML:")
		c, _ = j.ToToml()
		fmt.Println(string(c))
	}
}

func testSplitChar() {
	var v interface{}
	j := gjson.New(nil)
	t1 := gtime.Nanosecond()
	j.Set("a.b.c.d.e.f.g.h.i.j.k", 1)
	t2 := gtime.Nanosecond()
	fmt.Println(t2 - t1)

	t5 := gtime.Nanosecond()
	v = j.Get("a.b.c.d.e.f.g.h.i.j.k")
	t6 := gtime.Nanosecond()
	fmt.Println(v)
	fmt.Println(t6 - t5)

	j.SetSplitChar('#')

	t7 := gtime.Nanosecond()
	v = j.Get("a#b#c#d#e#f#g#h#i#j#k")
	t8 := gtime.Nanosecond()
	fmt.Println(v)
	fmt.Println(t8 - t7)
}

func testViolenceCheck() {
	j := gjson.New(nil)
	t1 := gtime.Nanosecond()
	j.Set("a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.a", 1)
	t2 := gtime.Nanosecond()
	fmt.Println(t2 - t1)

	t3 := gtime.Nanosecond()
	j.Set("a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.a", 1)
	t4 := gtime.Nanosecond()
	fmt.Println(t4 - t3)

	t5 := gtime.Nanosecond()
	j.Get("a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.a")
	t6 := gtime.Nanosecond()
	fmt.Println(t6 - t5)

	j.SetViolenceCheck(false)

	t7 := gtime.Nanosecond()
	j.Set("a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.a", 1)
	t8 := gtime.Nanosecond()
	fmt.Println(t8 - t7)

	t9 := gtime.Nanosecond()
	j.Get("a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.a")
	t10 := gtime.Nanosecond()
	fmt.Println(t10 - t9)
}

func main() {
	testViolenceCheck()
}
>>>>>>> upstream/master
