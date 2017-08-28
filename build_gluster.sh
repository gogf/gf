#!/bin/bash
# 编译
go build src/gapp/gluster/gluster.go
# 同步执行文件到测试节点
sshpass -p 123456 scp -P 22 ./gluster john@192.168.2.147:/home/john/gluster
sshpass -p 123456 scp -P 22 ./gluster john@192.168.2.186:/home/john/gluster
sshpass -p 123456 scp -P 22 ./gluster john@192.168.2.196:/home/john/gluster
# 同步配置文件到测试节点
sshpass -p 123456 scp -P 22 ./src/gapp/gluster/gluster_server.json john@192.168.2.147:/home/john/gluster.json
sshpass -p 123456 scp -P 22 ./src/gapp/gluster/gluster_server.json john@192.168.2.186:/home/john/gluster.json
sshpass -p 123456 scp -P 22 ./src/gapp/gluster/gluster_server.json john@192.168.2.196:/home/john/gluster.json
#删除本地执行文件
rm gluster