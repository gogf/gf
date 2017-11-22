package main

import (
    "fmt"
    "g/net/ghttp"
)

var kvUrl string      = "http://192.168.2.102:4168/kv"
var nodeUrl string    = "http://192.168.2.102:4168/node"
var serviceUrl string = "http://192.168.2.102:4168/service"

// kv操作
func addKV() {
    c := ghttp.NewClient()
    r := c.Put(kvUrl, "{\"name1\":\"john1\", \"name2\":\"john2\"}")
    fmt.Println("addKV:", r.ReadAll())
}

func getAllKV() {
    c := ghttp.NewClient()
    r := c.Get(kvUrl)
    fmt.Println("getAllKV:", r.ReadAll())
}

func getOneKV() {
    c := ghttp.NewClient()
    r := c.Get(kvUrl + "?k=name1")
    fmt.Println("getOneKV:", r.ReadAll())
}

func editKV() {
    c := ghttp.NewClient()
    r := c.Post(kvUrl, "{\"name1\":\"john3\", \"name2\":\"john4\"}")
    fmt.Println("editKV:", r.ReadAll())
}

func removeKV() {
    c := ghttp.NewClient()
    r := c.Delete(kvUrl, "[\"name1\"]")
    fmt.Println("removeKV:", r.ReadAll())
}


// node操作
func addNode() {
    c := ghttp.NewClient()
    r := c.Put(nodeUrl, "[\"172.17.42.1\"]")
    fmt.Println("addNode:", r.ReadAll())
}

func getAllNode() {
    c := ghttp.NewClient()
    r := c.Get(nodeUrl)
    fmt.Println("getAllNode:", r.ReadAll())
}

func removeNode() {
    c := ghttp.NewClient()
    r := c.Delete(nodeUrl, "[\"172.17.42.1\"]")
    fmt.Println("removeNode:", r.ReadAll())
}


// service操作
func getAllService() {
    c := ghttp.NewClient()
    r := c.Get(serviceUrl)
    fmt.Println("getAllService:", r.ReadAll())
}

func getOneService() {
    c := ghttp.NewClient()
    r := c.Get(serviceUrl + "?name=Site Database")
    fmt.Println("getOneService:", r.ReadAll())
}

func addDatabaseService() {
    c := ghttp.NewClient()
    s := `
{
    "name" : "Site Database",
    "type" : "mysql",
    "list" : [
        {"host":"192.168.2.102", "port":"3306", "user":"root", "pass":"123456", "database":"test"},
        {"host":"192.168.2.124", "port":"3306", "user":"root", "pass":"123456", "database":"tongwujie"}
    ]
}
    `
    r := c.Put(serviceUrl, s)
    fmt.Println("addDatabaseService:", r.ReadAll())
}

func editDatabaseService() {
    c := ghttp.NewClient()
    s := `
{
    "name" : "Site Database2",
    "type" : "mysql",
    "list" : [
        {"host":"192.168.2.102", "port":"3306", "user":"root", "pass":"123456", "database":"test"},
        {"host":"192.168.2.124", "port":"3306", "user":"root", "pass":"123456", "database":"tongwujie"}
    ]
}
    `
    r := c.Post(serviceUrl, s)
    fmt.Println("editDatabaseService:", r.ReadAll())
}

func removeDatabaseService() {
    c := ghttp.NewClient()
    r := c.Delete(serviceUrl, "[\"Site Database2\"]")
    fmt.Println("removeDatabaseService:", r.ReadAll())
}


func addWebService() {
    c := ghttp.NewClient()
    s := `
{
    "name" : "Site",
    "type" : "web",
    "list" : [
        {"url":"http://baidu.com", "check":"http://itsadeadlink.com"},
        {"url":"http://baidu.com"}
    ]
}
    `
    r := c.Put(serviceUrl, s)
    fmt.Println("addWebService:", r.ReadAll())
}

func editWebService() {
    c := ghttp.NewClient()
    s := `
{
    "name" : "Site2",
    "type" : "web",
    "list" : [
        {"url":"http://baidu.com"},
        {"url":"http://baidu.com"}
    ]
}
    `
    r := c.Post(serviceUrl, s)
    fmt.Println("editWebService:", r.ReadAll())
}

func removeWebService() {
    c := ghttp.NewClient()
    r := c.Delete(serviceUrl, "[\"Site2\"]")
    fmt.Println("removeWebService:", r.ReadAll())
}

func main() {
    addWebService()
}
