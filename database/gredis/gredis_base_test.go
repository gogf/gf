package gredis_test

import (
	"errors"
	"fmt"
	"github.com/gogf/gf/database/gredis"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/test/gtest"
	"github.com/gogf/gf/util/gconv"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var (
	Clusterip     = "192.168.0.55" //模拟的集群ip地址1  127.0.0.1   8220开始
	ClustersNodes = []string{Clusterip + ":7001", Clusterip + ":7002", Clusterip + ":7003", Clusterip + ":7004", Clusterip + ":7005", Clusterip + ":7006"}
)

func init() {
	gredis.FlagBanCluster = false
	config := `[rediscluster]
    [rediscluster.default]
        host = "` + strings.Join(ClustersNodes, ",") + `"
[redis]
     default = "` + Clusterip + `:8579,1"`
	err := createTestFile("config.toml", config)
	if err != nil {
		fmt.Println(err.Error())
	}
}

// 创建测试文件
func createTestFile(filename, content string) error {
	//TempDir := testpath()
	//fmt.Println(TempDir+"/"+filename)
	//err := ioutil.WriteFile(TempDir+filename, []byte(content), 0666)
	err := gfile.PutContents(filename, content)
	return err
}

// 测试完删除文件或目录
func delTestFiles(filenames string) {
	os.RemoveAll(testpath() + filenames)
	//os.RemoveAll(filenames)
}

// 统一格式化文件目录为"/"
func formatpaths(paths []string) []string {
	for k, v := range paths {
		paths[k] = filepath.ToSlash(v)
		paths[k] = strings.Replace(paths[k], "./", "/", 1)
	}
	return paths
}

// 指定返回要测试的目录
func testpath() string {
	return os.TempDir()
}

func Test_ClusterDo(t *testing.T) {
	gtest.Case(t, func() {
		redis := gredis.NewClusterClient(&gredis.ClusterOption{
			Nodes: ClustersNodes,
			Pwd:   "123456",
		})
		redis.Set("jname2", "jqrr2")
		r, err := redis.Get("jname2")
		gtest.Assert(err, nil)
		gtest.Assert(gconv.String(r), "jqrr2")
	})
}

func Test_Clustersg(t *testing.T) {
	gtest.Case(t, func() {
		var rr interface{}
		err := errors.New("")
		gredis.FlagBanCluster = false

		_, err = g.Redis().Set("jjname1", "jjqrr1")
		gtest.Assert(err, nil)
		_, err = g.Redis().Set("jjname2", "jjqrr2")
		_, err = g.Redis().Set("jjname3", "jjqrr3")
		gtest.Assert(err, nil)
		rr, err2 := g.Redis().Get("jjname2")
		gtest.Assert(err2, nil)
		gtest.Assert(gconv.String(rr), "jjqrr2")

		rr3, err3 := g.Redis().Get("jjname3")
		gtest.Assert(err3, nil)
		gtest.Assert(gconv.String(rr3), "jjqrr3")

		rr, err = g.Redis().Cluster("info")
		gtest.Assert(err, nil)
		str1 := gconv.String(rr)
		if !strings.Contains(str1, "cluster_state:ok") {
			t.Errorf("cluster errs.")
		}

	})
}
