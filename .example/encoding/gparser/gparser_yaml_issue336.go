package main

import (
	"fmt"

	"github.com/jin502437344/gf/encoding/gparser"
)

type Conf struct {
	Broker    Broker
	Database  Database
	Scheduler Scheduler
	Worker    Worker
	Common    Common
	Report    Report
}

type Broker struct {
	Ip        string
	Tcport    string
	Httpport  string
	User      string
	Pwd       string
	Clusterid string
}

type Database struct {
	Ip     string
	Port   string
	User   string
	Pwd    string
	Dbname string
}

type Scheduler struct {
	SleepTime   int
	Pidfilepath string
	A2rparallel int
}

type Worker struct {
	Name        string
	Domains     map[string]interface{}
	Temppath    string
	Pidfilepath string
}

type Common struct {
	Filepath string
	Logpath  string
	Logdebug bool
	Logtrace bool
}

type Report struct {
	Localip   string
	Localport string //暂不启用
}

func main() {
	_, err := gparser.Load("config.yaml")
	if err != nil {
		fmt.Println("oops,read config.yaml err:", err)
	}

	//fmt.Println("yaml.v3读取yaml文件")
	//f, err := os.Open("config.yaml")
	//if err != nil {
	//	panic(err)
	//}
	//var conf Conf
	//err = yaml.NewDecoder(f).Decode(&conf)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(conf)
}
