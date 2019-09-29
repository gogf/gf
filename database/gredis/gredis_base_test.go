package gredis_test

import (
	"errors"
	"fmt"
	"github.com/gogf/gf/database/gredis"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/test/gtest"
	"github.com/gogf/gf/util/gconv"
	"os"
	"strings"
	"testing"
)

var (
	Clusterip     = "192.168.0.55" //
	Pass1         = "123456"       //123456
	port          = 8579           //8579 6379
	ClustersNodes = []string{Clusterip + ":7001", Clusterip + ":7002", Clusterip + ":7003", Clusterip + ":7004", Clusterip + ":7005", Clusterip + ":7006"}
	config        = gredis.Config{
		Host: Clusterip, //192.168.0.55 127.0.0.1
		Port: port,      //8579 6379
		Db:   1,
		Pass: "yyb513941", // when is ci,no pass
	}
)

func init() {
	gredis.FlagBanCluster = false
	// pwd  = "123456"    when is ci,no pass
	config := `[rediscluster]
    [rediscluster.default]
        host = "` + strings.Join(ClustersNodes, ",") + `"
		pwd  ="` + Pass1 + `"
        
[redis]
     default = "` + Clusterip + `:` + gconv.String(port) + `,1"` // 8579  6379
	err := createTestFile("config.toml", config)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func createTestFile(filename, content string) error {
	//TempDir := testpath()
	err := gfile.PutContents(filename, content)
	return err
}

// get testdir
func testpath() string {
	return os.TempDir()
}





func Test_RedisDo(t *testing.T) {
	gtest.Case(t, func() {
		var(
			err = errors.New("")
			s string
			ss []string
			n int
		)

		redis := gredis.New(config)
		defer redis.Close()

		redis.Set("k1","kv1")
		redis.Set("k2","kv2")
		ss,err=redis.Mget("k1","k2")
		gtest.Assert(err,nil)
		gtest.Assert(ss,[]string{"kv1","kv2"})

		_,err=redis.Mget()
		gtest.AssertNE(err,nil)

		s,err=redis.Mset("k1","kv11","k2","kv22")
		gtest.Assert(err,nil)
		gtest.Assert(s,"OK")
		s,err=redis.Mset()
		gtest.AssertNE(err,nil)


		s,err=redis.Rename("k1","k11")
		gtest.Assert(err,nil)
		s,err=redis.Get("k11")
		gtest.Assert(s,"kv11")

		_,err=redis.Renamenx("k11","k113")
		gtest.Assert(err,nil)
		s,err=redis.Get("k113")
		gtest.Assert(err,nil)
		gtest.Assert(s,"kv11")

		n,err=redis.Msetnx("k_1","kv_1","k_2","kv_2")
		gtest.Assert(err,nil)
		gtest.Assert(n,1)
		n,err=redis.Msetnx("k_1","kv_1","k_3","kv_3")
		gtest.Assert(err,nil)
		gtest.Assert(n,0)
		redis.Del("k_1")
		redis.Del("k_2")


		// Msetnx


	})
}


func Test_Clustersg(t *testing.T) {
	gtest.Case(t, func() {
		var(
			n int
			n64 int64
			rr interface{}
			rrs []interface{}
			err = errors.New("")
			s string
			//ss []string
		)

		gredis.FlagBanCluster = false

		rdb:=g.Redis()

		rr, err = rdb.Cluster("info")
		gtest.Assert(err, nil)
		str1 := gconv.String(rr)
		if !strings.Contains(str1, "cluster_state:ok") {
			t.Errorf("cluster errs.")
		}

		_, err = rdb.Set("jjname1", "jjqrr1")
		gtest.Assert(err, nil)
		_, err = rdb.Set("jjname2", "jjqrr2")
		_, err = rdb.Set("jjname3", "jjqrr3")
		gtest.Assert(err, nil)
		rr, err2 := rdb.Get("jjname2")
		gtest.Assert(err2, nil)
		gtest.Assert(gconv.String(rr), "jjqrr2")

		rdb.Set("n1","10")

		rr, err= rdb.Get("jjname1")
		gtest.Assert(err, nil)
		gtest.Assert(gconv.String(rr), "jjqrr1")

		rr3, err3 :=rdb.Get("jjname3")
		gtest.Assert(err3, nil)
		gtest.Assert(gconv.String(rr3), "jjqrr3")

		n,_=rdb.Exists("jjname3")
		gtest.Assert(n,1)

		n64,_=rdb.Expire("jjname3",300)
		gtest.Assert(n,1)
		n64 ,_=rdb.Ttl("jjname3")
		gtest.AssertGT(n64,200)

		rr,_=rdb.Dump("jjname3")
		gtest.AssertNE(rr,nil)

		n,_=g.Redis().Expireat("jjname3",gtime.Now().Second()+120)
		gtest.Assert(n,1)

		rrs,_=rdb.Keys("*jjname*")
		gtest.AssertGT(len(rrs),0)

		rr,_=rdb.Object("REFCOUNT","jjname3")
		gtest.AssertGT(gconv.Int(rr),0)

		n,_=rdb.Persist("jjname3")
		gtest.Assert(n,1)
		n,_=rdb.Persist("jjname3_")
		gtest.Assert(n,0)

		n64,_=rdb.Pttl("jjname3")
		gtest.Assert(n64,-1)
		n64,_=rdb.Pttl("jjname3_")
		gtest.AssertLT(n64,0)
		g.Redis().Expire("jjname3",10)
		n64,_=rdb.Pttl("jjname3")
		gtest.AssertGT(n64,5)


		rr,_=rdb.RandomKey()
		gtest.AssertNE(rr,nil)






		_,err=rdb.ReStore("jjname2",100000,"servals")
		gtest.AssertNE(err,nil)

		s,err=rdb.Dump("jjname1")
		gtest.Assert(err,nil)

		s,err=rdb.ReStore("jjname1",100000,s,"replace")
		gtest.Assert(err,nil)
		gtest.Assert(s,"OK")


		n64,err=rdb.Lpush("numlist2",1,3)
		gtest.Assert(err,nil)
		gtest.AssertGT(n64,0)
		n64,err=rdb.Lpush("numlist2",2)

		rrs,err= rdb.Sort("numlist2","desc")
		gtest.Assert(err,nil)
		gtest.Assert(gconv.SliceStr(rrs),[]string{"3","2","1"})


		//=============================del this lists after test
		n,err=rdb.Del("numlist2")
		gtest.Assert(n,1)

		s,err=rdb.Get("numlist2")
		gtest.Assert(err,nil)
		gtest.Assert(s,"")
		// Sort

		rdb.Set("jname2","a")
		n64,err=rdb.Append("jname2","q")
		s,err=rdb.Get("jname2")
		gtest.Assert(s,"aq")

		s,err= rdb.Type("jname2")
		gtest.Assert(err,nil)
		gtest.Assert(s,"string")

		n,err=rdb.Setbit("jname2",3,1)
		gtest.Assert(err,nil)
		gtest.Assert(n,0)

		n,err=rdb.Getbit("jname2",3)
		gtest.Assert(err,nil)
		gtest.Assert(n,1)

		n,err=rdb.BitCount("jname2")
		gtest.Assert(err,nil)
		gtest.Assert(n,8)

		rdb.Set("jname22","tt22")
		n,err=rdb.Setbit("jname22",3,1)
		n,err=rdb.BiTop("and","and-result","jname2","jname22")
		gtest.AssertNE(err,nil)

		n,err=rdb.BitPos("jname22",1)
		gtest.Assert(err,nil)
		gtest.Assert(n,1)


		rrs,err=rdb.BitField("jname2 set a")
		gtest.Assert(err,nil)


		n64,err=rdb.Decr("n1")
		gtest.Assert(err,nil)
		gtest.Assert(n64,9)


		n64,err=rdb.Decrby("n1",2)
		gtest.Assert(err,nil)
		gtest.Assert(n64,7)

		n64,err=rdb.Decrby("n_21",2)
		gtest.Assert(err,nil)
		gtest.Assert(n64,-2)
		n,err=rdb.Del("n_21")
		gtest.Assert(err,nil)

		s,err=rdb.GetRange("jjname1",1,2)
		gtest.Assert(err,nil)
		gtest.Assert(s,"jq")

		rr,err=rdb.GetSet("jjname1","imjjname1")
		gtest.Assert(err,nil)
		gtest.Assert(rr,"jjqrr1")
		s,err=rdb.Get("jjname1")
		gtest.Assert(s,"imjjname1")


		n64,err=rdb.Incr("tn1")
		gtest.Assert(err,nil)
		gtest.Assert(n64,1)
		n64,err=rdb.Incr("tn1")
		gtest.Assert(n64,2)



		n64,err=rdb.IncrBy("tn1",2)
		gtest.Assert(err,nil)
		gtest.Assert(n64,4)
		rdb.Del("tn1")

		s,err=rdb.IncrByFloat("tn2",3.4)
		gtest.Assert(err,nil)
		gtest.Assert(s,"3.4")
		rdb.Del("tn2")



		// IncrByFloat


	})
}
