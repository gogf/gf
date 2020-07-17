package main

import (
	"fmt"

	"github.com/jin502437344/gf/container/gmap"
	"github.com/jin502437344/gf/frame/g"
	"github.com/jin502437344/gf/util/gutil"
)

func main() {
	array := g.Slice{2, 3, 1, 5, 4, 6, 8, 7, 9}
	hashMap := gmap.New()
	linkMap := gmap.NewListMap()
	treeMap := gmap.NewTreeMap(gutil.ComparatorInt)
	for _, v := range array {
		hashMap.Set(v, v)
	}
	for _, v := range array {
		linkMap.Set(v, v)
	}
	for _, v := range array {
		treeMap.Set(v, v)
	}
	fmt.Println("HashMap   Keys:", hashMap.Keys())
	fmt.Println("HashMap Values:", hashMap.Values())
	fmt.Println("LinkMap   Keys:", linkMap.Keys())
	fmt.Println("LinkMap Values:", linkMap.Values())
	fmt.Println("TreeMap   Keys:", treeMap.Keys())
	fmt.Println("TreeMap Values:", treeMap.Values())
}
