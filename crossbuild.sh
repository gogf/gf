#!/bin/bash
# 生成交叉编译的可执行文件cbuild
go build src/gapp/cbuild/cbuild.go
# 开始交叉编译gluster
./cbuild src/gapp/gluster/gluster.go --name=gluster --version=lastest
#./cbuild src/gapp/gluster/gluster.go --name=gluster --version=0.6 --arch=386,amd64 --os=linux,windows
# 删除临时可文件
rm cbuild