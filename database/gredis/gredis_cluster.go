// by jroam

package gredis

import (
	"errors"
	"github.com/gogf/gf/text/gregex"
	"github.com/gogf/gf/util/gconv"
	"strings"
)

type ClusterOption struct {
	Nodes []string // cluster nodes, for example: []string{":7001",":7002"}
	Pwd   string   // cluster password for AUTH
}

var (
	flagIsCluster  = false
	err            = errors.New("")
	reply          = new(interface{})
	clusterPasswd  = "" //cluster of passwd
	slotsMap       = map[string][]int{}
	clusterConnMap = map[string]*Redis{}
	FlagBanCluster = false // Disable cluster mode
)

// Get the coverage of slots
func (r *Redis) layoutSlots() {
	*reply, err = r.Do("cluster", "nodes")
	if err != nil {
		return
	}

	slotss, err2 := gregex.MatchAllString(`([\d]+\.[\d]+\.[\d]+\.[\d]+\:[\d]+).+?master.+?connected ([\d]+)-([\d]+)`, gconv.String(*reply))
	if err2 != nil {
		return
	}

	for _, v := range slotss {
		if len(v) != 4 {
			return
		}
		max := gconv.Int(v[3])
		if max == 0 {
			return
		}
		slotsMap[v[1]] = []int{gconv.Int(v[2]), max}
	}
}

func NewClusterClient(co *ClusterOption) *Redis {
	clusterPasswd = co.Pwd
	clusres := newClusterConn(co)
	clusres.layoutSlots()
	return clusres
}

func newClusterClientByHost(host string) *Redis {
	hosts := strings.Split(host, ":")
	if len(hosts) != 2 {
		return nil
	}
	config := Config{
		Host: hosts[0],
		Port: gconv.Int(hosts[1]),
		Db:   0, //cluster only  use of  0
		Pass: clusterPasswd,
	}

	clusterConnMap[host] = New(config)
	flagIsCluster = true
	return clusterConnMap[host]
}

func newClusterConn(co *ClusterOption) *Redis {

	for _, v := range co.Nodes {
		host := strings.Split(v, ":")
		if len(host) != 2 {
			continue
		}
		return newClusterClientByHost(v)
	}
	return nil
}

func choiceConn(key string) *Redis {
	slots := gconv.Int(getCRC16([]byte(key)))
	ks := ""
	for k, v := range slotsMap {
		if slots >= v[0] && slots <= v[1] {
			ks = k
			break
		}
	}

	if _, ok := clusterConnMap[ks]; !ok {
		newClusterClientByHost(gconv.String(ks))
	}
	return clusterConnMap[ks]
}

func (r *Redis) Cluster(key string) (interface{}, error) {
	return r.Do("cluster", key)
}

func (r *Redis) commandDo(action string, args ...interface{}) (interface{}, error) {

	if len(args) == 0 {
		conn := &Conn{r.pool.Get()}
		return conn.Do(action)
	}

	keys := gconv.String(args[0])
	if keys == "" {
		return nil, errors.New("key is empty")
	}
	ch := r
	if flagIsCluster {
		ch = choiceConn(keys)
	}

	conn := &Conn{ch.pool.Get()}
	defer conn.Close()

	*reply, err = ch.Do(action, args...)
	if err == nil {
		return *reply, nil
	}

	if err == nil {
		return reply, nil
	}

	if err != nil {
		if strings.Index(err.Error(), "MOVED") >= 0 {
			ch = r.movedConn(err.Error())
			return conn.Do(action, keys)
		}
		return nil, err
	}
	return nil, err
}

func (r *Redis) movedConn(errs string) *Redis {
	chs := strings.Split(errs, " ")
	r.layoutSlots()
	return choiceConn(chs[2])
}
